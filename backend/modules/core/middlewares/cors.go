// Package middlewares contains middleware functions for the web server.
//
// The functions are meant to be used with the net/http package.
//
// The functions are:
//
//   - AllowCors: adds CORS headers to the response.
package middlewares

import "net/http"

// AllowCors adds CORS headers to the response.
func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
