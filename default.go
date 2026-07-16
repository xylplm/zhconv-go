package zhconv

import "sync"

var (
	defaultOnce sync.Once
	defaultConv *Converter
)

// identity returns an empty converter that leaves text unchanged.
// All maps stay nil; Convert/ConvertBytes treat nil maps as empty.
func identity() *Converter {
	return &Converter{}
}

func defaultConverter() *Converter {
	defaultOnce.Do(func() {
		c, err := New(Options{})
		if err != nil || c == nil {
			// Embedded dictionaries are validated by tests. Fall back to a safe
			// no-op converter instead of panicking in production call paths.
			defaultConv = identity()
			return
		}
		defaultConv = c
	})
	// Defensive: Once should always assign, but never hand a nil pointer to callers.
	if defaultConv == nil {
		return identity()
	}
	return defaultConv
}

// ToSimplified converts traditional Chinese to simplified Chinese using the
// shared default converter.
func ToSimplified(s string) string {
	return Default().Convert(s)
}

// ToSimplifiedBytes converts traditional Chinese bytes to simplified Chinese.
func ToSimplifiedBytes(p []byte) []byte {
	return Default().ConvertBytes(p)
}

// NewDefault is an explicit constructor equivalent to New(Options{}).
// Prefer Default()/ToSimplified for normal use.
func NewDefault() (*Converter, error) {
	return New(Options{})
}
