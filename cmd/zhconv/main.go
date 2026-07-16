package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xylplm/zhconv-go"
)

func main() {
	inPath := flag.String("i", "", "input file path (default: stdin)")
	outPath := flag.String("o", "", "output file path (default: stdout)")
	demo := flag.Bool("demo", false, "run built-in demo samples and exit")
	flag.Parse()

	if *demo {
		runDemo()
		return
	}

	in, err := openInput(*inPath)
	if err != nil {
		fail(err)
	}
	defer in.Close()

	data, err := io.ReadAll(in)
	if err != nil {
		fail(err)
	}

	outText := zhconv.ToSimplified(string(data))

	if *outPath == "" {
		_, _ = os.Stdout.WriteString(outText)
		if !strings.HasSuffix(outText, "\n") && len(outText) > 0 {
			_, _ = os.Stdout.WriteString("\n")
		}
		return
	}
	if err := os.WriteFile(*outPath, []byte(outText), 0o644); err != nil {
		fail(err)
	}
}

func openInput(path string) (io.ReadCloser, error) {
	if strings.TrimSpace(path) == "" {
		return io.NopCloser(bufio.NewReader(os.Stdin)), nil
	}
	return os.Open(filepath.Clean(path))
}

func runDemo() {
	samples := []string{
		"軟體與網路連線",
		"請將繁體中文字幕轉成簡體",
		"資料庫程式設計師正在除錯",
		"影片訊息已儲存到記憶體",
		"這是一段已經是简体中文的文字",
		"Hello 世界 — 混合 ASCII",
	}
	fmt.Println("zhconv-go demo (traditional -> simplified)")
	fmt.Println(strings.Repeat("-", 48))
	for _, s := range samples {
		fmt.Printf("IN : %s\nOUT: %s\n\n", s, zhconv.ToSimplified(s))
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "zhconv: %v\n", err)
	os.Exit(1)
}
