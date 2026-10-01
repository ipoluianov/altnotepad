package editor

import (
	"math/rand"
	"strings"
	"testing"
)

// fullStates lexes the whole text from scratch
func fullStates(lexer Lexer, text string) ([]uint32, []foldInfo) {
	lines := strings.Split(text, "\n")
	states := []uint32{0}
	folds := make([]foldInfo, len(lines))
	var ls LineStyles
	for i, l := range lines {
		ls.reset()
		next := lexer.LexLine([]byte(l), states[i], &ls)
		folds[i] = packFold(ls.minDepth, ls.depth)
		states = append(states, next)
	}
	return states, folds
}

func TestHighlighterIncremental(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	pieces := []string{"int x = 1;", "/* comment", "end */", "\n", "\n", "{", "}", "\"str\"", "// line", "a", " ", "\"open"}
	var sb strings.Builder
	for i := 0; i < 3000; i++ {
		sb.WriteString(pieces[r.Intn(len(pieces))])
	}
	doc := NewDocument()
	doc.SetText([]byte(sb.String()))
	doc.SetLanguage(LanguageByID("c"))
	for step := 0; step < 400; step++ {
		// Look at a random part, as the view does
		doc.hl.ensure(r.Intn(doc.LineCount() + 1))
		pos := r.Intn(doc.Len() + 1)
		if r.Intn(2) == 0 {
			doc.Insert(pos, []byte(pieces[r.Intn(len(pieces))]))
		} else {
			doc.Delete(pos, r.Intn(20))
		}
		if step%10 == 0 {
			doc.hl.ensure(doc.LineCount())
			want, wantFolds := fullStates(doc.hl.lexer, doc.String())
			for i := 0; i < doc.LineCount(); i++ {
				if doc.hl.states[i] != want[i] {
					t.Fatalf("step %d: state of line %d = %x, want %x", step, i, doc.hl.states[i], want[i])
				}
				if doc.hl.folds[i] != wantFolds[i] {
					t.Fatalf("step %d: fold of line %d differs", step, i)
				}
			}
			// Fold levels
			lvl := 0
			for i := 0; i < doc.LineCount(); i += 7 {
				lvl = 0
				for k := 0; k < i; k++ {
					lvl += wantFolds[k].end()
				}
				if got := doc.FoldLevel(i); got != lvl {
					t.Fatalf("step %d: level of line %d = %d, want %d", step, i, got, lvl)
				}
			}
		}
	}
}

func TestLexersDoNotPanic(t *testing.T) {
	samples := []string{
		"<html><head><style>a { color: red; }</style><script>var x = /re/g; if (a < b) {}</script></head>",
		"<?php echo \"hi\"; $x = 1; ?> <div class=\"a\">&amp;</div>",
		"# Title\n* item **bold** `code` [link](url)\n```\ncode\n```",
		"[section]\nkey=value\n; comment",
		"diff --git a b\n--- a\n+++ b\n@@ -1 +1 @@\n-x\n+y",
		"key: value\n- item: 1\n  other: \"str\" # c",
		"def f(x):\n    return r'x' + f\"{x}\"  # c\n\"\"\"doc\nstring\"\"\"",
		"SELECT * FROM t WHERE a = 'it''s' -- c",
		"#!/bin/bash\necho $HOME ${X} $(pwd) # c\nif [ -f x ]; then\nfi",
		"@echo off\nrem c\n:label\nset X=%1\necho %X%",
	}
	for _, l := range Languages() {
		for _, s := range samples {
			var ls LineStyles
			state := uint32(0)
			for _, line := range strings.Split(s, "\n") {
				ls.reset()
				state = l.lexer.LexLine([]byte(line), state, &ls)
				if state == unknownState {
					t.Fatalf("%s: returned the unknown state", l.ID)
				}
			}
		}
	}
}

func TestLanguageForFile(t *testing.T) {
	cases := map[string]string{
		"main.go": "go", "a.PY": "python", "Makefile": "makefile", "x.html": "html",
		"CMakeLists.txt": "cmake", "notes.txt": "text", "a.json": "json",
	}
	for name, want := range cases {
		if got := LanguageForFile(name, nil).ID; got != want {
			t.Errorf("%s: %s, want %s", name, got, want)
		}
	}
	if got := LanguageForFile("script", []byte("#!/usr/bin/env python3")).ID; got != "python" {
		t.Errorf("shebang: %s", got)
	}
}
