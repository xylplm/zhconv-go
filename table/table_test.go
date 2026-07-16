package table

import (
	"strings"
	"testing"
)

func TestLoadTSVIgnoresCommentsAndBlank(t *testing.T) {
	raw := "# comment\n\n軟體\t软件\n網路\t网络\n"
	ms, err := LoadTSV([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if len(ms) != 2 || ms[0].From != "軟體" || ms[0].To != "软件" {
		t.Fatalf("unexpected mappings: %+v", ms)
	}
}

func TestLoadTSVCRLF(t *testing.T) {
	raw := "軟體\t软件\r\n網路\t网络\r\n"
	ms, err := LoadTSV([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if len(ms) != 2 {
		t.Fatalf("got %d mappings", len(ms))
	}
}

func TestLoadTSVFirstCandidateOnly(t *testing.T) {
	raw := "乾紅\t干红 乾红\n"
	ms, err := LoadTSV([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if len(ms) != 1 || ms[0].To != "干红" {
		t.Fatalf("expected first candidate, got %+v", ms)
	}
}

func TestLoadTSVRejectsMissingTab(t *testing.T) {
	if _, err := LoadTSV([]byte("bad-line-without-tab\n")); err == nil {
		t.Fatal("expected error")
	}
}

func TestDefaultTablesNonEmpty(t *testing.T) {
	chars, err := DefaultChars()
	if err != nil {
		t.Fatal(err)
	}
	phrases, err := DefaultPhrases()
	if err != nil {
		t.Fatal(err)
	}
	if len(chars) < 1000 {
		t.Fatalf("chars too few: %d", len(chars))
	}
	if len(phrases) < 200 {
		t.Fatalf("phrases too few: %d", len(phrases))
	}
	// Ensure no empty targets.
	for _, m := range append(chars[:3], phrases[:3]...) {
		if strings.TrimSpace(m.From) == "" || strings.TrimSpace(m.To) == "" {
			t.Fatalf("empty mapping: %+v", m)
		}
	}
}
