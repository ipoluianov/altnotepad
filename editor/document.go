package editor

import (
	"bytes"
	"os"
	"sort"
	"time"
)

// EOL is the line break of a file. Inside a document lines always end with
// "\n"; the EOL is used when the file is saved.
type EOL int

const (
	EOLCRLF EOL = iota // Windows
	EOLLF              // Unix
	EOLCR              // old Mac
)

func (e EOL) String() string {
	switch e {
	case EOLCRLF:
		return "Windows (CR LF)"
	case EOLCR:
		return "Macintosh (CR)"
	}
	return "Unix (LF)"
}

// Bytes returns the line break as saved to a file
func (e EOL) Bytes() []byte {
	switch e {
	case EOLCRLF:
		return []byte("\r\n")
	case EOLCR:
		return []byte("\r")
	}
	return []byte("\n")
}

// Change describes one edit of a document: Removed bytes at Pos were
// replaced with Inserted bytes. Line is the line of Pos before the edit.
type Change struct {
	Pos          int
	Removed      int
	Inserted     int
	Line         int
	LinesRemoved int
	LinesAdded   int
	// AtLineStart: Pos was at the start of Line
	AtLineStart bool

	// multi describes a change made of many edits (see ApplyEdits)
	multi *multiMap
}

// multiMap maps the positions and lines through many edits at once
type multiMap struct {
	pos, end   []int // old positions of the edits
	newPos     []int
	delta      []int // total length change of the edits before each one
	line, lend []int // old lines of the edit starts and ends
	ldelta     []int // total line count change of the edits before each one
}

func (m *multiMap) mapPos(p int) int {
	n := len(m.pos)
	k := sort.Search(n, func(i int) bool { return !(m.pos[i] < p && m.end[i] <= p) })
	if k < n && m.pos[k] < p && p < m.end[k] {
		return m.newPos[k]
	}
	return p + m.delta[k]
}

func (m *multiMap) mapLine(l int) int {
	n := len(m.line)
	k := sort.Search(n, func(i int) bool { return m.lend[i] >= l })
	if k < n && m.line[k] < l && l <= m.lend[k] {
		return m.line[k] + m.ldelta[k]
	}
	return l + m.ldelta[k]
}

// MapPos returns where a position before the change is after it
func (c *Change) MapPos(p int) int {
	if c.multi != nil {
		return c.multi.mapPos(p)
	}
	if p <= c.Pos {
		return p
	}
	if p < c.Pos+c.Removed {
		return c.Pos
	}
	return p - c.Removed + c.Inserted
}

// MapLine returns where the line before the change is after it
func (c *Change) MapLine(l int) int {
	if c.multi != nil {
		return c.multi.mapLine(l)
	}
	if l < c.Line {
		return l
	}
	if l == c.Line {
		if c.AtLineStart && c.LinesRemoved == 0 && c.LinesAdded > 0 {
			return l + c.LinesAdded
		}
		return l
	}
	if l <= c.Line+c.LinesRemoved {
		return c.Line
	}
	return l - c.LinesRemoved + c.LinesAdded
}

// DocListener is told about the changes of a document, e.g. a view showing it
type DocListener interface {
	DocChanged(doc *Document, c *Change)
}

// editOp is one change in the undo history
type editOp struct {
	pos int
	del []byte
	ins []byte
	// multi: many edits made at once, positions before them (see ApplyEdits)
	multi []Edit
}

// Edit replaces Len bytes at Pos with Text; Old is the replaced text
type Edit struct {
	Pos  int
	Len  int
	Text []byte
	Old  []byte
}

// undoGroup is what one undo step reverts
type undoGroup struct {
	ops []editOp
	// Selections of the view that made the change, before and after it
	selBefore, selAfter any
	typing              bool
	time                time.Time
}

// Document is a text with its undo history, file information, bookmarks
// and syntax highlighting state. It can be shown in several views at once.
type Document struct {
	buf *Buffer

	undo     []*undoGroup
	undoPos  int // groups before it are done, after it are undone
	savePos  int // undoPos when saved; -1 when that state is lost
	current  *undoGroup
	depth    int
	noUndo   bool
	inUndo   bool
	coalesce bool

	version int64

	listeners []DocListener

	// File
	Path     string
	EOL      EOL
	Encoding Encoding
	ModTime  time.Time
	FileSize int64
	ReadOnly bool
	// FileReadOnly: the file itself cannot be written
	FileReadOnly bool

	Bookmarks LineMarkers
	// Changed are the lines changed since the file was opened and not saved;
	// Saved are the ones changed and saved since (the change history)
	Changed LineMarkers
	Saved   LineMarkers
	// Marks are the ranges marked by the search, per mark style
	Marks [NumMarkStyles]RangeSet

	hl   *highlighter
	lang *Language
	// large: the file is above the large file limit: no highlighting, folding or wrapping
	large bool

	// OnModifiedChanged is called when the document becomes modified or
	// unmodified (saved, or undone to the saved state)
	OnModifiedChanged func()
	lastModified      bool
}

