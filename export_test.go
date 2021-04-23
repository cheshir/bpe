package bpe

import (
	"bytes"
	"io"
	"testing"

	"github.com/pkg/errors"
)

type encoderMock struct {
	data []byte
	err  error
}

func (e *encoderMock) Encode(w io.Writer, _ interface{}) error {
	_, _ = w.Write(e.data)

	return e.err
}

func TestExport(t *testing.T) {
	tt := []struct {
		name        string
		model       *BPE
		encoderMock *encoderMock
		expected    string
		withError   bool
	}{
		{
			name: "test with default encoder",
			model: newModelFromTokensFrequencyTable(
				tokensFrequencyTable{
					"foo": 1,
				},
				1,
				3,
			),
			expected: `{"max_token_length":3,"vocab":["foo"]}` + "\n",
		},
		{
			name:  "mocked encoder",
			model: &BPE{},
			encoderMock: &encoderMock{
				data: []byte("123123"),
			},
			expected: "123123",
		},
		{
			name:  "encoder returns error",
			model: &BPE{},
			encoderMock: &encoderMock{
				err: errors.New("Some error"),
			},
			expected:  "",
			withError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var opts []ExportOption
			if tc.encoderMock != nil {
				opts = append(opts, WithEncoder(tc.encoderMock))
			}

			buf := bytes.NewBuffer(nil)

			if err := Save(tc.model, buf, opts...); err != nil && !tc.withError {
				t.Errorf("Unexpected Save error %v", err)
			}

			actual := buf.String()

			if tc.expected != actual {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, actual)
			}
		})
	}
}
