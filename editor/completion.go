package editor

import (
	"image"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/ipoluianov/nui/ui"
)

// completion is the list of the words to complete the one at the caret,
// drawn over the text below the caret
type completion struct {
	active bool
	words  []string // the candidates
	shown  []string // those starting with the typed prefix
	sel    int
	top    int
	start  int // where the word being completed starts
	// the box of the list in the view, for the mouse
	x, y, w, h int
}

const completionRows = 10

// CompletionWords returns the extra words offered by the auto-completion,
// e.g. the keywords of the language; set by the application
var CompletionWords func(v *View) []string

// wordStartAt returns the start of the word ending at pos
func (v *View) wordStartAt(pos int) int {
	p := pos
	ls := v.doc.LineStart(v.doc.LineOfOffset(pos))
	for p > ls {
		r, size := v.doc.buf.DecodeLastRune(p)
		if charClass(r, v.opts.WordChars) != ccWord {
			break
		}
		p -= size
	}
	return p
}

// documentWords returns the words of the text around pos starting with prefix
func (v *View) documentWords(prefix string, pos int, skipStart int) []string {
	const around = 2 << 20
	from, to := max(0, pos-around), min(v.doc.Len(), pos+around)
	text := v.doc.buf.View(from, to)
	seen := map[string]bool{}
	var out []string
	lp := strings.ToLower(prefix)
	i := 0
	for i < len(text) {
		r, size := utf8.DecodeRune(text[i:])
		if charClass(r, v.opts.WordChars) != ccWord {
			i += size
			continue
		}
		j := i
		for j < len(text) {
			r2, s2 := utf8.DecodeRune(text[j:])
			if charClass(r2, v.opts.WordChars) != ccWord {
				break
			}
			j += s2
		}
		if from+i != skipStart && j-i > len(prefix) && j-i < 100 {
			w := string(text[i:j])
			if !seen[w] && strings.HasPrefix(strings.ToLower(w), lp) {
				seen[w] = true
				out = append(out, w)
			}
		}
		i = j
	}
	return out
}

// Complete opens the list of the words completing the one at the caret;
// with one candidate it is inserted at once. auto: opened while typing,
// then it is not opened for short prefixes and an only candidate is not inserted.
func (v *View) Complete(auto bool) {
	if v.readOnly() || len(v.st.sels.List) != 1 {
		return
	}
	pos := v.Caret()
	start := v.wordStartAt(pos)
	prefix := string(v.doc.Text(start, pos))
	if auto && utf8.RuneCountInString(prefix) < 3 {
		v.comp.active = false
		return
	}
	words := v.documentWords(prefix, pos, start)
	if CompletionWords != nil {
		lp := strings.ToLower(prefix)
		for _, w := range CompletionWords(v) {
			if len(w) > len(prefix) && strings.HasPrefix(strings.ToLower(w), lp) {
				words = append(words, w)
			}
		}
	}
	sort.Slice(words, func(i, j int) bool { return strings.ToLower(words[i]) < strings.ToLower(words[j]) })
	words = uniqueStrings(words)
	if len(words) == 0 {
		v.comp.active = false
		v.update()
		return
	}
	if len(words) == 1 && !auto {
		v.insertCompletion(start, words[0])
		return
	}
	v.comp = completion{active: true, words: words, start: start}
	v.filterCompletion()
	v.update()
}

func uniqueStrings(s []string) []string {
	out := s[:0]
	for i, w := range s {
		if i == 0 || w != s[i-1] {
			out = append(out, w)
		}
	}
	return out
}

// filterCompletion shows the words starting with what is typed now
func (v *View) filterCompletion() {
	c := &v.comp
	pos := v.Caret()
	if pos < c.start || v.doc.LineOfOffset(pos) != v.doc.LineOfOffset(c.start) {
		c.active = false
		return
	}
	prefix := strings.ToLower(string(v.doc.Text(c.start, pos)))
	c.shown = c.shown[:0]
	for _, w := range c.words {
		if strings.HasPrefix(strings.ToLower(w), prefix) && len(w) > len(prefix) {
			c.shown = append(c.shown, w)
		}
	}
	if len(c.shown) == 0 {
		c.active = false
		return
	}
	c.sel = 0
	c.top = 0
}

func (v *View) insertCompletion(start int, word string) {
	v.comp.active = false
	v.record("complete", word)
	v.edit(func() {
		pos := v.Caret()
		v.doc.Replace(start, pos-start, []byte(word))
		v.st.sels.setSingle(caretSel(start + len(word)))
	})
}

