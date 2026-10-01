package editor

import (
	"bytes"
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// SearchMode is how the search text is understood
type SearchMode int

const (
	SearchNormal   SearchMode = iota
	SearchExtended            // \n, \t, \0, \xHH... are the characters
	SearchRegex
)

// SearchOptions are the options of Find and Replace
type SearchOptions struct {
	Mode      SearchMode
	MatchCase bool
	WholeWord bool
	// DotAll: in a regular expression . also matches the line breaks
	DotAll    bool
	WordChars string
}

// Matcher finds a text or a regular expression in documents
type Matcher struct {
	opts    SearchOptions
	literal []byte // a case-sensitive text found with bytes.Index
	// folded: a text found ignoring the case, in lower case; the text
	// searched is lowered the same way, keeping the offsets
	folded []byte
	// foldASCII: the folded text is ASCII, the text needs only its ASCII letters lowered
	foldASCII bool
	re        *regexp.Regexp
	maxLen    int // the longest match of a literal search, for the overlap of the windows
}

// ErrEmptySearch: nothing to search for
var ErrEmptySearch = errors.New("the text to find is empty")

// NewMatcher prepares the search of the pattern
func NewMatcher(pattern string, opts SearchOptions) (*Matcher, error) {
	if pattern == "" {
		return nil, ErrEmptySearch
	}
	m := &Matcher{opts: opts}
	switch opts.Mode {
	case SearchRegex:
		expr := convertRegex(pattern)
		flags := "(?m"
		if !opts.MatchCase {
			flags += "i"
		}
		if opts.DotAll {
			flags += "s"
		}
		re, err := regexp.Compile(flags + ")" + expr)
		if err != nil {
			return nil, err
		}
		m.re = re
	default:
		text := pattern
		if opts.Mode == SearchExtended {
			text = ExpandExtended(pattern)
		}
		text = strings.ReplaceAll(text, "\r\n", "\n")
		text = strings.ReplaceAll(text, "\r", "\n")
		if text == "" {
			return nil, ErrEmptySearch
		}
		m.maxLen = len(text) * 3
		if opts.MatchCase {
			m.literal = []byte(text)
		} else if f, ok := foldPattern(text); ok {
			m.folded = f
			m.foldASCII = isASCII(f)
		} else {
			re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(text))
			if err != nil {
				return nil, err
			}
			m.re = re
		}
	}
	return m, nil
}

// IsRegex reports whether the matcher is a regular expression search
func (m *Matcher) IsRegex() bool { return m.opts.Mode == SearchRegex }

// convertRegex turns the Notepad++ (Boost) syntax into Go's: \r\n and \R
// are the line break, \< and \> the word boundaries
func convertRegex(p string) string {
	var b strings.Builder
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c != '\\' || i+1 >= len(p) {
			b.WriteByte(c)
			continue
		}
		n := p[i+1]
		switch n {
		case 'r':
			if strings.HasPrefix(p[i+2:], `\n`) {
				i += 3
			} else {
				i++
			}
			b.WriteString(`\n`)
		case 'R':
			b.WriteString(`\n`)
			i++
		case '<', '>':
			b.WriteString(`\b`)
			i++
		case 'h':
			b.WriteString(`[ \t]`)
			i++
		default:
			b.WriteByte(c)
			b.WriteByte(n)
			i++
		}
	}
	return b.String()
}

