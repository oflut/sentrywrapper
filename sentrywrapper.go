package sentrywrapper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryWrapper struct {
	client *sentry.Client
}

// New returns a wrapper type with given dns and options
func New(dsn string, opts ...Option) (*SentryWrapper, error) {
	clientOptions := sentry.ClientOptions{
		Dsn: dsn,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(&clientOptions)
	}

	client, err := sentry.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	return &SentryWrapper{client: client}, nil
}

func (s *SentryWrapper) Get() *sentry.Client {
	return s.client
}

func (sw *SentryWrapper) SetUser(user sentry.User) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(user)
	})
}

func (sw *SentryWrapper) SetTag(key, value string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

func (sw *SentryWrapper) SetTags(tags map[string]string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
}

func (sw *SentryWrapper) CaptureException(err error) {
	sentry.CaptureException(err)
}

func (sw *SentryWrapper) CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}

func (sw *SentryWrapper) AddBreadcrumb(breadcrumb *sentry.Breadcrumb) {
	sentry.AddBreadcrumb(breadcrumb)
}

func (sw *SentryWrapper) Flush(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}

func (sw *SentryWrapper) WithContext(ctx context.Context) context.Context {
	return sentry.SetHubOnContext(ctx, sentry.CurrentHub().Clone())
}

func (sw *SentryWrapper) Recover() {
	if err := recover(); err != nil {
		if e, ok := err.(error); ok {
			sw.CaptureException(e)
		} else {
			sw.CaptureMessage(fmt.Sprintf("%v", err))
		}
		log.Printf("Recovered from panic: %v", err)
	}
}
