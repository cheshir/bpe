package bpe

import (
	"unicode"
	"unicode/utf8"
)

// Scan sentences.
// Sentence starts from the beginning of string or from the previous sentence
// and continues up to the EOF, end of line or .!? symbols with several heuristics.
func scanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0

	// Skip leading spaces.
	for pos, symbol := range string(data) {
		if !unicode.IsSpace(symbol) {
			break
		}

		start = pos
	}

	// Scan until EOF, EOL or .!? symbol.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		if isEndOfSentence(r, data[start:i], data[i:]) {
			return i + width, data[start : i+width], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated sentence. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
}
