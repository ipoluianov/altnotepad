package editor

import (
	"bytes"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ipoluianov/nui/ui"
)

// edit runs f as one undo step; the selections before and after it are
// restored by undo and redo
func (v *View) edit(f func()) bool {
	if v.readOnly() {
		return false
	}
	before := v.st.sels.clone()
	ver := v.doc.version
	v.doc.BeginAction()
	f()
	v.doc.EndAction()
	v.st.sels.normalize()
	if v.doc.version != ver {
		after := v.st.sels.clone()
		v.doc.SetUndoSelection(&before, &after)
		v.lastEdit = time.Now()
		if v.OnModified != nil {
			v.OnModified()
		}
	}
	v.ensureCaretVisible()
	v.selectionChanged()
	return true
}

// selEdit replaces start..end with text; the selection after it is at
// anchor..caret, offsets from start in the new text
type selEdit struct {
	start, end    int
	text          []byte
	anchor, caret int
	vs            int // virtual space of the caret after the edit
}

// applySelEdits makes an edit for each selection at once and sets the selections after them
func (v *View) applySelEdits(edits []selEdit) {
	if len(edits) == 0 {
		return
	}
	docEdits := make([]Edit, 0, len(edits))
	kept := edits[:0]
	prevEnd := -1
	for _, e := range edits {
		if e.start < prevEnd {
			continue
		}
		prevEnd = e.end
		kept = append(kept, e)
		docEdits = append(docEdits, Edit{Pos: e.start, Len: e.end - e.start, Text: e.text})
	}
	edits = kept
	main := v.st.sels.Main
	if len(docEdits) == 1 {
		v.doc.Replace(docEdits[0].Pos, docEdits[0].Len, docEdits[0].Text)
	} else {
		v.doc.ApplyEdits(docEdits)
	}
	sels := make([]Sel, len(edits))
	delta := 0
	for i, e := range edits {
		base := e.start + delta
		sels[i] = Sel{Anchor: base + e.anchor, Caret: base + e.caret, WantCol: -1, CaretVS: e.vs, AnchorVS: e.vs}
		delta += len(e.text) - (e.end - e.start)
	}
	v.st.sels.List = sels
	v.st.sels.Main = min(main, len(sels)-1)
}

// sortedSels returns the selections sorted by position
func (v *View) sortedSels() []Sel {
	v.st.sels.normalize()
	return v.st.sels.List
}

// ---- Typing ----

var closingOf = map[string]string{"(": ")", "[": "]", "{": "}", `"`: `"`, "'": "'", "`": "`"}

// TypeText inserts the text at every caret as if typed
func (v *View) TypeText(text string) {
	if text == "" {
		return
	}
	v.record("type", text)
	v.edit(func() {
		v.doc.SetTyping()
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			start, end := s.Start(), s.End()
			ins := text
			pad := ""
			if s.IsEmptyText() {
				pad = strings.Repeat(" ", min(s.CaretVS, s.AnchorVS))
				if !v.st.sels.Rect {
					pad = strings.Repeat(" ", s.CaretVS)
				}
			}
			caret := len(pad) + len(ins)
			if v.overtype && s.IsEmptyText() && pad == "" {
				// Replace the character after the caret, not the line break
				if r, size := v.doc.buf.DecodeRune(end); size > 0 && r != '\n' {
					end += size
				}
			}
			if v.opts.AutoClose && utf8.RuneCountInString(text) == 1 {
				next, _ := v.doc.buf.DecodeRune(end)
				switch {
				case s.IsEmptyText() && isClosing(text) && string(next) == text && pad == "":
					// Type over the closing character
					edits = append(edits, selEdit{start: start, end: start, text: nil, anchor: len(text), caret: len(text)})
					continue
				case closingOf[text] != "" && !s.IsEmptyText() && !v.st.sels.Rect:
					// Enclose the selection
					inner := v.doc.Text(start, end)
					out := text + string(inner) + closingOf[text]
					edits = append(edits, selEdit{start: start, end: end, text: []byte(out), anchor: len(text), caret: len(text) + len(inner)})
					continue
				case closingOf[text] != "" && s.IsEmptyText() && v.canAutoClose(end, text):
					ins = text + closingOf[text]
				}
			}
			edits = append(edits, selEdit{start: start, end: end, text: []byte(pad + ins), anchor: caret, caret: caret})
		}
		v.applySelEdits(edits)
		v.keepRectAfterEdit()
	})
}