// ExpandExtended turns the escapes of the Extended search mode into characters:
// \n, \r, \t, \0, \\, \xHH, \uHHHH, \oOOO, \dDDD, \bBBBBBBBB
func ExpandExtended(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			b.WriteByte(c)
			continue
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case '0':
			b.WriteByte(0)
		case '\\':
			b.WriteByte('\\')
		case 'x', 'u', 'o', 'd', 'b':
			base, digits := 16, 2
			switch s[i] {
			case 'u':
				digits = 4
			case 'o':
				base, digits = 8, 3
			case 'd':
				base, digits = 10, 3
			case 'b':
				base, digits = 2, 8
			}
			if i+digits < len(s)+0 && i+1+digits <= len(s) {
				if n, err := strconv.ParseUint(s[i+1:i+1+digits], base, 32); err == nil {
					b.WriteRune(rune(n))
					i += digits
					continue
				}
			}
			b.WriteByte('\\')
			b.WriteByte(s[i])
		default:
			b.WriteByte('\\')
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// The search reads the text in windows, so a huge document is not copied whole
const (
	searchWindow  = 4 << 20
	searchOverlap = 64 << 10
)

// wholeWordAt reports whether the match from a to b is a whole word
func (m *Matcher) wholeWordAt(doc *Document, a, b int) bool {
	if !m.opts.WholeWord {
		return true
	}
	return isWordBoundary(doc, a, m.opts.WordChars) && isWordBoundary(doc, b, m.opts.WordChars)
}

// isWordBoundary: the characters on the two sides of pos are not both word characters
func isWordBoundary(doc *Document, pos int, wordChars string) bool {
	before, _ := doc.buf.DecodeLastRune(pos)
	after, _ := doc.buf.DecodeRune(pos)
	if pos <= 0 || pos >= doc.Len() {
		return true
	}
	return charClass(before, wordChars) != ccWord || charClass(after, wordChars) != ccWord
}

// findIn returns the matches in text (a window of the document starting at
// base), calling fn with document positions; fn returns false to stop
func (m *Matcher) findIn(doc *Document, text []byte, base int, fn func(a, b int) bool) bool {
	if m.literal != nil || m.folded != nil {
		lit := m.literal
		if m.folded != nil {
			lit = m.folded
			if m.foldASCII {
				text = foldASCII(text)
			} else {
				text = foldText(text)
			}
		}
		i := 0
		for {
			k := bytes.Index(text[i:], lit)
			if k < 0 {
				return true
			}
			a := base + i + k
			b := a + len(lit)
			if m.wholeWordAt(doc, a, b) {
				if !fn(a, b) {
					return false
				}
				i += k + len(lit)
			} else {
				_, size := utf8.DecodeRune(text[i+k:])
				i += k + max(1, size)
			}
			if i > len(text) {
				return true
			}
		}
	}
	for _, loc := range m.re.FindAllIndex(text, -1) {
		a, b := base+loc[0], base+loc[1]
		if !m.wholeWordAt(doc, a, b) {
			continue
		}
		if !fn(a, b) {
			return false
		}
	}
	return true
}

// windowEnd returns where a window from start should end: at a line break
// after the window size, so the regular expressions see whole lines
func windowEnd(doc *Document, start, to int) int {
	end := start + searchWindow
	if end >= to {
		return to
	}
	line := doc.LineOfOffset(end)
	le := doc.LineEnd(line)
	if le-start > 4*searchWindow || le >= to {
		return min(end, to)
	}
	return le + 1
}

// FindAll calls fn with the matches from `from` to `to` in order; fn returns false to stop
func (m *Matcher) FindAll(doc *Document, from, to int, fn func(a, b int) bool) {
	from = max(0, from)
	to = min(to, doc.Len())
	if m.re != nil && to-from > 4*searchWindow && runtime.NumCPU() > 1 {
		m.findAllParallel(doc, from, to, fn)
		return
	}
	pos := from
	lastEnd := -1
	for pos < to {
		end := windowEnd(doc, pos, to)
		// A literal can span the window end: read a bit more, report only
		// the matches that start within the window
		readEnd := end
		if m.re == nil && end < to {
			readEnd = min(to, end+m.maxLen)
		}
		text := doc.buf.View(pos, readEnd)
		cont := m.findIn(doc, text, pos, func(a, b int) bool {
			if a >= end && end < to {
				return true
			}
			if a < lastEnd || (a == b && a == lastEnd) {
				return true
			}
			lastEnd = b
			return fn(a, b)
		})
		if !cont {
			return
		}
		pos = end
		if m.re == nil && lastEnd > pos {
			pos = lastEnd
		}
	}
}

// FindNext returns the first match at or after from (before it when
// backward), within start..end of the document. An empty match at skipEmpty
// is skipped, so repeated searches move on.
func (m *Matcher) FindNext(doc *Document, from int, backward bool, start, end int, skipEmpty int) (int, int, bool) {
	start = max(0, start)
	end = min(end, doc.Len())
	if !backward {
		fa, fb, found := -1, -1, false
		m.FindAll(doc, max(from, start), end, func(a, b int) bool {
			if a == b && a == skipEmpty {
				return true
			}
			fa, fb, found = a, b, true
			return false
		})
		return fa, fb, found
	}
	// Backward: windows going back from `from`; the last match ending
	// before it wins
	winEnd := min(from, end)
	for winEnd > start {
		winStart := max(start, winEnd-searchWindow)
		if winStart > start {
			winStart = doc.LineStart(doc.LineOfOffset(winStart))
		}
		fa, fb, found := -1, -1, false
		readEnd := winEnd
		if m.re == nil {
			readEnd = min(end, winEnd+m.maxLen)
		}
		text := doc.buf.View(winStart, readEnd)
		m.findIn(doc, text, winStart, func(a, b int) bool {
			if b > from || a >= winEnd {
				return true
			}
			if a == b && a == skipEmpty {
				return true
			}
			fa, fb, found = a, b, true
			return true
		})
		if found {
			return fa, fb, true
		}
		if winStart == start {
			break
		}
		prevEnd := winEnd
		winEnd = winStart + min(searchOverlap, winEnd-winStart)
		if winEnd >= prevEnd {
			winEnd = winStart
		}
	}
	return -1, -1, false
}

// Replacement makes the replacement text of a match
type Replacement struct {
	parts []replPart
	regex bool
}

type replPart struct {
	text   string
	group  int    // >= 0: a group of the match
	name   string // a named group
	caseOp byte   // 'U', 'L', 'E', 'u', 'l': case conversion from here on
}

// NewReplacement prepares the replacement text: in the Extended mode its
// escapes are the characters, in the regular expression mode $1, \1, ${name},
// $& and \U, \L, \E, \u, \l are the groups and case conversions
func NewReplacement(s string, mode SearchMode) *Replacement {
	r := &Replacement{regex: mode == SearchRegex}
	switch mode {
	case SearchNormal:
		r.parts = []replPart{{text: normalizeEOL(s), group: -1}}
		return r
	case SearchExtended:
		r.parts = []replPart{{text: normalizeEOL(ExpandExtended(s)), group: -1}}
		return r
	}
	var lit strings.Builder
	flush := func() {
		if lit.Len() > 0 {
			r.parts = append(r.parts, replPart{text: lit.String(), group: -1})
			lit.Reset()
		}
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '$' && i+1 < len(s) {
			n := s[i+1]
			switch {
			case n == '&':
				flush()
				r.parts = append(r.parts, replPart{group: 0})
				i++
				continue
			case n >= '0' && n <= '9':
				j := i + 1
				for j < len(s) && j < i+3 && s[j] >= '0' && s[j] <= '9' {
					j++
				}
				g, _ := strconv.Atoi(s[i+1 : j])
				flush()
				r.parts = append(r.parts, replPart{group: g})
				i = j - 1
				continue
			case n == '{':
				k := strings.IndexByte(s[i:], '}')
				if k > 0 {
					name := s[i+2 : i+k]
					flush()
					if g, err := strconv.Atoi(name); err == nil {
						r.parts = append(r.parts, replPart{group: g})
					} else {
						r.parts = append(r.parts, replPart{group: -2, name: name})
					}
					i += k
					continue
				}
			case n == '$':
				lit.WriteByte('$')
				i++
				continue
			}
		}
		if c == '\\' && i+1 < len(s) {
			n := s[i+1]
			switch {
			case n >= '0' && n <= '9':
				flush()
				r.parts = append(r.parts, replPart{group: int(n - '0')})
			case n == 'n' || n == 'r':
				lit.WriteByte('\n')
				if n == 'r' && strings.HasPrefix(s[i+2:], `\n`) {
					i += 2
				}
			case n == 't':
				lit.WriteByte('\t')
			case n == '0':
				lit.WriteByte(0)
			case n == 'U' || n == 'L' || n == 'E' || n == 'u' || n == 'l':
				flush()
				r.parts = append(r.parts, replPart{group: -1, caseOp: n})
			default:
				lit.WriteByte(n)
			}
			i++
			continue
		}
		lit.WriteByte(c)
	}
	flush()
	return r
}

func normalizeEOL(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

// Expand returns the replacement of the match from a to b of the document
func (r *Replacement) Expand(m *Matcher, doc *Document, a, b int) []byte {
	if !r.regex || m.re == nil {
		if len(r.parts) == 0 {
			return nil
		}
		return []byte(r.parts[0].text)
	}
	// Match again within the found text to get the groups; the context
	// around it keeps the anchors and the lookaround right
	ctxStart := doc.LineStart(doc.LineOfOffset(a))
	ctxEnd := doc.LineEnd(doc.LineOfOffset(b))
	text := doc.buf.View(ctxStart, max(ctxEnd, b))
	var groups []int
	for _, loc := range m.re.FindAllSubmatchIndex(text, -1) {
		if ctxStart+loc[0] == a && ctxStart+loc[1] == b {
			groups = loc
			break
		}
	}
	if groups == nil {
		text = doc.buf.View(a, b)
		ctxStart = a
		groups = m.re.FindSubmatchIndex(text)
		if groups == nil {
			return nil
		}
	}
	group := func(g int) []byte {
		if g < 0 || 2*g+1 >= len(groups) || groups[2*g] < 0 {
			return nil
		}
		return text[groups[2*g]:groups[2*g+1]]
	}
	var out []byte
	mode := byte('E')
	once := byte(0)
	add := func(s []byte) {
		for _, c := range string(s) {
			switch {
			case once == 'u':
				c = unicode.ToUpper(c)
				once = 0
			case once == 'l':
				c = unicode.ToLower(c)
				once = 0
			case mode == 'U':
				c = unicode.ToUpper(c)
			case mode == 'L':
				c = unicode.ToLower(c)
			}
			out = utf8.AppendRune(out, c)
		}
	}
	for _, p := range r.parts {
		switch {
		case p.caseOp == 'u' || p.caseOp == 'l':
			once = p.caseOp
		case p.caseOp != 0:
			mode = p.caseOp
		case p.group >= 0:
			add(group(p.group))
		case p.group == -2:
			add(group(m.re.SubexpIndex(p.name)))
		default:
			add([]byte(p.text))
		}
	}
	return out
}

// isBrace reports whether the character is a bracket matched by Go to Brace
func isBrace(c byte) bool {
	return c == '(' || c == ')' || c == '[' || c == ']' || c == '{' || c == '}'
}

// FindMatchingBrace returns the position of the bracket matching the one
// at pos, -1 if none; brackets in comments and strings are skipped when
// the one at pos is not in one
func (v *View) FindMatchingBrace(pos int) int {
	doc := v.doc
	c := doc.buf.ByteAt(pos)
	var match byte
	forward := true
	switch c {
	case '(':
		match = ')'
	case '[':
		match = ']'
	case '{':
		match = '}'
	case ')':
		match, forward = '(', false
	case ']':
		match, forward = '[', false
	case '}':
		match, forward = '{', false
	default:
		return -1
	}
	line := doc.LineOfOffset(pos)
	styleOf := func(l int) func(off int) Style {
		text := doc.LineText(l)
		var ls LineStyles
		if len(text) <= maxStyledLine {
			doc.LineStyles(l, text, &ls)
		}
		return ls.StyleAt
	}
	lineStart := doc.LineStart(line)
	style := styleOf(line)
	ownStyle := style(pos - lineStart)
	depth := 0
	const maxLines = 20000
	for l, n := line, 0; l >= 0 && l < doc.LineCount() && n < maxLines; n++ {
		ls := doc.LineStart(l)
		text := doc.LineText(l)
		if l != line {
			style = styleOf(l)
		}
		i := len(text) - 1
		if forward {
			i = 0
		}
		if l == line {
			i = pos - ls
		}
		for i >= 0 && i < len(text) {
			ch := text[i]
			if (ch == c || ch == match) && style(i) == ownStyle {
				if ch == c {
					depth++
				} else {
					depth--
					if depth == 0 {
						return ls + i
					}
				}
			}
			if forward {
				i++
			} else {
				i--
			}
		}
		if forward {
			l++
		} else {
			l--
		}
	}
	return -1
}

// findAllParallel searches the windows of a large text with a regular
// expression on all the processors; the matches are reported in order
func (m *Matcher) findAllParallel(doc *Document, from, to int, fn func(a, b int) bool) {
	workers := runtime.NumCPU()
	type window struct {
		start int
		text  []byte
		locs  [][]int
	}
	pos := from
	lastEnd := -1
	for pos < to {
		// A batch of windows, copied here: the document is not safe for concurrent use
		batch := make([]*window, 0, workers*2)
		for len(batch) < cap(batch) && pos < to {
			end := windowEnd(doc, pos, to)
			batch = append(batch, &window{start: pos, text: doc.buf.Bytes(pos, end)})
			pos = end
		}
		var wg sync.WaitGroup
		for _, w := range batch {
			wg.Add(1)
			go func(w *window) {
				defer wg.Done()
				w.locs = m.re.FindAllIndex(w.text, -1)
			}(w)
		}
		wg.Wait()
		for _, w := range batch {
			for _, loc := range w.locs {
				a, b := w.start+loc[0], w.start+loc[1]
				if a < lastEnd || (a == b && a == lastEnd) || !m.wholeWordAt(doc, a, b) {
					continue
				}
				lastEnd = b
				if !fn(a, b) {
					return
				}
			}
		}
	}
}

// simpleFold lowers a character when its lower case has the same length in
// UTF-8, so the offsets of a lowered text are those of the text
func simpleFold(r rune, size int) (rune, bool) {
	l := unicode.ToLower(r)
	if l == r {
		return r, true
	}
	return l, utf8.RuneLen(l) == size
}

// foldPattern lowers the text to find; false if a character changes its length
func foldPattern(s string) ([]byte, bool) {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		size := utf8.RuneLen(r)
		l, ok := simpleFold(r, size)
		if !ok || unicode.ToUpper(l) != unicode.ToUpper(r) {
			return nil, false
		}
		// The upper case must have the same length too, or the text is not lowered to it
		if u := unicode.ToUpper(r); utf8.RuneLen(u) != size {
			return nil, false
		}
		out = utf8.AppendRune(out, l)
	}
	return out, true
}

// foldText returns the text lowered like foldPattern, keeping its length
func foldText(text []byte) []byte {
	out := make([]byte, len(text))
	for i := 0; i < len(text); {
		c := text[i]
		if c < utf8.RuneSelf {
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
			}
			out[i] = c
			i++
			continue
		}
		r, size := utf8.DecodeRune(text[i:])
		if l, ok := simpleFold(r, size); ok && r != utf8.RuneError {
			utf8.EncodeRune(out[i:], l)
		} else {
			copy(out[i:i+size], text[i:i+size])
		}
		i += size
	}
	return out
}

func isASCII(b []byte) bool {
	for _, c := range b {
		if c >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

// foldASCII returns the text with its ASCII letters lowered
func foldASCII(text []byte) []byte {
	out := make([]byte, len(text))
	for i, c := range text {
		if c-'A' < 26 {
			c += 'a' - 'A'
		}
		out[i] = c
	}
	return out
}
