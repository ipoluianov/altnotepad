package editor

import (
	"image"
	"image/color"
	"sort"
	"strconv"
	"unicode/utf8"

	"github.com/ipoluianov/nui/ui"
)

// Lines longer than this are shown without syntax highlighting
const maxStyledLine = longLine

// smartHighlight holds the occurrences of the selected word in the shown text
type smartHighlight struct {
	valid   bool
	version int64
	word    string
	from    int
	to      int
	ranges  []Range
}

// braceMatch holds the brace at the caret and its match
type braceMatch struct {
	valid   bool
	version int64
	a, b    int // positions; b < 0 - no match
}

type paintCtx struct {
	img      *image.RGBA
	clip     image.Rectangle
	textClip image.Rectangle
	ox, oy   int
	sels     []Sel // sorted by Start
	mainLine int
	prevInd  int
}

func (v *View) paint(cnv *ui.Canvas) {
	st := v.st
	v.layout()
	st.rebuildHidden(v.tabSize())
	v.clampScroll()

	sc := v.scheme
	pc := &paintCtx{img: cnv.RGBA(), ox: cnv.TranslatedX(), oy: cnv.TranslatedY()}
	pc.clip = image.Rect(cnv.ClipX(), cnv.ClipY(), cnv.ClipX()+cnv.ClipW(), cnv.ClipY()+cnv.ClipH())
	W, H := v.Width(), v.Height()
	fillRect(pc.img, pc.clip, pc.ox, pc.oy, W, H, sc.Background)
	if v.gutterW > 0 {
		fillRect(pc.img, pc.clip, pc.ox, pc.oy, v.gutterW, H, sc.GutterBack)
	}
	pc.textClip = image.Rect(pc.ox+v.gutterW, pc.oy, pc.ox+W-v.scrollW, pc.oy+v.textH).Intersect(pc.clip)

	pc.sels = append(pc.sels[:0], st.sels.List...)
	sort.Slice(pc.sels, func(i, j int) bool { return pc.sels[i].Start() < pc.sels[j].Start() })
	pc.mainLine = v.doc.LineOfOffset(st.sels.MainSel().Caret)

	lh := v.font.LineHeight
	v.rows = v.rows[:0]
	line, sub := st.topLine, st.topSub
	for y := 0; y < v.textH; y += lh {
		v.rows = append(v.rows, rowInfo{line, sub, y})
		l, s, ok := v.nextRow(line, sub)
		if !ok {
			break
		}
		line, sub = l, s
	}
	if len(v.rows) > 0 {
		v.updateSmartHighlight(v.rows[0].line, v.rows[len(v.rows)-1].line)
	}
	v.updateBraceMatch()

	pc.prevInd = 0
	for _, r := range v.rows {
		v.paintRow(pc, r)
	}

	// The long line marker
	if v.opts.EdgeColumn > 0 && !v.wrapping() {
		x := pc.ox + v.textLeft - st.scrollX + v.opts.EdgeColumn*v.font.CharWidth
		fillRect(pc.img, pc.textClip, x, pc.oy, 1, v.textH, sc.Edge)
	}
	v.paintScrollBars(pc)
	v.paintCompletion(pc)
}

