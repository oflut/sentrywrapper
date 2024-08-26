package sentrywrapper

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryWrapper struct {
	client *sentry.Client
}

// New returns a wrapper type with given dsn and options
func New(dsn string, opts ...Option) (*SentryWrapper, error) {
	if dsn == "" {
		return nil, errors.New("invalid DSN: cannot be empty")

	}

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
	if sw == nil || sw.client == nil {
		return nil
	}
	return sw.client
}

func (sw *SentryWrapper) SetUser(user sentry.User) {
	if sw == nil || sw.client == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(user)
	})
}

func (sw *SentryWrapper) SetTag(key, value string) {
	if sw == nil || sw.client == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

func (sw *SentryWrapper) SetTags(tags map[string]string) {
	if sw == nil || sw.client == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
}

func (sw *SentryWrapper) CaptureException(err error) *sentry.EventID {
	if sw == nil || sw.client == nil || err == nil {
		return nil
	}
	return sw.client.CaptureException(err, nil, nil)
}

func (sw *SentryWrapper) CaptureMessage(message string) *sentry.EventID {
	if sw == nil || sw.client == nil || message == "" {
		return nil
	}
	return sw.client.CaptureMessage(message, nil, nil)
}

func (sw *SentryWrapper) CaptureMessageWithContext(ctx context.Context, message string) *sentry.EventID {
	if sw == nil || sw.client == nil || message == "" {
		return nil
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	return hub.CaptureMessage(message)
}

func (sw *SentryWrapper) AddBreadcrumb(ctx context.Context, breadcrumb *sentry.Breadcrumb) {
	if sw == nil || sw.client == nil || breadcrumb == nil {
		return
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.AddBreadcrumb(breadcrumb, nil)
}

func (sw *SentryWrapper) Flush(timeout time.Duration) bool {
	if sw == nil || sw.client == nil {
		return false
	}
	return sw.client.Flush(timeout)
}

func (sw *SentryWrapper) WithContext(ctx context.Context) context.Context {
	if sw == nil || sw.client == nil {
		return ctx
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}
	hub.BindClient(sw.client)
	return ctx
}

func (sw *SentryWrapper) Recover(recoveredError interface{}) {
	if recoveredError != nil && sw != nil && sw.client != nil {
		eventID := sw.client.Recover(recoveredError, nil, nil)
		if eventID != nil {
			log.Printf("Captured panic (ID: %s): %v", *eventID, recoveredError)
		}
	}
}
