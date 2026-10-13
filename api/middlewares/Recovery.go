package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := "Caught panic: %v, Stack trace: %s"
				log.Printf(msg, err, string(debug.Stack()))

				statusCode := http.StatusInternalServerError
				http.Error(w, "Internal Server Error", statusCode)
			}
		}()

		next.ServeHTTP(w, r)
	})
}