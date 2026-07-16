package table

import (
	"bufio"
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
// Blank lines and # comments are ignored. Multi-candidate values are not expected
// in normalized dicts; if present, only the first field is used.
func LoadTSV(data []byte) ([]Mapping, error) {
	if len(data) == 0 {
		return nil, nil
	}
	out := make([]Mapping, 0, 1024)
	sc := bufio.NewScanner(bytes.NewReader(data))
	// Longest practical phrase is far below this; keep headroom for bad input.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		from, to, ok := strings.Cut(line, "\t")
		if !ok {
			return nil, fmt.Errorf("table: line %d: missing tab separator", lineNo)
		}
		from = strings.TrimSpace(from)
		to = strings.TrimSpace(to)
		if to == "" {
			return nil, fmt.Errorf("table: line %d: empty target", lineNo)
		}
		// Be tolerant if a raw OpenCC multi-value line sneaks in.
		if i := strings.IndexByte(to, ' '); i >= 0 {
			to = to[:i]
		}
		if from == "" || from == to {
			continue
		}
		if !utf8.ValidString(from) || !utf8.ValidString(to) {
			return nil, fmt.Errorf("table: line %d: invalid utf-8", lineNo)
		}
		out = append(out, Mapping{From: from, To: to})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("table: scan: %w", err)
	}
	return out, nil
}
