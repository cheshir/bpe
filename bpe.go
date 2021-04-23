package bpe

import (
	"sort"
)

type BPE struct {
	maxTokenLength int
	vocab          map[string]struct{} // Set with fast vocab search.
}

type weightedToken struct {
	Token  *string
	Weight int
}

func newModelFromTokensFrequencyTable(tft tokensFrequencyTable, tokensLimit, maxTokenLength int) *BPE {
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

	vocab := make(map[string]struct{}, len(tokensListWithWeights))

	for _, t := range tokensListWithWeights {
		vocab[*t.Token] = struct{}{}
	}

	return &BPE{
		maxTokenLength: maxTokenLength,
		vocab:          vocab,
	}
}