func NewDocument() *Document {
	d := &Document{buf: NewBuffer(), EOL: DefaultEOL(), Encoding: EncodingUTF8, lang: LangText}
	d.hl = newHighlighter(d)
	return d
}

// Language returns the language of the syntax highlighting
func (d *Document) Language() *Language { return d.lang }

// SetLanguage sets the language of the syntax highlighting
func (d *Document) SetLanguage(l *Language) {
	if l == nil {
		l = LangText
	}
	d.lang = l
	d.hl.setLanguage(l)
	d.hl.disabled = d.large
	d.version++
}

// IsLarge reports whether the document is over the large file limit, see SetLarge
func (d *Document) IsLarge() bool { return d.large }

// SetLarge turns off the syntax highlighting and folding of a large file
func (d *Document) SetLarge(large bool) {
	d.large = large
	d.hl.disabled = large
	d.hl.reset()
	d.version++
}

// LineStyles styles the line with the syntax highlighting of the language
func (d *Document) LineStyles(line int, text []byte, out *LineStyles) {
	d.hl.Styles(line, text, out)
}

// FoldLevel returns the fold depth at the start of the line
func (d *Document) FoldLevel(line int) int { return d.hl.Level(line) }

// DefaultEOL is the EOL of the new documents on this system
func DefaultEOL() EOL {
	if os.PathSeparator == '\\' {
		return EOLCRLF
	}
	return EOLLF
}

func (d *Document) Buffer() *Buffer { return d.buf }
func (d *Document) Len() int        { return d.buf.Len() }
func (d *Document) LineCount() int  { return d.buf.LineCount() }
func (d *Document) Version() int64  { return d.version }

func (d *Document) LineStart(line int) int     { return d.buf.LineStart(line) }
func (d *Document) LineEnd(line int) int       { return d.buf.LineEnd(line) }
func (d *Document) LineOfOffset(off int) int   { return d.buf.LineOfOffset(off) }
func (d *Document) Text(start, end int) []byte { return d.buf.Bytes(start, end) }
func (d *Document) String() string             { return d.buf.String() }

// LineText returns a copy of the line without its line break
func (d *Document) LineText(line int) []byte {
	return d.buf.Bytes(d.buf.LineStart(line), d.buf.LineEnd(line))
}

func (d *Document) AddListener(l DocListener) {
	for _, x := range d.listeners {
		if x == l {
			return
		}
	}
	d.listeners = append(d.listeners, l)
}

func (d *Document) RemoveListener(l DocListener) {
	for i, x := range d.listeners {
		if x == l {
			d.listeners = append(d.listeners[:i], d.listeners[i+1:]...)
			return
		}
	}
}

// SetText replaces the whole text; the history is cleared and the document is unmodified
func (d *Document) SetText(text []byte) {
	d.SetBuffer(NewBufferFromBytes(text))
}

// SetBuffer replaces the whole text with the buffer, e.g. a loaded file;
// the history is cleared and the document is unmodified
func (d *Document) SetBuffer(b *Buffer) {
	old := d.buf
	c := &Change{Pos: 0, Removed: old.Len(), Inserted: b.Len(), Line: 0,
		LinesRemoved: old.LineCount() - 1, LinesAdded: b.LineCount() - 1, AtLineStart: true}
	d.buf = b
	d.version++
	d.undo = nil
	d.undoPos = 0
	d.savePos = 0
	d.current = nil
	d.Bookmarks.Clear()
	d.Changed.Clear()
	d.Saved.Clear()
	for i := range d.Marks {
		d.Marks[i].Clear()
	}
	d.hl.reset()
	for _, l := range d.listeners {
		l.DocChanged(d, c)
	}
	d.checkModified()
}

// IsModified reports whether the text differs from the saved one
func (d *Document) IsModified() bool {
	return d.undoPos != d.savePos
}

