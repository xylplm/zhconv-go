package zhconv

import "sync"

var (
	defaultOnce sync.Once
	defaultConv *Converter
	defaultErr  error
)

func defaultConverter() *Converter {
	defaultOnce.Do(func() {
		defaultConv, defaultErr = New(Options{})
	})
	if defaultErr != nil {
		// Embedded dictionaries are part of the module build; a load failure is fatal.
		panic("zhconv: load default converter: " + defaultErr.Error())
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
