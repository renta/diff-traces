package main_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	trc "go.opentelemetry.io/otel/trace"
)

func TestTracingInjectionWithProperSpan(t *testing.T) {
	ctx := context.Background()
	//use os.Stdout for newExporter to examine a full traceSdk
	exp, err := newExporter(io.Discard)
	if err != nil {
		t.Fatalf("can not create tracing exporter with error '%s'", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			t.Fatalf(err.Error())
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("test")

	ctx, span := tracer.Start(ctx, "test")
	defer span.End()

	assert.True(t, span.IsRecording(), "here is a recordingSpan")
}

func TestTracingInjectionWithNonRecordingSpan(t *testing.T) {
	ctx := context.Background()
	//use os.Stdout for newExporter to examine a full traceSdk
	exp, err := newExporter(io.Discard)
	if err != nil {
		t.Fatalf("can not create tracing exporter with error '%s'", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			t.Fatalf(err.Error())
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("test")

	// the simplest way to bypass a tracing info
	traceID, err := trc.TraceIDFromHex("965b5a5c2dc31e184315335ad7c09b77")
	if err != nil {
		t.Fatalf(err.Error())
	}
	spanID, err := trc.SpanIDFromHex("c383b44aa60317c2")
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx, span := tracer.Start(trc.ContextWithSpanContext(
		ctx, trc.SpanContext{}.WithTraceID(traceID).WithSpanID(spanID)),
		"test",
	)
	defer span.End()

	assert.False(t, span.IsRecording(), "here is a nonRecordingSpan")
}

func TestTracingInjectionWithAlwaysSample(t *testing.T) {
	ctx := context.Background()
	//use os.Stdout for newExporter to examine a full traceSdk
	exp, err := newExporter(io.Discard)
	if err != nil {
		t.Fatalf("can not create tracing exporter with error '%s'", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.AlwaysSample()), // <- adding sampler creates a recordingSpan here
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			t.Fatalf(err.Error())
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("test")

	// the simplest way to bypass a tracing info
	traceID, err := trc.TraceIDFromHex("965b5a5c2dc31e184315335ad7c09b77")
	if err != nil {
		t.Fatalf(err.Error())
	}
	spanID, err := trc.SpanIDFromHex("c383b44aa60317c2")
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx, span := tracer.Start(trc.ContextWithSpanContext(
		ctx, trc.SpanContext{}.WithTraceID(traceID).WithSpanID(spanID)),
		"test",
	)
	defer span.End()

	assert.True(t, span.IsRecording(), "here is a ecordingSpan")
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
