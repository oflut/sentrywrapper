package sentrywrapper

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

// SentryWrapper is a wrapper around the Sentry client to simplify error reporting and context management.
type SentryWrapper struct {
	client *sentry.Client
}

// New initializes and returns a SentryWrapper with the provided DSN and options.
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

// Get retrieves the underlying Sentry client.
func (sw *SentryWrapper) Get() *sentry.Client {
	if sw == nil {
		return nil
	}
	return sw.client
}

// SetUser assigns the current user to the global Sentry scope.
func (sw *SentryWrapper) SetUser(user sentry.User) {
	if sw == nil || sw.client == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(user)
	})
}

// SetUserWithContext assigns the current user to the Sentry scope within the provided context.
func (sw *SentryWrapper) SetUserWithContext(ctx context.Context, user sentry.User) {
	if sw == nil || sw.client == nil {
		return
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(user)
	})
}

// SetTag sets a key-value pair tag on the global Sentry scope.
func (sw *SentryWrapper) SetTag(key, value string) {
	if sw == nil || sw.client == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

// SetTagWithContext sets a key-value pair tag on the Sentry scope within the provided context.
func (sw *SentryWrapper) SetTagWithContext(ctx context.Context, key, value string) {
	if sw == nil || sw.client == nil {
		return
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

// SetTags sets multiple key-value pair tags on the global Sentry scope.
func (sw *SentryWrapper) SetTags(tags map[string]string) {
	if sw == nil || sw.client == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
}

// SetTagsWithContext sets multiple key-value pair tags on the Sentry scope within the provided context.
func (sw *SentryWrapper) SetTagsWithContext(ctx context.Context, tags map[string]string) {
	if sw == nil || sw.client == nil {
		return
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
}

// CaptureException reports an error to Sentry.
func (sw *SentryWrapper) CaptureException(err error) *sentry.EventID {
	if sw == nil || sw.client == nil || err == nil {
		return nil
	}
	return sw.client.CaptureException(err, nil, nil)
}

// CaptureExceptionWithContext reports an error to Sentry using the Sentry hub within the provided context.
func (sw *SentryWrapper) CaptureExceptionWithContext(ctx context.Context, err error) *sentry.EventID {
	if sw == nil || sw.client == nil || err == nil {
		return nil
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	return hub.CaptureException(err)
}

// CaptureMessage sends a message to Sentry.
func (sw *SentryWrapper) CaptureMessage(message string) *sentry.EventID {
	if sw == nil || sw.client == nil || message == "" {
		return nil
	}
	return sw.client.CaptureMessage(message, nil, nil)
}

// CaptureMessageWithContext sends a message to Sentry using the Sentry hub within the provided context.
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

// AddBreadcrumb adds a breadcrumb to the Sentry scope within the provided context.
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

// Flush waits for all buffered events to be sent to Sentry, blocking up to the given timeout duration.
func (sw *SentryWrapper) Flush(timeout time.Duration) bool {
	if sw == nil || sw.client == nil {
		return false
	}
	return sw.client.Flush(timeout)
}

// WithContext returns a new context with a Sentry hub bound to it, ensuring the context is tracked.
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

// Recover captures and logs a panic, ensuring it is reported to Sentry before re-panicking.
func (sw *SentryWrapper) Recover(recoveredError interface{}) {
	if recoveredError == nil || sw == nil || sw.client == nil {
		return
	}

	if eventID := sw.client.Recover(recoveredError, nil, nil); eventID != nil {
		log.Printf("Captured panic (ID: %s): %v", *eventID, recoveredError)
	}
}
