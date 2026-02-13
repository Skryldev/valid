package validx

type Option func(*config)

type config struct {
	failFast bool
}

func WithFailFast() Option {
	return func(c *config) {
		c.failFast = true
	}
}
