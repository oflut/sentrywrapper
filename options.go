package sentrywrapper

import (
	"github.com/getsentry/sentry-go"
)

type Option func(*sentry.ClientOptions)

func WithEnvironment(environment string) Option {
	return func(o *sentry.ClientOptions) {
		o.Environment = environment
	}
}

func WithRelease(release string) Option {
	return func(o *sentry.ClientOptions) {
		o.Release = release
	}
}

func WithSampleRate(rate float64) Option {
	return func(o *sentry.ClientOptions) {
		o.SampleRate = rate
	}
}

func WithDebug(debug bool) Option {
	return func(o *sentry.ClientOptions) {
		o.Debug = debug
	}
}

func WithTracesSampleRate(rate float64) Option {
	return func(o *sentry.ClientOptions) {
		o.TracesSampleRate = rate
	}
}

func WithMaxBreadcrumbs(max int) Option {
	return func(o *sentry.ClientOptions) {
		o.MaxBreadcrumbs = max
	}
}
