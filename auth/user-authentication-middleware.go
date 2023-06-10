// Author: Maximilian Floto
package auth

import (
	"fmt"
	"net/http"

	"github.com/plantineers/plantbuddy-server/utils"
)

// Takes as parameters the function serving the endpoint, the minimum role
func UserAuthMiddleware(f func(http.ResponseWriter, *http.Request), role Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := http.HandlerFunc(f)

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
