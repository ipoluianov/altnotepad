package editor

import (
	"bytes"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MacroStep is a command of a recorded macro
type MacroStep struct {
	Cmd string
	Arg string
}

// Macro is a recorded sequence of commands
type Macro struct {
	Steps []MacroStep
}

// StartRecording records the commands of the view into the macro
func (v *View) StartRecording(m *Macro) { v.macro = &MacroRecorder{macro: m} }

// StopRecording stops recording
func (v *View) StopRecording() { v.macro = nil }

// IsRecording reports whether a macro is being recorded
func (v *View) IsRecording() bool { return v.macro != nil }

// MacroRecorder records the commands of a view
type MacroRecorder struct {
	macro *Macro
}

// RecordStep adds a step to the macro being recorded, e.g. a search made by the application
func (v *View) RecordStep(cmd, arg string) { v.record(cmd, arg) }

func (v *View) record(cmd, arg string) {
	if v.macro == nil || v.macro.macro == nil {
		return
	}
	steps := &v.macro.macro.Steps
	if cmd == "type" && len(*steps) > 0 && (*steps)[len(*steps)-1].Cmd == "type" {
		(*steps)[len(*steps)-1].Arg += arg
		return
	}
	*steps = append(*steps, MacroStep{cmd, arg})
}

// OnMacroCommand runs the steps of a macro the view does not know, e.g. a
// search; it returns false if it does not know them either
var OnMacroCommand func(v *View, cmd, arg string) bool

// Play runs the macro times times as one undo step (times < 0: up to the end of the document)
func (v *View) Play(m *Macro, times int) {
	if m == nil || len(m.Steps) == 0 {
		return
	}
	rec := v.macro
	v.macro = nil
	defer func() { v.macro = rec }()
	v.doc.BeginAction()
	defer v.doc.EndAction()
	untilEnd := times < 0
	for i := 0; untilEnd || i < times; i++ {
		before := v.Caret()
		ver := v.doc.version
		for _, s := range m.Steps {
			if !v.Exec(s.Cmd, s.Arg) && OnMacroCommand != nil {
				OnMacroCommand(v, s.Cmd, s.Arg)
			}
		}
		if untilEnd {
			// Stop at the end of the document, or when nothing moves any more
			if v.Caret() >= v.doc.Len() || (v.Caret() == before && v.doc.version == ver) || i > 1_000_000 {
				break
			}
		}
	}
}

// commands are the commands of the editor by name
var commands map[string]func(v *View, arg string)

func init() {
	moves := func(name string, f func(v *View, s Sel) (int, int)) {
		commands[name] = func(v *View, _ string) {
			if (name == "char-left" || name == "char-right") && v.collapse(name == "char-left") {
				return
			}
			v.moveCarets(false, func(s Sel) (int, int) { return f(v, s) })
		}
		commands[name+"-extend"] = func(v *View, _ string) {
			v.moveCarets(true, func(s Sel) (int, int) { return f(v, s) })
		}
	}
	commands = map[string]func(v *View, arg string){}
	moves("char-left", func(v *View, s Sel) (int, int) {
		if s.CaretVS > 0 {
			return s.Caret, s.CaretVS - 1
		}
		return v.charLeft(s.Caret), 0
	})
	moves("char-right", func(v *View, s Sel) (int, int) { return v.charRight(s.Caret), 0 })
	moves("word-left", func(v *View, s Sel) (int, int) { return v.WordLeft(s.Caret), 0 })
	moves("word-right", func(v *View, s Sel) (int, int) { return v.WordRight(s.Caret), 0 })
	moves("home", func(v *View, s Sel) (int, int) { return v.homePos(s.Caret), 0 })
	moves("end", func(v *View, s Sel) (int, int) { return v.endPos(s.Caret), 0 })
	moves("doc-start", func(v *View, s Sel) (int, int) { return 0, 0 })
	moves("doc-end", func(v *View, s Sel) (int, int) { return v.doc.Len(), 0 })
	commands["line-up"] = func(v *View, _ string) { v.moveVertical(-1, false) }
	commands["line-down"] = func(v *View, _ string) { v.moveVertical(1, false) }
	commands["line-up-extend"] = func(v *View, _ string) { v.moveVertical(-1, true) }
	commands["line-down-extend"] = func(v *View, _ string) { v.moveVertical(1, true) }
	page := func(v *View, dir int, extend bool) {
		n := v.visibleRows() - 1
		v.st.topLine, v.st.topSub = v.stepRows(v.st.topLine, v.st.topSub, dir*n)
		v.clampScroll()
		v.moveVertical(dir*n, extend)
	}
	commands["page-up"] = func(v *View, _ string) { page(v, -1, false) }
	commands["page-down"] = func(v *View, _ string) { page(v, 1, false) }
	commands["page-up-extend"] = func(v *View, _ string) { page(v, -1, true) }
	commands["page-down-extend"] = func(v *View, _ string) { page(v, 1, true) }
	commands["scroll-up"] = func(v *View, _ string) { v.ScrollLines(-1) }
	commands["scroll-down"] = func(v *View, _ string) { v.ScrollLines(1) }

	// Rectangular selection
	commands["char-left-rect"] = func(v *View, _ string) { v.extendRect(0, -1, false, false, 0) }
	commands["char-right-rect"] = func(v *View, _ string) { v.extendRect(0, 1, false, false, 0) }
	commands["line-up-rect"] = func(v *View, _ string) { v.extendRect(-1, 0, false, false, 0) }
	commands["line-down-rect"] = func(v *View, _ string) { v.extendRect(1, 0, false, false, 0) }
	commands["home-rect"] = func(v *View, _ string) { v.extendRect(0, 0, true, false, 0) }
	commands["end-rect"] = func(v *View, _ string) { v.extendRect(0, 0, false, true, 0) }
	commands["page-up-rect"] = func(v *View, _ string) { v.extendRect(0, 0, false, false, -(v.visibleRows() - 1)) }
	commands["page-down-rect"] = func(v *View, _ string) { v.extendRect(0, 0, false, false, v.visibleRows()-1) }
	commands["word-left-rect"] = commands["char-left-rect"]
	commands["word-right-rect"] = commands["char-right-rect"]
	commands["doc-start-rect"] = commands["home-rect"]
	commands["doc-end-rect"] = commands["end-rect"]

	// Editing
	commands["type"] = func(v *View, arg string) { v.TypeText(arg) }
	commands["insert-text"] = func(v *View, arg string) { v.typeRaw(normalizeEOL(arg)) }
	commands["newline"] = func(v *View, _ string) { v.Newline() }
	commands["backspace"] = func(v *View, _ string) { v.Backspace() }
	commands["delete"] = func(v *View, _ string) { v.DeleteForward() }
	commands["delete-word-left"] = func(v *View, _ string) { v.deleteTo(v.WordLeft) }
	commands["delete-word-right"] = func(v *View, _ string) { v.deleteTo(v.wordEnd) }
	commands["delete-line-left"] = func(v *View, _ string) {
		v.deleteTo(func(p int) int { return v.doc.LineStart(v.doc.LineOfOffset(p)) })
	}
	commands["delete-line-right"] = func(v *View, _ string) {
		v.deleteTo(func(p int) int { return v.doc.LineEnd(v.doc.LineOfOffset(p)) })
	}
	commands["tab"] = func(v *View, _ string) { v.Tab() }
	commands["backtab"] = func(v *View, _ string) { v.IndentLines(false) }
	commands["indent"] = func(v *View, _ string) { v.IndentLines(true) }
	commands["unindent"] = func(v *View, _ string) { v.IndentLines(false) }
	commands["toggle-overtype"] = func(v *View, _ string) {
		v.overtype = !v.overtype
		if v.OnOvertypeChanged != nil {
			v.OnOvertypeChanged()
		}
		v.update()
	}
	commands["cancel"] = func(v *View, _ string) {
		st := v.st
		if len(st.sels.List) > 1 || st.sels.Rect {
			m := st.sels.MainSel()
			st.sels.setSingle(caretSel(m.Caret))
			v.selectionChanged()
		}
	}
	commands["select-all"] = func(v *View, _ string) { v.SelectAll() }
	commands["copy"] = func(v *View, _ string) { v.Copy() }
	commands["cut"] = func(v *View, _ string) { v.Cut() }
	commands["paste"] = func(v *View, _ string) { v.Paste() }
	commands["paste-text"] = func(v *View, arg string) { v.PasteText(arg) }
	commands["complete"] = func(v *View, arg string) {
		if arg == "" {
			v.Complete(false)
			return
		}
		v.insertCompletion(v.wordStartAt(v.Caret()), arg)
	}
	commands["undo"] = func(v *View, _ string) { v.Undo() }
	commands["redo"] = func(v *View, _ string) { v.Redo() }

	// Lines
	commands["duplicate"] = func(v *View, _ string) { v.Duplicate() }
	commands["delete-lines"] = func(v *View, _ string) { v.DeleteLines() }
	commands["cut-lines"] = func(v *View, _ string) {
		text, _ := v.lineRangeText()
		lineClipboard = text
		v.toClipboard(text)
		v.DeleteLines()
	}
	commands["copy-lines"] = func(v *View, _ string) {
		text, _ := v.lineRangeText()
		lineClipboard = text
		v.toClipboard(text)
	}
	commands["move-lines-up"] = func(v *View, _ string) { v.MoveLines(-1) }
	commands["move-lines-down"] = func(v *View, _ string) { v.MoveLines(1) }
	commands["transpose"] = func(v *View, _ string) { v.TransposeLines() }
	commands["join-lines"] = func(v *View, _ string) { v.JoinLines() }
	commands["split-lines"] = func(v *View, _ string) { v.SplitLines() }
	commands["insert-line-above"] = func(v *View, _ string) { v.InsertBlankLine(true) }
	commands["insert-line-below"] = func(v *View, _ string) { v.InsertBlankLine(false) }
	for _, k := range []string{"asc", "desc", "asc-ci", "desc-ci", "int-asc", "int-desc", "dec-asc", "dec-desc", "len-asc", "len-desc"} {
		kind := k
		commands["sort-"+kind] = func(v *View, _ string) { v.SortLines(kind) }
	}
	commands["reverse-lines"] = func(v *View, _ string) {
		v.transformLines(func(lines [][]byte) [][]byte {
			for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
				lines[i], lines[j] = lines[j], lines[i]
			}
			return lines
		})
	}
	commands["shuffle-lines"] = func(v *View, _ string) {
		v.transformLines(func(lines [][]byte) [][]byte {
			rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
			return lines
		})
	}
	commands["remove-dup-lines"] = func(v *View, _ string) {
		v.transformLines(func(lines [][]byte) [][]byte {
			seen := make(map[string]bool, len(lines))
			out := lines[:0]
			for _, l := range lines {
				if !seen[string(l)] {
					seen[string(l)] = true
					out = append(out, l)
				}
			}
			return out
		})
	}
	commands["remove-consecutive-dup-lines"] = func(v *View, _ string) {
		v.transformLines(func(lines [][]byte) [][]byte {
			out := lines[:0]
			for i, l := range lines {
				if i == 0 || !bytes.Equal(l, lines[i-1]) {
					out = append(out, l)
				}
			}
			return out
		})
	}
	commands["remove-empty-lines"] = func(v *View, _ string) { v.removeLinesWhere(func(l []byte) bool { return len(l) == 0 }) }
	commands["remove-blank-lines"] = func(v *View, _ string) {
		v.removeLinesWhere(func(l []byte) bool { return len(bytes.TrimSpace(l)) == 0 })
	}

	// Blanks
	commands["trim-trailing"] = func(v *View, _ string) { v.trimLines(false, true) }
	commands["trim-leading"] = func(v *View, _ string) { v.trimLines(true, false) }
	commands["trim-both"] = func(v *View, _ string) { v.trimLines(true, true) }
	commands["eol-to-space"] = func(v *View, _ string) { v.EOLToSpace(false) }
	commands["trim-eol-to-space"] = func(v *View, _ string) { v.EOLToSpace(true) }
	commands["tabs-to-spaces"] = func(v *View, _ string) { v.TabsToSpaces() }
	commands["spaces-to-tabs-leading"] = func(v *View, _ string) { v.SpacesToTabs(true) }
	commands["spaces-to-tabs-all"] = func(v *View, _ string) { v.SpacesToTabs(false) }

	// Case
	for _, k := range []string{"upper", "lower", "proper", "proper-blend", "sentence", "sentence-blend", "invert", "random"} {
		kind := k
		commands["case-"+kind] = func(v *View, _ string) { v.ConvertCase(kind) }
	}

	// Comments
	commands["toggle-comment"] = func(v *View, _ string) { v.LineComment(0) }
	commands["comment"] = func(v *View, _ string) { v.LineComment(1) }
	commands["uncomment"] = func(v *View, _ string) { v.LineComment(-1) }
	commands["block-comment"] = func(v *View, _ string) { v.BlockComment(0) }
	commands["block-uncomment"] = func(v *View, _ string) { v.BlockComment(-1) }

	// Bookmarks
	commands["toggle-bookmark"] = func(v *View, _ string) {
		v.doc.Bookmarks.Toggle(v.doc.LineOfOffset(v.Caret()))
		v.update()
	}
	commands["toggle-bookmark-line"] = func(v *View, arg string) {
		if l, err := strconv.Atoi(arg); err == nil {
			v.doc.Bookmarks.Toggle(l)
			v.update()
		}
	}
	commands["next-bookmark"] = func(v *View, _ string) {
		if l := v.doc.Bookmarks.Next(v.doc.LineOfOffset(v.Caret())); l >= 0 {
			v.GotoLine(l)
		}
	}
	commands["prev-bookmark"] = func(v *View, _ string) {
		if l := v.doc.Bookmarks.Prev(v.doc.LineOfOffset(v.Caret())); l >= 0 {
			v.GotoLine(l)
		}
	}
	commands["clear-bookmarks"] = func(v *View, _ string) { v.doc.Bookmarks.Clear(); v.update() }
	commands["inverse-bookmarks"] = func(v *View, _ string) {
		var lines []int
		for l := 0; l < v.doc.LineCount(); l++ {
			if !v.doc.Bookmarks.Has(l) {
				lines = append(lines, l)
			}
		}
		v.doc.Bookmarks.Set(lines)
		v.update()
	}
	commands["copy-bookmarked-lines"] = func(v *View, _ string) { v.toClipboard(v.bookmarkedText()) }
	commands["cut-bookmarked-lines"] = func(v *View, _ string) {
		v.toClipboard(v.bookmarkedText())
		v.removeBookmarked(true)
	}
	commands["remove-bookmarked-lines"] = func(v *View, _ string) { v.removeBookmarked(true) }
	commands["remove-unbookmarked-lines"] = func(v *View, _ string) { v.removeBookmarked(false) }

	// Folding
	commands["fold-all"] = func(v *View, _ string) { v.FoldAll(true) }
	commands["unfold-all"] = func(v *View, _ string) { v.FoldAll(false) }
	commands["fold-current"] = func(v *View, _ string) { v.FoldCurrent(true) }
	commands["unfold-current"] = func(v *View, _ string) { v.FoldCurrent(false) }
	commands["toggle-fold"] = func(v *View, _ string) {
		v.toggleFoldAt(v.doc.LineOfOffset(v.Caret()), false)
	}
	commands["fold-level"] = func(v *View, arg string) {
		n, _ := strconv.Atoi(arg)
		v.FoldLevel(n, true)
	}
	commands["unfold-level"] = func(v *View, arg string) {
		n, _ := strconv.Atoi(arg)
		v.FoldLevel(n, false)
	}

	// Braces
	commands["goto-brace"] = func(v *View, _ string) { v.GotoBrace(false) }
	commands["select-to-brace"] = func(v *View, _ string) { v.GotoBrace(true) }
}

