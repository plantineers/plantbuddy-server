package user_management

import (
	"github.com/plantineers/plantbuddy-server/model"
	"net/http"
)

// Takes as parameters the function serving the endpoint, the minimum role, an array of functions that are not subject to authentication
func UserAuthMiddleware(f func(http.ResponseWriter, *http.Request), role model.Role, unrestrictedMethods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isUnrestricted := false
		for _, unrestrictedMethod := range unrestrictedMethods {
			if r.Method == unrestrictedMethod {
				isUnrestricted = true
			}
		}

		user := r.Header.Get("X-User-Name")
		password := r.Header.Get("X-User-Password")

		// Check if credentials were supplied
		if (user == "" || password == "") && !isUnrestricted {
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("No credentials supplied!"))
			if err != nil {
				return
			}
			return
		}

		// TODO: Check if user is in database
		// TODO: Check password
		// TODO: Check role
		handler := http.HandlerFunc(f)
		handler.ServeHTTP(w, r)
	})
}