func isClosing(s string) bool {
	return s == ")" || s == "]" || s == "}" || s == `"` || s == "'" || s == "`"
}

// canAutoClose: a bracket or quote is closed at once when nothing but
// spaces, closing brackets or the line end follow, and a quote is not
// after a word character
func (v *View) canAutoClose(pos int, open string) bool {
	next, _ := v.doc.buf.DecodeRune(pos)
	if !(pos >= v.doc.Len() || next == '\n' || next == ' ' || next == '\t' || strings.ContainsRune(")]};,", next)) {
		return false
	}
	if open == `"` || open == "'" || open == "`" {
		prev, _ := v.doc.buf.DecodeLastRune(pos)
		if charClass(prev, v.opts.WordChars) == ccWord || string(prev) == open {
			return false
		}
	}
	return true
}

// keepRectAfterEdit: typing in a rectangular selection goes on at the
// carets of all its rows
func (v *View) keepRectAfterEdit() {
	st := v.st
	if !st.sels.Rect {
		return
	}
	for _, s := range st.sels.List {
		if !s.IsEmptyText() {
			st.sels.Rect = false
			return
		}
	}
	m := st.sels.MainSel()
	col := v.colOfPos(m.Caret) + m.CaretVS
	st.sels.RectAnchorCol, st.sels.RectCaretCol = col, col
	for i := range st.sels.List {
		st.sels.List[i].AnchorVS = st.sels.List[i].CaretVS
	}
}

func (v *View) indentUnit(col int) string {
	if v.opts.InsertTabs {
		return "\t"
	}
	n := v.tabSize() - col%v.tabSize()
	return strings.Repeat(" ", n)
}

// leadingIndent returns the indentation at the start of the line text
func leadingIndent(text []byte) []byte {
	i := 0
	for i < len(text) && (text[i] == ' ' || text[i] == '\t') {
		i++
	}
	return text[:i]
}

// Newline breaks the line at every caret, indenting the new line
func (v *View) Newline() {
	v.record("newline", "")
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		eolWidth := 0
		for _, s := range sels {
			start, end := s.Start(), s.End()
			line := v.doc.LineOfOffset(start)
			ls := v.doc.LineStart(line)
			text := v.doc.buf.View(ls, v.doc.LineEnd(line))
			ind := ""
			extra := ""
			tail := ""
			if v.opts.AutoIndent {
				lead := leadingIndent(text)
				if len(lead) > start-ls {
					lead = lead[:start-ls]
				}
				ind = string(lead)
				before := bytes.TrimRight(text[:start-ls], " \t")
				after := bytes.TrimLeft(text[min(end-ls, len(text)):], " \t")
				if len(before) > 0 {
					last := before[len(before)-1]
					opensBlock := last == '{' || last == '[' || last == '(' ||
						(last == ':' && v.doc.lang != nil && v.doc.lang.FoldIndent)
					if opensBlock {
						extra = v.indentUnit(colOfOffset(text, len(ind), v.tabSize()))
						if len(after) > 0 && (after[0] == '}' || after[0] == ']' || after[0] == ')') {
							tail = "\n" + ind
						}
					}
				}
				// Spaces before the caret are not kept at the end of the line
				if s.IsEmptyText() {
					k := start
					for k > ls && (v.doc.buf.ByteAt(k-1) == ' ' || v.doc.buf.ByteAt(k-1) == '\t') && k-ls > len(ind) {
						k--
					}
					if k-ls >= len(ind) {
						start = k
					}
				}
			}
			ins := "\n" + ind + extra + tail
			caret := 1 + len(ind) + len(extra)
			edits = append(edits, selEdit{start: start, end: end, text: []byte(ins), anchor: caret, caret: caret})
		}
		_ = eolWidth
		v.st.sels.Rect = false
		v.applySelEdits(edits)
	})
}

