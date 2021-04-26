package bpe

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

type decoderMock struct {
	data *BPE
	err  error
}

func (d *decoderMock) Decode(_ io.Reader) (*BPE, error) {
	return d.data, d.err
}

func TestImport(t *testing.T) {
	tt := []struct {
		name        string
		source      io.Reader
		expected    *BPE
		decoderMock *decoderMock
		withError   bool
	}{
		{
			name:   "test with default decoder",
			source: strings.NewReader(`{"max_token_length":3,"vocab":["foo"]}`),
			expected: newModelFromTokensFrequencyTable(
				tokensFrequencyTable{
					"foo": 1,
				},
				1,
			),
		},
		{
			name:   "mocked decoder",
			source: strings.NewReader(""),
			decoderMock: &decoderMock{
				data: &BPE{
					maxTokenLength: 3,
					vocab: map[string]struct{}{
						"foo": {},
						"bar": {},
					},
				},
			},
			expected: &BPE{
				maxTokenLength: 3,
				vocab: map[string]struct{}{
					"foo": {},
					"bar": {},
				},
			},
			withError: false,
		},
		{
			name:     "encoder returns error",
			source:   strings.NewReader(""),
			expected: nil,
			decoderMock: &decoderMock{
				err: errors.New("Some error"),
			},
			withError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var opts []ImportOption
			if tc.decoderMock != nil {
				opts = append(opts, WithDecoder(tc.decoderMock))
			}

			actual, err := Import(tc.source, opts...)
			if err != nil {
				if !tc.withError {
					t.Errorf("Unexpected Export error %v", err)
				}

				return
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, actual)
			}
		})
	}
}
