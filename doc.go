// Package zhconv provides a lightweight, high-performance Traditional-to-Simplified
// Chinese converter implemented in pure Go.
//
// Design goals:
//   - One direction only: traditional -> simplified (zh-Hans oriented)
//   - Phrase-first longest match, then character fallback
//   - Embedded dictionaries, zero CGO, safe for concurrent use
//   - Small API surface for easy go get integration
//
// Quick start:
//
//	fmt.Println(zhconv.ToSimplified("軟體與網路連線"))
//	// 软件与网络连接
package zhconv
