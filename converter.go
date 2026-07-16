package zhconv

import (
	"unicode/utf8"

	"github.com/xylplm/zhconv-go/table"
)

// phrase is one traditional phrase rule, pre-tokenized for fast matching.
type phrase struct {
	from  []rune
	to    string
	bytes int // UTF-8 byte length of the traditional form
}

// Converter performs traditional-to-simplified Chinese conversion.
// It is safe for concurrent use after construction.
//
// Matching order:
//  1. longest phrase that starts at the current rune (phrase-first)
//  2. single-rune character map
//  3. original text
type Converter struct {
	// char1 holds 1:1 traditional->simplified rune replacements (hot path).
	char1 map[rune]rune
	// charN holds rare traditional runes that expand to multi-rune simplified text.
	charN map[rune]string
	// phrases indexes traditional phrases by their first rune, longest-first.
	phrases map[rune][]phrase
	// hasPhrase avoids map lookup when no phrases are installed.
	hasPhrase bool
}

// Options controls converter construction.
type Options struct {
	// Chars overrides the embedded character table when non-nil.
	Chars []table.Mapping
	// Phrases overrides the embedded phrase table when non-nil.
	Phrases []table.Mapping
	// DisablePhrases skips phrase matching and only applies character mapping.
	DisablePhrases bool
}

// New builds a converter from Options. Nil table fields fall back to embedded data.
func New(opts Options) (*Converter, error) {
	chars := opts.Chars
	if chars == nil {
		var err error
		chars, err = table.DefaultChars()
		if err != nil {
			return nil, err
		}
	}

	var phrases []table.Mapping
	if !opts.DisablePhrases {
		phrases = opts.Phrases
		if phrases == nil {
			var err error
			phrases, err = table.DefaultPhrases()
			if err != nil {
				return nil, err
			}
		}
	}

	c := &Converter{
		char1: make(map[rune]rune, len(chars)),
	}

	pendingMulti := make([]table.Mapping, 0)
	for _, m := range chars {
		r, size := utf8.DecodeRuneInString(m.From)
		if r == utf8.RuneError && size == 1 {
			continue
		}
		// Multi-rune sources belong in the phrase table.
		if size != len(m.From) {
			pendingMulti = append(pendingMulti, m)
			continue
		}
		// Prefer compact 1:1 rune map when target is a single rune.
		tr, tsize := utf8.DecodeRuneInString(m.To)
		if tsize == len(m.To) && tr != utf8.RuneError {
			c.char1[r] = tr
			continue
		}
		if c.charN == nil {
			c.charN = make(map[rune]string)
		}
		c.charN[r] = m.To
	}

	// Phrase targets may still contain traditional glyphs (OpenCC-style chains).
	// Normalize with the character map once at load time.
	allPhrases := make([]table.Mapping, 0, len(phrases)+len(pendingMulti))
	allPhrases = append(allPhrases, phrases...)
	allPhrases = append(allPhrases, pendingMulti...)
	for _, m := range allPhrases {
		to := c.simplifyWithChars(m.To)
		c.addPhrase(m.From, to)
	}
	c.finalizePhrases()
	return c, nil
}

// Default returns the process-wide shared converter.
// If embedded dictionaries fail to load (should be impossible after tests),
// an empty identity converter is returned and Convert becomes a no-op.
func Default() *Converter {
	return defaultConverter()
}

func (c *Converter) addPhrase(from, to string) {
	if c == nil || from == "" || to == "" || from == to {
		return
	}
	runes := []rune(from)
	if len(runes) == 0 {
		return
	}
	// Single-rune "phrase" can live in the char maps if missing.
	if len(runes) == 1 {
		r := runes[0]
		if _, ok := c.char1[r]; ok {
			return
		}
		if _, ok := c.charN[r]; ok {
			return
		}
		tr, tsize := utf8.DecodeRuneInString(to)
		if tsize == len(to) && tr != utf8.RuneError {
			c.char1[r] = tr
			return
		}
		if c.charN == nil {
			c.charN = make(map[rune]string)
		}
		c.charN[r] = to
		return
	}

	if c.phrases == nil {
		c.phrases = make(map[rune][]phrase)
	}
	start := runes[0]
	c.phrases[start] = append(c.phrases[start], phrase{
		from:  runes,
		to:    to,
		bytes: len(from),
	})
	c.hasPhrase = true
}