func (v *View) paintRow(pc *paintCtx, r rowInfo) {
	sc := v.scheme
	doc := v.doc
	st := v.st
	cw := v.font.CharWidth
	lh := v.font.LineHeight
	tab := v.tabSize()
	y := pc.oy + r.y
	line, sub := r.line, r.sub

	lineStart := doc.LineStart(line)
	lineEnd := doc.LineEnd(line)
	lineLen := lineEnd - lineStart
	long := lineLen > longLine
	rs, re := v.rowRange(line, sub, lineLen)
	lastRow := re == lineLen
	// text is the part of the line read, from ss to se of the line: all of
	// it, or of a long line what is on the screen
	ss, se := 0, lineLen
	var rowStartCol, segCol int
	switch {
	case !long:
		rowStartCol = -1
	case v.wrapping():
		ss, se = rs, re
		rowStartCol = v.colAt(line, rs)
		segCol = rowStartCol
	default:
		v.noteLineWidth(v.layoutOf(line).cols * cw)
		firstCol := max(0, st.scrollX/cw-1)
		ss, _ = v.offAt(line, firstCol, false)
		se, _ = v.offAt(line, firstCol+v.textW/cw+4, false)
		segCol = v.colAt(line, ss)
	}
	text := doc.buf.View(lineStart+ss, lineStart+se)
	tbase := ss // the offset in the line of text[0]
	if rowStartCol < 0 {
		rowStartCol = colOfOffset(text, rs, tab)
		segCol = rowStartCol
		ss, se = rs, re
	} else if long && !v.wrapping() {
		rowStartCol = 0
	}
	xBase := pc.ox + v.textLeft - st.scrollX
	if sub > 0 {
		xBase += v.wrap(line).indent * cw
	}
	xOfCol := func(col int) int { return xBase + (col-rowStartCol)*cw }
	baseline := y + (lh-v.font.Ascent-v.font.Descent)/2 + v.font.Ascent

	mainSel := st.sels.MainSel()
	// The current line
	if v.opts.CurrentLine && line == pc.mainLine && mainSel.Empty() && len(st.sels.List) == 1 {
		fillRect(pc.img, pc.textClip, pc.textClip.Min.X, y, pc.textClip.Dx(), lh, sc.CurrentLine)
	}

	// The syntax
	if !long {
		doc.LineStyles(line, text, &v.styles)
	} else {
		v.styles.reset()
	}
	runs := v.styles.Runs

	// Style backgrounds that fill the row, e.g. the added lines of a diff
	if len(runs) > 0 {
		last := runs[len(runs)-1].Style
		if b := sc.Styles[last].Back; b.A > 0 && lastRow {
			endX := xOfCol(colOfOffset(text, lineLen, tab))
			fillRect(pc.img, pc.textClip, endX, y, pc.textClip.Max.X-endX, lh, b)
		}
	}

	// Cells of the row: the columns of the characters
	type cell struct {
		off, size int
		r         rune
		col, w    int
	}
	// Offsets in the cells are from the line start
	cells := make([]cell, 0, se-ss)
	col := segCol
	for i := ss; i < se; {
		ch, size := rune(text[i-tbase]), 1
		if ch >= utf8.RuneSelf {
			ch, size = utf8.DecodeRune(text[i-tbase:])
		}
		w := runeCells(ch, size, col, tab)
		if xOfCol(col) > pc.textClip.Max.X && ch != '\t' {
			break
		}
		cells = append(cells, cell{i, size, ch, col, w})
		col += w
		i += size
	}
	endCol := col
	if (len(cells) > 0 && cells[len(cells)-1].off+cells[len(cells)-1].size < re) || se < re {
		endCol = -1 // cut at the right edge
	} else if !long {
		v.noteLineWidth(xOfCol(endCol) - xBase + (rowStartCol * cw))
	}
	colAtOff := func(off int) int {
		// Binary search in the cells
		k := sort.Search(len(cells), func(k int) bool { return cells[k].off >= off })
		if k < len(cells) {
			return cells[k].col
		}
		if endCol >= 0 {
			return endCol
		}
		return v.colAt(line, off)
	}

	// Style backgrounds
	for k, run := range runs {
		b := sc.Styles[run.Style].Back
		if b.A == 0 {
			continue
		}
		start := max(int(run.Start), rs)
		end := re
		if k+1 < len(runs) {
			end = min(int(runs[k+1].Start), re)
		}
		if start < end {
			x0, x1 := xOfCol(colAtOff(start)), xOfCol(colAtOff(end))
			fillRect(pc.img, pc.textClip, x0, y, x1-x0, lh, b)
		}
	}

	absRS, absRE := lineStart+rs, lineStart+re
	fillRange := func(a, b int, col color.RGBA, extraCells int) {
		a, b = max(a, absRS), min(b, absRE)
		if a > b || (a == b && extraCells == 0) {
			return
		}
		x0 := xOfCol(colAtOff(a - lineStart))
		x1 := xOfCol(colAtOff(b-lineStart)) + extraCells*cw
		fillRect(pc.img, pc.textClip, x0, y, x1-x0, lh, col)
	}

	// Marks of the search
	for m := 0; m < NumMarkStyles; m++ {
		for _, rg := range doc.Marks[m].InRange(absRS, absRE+1) {
			fillRange(rg.Start, rg.End, sc.Marks[m], 0)
		}
	}
	// Smart highlighting
	if v.smart.valid {
		k := sort.Search(len(v.smart.ranges), func(k int) bool { return v.smart.ranges[k].End > absRS })
		for ; k < len(v.smart.ranges) && v.smart.ranges[k].Start < absRE; k++ {
			rg := v.smart.ranges[k]
			fillRange(rg.Start, rg.End, sc.SmartHighlight, 0)
		}
	}
	// Matching braces
	if v.brace.valid && v.brace.a >= 0 {
		col := sc.BraceMatch
		if v.brace.b < 0 {
			col = sc.BraceBad
		}
		for _, p := range []int{v.brace.a, v.brace.b} {
			if p >= absRS && p < absRE {
				fillRange(p, p+1, col, 0)
			}
		}
	}

	// Selections
	selCol := sc.Selection
	if !v.focused {
		selCol = ui.MixColors(sc.Selection, sc.Background, 0.35)
	}
	k := sort.Search(len(pc.sels), func(k int) bool { return pc.sels[k].End() >= absRS })
	for ; k < len(pc.sels); k++ {
		s := pc.sels[k]
		if s.Start() > absRE {
			break
		}
		if s.IsEmptyText() && !(st.sels.Rect && s.AnchorVS != s.CaretVS) {
			continue
		}
		a, b := s.Start(), s.End()
		extra := 0
		if lastRow && b > absRE {
			extra = 1 // the line break is selected
		}
		if st.sels.Rect && lastRow && b == absRE {
			vsA, vsB := s.AnchorVS, s.CaretVS
			if s.Anchor > s.Caret {
				vsA, vsB = vsB, vsA
			}
			if a == absRE {
				// All of it in the virtual space
				x0 := xOfCol(colAtOff(re) + min(vsA, vsB))
				x1 := xOfCol(colAtOff(re) + max(vsA, vsB))
				fillRect(pc.img, pc.textClip, x0, y, x1-x0, lh, selCol)
				continue
			}
			extra = vsB
		}
		fillRange(a, b, selCol, extra)
	}

	// Indent guides
	if v.opts.ShowIndentGuides && sub == 0 && !long {
		ind := 0
		blank := true
		for _, c := range cells {
			if c.r == ' ' || c.r == '\t' {
				ind = c.col + c.w
				continue
			}
			blank = false
			break
		}
		if blank && len(cells) == len(text) {
			ind = pc.prevInd
		} else {
			if blank {
				ind = pc.prevInd
			}
			pc.prevInd = ind
		}
		for c := tab; c < ind; c += tab {
			x := xOfCol(c)
			for yy := y; yy < y+lh; yy += 2 {
				fillRect(pc.img, pc.textClip, x, yy, 1, 1, sc.IndentGuide)
			}
		}
	}

	// The text
	ri := 0
	style := StyleDefault
	for _, c := range cells {
		for ri < len(runs) && int(runs[ri].Start) <= c.off {
			style = runs[ri].Style
			ri++
		}
		x := xOfCol(c.col)
		if x+c.w*cw < pc.textClip.Min.X {
			continue
		}
		sty := sc.Styles[style]
		fore := sty.Fore
		if sc.SelectionText.A > 0 && v.inSelection(pc, lineStart+c.off) {
			fore = sc.SelectionText
		}
		switch {
		case c.r == ' ':
			if v.opts.ShowSpaces {
				d := max(2, cw/5)
				fillRect(pc.img, pc.textClip, x+(cw-d)/2, y+(lh-d)/2, d, d, sc.Whitespace)
			}
		case c.r == '\t':
			if v.opts.ShowSpaces {
				v.drawTabArrow(pc, x, y, c.w*cw)
			}
		case c.w > 0 && controlName(c.r, c.size, text[c.off-tbase]) != "":
			name := controlName(c.r, c.size, text[c.off-tbase])
			fillRect(pc.img, pc.textClip, x+1, y+1, c.w*cw-2, lh-2, sc.ControlCharBack)
			for i, ch := range name {
				v.font.DrawRune(pc.img, pc.textClip, x+i*cw, baseline, ch, sc.Background, false)
			}
		case c.w == 0:
			// A combining mark over the previous character
			v.font.DrawRune(pc.img, pc.textClip, x-cw, baseline, c.r, fore, sty.Bold)
		default:
			gx := x
			if c.w == 2 {
				if gw := v.font.GlyphWidth(c.r); gw < 2*cw {
					gx += (2*cw - gw) / 2
				}
			}
			v.font.DrawRune(pc.img, pc.textClip, gx, baseline, c.r, fore, sty.Bold)
		}
	}

	endX := 0
	if endCol >= 0 {
		endX = xOfCol(endCol)
	} else if lastRow {
		endX = xOfCol(v.colAt(line, re))
	}
	// The line break
	if lastRow && v.opts.ShowEOL && line < doc.LineCount()-1 {
		v.drawEOL(pc, endX+2, y)
	}
	// Wrapped row
	if !lastRow && v.opts.ShowWrapSymbol {
		v.drawWrapSymbol(pc, pc.textClip.Max.X-cw-3, y)
	}
	// Folded block
	if lastRow && st.folded.Has(line) {
		if _, hidden := st.hiddenRangeOf(line + 1); hidden {
			fillRect(pc.img, pc.textClip, pc.textClip.Min.X, y+lh-1, pc.textClip.Dx(), 1, sc.FoldedLine)
			bx := endX + cw
			if v.opts.ShowEOL {
				bx += 5 * cw
			}
			v.drawBox(pc, bx, y+2, 3*cw, lh-4, sc.FoldedLine)
			for i := 0; i < 3; i++ {
				fillRect(pc.img, pc.textClip, bx+i*cw+cw/2-1, y+lh/2, 2, 2, sc.FoldedLine)
			}
		}
	}

	// Carets
	if v.focused && v.caretOn {
		for _, s := range pc.sels {
			if s.Caret < absRS || s.Caret > absRE || (s.Caret == absRE && !lastRow) {
				continue
			}
			x := xOfCol(colAtOff(s.Caret-lineStart) + s.CaretVS)
			if v.overtype {
				fillRect(pc.img, pc.textClip, x, y+lh-2, cw, 2, sc.Caret)
			} else {
				fillRect(pc.img, pc.textClip, x, y, max(1, v.opts.CaretWidth), lh, sc.Caret)
			}
		}
	}

	if sub == 0 {
		v.paintGutter(pc, line, y)
	}
}

