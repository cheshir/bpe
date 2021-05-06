package bpe

import (
	"testing"
)

func TestIsEndOfSentence(t *testing.T) {
	tt := []struct {
		name     string
		prev     string
		lastRune rune
		next     string
		expected bool
	}{
		{
			name:     "line start",
			prev:     "",
			lastRune: 'E',
			next:     "xample",
			expected: false,
		},
		{
			name:     "single new line",
			prev:     "",
			lastRune: '\n',
			next:     "",
			expected: false,
		},
		{
			name:     "123.45",
			prev:     "123",
			lastRune: '.',
			next:     "45 text",
			expected: false,
		},
		{
			name:     "float with EOF",
			prev:     "123",
			lastRune: '.',
			next:     "",
			expected: true,
		},
		{
			name:     "list item number",
			prev:     "1",
			lastRune: '.',
			next:     " First",
			// Yeah. It's hard to understand is it a list item number or something like year et the end of string.
			expected: true,
		},
		{
			name:     "!",
			prev:     "Wow",
			lastRune: '!',
			next:     " ",
			expected: true,
		},
		{
			name:     "?",
			prev:     "Really",
			lastRune: '?',
			next:     " ",
			expected: true,
		},
		{
			name:     "eof",
			prev:     "Really",
			lastRune: '?',
			next:     " ",
			expected: true,
		},
		{
			name:     "abbreviation",
			prev:     "he's a Dr",
			lastRune: '.',
			next:     " of",
			expected: false,
		},
		{
			name:     "abbreviation 2",
			prev:     "Dr",
			lastRune: '.',
			next:     " Michael",
			expected: false,
		},
		{
			name:     "example.",
			prev:     "example",
			lastRune: '.',
			next:     " New sentence",
			expected: true,
		},
		{
			name:     "(lat",
			prev:     "(lat",
			lastRune: '.',
			next:     " Loren Ipsum",
			expected: false,
		},
		{
			name:     "4lat",
			prev:     "4lat",
			lastRune: '.',
			next:     " New sentence",
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := isEndOfSentence(tc.lastRune, []byte(tc.prev), []byte(tc.next))
			if tc.expected != actual {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, actual)
			}
		})
	}
}
