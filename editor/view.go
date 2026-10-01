package editor

import (
	"bytes"
	"strconv"
	"time"

	"github.com/ipoluianov/nui/ui"
)

// Options are the settings of an editor view
type Options struct {
	FontSize float64
	FontFile string
	Zoom     int // steps of 1 pixel added to the font size

	TabSize     int
	InsertTabs  bool // Tab inserts a tab, otherwise spaces
	AutoIndent  bool
	AutoClose   bool // brackets and quotes
	SmartHome   bool
	WordChars   string
	MultiEdit   bool // Ctrl+click adds carets
	ScrollPast  bool // scrolling beyond the last line
	CaretWidth  int
	CaretBlink  int // ms; 0 - no blinking
	EdgeColumn  int // the long line marker; 0 - none
	WrapIndent  bool
	WordWrap    bool
	LineNumbers bool
	Bookmarks   bool // the bookmark margin
	FoldMargin  bool

	ShowSpaces       bool
	ShowEOL          bool
	ShowIndentGuides bool
	ShowWrapSymbol   bool
	CurrentLine      bool
	SmartHighlight   bool
	SmartMatchCase   bool
	SmartWholeWord   bool
	BraceMatch       bool
	CopyLineNoSel    bool // Ctrl+C/X without a selection copy/cut the line
	AutoComplete     bool // the list of words opens while typing
	ChangeHistory    bool // marks of the changed lines in the margin
}

// DefaultOptions returns the options of a new view
func DefaultOptions() Options {
	return Options{
		FontSize: 14, TabSize: 4, InsertTabs: true, AutoIndent: true, SmartHome: true, MultiEdit: true,
		CaretWidth: 2, CaretBlink: 530, LineNumbers: true, Bookmarks: true, FoldMargin: true, CurrentLine: true,
		SmartHighlight: true, SmartWholeWord: true, BraceMatch: true, CopyLineNoSel: true, ShowWrapSymbol: true, WrapIndent: true,
		ChangeHistory: true,
	}
}

// View is the editor control: it shows a document and edits it
type View struct {
	ui.Widget

	st     *ViewState
	doc    *Document
	opts   Options
	font   *Font
	scheme *Scheme

	// Layout, in pixels
	lineNumW  int
	markerW   int
	foldW     int
	gutterW   int
	textLeft  int // the x of column 0 when not scrolled
	scrollW   int // the size of the scroll bars
	showHBar  bool
	textW     int
	textH     int
	maxLineW  int // the widest line seen, for the horizontal scroll bar
	overtype  bool
	focused   bool
	readOnlyV bool

	caretOn     bool
	caretTime   time.Time
	blinkPassed time.Duration

	mouse mouseState

	// Wrapped rows of the lines, for the current width
	wrapCache   map[int]wrapInfo
	wrapCols    int
	wrapVersion int64

	styles   LineStyles
	lineBuf  []byte
	smart    smartHighlight
	brace    braceMatch
	rows     []rowInfo // the rows painted last, for the mouse
	macro    *MacroRecorder
	lastEdit time.Time
	comp     completion

	// The checkpoints of the long lines, see longline.go
	longCache   map[int]*lineLayout
	longVersion int64

	// OnSelectionChanged is called when the caret moves or the selection changes
	OnSelectionChanged func()
	// OnContextMenu is called on a right click at (x, y) of the form
	OnContextMenu func(x, y int)
	// OnFocused is called when the view gets the focus
	OnFocused func()
	// OnModified is called after each change of the text made in this view
	OnModified func()
	// OnZoom is called when the zoom changes with Ctrl+wheel
	OnZoom func(zoom int)
	// OnDoubleClickLine is called on a double click, before selecting the word;
	// returning true stops the selection (used by the search results)
	OnDoubleClickLine func(line int) bool
	// OnOvertypeChanged is called when Insert switches the overtype mode
	OnOvertypeChanged func()
}

type wrapInfo struct {
	breaks []int
	indent int
	// rows of a long line cut at the width of the view (see longline.go); breaks is nil then
	hardRows int
}

// rowInfo is a row of text on the screen
type rowInfo struct {
	line, sub int
	y         int
}

