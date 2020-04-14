package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func RecoverFromPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recoverPanic := func() {
			if r := recover(); r != nil {
				log.Printf("error: %v, stack: %v\n", r, string(debug.Stack()))
				http.Error(w, "Something went wrong!", http.StatusInternalServerError)
				return
			}
		}

		defer recoverPanic()
		handler.ServeHTTP(w, r)
	})
}
