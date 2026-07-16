package main

import (
	"fmt"

	"github.com/xylplm/zhconv-go"
)

func main() {
	samples := []string{
		"軟體與網路連線",
		"資料庫伺服器已啟動",
		"螢幕上顯示系統訊息",
		"裏面與裡面都應轉成简体",
	}

	for _, s := range samples {
		fmt.Printf("%s\n  => %s\n\n", s, zhconv.ToSimplified(s))
	}
}