func (v *View) inSelection(pc *paintCtx, pos int) bool {
	k := sort.Search(len(pc.sels), func(k int) bool { return pc.sels[k].End() > pos })
	return k < len(pc.sels) && pc.sels[k].Start() <= pos && !pc.sels[k].IsEmptyText()
}

// noteLineWidth widens the horizontal scroll range for a long line
func (v *View) noteLineWidth(w int) {
	if w > v.maxLineW {
		v.maxLineW = w
	}
}

func (v *View) drawTabArrow(pc *paintCtx, x, y, w int) {
	lh := v.font.LineHeight
	col := v.scheme.Whitespace
	cy := y + lh/2
	x0, x1 := x+2, x+w-3
	if x1-x0 < 4 {
		return
	}
	fillRect(pc.img, pc.textClip, x0, cy, x1-x0, 1, col)
	for i := 1; i <= 3; i++ {
		fillRect(pc.img, pc.textClip, x1-i, cy-i, 1, 1, col)
		fillRect(pc.img, pc.textClip, x1-i, cy+i, 1, 1, col)
	}
}

// drawWrapSymbol draws a hooked arrow: the row goes on below
func (v *View) drawWrapSymbol(pc *paintCtx, x, y int) {
	cw, lh := v.font.CharWidth, v.font.LineHeight
	col := v.scheme.Whitespace
	right := x + cw
	top := y + lh/4
	mid := y + lh*3/5
	fillRect(pc.img, pc.textClip, right-1, top, 1, mid-top, col)
	fillRect(pc.img, pc.textClip, x+1, mid, right-x-1, 1, col)
	for i := 1; i <= 3; i++ {
		fillRect(pc.img, pc.textClip, x+1+i, mid-i, 1, 1, col)
		fillRect(pc.img, pc.textClip, x+1+i, mid+i, 1, 1, col)
	}
}

