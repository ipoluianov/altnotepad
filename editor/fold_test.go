package editor

import "testing"

func TestIndentFoldHeaders(t *testing.T) {
	doc := NewDocument()
	doc.SetText([]byte("import os\n\nclass Foo:\n    def f(self):\n        pass\n\ndef main():\n    for i in x:\n        print(i)\n"))
	doc.SetLanguage(LanguageByID("python"))
	st := NewViewState(doc)
	want := map[int]bool{2: true, 3: true, 6: true, 7: true}
	for l := 0; l < doc.LineCount(); l++ {
		if got := st.IsFoldHeader(l, 4); got != want[l] {
			t.Errorf("line %d: header %v, want %v", l+1, got, want[l])
		}
	}
	if e := st.foldBlockEnd(6, 4); e != 8 {
		t.Errorf("block of main ends at %d, want 8", e)
	}
}

func TestBraceFold(t *testing.T) {
	doc := NewDocument()
	doc.SetText([]byte("int f() {\n  if (x) {\n    a;\n  } else {\n    b;\n  }\n}\n"))
	doc.SetLanguage(LanguageByID("c"))
	st := NewViewState(doc)
	cases := map[int]int{0: 6, 1: 2, 3: 5}
	for h, end := range cases {
		if !st.IsFoldHeader(h, 4) {
			t.Errorf("line %d is not a header", h+1)
		}
		if e := st.foldBlockEnd(h, 4); e != end {
			t.Errorf("block of line %d ends at %d, want %d", h+1, e+1, end+1)
		}
	}
	if st.IsFoldHeader(2, 4) {
		t.Error("line 3 is a header")
	}
	if h := st.headerOf(4, 4); h != 3 {
		t.Errorf("header of line 5 is %d, want 4", h+1)
	}
}
