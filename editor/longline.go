package editor

import (
	"sort"
	"unicode/utf8"
)

// A line longer than longLine is never read whole to be shown or to move
// the caret in it: its columns are found from checkpoints taken every
// cpStep bytes, and only the part on the screen is read. Such a line has no
// syntax highlighting, and wrapped it is cut at the width of the view.
const (
	longLine = 64 << 10
	cpStep   = 16 << 10
)

// lineLayout holds the checkpoints of a long line: the offset (from the
// line start) and the column of every few kilobytes
type lineLayout struct {
	start  int
	length int
	cpOff  []int
	cpCol  []int
	cols   int
}

// lineLen returns the length of the line in bytes
func (v *View) lineLen(line int) int {
	return v.doc.LineEnd(line) - v.doc.LineStart(line)
}

func (v *View) isLong(line int) bool {
	return v.lineLen(line) > longLine
}

// lineHead returns up to n bytes from the start of the line
func (v *View) lineHead(line, n int) []byte {
	ls := v.doc.LineStart(line)
	return v.doc.buf.View(ls, min(v.doc.LineEnd(line), ls+n))
}

// layoutOf returns the checkpoints of a long line, making them once per version of the text
func (v *View) layoutOf(line int) *lineLayout {
	if v.longVersion != v.doc.version {
		v.longVersion = v.doc.version
		v.longCache = make(map[int]*lineLayout)
	}
	if l, ok := v.longCache[line]; ok {
		return l
	}
	tab := v.tabSize()
	start := v.doc.LineStart(line)
	end := v.doc.LineEnd(line)
	l := &lineLayout{start: start, length: end - start}
	col := 0
	rel := 0
	next := 0
	var carry []byte
	v.doc.buf.ForEachPiece(start, end, func(p []byte) bool {
		if len(carry) > 0 {
			// A character split between the pieces
			need := utf8.UTFMax - len(carry)
			joined := append(carry, p[:min(need, len(p))]...)
			r, size := utf8.DecodeRune(joined)
			if !utf8.FullRune(joined) && len(p) < need {
				carry = joined
				return true
			}
			col += runeCells(r, size, col, tab)
			used := size - len(carry)
			rel += size
			p = p[used:]
			carry = nil
		}
		i := 0
		for i < len(p) {
			if rel >= next {
				l.cpOff = append(l.cpOff, rel)
				l.cpCol = append(l.cpCol, col)
				next = rel + cpStep
			}
			c := p[i]
			if c < utf8.RuneSelf {
				switch {
				case c == '\t':
					col += tab - col%tab
				case c < 32 || c == 0x7F:
					col += runeCells(rune(c), 1, col, tab)
				default:
					col++
				}
				i++
				rel++
				continue
			}
			if !utf8.FullRune(p[i:]) {
				carry = append([]byte(nil), p[i:]...)
				break
			}
			r, size := utf8.DecodeRune(p[i:])
			col += runeCells(r, size, col, tab)
			i += size
			rel += size
		}
		return true
	})
	if len(carry) > 0 {
		col += 3 * len(carry)
	}
	if len(l.cpOff) == 0 {
		l.cpOff, l.cpCol = []int{0}, []int{0}
	}
	l.cols = col
	if len(v.longCache) > 16 {
		v.longCache = make(map[int]*lineLayout)
	}
	v.longCache[line] = l
	return l
}

// colAt returns the column of the offset rel of the line
func (v *View) colAt(line, rel int) int {
	if !v.isLong(line) {
		return colOfOffset(v.lineText(line), rel, v.tabSize())
	}
	l := v.layoutOf(line)
	rel = max(0, min(rel, l.length))
	k := sort.Search(len(l.cpOff), func(i int) bool { return l.cpOff[i] > rel }) - 1
	k = max(k, 0)
	seg := v.doc.buf.View(l.start+l.cpOff[k], l.start+rel)
	return colAdvance(seg, len(seg), l.cpCol[k], v.tabSize())
}

// offAt returns the offset of the line at the column col, and the columns
// beyond the end of the line
func (v *View) offAt(line, col int, nearest bool) (int, int) {
	if !v.isLong(line) {
		return offsetOfCol(v.lineText(line), col, v.tabSize(), nearest)
	}
	l := v.layoutOf(line)
	if col >= l.cols {
		return l.length, col - l.cols
	}
	k := sort.Search(len(l.cpCol), func(i int) bool { return l.cpCol[i] > col }) - 1
	k = max(k, 0)
	end := l.length
	if k+1 < len(l.cpOff) {
		end = l.cpOff[k+1] + utf8.UTFMax
	}
	seg := v.doc.buf.View(l.start+l.cpOff[k], l.start+min(end, l.length))
	off, beyond := offsetOfColFrom(seg, col, l.cpCol[k], v.tabSize(), nearest)
	return l.cpOff[k] + off, beyond
}

// lineColsOf returns the width of the line in columns
func (v *View) lineColsOf(line int) int {
	if !v.isLong(line) {
		return lineCols(v.lineText(line), v.tabSize())
	}
	return v.layoutOf(line).cols
}
