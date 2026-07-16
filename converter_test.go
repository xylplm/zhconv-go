package zhconv

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/xylplm/zhconv-go/table"
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

func TestTaiwanRegionalPhrases(t *testing.T) {
	cases := map[string]string{
		"伺服器": "服务器",
		"檔案":  "文件",
		"螢幕":  "屏幕",
		"印表機": "打印机",
		"光碟":  "光盘",
		"韌體":  "固件",
		"晶片":  "芯片",
		"迴圈":  "循环",
		"物件":  "对象",
		"介面":  "界面",
		"函式":  "函数",
		"變數":  "变量",
		"字串":  "字符串",
		"布林":  "布尔",
		"計程車": "出租车",
		"匯出":  "导出",
		"匯入":  "导入",
		"佇列":  "队列",
	}
	for in, want := range cases {
		if got := ToSimplified(in); got != want {
			t.Fatalf("TW phrase %q => %q, want %q", in, got, want)
		}
	}
}

func TestRegionalCharacterVariants(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"裏面", "里面"},
		{"裡面", "里面"},
		{"啓動", "启动"},
		{"啟動", "启动"},
		{"僞造", "伪造"},
		{"偽造", "伪造"},
		{"羣衆", "群众"},
		{"綫路", "线路"},
		{"線路", "线路"},
		{"說明", "说明"},
		{"説明", "说明"},
	}
	for _, tc := range cases {
		if got := ToSimplified(tc.in); got != tc.want {
			t.Fatalf("variant %q => %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSubtitleLikeParagraph(t *testing.T) {
	in := "主角安裝了最新軟體，透過網際網路連線到資料庫伺服器，螢幕上顯示系統訊息。"
	got := ToSimplified(in)
	for _, bad := range []string{"軟體", "網際", "網路", "資料庫", "伺服器", "螢幕", "訊息", "連線", "透過", "安裝", "顯示"} {
		if strings.Contains(got, bad) {
			t.Fatalf("traditional fragment %q still present in: %q", bad, got)
		}
	}
	for _, need := range []string{"软件", "互联", "数据库", "服务器", "屏幕", "消息"} {
		if !strings.Contains(got, need) {
			t.Fatalf("expected %q in output: %q", need, got)
		}
	}
}

func TestConvertKeepsASCIIAndAlreadySimplified(t *testing.T) {
	in := "Hello 世界 123 软件与网络"
	got := ToSimplified(in)
	if got != in {
		t.Fatalf("simplified/ascii text changed: %q -> %q", in, got)
	}
}

func TestConvertIsIdempotentOnSimplifiedOutput(t *testing.T) {
	in := "這是一段繁體測試，包含軟體、網路與資料庫。"
	once := ToSimplified(in)
	twice := ToSimplified(once)
	if once != twice {
		t.Fatalf("not idempotent:\n1: %q\n2: %q", once, twice)
	}
}

func TestConvertPhraseLongestMatch(t *testing.T) {
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
	if got := c.ConvertBytes(nil); got != nil {
		t.Fatalf("nil converter ConvertBytes(nil)=%v", got)
	}
	if got := c.ConvertBytes([]byte("軟體")); string(got) != "軟體" {
		t.Fatalf("nil converter should pass bytes through, got %q", got)
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

func TestConvertBytesRoundTripValidText(t *testing.T) {
	in := []byte("繁體中文軟體")
	got := ToSimplifiedBytes(in)
	if string(got) != "繁体中文软件" {
		t.Fatalf("ConvertBytes unexpected: %q", got)
	}
}

func TestConvertBytesNoChangeReturnsSameSlice(t *testing.T) {
	in := []byte("Hello 123")
	got := ToSimplifiedBytes(in)
	if &got[0] != &in[0] {
		t.Fatal("no-op ConvertBytes should return the input slice")
	}
}

func TestConvertBytesEmpty(t *testing.T) {
	if got := ToSimplifiedBytes(nil); got != nil {
		t.Fatalf("nil input => %v", got)
	}
	in := []byte{}
	got := ToSimplifiedBytes(in)
	if len(got) != 0 {
		t.Fatalf("empty input len=%d", len(got))
	}
}

func TestCustomTables(t *testing.T) {
	c, err := New(Options{
		Chars: []table.Mapping{
			{From: "測", To: "测"},
			{From: "試", To: "试"},
			{From: "組", To: "组"},
			{From: "", To: "x"},                   // ignored
			{From: "坏", To: "坏"},                  // noop ignored
			{From: string([]byte{0xff}), To: "x"}, // invalid ignored
		},
		Phrases: []table.Mapping{
			{From: "測試詞", To: "测试词"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Convert("測試詞"); got != "测试词" {
		t.Fatalf("custom phrase: %q", got)
	}
	if got := c.Convert("測試"); got != "测试" {
		t.Fatalf("custom chars: %q", got)
	}
	if got := c.Convert("測試詞組"); got != "测试词组" {
		t.Fatalf("phrase then char: %q", got)
	}
}

func TestMatchLongestPhraseOverShorter(t *testing.T) {
	c, err := New(Options{
		Chars: []table.Mapping{}, // no embedded fallbacks
		Phrases: []table.Mapping{
			{From: "軟體", To: "软件"},
			{From: "軟體工程", To: "软件工程"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Without char map, 師 stays traditional; longest phrase must still win.
	if got := c.Convert("軟體工程師"); got != "软件工程師" {
		t.Fatalf("longest phrase expected, got %q", got)
	}
	if got := c.Convert("軟體工程"); got != "软件工程" {
		t.Fatalf("exact longest phrase: %q", got)
	}
	if got := c.Convert("軟體"); got != "软件" {
		t.Fatalf("shorter phrase: %q", got)
	}
}

func TestIdentityFallbackNeverNils(t *testing.T) {
	c := identity()
	if c == nil {
		t.Fatal("identity nil")
	}
	if got := c.Convert("軟體"); got != "軟體" {
		t.Fatalf("identity should no-op, got %q", got)
	}
	in := []byte("abc")
	if got := c.ConvertBytes(in); &got[0] != &in[0] {
		t.Fatal("identity ConvertBytes should return input")
	}
}

func TestDisablePhrases(t *testing.T) {
	c, err := New(Options{DisablePhrases: true})
	if err != nil {
		t.Fatal(err)
	}
	got := c.Convert("軟體")
	if strings.ContainsAny(got, "軟體") {
		t.Fatalf("char-level should still convert glyphs: %q", got)
	}
	if got == "软件" {
		t.Fatalf("DisablePhrases unexpectedly produced phrase result %q", got)
	}
}

func TestDefaultSingleton(t *testing.T) {
	a := Default()
	b := Default()
	if a != b {
		t.Fatal("Default() should return the same instance")
	}
	if a == nil {
		t.Fatal("Default() returned nil")
	}
}

func TestNewDefault(t *testing.T) {
	c, err := NewDefault()
	if err != nil {
		t.Fatal(err)
	}
	if c.Convert("軟體") != "软件" {
		t.Fatalf("NewDefault convert failed: %q", c.Convert("軟體"))
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

func TestNoPanicOnAllRuneEdges(t *testing.T) {
	samples := []string{
		"",
		"a",
		"中",
		"繁體",
		strings.Repeat("軟體網路", 1000),
		string(rune(0)),
		string([]byte{0x80}),
	}
	for _, s := range samples {
		func(s string) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic on %q: %v", s, r)
				}
			}()
			_ = ToSimplified(s)
			_ = utf8.ValidString(ToSimplified(s[:min(len(s), 3)]))
		}(s)
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

func BenchmarkToSimplifiedNoChange(b *testing.B) {
	in := strings.Repeat("Hello 世界 software network database 12345 ", 40)
	c := Default()
	b.ReportAllocs()
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for b.Loop() {
		_ = c.Convert(in)
	}
}

func BenchmarkConvertBytes(b *testing.B) {
	in := []byte(strings.Repeat("這是一段用於基準測試的繁體中文字幕，包含軟體、網路、資料庫、影片與訊息。", 20))
	c := Default()
	b.ReportAllocs()
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for b.Loop() {
		_ = c.ConvertBytes(in)
	}
}

func BenchmarkConvertBytesNoChange(b *testing.B) {
	in := []byte(strings.Repeat("Hello 世界 software network database 12345 ", 40))
	c := Default()
	b.ReportAllocs()
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for b.Loop() {
		_ = c.ConvertBytes(in)
	}
}
