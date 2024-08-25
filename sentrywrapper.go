package sentrywrapper

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func New(dsn string, options ...Option) (*Wrapper, error) {
	if dsn == "" {
		return nil, errors.New("dsn must be provided")
	}
	w := &Wrapper{
		dsn:              dsn,
		environment:      "production",
		sampleRate:       1.0,
		maxBreadcrumbs:   100,
		attachStacktrace: true,
	}

	for _, option := range options {
		option(w)
	}

	return w, nil
}

func (w *Wrapper) Initialize() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              w.dsn,
		Environment:      w.environment,
		Release:          w.release,
		Debug:            w.debug,
		SampleRate:       w.sampleRate,
		MaxBreadcrumbs:   w.maxBreadcrumbs,
		AttachStacktrace: w.attachStacktrace,
	})
}

func (w *Wrapper) CaptureException(err error) {
	sentry.CaptureException(err)
}

func (w *Wrapper) CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}

func (w *Wrapper) AddBreadcrumb(breadcrumb *sentry.Breadcrumb) {
	sentry.AddBreadcrumb(breadcrumb)
}

func (w *Wrapper) Flush(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}

func (w *Wrapper) WithContext(ctx context.Context) context.Context {
	return sentry.SetHubOnContext(ctx, sentry.CurrentHub().Clone())
}

func (w *Wrapper) Recover() {
	if err := recover(); err != nil {
		w.CaptureException(err.(error))
		log.Printf("Recovered from panic: %v", err)
	}
}