// Exec runs a command of the editor by name; returns false for an unknown one
func (v *View) Exec(cmd, arg string) bool {
	f, ok := commands[cmd]
	if !ok {
		return false
	}
	v.record(cmd, arg)
	rec := v.macro
	v.macro = nil // the command itself is not recorded again
	f(v, arg)
	v.macro = rec
	return true
}

// HasCommand reports whether the editor knows the command
func HasCommand(cmd string) bool {
	_, ok := commands[cmd]
	return ok
}

// ---- Lines ----

// lineRange returns the lines of all the selections
func (v *View) lineRange() (int, int) {
	l1, l2 := v.doc.LineCount(), 0
	for _, s := range v.st.sels.List {
		a, b := v.selLines(s)
		l1, l2 = min(l1, a), max(l2, b)
	}
	return l1, l2
}

func (v *View) lineRangeText() (string, bool) {
	l1, l2 := v.lineRange()
	return string(v.doc.Text(v.doc.LineStart(l1), v.lineEndWithBreak(l2))) + v.missingEOL(l2), true
}

// missingEOL returns a line break when the line is the last one, which has none
func (v *View) missingEOL(line int) string {
	if line == v.doc.LineCount()-1 {
		return "\n"
	}
	return ""
}

// Duplicate duplicates the selections, or the lines of the carets without one
func (v *View) Duplicate() {
	v.edit(func() {
		var edits []Edit
		seen := map[int]bool{}
		for _, s := range v.sortedSels() {
			if !s.IsEmptyText() {
				edits = append(edits, Edit{Pos: s.End(), Text: v.doc.Text(s.Start(), s.End())})
				continue
			}
			line := v.doc.LineOfOffset(s.Caret)
			if seen[line] {
				continue
			}
			seen[line] = true
			edits = append(edits, Edit{Pos: v.doc.LineEnd(line), Text: append([]byte{'\n'}, v.doc.LineText(line)...)})
		}
		v.doc.ApplyEdits(edits)
	})
}

