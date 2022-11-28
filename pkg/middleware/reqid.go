package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

// requestIDKey points to the value in the context where the request id  is stored.
const requestIDKey = "X-Request-ID"

func WithXRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := RequestID(r)
		// Set the id to ensure that the same request id is in the response
		w.Header().Set(requestIDKey, id)
		next.ServeHTTP(w, r)
	})
}

func RequestID(r *http.Request) string {
	requestID := r.Header.Get(requestIDKey)
	if requestID == "" {
		requestID = uuid.New().String()
		r.Header.Add(requestIDKey, requestID)
	}
	return requestID
}
