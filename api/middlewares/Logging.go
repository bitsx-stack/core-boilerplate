package middlewares

import (
	"log"
	"net/http"
	"time"
)


func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		log.Printf("Start %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		next.ServeHTTP(w, r)
		
		log.Printf("Complete %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}