// NewView makes an editor view with an empty document
func NewView() *View {
	v := &View{}
	v.InitWidget()
	v.SetTypeName("Editor")
	v.SetCanBeFocused(true)
	v.SetXExpandable(true)
	v.SetYExpandable(true)
	v.SetMouseCursor(ui.MouseCursorIBeam)
	v.opts = DefaultOptions()
	v.scheme = Schemes[0]
	v.font = NewFont(v.opts.FontSize, "")
	v.wrapCache = make(map[int]wrapInfo)
	v.scrollW = 12

	doc := NewDocument()
	v.SetState(NewViewState(doc))

	v.SetOnPaint(v.paint)
	v.SetOnKeyDown(v.keyDown)
	v.SetOnChar(v.charTyped)
	v.SetOnKeyUp(func(key ui.Key, mods ui.KeyModifiers) bool {
		TrackModifiers(key, mods, false)
		return false
	})
	v.SetOnMouseDown(v.mouseDown)
	v.SetOnMouseUp(v.mouseUp)
	v.SetOnMouseMove(v.mouseMove)
	v.SetOnMouseDblClick(v.mouseDblClick)
	v.SetOnMouseWheel(v.mouseWheel)
	v.SetOnMouseLeave(func() { v.mouse.overBar = 0; v.update() })
	v.SetOnFocused(func() {
		v.focused = true
		v.resetBlink()
		if v.OnFocused != nil {
			v.OnFocused()
		}
	})
	v.SetOnFocusLost(func() {
		v.focused = false
		v.update()
	})
	v.AddTimer(30, v.timer)
	return v
}

func (v *View) update() {
	if v.Form() != nil {
		v.Form().Update()
	}
}

// Document returns the document shown
func (v *View) Document() *Document { return v.doc }

// State returns the state of the document shown
func (v *View) State() *ViewState { return v.st }

// SetState shows the document of the state, where the state was left
func (v *View) SetState(st *ViewState) {
	if v.st == st {
		return
	}
	if v.st != nil {
		v.st.view = nil
	}
	v.st = st
	v.doc = st.doc
	st.view = v
	v.wrapCache = make(map[int]wrapInfo)
	v.maxLineW = 0
	v.smart = smartHighlight{}
	v.brace = braceMatch{}
	v.layout()
	v.clampScroll()
	v.selectionChanged()
	v.update()
}

// Options returns the options of the view
func (v *View) Options() Options { return v.opts }

// SetOptions applies the options
func (v *View) SetOptions(o Options) {
	if o.TabSize < 1 {
		o.TabSize = 4
	}
	fontChanged := o.FontSize != v.opts.FontSize || o.FontFile != v.opts.FontFile || o.Zoom != v.opts.Zoom
	v.opts = o
	if fontChanged || v.font == nil {
		v.font = NewFont(o.FontSize+float64(o.Zoom), o.FontFile)
	}
	v.wrapCache = make(map[int]wrapInfo)
	v.st.hiddenDirty = true
	v.layout()
	v.clampScroll()
	v.update()
}

// SetScheme sets the color scheme
func (v *View) SetScheme(s *Scheme) {
	v.scheme = s
	v.update()
}

func (v *View) Scheme() *Scheme { return v.scheme }

// Zoom returns the zoom steps
func (v *View) Zoom() int { return v.opts.Zoom }

// SetZoom changes the font size by steps of a pixel
func (v *View) SetZoom(z int) {
	z = max(-10, min(z, 40))
	if z == v.opts.Zoom {
		return
	}
	o := v.opts
	o.Zoom = z
	top := v.st.topLine
	v.SetOptions(o)
	v.st.topLine = top
	v.ensureCaretVisible()
}

// SetReadOnlyView makes the view not edit the document, e.g. the search results
func (v *View) SetReadOnlyView(ro bool) { v.readOnlyV = ro }

// IsOvertype reports whether typed characters replace the text
func (v *View) IsOvertype() bool { return v.overtype }

// Font returns the font of the text
func (v *View) Font() *Font { return v.font }

func (v *View) readOnly() bool {
	return v.readOnlyV || v.doc.ReadOnly
}