// DeleteLines deletes the lines of the selections
func (v *View) DeleteLines() {
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
				start, end := v.doc.LineStart(l), v.lineEndWithBreak(l)
				if l == v.doc.LineCount()-1 && l > 0 {
					start-- // the last line: remove the break before it
				}
				edits = append(edits, Edit{Pos: start, Len: end - start})
			}
		}
		v.doc.ApplyEdits(edits)
		st := v.st
		m := st.sels.MainSel()
		st.sels.setSingle(caretSel(v.doc.LineStart(v.doc.LineOfOffset(m.Caret))))
	})
}

// MoveLines moves the lines of the selections up (dir -1) or down
func (v *View) MoveLines(dir int) {
	v.edit(func() {
		st := v.st
		l1, l2 := v.lineRange()
		if (dir < 0 && l1 == 0) || (dir > 0 && l2 >= v.doc.LineCount()-1) {
			return
		}
		keep := st.sels.clone()
		block := v.doc.Text(v.doc.LineStart(l1), v.doc.LineEnd(l2))
		shift := 0
		if dir < 0 {
			other := v.doc.LineText(l1 - 1)
			start := v.doc.LineStart(l1 - 1)
			v.doc.Replace(start, v.doc.LineEnd(l2)-start, bytes.Join([][]byte{block, other}, nlBytes))
			shift = -(len(other) + 1)
		} else {
			other := v.doc.LineText(l2 + 1)
			start := v.doc.LineStart(l1)
			v.doc.Replace(start, v.doc.LineEnd(l2+1)-start, bytes.Join([][]byte{other, block}, nlBytes))
			shift = len(other) + 1
		}
		for i := range keep.List {
			s := &keep.List[i]
			s.Anchor = max(0, min(v.doc.Len(), s.Anchor+shift))
			s.Caret = max(0, min(v.doc.Len(), s.Caret+shift))
		}
		if keep.Rect {
			keep.RectAnchorLine += dir
			keep.RectCaretLine += dir
		}
		st.sels = keep
	})
}

