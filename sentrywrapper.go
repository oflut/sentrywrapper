package sentrywrapper

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryWrapper struct {
	client *sentry.Client
}

// New returns a wrapper type with given dsn and options
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

func (sw *SentryWrapper) Get() *sentry.Client {
	return sw.client
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

func (sw *SentryWrapper) CaptureException(err error) *sentry.EventID {
	return sw.client.CaptureException(err, nil, nil)
}

func (sw *SentryWrapper) CaptureMessage(message string) *sentry.EventID {
	return sw.client.CaptureMessage(message, nil, nil)
}

func (sw *SentryWrapper) AddBreadcrumb(ctx context.Context, breadcrumb *sentry.Breadcrumb) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.AddBreadcrumb(breadcrumb, nil)
}

func (sw *SentryWrapper) Flush(timeout time.Duration) bool {
	return sw.client.Flush(timeout)
}

func (sw *SentryWrapper) WithContext(ctx context.Context) context.Context {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}
	hub.BindClient(sw.client)
	return ctx
}

func (sw *SentryWrapper) Recover() {
	if err := recover(); err != nil {
		eventID := sw.client.Recover(err, nil, nil)
		if eventID != nil {
			log.Printf("Captured panic (ID: %s): %v", *eventID, err)
		}
		panic(err)
	}
}