// SetSavePoint marks the current text as saved
func (d *Document) SetSavePoint() {
	d.savePos = d.undoPos
	if d.Changed.Len() > 0 {
		d.Saved.Set(append(d.Saved.Lines(), d.Changed.lines...))
		d.Changed.Clear()
	}
	d.checkModified()
}

// markChanged adds the lines from..to to the change history
func (d *Document) markChanged(from, to int) {
	if to-from > 1_000_000 {
		return
	}
	if from == to {
		if !d.Changed.Has(from) {
			d.Changed.Add(from)
		}
		d.Saved.Remove(from)
		return
	}
	lines := make([]int, 0, to-from+1)
	for l := from; l <= to; l++ {
		lines = append(lines, l)
	}
	d.addChanged(lines)
}

func (d *Document) addChanged(lines []int) {
	if len(lines) == 0 {
		return
	}
	d.Changed.Set(append(d.Changed.Lines(), lines...))
	if d.Saved.Len() > 0 {
		var keep []int
		for _, l := range d.Saved.lines {
			if !d.Changed.Has(l) {
				keep = append(keep, l)
			}
		}
		d.Saved.lines = keep
	}
}

// SetModified marks the document modified without a change, e.g. after the
// encoding changed
func (d *Document) SetModified() {
	d.savePos = -1
	d.checkModified()
}

func (d *Document) checkModified() {
	m := d.IsModified()
	if m != d.lastModified {
		d.lastModified = m
		if d.OnModifiedChanged != nil {
			d.OnModifiedChanged()
		}
	}
}

// BeginAction starts a group of changes undone at once; calls nest
func (d *Document) BeginAction() {
	if d.depth == 0 {
		d.current = nil
	}
	d.depth++
}

// EndAction ends the group started by BeginAction
func (d *Document) EndAction() {
	if d.depth > 0 {
		d.depth--
	}
	if d.depth == 0 {
		d.current = nil
		d.coalesce = false
	}
}

// SetTyping marks the next change as typing: typed characters join the
// previous typing group, so undo removes a word at a time and not a letter
func (d *Document) SetTyping() {
	d.coalesce = true
}

// SetUndoSelection stores the selections of the view before and after the
// current group of changes, restored by undo and redo
func (d *Document) SetUndoSelection(before, after any) {
	g := d.lastGroup()
	if g == nil {
		return
	}
	if g.selBefore == nil && before != nil {
		g.selBefore = before
	}
	if after != nil {
		g.selAfter = after
	}
}

func (d *Document) lastGroup() *undoGroup {
	if d.undoPos == 0 || d.undoPos > len(d.undo) {
		return nil
	}
	return d.undo[d.undoPos-1]
}

// Insert inserts text at the position
func (d *Document) Insert(pos int, text []byte) {
	d.Replace(pos, 0, text)
}

// Delete removes n bytes at the position
func (d *Document) Delete(pos, n int) {
	d.Replace(pos, n, nil)
}

// Replace replaces n bytes at pos with text
func (d *Document) Replace(pos, n int, text []byte) {
	if d.ReadOnly && !d.inUndo {
		return
	}
	pos = max(0, min(pos, d.buf.Len()))
	n = max(0, min(n, d.buf.Len()-pos))
	if n == 0 && len(text) == 0 {
		return
	}
	var del []byte
	if n > 0 {
		del = d.buf.Bytes(pos, pos+n)
	}
	if n > 0 && bytes.Equal(del, text) {
		return
	}
	d.record(editOp{pos: pos, del: del, ins: append([]byte(nil), text...)})
	d.apply(pos, del, text)
}

func (d *Document) record(op editOp) {
	if d.noUndo || d.inUndo {
		return
	}
	// A new change drops what was undone
	if d.undoPos < len(d.undo) {
		if d.savePos > d.undoPos {
			d.savePos = -1
		}
		d.undo = d.undo[:d.undoPos]
	}
	typing := d.coalesce
	g := d.current
	if g == nil && typing && d.depth <= 1 {
		// Join the previous typing group while typing goes on
		if last := d.lastGroup(); last != nil && last.typing && d.savePos != d.undoPos &&
			time.Since(last.time) < 2*time.Second && canCoalesce(last, op) {
			g = last
		}
	}
	if g == nil {
		g = &undoGroup{typing: typing}
		d.undo = append(d.undo, g)
		d.undoPos = len(d.undo)
	}
	g.time = time.Now()
	g.ops = append(g.ops, op)
	if d.depth > 0 {
		d.current = g
	}
	if d.depth == 0 {
		d.coalesce = false
	}
}

