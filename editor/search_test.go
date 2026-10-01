package editor

import (
	"regexp"
	"strings"
	"testing"
)

func findAll(t *testing.T, doc *Document, pattern string, opts SearchOptions) [][2]int {
	t.Helper()
	m, err := NewMatcher(pattern, opts)
	if err != nil {
		t.Fatal(err)
	}
	var out [][2]int
	m.FindAll(doc, 0, doc.Len(), func(a, b int) bool {
		out = append(out, [2]int{a, b})
		return true
	})
	return out
}

func TestSearchModes(t *testing.T) {
	doc := NewDocument()
	doc.SetText([]byte("Hello hello HELLO helloworld\nПривет привет ПРИВЕТ\nline2\tend"))
	cases := []struct {
		pattern string
		opts    SearchOptions
		want    int
	}{
		{"hello", SearchOptions{MatchCase: true}, 2},
		{"hello", SearchOptions{}, 4},
		{"hello", SearchOptions{WholeWord: true}, 3},
		{"привет", SearchOptions{}, 3},
		{"Привет", SearchOptions{MatchCase: true, WholeWord: true}, 1},
		{`h\w+`, SearchOptions{Mode: SearchRegex, MatchCase: true}, 2},
		{`^line\d`, SearchOptions{Mode: SearchRegex}, 1},
		{`\tend`, SearchOptions{Mode: SearchExtended}, 1},
		{`ЕТ\n`, SearchOptions{Mode: SearchExtended}, 1},
		{`D\r\nП`, SearchOptions{Mode: SearchRegex}, 1},
	}
	for _, c := range cases {
		if got := findAll(t, doc, c.pattern, c.opts); len(got) != c.want {
			t.Errorf("%q %+v: %d matches, want %d", c.pattern, c.opts, len(got), c.want)
		}
	}
}

func TestSearchAcrossWindows(t *testing.T) {
	// Matches around the window boundaries of a text larger than a window
	var sb strings.Builder
	for sb.Len() < 3*searchWindow {
		sb.WriteString("abcdefghij klmnopqrst ")
		if sb.Len()%1000 < 30 {
			sb.WriteString("\n")
		}
	}
	text := sb.String()
	doc := NewDocument()
	doc.SetText([]byte(text))
	for _, opts := range []SearchOptions{{MatchCase: true}, {}, {Mode: SearchRegex}} {
		pattern := "j klm"
		if opts.Mode == SearchRegex {
			pattern = `j\sklm`
		}
		got := findAll(t, doc, pattern, opts)
		want := strings.Count(text, "j klm")
		if len(got) != want {
			t.Errorf("%+v: %d matches, want %d", opts, len(got), want)
		}
		for i := 1; i < len(got); i++ {
			if got[i][0] <= got[i-1][0] {
				t.Fatalf("%+v: matches out of order", opts)
			}
		}
	}
}

func TestFindNextBackward(t *testing.T) {
	doc := NewDocument()
	doc.SetText([]byte("one two one two one"))
	m, _ := NewMatcher("one", SearchOptions{MatchCase: true})
	a, _, ok := m.FindNext(doc, 11, true, 0, doc.Len(), -1)
	if !ok || a != 8 {
		t.Errorf("backward from 11: %d %v", a, ok)
	}
	// A match ending after the position is not before it
	if a, _, ok = m.FindNext(doc, 10, true, 0, doc.Len(), -1); !ok || a != 0 {
		t.Errorf("backward from 10: %d %v", a, ok)
	}
	a, _, ok = m.FindNext(doc, 9, false, 0, doc.Len(), -1)
	if !ok || a != 16 {
		t.Errorf("forward from 9: %d %v", a, ok)
	}
	if _, _, ok = m.FindNext(doc, 17, false, 0, doc.Len(), -1); ok {
		t.Error("found past the last")
	}
}

func TestReplacement(t *testing.T) {
	doc := NewDocument()
	doc.SetText([]byte("key1=value1\nkey2=value2"))
	m, _ := NewMatcher(`(\w+)=(\w+)`, SearchOptions{Mode: SearchRegex, MatchCase: true})
	r := NewReplacement(`\U$2\E: ${1}\n`, SearchRegex)
	var edits []Edit
	m.FindAll(doc, 0, doc.Len(), func(a, b int) bool {
		edits = append(edits, Edit{Pos: a, Len: b - a, Text: r.Expand(m, doc, a, b)})
		return true
	})
	doc.ApplyEdits(edits)
	if got := doc.String(); got != "VALUE1: key1\n\nVALUE2: key2\n" {
		t.Errorf("got %q", got)
	}
	doc.Undo()
	if got := doc.String(); got != "key1=value1\nkey2=value2" {
		t.Errorf("undo: %q", got)
	}
	if got := ExpandExtended(`a\tb\x41Ж\\`); got != "a\tbAЖ\\" {
		t.Errorf("extended: %q", got)
	}
	if _, err := regexp.Compile(convertRegex(`\<word\>\r\n`)); err != nil {
		t.Error(err)
	}
}
