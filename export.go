package bpe

import (
	"encoding/json"
	"io"
)

var defaultExportOptions = ExportOptions{
	Encoder: &defaultEncoder{},
}

type ModelEncoder interface {
	Encode(w io.Writer, model interface{}) error
}

type ExportOptions struct {
	Encoder ModelEncoder
}

func (o *ExportOptions) Apply(opts ...ExportOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type ExportOption func(opts *ExportOptions)

func WithEncoder(enc ModelEncoder) ExportOption {
	return func(opts *ExportOptions) {
		opts.Encoder = enc
	}
}

type dto struct {
	MaxTokenLength int      `json:"max_token_length"`
	Vocab          []string `json:"vocab"`
}

func Save(model *BPE, w io.Writer, opts ...ExportOption) error {
	options := defaultExportOptions
	options.Apply(opts...)

	m := dto{
		MaxTokenLength: model.maxTokenLength,
		Vocab:          make([]string, 0, len(model.vocab)),
	}

	for t := range model.vocab {
		m.Vocab = append(m.Vocab, t)
	}

	return options.Encoder.Encode(w, m)
}

type defaultEncoder struct{}

func (e *defaultEncoder) Encode(w io.Writer, model interface{}) error {
	return json.NewEncoder(w).Encode(model)
}