func (v *View) drawEOL(pc *paintCtx, x, y int) {
	name := "LF"
	switch v.doc.EOL {
	case EOLCRLF:
		name = "CRLF"
	case EOLCR:
		name = "CR"
	}
	cw := v.font.CharWidth
	lh := v.font.LineHeight
	w := len(name)*cw + 2
	fillRect(pc.img, pc.textClip, x, y+1, w, lh-2, v.scheme.Whitespace)
	baseline := y + (lh-v.font.Ascent-v.font.Descent)/2 + v.font.Ascent
	for i, ch := range name {
		v.font.DrawRune(pc.img, pc.textClip, x+1+i*cw, baseline, ch, v.scheme.Background, false)
	}
}

func (v *View) drawBox(pc *paintCtx, x, y, w, h int, col color.RGBA) {
	fillRect(pc.img, pc.textClip, x, y, w, 1, col)
	fillRect(pc.img, pc.textClip, x, y+h-1, w, 1, col)
	fillRect(pc.img, pc.textClip, x, y, 1, h, col)
	fillRect(pc.img, pc.textClip, x+w-1, y, 1, h, col)
}

func (v *View) paintGutter(pc *paintCtx, line, y int) {
	sc := v.scheme
	lh := v.font.LineHeight
	cw := v.font.CharWidth
	gclip := image.Rect(pc.ox, pc.oy, pc.ox+v.gutterW, pc.oy+v.textH).Intersect(pc.clip)
	baseline := y + (lh-v.font.Ascent-v.font.Descent)/2 + v.font.Ascent
	x := pc.ox
	if v.lineNumW > 0 {
		num := strconv.Itoa(line + 1)
		col := sc.LineNumber
		if line == pc.mainLine {
			col = sc.LineNumberCur
		}
		nx := x + v.lineNumW - cw/2 - len(num)*cw
		for i, ch := range num {
			v.font.DrawRune(pc.img, gclip, nx+i*cw, baseline, ch, col, false)
		}
		x += v.lineNumW
	}
	if v.opts.ChangeHistory {
		switch {
		case v.doc.Changed.Has(line):
			fillRect(pc.img, gclip, x-3, y, 3, lh, sc.ChangedLine)
		case v.doc.Saved.Has(line):
			fillRect(pc.img, gclip, x-3, y, 3, lh, sc.SavedLine)
		}
	}
	if v.markerW > 0 {
		if v.doc.Bookmarks.Has(line) {
			d := min(v.markerW-4, lh-4)
			bx, by := x+(v.markerW-d)/2, y+(lh-d)/2
			drawDisc(pc.img, gclip, bx, by, d, sc.Bookmark)
		}
		x += v.markerW
	}
	if v.foldW > 0 {
		v.paintFoldMargin(pc, gclip, line, x, y)
	}
}