// Backspace deletes the selections or the characters before the carets
func (v *View) Backspace() {
	v.record("backspace", "")
	st := v.st
	// In the virtual space the carets just move back
	allVirtual := true
	for _, s := range st.sels.List {
		if !(s.IsEmptyText() && s.CaretVS > 0) {
			allVirtual = false
		}
	}
	if allVirtual {
		for i := range st.sels.List {
			st.sels.List[i].CaretVS--
			st.sels.List[i].AnchorVS = st.sels.List[i].CaretVS
		}
		v.keepRectAfterEdit()
		v.selectionChanged()
		return
	}
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			if !s.IsEmptyText() {
				edits = append(edits, selEdit{start: s.Start(), end: s.End()})
				continue
			}
			pos := s.Caret
			if s.CaretVS > 0 {
				edits = append(edits, selEdit{start: pos, end: pos, text: []byte(strings.Repeat(" ", s.CaretVS-1)), anchor: s.CaretVS - 1, caret: s.CaretVS - 1})
				continue
			}
			if pos == 0 {
				edits = append(edits, selEdit{start: 0, end: 0})
				continue
			}
			r, size := v.doc.buf.DecodeLastRune(pos)
			start := pos - size
			// Unindent to the previous tab stop in the leading spaces
			if r == ' ' && !v.opts.InsertTabs && pos-v.doc.LineStart(v.doc.LineOfOffset(pos)) < 4096 {
				line := v.doc.LineOfOffset(pos)
				ls := v.doc.LineStart(line)
				text := v.doc.buf.View(ls, pos)
				if len(leadingIndent(text)) == len(text) && bytes.IndexByte(text, '\t') < 0 {
					col := len(text)
					stop := (col - 1) / v.tabSize() * v.tabSize()
					start = ls + stop
				}
			}
			// Delete the pair of an auto-closed bracket
			end := pos
			if v.opts.AutoClose && size == 1 {
				next := v.doc.buf.ByteAt(pos)
				if cl, ok := closingOf[string(rune(r))]; ok && pos < v.doc.Len() && string(next) == cl {
					end = pos + 1
				}
			}
			edits = append(edits, selEdit{start: start, end: end})
		}
		v.applySelEdits(edits)
		v.keepRectAfterEdit()
	})
}

// DeleteForward deletes the selections or the characters after the carets
func (v *View) DeleteForward() {
	v.record("delete", "")
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			if !s.IsEmptyText() {
				edits = append(edits, selEdit{start: s.Start(), end: s.End()})
				continue
			}
			pos := s.Caret
			_, size := v.doc.buf.DecodeRune(pos)
			if s.CaretVS > 0 && v.doc.buf.ByteAt(pos) == '\n' {
				// Join the next line at the caret in the virtual space
				pad := strings.Repeat(" ", s.CaretVS)
				edits = append(edits, selEdit{start: pos, end: pos + 1, text: []byte(pad), anchor: len(pad), caret: len(pad)})
				continue
			}
			edits = append(edits, selEdit{start: pos, end: pos + size})
		}
		v.applySelEdits(edits)
		v.keepRectAfterEdit()
	})
}

// deleteTo deletes from each caret to the position target returns
func (v *View) deleteTo(target func(pos int) int) {
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			if !s.IsEmptyText() {
				edits = append(edits, selEdit{start: s.Start(), end: s.End()})
				continue
			}
			t := target(s.Caret)
			edits = append(edits, selEdit{start: min(t, s.Caret), end: max(t, s.Caret)})
		}
		v.st.sels.Rect = false
		v.applySelEdits(edits)
	})
}

// ---- Indentation ----

// selLines returns the lines of the selection; a selection ending at the
// start of a line does not take that line
func (v *View) selLines(s Sel) (int, int) {
	l1 := v.doc.LineOfOffset(s.Start())
	l2 := v.doc.LineOfOffset(s.End())
	if l2 > l1 && s.End() == v.doc.LineStart(l2) {
		l2--
	}
	return l1, l2
}

