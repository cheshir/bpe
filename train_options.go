package bpe

var defaultTrainOptions = TrainOptions{
	MaxNumberOfTokens: 50000,
	MaxTokenLength:    5,
	ScanBufferSize:    64 * 1024,
}

type TrainOptions struct {
	MaxNumberOfTokens int
	MaxTokenLength    int
	ScanBufferSize    int
}

func (o *TrainOptions) Apply(opts ...TrainOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type TrainOption func(opts *TrainOptions)

func WithDefaultTrainOptions() TrainOption {
	return func(opts *TrainOptions) {
		*opts = defaultTrainOptions
	}
}

func WithMaxNumberOfTokensTrainOption(n int) TrainOption {
	return func(opts *TrainOptions) {
		opts.MaxNumberOfTokens = n
	}
}

func WithMaxTokenLengthTrainOption(length int) TrainOption {
	return func(opts *TrainOptions) {
		opts.MaxTokenLength = length
	}
}

func WithScanBufferSizeTrainOption(size int) TrainOption {
	return func(opts *TrainOptions) {
		opts.ScanBufferSize = size
	}
}