func drawDisc(img *image.RGBA, clip image.Rectangle, x, y, d int, col color.RGBA) {
	r := float64(d) / 2
	for yy := 0; yy < d; yy++ {
		for xx := 0; xx < d; xx++ {
			dx, dy := float64(xx)+0.5-r, float64(yy)+0.5-r
			dist := dx*dx + dy*dy
			if dist <= (r-0.5)*(r-0.5) {
				fillRect(img, clip, x+xx, y+yy, 1, 1, col)
			} else if dist <= (r+0.5)*(r+0.5) {
				fillRect(img, clip, x+xx, y+yy, 1, 1, withA(col, 110))
			}
		}
	}
}

func (v *View) paintFoldMargin(pc *paintCtx, clip image.Rectangle, line, x, y int) {
	sc := v.scheme
	st := v.st
	lh := v.font.LineHeight
	tab := v.tabSize()
	cx := x + v.foldW/2
	col := sc.FoldMarker
	header := st.IsFoldHeader(line, tab)
	inside, closing := false, false
	if st.indentFold() {
		ind := st.lineIndent(line, tab)
		inside = ind > 0 || (ind < 0 && pc.prevInd > 0)
	} else {
		lvl := v.doc.hl.Level(line)
		f := v.doc.hl.fold(line)
		inside = lvl > 0
		closing = lvl+f.min() < lvl && !header
		if header && lvl+f.min() < lvl {
			inside = true
		}
	}
	box := min(v.foldW-2, lh-4) | 1
	by := y + (lh-box)/2
	if inside {
		if closing {
			fillRect(pc.img, clip, cx, y, 1, lh/2, col)
			fillRect(pc.img, clip, cx, y+lh/2, v.foldW/2-1, 1, col)
		} else {
			fillRect(pc.img, clip, cx, y, 1, lh, col)
		}
	}
	if header {
		folded := st.folded.Has(line)
		if !folded {
			fillRect(pc.img, clip, cx, by+box, 1, y+lh-by-box, col)
		}
		bx := cx - box/2
		fillRect(pc.img, clip, bx, by, box, box, sc.GutterBack)
		fillRect(pc.img, clip, bx, by, box, 1, col)
		fillRect(pc.img, clip, bx, by+box-1, box, 1, col)
		fillRect(pc.img, clip, bx, by, 1, box, col)
		fillRect(pc.img, clip, bx+box-1, by, 1, box, col)
		fillRect(pc.img, clip, bx+2, by+box/2, box-4, 1, col)
		if folded {
			fillRect(pc.img, clip, cx, by+2, 1, box-4, col)
		}
	}
}

