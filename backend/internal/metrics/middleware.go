package metrics

import (
	"net/http"
	"strconv"
	"time"
)

// responseWriter оборачивает http.ResponseWriter для перехвата статуса
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// PrometheusMiddleware возвращает middleware, который записывает метрики запросов
func PrometheusMiddleware(handlerName string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// Нормализуем статус: 0 → 200 (если WriteHeader не вызывался)
		status := wrapped.statusCode
		if status == 0 {
			status = http.StatusOK
		}

		// Записываем метрики
		HttpRequestsTotal.WithLabelValues(r.Method, handlerName, strconv.Itoa(status)).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, handlerName, strconv.Itoa(status)).
			Observe(time.Since(start).Seconds())
	}
}
