package utils

import (
	"unicode/utf8"
	"strings"
)

var asteriskRune rune

func init() {
	asteriskRune, _ = utf8.DecodeRuneInString("*")
}

// MatchAsterisk match pattern against a utf8 string, support asterisk (*) only
func MatchAsterisk(pattern string, s string) bool {
	// trim space s
	pattern = strings.TrimSpace(pattern)
	s = strings.TrimSpace(s)
	// no asterisk in pattern
	if len(pattern) == 0 || !strings.Contains(pattern, "*") {
		return pattern == s
	}
	// consume pattern and compare string
	var flag bool
	for len(pattern) > 0 {
		// find a rune in pattern and next
		r, size := utf8.DecodeRuneInString(pattern)
		if r == utf8.RuneError {
			return false
		}
		pattern = pattern[size:]

		// if asterisk rune, set flag
		if r == asteriskRune {
			flag = true
		} else {
			// target string is empty, returns false
			if len(s) == 0 {
				return false
			}
			// iterate target string find the rune
			for len(s) > 0 {
				// find a rune in target string and next
				r2, size2 := utf8.DecodeRuneInString(s)
				if r2 == utf8.RuneError {
					return false
				}
				s = s[size2:]
				// flag set, continue till rune found, conumes all same runes
				if flag {
					if r2 == r {
						// consume all same runes, makes * greedy
						for len(s) > 0 {
							r3, size3 := utf8.DecodeRuneInString(s)
							if r3 == utf8.RuneError {
								return false
							}
							// break before actually cut the target string
							if r3 != r2 {
								break
							}
							s = s[size3:]
						}
						flag = false
						break
					}
				} else {
					// no flag, not equal, returns false
					// no flag, equal, break immediately
					if r2 != r {
						return false
					}
					break
				}
			}
			// if flag is still true, rune not found
			if flag {
				return false
			}
		}
	}
	return true
}
