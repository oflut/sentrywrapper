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

// GetClient retrieves the underlying Sentry client.
func (sw *SentryWrapper) GetClient() *sentry.Client {
	return sw.client
}

// SetUser assigns the current user to the Sentry scope within the provided context.
func (sw *SentryWrapper) SetUser(ctx context.Context, user sentry.User) {
	if sw.client == nil {
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

// SetTag sets a key-value pair tag on the Sentry scope within the provided context.
func (sw *SentryWrapper) SetTag(ctx context.Context, key, value string) {
	if sw.client == nil {
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

// SetTags sets multiple key-value pair tags on the Sentry scope within the provided context.
func (sw *SentryWrapper) SetTags(ctx context.Context, tags map[string]string) {
	if sw.client == nil {
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

// CaptureException reports an error to Sentry using the Sentry hub within the provided context.
func (sw *SentryWrapper) CaptureException(ctx context.Context, err error, tags map[string]string) *sentry.EventID {
	if sw.client == nil || err == nil {
		return nil
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	event := sentry.NewEvent()
	event.Exception = []sentry.Exception{{
		Value:      err.Error(),
		Type:       "error",
		Stacktrace: sentry.NewStacktrace(),
	}}
	event.Level = sentry.LevelError

	if tags != nil {
		event.Tags = tags
	}

	return hub.CaptureEvent(event)
}

// CaptureMessage sends a message to Sentry using the Sentry hub within the provided context.
func (sw *SentryWrapper) CaptureMessage(ctx context.Context, message string, tags map[string]string) *sentry.EventID {
	if sw.client == nil || message == "" {
		return nil
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	event := sentry.NewEvent()
	event.Message = message
	event.Level = sentry.LevelInfo // Default level for messages

	if tags != nil {
		event.Tags = tags
	}

	return hub.CaptureEvent(event)
}

// AddBreadcrumb adds a breadcrumb to the Sentry scope within the provided context.
func (sw *SentryWrapper) AddBreadcrumb(ctx context.Context, breadcrumb *sentry.Breadcrumb) {
	if sw.client == nil || breadcrumb == nil {
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
	if sw.client == nil {
		return false
	}
	return sw.client.Flush(timeout)
}

// WithContext returns a new context with a Sentry hub bound to it.
func (sw *SentryWrapper) WithContext(ctx context.Context) context.Context {
	if sw.client == nil {
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
func (sw *SentryWrapper) Recover(ctx context.Context, recoveredError interface{}, additionalTags map[string]string) {
	if recoveredError == nil || sw == nil || sw.client == nil {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)

	sw.SetTag(ctx, "timestamp", timestamp)
	sw.SetTags(ctx, additionalTags)

	if eventID := sw.client.Recover(recoveredError, nil, nil); eventID != nil {
		log.Printf("Captured panic (ID: %s): %v", *eventID, recoveredError)
	}

}
