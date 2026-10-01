package editor

import (
	"strings"
	"testing"
)

func TestWrapLine(t *testing.T) {
	text := []byte("aaa bbb ccc ddd eee")
	breaks, _ := wrapLine(text, 8, 4, false)
	var rows []string
	prev := 0
	for _, b := range append(breaks, len(text)) {
		rows = append(rows, string(text[prev:b]))
		prev = b
	}
	if got := strings.Join(rows, "|"); got != "aaa bbb |ccc ddd |eee" {
		t.Errorf("rows %q", got)
	}
	// A word longer than the row is cut
	breaks, _ = wrapLine([]byte(strings.Repeat("x", 20)), 8, 4, false)
	if len(breaks) != 2 || breaks[0] != 8 || breaks[1] != 16 {
		t.Errorf("breaks %v", breaks)
	}
	// The indent of the wrapped rows
	_, ind := wrapLine([]byte("    "+strings.Repeat("w ", 30)), 20, 4, true)
	if ind != 4 {
		t.Errorf("indent %d", ind)
	}
}

func TestColumns(t *testing.T) {
	text := []byte("\tab中c")
	if c := colOfOffset(text, len(text), 4); c != 4+2+2+1 {
		t.Errorf("width %d", c)
	}
	if off, _ := offsetOfCol(text, 6, 4, false); off != 3 {
		t.Errorf("offset of col 6: %d", off)
	}
	if off, beyond := offsetOfCol(text, 12, 4, false); off != len(text) || beyond != 3 {
		t.Errorf("past the end: %d %d", off, beyond)
	}
}
