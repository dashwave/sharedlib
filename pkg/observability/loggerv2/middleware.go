package loggerv2

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		L.Info().Msg("LoggerMiddleWare")
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		ctx := r.Context()
		trace := trace.SpanFromContext(ctx)
		defer func() {

			panicVal := recover()
			if panicVal != nil {
				lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
				panic(panicVal)                                 // continue panicking
			}
			L.
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Str("span_id", trace.SpanContext().SpanID().String()).
				Str("trace_id", trace.SpanContext().TraceID().String()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("incoming request")
		}()

		next.ServeHTTP(w, r)
	})
}
