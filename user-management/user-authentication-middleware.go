package user_management

import "net/http"

func UserAuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-Name")
		password := r.Header.Get("X-User-Password")

		// Check if credentials were supplied
		if user == "" || password == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("No credentials supplied!"))
			return
		}

		// TODO: Check if user is in database

		handler.ServeHTTP(w, r)
	})
}
