package utils

import (
	"unicode"
	"unicode/utf8"
)

func ExtractReadableString(buf []byte) (s string) {
	var r rune
	var l, i int
	p := make([]byte, utf8.MaxRune, utf8.MaxRune)
	for {
		if r, l = utf8.DecodeRune(buf); r == utf8.RuneError {
			break
		}
		buf = buf[l:]
		if !unicode.IsControl(r) {
			if i = utf8.EncodeRune(p, r); i != 0 {
				s = s + string(p[0:i])
			}
		}
	}
	return
}
