package editor

import (
	"bytes"
	"strings"
)

// subLexer styles a part of a line: a language embedded in another one
type subLexer interface {
	lexRange(text []byte, from, to int, state uint32, out *LineStyles) uint32
}

func firstNonSpace(text []byte) int {
	i := 0
	for i < len(text) && isSpace(text[i]) {
		i++
	}
	return i
}

// markdownLexer: headings, emphasis, code, links, quotes and lists
type markdownLexer struct{}

const (
	mdFence      = 1 // inside ``` code
	mdFenceTilde = 2 // inside ~~~ code
)

func (markdownLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	i := firstNonSpace(text)
	rest := text[i:]
	if state == mdFence || state == mdFenceTilde {
		out.Set(0, StyleCode)
		if (state == mdFence && bytes.HasPrefix(rest, []byte("```"))) || (state == mdFenceTilde && bytes.HasPrefix(rest, []byte("~~~"))) {
			out.Close()
			return 0
		}
		return state
	}
	switch {
	case bytes.HasPrefix(rest, []byte("```")):
		out.Set(0, StyleCode)
		out.Open()
		return mdFence
	case bytes.HasPrefix(rest, []byte("~~~")):
		out.Set(0, StyleCode)
		out.Open()
		return mdFenceTilde
	case len(rest) > 0 && rest[0] == '#':
		n := 0
		for n < len(rest) && rest[n] == '#' {
			n++
		}
		if n <= 6 && (n == len(rest) || rest[n] == ' ' || rest[n] == '\t') {
			out.Set(0, StyleHeading)
			return 0
		}
	case len(rest) > 0 && rest[0] == '>':
		out.Set(0, StyleQuote)
		return 0
	case isRule(rest):
		out.Set(0, StyleOperator)
		return 0
	}
	if n := listMarkerLen(rest); n > 0 {
		out.Span(i, i+n, StyleKeyword)
		i += n
	}
	mdInline(text, i, out)
	return 0
}

func isRule(s []byte) bool {
	s = bytes.TrimRight(s, " \t")
	if len(s) < 3 {
		return false
	}
	c := s[0]
	if c != '-' && c != '*' && c != '_' && c != '=' {
		return false
	}
	for _, x := range s {
		if x != c && x != ' ' {
			return false
		}
	}
	return true
}

func isListMarker(s []byte) bool { return listMarkerLen(s) > 0 }

func listMarkerLen(s []byte) int {
	if len(s) >= 2 && (s[0] == '-' || s[0] == '*' || s[0] == '+') && (s[1] == ' ' || s[1] == '\t') {
		return 1
	}
	j := 0
	for j < len(s) && j < 9 && isDigit(s[j]) {
		j++
	}
	if j > 0 && j+1 < len(s) && (s[j] == '.' || s[j] == ')') && (s[j+1] == ' ' || s[j+1] == '\t') {
		return j + 1
	}
	return 0
}

// mdInline styles `code`, **strong**, *emphasis*, [links](url) and <urls>
func mdInline(text []byte, i int, out *LineStyles) {
	for i < len(text) {
		c := text[i]
		switch {
		case c == '\\':
			i += 2
			continue
		case c == '`':
			n := 0
			for i+n < len(text) && text[i+n] == '`' {
				n++
			}
			fence := text[i : i+n]
			k := bytes.Index(text[i+n:], fence)
			if k >= 0 {
				out.Span(i, i+n+k+n, StyleCode)
				i += n + k + n
				continue
			}
			i += n
			continue
		case (c == '*' || c == '_') && i+1 < len(text) && text[i+1] == c:
			k := bytes.Index(text[i+2:], []byte{c, c})
			if k > 0 {
				out.Span(i, i+2+k+2, StyleStrong)
				i += k + 4
				continue
			}
		case c == '*' || (c == '_' && (i == 0 || !isIdentByte(text[i-1]))):
			k := bytes.IndexByte(text[i+1:], c)
			if k > 0 && text[i+1] != ' ' {
				out.Span(i, i+1+k+1, StyleEmphasis)
				i += k + 2
				continue
			}
		case c == '[':
			k := bytes.IndexByte(text[i:], ']')
			if k > 0 && i+k+1 < len(text) && text[i+k+1] == '(' {
				e := bytes.IndexByte(text[i+k+1:], ')')
				if e > 0 {
					out.Span(i, i+k+1+e+1, StyleLink)
					i += k + e + 2
					continue
				}
			}
		case c == '<' && (hasPrefixAt(text, i+1, "http") || hasPrefixAt(text, i+1, "mailto:")):
			k := bytes.IndexByte(text[i:], '>')
			if k > 0 {
				out.Span(i, i+k+1, StyleLink)
				i += k + 1
				continue
			}
		}
		i++
	}
}

