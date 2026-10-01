package editor

// Lexer styles the text of a language line by line. The state carries what
// a line leaves open for the next one, e.g. a block comment.
type Lexer interface {
	// LexLine styles text (a line without its line break) that starts in
	// state, and returns the state at the start of the next line
	LexLine(text []byte, state uint32, out *LineStyles) uint32
}

// StyleRun: the text from Start to the start of the next run has Style
type StyleRun struct {
	Start int32
	Style Style
}

// LineStyles collects the styles of a line and its folding: how many blocks
// (braces, tags...) the line opens and closes
type LineStyles struct {
	Runs []StyleRun

	depth    int // blocks opened minus closed so far
	minDepth int // the lowest depth reached in the line
}

func (ls *LineStyles) reset() {
	ls.Runs = ls.Runs[:0]
	ls.depth = 0
	ls.minDepth = 0
}

// Set styles the text from pos on with st
func (ls *LineStyles) Set(pos int, st Style) {
	n := len(ls.Runs)
	if n > 0 {
		last := &ls.Runs[n-1]
		if last.Style == st {
			return
		}
		if int(last.Start) >= pos {
			last.Style = st
			// Joins with the run before it when the same
			if n > 1 && ls.Runs[n-2].Style == st {
				ls.Runs = ls.Runs[:n-1]
			}
			return
		}
	}
	ls.Runs = append(ls.Runs, StyleRun{int32(pos), st})
}

// Span styles text from start to end with st; after it the style is back
// to what it was before
func (ls *LineStyles) Span(start, end int, st Style) {
	ls.Set(start, st)
	ls.Set(end, StyleDefault)
}

// Open counts a block opened on the line, e.g. "{"
func (ls *LineStyles) Open() { ls.depth++ }

// Close counts a block closed on the line, e.g. "}"
func (ls *LineStyles) Close() {
	ls.depth--
	if ls.depth < ls.minDepth {
		ls.minDepth = ls.depth
	}
}

// StyleAt returns the style of the byte at pos of the line
func (ls *LineStyles) StyleAt(pos int) Style {
	st := StyleDefault
	for _, r := range ls.Runs {
		if int(r.Start) > pos {
			break
		}
		st = r.Style
	}
	return st
}

// foldInfo is packed: the depth at the end of the line and the lowest one
// within it, both relative to its start
type foldInfo int32

func packFold(minDepth, endDepth int) foldInfo {
	minDepth = max(-32768, min(minDepth, 0))
	endDepth = max(-32768, min(endDepth, 32767))
	return foldInfo(uint32(uint16(int16(minDepth)))<<16 | uint32(uint16(int16(endDepth))))
}

func (f foldInfo) min() int { return int(int16(uint16(uint32(f) >> 16))) }
func (f foldInfo) end() int { return int(int16(uint16(uint32(f)))) }

// plainLexer leaves the text as it is
type plainLexer struct{}

