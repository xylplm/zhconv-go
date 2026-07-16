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

// version is injected by -ldflags at release build time.
var version = "dev"

func main() {
	inPath := flag.String("i", "", "input file path (default: stdin)")
	outPath := flag.String("o", "", "output file path (default: stdout)")
	demo := flag.Bool("demo", false, "run built-in demo samples and exit")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}
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

	// ConvertBytes keeps the no-change path allocation-free for pure ASCII/simplified input.
	outBytes := zhconv.ToSimplifiedBytes(data)

	if *outPath == "" {
		_, _ = os.Stdout.Write(outBytes)
		if len(outBytes) > 0 && outBytes[len(outBytes)-1] != '\n' {
			_, _ = os.Stdout.Write([]byte{'\n'})
		}
		return
	}
	if err := os.WriteFile(*outPath, outBytes, 0o644); err != nil {
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
	fmt.Printf("zhconv-go demo (traditional -> simplified) [%s]\n", version)
	fmt.Println(strings.Repeat("-", 48))
	for _, s := range samples {
		fmt.Printf("IN : %s\nOUT: %s\n\n", s, zhconv.ToSimplified(s))
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "zhconv: %v\n", err)
	os.Exit(1)
}
