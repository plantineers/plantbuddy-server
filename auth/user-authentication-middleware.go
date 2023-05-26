// Author: Maximilian Floto
package auth

import (
	"net/http"
)

// Takes as parameters the function serving the endpoint, the minimum role, an array of functions that are not subject to authentication
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

		user, err := authUser(r)
		switch err {
		case ErrWrongCredentials:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Wrong credentials"))
		case ErrNoCredentials:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No credentials supplied"))
		case nil:
			// Check role
			if user.Role > role {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Insufficient permissions"))
				return
			}

			handler.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

	})
}