// iniLexer: [sections], key = value and comments; a section folds up to the next one
type iniLexer struct {
	props bool // .properties: ':' separates too, '!' starts comments
}

func (l iniLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	i := firstNonSpace(text)
	if i >= len(text) {
		return state
	}
	c := text[i]
	switch {
	case c == ';' || c == '#' || (l.props && c == '!'):
		out.Set(i, StyleCommentLine)
		return state
	case c == '[':
		out.Set(i, StyleSection)
		if k := bytes.IndexByte(text[i:], ']'); k >= 0 {
			out.Set(i+k+1, StyleDefault)
		}
		if state == 1 {
			out.Close()
		}
		out.Open()
		return 1
	}
	sep := bytes.IndexByte(text[i:], '=')
	if l.props {
		if k := bytes.IndexByte(text[i:], ':'); k >= 0 && (sep < 0 || k < sep) {
			sep = k
		}
	}
	if sep >= 0 {
		out.Span(i, i+sep, StyleKey)
		out.Span(i+sep, i+sep+1, StyleOperator)
		out.Set(i+sep+1, StyleValue)
	}
	return state
}

// diffLexer: added and removed lines, hunks and file headers; files and
// hunks fold
type diffLexer struct{}

const (
	diffFileOpen = 1
	diffHunkOpen = 2
)

func (diffLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	closeHunk := func() {
		if state&diffHunkOpen != 0 {
			out.Close()
			state &^= diffHunkOpen
		}
	}
	switch {
	case bytes.HasPrefix(text, []byte("diff ")) || bytes.HasPrefix(text, []byte("Index: ")):
		closeHunk()
		if state&diffFileOpen != 0 {
			out.Close()
		}
		out.Open()
		state |= diffFileOpen
		out.Set(0, StyleDiffHeader)
	case bytes.HasPrefix(text, []byte("+++")) || bytes.HasPrefix(text, []byte("---")) ||
		bytes.HasPrefix(text, []byte("***")) || bytes.HasPrefix(text, []byte("index ")) ||
		bytes.HasPrefix(text, []byte("new file")) || bytes.HasPrefix(text, []byte("deleted file")) ||
		bytes.HasPrefix(text, []byte("similarity")) || bytes.HasPrefix(text, []byte("rename ")):
		out.Set(0, StyleDiffHeader)
	case bytes.HasPrefix(text, []byte("@@")):
		closeHunk()
		out.Open()
		state |= diffHunkOpen
		out.Set(0, StyleDiffPosition)
	case len(text) > 0 && (text[0] == '+' || text[0] == '>'):
		out.Set(0, StyleDiffAdded)
	case len(text) > 0 && (text[0] == '-' || text[0] == '<'):
		out.Set(0, StyleDiffRemoved)
	case len(text) > 0 && text[0] == '!':
		out.Set(0, StyleChar)
	}
	return state
}

// yamlLexer: keys, comments, strings, anchors and tags
type yamlLexer struct{}

var yamlConstants = wordSet("true false yes no on off null True False Yes No On Off Null TRUE FALSE YES NO ON OFF NULL ~")

