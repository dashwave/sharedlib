package loggerv2

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

type loggingResponseWriterGin struct {
	gin.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func newLoggingResponseWriterGin(w gin.ResponseWriter) *loggingResponseWriterGin {
	return &loggingResponseWriterGin{w, http.StatusOK}
}

func (lrw *loggingResponseWriterGin) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func SetContextMiddlewareMux(ctx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqContext := r.Context()
			newContext := ZCtx(ctx).Logger().WithContext(reqContext)
			r = r.WithContext(newContext)
			next.ServeHTTP(w, r)
		})
	}
}

func LoggerMiddleWareMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		ctx := r.Context()
		trace := trace.SpanFromContext(ctx)
		ctx = ZCtx(ctx).Str("span_id", trace.SpanContext().SpanID().String()).
			Str("trace_id", trace.SpanContext().TraceID().String()).
			Logger().WithContext(ctx)
		r = r.WithContext(ctx)
		defer func() {

			panicVal := recover()
			if panicVal != nil {
				lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
				panic(panicVal)                                 // continue panicking
			}
			Ctx(ctx).
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Int("status_code", lrw.statusCode).
				Dur("elapsed_ms", time.Since(start)).
				Msg("incoming request")
		}()

		next.ServeHTTP(lrw, r)
	})
}

func LoggerMiddleWareGin(c *gin.Context) {
	start := time.Now()
	lrw := newLoggingResponseWriterGin(c.Writer)
	ctx := c.Request.Context()
	trace := trace.SpanFromContext(ctx)
	logger := Get().With().
		Str("span_id", trace.SpanContext().SpanID().String()).
		Str("trace_id", trace.SpanContext().TraceID().String()).
		Logger()
	c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context()))
	defer func() {
		panicVal := recover()
		if panicVal != nil {
			lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
			panic(panicVal)                                 // continue panicking
		}
		logger.
			Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.RequestURI()).
			Str("user_agent", c.Request.UserAgent()).
			Int("status_code", lrw.statusCode).
			Dur("elapsed_ms", time.Since(start)).
			Msg("incoming request")
	}()
	c.Writer = lrw
	c.Next()
}
