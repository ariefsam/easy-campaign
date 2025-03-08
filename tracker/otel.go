package tracker

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func Init() {
	dsn := os.Getenv("SENTRY_DSN")
	environment := os.Getenv("ENVIRONMENT")
	if dsn == "" {
		log.Println("SENTRY_DSN is empty")
		return
	}
	if environment == "" {
		environment = "not-set"
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		Environment:      environment,
	})
	if err != nil {
		log.Println("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

}