// Tab indents the selected lines, or inserts a tab at the carets
func (v *View) Tab() {
	multiLine := false
	for _, s := range v.st.sels.List {
		l1, l2 := v.doc.LineOfOffset(s.Start()), v.doc.LineOfOffset(s.End())
		if l1 != l2 && !v.st.sels.Rect {
			multiLine = true
		}
	}
	if multiLine {
		v.IndentLines(true)
		return
	}
	v.record("tab", "")
	if v.opts.InsertTabs {
		v.typeRaw("\t")
		return
	}
	// Spaces up to the next tab stop of each caret
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			col := v.colOfPos(s.Start()) + s.CaretVS
			if !s.IsEmptyText() {
				col = v.colOfPos(s.Start())
			}
			pad := strings.Repeat(" ", s.CaretVS)
			if !s.IsEmptyText() {
				pad = ""
			}
			ins := pad + v.indentUnit(col)
			edits = append(edits, selEdit{start: s.Start(), end: s.End(), text: []byte(ins), anchor: len(ins), caret: len(ins)})
		}
		v.applySelEdits(edits)
		v.keepRectAfterEdit()
	})
}

// typeRaw inserts text at the carets without auto-closing
func (v *View) typeRaw(text string) {
	v.edit(func() {
		v.doc.SetTyping()
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			pad := ""
			if s.IsEmptyText() {
				pad = strings.Repeat(" ", s.CaretVS)
			}
			ins := pad + text
			edits = append(edits, selEdit{start: s.Start(), end: s.End(), text: []byte(ins), anchor: len(ins), caret: len(ins)})
		}
		v.applySelEdits(edits)
		v.keepRectAfterEdit()
	})
}

// IndentLines indents (or unindents) the lines of the selections
func (v *View) IndentLines(indent bool) {
	if indent {
		v.record("indent", "")
	} else {
		v.record("unindent", "")
	}
	v.edit(func() {
		seen := map[int]bool{}
		var edits []Edit
		for _, s := range v.sortedSels() {
			l1, l2 := v.selLines(s)
			for l := l1; l <= l2; l++ {
				if seen[l] {
					continue
				}
				seen[l] = true
				ls := v.doc.LineStart(l)
				text := v.doc.buf.View(ls, v.doc.LineEnd(l))
				lead := leadingIndent(text)
				if indent {
					if len(text) == 0 {
						continue
					}
					edits = append(edits, Edit{Pos: ls, Text: []byte(v.indentUnit(0))})
					continue
				}
				// Remove one level: a tab or up to a tab size of spaces
				n := 0
				col := 0
				for n < len(lead) && col < v.tabSize() {
					if lead[n] == '\t' {
						n++
						break
					}
					n++
					col++
				}
				if n > 0 {
					edits = append(edits, Edit{Pos: ls, Len: n})
				}
			}
		}
		v.doc.ApplyEdits(edits)
		// A selection starting at the start of a line keeps starting there
		for i := range v.st.sels.List {
			s := &v.st.sels.List[i]
			if s.IsEmptyText() {
				continue
			}
			if s.Anchor < s.Caret {
				s.Anchor = v.doc.LineStart(v.doc.LineOfOffset(s.Anchor))
			} else {
				s.Caret = v.doc.LineStart(v.doc.LineOfOffset(s.Caret))
			}
		}
	})
}

// ---- Clipboard ----

// rectClipboard is the text last copied from a rectangular selection: pasted
// back, it is pasted as a rectangle
var rectClipboard string

// lineClipboard is the text last copied as whole lines (a copy without a selection)
var lineClipboard string

func (v *View) toClipboard(text string) {
	if v.doc.EOL != EOLLF {
		text = strings.ReplaceAll(text, "\n", string(v.doc.EOL.Bytes()))
	}
	ui.ClipboardSetText(text)
}

// Copy copies the selections; without a selection the line of the caret
func (v *View) Copy() {
	v.record("copy", "")
	text, whole := v.copyText()
	rectClipboard, lineClipboard = "", ""
	if v.st.sels.Rect {
		rectClipboard = text
	}
	if whole {
		lineClipboard = text
	}
	v.toClipboard(text)
}

// copyText returns the text to copy and whether it is whole lines
func (v *View) copyText() (string, bool) {
	st := v.st
	if !v.HasSelection() && !st.sels.Rect {
		if !v.opts.CopyLineNoSel {
			return "", false
		}
		var b strings.Builder
		seen := map[int]bool{}
		for _, s := range v.sortedSels() {
			line := v.doc.LineOfOffset(s.Caret)
			if seen[line] {
				continue
			}
			seen[line] = true
			b.Write(v.doc.LineText(line))
			b.WriteByte('\n')
		}
		return b.String(), true
	}
	var b strings.Builder
	for i, s := range v.sortedSels() {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.Write(v.doc.Text(s.Start(), s.End()))
	}
	if st.sels.Rect {
		b.WriteByte('\n')
	}
	return b.String(), false
}

