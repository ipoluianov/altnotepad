package editor

import "sort"

// Sel is a selection: the text from Anchor to Caret; empty when they are
// equal, then it is just the caret. AnchorVS and CaretVS are the columns of
// virtual space after the end of the line (in a rectangular selection).
type Sel struct {
	Anchor, Caret     int
	AnchorVS, CaretVS int
	// WantCol is the column the caret goes to when moved up or down; -1 - its own
	WantCol int
}

func caretSel(pos int) Sel { return Sel{Anchor: pos, Caret: pos, WantCol: -1} }

// Start and End return the selected text range
func (s Sel) Start() int { return min(s.Anchor, s.Caret) }
func (s Sel) End() int   { return max(s.Anchor, s.Caret) }
func (s Sel) Empty() bool {
	return s.Anchor == s.Caret && s.AnchorVS == s.CaretVS
}

// IsEmptyText: no text is selected (virtual space may be)
func (s Sel) IsEmptyText() bool { return s.Anchor == s.Caret }

// Selections are the selections of a view: one usually, several with
// multi-editing or a rectangular selection
type Selections struct {
	List []Sel
	Main int
	// Rect: the selections are the rows of a rectangular selection from
	// (RectAnchorLine, RectAnchorCol) to (RectCaretLine, RectCaretCol)
	Rect           bool
	RectAnchorLine int
	RectAnchorCol  int
	RectCaretLine  int
	RectCaretCol   int
}

func newSelections() Selections {
	return Selections{List: []Sel{caretSel(0)}}
}

func (s *Selections) MainSel() Sel {
	if s.Main >= len(s.List) {
		s.Main = 0
	}
	return s.List[s.Main]
}

func (s *Selections) clone() Selections {
	c := *s
	c.List = append([]Sel(nil), s.List...)
	return c
}

// normalize sorts the selections and joins the ones that overlap
func (s *Selections) normalize() {
	if len(s.List) <= 1 {
		s.Main = 0
		return
	}
	main := s.List[s.Main]
	sort.SliceStable(s.List, func(i, j int) bool {
		a, b := s.List[i], s.List[j]
		if a.Start() != b.Start() {
			return a.Start() < b.Start()
		}
		return a.End() < b.End()
	})
	out := s.List[:1]
	for _, sel := range s.List[1:] {
		last := &out[len(out)-1]
		overlap := sel.Start() < last.End() || (sel.Start() == last.End() && (sel.IsEmptyText() || last.IsEmptyText()) && !s.Rect)
		if sel.Start() == last.Start() && sel.End() == last.End() && sel.CaretVS == last.CaretVS {
			overlap = true
		}
		if s.Rect && sel.Start() == last.Start() && sel.CaretVS != last.CaretVS {
			overlap = false
		}
		if overlap {
			start := min(last.Start(), sel.Start())
			end := max(last.End(), sel.End())
			forward := last.Caret >= last.Anchor
			if forward {
				last.Anchor, last.Caret = start, end
			} else {
				last.Anchor, last.Caret = end, start
			}
			last.CaretVS = max(last.CaretVS, sel.CaretVS)
			continue
		}
		out = append(out, sel)
	}
	s.List = out
	s.Main = 0
	for i, sel := range s.List {
		if sel.Start() <= main.Caret && main.Caret <= sel.End() {
			s.Main = i
			break
		}
	}
}

// setSingle leaves one selection
func (s *Selections) setSingle(sel Sel) {
	s.List = append(s.List[:0], sel)
	s.Main = 0
	s.Rect = false
}
