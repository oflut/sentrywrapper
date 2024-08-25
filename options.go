package sentrywrapper

type Wrapper struct {
	dsn              string
	environment      string
	release          string
	debug            bool
	sampleRate       float64
	maxBreadcrumbs   int
	attachStacktrace bool
}

type Option func(*Wrapper)

func WithEnvironment(env string) Option {
	return func(w *Wrapper) {
		w.environment = env
	}
}

func WithRelease(release string) Option {
	return func(w *Wrapper) {
		w.release = release
	}
}

func WithDebug(debug bool) Option {
	return func(w *Wrapper) {
		w.debug = debug
	}
}

func WithSampleRate(rate float64) Option {
	return func(w *Wrapper) {
		w.sampleRate = rate
	}
}

func WithMaxBreadcrumbs(max int) Option {
	return func(w *Wrapper) {
		w.maxBreadcrumbs = max
	}
}

func WithAttachStacktrace(attach bool) Option {
	return func(w *Wrapper) {
		w.attachStacktrace = attach
	}
}
