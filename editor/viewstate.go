package editor

import "sort"

// ViewState is how a view shows a document: the selections, the scroll
// position and the folds. A view keeps one for each of its documents; it
// follows the edits of the document, also made in other views.
type ViewState struct {
	doc  *Document
	sels Selections

	topLine int
	topSub  int
	scrollX int

	folded LineMarkers
	// hidden are the line ranges (inclusive) hidden by the folds, sorted
	hidden      []Range
	hiddenDirty bool

	view *View
}

// NewViewState makes the state of a view showing the document from its start
func NewViewState(doc *Document) *ViewState {
	s := &ViewState{doc: doc, sels: newSelections()}
	doc.AddListener(s)
	return s
}

// Close stops following the document
func (s *ViewState) Close() {
	s.doc.RemoveListener(s)
}

func (s *ViewState) Document() *Document { return s.doc }

// Caret returns the position of the main caret
func (s *ViewState) Caret() int { return s.sels.MainSel().Caret }

// TopLine returns the first line shown
func (s *ViewState) TopLine() int { return s.topLine }

// SetPosition sets the caret and the first line shown, e.g. restoring a session
func (s *ViewState) SetPosition(caret, anchor, topLine int) {
	n := s.doc.Len()
	s.sels.setSingle(Sel{Anchor: max(0, min(anchor, n)), Caret: max(0, min(caret, n)), WantCol: -1})
	s.topLine = max(0, min(topLine, s.doc.LineCount()-1))
	s.topSub = 0
}

// Anchor returns the anchor of the main selection
func (s *ViewState) Anchor() int { return s.sels.MainSel().Anchor }

// FoldedLines returns the lines with folded blocks
func (s *ViewState) FoldedLines() []int { return s.folded.Lines() }

// SetFoldedLines folds the blocks of the lines, e.g. restoring a session
func (s *ViewState) SetFoldedLines(lines []int) {
	s.folded.Set(lines)
	s.hiddenDirty = true
}

func (s *ViewState) DocChanged(doc *Document, c *Change) {
	for i := range s.sels.List {
		sel := &s.sels.List[i]
		na, nc := c.MapPos(sel.Anchor), c.MapPos(sel.Caret)
		if na != sel.Anchor {
			sel.AnchorVS = 0
		}
		if nc != sel.Caret {
			sel.CaretVS = 0
		}
		sel.Anchor, sel.Caret = na, nc
	}
	if c.Pos == 0 && c.Removed > 0 && c.Removed >= doc.Len()-c.Inserted+c.Removed {
		// The whole text was replaced
		s.topLine, s.topSub = min(s.topLine, doc.LineCount()-1), 0
	} else {
		top := s.topLine
		s.topLine = c.MapLine(s.topLine)
		if s.topLine != top || c.Line == top {
			s.topSub = 0
		}
	}
	s.topLine = max(0, min(s.topLine, doc.LineCount()-1))
	if s.folded.Len() > 0 {
		s.folded.Adjust(c)
	}
	s.hiddenDirty = true
	if s.view != nil {
		s.view.docChanged(c)
	}
}

// ---- Folding ----

func (s *ViewState) indentFold() bool {
	return s.doc.lang != nil && s.doc.lang.FoldIndent
}

func (s *ViewState) noFold() bool {
	return s.doc.large || s.doc.lang == nil || s.doc.lang.NoFold
}

// lineIndent returns the indentation of the line in columns, -1 for a blank line
func (s *ViewState) lineIndent(line int, tabSize int) int {
	start := s.doc.LineStart(line)
	end := s.doc.LineEnd(line)
	text := s.doc.buf.View(start, min(end, start+4096))
	col := 0
	for _, c := range text {
		switch c {
		case ' ':
			col++
		case '\t':
			col += tabSize - col%tabSize
		default:
			return col
		}
	}
	if end-start > 4096 {
		return col
	}
	return -1
}

// IsFoldHeader reports whether a block that can be folded starts at the line
func (s *ViewState) IsFoldHeader(line int, tabSize int) bool {
	if s.noFold() {
		return false
	}
	if s.indentFold() {
		ind := s.lineIndent(line, tabSize)
		if ind < 0 {
			return false
		}
		last := min(s.doc.LineCount()-1, line+1000)
		for l := line + 1; l <= last; l++ {
			ni := s.lineIndent(l, tabSize)
			if ni >= 0 {
				return ni > ind
			}
		}
		return false
	}
	lvl := s.doc.hl.Level(line)
	f := s.doc.hl.fold(line)
	return lvl+f.end() > lvl+f.min()
}

// foldBlockEnd returns the last line hidden when the block of the header
// line is folded, -1 if the line is not a header
func (s *ViewState) foldBlockEnd(line int, tabSize int) int {
	if s.noFold() {
		return -1
	}
	last := s.doc.LineCount() - 1
	if s.indentFold() {
		ind := s.lineIndent(line, tabSize)
		if ind < 0 {
			return -1
		}
		end := -1
		for l := line + 1; l <= last; l++ {
			ni := s.lineIndent(l, tabSize)
			if ni < 0 {
				continue
			}
			if ni <= ind {
				break
			}
			end = l
		}
		return end
	}
	h := s.doc.hl
	lvl := h.Level(line)
	f := h.fold(line)
	lowest, target := lvl+f.min(), lvl+f.end()
	if target <= lowest {
		return -1
	}
	cur := target
	for l := line + 1; l <= last; l++ {
		fl := h.fold(l)
		if cur+fl.min() < target {
			// The closing line: kept shown when it opens the next block ("} else {")
			if cur+fl.end() >= target {
				if l-1 <= line {
					return -1
				}
				return l - 1
			}
			return l
		}
		cur += fl.end()
	}
	return last
}