func (v *View) tabSize() int {
	return max(1, v.opts.TabSize)
}

// wrapping reports whether the lines are wrapped now
func (v *View) wrapping() bool {
	return v.opts.WordWrap && !v.doc.large
}

// layout computes the sizes of the margins
func (v *View) layout() {
	cw := v.font.CharWidth
	v.lineNumW = 0
	if v.opts.LineNumbers {
		digits := len(strconv.Itoa(v.doc.LineCount()))
		v.lineNumW = max(digits, 3)*cw + cw
	}
	v.markerW = 0
	if v.opts.Bookmarks {
		v.markerW = max(14, v.font.LineHeight-4)
	}
	v.foldW = 0
	if v.opts.FoldMargin && !v.st.noFold() {
		v.foldW = max(12, v.font.LineHeight*2/3)
	}
	v.gutterW = v.lineNumW + v.markerW + v.foldW
	v.textLeft = v.gutterW + max(3, cw/2)
	v.textW = max(1, v.Width()-v.textLeft-v.scrollW)
	v.showHBar = !v.wrapping()
	v.textH = v.Height()
	if v.showHBar {
		v.textH -= v.scrollW
	}
	v.textH = max(v.font.LineHeight, v.textH)
	if v.wrapping() {
		cols := max(8, v.textW/cw-1)
		if v.opts.ShowWrapSymbol {
			cols = max(8, cols-2)
		}
		if cols != v.wrapCols {
			v.wrapCols = cols
			v.wrapCache = make(map[int]wrapInfo)
		}
	}
}

// SetSize lays the view out for its new size
func (v *View) SetSize(w, h int) {
	v.Widget.SetSize(w, h)
	if v.st != nil {
		v.layout()
		v.clampScroll()
	}
}

// visibleRows returns the number of rows that fit the view
func (v *View) visibleRows() int {
	return max(1, v.textH/v.font.LineHeight)
}

// ---- Rows ----

func (v *View) lineText(line int) []byte {
	return v.doc.buf.View(v.doc.LineStart(line), v.doc.LineEnd(line))
}

// wrap returns the wrapped rows of the line
func (v *View) wrap(line int) wrapInfo {
	if !v.wrapping() {
		return wrapInfo{}
	}
	if v.wrapVersion != v.doc.version {
		v.wrapVersion = v.doc.version
		v.wrapCache = make(map[int]wrapInfo)
	}
	if w, ok := v.wrapCache[line]; ok {
		return w
	}
	var w wrapInfo
	if v.isLong(line) {
		cols := v.layoutOf(line).cols
		w.hardRows = max(1, (cols+v.wrapCols-1)/v.wrapCols)
	} else {
		w.breaks, w.indent = wrapLine(v.lineText(line), v.wrapCols, v.tabSize(), v.opts.WrapIndent)
	}
	if len(v.wrapCache) > 10000 {
		v.wrapCache = make(map[int]wrapInfo)
	}
	v.wrapCache[line] = w
	return w
}

// rowCount returns the number of rows of the line
func (v *View) rowCount(line int) int {
	if !v.wrapping() {
		return 1
	}
	w := v.wrap(line)
	if w.hardRows > 0 {
		return w.hardRows
	}
	return len(w.breaks) + 1
}

// rowOfOffset returns the row of the line holding the offset rel within the line
func (v *View) rowOfOffset(line, rel int) int {
	if !v.wrapping() {
		return 0
	}
	w := v.wrap(line)
	if w.hardRows > 0 {
		return min(v.colAt(line, rel)/v.wrapCols, w.hardRows-1)
	}
	sub := 0
	for sub < len(w.breaks) && rel >= w.breaks[sub] {
		sub++
	}
	return sub
}

// rowRange returns the offsets within the line of its row sub
func (v *View) rowRange(line, sub int, lineLen int) (int, int) {
	if !v.wrapping() {
		return 0, lineLen
	}
	w := v.wrap(line)
	if w.hardRows > 0 {
		start, end := 0, lineLen
		if sub > 0 {
			start, _ = v.offAt(line, sub*v.wrapCols, false)
		}
		if sub+1 < w.hardRows {
			end, _ = v.offAt(line, (sub+1)*v.wrapCols, false)
		}
		return start, end
	}
	start, end := 0, lineLen
	if sub > 0 && sub-1 < len(w.breaks) {
		start = w.breaks[sub-1]
	}
	if sub < len(w.breaks) {
		end = w.breaks[sub]
	}
	return start, end
}