// Cut copies and deletes the selections; without a selection the line of the caret
func (v *View) Cut() {
	if v.readOnly() {
		v.Copy()
		return
	}
	v.record("cut", "")
	text, whole := v.copyText()
	rectClipboard, lineClipboard = "", ""
	if v.st.sels.Rect {
		rectClipboard = text
	}
	if whole {
		lineClipboard = text
		v.toClipboard(text)
		v.DeleteLines()
		return
	}
	v.toClipboard(text)
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			edits = append(edits, selEdit{start: s.Start(), end: s.End()})
		}
		v.applySelEdits(edits)
		v.keepRectAfterEdit()
	})
}

// Paste pastes the clipboard at the carets
func (v *View) Paste() {
	text, err := ui.ClipboardGetText()
	if err != nil || text == "" {
		return
	}
	v.PasteText(text)
}

// PasteText pastes the text as if from the clipboard
func (v *View) PasteText(text string) {
	text = normalizeEOL(text)
	v.record("paste-text", text)
	st := v.st
	if rectClipboard != "" && normalizeEOL(rectClipboard) == text && (len(st.sels.List) == 1 || st.sels.Rect) {
		v.pasteRect(text)
		return
	}
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		distribute := len(sels) > 1 && len(lines) == len(sels)
		for i, s := range sels {
			ins := text
			if distribute {
				ins = lines[i]
			}
			// Whole lines copied without a selection go above the caret line
			if !distribute && lineClipboard != "" && normalizeEOL(lineClipboard) == text && s.IsEmptyText() {
				ls := v.doc.LineStart(v.doc.LineOfOffset(s.Caret))
				off := s.Caret - ls + len(ins)
				edits = append(edits, selEdit{start: ls, end: ls, text: []byte(ins), anchor: off, caret: off})
				continue
			}
			pad := ""
			if s.IsEmptyText() {
				pad = strings.Repeat(" ", s.CaretVS)
			}
			ins = pad + ins
			edits = append(edits, selEdit{start: s.Start(), end: s.End(), text: []byte(ins), anchor: len(ins), caret: len(ins)})
		}
		st.sels.Rect = false
		v.applySelEdits(edits)
	})
}

// pasteRect pastes the lines of the text as a rectangle at the caret column
func (v *View) pasteRect(text string) {
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	v.edit(func() {
		st := v.st
		m := st.sels.MainSel()
		if st.sels.Rect {
			m = st.sels.List[0]
			// Replace the rectangle first
			sels := v.sortedSels()
			edits := make([]selEdit, 0, len(sels))
			for _, s := range sels {
				edits = append(edits, selEdit{start: s.Start(), end: s.End(), vs: 0})
			}
			v.applySelEdits(edits)
			m = v.st.sels.List[0]
		} else if !m.IsEmptyText() {
			v.doc.Delete(m.Start(), m.End()-m.Start())
			m = caretSel(m.Start())
		}
		line := v.doc.LineOfOffset(m.Caret)
		col := v.colOfPos(m.Caret) + m.CaretVS
		var edits []Edit
		extra := ""
		for i, l := range lines {
			ln := line + i
			if ln >= v.doc.LineCount() {
				// Past the end: new lines
				extra += "\n" + strings.Repeat(" ", col) + l
				continue
			}
			pos, beyond := v.posOfCol(ln, col, false)
			edits = append(edits, Edit{Pos: pos, Text: []byte(strings.Repeat(" ", beyond) + l)})
		}
		if extra != "" {
			edits = append(edits, Edit{Pos: v.doc.Len(), Text: []byte(extra)})
		}
		v.doc.ApplyEdits(edits)
		last := min(line+len(lines)-1, v.doc.LineCount()-1)
		p, _ := v.posOfCol(last, col+utf8.RuneCountInString(lines[len(lines)-1]), false)
		st.sels.setSingle(caretSel(p))
	})
}

