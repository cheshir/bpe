package bpe

import (
	"sort"
	"strings"
)

type BPE struct {
	vocab map[string]struct{} // Set with fast vocab search.
}

func (m *BPE) String() string {
	builder := strings.Builder{}
	builder.WriteString("Vocabulary: ")

	if len(m.vocab) == 0 {
		return "empty"
	}

	for t := range m.vocab {
		builder.WriteByte('"')
		builder.WriteString(t)
		builder.WriteString(`", `)
	}

	result := builder.String()

	return result[:len(result)-2]
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

	vocab := make(map[string]struct{}, len(tokensListWithWeights))

	for _, t := range tokensListWithWeights {
		vocab[*t.Token] = struct{}{}
	}

	return &BPE{
		vocab: vocab,
	}
}
