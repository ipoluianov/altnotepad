package editor

import (
	"time"

	"github.com/ipoluianov/nui/ui"
)

// mouseState is the mouse gesture in progress
type mouseState struct {
	selecting bool
	mode      int // selDrag*: what a drag selects
	rect      bool
	anchorA   int // the selected unit the drag started in: from..to
	anchorB   int
	lineSel   bool // a drag in the line number margin selects lines
	x, y      int  // last position

	dragBar      int // 1 - dragging the vertical thumb, 2 - the horizontal one
	dragStart    int
	dragStartPos int
	overBar      int

	lastClick     time.Time
	lastClickX    int
	lastClickY    int
	clicks        int
	suppressClick bool
}

const (
	selDragChar = iota
	selDragWord
	selDragLine
)

// keyDown handles the keys of the editor; the other keys go on to the form
func (v *View) keyDown(key ui.Key, mods ui.KeyModifiers) bool {
	TrackModifiers(key, mods, true)
	if v.completionKey(key, mods) {
		return true
	}
	cmd := keyCommand(key, mods)
	if cmd == "" {
		return false
	}
	v.Exec(cmd, "")
	if v.comp.active {
		if cmd == "backspace" {
			v.filterCompletion()
		} else {
			v.comp.active = false
		}
		v.update()
	}
	return true
}

// keyCommand returns the editor command of a key, "" if none
func keyCommand(key ui.Key, m ui.KeyModifiers) string {
	shift, ctrl, alt := m.Shift, m.Ctrl || m.Cmd, m.Alt
	plain := !ctrl && !alt
	sfx := func(base string) string {
		switch {
		case alt && shift && !ctrl:
			return base + "-rect"
		case shift && !alt:
			return base + "-extend"
		case !shift && !alt:
			return base
		}
		return ""
	}
	switch key {
	case ui.KeyArrowLeft:
		if ctrl {
			return sfx("word-left")
		}
		return sfx("char-left")
	case ui.KeyArrowRight:
		if ctrl {
			return sfx("word-right")
		}
		return sfx("char-right")
	case ui.KeyArrowUp:
		switch {
		case ctrl && shift:
			return "move-lines-up"
		case ctrl:
			return "scroll-up"
		}
		return sfx("line-up")
	case ui.KeyArrowDown:
		switch {
		case ctrl && shift:
			return "move-lines-down"
		case ctrl:
			return "scroll-down"
		}
		return sfx("line-down")
	case ui.KeyHome:
		if ctrl {
			return sfx("doc-start")
		}
		return sfx("home")
	case ui.KeyEnd:
		if ctrl {
			return sfx("doc-end")
		}
		return sfx("end")
	case ui.KeyPageUp:
		if ctrl {
			return ""
		}
		return sfx("page-up")
	case ui.KeyPageDown:
		if ctrl {
			return ""
		}
		return sfx("page-down")
	case ui.KeyBackspace:
		switch {
		case ctrl && shift:
			return "delete-line-left"
		case ctrl:
			return "delete-word-left"
		case alt:
			return "undo"
		}
		return "backspace"
	case ui.KeyDelete:
		switch {
		case ctrl && shift:
			return "delete-line-right"
		case ctrl:
			return "delete-word-right"
		case shift:
			return "cut"
		}
		return "delete"
	case ui.KeyEnter:
		if plain || (shift && !ctrl && !alt) {
			return "newline"
		}
	case ui.KeyTab:
		if !ctrl && !alt {
			if shift {
				return "backtab"
			}
			return "tab"
		}
	case ui.KeyInsert:
		switch {
		case ctrl && !shift:
			return "copy"
		case shift && !ctrl:
			return "paste"
		case plain:
			return "toggle-overtype"
		}
	case ui.KeyEsc:
		if plain {
			return "cancel"
		}
	}
	if ctrl && !alt {
		switch key {
		case ui.KeyA:
			if !shift {
				return "select-all"
			}
		case ui.KeyC:
			if !shift {
				return "copy"
			}
		case ui.KeyX:
			if !shift {
				return "cut"
			}
		case ui.KeyV:
			if !shift {
				return "paste"
			}
		case ui.KeyZ:
			if shift {
				return "redo"
			}
			return "undo"
		case ui.KeyY:
			if !shift {
				return "redo"
			}
		}
	}
	return ""
}

