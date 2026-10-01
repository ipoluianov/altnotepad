package editor

import (
	"testing"
)

func newTestView(text, lang string) *View {
	v := NewView()
	doc := NewDocument()
	doc.SetText([]byte(text))
	doc.SetLanguage(LanguageByID(lang))
	v.SetState(NewViewState(doc))
	o := v.Options()
	o.InsertTabs = false
	o.TabSize = 4
	v.SetOptions(o)
	return v
}

func check(t *testing.T, v *View, want string) {
	t.Helper()
	if got := v.Document().String(); got != want {
		t.Errorf("text:\n%q\nwant\n%q", got, want)
	}
}

func TestLineCommands(t *testing.T) {
	v := newTestView("b\na\nc\na\n", "text")
	v.Exec("sort-asc", "")
	check(t, v, "a\na\nb\nc\n")
	v.Exec("remove-dup-lines", "")
	check(t, v, "a\nb\nc\n")
	v.Exec("sort-desc", "")
	check(t, v, "c\nb\na\n")
	v.Exec("undo", "")
	check(t, v, "a\nb\nc\n")

	v = newTestView("10\n9\n100\nx\n", "text")
	v.Exec("sort-int-asc", "")
	check(t, v, "x\n9\n10\n100\n")

	v = newTestView("one\ntwo\nthree", "text")
	v.SetCaret(5) // "two"
	v.Exec("duplicate", "")
	check(t, v, "one\ntwo\ntwo\nthree")
	v.Exec("move-lines-up", "")
	check(t, v, "two\none\ntwo\nthree")
	v.Exec("move-lines-down", "")
	v.Exec("move-lines-down", "")
	v.Exec("move-lines-down", "")
	check(t, v, "one\ntwo\nthree\ntwo")
	v.Exec("delete-lines", "")
	check(t, v, "one\ntwo\nthree")

	v = newTestView("a  \n  b\t\n\n c ", "text")
	v.Exec("trim-trailing", "")
	check(t, v, "a\n  b\n\n c")
	v.Exec("remove-empty-lines", "")
	check(t, v, "a\n  b\n c")
	v.Exec("trim-both", "")
	check(t, v, "a\nb\nc")
	v.SelectAll()
	v.Exec("join-lines", "")
	check(t, v, "a b c")
}

func TestCaseAndComments(t *testing.T) {
	v := newTestView("hello World", "text")
	v.SelectAll()
	v.Exec("case-upper", "")
	check(t, v, "HELLO WORLD")
	v.Exec("case-proper", "")
	check(t, v, "Hello World")
	v.Exec("case-invert", "")
	check(t, v, "hELLO wORLD")

	v = newTestView("int a;\n  int b;\n", "c")
	v.SetSelection(0, 12)
	v.Exec("toggle-comment", "")
	check(t, v, "// int a;\n//   int b;\n")
	v.Exec("toggle-comment", "")
	check(t, v, "int a;\n  int b;\n")

	v = newTestView("x = 1\n", "python")
	v.Exec("toggle-comment", "")
	check(t, v, "# x = 1\n")

	v = newTestView("<p>hi</p>", "html")
	v.Exec("toggle-comment", "")
	check(t, v, "<!-- <p>hi</p> -->")
}

func TestIndentAndTabs(t *testing.T) {
	v := newTestView("a\nb\n", "text")
	v.SetSelection(0, 3)
	v.Exec("tab", "")
	check(t, v, "    a\n    b\n")
	v.Exec("backtab", "")
	check(t, v, "a\nb\n")

	v = newTestView("\tx\ty", "text")
	v.Exec("tabs-to-spaces", "")
	check(t, v, "    x   y")
	v.Exec("spaces-to-tabs-leading", "")
	check(t, v, "\tx   y")
}

func TestTypingAndUndo(t *testing.T) {
	v := newTestView("", "c")
	v.TypeText("i")
	v.TypeText("f")
	v.TypeText(" ")
	v.TypeText("(")
	v.TypeText("x")
	check(t, v, "if (x")
	v.Exec("newline", "")
	v.TypeText("{")
	v.Exec("newline", "")
	check(t, v, "if (x\n{\n    ")
	v.Undo()
	v.Undo()
	v.Undo()
	check(t, v, "if (x")

	// Auto-close and type over
	v = newTestView("", "c")
	o := v.Options()
	o.AutoClose = true
	v.SetOptions(o)
	v.TypeText("(")
	check(t, v, "()")
	v.TypeText("a")
	v.TypeText(")")
	check(t, v, "(a)")
}

func TestMultiCaret(t *testing.T) {
	v := newTestView("aa\nbb\ncc", "text")
	v.SetSelections([]Sel{caretSel(1), caretSel(4), caretSel(7)}, 0)
	v.TypeText("X")
	v.TypeText("Y")
	check(t, v, "aXYa\nbXYb\ncXYc")
	v.Backspace()
	check(t, v, "aXa\nbXb\ncXc")
	v.Undo()
	check(t, v, "aXYa\nbXYb\ncXYc")

	// Rectangular selection with virtual space
	v = newTestView("abcd\nx\nefgh", "text")
	v.SetCaret(1)
	for i := 0; i < 2; i++ {
		v.Exec("line-down-rect", "")
	}
	v.Exec("char-right-rect", "")
	v.Exec("char-right-rect", "")
	v.TypeText("-")
	check(t, v, "a-d\nx-\ne-h")

	// Select all the occurrences, then type over them
	v = newTestView("foo bar foo baz foo", "text")
	v.SetCaret(1)
	v.Exec("multi-select-all", "")
	if n := len(v.Selections()); n != 3 {
		t.Fatalf("%d selections", n)
	}
	v.TypeText("q")
	check(t, v, "q bar q baz q")
}

func TestClipboardRect(t *testing.T) {
	v := newTestView("ab\ncd\nef", "text")
	v.SetCaret(0)
	v.Exec("char-right-rect", "")
	v.Exec("line-down-rect", "")
	text, _ := v.copyText()
	if text != "a\nc\n" {
		t.Fatalf("rect copy %q", text)
	}
	rectClipboard = text
	v.SetCaret(v.Document().Len())
	v.pasteRect(text)
	check(t, v, "ab\ncd\nefa\n  c")
}
