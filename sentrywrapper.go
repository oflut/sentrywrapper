package sentrywrapper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

// New returns a wrapper type with given dns and options
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

func (w *Wrapper) Initialize() (*sentry.Client, error) {
	return sentry.NewClient(sentry.ClientOptions{
		Dsn:              w.dsn,
		Environment:      w.environment,
		Release:          w.release,
		Debug:            w.debug,
		SampleRate:       w.sampleRate,
		MaxBreadcrumbs:   w.maxBreadcrumbs,
		AttachStacktrace: w.attachStacktrace,
	})
}

func (w *Wrapper) SetUser(user sentry.User) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(user)
	})
}

func (w *Wrapper) SetTag(key, value string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

func (w *Wrapper) SetTags(tags map[string]string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
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
		if e, ok := err.(error); ok {
			w.CaptureException(e)
		} else {
			w.CaptureMessage(fmt.Sprintf("%v", err))
		}
		log.Printf("Recovered from panic: %v", err)
	}
}
