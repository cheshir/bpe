package bpe

import (
	"bufio"
	"io"
	"os"
	"unicode"

	"github.com/pkg/errors"
)

const maxTokenLength = 5
const scanBufferSize = 64 * 1024

func Train(file string, tokensLimit int) (*BPE, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	tft, err := buildFrequencyTable(f, tokensLimit)
	if err != nil {
		return nil, err
	}

	model := newModelFromTokensFrequencyTable(tft, tokensLimit)

	return model, nil
}

type tokenFrequencyTable map[string]int

func buildFrequencyTable(r io.Reader, tokensLimit int) (tokenFrequencyTable, error) {
	tokensFrequency := make(tokenFrequencyTable, tokensLimit) // Approximated size.
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 0, scanBufferSize), scanBufferSize)

	// TODO read in separate threads.
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				break
			}

			return nil, errors.Wrap(err, "file scan")
		}

		word := scanner.Text()
		tokenize(tokensFrequency, word, maxTokenLength)
	}

	return tokensFrequency, nil
}

// Preserve Unicode symbols.
func tokenize(tokens tokenFrequencyTable, word string, maxTokenLength int) {
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

func isWordChar(char rune) bool {
	if !unicode.IsLetter(char) && char != '-' {
		return false
	}

	return true
}
