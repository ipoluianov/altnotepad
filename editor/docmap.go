package editor

import (
	"image"
	"image/color"
	"unicode/utf8"

	"github.com/ipoluianov/nui/ui"
)

// DocMap is the document map: the text of a view drawn small, a line a
// couple of pixels high, with the part shown in the view framed. A click or
// a drag moves the view there.
type DocMap struct {
	ui.Widget

	view     *View
	top      int // the first line drawn
	dragging bool
	styles   LineStyles
}

const (
	mapLineH  = 2
	mapColW   = 1
	mapMaxCol = 400
)

func NewDocMap() *DocMap {
	m := &DocMap{}
	m.InitWidget()
	m.SetTypeName("DocMap")
	m.SetXExpandable(true)
	m.SetYExpandable(true)
	m.SetMouseCursor(ui.MouseCursorPointer)
	m.SetOnPaint(m.paint)
	m.SetOnMouseDown(func(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
		if button == ui.MouseButtonLeft {
			m.dragging = true
			m.scrollTo(y)
		}
		return true
	})
	m.SetOnMouseMove(func(x, y int, mods ui.KeyModifiers) bool {
		if m.dragging {
			m.scrollTo(y)
		}
		return true
	})
	m.SetOnMouseUp(func(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
		m.dragging = false
		return true
	})
	m.SetOnMouseWheel(func(dx, dy int) bool {
		if m.view != nil {
			m.view.ScrollLines(-dy * 3)
		}
		return true
	})
	return m
}

// SetView shows the text of the view
func (m *DocMap) SetView(v *View) {
	m.view = v
	if m.Form() != nil {
		m.Form().Update()
	}
}

// scrollTo scrolls the view so the line at y of the map is in its middle
func (m *DocMap) scrollTo(y int) {
	v := m.view
	if v == nil {
		return
	}
	line := m.lineAt(y)
	v.ScrollToLine(max(0, line-v.visibleRows()/2))
}

// lineAt returns the line drawn at y
func (m *DocMap) lineAt(y int) int {
	st := m.view.st
	k := st.visibleIndexOf(m.top) + max(0, y)/mapLineH
	return st.lineOfVisibleIndex(k)
}

func (m *DocMap) paint(cnv *ui.Canvas) {
	v := m.view
	if v == nil {
		return
	}
	img := cnv.RGBA()
	ox, oy := cnv.TranslatedX(), cnv.TranslatedY()
	clip := image.Rect(cnv.ClipX(), cnv.ClipY(), cnv.ClipX()+cnv.ClipW(), cnv.ClipY()+cnv.ClipH())
	sc := v.scheme
	w, h := m.Width(), m.Height()
	fillRect(img, clip, ox, oy, w, h, sc.Background)
	st := v.st
	doc := v.doc
	st.rebuildHidden(v.tabSize())
	rows := h / mapLineH
	// The lines shown in the view
	first := st.topLine
	last := first
	if n := len(v.rows); n > 0 {
		last = v.rows[n-1].line
	}
	// Keep the frame of the view within the map
	fi, li := st.visibleIndexOf(first), st.visibleIndexOf(last)
	ti := st.visibleIndexOf(m.top)
	total := st.visibleCount()
	switch {
	case total <= rows:
		ti = 0
	case fi < ti:
		ti = fi
	case li >= ti+rows:
		ti = li - rows + 1
	}
	ti = max(0, min(ti, total-rows))
	m.top = st.lineOfVisibleIndex(ti)
	tab := v.tabSize()
	line := m.top
	for r := 0; r < rows && line >= 0; r++ {
		y := oy + r*mapLineH
		ls := doc.LineStart(line)
		le := doc.LineEnd(line)
		text := doc.buf.View(ls, min(le, ls+mapMaxCol*2))
		if le-ls <= maxStyledLine && !doc.large {
			doc.LineStyles(line, doc.buf.View(ls, le), &m.styles)
		} else {
			m.styles.reset()
		}
		runs := m.styles.Runs
		ri := 0
		style := StyleDefault
		col := 0
		for i := 0; i < len(text) && col < mapMaxCol; {
			c, size := rune(text[i]), 1
			if c >= utf8.RuneSelf {
				c, size = utf8.DecodeRune(text[i:])
			}
			for ri < len(runs) && int(runs[ri].Start) <= i {
				style = runs[ri].Style
				ri++
			}
			cells := runeCells(c, size, col, tab)
			if c != ' ' && c != '\t' {
				fg := sc.Styles[style].Fore
				fillRect(img, clip, ox+2+col*mapColW, y, max(1, cells)*mapColW, mapLineH-1, color.RGBA{fg.R, fg.G, fg.B, 190})
			}
			col += cells
			i += size
		}
		line = st.nextVisible(line)
	}
	// The frame of the view
	fy := oy + (fi-ti)*mapLineH
	fh := max(4, (li-fi+1)*mapLineH)
	fillRect(img, clip, ox, fy, w, fh, withA(sc.Selection, 70))
	border := withA(sc.Foreground, 110)
	fillRect(img, clip, ox, fy, w, 1, border)
	fillRect(img, clip, ox, fy+fh-1, w, 1, border)
}
