package bpe

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
)

const (
	defaultMaxNumberOfTokens = 50000
	defaultMaxTokenLength    = 32
	maxScanBufferSize        = 64 * 1024

	BeginOfWord     = "<w>"
	EndOfWord       = "</w>"
	BeginOfSentence = "<s>"
	EndOfSentence   = "</s>"
	UnknownToken    = "<u>"
)

// Train returns BPE instance with vocabulary learned from source.
func Train(source io.Reader, opts ...TrainOption) (*BPE, error) {
	options := defaultTrainOptions()
	options.Apply(opts...)

	tft, err := calculateTokensFrequency(source, options)
	if err != nil {
		return nil, err
	}

	model := newModelFromTokensFrequencyTable(tft, options.MaxNumberOfTokens)

	return model, nil
}

func defaultTrainOptions() *trainOptions {
	return &trainOptions{
		MaxNumberOfTokens: defaultMaxNumberOfTokens,
		MaxTokenLength:    defaultMaxTokenLength,
		ScanBufferSize:    maxScanBufferSize,
	}
}

type trainOptions struct {
	MaxNumberOfTokens int
	MaxTokenLength    int
	ScanBufferSize    int
	WordsOnly         bool
}

func (o *trainOptions) Apply(opts ...TrainOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type TrainOption func(opts *trainOptions)

func WithMaxNumberOfTokens(n int) TrainOption {
	return func(opts *trainOptions) {
		opts.MaxNumberOfTokens = n
	}
}

func WithMaxTokenLength(length int) TrainOption {
	return func(opts *trainOptions) {
		opts.MaxTokenLength = length
	}
}

func WithScanBufferSize(size int) TrainOption {
	return func(opts *trainOptions) {
		opts.ScanBufferSize = size
	}
}

func WithWordsOnly() TrainOption {
	return func(opts *trainOptions) {
		opts.WordsOnly = true
	}
}

type tokensFrequencyTable map[string]int

func calculateTokensFrequency(r io.Reader, options *trainOptions) (tokensFrequencyTable, error) {
	tokensFrequency := make(tokensFrequencyTable, options.MaxNumberOfTokens) // Approximate size. Avoid extra allocations.
	scanner := bufio.NewScanner(r)
	scanner.Split(scanSentences)
	scanner.Buffer(make([]byte, 0, options.ScanBufferSize), options.ScanBufferSize)

	// TODO read in separate threads.
	for scanner.Scan() {
		sentence := scanner.Text()
		tokenize(tokensFrequency, sentence, options.MaxTokenLength, options.WordsOnly)
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "file scan")
	}

	return tokensFrequency, nil
}

const abbreviationLength = 4

// isEndOfSentence checks some heuristics to understand whether it's the end of sentence or not.
func isEndOfSentence(lastSymbol rune, prev, next []byte) bool {
	if len(prev) == 0 {
		return false
	}

	switch lastSymbol {
	case '\r', '\n', '!', '?':
		return true
	case '.':
		// Sad but dot isn't explicit marker of the end of sentence.
		// It can be used for name or other abbreviation.

		// Check last runes.
		// If they're \w{abbreviationLength} – looks like a Mrs., Dr. or other abbreviation.
		// It previous and next rune are numbers – it's a float.
		// Otherwise it looks like an end of sentence.
		prevRune, width := utf8.DecodeLastRune(prev)
		if unicode.IsLetter(prevRune) {
			nextAfterCurrent := prevRune // Is needed to check letter capitalization after getting first space.
			margin := width

			// Let's find the space.
			for i := 0; i < abbreviationLength; i++ {
				if len(prev)-margin <= 0 {
					return false
				}

				currentRune, currentWidth := utf8.DecodeLastRune(prev[:len(prev)-margin])

				// We've found the space. Let's check that the next character after space is a capitalized letter.
				if unicode.IsSpace(currentRune) {
					return !unicode.IsUpper(nextAfterCurrent)
				}

				// Is token is inside some group?
				if strings.ContainsAny(string(currentRune), `[({"'`) {
					return false
				}

				// If it's not letter it's not an abbreviation.
				if !unicode.IsLetter(currentRune) {
					return true
				}

				nextAfterCurrent = currentRune
				margin += currentWidth
			}

			// If the last n characters was letters it's probably is the end of string.
			return true
		}

		if unicode.IsDigit(prevRune) {
			if len(next) == 0 {
				return true
			}

			nextRune, _ := utf8.DecodeRune(next)
			// Looks like it's a float number.
			if unicode.IsDigit(nextRune) {
				return false
			}
		}

		return true
	}

	return false
}

// Preserve Unicode symbols.
func tokenize(tft tokensFrequencyTable, sentence string, maxTokenLength int, wordsOnly bool) {
	words := strings.Fields(sentence)

	for _, word := range words {
		if wordsOnly && !isWord(word) {
			continue
		}

		wordTokens := strings.Split(word, "")

		// Add special tokens.
		wordTokens[0] = BeginOfWord + wordTokens[0]
		wordTokens[len(wordTokens)-1] = wordTokens[len(wordTokens)-1] + EndOfWord
		tokenizeWord(tft, wordTokens, maxTokenLength)
	}
}

func tokenizeWord(tft tokensFrequencyTable, word []string, maxTokenLength int) {
	for i, firstToken := range word {
		tft[firstToken]++

		b := strings.Builder{}
		b.WriteString(firstToken)

		for i2, token := range word[i+1:] {
			// Current index plus first token.
			if i2+1 >= maxTokenLength {
				break
			}

			b.WriteString(token)
			tft[b.String()]++
		}
	}
}

// isWord checks that given string contains only letters or dashes.
func isWord(word string) bool {
	for _, r := range word {
		if !unicode.IsLetter(r) && r != '-' {
			return false
		}
	}

	return true
}