// headerOf returns the header of the innermost block holding the line (the
// line itself if it is a header), -1 if none
func (s *ViewState) headerOf(line int, tabSize int) int {
	if s.noFold() {
		return -1
	}
	if s.IsFoldHeader(line, tabSize) {
		return line
	}
	limit := max(0, line-100000)
	if s.indentFold() {
		ind := s.lineIndent(line, tabSize)
		for l := line - 1; l >= limit; l-- {
			li := s.lineIndent(l, tabSize)
			if li < 0 {
				continue
			}
			if ind < 0 || li < ind {
				if s.IsFoldHeader(l, tabSize) {
					return l
				}
				ind = li
			}
		}
		return -1
	}
	h := s.doc.hl
	m := h.Level(line)
	lvl := m
	for l := line - 1; l >= limit; l-- {
		f := h.fold(l)
		lvl -= f.end() // the level at the start of l
		end, lowest := lvl+f.end(), lvl+f.min()
		if end > lowest && end <= m {
			return l
		}
		m = min(m, lowest)
	}
	return -1
}

func (s *ViewState) rebuildHidden(tabSize int) {
	if !s.hiddenDirty {
		return
	}
	s.hiddenDirty = false
	s.hidden = s.hidden[:0]
	if s.folded.Len() == 0 {
		return
	}
	var keep []int
	for _, h := range s.folded.lines {
		if h >= s.doc.LineCount() {
			continue
		}
		e := s.foldBlockEnd(h, tabSize)
		if e <= h {
			continue // not a block any more
		}
		keep = append(keep, h)
		n := len(s.hidden)
		if n > 0 && s.hidden[n-1].End >= h {
			// Inside a folded block
			if e > s.hidden[n-1].End {
				s.hidden[n-1].End = e
			}
			continue
		}
		s.hidden = append(s.hidden, Range{h + 1, e})
	}
	s.folded.lines = keep
}

// isHidden reports whether the line is hidden in a folded block
func (s *ViewState) isHidden(line int) bool {
	_, ok := s.hiddenRangeOf(line)
	return ok
}

func (s *ViewState) hiddenRangeOf(line int) (Range, bool) {
	i := sort.Search(len(s.hidden), func(i int) bool { return s.hidden[i].End >= line })
	if i < len(s.hidden) && s.hidden[i].Start <= line {
		return s.hidden[i], true
	}
	return Range{}, false
}

// nextVisible returns the first shown line after the line, -1 if none
func (s *ViewState) nextVisible(line int) int {
	l := line + 1
	if r, ok := s.hiddenRangeOf(l); ok {
		l = r.End + 1
	}
	if l >= s.doc.LineCount() {
		return -1
	}
	return l
}

// prevVisible returns the last shown line before the line, -1 if none
func (s *ViewState) prevVisible(line int) int {
	l := line - 1
	if r, ok := s.hiddenRangeOf(l); ok {
		l = r.Start - 1
	}
	return l
}

// visibleLine returns the line or, if it is hidden, the header of its block
func (s *ViewState) visibleLine(line int) int {
	if r, ok := s.hiddenRangeOf(line); ok {
		return r.Start - 1
	}
	return line
}

// hiddenBefore returns the number of hidden lines before the line
func (s *ViewState) hiddenBefore(line int) int {
	n := 0
	for _, r := range s.hidden {
		if r.Start >= line {
			break
		}
		n += min(r.End, line-1) - r.Start + 1
	}
	return n
}

// visibleCount returns the number of shown lines
func (s *ViewState) visibleCount() int {
	return s.doc.LineCount() - s.hiddenBefore(s.doc.LineCount())
}

// lineOfVisibleIndex returns the line shown as the k-th one
func (s *ViewState) lineOfVisibleIndex(k int) int {
	line := k
	for _, r := range s.hidden {
		if r.Start > line {
			break
		}
		line += r.End - r.Start + 1
	}
	return max(0, min(line, s.doc.LineCount()-1))
}

// unfoldLine shows the line, unfolding the blocks that hide it
func (s *ViewState) unfoldLine(line int, tabSize int) bool {
	s.rebuildHidden(tabSize)
	if !s.isHidden(line) {
		return false
	}
	for _, h := range s.folded.Lines() {
		if h < line {
			if e := s.foldBlockEnd(h, tabSize); e >= line {
				s.folded.Remove(h)
			}
		}
	}
	s.hiddenDirty = true
	s.rebuildHidden(tabSize)
	return true
}

// visibleIndexOf returns the index of the line among the shown lines
func (s *ViewState) visibleIndexOf(line int) int {
	return line - s.hiddenBefore(line)
}
