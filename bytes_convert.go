package zhconv

import "unsafe"

// bytesToString reinterprets b as a string without copying.
// The returned string must not outlive a mutable b; only use when:
//   - b is a temporary read-only view of caller input (ConvertBytes), or
//   - b is a freshly allocated buffer that will not be mutated (Convert).
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
