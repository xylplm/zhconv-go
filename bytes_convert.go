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

// stringToBytesCopy returns a newly allocated mutable copy of s.
func stringToBytesCopy(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	out := make([]byte, len(s))
	copy(out, s)
	return out
}