// ---- Undo ----

// Undo reverts the last change
func (v *View) Undo() {
	v.record("undo", "")
	sel, pos, ok := v.doc.Undo()
	if !ok {
		return
	}
	v.restoreAfterUndo(sel, pos)
}

// Redo makes the last undone change again
func (v *View) Redo() {
	v.record("redo", "")
	sel, pos, ok := v.doc.Redo()
	if !ok {
		return
	}
	v.restoreAfterUndo(sel, pos)
}

func (v *View) restoreAfterUndo(sel any, pos int) {
	if s, ok := sel.(*Selections); ok && s != nil {
		c := s.clone()
		n := v.doc.Len()
		for i := range c.List {
			c.List[i].Anchor = max(0, min(c.List[i].Anchor, n))
			c.List[i].Caret = max(0, min(c.List[i].Caret, n))
		}
		v.st.sels = c
	} else {
		v.st.sels.setSingle(caretSel(pos))
	}
	v.EnsureCaretVisible(true)
	v.selectionChanged()
	if v.OnModified != nil {
		v.OnModified()
	}
}

// ---- Caret moves ----

// moveCarets moves every caret with f; extend keeps the anchors
func (v *View) moveCarets(extend bool, f func(s Sel) (pos, vs int)) {
	st := v.st
	if st.sels.Rect && !extend {
		// Leaving a rectangular selection: one caret
		m := st.sels.MainSel()
		st.sels.setSingle(m)
	}
	st.sels.Rect = false
	for i := range st.sels.List {
		s := &st.sels.List[i]
		pos, vs := f(*s)
		s.Caret, s.CaretVS = pos, vs
		if !extend {
			s.Anchor, s.AnchorVS = pos, vs
		}
	}
	st.sels.normalize()
	v.ensureCaretVisible()
	v.selectionChanged()
}

// collapse moves the carets of non-empty selections to one of their ends
func (v *View) collapse(toStart bool) bool {
	if !v.HasSelection() || v.st.sels.Rect {
		return false
	}
	v.moveCarets(false, func(s Sel) (int, int) {
		if toStart {
			return s.Start(), 0
		}
		return s.End(), 0
	})
	return true
}

func (v *View) charLeft(pos int) int {
	if pos <= 0 {
		return 0
	}
	_, size := v.doc.buf.DecodeLastRune(pos)
	p := pos - max(1, size)
	// Skip the combining marks
	for p > 0 {
		r, s := v.doc.buf.DecodeRune(p)
		if r < 0x300 || !isZeroWidth(r) {
			break
		}
		_ = s
		_, size = v.doc.buf.DecodeLastRune(p)
		p -= max(1, size)
	}
	return p
}

func (v *View) charRight(pos int) int {
	if pos >= v.doc.Len() {
		return v.doc.Len()
	}
	_, size := v.doc.buf.DecodeRune(pos)
	p := pos + max(1, size)
	for p < v.doc.Len() {
		r, s := v.doc.buf.DecodeRune(p)
		if r < 0x300 || !isZeroWidth(r) {
			break
		}
		p += s
	}
	return p
}

// WordLeft returns the start of the word before pos
func (v *View) WordLeft(pos int) int {
	buf := v.doc.buf
	if pos <= 0 {
		return 0
	}
	p := pos
	r, size := buf.DecodeLastRune(p)
	if r == '\n' {
		return p - 1
	}
	for p > 0 {
		r, size = buf.DecodeLastRune(p)
		if charClass(r, v.opts.WordChars) != ccSpace {
			break
		}
		p -= size
	}
	if p == 0 {
		return 0
	}
	r, _ = buf.DecodeLastRune(p)
	cls := charClass(r, v.opts.WordChars)
	if cls == ccLineEnd {
		return p
	}
	for p > 0 {
		r, size = buf.DecodeLastRune(p)
		if charClass(r, v.opts.WordChars) != cls {
			break
		}
		p -= size
	}
	return p
}