// canCoalesce: typing joins the group when it goes on at the same place
// and does not start a new word after a space
func canCoalesce(g *undoGroup, op editOp) bool {
	if len(g.ops) == 0 || len(op.del) > 0 {
		return false
	}
	last := g.ops[len(g.ops)-1]
	if op.multi != nil || last.multi != nil {
		// Typing with many carets: each keystroke is an edit at every caret
		return op.multi != nil && last.multi != nil && len(op.multi) == len(last.multi) && onlyInserts(op.multi)
	}
	if last.pos+len(last.ins) != op.pos {
		return false
	}
	if len(op.ins) == 1 && len(last.ins) > 0 {
		prev := last.ins[len(last.ins)-1]
		cur := op.ins[0]
		if (prev == ' ' || prev == '\t' || prev == '\n') && cur != ' ' && cur != '\t' && cur != '\n' {
			return false
		}
	}
	return true
}

func (d *Document) apply(pos int, del, ins []byte) {
	line := d.buf.LineOfOffset(pos)
	c := &Change{
		Pos:          pos,
		Removed:      len(del),
		Inserted:     len(ins),
		Line:         line,
		LinesRemoved: bytes.Count(del, nlBytes),
		LinesAdded:   bytes.Count(ins, nlBytes),
		AtLineStart:  pos == d.buf.LineStart(line),
	}
	if len(del) > 0 {
		d.buf.Delete(pos, len(del))
	}
	if len(ins) > 0 {
		d.buf.Insert(pos, ins)
	}
	d.version++
	d.Bookmarks.Adjust(c)
	d.Changed.Adjust(c)
	d.Saved.Adjust(c)
	d.markChanged(c.Line, c.Line+c.LinesAdded)
	for i := range d.Marks {
		d.Marks[i].Adjust(c)
	}
	d.hl.changed(c)
	for _, l := range d.listeners {
		l.DocChanged(d, c)
	}
	d.checkModified()
}

func (d *Document) CanUndo() bool { return d.undoPos > 0 }
func (d *Document) CanRedo() bool { return d.undoPos < len(d.undo) }

// Undo reverts the last group of changes; returns the selections stored
// before it (nil if none) and the position of the first change
func (d *Document) Undo() (sel any, pos int, ok bool) {
	if !d.CanUndo() {
		return nil, 0, false
	}
	d.current = nil
	g := d.undo[d.undoPos-1]
	d.inUndo = true
	for i := len(g.ops) - 1; i >= 0; i-- {
		op := g.ops[i]
		if op.multi != nil {
			inv := invertEdits(op.multi)
			d.applyMulti(inv)
			if i == 0 {
				pos = inv[0].Pos
			}
			continue
		}
		d.apply(op.pos, op.ins, op.del)
		pos = op.pos + len(op.del)
		if i == 0 {
			pos = op.pos
		}
	}
	d.inUndo = false
	d.undoPos--
	d.checkModified()
	return g.selBefore, pos, true
}

// Redo applies again the last undone group; returns the selections stored after it
func (d *Document) Redo() (sel any, pos int, ok bool) {
	if !d.CanRedo() {
		return nil, 0, false
	}
	d.current = nil
	g := d.undo[d.undoPos]
	d.inUndo = true
	for _, op := range g.ops {
		if op.multi != nil {
			d.applyMulti(op.multi)
			last := op.multi[len(op.multi)-1]
			pos = last.Pos + len(last.Text)
			continue
		}
		d.apply(op.pos, op.del, op.ins)
		pos = op.pos + len(op.ins)
	}
	d.inUndo = false
	d.undoPos++
	d.checkModified()
	return g.selAfter, pos, true
}

// ClearUndo forgets the history
func (d *Document) ClearUndo() {
	if d.savePos == d.undoPos {
		d.savePos = 0
	} else {
		d.savePos = -1
	}
	d.undo = nil
	d.undoPos = 0
	d.current = nil
	d.checkModified()
}

