package tracer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dashwave/sharedlib/pkg/observability/loggerv2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type CustomTracer struct {
	Tracer trace.Tracer
}

type Span trace.Span

func Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	ctx, span := t.Tracer.Start(ctx, spanName)
	spanID := span.SpanContext().SpanID().String()
	traceID := span.SpanContext().TraceID().String()
	ctx = loggerv2.ZCtx(ctx).Str("span_id", spanID).Str("trace_id", traceID).Logger().WithContext(ctx)
	return ctx, span
}

func SetSpanInContext(span trace.Span, ctx context.Context) context.Context {
	return trace.ContextWithSpan(ctx, span)
}

func GetSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

var t CustomTracer
var tp *sdktrace.TracerProvider

func GetMuxMiddleware() mux.MiddlewareFunc {
	return otelmux.Middleware(os.Getenv("SERVICE_NAME"))
}

func GetGinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware(os.Getenv("SERVICE_NAME"))
}

func Shutdown(ctx context.Context) error {
	return tp.Shutdown(ctx)
}

func InitTracer() error {
	insecureMode := true

	if os.Getenv("OTEL_EXPORTER_ENDPOINT") == "" {
		return fmt.Errorf("OTEL_EXPORTER_ENDPOINT is not set")
	}

	if os.Getenv("SERVICE_NAME") == "" {
		return fmt.Errorf("SERVICE_NAME is not set")
	}

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecureMode {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_ENDPOINT")),
		),
	)

	if err != nil {
		return err
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", os.Getenv("SERVICE_NAME")),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tp)

	t.Tracer = otel.Tracer(os.Getenv("SERVICE_NAME"))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

func SetSpanID(l zerolog.Logger, s trace.Span) zerolog.Logger {
	return l.With().Str("span_id", s.SpanContext().SpanID().String()).Logger()
}

func SetTraceID(l zerolog.Logger, s trace.Span) zerolog.Logger {
	return l.With().Str("trace_id", s.SpanContext().TraceID().String()).Logger()
}

func GetTransport() *otelhttp.Transport {
	return otelhttp.NewTransport(http.DefaultTransport)
}

func GetSpanIDFromCtx(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().SpanID().String()
}

func GetTraceIDFromCtx(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
