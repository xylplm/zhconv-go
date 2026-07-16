package table

import _ "embed"

//go:embed chars.tsv
var charsTSV []byte

//go:embed phrases.tsv
var phrasesTSV []byte

// DefaultChars returns the embedded traditional->simplified character mappings.
func DefaultChars() ([]Mapping, error) {
	return LoadTSV(charsTSV)
}

// DefaultPhrases returns the embedded traditional->simplified phrase mappings.
func DefaultPhrases() ([]Mapping, error) {
	return LoadTSV(phrasesTSV)
}
