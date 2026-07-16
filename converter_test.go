package zhconv

import (
	"strings"
	"testing"
)

func TestToSimplifiedCommonPhrases(t *testing.T) {
	cases := map[string]string{
		"軟體":   "软件",
		"網路":   "网络",
		"程式":   "程序",
		"資料庫":  "数据库",
		"影片":   "视频",
		"訊息":   "消息",
		"記憶體":  "内存",
		"硬體":   "硬件",
		"滑鼠":   "鼠标",
		"繁體中文": "繁体中文",
	}
	for in, want := range cases {
		if got := ToSimplified(in); got != want {
			t.Fatalf("ToSimplified(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestConvertKeepsASCIIAndAlreadySimplified(t *testing.T) {
	in := "Hello 世界 123"
	if got := ToSimplified(in); got != in {
		t.Fatalf("simplified/ascii text changed: %q -> %q", in, got)
	}
}

func TestConvertPhraseLongestMatch(t *testing.T) {
	// 網際網路 should prefer the whole phrase over 網路 if both exist.
	in := "網際網路連線"
	got := ToSimplified(in)
	if got != "互联网连接" {
		t.Fatalf("unexpected conversion: %q", got)
	}
}

func TestConvertEmptyAndNilReceiver(t *testing.T) {
	if got := ToSimplified(""); got != "" {
		t.Fatalf("empty input: got %q", got)
	}
	var c *Converter
	if got := c.Convert("軟體"); got != "軟體" {
		t.Fatalf("nil converter should return input, got %q", got)
	}
}

func TestConvertInvalidUTF8Passthrough(t *testing.T) {
	in := string([]byte{0xff, 0xfe, 0xe8}) + "軟體"
	got := ToSimplified(in)
	if !strings.HasSuffix(got, "软件") {
		t.Fatalf("expected suffix simplified, got %q", got)
	}
	if got[0] != 0xff || got[1] != 0xfe {
		t.Fatalf("invalid prefix bytes were altered: %q", got)
	}
}

func TestDisablePhrases(t *testing.T) {
	c, err := New(Options{DisablePhrases: true})
	if err != nil {
		t.Fatal(err)
	}
	// Without phrases, 軟體 still becomes 软体 via chars, not necessarily 软件.
	got := c.Convert("軟體")
	if strings.ContainsAny(got, "軟體") {
		t.Fatalf("char-level should still convert glyphs: %q", got)
	}
}

func TestConcurrentConvert(t *testing.T) {
	c := Default()
	done := make(chan struct{}, 8)
	for range 8 {
		go func() {
			defer func() { done <- struct{}{} }()
			for range 200 {
				_ = c.Convert("這是一段繁體中文測試，包含軟體、網路與資料庫。")
			}
		}()
	}
	for range 8 {
		<-done
	}
}

func BenchmarkToSimplified(b *testing.B) {
	in := strings.Repeat("這是一段用於基準測試的繁體中文字幕，包含軟體、網路、資料庫、影片與訊息。", 20)
	c := Default()
	b.ReportAllocs()
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for b.Loop() {
		_ = c.Convert(in)
	}
}
