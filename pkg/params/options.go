package params

type Option func(params Params)

func WithParams(params bool) Option {
	return func(params Params) {}
}
