package bpe

import (
	"bufio"
	"io"
	"unicode"

	"github.com/pkg/errors"
)

// Train returns BPE instance with vocabulary learned from source.
func Train(source io.Reader, opts ...TrainOption) (*BPE, error) {
	options := &TrainOptions{}
	options.Apply(opts...)

	tft, err := calculateTokensFrequency(source, options)
	if err != nil {
		return nil, err
	}

	model := newModelFromTokensFrequencyTable(tft, options.MaxNumberOfTokens)

	return model, nil
}

type tokensFrequencyTable map[string]int

func calculateTokensFrequency(r io.Reader, options *TrainOptions) (tokensFrequencyTable, error) {
	tokensFrequency := make(tokensFrequencyTable, options.MaxNumberOfTokens) // Approximate size. Avoid extra allocations.
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 0, options.ScanBufferSize), options.ScanBufferSize)

	// TODO read in separate threads.
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				break
			}

			return nil, errors.Wrap(err, "file scan")
		}

		word := scanner.Text()
		tokenize(tokensFrequency, word, options.MaxTokenLength)
	}

	return tokensFrequency, nil
}

// Preserve Unicode symbols.
func tokenize(tokens tokensFrequencyTable, word string, maxTokenLength int) {
	for i, char := range word {
		tokens[string(char)]++

		if !isWordChar(char) {
			continue
		}

		for j, char2 := range word[i+1:] {
			i2 := i + j + 2

			if i2-i > maxTokenLength {
				break
			}

			if !isWordChar(char2) {
				break
			}

			tokens[word[i:i2]]++
		}
	}
}

// isWordChar checks that given word contains only letters and hyphens.
func isWordChar(char rune) bool {
	if !unicode.IsLetter(char) && char != '-' {
		return false
	}

	return true
}