// charTyped inserts a typed character
func (v *View) charTyped(r rune, mods ui.KeyModifiers) bool {
	// Ctrl and Alt make shortcuts; both together are AltGr on Windows
	if (mods.Ctrl && !mods.Alt) || (mods.Alt && !mods.Ctrl) || mods.Cmd {
		return false
	}
	if r < 32 || r == 127 {
		return r != 0
	}
	v.Exec("type", string(r))
	word := charClass(r, v.opts.WordChars) == ccWord
	switch {
	case v.comp.active && word:
		v.filterCompletion()
	case v.comp.active:
		v.comp.active = false
	case word && v.opts.AutoComplete:
		v.Complete(true)
	}
	return true
}

// ---- Mouse ----

// hitTest returns the position under (x, y) of the view: the nearest
// character boundary, the virtual space beyond the line end, the line and
// the column in it
func (v *View) hitTest(x, y int) (pos, vs, line, col int) {
	lh := v.font.LineHeight
	cw := v.font.CharWidth
	var row rowInfo
	switch {
	case len(v.rows) == 0:
		row = rowInfo{line: v.st.topLine}
	case y < 0:
		l, s := v.stepRows(v.rows[0].line, v.rows[0].sub, (y-lh+1)/lh)
		row = rowInfo{line: l, sub: s}
	default:
		k := y / lh
		if k < len(v.rows) {
			row = v.rows[k]
		} else {
			last := v.rows[len(v.rows)-1]
			l, s := v.stepRows(last.line, last.sub, k-len(v.rows)+1)
			row = rowInfo{line: l, sub: s}
		}
	}
	lineLen := v.lineLen(row.line)
	rs, re := v.rowRange(row.line, row.sub, lineLen)
	anchorCol := 0
	if v.wrapping() {
		anchorCol = v.colAt(row.line, rs)
	}
	px := x - v.textLeft + v.st.scrollX
	if row.sub > 0 {
		px -= v.wrap(row.line).indent * cw
	}
	colF := px
	if colF < 0 {
		colF = 0
	}
	col = anchorCol + (colF+cw/2)/cw
	off, beyond := v.offAt(row.line, col, true)
	if off > re || (off == re && re < lineLen && v.wrapping()) {
		// Past the end of a wrapped row: the end of the row
		off, beyond = re, 0
		if re < lineLen && re > rs {
			off = v.charLeft(v.doc.LineStart(row.line)+re) - v.doc.LineStart(row.line)
			off = max(off, rs)
		}
	}
	return v.doc.LineStart(row.line) + off, beyond, row.line, col
}

func decodeLastRune(b []byte) (rune, int) {
	if len(b) == 0 {
		return 0, 0
	}
	for i := len(b) - 1; i >= 0 && i >= len(b)-4; i-- {
		if b[i] < 0x80 || b[i] >= 0xC0 {
			return 0, len(b) - i
		}
	}
	return 0, 1
}

// region returns what is under x: 0 - text, 1 - line numbers, 2 - bookmarks,
// 3 - folds, 4 - vertical scroll bar, 5 - horizontal scroll bar
func (v *View) region(x, y int) int {
	switch {
	case x >= v.Width()-v.scrollW:
		return 4
	case v.showHBar && y >= v.Height()-v.scrollW && x >= v.gutterW:
		return 5
	case x < v.lineNumW:
		return 1
	case x < v.lineNumW+v.markerW:
		return 2
	case x < v.gutterW:
		return 3
	}
	return 0
}