// TransposeLines swaps the line of the caret with the line above it
func (v *View) TransposeLines() {
	line := v.doc.LineOfOffset(v.Caret())
	if line == 0 {
		return
	}
	v.edit(func() {
		a, b := v.doc.LineText(line-1), v.doc.LineText(line)
		start := v.doc.LineStart(line - 1)
		col := v.Caret() - v.doc.LineStart(line)
		v.doc.Replace(start, len(a)+1+len(b), bytes.Join([][]byte{b, a}, nlBytes))
		v.st.sels.setSingle(caretSel(v.doc.LineStart(line) + min(col, len(a))))
	})
}

// JoinLines joins the selected lines (or the line with the next one) with spaces
func (v *View) JoinLines() {
	v.edit(func() {
		l1, l2 := v.lineRange()
		if l2 == l1 {
			l2 = min(l1+1, v.doc.LineCount()-1)
		}
		if l1 == l2 {
			return
		}
		var out []byte
		for l := l1; l <= l2; l++ {
			t := v.doc.LineText(l)
			if l > l1 {
				t = bytes.TrimLeft(t, " \t")
				out = bytes.TrimRight(out, " \t")
				if len(out) > 0 && len(t) > 0 {
					out = append(out, ' ')
				}
			}
			out = append(out, t...)
		}
		start := v.doc.LineStart(l1)
		v.doc.Replace(start, v.doc.LineEnd(l2)-start, out)
		v.st.sels.setSingle(caretSel(start + len(out)))
	})
}

// SplitLines breaks the selected lines at the long line marker column (or
// the width of the view) between words
func (v *View) SplitLines() {
	width := v.opts.EdgeColumn
	if width <= 0 {
		width = max(20, v.textW/v.font.CharWidth-1)
	}
	v.transformLines(func(lines [][]byte) [][]byte {
		var out [][]byte
		for _, l := range lines {
			for {
				if lineCols(l, v.tabSize()) <= width {
					out = append(out, l)
					break
				}
				off, _ := offsetOfCol(l, width, v.tabSize(), false)
				cut := bytes.LastIndexAny(l[:off+1], " \t")
				if cut <= 0 {
					out = append(out, l)
					break
				}
				out = append(out, l[:cut])
				l = bytes.TrimLeft(l[cut:], " \t")
				if len(l) == 0 {
					break
				}
			}
		}
		return out
	})
}

// InsertBlankLine inserts an empty line above or below the caret line
func (v *View) InsertBlankLine(above bool) {
	v.edit(func() {
		line := v.doc.LineOfOffset(v.Caret())
		ind := leadingIndent(v.doc.LineText(line))
		if !v.opts.AutoIndent {
			ind = nil
		}
		if above {
			pos := v.doc.LineStart(line)
			v.doc.Insert(pos, append(append([]byte{}, ind...), '\n'))
			v.st.sels.setSingle(caretSel(pos + len(ind)))
		} else {
			pos := v.doc.LineEnd(line)
			v.doc.Insert(pos, append([]byte{'\n'}, ind...))
			v.st.sels.setSingle(caretSel(pos + 1 + len(ind)))
		}
	})
}

// targetLines returns the lines an operation works on: those of the
// selection when it spans lines, otherwise all of them
func (v *View) targetLines() (int, int, bool) {
	if v.HasSelection() {
		l1, l2 := v.lineRange()
		if l2 > l1 || len(v.st.sels.List) > 1 {
			return l1, l2, false
		}
	}
	return 0, v.doc.LineCount() - 1, true
}

// transformLines replaces the target lines with what f makes of them
func (v *View) transformLines(f func(lines [][]byte) [][]byte) {
	v.edit(func() {
		l1, l2, all := v.targetLines()
		// The empty last line of the document stays last
		if all && l2 > 0 && v.doc.LineStart(l2) == v.doc.LineEnd(l2) {
			l2--
		}
		start, end := v.doc.LineStart(l1), v.doc.LineEnd(l2)
		text := v.doc.Text(start, end)
		lines := bytes.Split(text, nlBytes)
		out := bytes.Join(f(lines), nlBytes)
		v.doc.Replace(start, end-start, out)
		if all {
			v.st.sels.setSingle(caretSel(min(v.Caret(), v.doc.Len())))
		} else {
			v.st.sels.setSingle(Sel{Anchor: start, Caret: start + len(out), WantCol: -1})
		}
	})
}