func (v *View) nextRow(line, sub int) (int, int, bool) {
	if sub+1 < v.rowCount(line) {
		return line, sub + 1, true
	}
	n := v.st.nextVisible(line)
	if n < 0 {
		return line, sub, false
	}
	return n, 0, true
}

func (v *View) prevRow(line, sub int) (int, int, bool) {
	if sub > 0 {
		return line, sub - 1, true
	}
	p := v.st.prevVisible(line)
	if p < 0 {
		return line, sub, false
	}
	return p, v.rowCount(p) - 1, true
}

// stepRows moves (line, sub) by n rows, down for n > 0
func (v *View) stepRows(line, sub, n int) (int, int) {
	for n > 0 {
		l, s, ok := v.nextRow(line, sub)
		if !ok {
			break
		}
		line, sub = l, s
		n--
	}
	for n < 0 {
		l, s, ok := v.prevRow(line, sub)
		if !ok {
			break
		}
		line, sub = l, s
		n++
	}
	return line, sub
}

// rowsBetween counts the rows from (l1, s1) to (l2, s2), up to limit
func (v *View) rowsBetween(l1, s1, l2, s2, limit int) int {
	n := 0
	for (l1 < l2 || (l1 == l2 && s1 < s2)) && n < limit {
		l, s, ok := v.nextRow(l1, s1)
		if !ok {
			break
		}
		l1, s1 = l, s
		n++
	}
	return n
}

// ---- Scrolling ----

// lastTop returns the highest top row: the last row at the bottom of the view
func (v *View) lastTop() (int, int) {
	last := v.st.visibleLine(v.doc.LineCount() - 1)
	sub := v.rowCount(last) - 1
	if v.opts.ScrollPast {
		return last, sub
	}
	return v.stepRows(last, sub, -(v.visibleRows() - 1))
}

func (v *View) clampScroll() {
	st := v.st
	st.rebuildHidden(v.tabSize())
	st.topLine = max(0, min(st.topLine, v.doc.LineCount()-1))
	st.topLine = st.visibleLine(st.topLine)
	st.topSub = max(0, min(st.topSub, v.rowCount(st.topLine)-1))
	ll, ls := v.lastTop()
	if st.topLine > ll || (st.topLine == ll && st.topSub > ls) {
		st.topLine, st.topSub = ll, ls
	}
	if v.wrapping() {
		st.scrollX = 0
	} else {
		maxX := max(0, v.maxLineW-v.textW+v.font.CharWidth*4)
		st.scrollX = max(0, min(st.scrollX, maxX))
	}
}

// ScrollLines scrolls by n rows, down for n > 0
func (v *View) ScrollLines(n int) {
	v.st.topLine, v.st.topSub = v.stepRows(v.st.topLine, v.st.topSub, n)
	v.clampScroll()
	v.update()
}

// ScrollToLine shows the line at the top of the view
func (v *View) ScrollToLine(line int) {
	v.st.rebuildHidden(v.tabSize())
	v.st.topLine = v.st.visibleLine(max(0, min(line, v.doc.LineCount()-1)))
	v.st.topSub = 0
	v.clampScroll()
	v.update()
}

// caretRow returns the line and row of a position
func (v *View) posRow(pos int) (line, sub int) {
	line = v.doc.LineOfOffset(pos)
	v.st.rebuildHidden(v.tabSize())
	if v.st.isHidden(line) {
		return v.st.visibleLine(line), 0
	}
	return line, v.rowOfOffset(line, pos-v.doc.LineStart(line))
}

// ensureCaretVisible scrolls so that the main caret is shown
func (v *View) ensureCaretVisible() {
	v.ensureVisible(v.st.sels.MainSel().Caret, v.st.sels.MainSel().CaretVS)
}

