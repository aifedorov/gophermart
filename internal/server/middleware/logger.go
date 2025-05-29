package middleware

import (
	"github.com/aifedorov/gophermart/internal/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
		body   []byte
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	r.responseData.body = append(r.responseData.body, b...)
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Log.Info("HTTP request ==>",
			zap.String("method", r.Method),
			zap.String("URL", r.URL.String()),
			zap.Any("headers", r.Header),
			zap.Duration("duration", duration),
		)
	})
}

func ResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd := &responseData{}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   rd,
		}

		start := time.Now()
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)

		logger.Log.Info("HTTP response <==",
			zap.Int("status", rd.status),
			zap.Any("headers", r.Header),
			zap.ByteString("body", rd.body),
			zap.Duration("duration", duration),
		)
	})
}