// SortLines sorts the target lines; kind is asc, desc, asc-ci, desc-ci,
// int-asc, int-desc, dec-asc, dec-desc, len-asc, len-desc
func (v *View) SortLines(kind string) {
	v.transformLines(func(lines [][]byte) [][]byte {
		desc := strings.HasSuffix(kind, "desc")
		var less func(a, b []byte) bool
		switch {
		case strings.HasPrefix(kind, "int") || strings.HasPrefix(kind, "dec"):
			num := func(l []byte) (float64, bool) {
				f := strings.Fields(string(l))
				if len(f) == 0 {
					return 0, false
				}
				s := strings.ReplaceAll(f[0], ",", ".")
				x, err := strconv.ParseFloat(s, 64)
				return x, err == nil
			}
			less = func(a, b []byte) bool {
				x, okA := num(a)
				y, okB := num(b)
				if okA != okB {
					return !okA // lines without a number first
				}
				return x < y
			}
		case strings.HasPrefix(kind, "len"):
			less = func(a, b []byte) bool { return utf8.RuneCount(a) < utf8.RuneCount(b) }
		case strings.HasSuffix(kind, "-ci"):
			less = func(a, b []byte) bool { return bytes.Compare(bytes.ToLower(a), bytes.ToLower(b)) < 0 }
			desc = strings.HasPrefix(kind, "desc")
		default:
			less = func(a, b []byte) bool { return bytes.Compare(a, b) < 0 }
		}
		sort.SliceStable(lines, func(i, j int) bool {
			if desc {
				return less(lines[j], lines[i])
			}
			return less(lines[i], lines[j])
		})
		return lines
	})
}

// removeLinesWhere removes the target lines f is true for
func (v *View) removeLinesWhere(f func(l []byte) bool) {
	v.edit(func() {
		l1, l2, all := v.targetLines()
		var edits []Edit
		v.doc.buf.forEachLine(l1, func(line int, text []byte) bool {
			if line > l2 {
				return false
			}
			if f(text) {
				start, end := v.doc.LineStart(line), v.lineEndWithBreak(line)
				if line == v.doc.LineCount()-1 && line > 0 {
					start--
				}
				if end > start {
					edits = append(edits, Edit{Pos: start, Len: end - start})
				}
			}
			return true
		})
		// The edits of the last two lines may overlap: the later wins
		for i := 1; i < len(edits); i++ {
			if edits[i].Pos < edits[i-1].Pos+edits[i-1].Len {
				edits[i-1].Len = edits[i].Pos - edits[i-1].Pos
			}
		}
		v.doc.ApplyEdits(edits)
		_ = all
		v.st.sels.setSingle(caretSel(min(v.Caret(), v.doc.Len())))
	})
}

// trimLines removes the spaces at the start and/or the end of the target lines
func (v *View) trimLines(leading, trailing bool) {
	v.edit(func() {
		l1, l2, _ := v.targetLines()
		var edits []Edit
		v.doc.buf.forEachLine(l1, func(line int, text []byte) bool {
			if line > l2 {
				return false
			}
			ls := v.doc.LineStart(line)
			if trailing {
				t := bytes.TrimRight(text, " \t")
				if len(t) < len(text) {
					edits = append(edits, Edit{Pos: ls + len(t), Len: len(text) - len(t)})
				}
			}
			if leading {
				n := len(leadingIndent(text))
				if n > 0 && n < len(text) || (n > 0 && !trailing) {
					edits = append(edits, Edit{Pos: ls, Len: n})
				}
			}
			return true
		})
		v.doc.ApplyEdits(edits)
	})
}

// EOLToSpace joins the target lines with spaces; trim removes their spaces too
func (v *View) EOLToSpace(trim bool) {
	v.transformLines(func(lines [][]byte) [][]byte {
		parts := make([][]byte, 0, len(lines))
		for _, l := range lines {
			if trim {
				l = bytes.TrimSpace(l)
				if len(l) == 0 {
					continue
				}
			}
			parts = append(parts, l)
		}
		return [][]byte{bytes.Join(parts, []byte(" "))}
	})
}

// TabsToSpaces expands the tabs of the target lines
func (v *View) TabsToSpaces() {
	tab := v.tabSize()
	v.edit(func() {
		l1, l2, _ := v.targetLines()
		var edits []Edit
		v.doc.buf.forEachLine(l1, func(line int, text []byte) bool {
			if line > l2 {
				return false
			}
			if bytes.IndexByte(text, '\t') < 0 {
				return true
			}
			var out []byte
			col := 0
			for _, r := range string(text) {
				if r == '\t' {
					n := tab - col%tab
					out = append(out, bytes.Repeat([]byte{' '}, n)...)
					col += n
					continue
				}
				out = utf8.AppendRune(out, r)
				col += runeCells(r, utf8.RuneLen(r), col, tab)
			}
			edits = append(edits, Edit{Pos: v.doc.LineStart(line), Len: len(text), Text: out})
			return true
		})
		v.doc.ApplyEdits(edits)
	})
}

// SpacesToTabs turns the spaces reaching tab stops into tabs, in the leading
// indentation only or everywhere
func (v *View) SpacesToTabs(leadingOnly bool) {
	tab := v.tabSize()
	v.edit(func() {
		l1, l2, _ := v.targetLines()
		var edits []Edit
		v.doc.buf.forEachLine(l1, func(line int, text []byte) bool {
			if line > l2 {
				return false
			}
			if bytes.IndexByte(text, ' ') < 0 {
				return true
			}
			var out []byte
			col := 0
			spaces := 0
			inLead := true
			flush := func(toTabStop bool) {
				out = append(out, bytes.Repeat([]byte{' '}, spaces)...)
				spaces = 0
			}
			for _, r := range string(text) {
				if r == ' ' && (inLead || !leadingOnly) {
					spaces++
					col++
					if col%tab == 0 {
						if spaces > 1 {
							out = append(out, '\t')
						} else {
							out = append(out, ' ')
						}
						spaces = 0
					}
					continue
				}
				if r == '\t' {
					spaces = 0
					out = append(out, '\t')
					col += tab - col%tab
					continue
				}
				flush(false)
				inLead = false
				out = utf8.AppendRune(out, r)
				col += runeCells(r, utf8.RuneLen(r), col, tab)
			}
			flush(false)
			if !bytes.Equal(out, text) {
				edits = append(edits, Edit{Pos: v.doc.LineStart(line), Len: len(text), Text: out})
			}
			return true
		})
		v.doc.ApplyEdits(edits)
	})
}

// ---- Case ----

