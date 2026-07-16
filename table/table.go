package table

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Mapping is one traditional -> simplified rewrite rule.
type Mapping struct {
	From string
	To   string
}

// LoadTSV parses traditional\tsimplified lines.
// Blank lines and # comments are ignored. Multi-candidate values keep the first field.
func LoadTSV(data []byte) ([]Mapping, error) {
	if len(data) == 0 {
		return nil, nil
	}
	// Bound growth; embedded tables are a few thousand lines.
	out := make([]Mapping, 0, 4096)
	lineNo := 0
	for len(data) > 0 {
		lineNo++
		var line []byte
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			line = data[:i]
			data = data[i+1:]
		} else {
			line = data
			data = nil
		}
		// Trim CR for Windows-edited files.
		if n := len(line); n > 0 && line[n-1] == '\r' {
			line = line[:n-1]
		}
		// Fast skip blanks / comments without string alloc when possible.
		trim := bytes.TrimSpace(line)
		if len(trim) == 0 || trim[0] == '#' {
			continue
		}
		tab := bytes.IndexByte(trim, '\t')
		if tab < 0 {
			return nil, fmt.Errorf("table: line %d: missing tab separator", lineNo)
		}
		from := string(bytes.TrimSpace(trim[:tab]))
		toRaw := bytes.TrimSpace(trim[tab+1:])
		if len(toRaw) == 0 {
			return nil, fmt.Errorf("table: line %d: empty target", lineNo)
		}
		// OpenCC multi-candidate: keep first token.
		if sp := bytes.IndexByte(toRaw, ' '); sp >= 0 {
			toRaw = toRaw[:sp]
		}
		to := string(toRaw)
		if from == "" || from == to {
			continue
		}
		if !utf8.ValidString(from) || !utf8.ValidString(to) {
			return nil, fmt.Errorf("table: line %d: invalid utf-8", lineNo)
		}
		out = append(out, Mapping{From: from, To: to})
	}
	return out, nil
}

// MustLoadTSV is like LoadTSV but panics on error (for tests/tooling only).
func MustLoadTSV(data []byte) []Mapping {
	ms, err := LoadTSV(data)
	if err != nil {
		panic(err)
	}
	return ms
}

// JoinPreview is a tiny helper for tests/debug.
func JoinPreview(ms []Mapping, n int) string {
	if n <= 0 || len(ms) == 0 {
		return ""
	}
	if n > len(ms) {
		n = len(ms)
	}
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, ms[i].From+"→"+ms[i].To)
	}
	return strings.Join(parts, ", ")
}