// finalizePhrases sorts each bucket longest-first so the first hit is maximal.
func (c *Converter) finalizePhrases() {
	if c == nil || !c.hasPhrase {
		return
	}
	for k, list := range c.phrases {
		// Insertion-sort by rune length desc (lists are small).
		for i := 1; i < len(list); i++ {
			j := i
			for j > 0 && len(list[j-1].from) < len(list[j].from) {
				list[j-1], list[j] = list[j], list[j-1]
				j--
			}
		}
		// First mapping wins for equal-length duplicates.
		dedup := list[:0]
		seen := make(map[string]struct{}, len(list))
		for _, p := range list {
			key := string(p.from)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			dedup = append(dedup, p)
		}
		c.phrases[k] = dedup
	}
}

// simplifyWithChars applies character-level mapping only (no phrase recursion).
func (c *Converter) simplifyWithChars(s string) string {
	if c == nil || s == "" {
		return s
	}
	out, changed := c.mapChars(s)
	if !changed {
		return s
	}
	return out
}

func (c *Converter) mapChars(s string) (string, bool) {
	if len(c.char1) == 0 && len(c.charN) == 0 {
		return s, false
	}
	buf := make([]byte, 0, len(s))
	changed := false
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			buf = append(buf, s[i])
			i++
			continue
		}
		if repl, ok := c.char1[r]; ok {
			buf = utf8.AppendRune(buf, repl)
			changed = true
		} else if repl, ok := c.charN[r]; ok {
			buf = append(buf, repl...)
			changed = true
		} else {
			buf = append(buf, s[i:i+size]...)
		}
		i += size
	}
	if !changed {
		return s, false
	}
	return string(buf), true
}

// Convert rewrites traditional Chinese in s to simplified Chinese.
// Non-Chinese text is preserved. Invalid UTF-8 bytes pass through unchanged.
// If no replacement occurs, the original string is returned (no allocation).
func (c *Converter) Convert(s string) string {
	if c == nil || s == "" {
		return s
	}

	// Lazy builder: only allocate when the first replacement happens.
	var buf []byte
	started := false
	i := 0
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			if started {
				buf = append(buf, s[i])
			}
			i++
			continue
		}

		// 1) Longest phrase starting with this rune.
		if c.hasPhrase {
			if to, nBytes := c.matchPhraseAt(s, i, r); nBytes > 0 {
				if !started {
					buf = make([]byte, 0, len(s))
					buf = append(buf, s[:i]...)
					started = true
				}
				buf = append(buf, to...)
				i += nBytes
				continue
			}
		}

		// 2) Single-rune character map.
		if repl, ok := c.char1[r]; ok {
			if !started {
				buf = make([]byte, 0, len(s))
				buf = append(buf, s[:i]...)
				started = true
			}
			buf = utf8.AppendRune(buf, repl)
			i += size
			continue
		}
		if repl, ok := c.charN[r]; ok {
			if !started {
				buf = make([]byte, 0, len(s))
				buf = append(buf, s[:i]...)
				started = true
			}
			buf = append(buf, repl...)
			i += size
			continue
		}

		// 3) Keep original.
		if started {
			buf = append(buf, s[i:i+size]...)
		}
		i += size
	}
	if !started {
		return s
	}
	return string(buf)
}

// matchPhraseAt returns replacement and matched traditional byte length.
// Zero-allocation: compares candidate runes directly against s without a temp buffer.
func (c *Converter) matchPhraseAt(s string, byteIndex int, first rune) (string, int) {
	list := c.phrases[first]
	if len(list) == 0 {
		return "", 0
	}
	// list is sorted longest-first; first hit is maximal.
	for _, p := range list {
		// Remaining input must hold at least the phrase byte length.
		if byteIndex+p.bytes > len(s) {
			continue
		}
		// First rune already matched via bucket key.
		j := byteIndex + utf8.RuneLen(first)
		ok := true
		for k := 1; k < len(p.from); k++ {
			if j >= len(s) {
				ok = false
				break
			}
			rr, sz := utf8.DecodeRuneInString(s[j:])
			if rr != p.from[k] {
				ok = false
				break
			}
			j += sz
		}
		if ok {
			return p.to, p.bytes
		}
	}
	return "", 0
}

// ConvertBytes converts traditional Chinese bytes to simplified Chinese.
// When no change is required, the input slice is returned as-is.
func (c *Converter) ConvertBytes(p []byte) []byte {
	if c == nil || len(p) == 0 {
		return p
	}
	s := string(p)
	out := c.Convert(s)
	// Convert returns s when unchanged; keep the caller's slice to avoid a copy.
	if out == s {
		return p
	}
	return []byte(out)
}