func (plainLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 { return 0 }

// highlighter keeps the lexer states at the line starts of a document, so
// only the lines that changed are styled again
type highlighter struct {
	doc   *Document
	lexer Lexer
	lang  *Language

	// states[i] is the state at the start of line i, folds[i] the folding of
	// line i (known when states[i+1] is)
	states []uint32
	folds  []foldInfo
	// states before valid are right; those from valid to stale were right
	// before the last edits and are right again once a line gets its old state
	valid int
	stale int

	// levels[k] is the fold depth at the start of line k*levelStep; valid for k < levelsValid
	levels      []int32
	levelsValid int

	disabled bool
	scratch  LineStyles
}

const levelStep = 1024

func newHighlighter(d *Document) *highlighter {
	h := &highlighter{doc: d, lexer: plainLexer{}}
	h.reset()
	return h
}

func (h *highlighter) reset() {
	h.states = h.states[:0]
	h.states = append(h.states, 0)
	h.folds = h.folds[:0]
	h.valid = 1
	h.stale = 1
	h.levels = h.levels[:0]
	h.levelsValid = 0
}

func (h *highlighter) setLanguage(lang *Language) {
	h.lang = lang
	if lang == nil || lang.lexer == nil {
		h.lexer = plainLexer{}
	} else {
		h.lexer = lang.lexer
	}
	h.reset()
}

// unknownState marks the states of added lines; no lexer returns it
const unknownState = 0xFFFFFFFF

// changed updates the states after an edit: the lines from the edited one
// on must be styled again
func (h *highlighter) changed(c *Change) {
	line := c.Line
	h.levelsValid = min(h.levelsValid, line/levelStep+1)
	if line+1 < len(h.states) && c.LinesRemoved != c.LinesAdded {
		// Remove the states of the removed lines, add unknown ones for the added lines
		from := line + 1
		to := min(line+1+c.LinesRemoved, len(h.states))
		h.states = spliceStates(h.states, from, to, c.LinesAdded)
		if from < len(h.folds) {
			h.folds = spliceFolds(h.folds, from, min(to, len(h.folds)), c.LinesAdded)
		}
		if h.stale > line+1+c.LinesRemoved {
			h.stale += c.LinesAdded - c.LinesRemoved
		} else if h.stale > line+1 {
			h.stale = line + 1
		}
	}
	// The state after the edited line is unknown: a later convergence must
	// not jump over this line
	if line+1 < len(h.states) {
		h.states[line+1] = unknownState
	}
	h.valid = min(h.valid, line+1)
	h.stale = max(min(h.stale, len(h.states)), h.valid)
}

func spliceStates(s []uint32, from, to, add int) []uint32 {
	if from > len(s) {
		return s
	}
	out := make([]uint32, 0, len(s)-(to-from)+add)
	out = append(out, s[:from]...)
	for i := 0; i < add; i++ {
		out = append(out, unknownState)
	}
	return append(out, s[to:]...)
}

func spliceFolds(s []foldInfo, from, to, add int) []foldInfo {
	out := make([]foldInfo, 0, len(s)-(to-from)+add)
	out = append(out, s[:from]...)
	for i := 0; i < add; i++ {
		out = append(out, 0)
	}
	return append(out, s[to:]...)
}

// ensure makes the states known up to the start of line upTo
func (h *highlighter) ensure(upTo int) {
	lines := h.doc.LineCount()
	upTo = min(upTo, lines)
	if upTo < h.valid {
		return
	}
	if len(h.states) < lines+1 {
		grow := make([]uint32, lines+1)
		n := copy(grow, h.states)
		for i := n; i < len(grow); i++ {
			grow[i] = unknownState
		}
		h.states = grow
	}
	if len(h.folds) < lines {
		grow := make([]foldInfo, lines)
		copy(grow, h.folds)
		h.folds = grow
	}
	if h.disabled {
		for i := h.valid; i <= upTo; i++ {
			h.states[i] = 0
			if i > 0 {
				h.folds[i-1] = 0
			}
		}
		h.valid = upTo + 1
		return
	}
	b := h.doc.buf
	for h.valid <= upTo {
		start := h.valid - 1
		state := h.states[start]
		b.forEachLine(start, func(line int, text []byte) bool {
			h.scratch.reset()
			next := h.lexer.LexLine(text, state, &h.scratch)
			h.folds[line] = packFold(h.scratch.minDepth, h.scratch.depth)
			// Back to the state it had before the edit: the states after it
			// are right up to the lines added since
			if line+1 < h.stale && h.states[line+1] == next {
				j := line + 1
				for j < h.stale && h.states[j] != unknownState {
					j++
				}
				h.valid = j
				return false
			}
			h.states[line+1] = next
			h.valid = line + 2
			state = next
			return line+1 < upTo
		})
		if h.valid > lines {
			break
		}
	}
	// The states after the ones just made are from an older pass: whether
	// they follow from the new ones is not known, so a later convergence
	// must not jump over here
	if h.valid < h.stale {
		h.states[h.valid] = unknownState
	}
	if h.valid > h.stale {
		h.stale = h.valid
	}
}

// StartState returns the lexer state at the start of the line
func (h *highlighter) StartState(line int) uint32 {
	h.ensure(line)
	if line < len(h.states) {
		return h.states[line]
	}
	return 0
}

// Styles styles the line; the result is valid until the next call
func (h *highlighter) Styles(line int, text []byte, out *LineStyles) {
	out.reset()
	if h.disabled {
		return
	}
	state := h.StartState(line)
	h.lexer.LexLine(text, state, out)
}

// fold returns the folding of the line
func (h *highlighter) fold(line int) foldInfo {
	h.ensure(line + 1)
	if line < len(h.folds) {
		return h.folds[line]
	}
	return 0
}

// Level returns the fold depth at the start of the line
func (h *highlighter) Level(line int) int {
	if line <= 0 {
		return 0
	}
	h.ensure(line)
	k := line / levelStep
	if len(h.levels) < k+1 {
		grow := make([]int32, k+1)
		copy(grow, h.levels)
		h.levels = grow
	}
	if h.levelsValid == 0 {
		h.levels[0] = 0
		h.levelsValid = 1
	}
	for h.levelsValid <= k {
		i := h.levelsValid
		lvl := int(h.levels[i-1])
		for l := (i - 1) * levelStep; l < i*levelStep; l++ {
			lvl += h.folds[l].end()
		}
		h.levels[i] = int32(lvl)
		h.levelsValid++
	}
	lvl := int(h.levels[k])
	for l := k * levelStep; l < line; l++ {
		lvl += h.folds[l].end()
	}
	return lvl
}