// completionKey handles the keys of the open list; false for the other keys
func (v *View) completionKey(key ui.Key, mods ui.KeyModifiers) bool {
	c := &v.comp
	if !c.active {
		return false
	}
	if mods.Ctrl || mods.Alt {
		c.active = false
		return false
	}
	move := func(n int) {
		c.sel = max(0, min(c.sel+n, len(c.shown)-1))
		if c.sel < c.top {
			c.top = c.sel
		} else if c.sel >= c.top+completionRows {
			c.top = c.sel - completionRows + 1
		}
		v.update()
	}
	switch key {
	case ui.KeyArrowUp:
		move(-1)
	case ui.KeyArrowDown:
		move(1)
	case ui.KeyPageUp:
		move(-completionRows)
	case ui.KeyPageDown:
		move(completionRows)
	case ui.KeyEnter, ui.KeyTab:
		v.insertCompletion(c.start, c.shown[c.sel])
	case ui.KeyEsc:
		c.active = false
		v.update()
	case ui.KeyBackspace:
		return false // handled, then the list is filtered again
	case ui.KeyArrowLeft, ui.KeyArrowRight, ui.KeyHome, ui.KeyEnd, ui.KeyDelete:
		c.active = false
		return false
	default:
		return false
	}
	return true
}

// paintCompletion draws the list below (or above) the caret
func (v *View) paintCompletion(pc *paintCtx) {
	c := &v.comp
	if !c.active || len(c.shown) == 0 {
		return
	}
	cw, lh := v.font.CharWidth, v.font.LineHeight
	// The caret row
	cl, cs := v.posRow(c.start)
	rowY := -1
	for _, r := range v.rows {
		if r.line == cl && r.sub == cs {
			rowY = r.y
		}
	}
	if rowY < 0 {
		c.active = false
		return
	}
	maxLen := 0
	for _, w := range c.shown {
		maxLen = max(maxLen, utf8.RuneCountInString(w))
	}
	rows := min(completionRows, len(c.shown))
	w := min(maxLen+2, 60)*cw + 8
	h := rows*lh + 4
	x := v.textLeft - v.st.scrollX + v.xOfPos(c.start, 0)
	y := rowY + lh
	if y+h > v.textH && rowY-h >= 0 {
		y = rowY - h
	}
	x = max(v.gutterW, min(x, v.Width()-v.scrollW-w))
	c.x, c.y, c.w, c.h = x, y, w, h
	sc := v.scheme
	bg := ui.MixColors(sc.Background, sc.Foreground, 0.06)
	border := ui.MixColors(sc.Background, sc.Foreground, 0.4)
	clip := image.Rect(pc.ox, pc.oy, pc.ox+v.Width(), pc.oy+v.Height()).Intersect(pc.clip)
	fillRect(pc.img, clip, pc.ox+x, pc.oy+y, w, h, border)
	fillRect(pc.img, clip, pc.ox+x+1, pc.oy+y+1, w-2, h-2, bg)
	prefixLen := utf8.RuneCount(v.doc.Text(c.start, v.Caret()))
	for i := 0; i < rows; i++ {
		k := c.top + i
		if k >= len(c.shown) {
			break
		}
		iy := pc.oy + y + 2 + i*lh
		fore := sc.Foreground
		if k == c.sel {
			fillRect(pc.img, clip, pc.ox+x+1, iy, w-2, lh, sc.Selection)
		}
		baseline := iy + (lh-v.font.Ascent-v.font.Descent)/2 + v.font.Ascent
		n := 0
		for j, r := range []rune(c.shown[k]) {
			if j >= 60 {
				break
			}
			col := fore
			bold := j < prefixLen
			v.font.DrawRune(pc.img, clip, pc.ox+x+4+j*cw, baseline, r, col, bold)
			n++
		}
	}
	// A scroll mark when the list is longer
	if len(c.shown) > rows {
		th := max(6, (h-4)*rows/len(c.shown))
		ty := (h - 4 - th) * c.top / max(1, len(c.shown)-rows)
		fillRect(pc.img, clip, pc.ox+x+w-4, pc.oy+y+2+ty, 2, th, border)
	}
}

// completionMouse handles a click on the list; false if it is not on it
func (v *View) completionMouse(x, y int) bool {
	c := &v.comp
	if !c.active {
		return false
	}
	if x < c.x || x >= c.x+c.w || y < c.y || y >= c.y+c.h {
		c.active = false
		return false
	}
	k := c.top + (y-c.y-2)/v.font.LineHeight
	if k >= 0 && k < len(c.shown) {
		v.insertCompletion(c.start, c.shown[k])
	}
	return true
}

// completionWheel scrolls the list under the mouse
func (v *View) completionWheel(dy int) bool {
	c := &v.comp
	if !c.active || v.mouse.x < c.x || v.mouse.x >= c.x+c.w || v.mouse.y < c.y || v.mouse.y >= c.y+c.h {
		return false
	}
	c.top = max(0, min(c.top-dy*3, len(c.shown)-completionRows))
	v.update()
	return true
}