// ConvertCase converts the case of the selections (of the words at the
// carets without a selection)
func (v *View) ConvertCase(kind string) {
	v.edit(func() {
		sels := v.sortedSels()
		edits := make([]selEdit, 0, len(sels))
		for _, s := range sels {
			a, b := s.Start(), s.End()
			empty := a == b
			if empty {
				a, b = v.wordAt(a)
			}
			if a == b {
				edits = append(edits, selEdit{start: s.Caret, end: s.Caret, anchor: 0, caret: 0})
				continue
			}
			conv := []byte(convertCase(string(v.doc.Text(a, b)), kind))
			e := selEdit{start: a, end: b, text: conv, anchor: 0, caret: len(conv)}
			switch {
			case empty:
				off := min(s.Caret-a, len(conv))
				e.anchor, e.caret = off, off
			case s.Caret < s.Anchor:
				e.anchor, e.caret = len(conv), 0
			}
			edits = append(edits, e)
		}
		v.applySelEdits(edits)
	})
}

func convertCase(s, kind string) string {
	switch kind {
	case "upper":
		return strings.ToUpper(s)
	case "lower":
		return strings.ToLower(s)
	case "invert":
		return strings.Map(func(r rune) rune {
			if unicode.IsUpper(r) {
				return unicode.ToLower(r)
			}
			return unicode.ToUpper(r)
		}, s)
	case "random":
		return strings.Map(func(r rune) rune {
			if rand.Intn(2) == 0 {
				return unicode.ToLower(r)
			}
			return unicode.ToUpper(r)
		}, s)
	case "proper", "proper-blend":
		var b strings.Builder
		start := true
		for _, r := range s {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'' {
				if start {
					b.WriteRune(unicode.ToUpper(r))
				} else if kind == "proper" {
					b.WriteRune(unicode.ToLower(r))
				} else {
					b.WriteRune(r)
				}
				start = false
				continue
			}
			start = true
			b.WriteRune(r)
		}
		return b.String()
	case "sentence", "sentence-blend":
		var b strings.Builder
		start := true
		for _, r := range s {
			if unicode.IsLetter(r) {
				if start {
					b.WriteRune(unicode.ToUpper(r))
					start = false
				} else if kind == "sentence" {
					b.WriteRune(unicode.ToLower(r))
				} else {
					b.WriteRune(r)
				}
				continue
			}
			if r == '.' || r == '!' || r == '?' || r == '\n' {
				start = true
			}
			b.WriteRune(r)
		}
		return b.String()
	}
	return s
}

// ---- Comments ----

// LineComment comments (1), uncomments (-1) or toggles (0) the lines of the selections
func (v *View) LineComment(mode int) {
	lang := v.doc.lang
	if lang == nil {
		return
	}
	token := lang.LineComment
	if token == "" {
		if lang.BlockComment[0] != "" {
			v.blockCommentLines(mode)
		}
		return
	}
	v.edit(func() {
		l1, l2 := v.lineRange()
		lines := make([][]byte, 0, l2-l1+1)
		for l := l1; l <= l2; l++ {
			lines = append(lines, v.doc.LineText(l))
		}
		// Toggle: uncomment when all the non-blank lines are comments
		trimTok := strings.TrimRight(token, " ")
		isComment := func(t []byte) bool {
			return bytes.HasPrefix(bytes.TrimLeft(t, " \t"), []byte(trimTok))
		}
		if mode == 0 {
			mode = -1
			for _, t := range lines {
				if len(bytes.TrimSpace(t)) > 0 && !isComment(t) {
					mode = 1
					break
				}
			}
		}
		minInd := -1
		for _, t := range lines {
			if len(bytes.TrimSpace(t)) == 0 {
				continue
			}
			if n := len(leadingIndent(t)); minInd < 0 || n < minInd {
				minInd = n
			}
		}
		minInd = max(minInd, 0)
		var edits []Edit
		for i, t := range lines {
			ls := v.doc.LineStart(l1 + i)
			if len(bytes.TrimSpace(t)) == 0 {
				continue
			}
			if mode > 0 {
				ins := strings.TrimRight(token, " ") + " "
				edits = append(edits, Edit{Pos: ls + min(minInd, len(t)), Text: []byte(ins)})
				continue
			}
			if !isComment(t) {
				continue
			}
			k := len(leadingIndent(t))
			n := len(trimTok)
			if k+n < len(t) && t[k+n] == ' ' {
				n++
			}
			edits = append(edits, Edit{Pos: ls + k, Len: n})
		}
		v.doc.ApplyEdits(edits)
	})
}

// blockCommentLines comments each line with the block comment, for the
// languages without line comments (HTML, CSS)
func (v *View) blockCommentLines(mode int) {
	open, cls := v.doc.lang.BlockComment[0], v.doc.lang.BlockComment[1]
	v.edit(func() {
		l1, l2 := v.lineRange()
		var edits []Edit
		isComment := func(t []byte) bool {
			t = bytes.TrimSpace(t)
			return bytes.HasPrefix(t, []byte(open)) && bytes.HasSuffix(t, []byte(cls))
		}
		if mode == 0 {
			mode = -1
			for l := l1; l <= l2; l++ {
				t := v.doc.LineText(l)
				if len(bytes.TrimSpace(t)) > 0 && !isComment(t) {
					mode = 1
				}
			}
		}
		for l := l1; l <= l2; l++ {
			t := v.doc.LineText(l)
			ls := v.doc.LineStart(l)
			if len(bytes.TrimSpace(t)) == 0 {
				continue
			}
			k := len(leadingIndent(t))
			e := len(bytes.TrimRight(t, " \t"))
			if mode > 0 {
				edits = append(edits, Edit{Pos: ls + k, Text: []byte(open + " ")}, Edit{Pos: ls + e, Text: []byte(" " + cls)})
			} else if isComment(t) {
				n := len(open)
				if k+n < len(t) && t[k+n] == ' ' {
					n++
				}
				m := len(cls)
				if e-m-1 >= 0 && t[e-m-1] == ' ' {
					m++
				}
				edits = append(edits, Edit{Pos: ls + k, Len: n}, Edit{Pos: ls + e - m, Len: m})
			}
		}
		v.doc.ApplyEdits(edits)
	})
}

