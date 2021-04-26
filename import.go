package bpe

import (
	"encoding/json"
	"io"
)

func Import(r io.Reader, opts ...ImportOption) (*BPE, error) {
	options := defaultImportOptions()
	options.Apply(opts...)

	return options.Decoder.Decode(r)
}

func defaultImportOptions() *importOptions {
	return &importOptions{
		Decoder: &defaultDecoder{},
	}
}

type ModelDecoder interface {
	Decode(r io.Reader) (*BPE, error)
}

type importOptions struct {
	Decoder ModelDecoder
}

func (o *importOptions) Apply(opts ...ImportOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type ImportOption func(opts *importOptions)

func WithDecoder(decoder ModelDecoder) ImportOption {
	return func(opts *importOptions) {
		opts.Decoder = decoder
	}
}

type defaultDecoder struct{}

func (e *defaultDecoder) Decode(r io.Reader) (*BPE, error) {
	dto := &exportedModel{}
	err := json.NewDecoder(r).Decode(dto)
	if err != nil {
		return nil, err
	}

	vocab := make(map[string]struct{}, len(dto.Vocab))

	for _, token := range dto.Vocab {
		vocab[token] = struct{}{}
	}

	model := &BPE{
		maxTokenLength: dto.MaxTokenLength,
		vocab:          vocab,
	}

	return model, nil
}