// WordRight returns the start of the word after pos
func (v *View) WordRight(pos int) int {
	buf := v.doc.buf
	n := buf.Len()
	if pos >= n {
		return n
	}
	p := pos
	r, size := buf.DecodeRune(p)
	if r == '\n' {
		return p + 1
	}
	cls := charClass(r, v.opts.WordChars)
	if cls != ccSpace {
		for p < n {
			r, size = buf.DecodeRune(p)
			if charClass(r, v.opts.WordChars) != cls {
				break
			}
			p += size
		}
	}
	for p < n {
		r, size = buf.DecodeRune(p)
		if charClass(r, v.opts.WordChars) != ccSpace {
			break
		}
		p += size
	}
	return p
}

// wordEnd returns the end of the word at or after pos, for Ctrl+Delete
func (v *View) wordEnd(pos int) int {
	buf := v.doc.buf
	n := buf.Len()
	if pos >= n {
		return n
	}
	r, size := buf.DecodeRune(pos)
	if r == '\n' {
		return pos + 1
	}
	cls := charClass(r, v.opts.WordChars)
	p := pos
	for p < n {
		r, size = buf.DecodeRune(p)
		if charClass(r, v.opts.WordChars) != cls {
			break
		}
		p += size
	}
	if cls == ccSpace {
		return p
	}
	for p < n {
		r, size = buf.DecodeRune(p)
		if charClass(r, v.opts.WordChars) != ccSpace {
			break
		}
		p += size
	}
	return p
}

// homePos returns where Home moves the caret: the first non-blank
// character, or the line start when already there (or the wrapped row start)
func (v *View) homePos(pos int) int {
	line := v.doc.LineOfOffset(pos)
	ls := v.doc.LineStart(line)
	if v.wrapping() {
		sub := v.rowOfOffset(line, pos-ls)
		if sub > 0 {
			rs, _ := v.rowRange(line, sub, v.lineLen(line))
			if pos != ls+rs {
				return ls + rs
			}
		}
	}
	if !v.opts.SmartHome {
		return ls
	}
	first := ls + len(leadingIndent(v.lineHead(line, 4096)))
	if pos == first {
		return ls
	}
	return first
}

// endPos returns where End moves the caret: the end of the wrapped row, then of the line
func (v *View) endPos(pos int) int {
	line := v.doc.LineOfOffset(pos)
	ls := v.doc.LineStart(line)
	le := v.doc.LineEnd(line)
	if v.wrapping() {
		sub := v.rowOfOffset(line, pos-ls)
		_, re := v.rowRange(line, sub, le-ls)
		if ls+re != le && pos != ls+re {
			// Before the first character of the next row
			return ls + re
		}
	}
	return le
}

// vertical returns the position n rows below (above for n < 0) the caret,
// keeping its column; virtual allows the virtual space
func (v *View) vertical(s Sel, n int, virtual bool) (int, int, int) {
	pos := s.Caret
	line := v.doc.LineOfOffset(pos)
	ls := v.doc.LineStart(line)
	sub := v.rowOfOffset(line, pos-ls)
	rs, _ := v.rowRange(line, sub, v.lineLen(line))
	want := s.WantCol
	if want < 0 {
		want = v.colAt(line, pos-ls) - v.colAt(line, rs) + s.CaretVS
		if !v.wrapping() {
			want = v.colAt(line, pos-ls) + s.CaretVS
		}
		if sub > 0 {
			want += v.wrap(line).indent
		}
	}
	v.st.rebuildHidden(v.tabSize())
	l2, s2 := v.stepRows(v.st.visibleLine(line), sub, n)
	if l2 == line && s2 == sub {
		// No row to go to: the start or the end of the document
		if n < 0 {
			return 0, 0, want
		}
		return v.doc.Len(), 0, want
	}
	len2 := v.lineLen(l2)
	rs2, re2 := v.rowRange(l2, s2, len2)
	col := want
	if s2 > 0 {
		col -= v.wrap(l2).indent
	}
	col = max(0, col)
	if v.wrapping() {
		col += v.colAt(l2, rs2)
	}
	off, beyond := v.offAt(l2, col, false)
	if off > re2 || (off == re2 && re2 < len2) {
		// Stay in the row: before the first character of the next one
		off = re2
		if re2 < len2 {
			off = v.charLeft(v.doc.LineStart(l2)+re2) - v.doc.LineStart(l2)
			off = max(off, rs2)
		}
		beyond = 0
	}
	if !virtual {
		beyond = 0
	}
	return v.doc.LineStart(l2) + off, beyond, want
}

