package bpe

import (
	"io"
	"strings"
	"testing"
)

func TestBPE_Encode(t *testing.T) {
	tt := []struct {
		name      string
		b         *BPE
		in        io.Reader
		expected  []string
		withError bool
	}{
		{
			name: "Foo. Bar\n",
			in:   strings.NewReader("Foo foo. Bar\n"),
			expected: []string{
				BeginOfSentence,
				BeginOfWord + "Fo",
				"o" + EndOfWord,
				BeginOfWord + "fo",
				"o." + EndOfWord,
				EndOfSentence,

				BeginOfSentence,
				BeginOfWord + "Bar" + EndOfWord,
				EndOfSentence,
			},
			b: &BPE{
				maxTokenLength: 16,
				vocab: map[string]struct{}{
					BeginOfWord + "F":               {},
					BeginOfWord + "Fo":              {},
					BeginOfWord + "fo":              {},
					"o":                             {},
					"o" + EndOfWord:                 {},
					"o." + EndOfWord:                {},
					"." + EndOfWord:                 {},
					BeginOfWord + "B":               {},
					BeginOfWord + "Ba":              {},
					BeginOfWord + "Bar" + EndOfWord: {},
				},
			},
		},
		{
			name:     "empty",
			in:       strings.NewReader(""),
			b:        &BPE{},
			expected: []string{},
		},
		{
			name: "foo",
			in:   strings.NewReader("foo"),
			expected: []string{
				BeginOfSentence,
				BeginOfWord + "foo" + EndOfWord,
				EndOfSentence,
			},
			b: &BPE{
				maxTokenLength: 64,
				vocab: map[string]struct{}{
					BeginOfWord + "f":               {},
					"o":                             {},
					"o" + EndOfWord:                 {},
					BeginOfWord + "fo":              {},
					"oo" + EndOfWord:                {},
					BeginOfWord + "foo" + EndOfWord: {},
				},
			},
		},
		{
			name: "foo",
			in:   strings.NewReader("foo"),
			expected: []string{
				BeginOfSentence,
				BeginOfWord,
				"fo",
				"o",
				EndOfWord,
				EndOfSentence,
			},
			b: &BPE{
				maxTokenLength: 4,
				vocab: map[string]struct{}{
					BeginOfWord:                     {},
					EndOfWord:                       {},
					"fo":                            {},
					"o":                             {},
					"o" + EndOfWord:                 {},
					BeginOfWord + "fo":              {},
					"oo" + EndOfWord:                {},
					BeginOfWord + "foo" + EndOfWord: {},
				},
			},
		},
		{
			name: "empty vocab",
			in:   strings.NewReader("foo"),
			// <w>foo</w>
			expected: []string{
				BeginOfSentence,
				UnknownToken, UnknownToken, UnknownToken, UnknownToken, UnknownToken,
				UnknownToken, UnknownToken, UnknownToken, UnknownToken, UnknownToken,
				EndOfSentence,
			},
			b: &BPE{
				maxTokenLength: 64,
				vocab:          map[string]struct{}{},
			},
		},
		{
			name: "тест.",
			in:   strings.NewReader("тест."),
			expected: []string{
				BeginOfSentence,
				BeginOfWord + "тес",
				"т",
				"." + EndOfWord,
				EndOfSentence,
			},
			b: &BPE{
				maxTokenLength: 64,
				vocab: map[string]struct{}{
					BeginOfWord + "т":   {},
					BeginOfWord + "те":  {},
					BeginOfWord + "тес": {},
					"т":                 {},
					"." + EndOfWord:     {},
					"ест." + EndOfWord:  {},
				},
			},
		},
		{
			name: "Lorem Ipsum",
			in:   strings.NewReader("Lorem Ipsum"),
			expected: []string{
				BeginOfSentence,
				BeginOfWord + "Lorem" + EndOfWord,
				BeginOfWord + "Ipsum" + EndOfWord,
				EndOfSentence,
			},
			b: &BPE{
				maxTokenLength: 64,
				vocab: map[string]struct{}{
					BeginOfWord + "Lorem" + EndOfWord: {},
					BeginOfWord + "Ipsum" + EndOfWord: {},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.b.Encode(tc.in)
			if err != nil && !tc.withError {
				t.Fatalf("Unexpected error: %v\n", err)
				return
			}

			if tc.withError {
				t.Fatalf("Error expected got: %v\n", actual)
			}

			if len(tc.expected) != len(actual) {
				t.Fatalf("Expected: %v\nGot: %v\n", tc.expected, actual)
				return
			}

			for i, token := range tc.expected {
				if token != actual[i] {
					t.Errorf(
						"Expected token at index %d is not matched with actual: \nExpected: %v\nGot: %v\n",
						i,
						token,
						actual[i],
					)
				}
			}
		})
	}
}
