package editor

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/width"
)

// A line is laid out in cells of the width of a character: a tab takes the
// cells up to the next tab stop, a wide (e.g. Chinese) character takes two,
// a control character is shown as its name in a box, like NUL.

var controlNames = [32]string{
	"NUL", "SOH", "STX", "ETX", "EOT", "ENQ", "ACK", "BEL", "BS", "HT", "LF", "VT", "FF", "CR", "SO", "SI",
	"DLE", "DC1", "DC2", "DC3", "DC4", "NAK", "SYN", "ETB", "CAN", "EM", "SUB", "ESC", "FS", "GS", "RS", "US",
}

// controlName returns the name shown for a control character or an invalid
// byte, "" for a usual character
func controlName(r rune, size int, b byte) string {
	if r == utf8.RuneError && size == 1 {
		return fmt.Sprintf("x%02X", b)
	}
	if r < 32 && r != '\t' {
		return controlNames[r]
	}
	if r == 0x7F {
		return "DEL"
	}
	if r >= 0x80 && r < 0xA0 {
		return fmt.Sprintf("x%02X", r)
	}
	return ""
}

// runeCells returns the number of cells of a character at the column col
func runeCells(r rune, size int, col, tabSize int) int {
	switch {
	case r == '\t':
		return tabSize - col%tabSize
	case r < 32 || r == 0x7F || (r == utf8.RuneError && size == 1) || (r >= 0x80 && r < 0xA0):
		if r == utf8.RuneError && size == 1 {
			return 3
		}
		if r < 32 {
			return len(controlNames[r])
		}
		return 3
	case r < 0x300:
		return 1
	case isZeroWidth(r):
		return 0
	case isWide(r):
		return 2
	}
	return 1
}

func isZeroWidth(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF ||
		(r >= 0xFE00 && r <= 0xFE0F)
}

func isWide(r rune) bool {
	if r < 0x1100 {
		return false
	}
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return true
	}
	// Emoji
	return r >= 0x1F300 && r <= 0x1FAFF
}

// colOfOffset returns the column of the byte offset off of the line text
func colOfOffset(text []byte, off, tabSize int) int {
	return colAdvance(text, off, 0, tabSize)
}

// colAdvance returns the column reached at the offset off of text, a part
// of a line starting at the column startCol
func colAdvance(text []byte, off, startCol, tabSize int) int {
	col := startCol
	i := 0
	for i < off && i < len(text) {
		c := text[i]
		if c < utf8.RuneSelf {
			if c == '\t' {
				col += tabSize - col%tabSize
			} else if c < 32 || c == 0x7F {
				col += runeCells(rune(c), 1, col, tabSize)
			} else {
				col++
			}
			i++
			continue
		}
		r, size := utf8.DecodeRune(text[i:])
		col += runeCells(r, size, col, tabSize)
		i += size
	}
	return col
}

// offsetOfCol returns the offset of the character at the column col of the
// line text; a column inside a wide character or a tab goes to its start,
// or to its end when nearest is set and the column is in its second half.
// Past the end of the line it returns the end and the columns beyond it.
func offsetOfCol(text []byte, col, tabSize int, nearest bool) (off int, beyond int) {
	return offsetOfColFrom(text, col, 0, tabSize, nearest)
}

// offsetOfColFrom is offsetOfCol for text, a part of a line starting at the column startCol
func offsetOfColFrom(text []byte, col, startCol, tabSize int, nearest bool) (off int, beyond int) {
	c := startCol
	i := 0
	for i < len(text) {
		r, size := rune(text[i]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRune(text[i:])
		}
		w := runeCells(r, size, c, tabSize)
		if c+w > col {
			if nearest && col-c >= (w+1)/2 && w > 0 {
				return i + size, 0
			}
			return i, 0
		}
		c += w
		i += size
	}
	return len(text), max(0, col-c)
}

// lineCols returns the width of the line in columns
func lineCols(text []byte, tabSize int) int {
	return colOfOffset(text, len(text), tabSize)
}

// wrapLine returns the offsets of the line text where its wrapped rows
// start, after the first one (nil if the line fits); it breaks after spaces
// when it can. indent is the indentation of the rows after the first one in columns.
func wrapLine(text []byte, cols, tabSize int, indentWrapped bool) (breaks []int, indent int) {
	if cols < 8 {
		cols = 8
	}
	if len(text) <= cols && !containsTabOrWide(text) {
		return nil, 0
	}
	if indentWrapped {
		// The rows after the first are indented like the line
		i := 0
		for i < len(text) && (text[i] == ' ' || text[i] == '\t') {
			i++
		}
		indent = colOfOffset(text, i, tabSize)
		if indent > cols/2 {
			indent = 0
		}
	}
	rowStart := 0
	rowCol := 0 // column of the row start in the line
	col := 0
	lastBreak := -1 // offset after the last space in the row
	lastBreakCol := 0
	avail := cols
	i := 0
	for i < len(text) {
		r, size := rune(text[i]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRune(text[i:])
		}
		w := runeCells(r, size, col, tabSize)
		// Spaces stay at the end of the row even past the edge: a row does not start with them
		if col+w-rowCol > avail && i > rowStart && r != ' ' && r != '\t' {
			b := i
			bcol := col
			if lastBreak > rowStart {
				b = lastBreak
				bcol = lastBreakCol
			}
			breaks = append(breaks, b)
			rowStart = b
			rowCol = bcol
			avail = cols - indent
			lastBreak = -1
			if b != i {
				continue
			}
		}
		col += w
		i += size
		if r == ' ' || r == '\t' {
			lastBreak = i
			lastBreakCol = col
		}
		// Keep zero width characters with the one before them
		for i < len(text) {
			r2, s2 := utf8.DecodeRune(text[i:])
			if !isZeroWidth(r2) || r2 < 0x300 {
				break
			}
			i += s2
		}
	}
	return breaks, indent
}

func containsTabOrWide(text []byte) bool {
	for _, c := range text {
		if c == '\t' || c >= 0x80 || c < 32 {
			return true
		}
	}
	return false
}

// Classes of characters for the word moves and the double click selection
const (
	ccSpace = iota
	ccWord
	ccPunct
	ccLineEnd
)

func charClass(r rune, extraWordChars string) int {
	switch {
	case r == '\n':
		return ccLineEnd
	case r == ' ' || r == '\t' || r == 0xA0 || r == 0x3000:
		return ccSpace
	case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.Is(unicode.Mn, r):
		return ccWord
	}
	for _, c := range extraWordChars {
		if c == r {
			return ccWord
		}
	}
	return ccPunct
}
