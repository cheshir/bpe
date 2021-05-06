package bpe

import (
	"reflect"
	"strings"
	"testing"
)

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
			option: WithMaxTokenLength(3),
			expected: map[string]struct{}{
				BeginOfWord + "a":   {},
				"p":                 {},
				"l":                 {},
				"e" + EndOfWord:     {},
				BeginOfWord + "ap":  {},
				BeginOfWord + "app": {},
				"pp":                {},
				"ppl":               {},
				"pl":                {},
				"ple" + EndOfWord:   {},
				"le" + EndOfWord:    {},
			},
		},
		{
			name:  "words",
			input: "foo bar",
			expected: map[string]struct{}{
				BeginOfWord + "f":               {},
				BeginOfWord + "fo":              {},
				BeginOfWord + "foo" + EndOfWord: {},
				"o":                             {},
				"o" + EndOfWord:                 {},
				"oo" + EndOfWord:                {},
				BeginOfWord + "b":               {},
				BeginOfWord + "ba":              {},
				BeginOfWord + "bar" + EndOfWord: {},
				"a":                             {},
				"ar" + EndOfWord:                {},
				"r" + EndOfWord:                 {},
			},
		},
		{
			name:   "not word",
			input:  "[a]=1",
			option: WithMaxTokenLength(3),
			expected: map[string]struct{}{
				BeginOfWord + "[":   {},
				BeginOfWord + "[a":  {},
				BeginOfWord + "[a]": {},
				"a":                 {},
				"a]":                {},
				"a]=":               {},
				"]":                 {},
				"]=":                {},
				"]=1" + EndOfWord:   {},
				"=":                 {},
				"=1" + EndOfWord:    {},
				"1" + EndOfWord:     {},
			},
		},
		{
			name:     "not word with filter",
			input:    "[a]=1",
			option:   WithWordsOnly(),
			expected: map[string]struct{}{},
		},
		{
			name:   "max token length",
			input:  "aaaaaaaaa",
			option: WithMaxTokenLength(3),
			expected: map[string]struct{}{
				BeginOfWord + "a":   {},
				BeginOfWord + "aa":  {},
				BeginOfWord + "aaa": {},
				"a":                 {},
				"aa":                {},
				"aaa":               {},
				"aaa" + EndOfWord:   {},
				"aa" + EndOfWord:    {},
				"a" + EndOfWord:     {},
			},
		},
		{
			name:   "Mix words and not words",
			input:  "foo foo2",
			option: WithWordsOnly(),
			expected: map[string]struct{}{
				BeginOfWord + "f":               {},
				BeginOfWord + "fo":              {},
				BeginOfWord + "foo" + EndOfWord: {},
				"o":                             {},
				"o" + EndOfWord:                 {},
				"oo" + EndOfWord:                {},
			},
		},
		{
			name:   "max tokens",
			input:  "aaaaaaaaa",
			option: WithMaxNumberOfTokens(1),
			expected: map[string]struct{}{
				"a": {},
			},
		},
		{
			name:   "Leading space",
			input:  "  foo",
			option: WithWordsOnly(),
			expected: map[string]struct{}{
				BeginOfWord + "f":               {},
				BeginOfWord + "fo":              {},
				BeginOfWord + "foo" + EndOfWord: {},
				"o":                             {},
				"o" + EndOfWord:                 {},
				"oo" + EndOfWord:                {},
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
			option:    WithScanBufferSize(1),
			withError: true, // bufio.Scanner: token too long
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			options := make([]TrainOption, 0)
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
		wordsOnly    bool
		expected     tokensFrequencyTable
	}{
		{
			word:         "12114",
			maxTokenSize: 3,
			wordsOnly:    false,
			expected: tokensFrequencyTable{
				BeginOfWord + "1":   1,
				BeginOfWord + "12":  1,
				BeginOfWord + "121": 1,
				"2":                 1,
				"21":                1,
				"211":               1,
				"1":                 2,
				"11":                1,
				"114" + EndOfWord:   1,
				"14" + EndOfWord:    1,
				"4" + EndOfWord:     1,
			},
		},
		{
			word:         "5678",
			maxTokenSize: 3,
			wordsOnly:    true,
			expected:     tokensFrequencyTable{},
		},
		{
			word:         "keep",
			maxTokenSize: 2,
			expected: tokensFrequencyTable{
				BeginOfWord + "k":  1,
				BeginOfWord + "ke": 1,
				"e":                2,
				"ee":               1,
				"ep" + EndOfWord:   1,
				"p" + EndOfWord:    1,
			},
		},
		{
			word:         "words",
			maxTokenSize: 2,
			wordsOnly:    true,
			expected: tokensFrequencyTable{
				BeginOfWord + "w":  1,
				BeginOfWord + "wo": 1,
				"o":                1,
				"or":               1,
				"r":                1,
				"rd":               1,
				"d":                1,
				"ds" + EndOfWord:   1,
				"s" + EndOfWord:    1,
			},
		},
		{
			word:         "a-b",
			maxTokenSize: 3,
			expected: tokensFrequencyTable{
				BeginOfWord + "a":               1,
				BeginOfWord + "a-":              1,
				BeginOfWord + "a-b" + EndOfWord: 1,
				"-":                             1,
				"-b" + EndOfWord:                1,
				"b" + EndOfWord:                 1,
			},
		},
		{
			word:         "a-w",
			maxTokenSize: 3,
			wordsOnly:    true,
			expected: tokensFrequencyTable{
				BeginOfWord + "a":               1,
				BeginOfWord + "a-":              1,
				BeginOfWord + "a-w" + EndOfWord: 1,
				"-":                             1,
				"-w" + EndOfWord:                1,
				"w" + EndOfWord:                 1,
			},
		},
		{
			word:         "[xx]",
			maxTokenSize: 2,
			expected: tokensFrequencyTable{
				BeginOfWord + "[":  1,
				BeginOfWord + "[x": 1,
				"x":                2,
				"xx":               1,
				"x]" + EndOfWord:   1,
				"]" + EndOfWord:    1,
			},
		},
		{
			word:         "(foo)",
			maxTokenSize: 1,
			expected: tokensFrequencyTable{
				BeginOfWord + "(": 1,
				"f":               1,
				"o":               2,
				")" + EndOfWord:   1,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.word, func(t *testing.T) {
			actualTokens := make(tokensFrequencyTable)
			tokenize(actualTokens, tc.word, tc.maxTokenSize, tc.wordsOnly)

			if !reflect.DeepEqual(tc.expected, actualTokens) {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, actualTokens)
			}
		})
	}
}