// ---- Scroll bars ----

// vScrollGeometry returns the thumb of the vertical scroll bar: its y and
// height, and the scroll range in lines
func (v *View) vScrollGeometry() (thumbY, thumbH, pos, maxPos int) {
	st := v.st
	total := st.visibleCount()
	rows := v.visibleRows()
	pos = st.topLine - st.hiddenBefore(st.topLine)
	maxPos = total - rows
	if v.opts.ScrollPast || v.wrapping() {
		maxPos = total - 1
	}
	if v.wrapping() && total <= rows && !v.opts.ScrollPast {
		// All the rows may fit
		n := 0
		for l := 0; l >= 0 && n <= rows; l = st.nextVisible(l) {
			n += v.rowCount(l)
		}
		if n <= rows {
			maxPos = 0
		}
	}
	maxPos = max(maxPos, 0)
	track := v.textH
	if maxPos == 0 {
		return 0, track, pos, 0
	}
	thumbH = max(20, track*rows/max(rows, total+rows))
	thumbY = (track - thumbH) * min(pos, maxPos) / maxPos
	return thumbY, thumbH, pos, maxPos
}

// hScrollGeometry returns the thumb of the horizontal scroll bar and the content width
func (v *View) hScrollGeometry() (thumbX, thumbW, content int) {
	content = max(v.maxLineW+v.font.CharWidth*4, v.st.scrollX+v.textW)
	track := v.textW + v.textLeft - v.gutterW
	if content <= v.textW {
		return 0, track, content
	}
	thumbW = max(20, track*v.textW/content)
	thumbX = (track - thumbW) * v.st.scrollX / max(1, content-v.textW)
	return thumbX, thumbW, content
}

