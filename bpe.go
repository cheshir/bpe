package bpe

import (
	"sort"
	"strings"
)

type BPE struct {
	maxTokenLength int
	vocab          map[string]struct{} // Set with fast vocab search.
}

type weightedToken struct {
	Token  *string
	Weight int
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
		if strings.HasPrefix(token, BeginOfWord) {
			tokenLength -= len(BeginOfWord)
		}
		if strings.HasSuffix(token, EndOfWord) {
			tokenLength -= len(EndOfWord)
		}

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
