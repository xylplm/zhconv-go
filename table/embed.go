package table

import (
	_ "embed"
	"sync"
)

//go:embed chars.tsv
var charsTSV []byte

//go:embed phrases.tsv
var phrasesTSV []byte

var (
	charsOnce    sync.Once
	charsCache   []Mapping
	charsLoadErr error

	phrasesOnce    sync.Once
	phrasesCache   []Mapping
	phrasesLoadErr error
)

// DefaultChars returns the embedded traditional->simplified character mappings.
// The returned slice is shared and must be treated as read-only.
func DefaultChars() ([]Mapping, error) {
	charsOnce.Do(func() {
		charsCache, charsLoadErr = LoadTSV(charsTSV)
	})
	return charsCache, charsLoadErr
}

// DefaultPhrases returns the embedded traditional->simplified phrase mappings.
// The returned slice is shared and must be treated as read-only.
func DefaultPhrases() ([]Mapping, error) {
	phrasesOnce.Do(func() {
		phrasesCache, phrasesLoadErr = LoadTSV(phrasesTSV)
	})
	return phrasesCache, phrasesLoadErr
}
