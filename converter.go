package zhconv

import (
	"strings"
	"unicode/utf8"

	"github.com/xylplm/zhconv-go/table"
)

// trieNode is a compact UTF-8 byte trie for longest-match phrase replacement.
type trieNode struct {
	// next maps the next UTF-8 byte to a child. nil until needed.
	next map[byte]*trieNode
	// replacement is set when a phrase ends at this node.
	// Using pointer avoids confusing empty-string replacements with "not terminal".
	replacement *string
}

func (n *trieNode) child(b byte) *trieNode {
	if n.next == nil {
		return nil
	}
	return n.next[b]
}

func (n *trieNode) ensureChild(b byte) *trieNode {
	if n.next == nil {
		n.next = make(map[byte]*trieNode)
	}
	ch := n.next[b]
	if ch == nil {
		ch = &trieNode{}
		n.next[b] = ch
	}
	return ch
}

// Converter performs traditional-to-simplified Chinese conversion.
// It is safe for concurrent use after construction.
type Converter struct {
	// chars maps a single traditional rune to simplified text (usually one rune).
	chars map[rune]string
	// phrases is the root of the phrase trie over UTF-8 bytes.
	phrases trieNode
	// hasPhrase is true if at least one phrase was installed.
	hasPhrase bool
	// maxPhraseBytes is the longest phrase length in bytes (upper bound for scanners).
	maxPhraseBytes int
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

// New builds a converter from Options. Nil fields fall back to embedded tables.
func New(opts Options) (*Converter, error) {
	chars := opts.Chars
	if chars == nil {
		var err error
		chars, err = table.DefaultChars()
		if err != nil {
			return nil, err
		}
	}
	phrases := opts.Phrases
	if phrases == nil && !opts.DisablePhrases {
		var err error
		phrases, err = table.DefaultPhrases()
		if err != nil {
			return nil, err
		}
	}
	if opts.DisablePhrases {
		phrases = nil
	}

	c := &Converter{
		chars: make(map[rune]string, len(chars)),
	}
	pendingMulti := make([]table.Mapping, 0)
	for _, m := range chars {
		// Prefer single-rune sources for the char map; multi-rune entries are rare
		// in normalized tables and are better expressed as phrases.
		r, size := utf8.DecodeRuneInString(m.From)
		if r == utf8.RuneError && size == 1 {
			continue
		}
		if size != len(m.From) {
			pendingMulti = append(pendingMulti, m)
			continue
		}
		c.chars[r] = m.To
	}
	// Phrase targets from regional tables may still contain traditional glyphs
	// (OpenCC chain: phrases then characters). Normalize targets with the char map.
	for _, m := range phrases {
		c.addPhrase(m.From, c.simplifyWithChars(m.To))
	}
	for _, m := range pendingMulti {
		c.addPhrase(m.From, c.simplifyWithChars(m.To))
	}
	return c, nil
}

// Default returns the process-wide shared traditional->simplified converter.
// It panics only if embedded dictionaries are corrupt, which is a build-time fault.
func Default() *Converter {
	return defaultConverter()
}

func (c *Converter) addPhrase(from, to string) {
	if from == "" || to == "" || from == to {
		return
	}
	node := &c.phrases
	for i := 0; i < len(from); i++ {
		node = node.ensureChild(from[i])
	}
	// First write wins to keep load order stable and avoid silent churn.
	if node.replacement == nil {
		repl := to
		node.replacement = &repl
	}
	c.hasPhrase = true
	if len(from) > c.maxPhraseBytes {
		c.maxPhraseBytes = len(from)
	}
}

// simplifyWithChars applies character-level mapping only (no phrase recursion).
func (c *Converter) simplifyWithChars(s string) string {
	if s == "" || len(c.chars) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			b.WriteByte(s[i])
			i++
			continue
		}
		if repl, ok := c.chars[r]; ok {
			b.WriteString(repl)
		} else {
			b.WriteString(s[i : i+size])
		}
		i += size
	}
	return b.String()
}

// Convert rewrites traditional Chinese in s to simplified Chinese.
// Non-Chinese text is preserved. Invalid UTF-8 bytes are passed through unchanged.
func (c *Converter) Convert(s string) string {
	if c == nil || s == "" {
		return s
	}
	// Fast path: no phrase table and no traditional hit is still O(n),
	// but we avoid builder growth guesses being too small.
	var b strings.Builder
	b.Grow(len(s))

	i := 0
	for i < len(s) {
		// Phrase longest-match over UTF-8 bytes.
		if c.hasPhrase {
			if to, n := c.matchPhrase(s, i); n > 0 {
				b.WriteString(to)
				i += n
				continue
			}
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			// Invalid byte: copy as-is.
			b.WriteByte(s[i])
			i++
			continue
		}
		if repl, ok := c.chars[r]; ok {
			b.WriteString(repl)
		} else {
			b.WriteString(s[i : i+size])
		}
		i += size
	}
	return b.String()
}

// matchPhrase returns the simplified replacement and matched byte length
// for the longest phrase starting at s[i:]. ok/n==0 means no match.
func (c *Converter) matchPhrase(s string, i int) (string, int) {
	node := &c.phrases
	bestTo := ""
	bestN := 0
	// Walk as far as the trie allows; keep the last terminal node.
	for j := i; j < len(s); j++ {
		node = node.child(s[j])
		if node == nil {
			break
		}
		if node.replacement != nil {
			bestTo = *node.replacement
			bestN = j - i + 1
		}
	}
	return bestTo, bestN
}

// ConvertBytes is a convenience wrapper around Convert.
func (c *Converter) ConvertBytes(p []byte) []byte {
	if c == nil || len(p) == 0 {
		return p
	}
	// Avoid holding two giant strings longer than needed.
	return []byte(c.Convert(string(p)))
}
