package tracer

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func GetTracerProvider() (*sdktrace.TracerProvider, error) {
	insecureMode := true

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecureMode {
		secureOption = otlptracegrpc.WithInsecure()
	}

	otelExporterEndoint := os.Getenv("OTEL_EXPORTER_ENDPOINT")
	if otelExporterEndoint == "" {
		return nil, errors.New("OTEL_EXPORTER_ENDPOINT is not set")
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		return nil, errors.New("SERVICE_NAME is not set")
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(otelExporterEndoint),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create otel exporter: %w", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel resources: %w", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return traceProvider, nil
}

func GetTracer() (trace.Tracer, error) {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		return nil, errors.New("SERVICE_NAME is not set")
	}
	return otel.Tracer(serviceName), nil
}
