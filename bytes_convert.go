package zhconv

import "unsafe"

// bytesToStringRO returns a string view over b without copying.
// The caller must not mutate b while the string is in use.
func bytesToStringRO(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