// ApplyEdits makes many edits at once, as one undo step: e.g. Replace All
// or typing with many carets. The edits must not overlap; their positions
// are before any of them is made. It takes one pass over the text however
// many edits there are.
func (d *Document) ApplyEdits(edits []Edit) {
	if d.ReadOnly || len(edits) == 0 {
		return
	}
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].Pos < edits[j].Pos })
	out := make([]Edit, 0, len(edits))
	prevEnd := -1
	for _, e := range edits {
		e.Pos = max(0, min(e.Pos, d.buf.Len()))
		e.Len = max(0, min(e.Len, d.buf.Len()-e.Pos))
		if e.Pos < prevEnd {
			continue // overlaps the previous one
		}
		if e.Len == 0 && len(e.Text) == 0 {
			continue
		}
		e.Old = d.buf.Bytes(e.Pos, e.Pos+e.Len)
		if bytes.Equal(e.Old, e.Text) {
			continue
		}
		prevEnd = e.Pos + e.Len
		out = append(out, e)
	}
	if len(out) == 0 {
		return
	}
	if len(out) == 1 {
		d.Replace(out[0].Pos, out[0].Len, out[0].Text)
		return
	}
	d.record(editOp{multi: out})
	d.applyMulti(out)
}

// invertEdits returns the edits that undo the edits
func invertEdits(edits []Edit) []Edit {
	inv := make([]Edit, len(edits))
	delta := 0
	for i, e := range edits {
		inv[i] = Edit{Pos: e.Pos + delta, Len: len(e.Text), Text: e.Old, Old: e.Text}
		delta += len(e.Text) - len(e.Old)
	}
	return inv
}

func (d *Document) applyMulti(edits []Edit) {
	first := edits[0].Pos
	lastE := edits[len(edits)-1]
	last := lastE.Pos + len(lastE.Old)
	old := d.buf.Bytes(first, last)

	m := &multiMap{}
	n := len(edits)
	m.pos, m.end, m.newPos, m.delta = make([]int, n), make([]int, n), make([]int, n), make([]int, n+1)
	m.line, m.lend, m.ldelta = make([]int, n), make([]int, n), make([]int, n+1)
	firstLine := d.buf.LineOfOffset(first)

	var region bytes.Buffer
	region.Grow(len(old) + 64)
	at := first
	line := firstLine
	delta, ldelta := 0, 0
	for i, e := range edits {
		gap := old[at-first : e.Pos-first]
		region.Write(gap)
		line += bytes.Count(gap, nlBytes)
		m.pos[i], m.end[i] = e.Pos, e.Pos+len(e.Old)
		m.delta[i] = delta
		m.newPos[i] = e.Pos + delta
		dl := bytes.Count(e.Old, nlBytes)
		il := bytes.Count(e.Text, nlBytes)
		m.line[i], m.lend[i] = line, line+dl
		m.ldelta[i] = ldelta
		region.Write(e.Text)
		delta += len(e.Text) - len(e.Old)
		ldelta += il - dl
		line += dl
		at = e.Pos + len(e.Old)
	}
	m.delta[n] = delta
	m.ldelta[n] = ldelta
	newRegion := region.Bytes()

	c := &Change{
		Pos:          first,
		Removed:      len(old),
		Inserted:     len(newRegion),
		Line:         firstLine,
		LinesRemoved: bytes.Count(old, nlBytes),
		LinesAdded:   bytes.Count(newRegion, nlBytes),
		AtLineStart:  first == d.buf.LineStart(firstLine),
		multi:        m,
	}
	d.buf.Delete(first, len(old))
	d.buf.Insert(first, newRegion)
	d.version++
	d.Bookmarks.Adjust(c)
	d.Changed.Adjust(c)
	d.Saved.Adjust(c)
	if n <= 100000 {
		var lines []int
		for i, e := range edits {
			from := m.line[i] + m.ldelta[i]
			for l := from; l <= from+bytes.Count(e.Text, nlBytes) && len(lines) < 1_000_000; l++ {
				lines = append(lines, l)
			}
		}
		d.addChanged(lines)
	}
	for i := range d.Marks {
		d.Marks[i].AdjustMulti(c)
	}
	d.hl.changed(c)
	for _, l := range d.listeners {
		l.DocChanged(d, c)
	}
	d.checkModified()
}

// AppendNoUndo adds text at the end without undo history and without
// making the document modified, e.g. the new lines of a monitored log
func (d *Document) AppendNoUndo(text []byte) {
	if len(text) == 0 {
		return
	}
	saved := !d.IsModified()
	d.apply(d.buf.Len(), nil, text)
	d.ClearUndo()
	if saved {
		d.SetSavePoint()
	}
}

func onlyInserts(edits []Edit) bool {
	for _, e := range edits {
		if len(e.Old) > 0 {
			return false
		}
	}
	return true
}
