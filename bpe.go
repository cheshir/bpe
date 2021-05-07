package bpe

import (
	"bufio"
	"io"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var defaultTokensCap = 8

type BPE struct {
	maxTokenLength int
	vocab          map[string]struct{} // Set with fast vocab search.
}

type weightedToken struct {
	Token  *string
	Weight int
}

func (b *BPE) Encode(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanSentences)
	tokens := make([]string, 0, defaultTokensCap)

	for scanner.Scan() {
		sentence := scanner.Text()
		b.encodeSentence(&tokens, sentence)
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "file scan")
	}

	return tokens, nil
}

// Target is a pointer to slice of tokens because it helps avoid unnecessary memory allocations.
func (b *BPE) encodeSentence(target *[]string, sentence string) {
	*target = append(*target, BeginOfSentence)
	words := strings.Fields(sentence)
	for _, word := range words {
		b.encodeWord(target, word)
	}
	*target = append(*target, EndOfSentence)
}

func (b *BPE) encodeWord(target *[]string, word string) {
	word = BeginOfWord + word + EndOfWord // TODO use special tokens from BPE.
	tokenStart := 0

tokenLoop:
	for tokenStart < len(word) {
		tokenEnd := len(word)

		if tokenEnd-tokenStart > b.maxTokenLength {
			tokenEnd = tokenStart + b.maxTokenLength
		}

		for ; tokenEnd != tokenStart; tokenEnd-- {
			token := word[tokenStart:tokenEnd]
			_, ok := b.vocab[token]

			if ok {
				*target = append(*target, token)
				tokenStart += len(token)
				continue tokenLoop
			}
		}

		*target = append(*target, UnknownToken)
		tokenStart++
	}
}

// Decode todo description.
// Error in response added for potential future usages to keep backward compatibility.
func (b *BPE) Decode(tokens []string) (string, error) {
	builder := strings.Builder{}

	for _, token := range tokens {
		// Skip special tokens.
		// TODO Use special tokens from BPE.
		token = strings.TrimSuffix(token, BeginOfSentence)
		token = strings.TrimSuffix(token, EndOfSentence)
		token = strings.TrimSuffix(token, EndOfWord)

		if strings.HasPrefix(token, BeginOfWord) {
			builder.WriteByte(' ')
			token = token[len(BeginOfWord):]
		}

		_, err := builder.WriteString(token)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(builder.String()), nil
}

func newModelFromTokensFrequencyTable(tft tokensFrequencyTable, tokensLimit int) *BPE {
	tokensListWithWeights := make([]weightedToken, 0, len(tft))

	for t, w := range tft {
		token := t
		tokensListWithWeights = append(tokensListWithWeights, weightedToken{
			Token:  &token,
			Weight: w,
		})
	}

	sort.Slice(tokensListWithWeights, func(i, j int) bool {
		return tokensListWithWeights[i].Weight > tokensListWithWeights[j].Weight
	})

	if len(tokensListWithWeights) > tokensLimit {
		tokensListWithWeights = tokensListWithWeights[:tokensLimit]
	}

	var maxTokenLength int
	vocab := make(map[string]struct{}, len(tokensListWithWeights))

	for _, t := range tokensListWithWeights {
		token := *t.Token

		// TODO consider removing it and using value from config.
		// Need to check necessity for this change with benchmarks.
		tokenLength := len(token)
		if len(token) > maxTokenLength {
			maxTokenLength = tokenLength
		}

		vocab[token] = struct{}{}
	}

	return &BPE{
		maxTokenLength: maxTokenLength,
		vocab:          vocab,
	}
}
