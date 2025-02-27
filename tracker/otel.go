package tracker

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

func Init() func() {
	ctx := context.Background()

	// Configure OTLP Exporter for Sentry
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("https://o299459.ingest.us.sentry.io/api/4508888280989696/security/?sentry_key=dd73ffa13055e661a88d441c0b73b64f"), // Replace with your Sentry org number
		// otlptracegrpc.WithHeaders(map[string]string{
		// 	"X-Sentry-Auth": "Sentry sentry_key=<your-public-key>, sentry_version=7",
		// }),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")), // Use TLS for secure connection
	)
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v", err)
	}

	// Set up Tracer Provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("easy-campaign")
	ctx, span := tracer.Start(ctx, "Init")
	span.AddEvent("Starting the application...")
	time.Sleep(1 * time.Second)

	span.End()

	err = sentry.Init(sentry.ClientOptions{
		Dsn: "https://dd73ffa13055e661a88d441c0b73b64f@o299459.ingest.us.sentry.io/4508888280989696",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

	// Shutdown function
	return func() {
		_ = tp.Shutdown(ctx)
	}
}