// BlockComment wraps the selections in the block comment (mode 0 toggles, -1 removes it)
func (v *View) BlockComment(mode int) {
	lang := v.doc.lang
	if lang == nil || lang.BlockComment[0] == "" {
		v.LineComment(mode)
		return
	}
	open, cls := lang.BlockComment[0], lang.BlockComment[1]
	v.edit(func() {
		var edits []Edit
		for _, s := range v.sortedSels() {
			a, b := s.Start(), s.End()
			if a == b {
				line := v.doc.LineOfOffset(a)
				t := v.doc.LineText(line)
				a = v.doc.LineStart(line) + len(leadingIndent(t))
				b = v.doc.LineStart(line) + len(bytes.TrimRight(t, " \t"))
			}
			text := v.doc.Text(a, b)
			inner := bytes.TrimSpace(text)
			wrapped := bytes.HasPrefix(inner, []byte(open)) && bytes.HasSuffix(inner, []byte(cls)) && len(inner) >= len(open)+len(cls)
			if mode == -1 || (mode == 0 && wrapped) {
				if !wrapped {
					continue
				}
				k := bytes.Index(text, []byte(open))
				e := bytes.LastIndex(text, []byte(cls))
				n, m := len(open), len(cls)
				if k+n < len(text) && text[k+n] == ' ' {
					n++
				}
				if e > 0 && text[e-1] == ' ' {
					e--
					m++
				}
				edits = append(edits, Edit{Pos: a + k, Len: n}, Edit{Pos: a + e, Len: m})
				continue
			}
			edits = append(edits, Edit{Pos: a, Text: []byte(open + " ")}, Edit{Pos: b, Text: []byte(" " + cls)})
		}
		v.doc.ApplyEdits(edits)
	})
}

// ---- Bookmarks ----