func (v *View) ensureVisible(pos, vs int) {
	st := v.st
	line := v.doc.LineOfOffset(pos)
	if st.unfoldLine(line, v.tabSize()) {
		v.update()
	}
	cl, cs := v.posRow(pos)
	// Vertically
	if cl < st.topLine || (cl == st.topLine && cs < st.topSub) {
		st.topLine, st.topSub = cl, cs
	} else {
		rows := v.visibleRows()
		if v.rowsBetween(st.topLine, st.topSub, cl, cs, rows+1) >= rows {
			st.topLine, st.topSub = v.stepRows(cl, cs, -(rows - 1))
		}
	}
	// Horizontally
	if !v.wrapping() {
		x := v.xOfPos(pos, vs)
		cw := v.font.CharWidth
		if x < st.scrollX {
			st.scrollX = max(0, x-cw*4)
		} else if x > st.scrollX+v.textW-cw*2 {
			st.scrollX = x - v.textW + cw*6
		}
		if x+cw*6 > v.maxLineW {
			v.maxLineW = x + cw*6
		}
	}
	v.clampScroll()
}

// EnsureCaretVisible scrolls to the caret, centering it when it is far
func (v *View) EnsureCaretVisible(center bool) {
	pos := v.st.sels.MainSel().Caret
	if center {
		line := v.doc.LineOfOffset(pos)
		v.st.unfoldLine(line, v.tabSize())
		cl, cs := v.posRow(pos)
		rows := v.visibleRows()
		st := v.st
		above := cl < st.topLine || (cl == st.topLine && cs < st.topSub)
		below := v.rowsBetween(st.topLine, st.topSub, cl, cs, rows+1) >= rows
		if above || below {
			st.topLine, st.topSub = v.stepRows(cl, cs, -rows/3)
		}
	}
	v.ensureCaretVisible()
	v.update()
}

// xOfPos returns the x of a position from the start of its row's text, in pixels
func (v *View) xOfPos(pos, vs int) int {
	line := v.doc.LineOfOffset(pos)
	rel := pos - v.doc.LineStart(line)
	col := v.colAt(line, rel)
	if v.wrapping() {
		sub := v.rowOfOffset(line, rel)
		rs, _ := v.rowRange(line, sub, v.lineLen(line))
		col -= v.colAt(line, rs)
		if sub > 0 {
			col += v.wrap(line).indent
		}
	}
	return (col + vs) * v.font.CharWidth
}

// colOfPos returns the column of the position in its line
func (v *View) colOfPos(pos int) int {
	line := v.doc.LineOfOffset(pos)
	return v.colAt(line, pos-v.doc.LineStart(line))
}

// posOfCol returns the position at the column of the line and the virtual
// space beyond the end of the line
func (v *View) posOfCol(line, col int, nearest bool) (int, int) {
	off, beyond := v.offAt(line, col, nearest)
	return v.doc.LineStart(line) + off, beyond
}

// ---- Changes ----

// docChanged is called by the state when the document changed
func (v *View) docChanged(c *Change) {
	if c.LinesAdded != c.LinesRemoved {
		oldW := v.lineNumW
		v.layout()
		if v.lineNumW != oldW {
			v.update()
		}
	}
	if c.LinesAdded == 0 && c.LinesRemoved == 0 && c.multi == nil {
		delete(v.wrapCache, c.Line)
		v.wrapVersion = v.doc.version
	}
	v.smart.valid = false
	v.brace.valid = false
	v.update()
}

// selectionChanged is called after the selections changed
func (v *View) selectionChanged() {
	v.resetBlink()
	v.smart.valid = false
	v.brace.valid = false
	if v.OnSelectionChanged != nil {
		v.OnSelectionChanged()
	}
	v.update()
}

func (v *View) resetBlink() {
	v.caretOn = true
	v.caretTime = time.Now()
	v.update()
}

func (v *View) timer() {
	if v.focused && v.opts.CaretBlink > 0 {
		if time.Since(v.caretTime) > time.Duration(v.opts.CaretBlink)*time.Millisecond {
			v.caretOn = !v.caretOn
			v.caretTime = time.Now()
			v.update()
		}
	} else if !v.caretOn {
		v.caretOn = true
		v.update()
	}
	v.autoScroll()
}

// ---- Public API ----

