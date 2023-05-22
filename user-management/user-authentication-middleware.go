package user_management

import (
	"github.com/plantineers/plantbuddy-server/model"
	"net/http"
)

func UserAuthMiddleware(f func(http.ResponseWriter, *http.Request), role model.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-Name")
		password := r.Header.Get("X-User-Password")

		// Check if credentials were supplied
		if user == "" || password == "" {
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
