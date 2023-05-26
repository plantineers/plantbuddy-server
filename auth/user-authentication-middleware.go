// Author: Maximilian Floto
package auth

import (
	"database/sql"
	"net/http"

	"github.com/plantineers/plantbuddy-server/model"
	"github.com/plantineers/plantbuddy-server/utils"
)

// Takes as parameters the function serving the endpoint, the minimum role, an array of functions that are not subject to authentication
func UserAuthMiddleware(f func(http.ResponseWriter, *http.Request), role model.Role, unrestrictedMethods []string) http.Handler {
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

		userName := r.Header.Get("X-User-Name")
		password := r.Header.Get("X-User-Password")

		// Check if credentials were supplied
		if userName == "" || password == "" {
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("No credentials supplied!"))
			if err != nil {
				return
			}
			return
		}

		// Get user from db
		user, err := getUserByName(userName)
		if err == sql.ErrNoRows { // User not found
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("User not found"))
			return
		} else if err != nil { // Database error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Check password
		password = utils.HashPassword(password)
		if password != user.Password {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid password"))
			return
		}

		// Check role
		if user.Role > role {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Insufficient permissions"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
