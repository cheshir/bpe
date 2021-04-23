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

func TestTokenize(t *testing.T) {
	const maxTokenSize = 3

	tt := []struct {
		word     string
		expected tokensFrequencyTable
	}{
		{
			word: "12314",
			expected: tokensFrequencyTable{
				"1": 2,
				"2": 1,
				"3": 1,
				"4": 1,
			},
		},
		{
			word: "abcad",
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
			word: "a-b",
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
			word: "[xxx]",
			expected: tokensFrequencyTable{
				"[":   1,
				"]":   1,
				"x":   3,
				"xx":  2,
				"xxx": 1,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.word, func(t *testing.T) {
			actualTokens := make(tokensFrequencyTable, 0)
			tokenize(actualTokens, tc.word, maxTokenSize)

			if !reflect.DeepEqual(tc.expected, actualTokens) {
				t.Errorf("Expected: %v\nGot: %v\n", tc.expected, actualTokens)
			}
		})
	}
}
