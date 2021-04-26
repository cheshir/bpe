package bpe

import (
	"encoding/json"
	"io"
)

func Export(model *BPE, w io.Writer, opts ...ExportOption) error {
	options := defaultExportOptions()
	options.Apply(opts...)

	m := exportedModel{
		MaxTokenLength: model.maxTokenLength,
		Vocab:          make([]string, 0, len(model.vocab)),
	}

	for t := range model.vocab {
		m.Vocab = append(m.Vocab, t)
	}

	return options.Encoder.Encode(w, m)
}

func defaultExportOptions() *exportOptions {
	return &exportOptions{
		Encoder: &defaultEncoder{},
	}
}

type ModelEncoder interface {
	Encode(w io.Writer, model interface{}) error
}

type exportOptions struct {
	Encoder ModelEncoder
}

func (o *exportOptions) Apply(opts ...ExportOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type ExportOption func(opts *exportOptions)

func WithEncoder(enc ModelEncoder) ExportOption {
	return func(opts *exportOptions) {
		opts.Encoder = enc
	}
}

type exportedModel struct {
	MaxTokenLength int      `json:"max_token_length"`
	Vocab          []string `json:"vocab"`
}

type defaultEncoder struct{}

func (e *defaultEncoder) Encode(w io.Writer, model interface{}) error {
	return json.NewEncoder(w).Encode(model)
}