func (v *View) mouseDown(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
	v.Focus()
	v.mouse.x, v.mouse.y = x, y
	if button == ui.MouseButtonLeft && v.completionMouse(x, y) {
		return true
	}
	v.comp.active = false
	reg := v.region(x, y)
	if button == ui.MouseButtonRight {
		if reg == 0 {
			pos, _, _, _ := v.hitTest(x, y)
			inSel := false
			for _, s := range v.st.sels.List {
				if pos >= s.Start() && pos <= s.End() && !s.IsEmptyText() {
					inSel = true
				}
			}
			if !inSel {
				v.st.sels.setSingle(caretSel(pos))
				v.selectionChanged()
			}
		}
		if v.OnContextMenu != nil {
			wx, wy := v.RectClientAreaOnWindow()
			v.OnContextMenu(wx+x, wy+y)
		}
		return true
	}
	if button != ui.MouseButtonLeft {
		return true
	}
	switch reg {
	case 4:
		v.scrollBarDown(y, true)
		return true
	case 5:
		v.scrollBarDown(x-v.gutterW, false)
		return true
	case 2, 3:
		_, _, line, _ := v.hitTest(v.textLeft, y)
		if reg == 2 {
			v.Exec("toggle-bookmark-line", itoa(line))
		} else {
			v.toggleFoldAt(line, mods.Shift)
		}
		return true
	}

	// Clicks in a row: double and triple ones select words and lines
	now := time.Now()
	if now.Sub(v.mouse.lastClick) < 500*time.Millisecond && abs(x-v.mouse.lastClickX) < 5 && abs(y-v.mouse.lastClickY) < 5 {
		v.mouse.clicks++
	} else {
		v.mouse.clicks = 1
	}
	v.mouse.lastClick, v.mouse.lastClickX, v.mouse.lastClickY = now, x, y
	if v.mouse.clicks == 2 {
		return true // handled by mouseDblClick
	}

	pos, vs, line, col := v.hitTest(x, y)
	st := v.st
	v.mouse.selecting = true
	v.mouse.rect = false
	v.mouse.lineSel = reg == 1
	switch {
	case reg == 1 || v.mouse.clicks >= 3:
		// Line selection
		v.mouse.mode = selDragLine
		a := v.doc.LineStart(line)
		b := v.lineEndWithBreak(line)
		if mods.Shift {
			m := st.sels.MainSel()
			v.mouse.anchorA, v.mouse.anchorB = v.doc.LineStart(v.doc.LineOfOffset(m.Anchor)), v.lineEndWithBreak(v.doc.LineOfOffset(m.Anchor))
			v.extendDrag(pos)
			return true
		}
		v.mouse.anchorA, v.mouse.anchorB = a, b
		st.sels.setSingle(Sel{Anchor: a, Caret: b, WantCol: -1})
	case mods.Alt && !mods.Ctrl:
		// Rectangular selection
		v.mouse.mode = selDragChar
		v.mouse.rect = true
		if !(mods.Shift && st.sels.Rect) {
			m := st.sels.MainSel()
			if mods.Shift {
				st.sels.RectAnchorLine = v.doc.LineOfOffset(m.Anchor)
				st.sels.RectAnchorCol = v.colOfPos(m.Anchor) + m.AnchorVS
			} else {
				st.sels.RectAnchorLine, st.sels.RectAnchorCol = line, col
			}
		}
		st.sels.Rect = true
		st.sels.RectCaretLine, st.sels.RectCaretCol = line, col
		v.rebuildRect()
	case mods.Ctrl && !mods.Shift && v.opts.MultiEdit:
		// Another caret
		v.mouse.mode = selDragChar
		if st.sels.Rect {
			st.sels.Rect = false
		}
		st.sels.List = append(st.sels.List, caretSel(pos))
		st.sels.Main = len(st.sels.List) - 1
		v.mouse.anchorA, v.mouse.anchorB = pos, pos
		st.sels.normalize()
	case mods.Shift:
		v.mouse.mode = selDragChar
		m := st.sels.MainSel()
		v.mouse.anchorA, v.mouse.anchorB = m.Anchor, m.Anchor
		st.sels.setSingle(Sel{Anchor: m.Anchor, Caret: pos, WantCol: -1})
	default:
		v.mouse.mode = selDragChar
		v.mouse.anchorA, v.mouse.anchorB = pos, pos
		s := caretSel(pos)
		_ = vs
		st.sels.setSingle(s)
	}
	v.selectionChanged()
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// lineEndWithBreak returns the end of the line including its line break
func (v *View) lineEndWithBreak(line int) int {
	if line+1 < v.doc.LineCount() {
		return v.doc.LineStart(line + 1)
	}
	return v.doc.Len()
}

func (v *View) mouseDblClick(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
	if button != ui.MouseButtonLeft || v.region(x, y) != 0 {
		if button == ui.MouseButtonLeft && v.region(x, y) == 3 {
			_, _, line, _ := v.hitTest(v.textLeft, y)
			v.toggleFoldAt(line, mods.Shift)
		}
		return true
	}
	v.mouse.clicks = 2
	v.mouse.lastClick = time.Now()
	pos, _, line, _ := v.hitTest(x, y)
	if v.OnDoubleClickLine != nil && v.OnDoubleClickLine(line) {
		v.mouse.selecting = false
		return true
	}
	a, b := v.wordAt(pos)
	st := v.st
	if mods.Ctrl && v.opts.MultiEdit {
		st.sels.List = append(st.sels.List, Sel{Anchor: a, Caret: b, WantCol: -1})
		st.sels.Main = len(st.sels.List) - 1
		st.sels.normalize()
	} else {
		st.sels.setSingle(Sel{Anchor: a, Caret: b, WantCol: -1})
	}
	v.mouse.selecting = true
	v.mouse.mode = selDragWord
	v.mouse.anchorA, v.mouse.anchorB = a, b
	v.selectionChanged()
	return true
}

// wordAt returns the word (or the run of spaces or of other characters) at pos
func (v *View) wordAt(pos int) (int, int) {
	doc := v.doc
	line := doc.LineOfOffset(pos)
	ls, le := doc.LineStart(line), doc.LineEnd(line)
	if ls == le {
		return pos, pos
	}
	r, _ := doc.buf.DecodeRune(pos)
	cls := charClass(r, v.opts.WordChars)
	if pos >= le || cls == ccLineEnd {
		r, _ = doc.buf.DecodeLastRune(pos)
		cls = charClass(r, v.opts.WordChars)
		if cls != ccWord {
			return pos, pos
		}
	}
	a := pos
	for a > ls {
		r, size := doc.buf.DecodeLastRune(a)
		if charClass(r, v.opts.WordChars) != cls {
			break
		}
		a -= size
	}
	b := pos
	for b < le {
		r, size := doc.buf.DecodeRune(b)
		if charClass(r, v.opts.WordChars) != cls {
			break
		}
		b += size
	}
	return a, b
}

// isWholeWord reports whether a..b is a word with non-word characters around it
func (v *View) isWholeWord(a, b int) bool {
	for p := a; p < b; {
		r, size := v.doc.buf.DecodeRune(p)
		if charClass(r, v.opts.WordChars) != ccWord {
			return false
		}
		p += size
	}
	return isWordBoundary(v.doc, a, v.opts.WordChars) && isWordBoundary(v.doc, b, v.opts.WordChars)
}

func (v *View) mouseMove(x, y int, mods ui.KeyModifiers) bool {
	v.mouse.x, v.mouse.y = x, y
	if v.mouse.dragBar != 0 {
		v.scrollBarDrag(x, y)
		return true
	}
	reg := v.region(x, y)
	over := 0
	if reg == 4 {
		over = 1
	} else if reg == 5 {
		over = 2
	}
	if over != v.mouse.overBar {
		v.mouse.overBar = over
		v.update()
	}
	if !v.mouse.selecting {
		if reg == 0 {
			v.SetMouseCursor(ui.MouseCursorIBeam)
		} else {
			v.SetMouseCursor(ui.MouseCursorArrow)
		}
		return true
	}
	v.dragTo(x, y)
	return true
}

func (v *View) dragTo(x, y int) {
	pos, _, line, col := v.hitTest(x, y)
	if v.mouse.rect {
		st := v.st
		st.sels.RectCaretLine, st.sels.RectCaretCol = line, col
		v.rebuildRect()
		v.selectionChanged()
		return
	}
	v.extendDrag(pos)
}

// extendDrag extends the selection being dragged to pos, by the unit of the drag
func (v *View) extendDrag(pos int) {
	st := v.st
	a, b := v.mouse.anchorA, v.mouse.anchorB
	var anchor, caret int
	switch v.mouse.mode {
	case selDragWord:
		wa, wb := v.wordAt(pos)
		if pos < a {
			anchor, caret = b, wa
		} else {
			anchor, caret = a, max(wb, b)
		}
	case selDragLine:
		line := v.doc.LineOfOffset(pos)
		la, lb := v.doc.LineStart(line), v.lineEndWithBreak(line)
		if la < a {
			anchor, caret = b, la
		} else {
			anchor, caret = a, lb
		}
	default:
		anchor, caret = a, pos
	}
	m := &st.sels.List[st.sels.Main]
	m.Anchor, m.Caret = anchor, caret
	m.WantCol = -1
	m.AnchorVS, m.CaretVS = 0, 0
	st.sels.normalize()
	v.selectionChanged()
}

func (v *View) mouseUp(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
	if v.mouse.dragBar != 0 {
		v.mouse.dragBar = 0
		v.update()
		return true
	}
	if v.mouse.selecting {
		v.mouse.selecting = false
		v.ensureCaretVisible()
		v.update()
	}
	return true
}

// autoScroll scrolls while a selection is dragged beyond the text
func (v *View) autoScroll() {
	if !v.mouse.selecting {
		return
	}
	x, y := v.mouse.x, v.mouse.y
	lh := v.font.LineHeight
	scrolled := false
	if y < 0 {
		v.ScrollLines(-max(1, -y/lh))
		scrolled = true
	} else if y > v.textH {
		v.ScrollLines(max(1, (y-v.textH)/lh))
		scrolled = true
	}
	if !v.wrapping() {
		cw := v.font.CharWidth
		if x < v.textLeft && v.st.scrollX > 0 {
			v.st.scrollX = max(0, v.st.scrollX-cw*max(1, (v.textLeft-x)/cw))
			scrolled = true
		} else if x > v.textLeft+v.textW {
			v.st.scrollX += cw * max(1, (x-v.textLeft-v.textW)/cw)
			if v.st.scrollX+v.textW > v.maxLineW {
				v.maxLineW = v.st.scrollX + v.textW
			}
			scrolled = true
		}
	}
	if scrolled {
		v.dragTo(x, y)
	}
}

func (v *View) mouseWheel(dx, dy int) bool {
	if v.completionWheel(dy) {
		return true
	}
	if v.Form() != nil {
		mods := v.formMods()
		if mods.Ctrl {
			if dy > 0 {
				v.SetZoom(v.opts.Zoom + 1)
			} else if dy < 0 {
				v.SetZoom(v.opts.Zoom - 1)
			}
			if v.OnZoom != nil {
				v.OnZoom(v.opts.Zoom)
			}
			return true
		}
	}
	if dy != 0 {
		v.ScrollLines(-dy * 3)
	}
	if dx != 0 && !v.wrapping() {
		v.st.scrollX = max(0, v.st.scrollX-dx*v.font.CharWidth*6)
		v.clampScroll()
		v.update()
	}
	if v.mouse.selecting {
		v.dragTo(v.mouse.x, v.mouse.y)
	}
	return true
}

// formMods returns the modifier keys held now
func (v *View) formMods() ui.KeyModifiers {
	return currentMods
}

// currentMods are the modifier keys held now, for the mouse wheel
var currentMods ui.KeyModifiers

// TrackModifiers follows the modifier keys held, as the wheel events do not
// carry them; the application calls it from its key handlers too
func TrackModifiers(key ui.Key, mods ui.KeyModifiers, down bool) {
	currentMods = mods
	// The state of a key event is from before it
	switch key {
	case ui.KeyCtrl:
		currentMods.Ctrl = down
	case ui.KeyShift:
		currentMods.Shift = down
	case ui.KeyAlt:
		currentMods.Alt = down
	}
}

// ---- Scroll bars ----

func (v *View) scrollBarDown(p int, vertical bool) {
	if vertical {
		ty, th, pos, maxPos := v.vScrollGeometry()
		if maxPos == 0 {
			return
		}
		switch {
		case p < ty:
			v.ScrollLines(-(v.visibleRows() - 1))
		case p >= ty+th:
			v.ScrollLines(v.visibleRows() - 1)
		default:
			v.mouse.dragBar = 1
			v.mouse.dragStart = p
			v.mouse.dragStartPos = pos
		}
		return
	}
	tx, tw, content := v.hScrollGeometry()
	if content <= v.textW {
		return
	}
	switch {
	case p < tx:
		v.st.scrollX = max(0, v.st.scrollX-v.textW*3/4)
	case p >= tx+tw:
		v.st.scrollX += v.textW * 3 / 4
	default:
		v.mouse.dragBar = 2
		v.mouse.dragStart = p
		v.mouse.dragStartPos = v.st.scrollX
	}
	v.clampScroll()
	v.update()
}

func (v *View) scrollBarDrag(x, y int) {
	if v.mouse.dragBar == 1 {
		_, th, _, maxPos := v.vScrollGeometry()
		track := v.textH - th
		if track <= 0 || maxPos == 0 {
			return
		}
		pos := v.mouse.dragStartPos + (y-v.mouse.dragStart)*maxPos/track
		pos = max(0, min(pos, maxPos))
		v.st.topLine = v.st.lineOfVisibleIndex(pos)
		v.st.topSub = 0
		v.clampScroll()
		v.update()
		return
	}
	_, tw, content := v.hScrollGeometry()
	track := v.textW + v.textLeft - v.gutterW - tw
	if track <= 0 {
		return
	}
	v.st.scrollX = v.mouse.dragStartPos + (x-v.gutterW-v.mouse.dragStart)*(content-v.textW)/track
	v.clampScroll()
	v.update()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
