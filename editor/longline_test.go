package editor

import (
	"math/rand"
	"strings"
	"testing"
)

func TestLongLineColumns(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	pieces := []string{"a", "bc", "\t", "Ж", "中", " ", "xyz", "é"}
	var sb strings.Builder
	for sb.Len() < 3*longLine {
		sb.WriteString(pieces[r.Intn(len(pieces))])
	}
	line := sb.String()
	v := NewView()
	doc := NewDocument()
	doc.SetText([]byte("short\n" + line + "\nend"))
	v.SetState(NewViewState(doc))
	text := []byte(line)
	if !v.isLong(1) {
		t.Fatal("not long")
	}
	if got, want := v.lineColsOf(1), lineCols(text, 4); got != want {
		t.Fatalf("cols %d, want %d", got, want)
	}
	for k := 0; k < 300; k++ {
		off := r.Intn(len(text) + 1)
		// A character boundary
		for off < len(text) && (text[off]&0xC0) == 0x80 {
			off++
		}
		if got, want := v.colAt(1, off), colOfOffset(text, off, 4); got != want {
			t.Fatalf("colAt(%d) = %d, want %d", off, got, want)
		}
		col := r.Intn(lineCols(text, 4) + 10)
		o1, b1 := v.offAt(1, col, false)
		o2, b2 := offsetOfCol(text, col, 4, false)
		if o1 != o2 || b1 != b2 {
			t.Fatalf("offAt(%d) = %d %d, want %d %d", col, o1, b1, o2, b2)
		}
	}
}