// Selections returns the selections
func (v *View) Selections() []Sel { return append([]Sel(nil), v.st.sels.List...) }

// MainSelection returns the main selection
func (v *View) MainSelection() Sel { return v.st.sels.MainSel() }

// IsRectSelection reports whether the selection is rectangular
func (v *View) IsRectSelection() bool { return v.st.sels.Rect }

// SetSelection selects the text from anchor to caret, as the only selection
func (v *View) SetSelection(anchor, caret int) {
	n := v.doc.Len()
	v.st.sels.setSingle(Sel{Anchor: max(0, min(anchor, n)), Caret: max(0, min(caret, n)), WantCol: -1})
	v.selectionChanged()
}

// SetSelections sets several selections; main is the one scrolled to
func (v *View) SetSelections(sels []Sel, main int) {
	if len(sels) == 0 {
		return
	}
	v.st.sels.List = append(v.st.sels.List[:0], sels...)
	v.st.sels.Main = max(0, min(main, len(sels)-1))
	v.st.sels.Rect = false
	v.st.sels.normalize()
	v.selectionChanged()
}

// SetCaret moves the caret, dropping the selection
func (v *View) SetCaret(pos int) {
	v.SetSelection(pos, pos)
}

// Caret returns the position of the main caret
func (v *View) Caret() int { return v.st.sels.MainSel().Caret }

// SelectedText returns the text of the main selection
func (v *View) SelectedText() string {
	s := v.st.sels.MainSel()
	return string(v.doc.Text(s.Start(), s.End()))
}

// AllSelectedText returns the text of all the selections, one per line
func (v *View) AllSelectedText() string {
	var b bytes.Buffer
	for i, s := range v.st.sels.List {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.Write(v.doc.Text(s.Start(), s.End()))
	}
	return b.String()
}

// HasSelection reports whether any text is selected
func (v *View) HasSelection() bool {
	for _, s := range v.st.sels.List {
		if !s.IsEmptyText() {
			return true
		}
	}
	return false
}

// CaretLineCol returns the line and the column (1-based, tabs expanded) of the main caret
func (v *View) CaretLineCol() (line, col int) {
	s := v.st.sels.MainSel()
	line = v.doc.LineOfOffset(s.Caret)
	return line + 1, v.colOfPos(s.Caret) + s.CaretVS + 1
}

// GotoLine moves the caret to the start of the line (0-based) and shows it
func (v *View) GotoLine(line int) {
	line = max(0, min(line, v.doc.LineCount()-1))
	v.SetCaret(v.doc.LineStart(line))
	v.EnsureCaretVisible(true)
}

// GotoPos moves the caret to the position and shows it
func (v *View) GotoPos(pos int) {
	v.SetCaret(pos)
	v.EnsureCaretVisible(true)
}

// SelectRange selects the range and shows it
func (v *View) SelectRange(start, end int) {
	v.SetSelection(start, end)
	v.ensureVisible(start, 0)
	v.ensureVisible(end, 0)
	v.EnsureCaretVisible(true)
}

// Refresh repaints the view, e.g. after the document changed its language
func (v *View) Refresh() {
	v.wrapCache = make(map[int]wrapInfo)
	v.st.hiddenDirty = true
	v.layout()
	v.clampScroll()
	v.smart.valid = false
	v.update()
}

// PosOfColumn returns the position at the column of the line of the document
// and the columns of virtual space beyond the end of the line
func PosOfColumn(doc *Document, line, col, tabSize int) (int, int) {
	text := doc.LineText(line)
	off, beyond := offsetOfCol(text, col, max(1, tabSize), false)
	return doc.LineStart(line) + off, beyond
}

// ColumnOf returns the column of the position in its line
func ColumnOf(doc *Document, pos, tabSize int) int {
	line := doc.LineOfOffset(pos)
	return colOfOffset(doc.LineText(line), pos-doc.LineStart(line), max(1, tabSize))
}

// WordAtCaret returns the word at the main caret
func (v *View) WordAtCaret() (int, int) {
	return v.wordAt(v.Caret())
}

// LastEdit returns the time of the last change made in the view
func (v *View) LastEdit() time.Time { return v.lastEdit }
