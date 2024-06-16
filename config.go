package funky

type Config struct {
	name    string
	onError func(error)
	logger  Logger
}

func (p *Config) Error(err error) {
	if p.onError != nil {
		if p.logger != nil {
			p.logger.Error("encountered error", "error", err, "name", p.name)
		}
		p.onError(err)
	}
}

type Opt func(p *Config)

func With(c *Config) Opt {
	return func(p *Config) {
		p.name = c.name
		p.onError = c.onError
		p.logger = c.logger
	}
}

func WithName(name string) Opt {
	return func(p *Config) {
		p.name = name
	}
}

func OnError(onError func(error)) Opt {
	return func(p *Config) {
		p.onError = onError
	}
}

func WithLogger(logger Logger) Opt {
	return func(p *Config) {
		p.logger = logger
	}
}

func buildConfig(opts ...Opt) *Config {
	p := &Config{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

type Logger interface {
	Error(msg string, args ...interface{})
}