func (v *View) moveVertical(n int, extend bool) {
	st := v.st
	if !extend && v.HasSelection() && len(st.sels.List) == 1 && !st.sels.Rect {
		// Up from the start of the selection, down from its end
		s := st.sels.MainSel()
		if n < 0 {
			st.sels.List[0] = Sel{Anchor: s.Start(), Caret: s.Start(), WantCol: -1}
		} else {
			st.sels.List[0] = Sel{Anchor: s.End(), Caret: s.End(), WantCol: -1}
		}
	}
	if st.sels.Rect && !extend {
		st.sels.setSingle(st.sels.MainSel())
	}
	st.sels.Rect = false
	for i := range st.sels.List {
		s := &st.sels.List[i]
		pos, vs, want := v.vertical(*s, n, false)
		s.Caret, s.CaretVS, s.WantCol = pos, vs, want
		if !extend {
			s.Anchor, s.AnchorVS = pos, vs
		}
	}
	st.sels.normalize()
	v.ensureCaretVisible()
	v.selectionChanged()
}

// ---- Rectangular selection ----

// rebuildRect makes the selections of the rows of the rectangular selection
func (v *View) rebuildRect() {
	st := v.st
	sl := &st.sels
	l1, l2 := min(sl.RectAnchorLine, sl.RectCaretLine), max(sl.RectAnchorLine, sl.RectCaretLine)
	c1, c2 := min(sl.RectAnchorCol, sl.RectCaretCol), max(sl.RectAnchorCol, sl.RectCaretCol)
	caretRight := sl.RectCaretCol >= sl.RectAnchorCol
	sl.List = sl.List[:0]
	st.rebuildHidden(v.tabSize())
	for l := l1; l <= l2 && l < v.doc.LineCount(); l++ {
		if st.isHidden(l) {
			continue
		}
		a, avs := v.posOfCol(l, c1, false)
		b, bvs := v.posOfCol(l, c2, false)
		s := Sel{Anchor: a, Caret: b, AnchorVS: avs, CaretVS: bvs, WantCol: -1}
		if !caretRight {
			s = Sel{Anchor: b, Caret: a, AnchorVS: bvs, CaretVS: avs, WantCol: -1}
		}
		sl.List = append(sl.List, s)
		if l == sl.RectCaretLine {
			sl.Main = len(sl.List) - 1
		}
	}
	if len(sl.List) == 0 {
		sl.List = append(sl.List, caretSel(0))
	}
	sl.Main = max(0, min(sl.Main, len(sl.List)-1))
}

// extendRect moves the caret corner of the rectangular selection
func (v *View) extendRect(dLine, dCol int, toHome, toEnd bool, page int) {
	st := v.st
	sl := &st.sels
	if !sl.Rect {
		m := sl.MainSel()
		sl.RectAnchorLine = v.doc.LineOfOffset(m.Anchor)
		sl.RectAnchorCol = v.colOfPos(m.Anchor) + m.AnchorVS
		sl.RectCaretLine = v.doc.LineOfOffset(m.Caret)
		sl.RectCaretCol = v.colOfPos(m.Caret) + m.CaretVS
		sl.Rect = true
	}
	line := sl.RectCaretLine
	for i := 0; i < abs(dLine+page); i++ {
		var n int
		if dLine+page > 0 {
			n = st.nextVisible(line)
		} else {
			n = st.prevVisible(line)
		}
		if n < 0 {
			break
		}
		line = n
	}
	sl.RectCaretLine = line
	sl.RectCaretCol = max(0, sl.RectCaretCol+dCol)
	if toHome {
		sl.RectCaretCol = 0
	}
	if toEnd {
		sl.RectCaretCol = v.lineColsOf(line)
	}
	v.rebuildRect()
	v.ensureVisible(sl.MainSel().Caret, sl.MainSel().CaretVS)
	v.selectionChanged()
}

// SelectAll selects the whole text
func (v *View) SelectAll() {
	v.st.sels.setSingle(Sel{Anchor: 0, Caret: v.doc.Len(), WantCol: -1})
	v.selectionChanged()
}
