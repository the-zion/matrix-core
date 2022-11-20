package trace

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/metadata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"strings"
)

const serviceHeader = "x-md-service-name"

// Metadata is tracing metadata propagator
type Metadata struct{}

var _ propagation.TextMapPropagator = Metadata{}

// Inject sets metadata key-values from ctx into the carrier.
func (b Metadata) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	app, ok := kratos.FromContext(ctx)
	if ok {
		var builder strings.Builder
		for _, str := range []string{app.Name(), ".", app.ID()} {
			builder.WriteString(str)
		}
		carrier.Set(serviceHeader, builder.String())
	}
}

// Extract returns a copy of parent with the metadata from the carrier added.
func (b Metadata) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	name := carrier.Get(serviceHeader)
	if name == "" {
		return parent
	}
	if md, ok := metadata.FromServerContext(parent); ok {
		md.Set(serviceHeader, name)
		return parent
	}
	md := metadata.New()
	md.Set(serviceHeader, name)
	parent = metadata.NewServerContext(parent, md)
	return parent
}

// Fields returns the keys who's values are set with Inject.
func (b Metadata) Fields() []string {
	return []string{serviceHeader}
}

func SetTracerProvider(url, token, service, hostname string) (*tracesdk.TracerProvider, error) {
	var builder strings.Builder
	for _, str := range []string{service, ".", hostname} {
		builder.WriteString(str)
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(0.5))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(builder.String()),
			attribute.String("exporter", "jaeger"),
			attribute.Float64("float", 312.23),
			attribute.KeyValue{
				Key: "token", Value: attribute.StringValue(token),
			},
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
