package middleware

import "net/http"

func UserAgent() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = SetUserAgentToContext(ctx, r.Header.Get("User-Agent"))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
