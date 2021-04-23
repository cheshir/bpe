package bpe

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

func ExampleTrain() {
	source := strings.NewReader("Lorem Ipsum")
	m, err := Train(source, WithDefaultTrainOptions())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%d", len(m.vocab))
	// Output: 29
}

func TestTrain(t *testing.T) {
	tt := []struct {
		name      string
		input     string
		option    TrainOption
		expected  map[string]struct{}
		withError bool
	}{
		{
			name:   "word",
			input:  "apple",
			option: WithMaxTokenLengthTrainOption(3),
			expected: map[string]struct{}{
				"a":   {},
				"p":   {},
				"l":   {},
				"e":   {},
				"ap":  {},
				"app": {},
				"pp":  {},
				"ppl": {},
				"pl":  {},
				"ple": {},
				"le":  {},
			},
		},
		{
			name:  "words",
			input: "foo bar",
			expected: map[string]struct{}{
				"f":   {},
				"fo":  {},
				"foo": {},
				"o":   {},
				"oo":  {},
				"b":   {},
				"ba":  {},
				"bar": {},
				"a":   {},
				"ar":  {},
				"r":   {},
			},
		},
		{
			name:  "not word",
			input: "[a]=1",
			expected: map[string]struct{}{
				"[": {},
				"a": {},
				"]": {},
				"=": {},
				"1": {},
			},
		},
		{
			name:   "max token length",
			input:  "aaaaaaaaa",
			option: WithMaxTokenLengthTrainOption(5),
			expected: map[string]struct{}{
				"a":     {},
				"aa":    {},
				"aaa":   {},
				"aaaa":  {},
				"aaaaa": {},
			},
		},
		{
			name:   "max tokens",
			input:  "aaaaaaaaa",
			option: WithMaxNumberOfTokensTrainOption(1),
			expected: map[string]struct{}{
				"a": {},
			},
		},
		{
			name:     "empty",
			input:    "",
			expected: map[string]struct{}{},
		},
		{
			name:      "error",
			input:     "asdasdasd",
			option:    WithScanBufferSizeTrainOption(1),
			withError: true, // bufio.Scanner: token too long
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			options := []TrainOption{WithDefaultTrainOptions()}
			if tc.option != nil {
				options = append(options, tc.option)
			}

			m, err := Train(strings.NewReader(tc.input), options...)
			if err != nil && !tc.withError {
				t.Errorf("Unexpected error: %v", err)
			}

			if m == nil {
				if !tc.withError {
					t.Error("Model should not be nil")
				}

				return
			}

			if !reflect.DeepEqual(tc.expected, m.vocab) {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, m.vocab)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tt := []struct {
		word         string
		maxTokenSize int
		expected     tokensFrequencyTable
	}{
		{
			word:         "12314",
			maxTokenSize: 3,
			expected: tokensFrequencyTable{
				"1": 2,
				"2": 1,
				"3": 1,
				"4": 1,
			},
		},
		{
			word:         "abcad",
			maxTokenSize: 3,
			expected: tokensFrequencyTable{
				"a":   2,
				"b":   1,
				"c":   1,
				"d":   1,
				"ab":  1,
				"ad":  1,
				"abc": 1,
				"bc":  1,
				"bca": 1,
				"ca":  1,
				"cad": 1,
			},
		},
		{
			word:         "a-b",
			maxTokenSize: 3,
			expected: tokensFrequencyTable{
				"a":   1,
				"a-":  1,
				"a-b": 1,
				"-":   1,
				"-b":  1,
				"b":   1,
			},
		},
		{
			word:         "[xxx]",
			maxTokenSize: 3,
			expected: tokensFrequencyTable{
				"[":   1,
				"]":   1,
				"x":   3,
				"xx":  2,
				"xxx": 1,
			},
		},
		{
			word:         "(foo)",
			maxTokenSize: 1,
			expected: tokensFrequencyTable{
				"(": 1,
				"f": 1,
				"o": 2,
				")": 1,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.word, func(t *testing.T) {
			actualTokens := make(tokensFrequencyTable, 0)
			tokenize(actualTokens, tc.word, tc.maxTokenSize)

			if !reflect.DeepEqual(tc.expected, actualTokens) {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, actualTokens)
			}
		})
	}
}
