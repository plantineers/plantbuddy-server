// Author: Maximilian Floto
package auth

import (
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/utils"
)

// UserAuthMiddleware takes as parameters the function serving the endpoint, the minimum role, an array of methods that are not subject to authentication
func UserAuthMiddleware(f func(http.ResponseWriter, *http.Request), role Role, unrestrictedMethods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := http.HandlerFunc(f)

		// Check if method is unrestricted
		for _, unrestrictedMethod := range unrestrictedMethods {
			// Unrestricted methods skip the authentication process
			if r.Method == unrestrictedMethod {
				handler.ServeHTTP(w, r)
				return
			}
		}

		// Check if user is authorized
		user, err := authUser(r)
		switch err {
		case ErrWrongCredentials:
			utils.HttpForbiddenResponse(w, "Wrong credentials")
		case ErrNoCredentials:
			utils.HttpBadRequestResponse(w, "No credentials supplied")
		case nil:
			if user.Role > role {
				utils.HttpForbiddenResponse(w, "Insufficient permissions")
				return
			}

			handler.ServeHTTP(w, r)
		default:
			msg := fmt.Sprintf("Error authenticating user: %s", err.Error())
			utils.HttpInternalServerErrorResponse(w, msg)
		}
	})
}