func (v *View) bookmarkedText() string {
	var b strings.Builder
	for _, l := range v.doc.Bookmarks.lines {
		if l < v.doc.LineCount() {
			b.Write(v.doc.LineText(l))
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// removeBookmarked removes the bookmarked lines, or the others
func (v *View) removeBookmarked(marked bool) {
	v.removeLinesWhere2(func(line int) bool { return v.doc.Bookmarks.Has(line) == marked })
}

func (v *View) removeLinesWhere2(f func(line int) bool) {
	v.edit(func() {
		var edits []Edit
		for line := 0; line < v.doc.LineCount(); line++ {
			if !f(line) {
				continue
			}
			start, end := v.doc.LineStart(line), v.lineEndWithBreak(line)
			if line == v.doc.LineCount()-1 && line > 0 {
				start--
			}
			edits = append(edits, Edit{Pos: start, Len: end - start})
		}
		for i := 1; i < len(edits); i++ {
			if edits[i].Pos < edits[i-1].Pos+edits[i-1].Len {
				edits[i-1].Len = edits[i].Pos - edits[i-1].Pos
			}
		}
		v.doc.ApplyEdits(edits)
		v.st.sels.setSingle(caretSel(min(v.Caret(), v.doc.Len())))
	})
}

// ---- Folding ----

func (v *View) foldChanged() {
	v.st.hiddenDirty = true
	v.st.rebuildHidden(v.tabSize())
	// The caret must not be hidden
	m := v.st.sels.MainSel()
	line := v.doc.LineOfOffset(m.Caret)
	if v.st.isHidden(line) {
		h := v.st.visibleLine(line)
		v.st.sels.setSingle(caretSel(v.doc.LineEnd(h)))
	}
	v.clampScroll()
	v.selectionChanged()
}

// toggleFoldAt folds or unfolds the block of the line (all the blocks in it with children)
func (v *View) toggleFoldAt(line int, children bool) {
	st := v.st
	tab := v.tabSize()
	if !st.IsFoldHeader(line, tab) {
		h := st.headerOf(line, tab)
		if h < 0 {
			return
		}
		line = h
	}
	folding := !st.folded.Has(line)
	if folding {
		st.folded.Add(line)
	} else {
		st.folded.Remove(line)
	}
	if children {
		end := st.foldBlockEnd(line, tab)
		for l := line + 1; l <= end; l++ {
			if st.IsFoldHeader(l, tab) {
				if folding {
					st.folded.Add(l)
				} else {
					st.folded.Remove(l)
				}
			}
		}
	}
	v.foldChanged()
}

// FoldAll folds (or unfolds) all the blocks
func (v *View) FoldAll(fold bool) {
	st := v.st
	if !fold {
		st.folded.Clear()
		v.foldChanged()
		return
	}
	var lines []int
	tab := v.tabSize()
	for l := 0; l < v.doc.LineCount(); l++ {
		if st.IsFoldHeader(l, tab) {
			lines = append(lines, l)
		}
	}
	st.folded.Set(lines)
	v.foldChanged()
}

// FoldCurrent folds (or unfolds) the block holding the caret
func (v *View) FoldCurrent(fold bool) {
	st := v.st
	line := v.doc.LineOfOffset(v.Caret())
	h := st.headerOf(line, v.tabSize())
	if h < 0 {
		return
	}
	if fold {
		st.folded.Add(h)
	} else {
		st.folded.Remove(h)
	}
	v.foldChanged()
}

// FoldLevel folds (or unfolds) the blocks at the depth level (1 - the outermost)
func (v *View) FoldLevel(level int, fold bool) {
	st := v.st
	tab := v.tabSize()
	if st.indentFold() {
		// The depth of a block: the number of the headers it is in
		var stack []int // indents of the open headers
		for l := 0; l < v.doc.LineCount(); l++ {
			ind := st.lineIndent(l, tab)
			if ind < 0 {
				continue
			}
			for len(stack) > 0 && stack[len(stack)-1] >= ind {
				stack = stack[:len(stack)-1]
			}
			if st.IsFoldHeader(l, tab) {
				if len(stack)+1 == level {
					if fold {
						st.folded.Add(l)
					} else {
						st.folded.Remove(l)
					}
				}
				stack = append(stack, ind)
			}
		}
	} else {
		for l := 0; l < v.doc.LineCount(); l++ {
			if !st.IsFoldHeader(l, tab) {
				continue
			}
			lvl := v.doc.hl.Level(l)
			f := v.doc.hl.fold(l)
			if lvl+f.min()+1 == level {
				if fold {
					st.folded.Add(l)
				} else {
					st.folded.Remove(l)
				}
			}
		}
	}
	v.foldChanged()
}

// GotoBrace moves the caret to the bracket matching the one at it; select
// selects the text between them too
func (v *View) GotoBrace(sel bool) {
	pos := v.Caret()
	a := -1
	for _, p := range []int{pos - 1, pos} {
		if p >= 0 && p < v.doc.Len() && isBrace(v.doc.buf.ByteAt(p)) {
			a = p
			break
		}
	}
	if a < 0 {
		return
	}
	b := v.FindMatchingBrace(a)
	if b < 0 {
		return
	}
	if sel {
		v.SetSelection(min(a, b), max(a, b)+1)
	} else if b > a {
		v.SetCaret(b + 1)
	} else {
		v.SetCaret(b)
	}
	v.EnsureCaretVisible(false)
}

// ---- Multi-selection of the occurrences ----

// occurrenceText returns the text whose occurrences are selected: the main
// selection or the word at the caret (then selected first); whole tells
// whether only whole words match
func (v *View) occurrenceText() (string, bool) {
	s := v.st.sels.MainSel()
	if s.IsEmptyText() {
		a, b := v.wordAt(s.Caret)
		if a == b {
			return "", false
		}
		v.st.sels.List[v.st.sels.Main] = Sel{Anchor: a, Caret: b, WantCol: -1}
		return string(v.doc.Text(a, b)), true
	}
	return string(v.doc.Text(s.Start(), s.End())), false
}

// SelectAllOccurrences selects all the occurrences of the selected text (of the word at the caret)
func (v *View) SelectAllOccurrences(matchCase bool) {
	text, whole := v.occurrenceText()
	if text == "" {
		return
	}
	m, err := NewMatcher(text, SearchOptions{MatchCase: matchCase, WholeWord: whole, WordChars: v.opts.WordChars})
	if err != nil {
		return
	}
	main := v.st.sels.MainSel()
	var sels []Sel
	mainIdx := 0
	m.FindAll(v.doc, 0, v.doc.Len(), func(a, b int) bool {
		if a == main.Start() {
			mainIdx = len(sels)
		}
		sels = append(sels, Sel{Anchor: a, Caret: b, WantCol: -1})
		return len(sels) < 100000
	})
	if len(sels) == 0 {
		return
	}
	v.SetSelections(sels, mainIdx)
}

// AddNextOccurrence adds a selection of the next occurrence of the main selection
func (v *View) AddNextOccurrence(matchCase bool) {
	text, whole := v.occurrenceText()
	if text == "" {
		return
	}
	if whole {
		// The word itself is selected first
		v.selectionChanged()
		return
	}
	m, err := NewMatcher(text, SearchOptions{MatchCase: matchCase, WholeWord: v.isWholeWord(v.MainSelection().Start(), v.MainSelection().End()), WordChars: v.opts.WordChars})
	if err != nil {
		return
	}
	last := v.st.sels.List[len(v.st.sels.List)-1]
	for _, s := range v.st.sels.List {
		if s.End() > last.End() {
			last = s
		}
	}
	a, b, ok := m.FindNext(v.doc, last.End(), false, 0, v.doc.Len(), -1)
	if !ok {
		a, b, ok = m.FindNext(v.doc, 0, false, 0, v.doc.Len(), -1)
	}
	if !ok {
		return
	}
	for _, s := range v.st.sels.List {
		if s.Start() == a && s.End() == b {
			return // all are selected
		}
	}
	v.st.sels.List = append(v.st.sels.List, Sel{Anchor: a, Caret: b, WantCol: -1})
	v.st.sels.Main = len(v.st.sels.List) - 1
	v.st.sels.normalize()
	v.ensureVisible(b, 0)
	v.selectionChanged()
}

// UndoLastOccurrence drops the selection added last
func (v *View) UndoLastOccurrence() {
	if len(v.st.sels.List) < 2 {
		return
	}
	v.st.sels.List = append(v.st.sels.List[:v.st.sels.Main], v.st.sels.List[v.st.sels.Main+1:]...)
	v.st.sels.Main = len(v.st.sels.List) - 1
	v.ensureCaretVisible()
	v.selectionChanged()
}

// SkipOccurrence drops the selection added last and adds the next occurrence instead
func (v *View) SkipOccurrence(matchCase bool) {
	if len(v.st.sels.List) < 2 {
		v.AddNextOccurrence(matchCase)
		return
	}
	cur := v.st.sels.MainSel()
	v.AddNextOccurrence(matchCase)
	for i, s := range v.st.sels.List {
		if s.Start() == cur.Start() && s.End() == cur.End() {
			v.st.sels.List = append(v.st.sels.List[:i], v.st.sels.List[i+1:]...)
			break
		}
	}
	v.st.sels.Main = 0
	for i, s := range v.st.sels.List {
		if s.Start() > cur.Start() {
			v.st.sels.Main = i
			break
		}
	}
	v.selectionChanged()
}

func init() {
	commands["multi-select-all"] = func(v *View, _ string) { v.SelectAllOccurrences(false) }
	commands["multi-select-all-case"] = func(v *View, _ string) { v.SelectAllOccurrences(true) }
	commands["multi-select-next"] = func(v *View, _ string) { v.AddNextOccurrence(false) }
	commands["multi-select-undo"] = func(v *View, _ string) { v.UndoLastOccurrence() }
	commands["multi-select-skip"] = func(v *View, _ string) { v.SkipOccurrence(false) }
}

// nextChange returns the first line of the next (dir 1) or previous block
// of changed lines from the line, -1 if none
func (v *View) nextChange(line, dir int) int {
	var all LineMarkers
	all.Set(append(v.doc.Changed.Lines(), v.doc.Saved.Lines()...))
	if all.Len() == 0 {
		return -1
	}
	l := line
	// Leave the block the caret is in
	for all.Has(l) {
		l += dir
	}
	var target int
	if dir > 0 {
		target = all.Next(l - 1)
	} else {
		target = all.Prev(l + 1)
		for all.Has(target - 1) {
			target--
		}
	}
	return target
}

func init() {
	commands["next-change"] = func(v *View, _ string) {
		if l := v.nextChange(v.doc.LineOfOffset(v.Caret()), 1); l >= 0 {
			v.GotoLine(l)
		}
	}
	commands["prev-change"] = func(v *View, _ string) {
		if l := v.nextChange(v.doc.LineOfOffset(v.Caret()), -1); l >= 0 {
			v.GotoLine(l)
		}
	}
	commands["clear-change-history"] = func(v *View, _ string) {
		v.doc.Changed.Clear()
		v.doc.Saved.Clear()
		v.update()
	}
}