func (v *View) paintScrollBars(pc *paintCtx) {
	sc := v.scheme
	W := v.Width()
	bg := ui.MixColors(sc.Background, sc.Foreground, 0.05)
	thumb := ui.MixColors(sc.Background, sc.Foreground, 0.3)
	thumbHover := ui.MixColors(sc.Background, sc.Foreground, 0.45)
	// Vertical
	x := pc.ox + W - v.scrollW
	fillRect(pc.img, pc.clip, x, pc.oy, v.scrollW, v.Height(), bg)
	ty, th, _, maxPos := v.vScrollGeometry()
	if maxPos > 0 {
		c := thumb
		if v.mouse.overBar == 1 || v.mouse.dragBar == 1 {
			c = thumbHover
		}
		fillRect(pc.img, pc.clip, x+2, pc.oy+ty, v.scrollW-4, th, c)
	}
	// Bookmarks and the caret on the vertical bar
	total := max(1, v.doc.LineCount())
	for _, l := range v.doc.Bookmarks.lines {
		my := pc.oy + l*v.textH/total
		fillRect(pc.img, pc.clip, x, my, v.scrollW, 2, sc.Bookmark)
	}
	cy := pc.oy + pc.mainLine*v.textH/total
	fillRect(pc.img, pc.clip, x, cy, v.scrollW, 2, sc.Caret)
	// Horizontal
	if v.showHBar {
		y := pc.oy + v.Height() - v.scrollW
		fillRect(pc.img, pc.clip, pc.ox+v.gutterW, y, W-v.gutterW, v.scrollW, bg)
		tx, tw, content := v.hScrollGeometry()
		if content > v.textW {
			c := thumb
			if v.mouse.overBar == 2 || v.mouse.dragBar == 2 {
				c = thumbHover
			}
			fillRect(pc.img, pc.clip, pc.ox+v.gutterW+tx, y+2, tw, v.scrollW-4, c)
		}
	}
}

// updateSmartHighlight finds the occurrences of the selected word in the lines shown
func (v *View) updateSmartHighlight(firstLine, lastLine int) {
	if !v.opts.SmartHighlight || len(v.st.sels.List) != 1 {
		v.smart = smartHighlight{}
		return
	}
	s := v.st.sels.MainSel()
	if s.IsEmptyText() || s.End()-s.Start() > 200 {
		v.smart = smartHighlight{}
		return
	}
	from, to := v.doc.LineStart(firstLine), v.doc.LineEnd(lastLine)
	if to-from > 2<<20 {
		// Long lines on the screen: the text around the selection only
		from, to = max(from, s.Start()-64<<10), min(to, s.End()+64<<10)
	}
	word := string(v.doc.Text(s.Start(), s.End()))
	if v.smart.valid && v.smart.version == v.doc.version && v.smart.word == word && v.smart.from == from && v.smart.to == to {
		return
	}
	v.smart = smartHighlight{valid: true, version: v.doc.version, word: word, from: from, to: to}
	for _, c := range word {
		if c == '\n' {
			v.smart.valid = false
			return
		}
	}
	if v.opts.SmartWholeWord && !v.isWholeWord(s.Start(), s.End()) {
		v.smart.valid = false
		return
	}
	opts := SearchOptions{MatchCase: v.opts.SmartMatchCase, WholeWord: v.opts.SmartWholeWord, WordChars: v.opts.WordChars}
	m, err := NewMatcher(word, opts)
	if err != nil {
		v.smart.valid = false
		return
	}
	m.FindAll(v.doc, from, to, func(a, b int) bool {
		if a != s.Start() {
			v.smart.ranges = append(v.smart.ranges, Range{a, b})
		}
		return len(v.smart.ranges) < 5000
	})
}

func (v *View) updateBraceMatch() {
	if !v.opts.BraceMatch || v.doc.large {
		v.brace = braceMatch{a: -1}
		return
	}
	if v.brace.valid && v.brace.version == v.doc.version {
		return
	}
	v.brace = braceMatch{valid: true, version: v.doc.version, a: -1, b: -1}
	s := v.st.sels.MainSel()
	if !s.IsEmptyText() || len(v.st.sels.List) != 1 {
		return
	}
	pos := s.Caret
	if v.isLong(v.doc.LineOfOffset(pos)) {
		return
	}
	for _, p := range []int{pos - 1, pos} {
		if p < 0 || p >= v.doc.Len() {
			continue
		}
		c := v.doc.buf.ByteAt(p)
		if isBrace(c) {
			v.brace.a = p
			v.brace.b = v.FindMatchingBrace(p)
			return
		}
	}
}