func (yamlLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	i := firstNonSpace(text)
	if hasPrefixAt(text, 0, "---") || hasPrefixAt(text, 0, "...") {
		out.Set(0, StylePreprocessor)
		if k := bytes.IndexByte(text, '#'); k >= 0 {
			out.Set(k, StyleCommentLine)
		}
		return 0
	}
	// "- " list items
	for i+1 < len(text) && text[i] == '-' && (text[i+1] == ' ' || text[i+1] == '\t') {
		out.Span(i, i+1, StyleOperator)
		i += 2
		for i < len(text) && isSpace(text[i]) {
			i++
		}
	}
	if i < len(text) && text[i] == '#' {
		out.Set(i, StyleCommentLine)
		return 0
	}
	// The key
	keyEnd := yamlKeyEnd(text, i)
	if keyEnd > i {
		out.Span(i, keyEnd, StyleKey)
		out.Span(keyEnd, keyEnd+1, StyleOperator)
		i = keyEnd + 1
	}
	// The value
	for i < len(text) {
		c := text[i]
		switch {
		case isSpace(c):
			i++
			continue
		case c == '#' && (i == 0 || isSpace(text[i-1])):
			out.Set(i, StyleCommentLine)
			return 0
		case c == '"' || c == '\'':
			j := i + 1
			for j < len(text) && text[j] != c {
				if c == '"' && text[j] == '\\' {
					j++
				}
				j++
			}
			j = min(j+1, len(text))
			out.Span(i, j, StyleString)
			i = j
			continue
		case c == '&' || c == '*':
			j := i + 1
			for j < len(text) && !isSpace(text[j]) {
				j++
			}
			out.Span(i, j, StyleVariable)
			i = j
			continue
		case c == '!':
			j := i + 1
			for j < len(text) && !isSpace(text[j]) {
				j++
			}
			out.Span(i, j, StyleAnnotation)
			i = j
			continue
		case c == '|' || c == '>' || c == '[' || c == ']' || c == '{' || c == '}' || c == ',':
			out.Span(i, i+1, StyleOperator)
			i++
			continue
		}
		// A plain word up to a comment
		j := i
		for j < len(text) && !(text[j] == '#' && isSpace(text[j-1])) && text[j] != ',' && text[j] != ']' && text[j] != '}' {
			j++
		}
		word := bytes.TrimRight(text[i:j], " \t")
		switch {
		case yamlConstants[string(word)]:
			out.Span(i, i+len(word), StyleConstant)
		case len(word) > 0 && (isDigit(word[0]) || (word[0] == '-' && len(word) > 1 && isDigit(word[1]))) && isNumberLike(word):
			out.Span(i, i+len(word), StyleNumber)
		}
		i = j
		if i < len(text) && text[i] != '#' {
			out.Span(i, i+1, StyleOperator)
			i++
		}
	}
	return 0
}

func isNumberLike(w []byte) bool {
	for _, c := range w {
		if !(isDigit(c) || c == '.' || c == '-' || c == '+' || c == 'e' || c == 'E' || c == '_' || c == 'x' || c == 'o') {
			return false
		}
	}
	return true
}

// yamlKeyEnd returns where the key of "key: value" at i ends, or i if none
func yamlKeyEnd(text []byte, i int) int {
	j := i
	if j < len(text) && (text[j] == '"' || text[j] == '\'') {
		q := text[j]
		k := bytes.IndexByte(text[j+1:], q)
		if k < 0 {
			return i
		}
		j += k + 2
		if j < len(text) && text[j] == ':' {
			return j
		}
		return i
	}
	for j < len(text) {
		c := text[j]
		if c == ':' && (j+1 == len(text) || text[j+1] == ' ' || text[j+1] == '\t') {
			return j
		}
		if c == '#' && j > i && isSpace(text[j-1]) {
			return i
		}
		if c == '{' || c == '[' || c == '"' || c == '\'' {
			return i
		}
		j++
	}
	return i
}

// cssLexer: selectors, properties, values, comments and @-rules
type cssLexer struct{}

var cssAtRules = wordSet("@media @import @font-face @keyframes @charset @supports @page @namespace @layer @container @property")

func (l cssLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	return l.lexRange(text, 0, len(text), state, out)
}

// The state: the depth of the braces, and whether a comment is open
const cssInComment = 1 << 8

