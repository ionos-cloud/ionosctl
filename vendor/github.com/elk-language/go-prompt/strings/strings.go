package strings

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// Get the length of the string in bytes.
func Len(s string) ByteNumber {
	return ByteNumber(len(s))
}

// Get the length of the string in runes.
func RuneCountInString(s string) RuneNumber {
	return RuneNumber(utf8.RuneCountInString(s))
}

// Get the length of the byte slice in runes.
func RuneCount(b []byte) RuneNumber {
	return RuneNumber(utf8.RuneCount(b))
}

// Returns the number of horizontal cells needed to print the given
// text. It splits the text into its grapheme clusters, calculates each
// cluster's width, and adds them up to a total.
func GetWidth(text string) Width {
	return Width(uniseg.StringWidth(text))
}

// Returns the number of horizontal cells needed to print the given
// text. It splits the text into its grapheme clusters, calculates each
// cluster's width, and adds them up to a total.
func GraphemeCountInString(text string) GraphemeNumber {
	return GraphemeNumber(uniseg.GraphemeClusterCount(text))
}

// Get the width of the rune (how many columns it takes upt in the terminal).
func GetRuneWidth(char rune) Width {
	return Width(runewidth.RuneWidth(char))
}

// Returns the rune index of the nth grapheme in the given text.
func RuneIndexNthGrapheme(text string, n GraphemeNumber) RuneNumber {
	g := uniseg.NewGraphemes(text)
	var currentGraphemeIndex GraphemeNumber
	var currentPosition RuneNumber

	for g.Next() {
		if currentGraphemeIndex >= n {
			break
		}

		currentPosition += RuneNumber(len(g.Runes()))
		currentGraphemeIndex++
	}
	return currentPosition
}

// Returns the rune index of the nth column (in terms of char width) in the given text.
func RuneIndexNthColumn(text string, n Width) RuneNumber {
	g := uniseg.NewGraphemes(text)
	var currentColumnIndex Width
	var currentPosition RuneNumber
	var previousPosition RuneNumber

	for g.Next() {
		if currentColumnIndex > n {
			currentPosition = previousPosition
			break
		}
		if currentColumnIndex == n {
			break
		}
		previousPosition = currentPosition
		currentPosition += RuneNumber(len(g.Runes()))
		currentColumnIndex += Width(g.Width())
	}
	return currentPosition
}

// IndexNotByte is similar with strings.IndexByte but showing the opposite behavior.
func IndexNotByte(s string, c byte) ByteNumber {
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] != c {
			return ByteNumber(i)
		}
	}
	return -1
}

// LastIndexNotByte is similar with strings.LastIndexByte but showing the opposite behavior.
func LastIndexNotByte(s string, c byte) ByteNumber {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != c {
			return ByteNumber(i)
		}
	}
	return -1
}

type asciiSet [8]uint32

func (as *asciiSet) notContains(c byte) bool {
	return (as[c>>5] & (1 << uint(c&31))) == 0
}

func makeASCIISet(chars string) (as asciiSet, ok bool) {
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c >= utf8.RuneSelf {
			return as, false
		}
		as[c>>5] |= 1 << uint(c&31)
	}
	return as, true
}

// IndexNotAny is similar with strings.IndexAny but showing the opposite behavior.
func IndexNotAny(s, chars string) ByteNumber {
	if len(chars) > 0 {
		if len(s) > 8 {
			if as, isASCII := makeASCIISet(chars); isASCII {
				for i := 0; i < len(s); i++ {
					if as.notContains(s[i]) {
						return ByteNumber(i)
					}
				}
				return -1
			}
		}

	LabelFirstLoop:
		for i, c := range s {
			for j, m := range chars {
				if c != m && j == len(chars)-1 {
					return ByteNumber(i)
				} else if c != m {
					continue
				} else {
					continue LabelFirstLoop
				}
			}
		}
	}
	return -1
}

// LastIndexNotAny is similar with strings.LastIndexAny but showing the opposite behavior.
func LastIndexNotAny(s, chars string) ByteNumber {
	if len(chars) > 0 {
		if len(s) > 8 {
			if as, isASCII := makeASCIISet(chars); isASCII {
				for i := len(s) - 1; i >= 0; i-- {
					if as.notContains(s[i]) {
						return ByteNumber(i)
					}
				}
				return -1
			}
		}
	LabelFirstLoop:
		for i := len(s); i > 0; {
			r, size := utf8.DecodeLastRuneInString(s[:i])
			i -= size
			for j, m := range chars {
				if r != m && j == len(chars)-1 {
					return ByteNumber(i)
				} else if r != m {
					continue
				} else {
					continue LabelFirstLoop
				}
			}
		}
	}
	return -1
}
