package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitTracer initializes the OpenTelemetry tracer provider.
func InitTracer(serviceName string) (*trace.TracerProvider, error) {
	// Create stdout exporter to be able to retrieve shared spans for now
	// In production, you would use otlptracegrpc or similar
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(), // Make it readable in logs for now
	)
	if err != nil {
		return nil, err
	}

	// Create Resource to identify this service
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Register the Trace Provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	// Register global Trace Provider and Propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Println("Initialized OpenTelemetry Tracer")
	return tp, nil
}