func (cssLexer) lexRange(text []byte, from, to int, state uint32, out *LineStyles) uint32 {
	depth := int(state & 0xFF)
	inComment := state&cssInComment != 0
	i := from
	pack := func() uint32 {
		s := uint32(min(depth, 255))
		if inComment {
			s |= cssInComment
		}
		return s
	}
	for i < to {
		if inComment {
			out.Set(i, StyleComment)
			k := bytes.Index(text[i:to], []byte("*/"))
			if k < 0 {
				return pack()
			}
			i += k + 2
			out.Set(i, StyleDefault)
			out.Close()
			inComment = false
			continue
		}
		c := text[i]
		switch {
		case isSpace(c):
			i++
		case hasPrefixAt(text, i, "/*"):
			out.Set(i, StyleComment)
			out.Open()
			inComment = true
			i += 2
		case c == '{':
			out.Span(i, i+1, StyleOperator)
			out.Open()
			depth++
			i++
		case c == '}':
			out.Span(i, i+1, StyleOperator)
			out.Close()
			depth = max(0, depth-1)
			i++
		case c == '"' || c == '\'':
			j := i + 1
			for j < to && text[j] != c {
				if text[j] == '\\' {
					j++
				}
				j++
			}
			j = min(j+1, to)
			out.Span(i, j, StyleString)
			i = j
		case c == '@':
			j := i + 1
			for j < to && (isIdentByte(text[j]) || text[j] == '-') {
				j++
			}
			st := StyleKeyword
			if cssAtRules[string(text[i:j])] {
				st = StylePreprocessor
			}
			out.Span(i, j, st)
			i = j
		case c == '!' && hasPrefixFoldAt(text, i, "!important"):
			out.Span(i, i+10, StyleKeyword)
			i += 10
		case c == '#' && depth > 0:
			j := i + 1
			for j < to && isIdentByte(text[j]) {
				j++
			}
			out.Span(i, j, StyleNumber)
			i = j
		case isDigit(c) || (c == '.' && i+1 < to && isDigit(text[i+1])) || (c == '-' && i+1 < to && isDigit(text[i+1]) && depth > 0):
			j := i + 1
			for j < to && (isIdentByte(text[j]) || text[j] == '.' || text[j] == '%') {
				j++
			}
			out.Span(i, j, StyleNumber)
			i = j
		case isIdentByte(c) || c == '-' || ((c == '.' || c == '#' || c == ':') && depth == 0):
			j := i + 1
			for j < to && (isIdentByte(text[j]) || text[j] == '-') {
				j++
			}
			word := text[i:j]
			switch {
			case depth > 0 && followedByColon(text, j, to):
				out.Span(i, j, StyleKey)
			case depth > 0 && followedBy(text, j, to, '('):
				out.Span(i, j, StyleBuiltin)
			case depth > 0:
				out.Span(i, j, StyleValue)
			case c == '.' || c == '#':
				out.Span(i, j, StyleAttribute)
			case c == ':':
				out.Span(i, j, StyleKeyword)
			default:
				out.Span(i, j, StyleTag)
			}
			_ = word
			i = j
		default:
			if c < 0x80 {
				out.Span(i, i+1, StyleOperator)
			}
			i++
		}
	}
	return pack()
}

// searchResultsLexer styles the results of Find All: the search, the files
// and the lines found; each folds
type searchResultsLexer struct{}

const (
	srSearchOpen = 1
	srFileOpen   = 2
)

func (searchResultsLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	closeFile := func() {
		if state&srFileOpen != 0 {
			out.Close()
			state &^= srFileOpen
		}
	}
	switch {
	case len(text) > 0 && text[0] != ' ' && text[0] != '\t':
		closeFile()
		if state&srSearchOpen != 0 {
			out.Close()
		}
		out.Open()
		state |= srSearchOpen
		out.Set(0, StyleSearchHeader)
	case bytes.HasPrefix(text, []byte("  ")) && len(text) > 2 && text[2] != ' ':
		closeFile()
		out.Open()
		state |= srFileOpen
		out.Set(0, StyleSearchFile)
	case bytes.HasPrefix(text, []byte("\t")):
		if k := bytes.IndexByte(text, ':'); k > 0 {
			out.Span(0, k+1, StyleSearchLineNo)
		}
	}
	return state
}

// parseSearchResultLine returns the line number of a "\tLine 12: text" line of the results
func parseSearchResultLine(text []byte) (int, bool) {
	if !bytes.HasPrefix(text, []byte("\t")) {
		return 0, false
	}
	k := bytes.IndexByte(text, ':')
	if k < 0 {
		return 0, false
	}
	fields := strings.Fields(string(text[1:k]))
	if len(fields) == 0 {
		return 0, false
	}
	n := 0
	for _, c := range fields[len(fields)-1] {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}
