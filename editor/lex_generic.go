package editor

import (
	"bytes"
	"strings"
)

// genericDef describes a language for genericLexer: most programming
// languages differ only in their comments, strings and words
type genericDef struct {
	lineComments []string
	// docLineComments are styled as doc comments, e.g. /// of C#
	docLineComments []string
	// lineCommentsAtStart are comments only at the start of a line, e.g. REM of batch files
	lineCommentsAtStart []string
	// lineCommentAfterSpace: a line comment only starts after a space, e.g. # of shell scripts
	lineCommentAfterSpace bool
	blockComments         []blockDef
	strings               []stringDef
	// stringPrefixes are letters that may come right before a string: r"..." of Python
	stringPrefixes string

	words           map[string]Style
	caseInsensitive bool
	identExtra      string // characters allowed in words besides letters, digits and _
	identStartExtra string // characters words may start with

	// preprocessor: the character of preprocessor lines, e.g. '#' of C
	preprocessor byte
	// varPrefixes: "$" styles $name, ${name}; "%" styles %name% of batch files
	varPrefixes string
	// annotation: '@' styles @Name of Java and Python decorators
	annotation byte

	foldBraces      string // pairs of fold characters, e.g. "{}"
	foldOpen        map[string]bool
	foldClose       map[string]bool
	foldMid         map[string]bool // close and open again, e.g. else
	foldNoOpenAfter map[string]bool // no fold opens after these words, e.g. End If

	// keyBeforeColon: a string or word followed by ':' is a key, as in JSON
	keyBeforeColon bool
	regexLiterals  bool
	// functions: a word followed by '(' is styled as a function
	functions bool
	// labels: a word at the start of a line followed by ':' is a label, e.g. in assembler
	labels bool
	// labelPrefix: a line starting with it is a label, e.g. ':' of batch files
	labelPrefix byte
	// commandPrefix: a word starting with it is a keyword, e.g. \section of TeX
	commandPrefix byte
	// varParens: $(name) is a variable, as in makefiles
	varParens bool
}

type blockDef struct {
	open, close string
	style       Style
	nested      bool
	// notFollowedBy: the open does not count when this character follows,
	// e.g. "/**/" is not a doc comment
	notFollowedBy byte
}

type stringDef struct {
	open, close string
	escape      byte // 0 - no escapes
	multiline   bool
	style       Style
	// doubledClose: the close twice is an escaped close, e.g. '' in SQL
	doubledClose bool
}

// Generic lexer states
const (
	gsDefault = iota
	gsComment
	gsString
)

func gsPack(kind, index, depth int) uint32 {
	return uint32(kind) | uint32(index)<<4 | uint32(depth)<<8
}

func gsUnpack(s uint32) (kind, index, depth int) {
	return int(s & 0xF), int(s>>4) & 0xF, int(s>>8) & 0xFFFF
}

type genericLexer struct {
	def *genericDef
}

func newGeneric(def *genericDef) *genericLexer {
	if def.caseInsensitive {
		lower := func(m map[string]bool) map[string]bool {
			out := make(map[string]bool, len(m))
			for k, v := range m {
				out[strings.ToLower(k)] = v
			}
			return out
		}
		def.foldOpen, def.foldClose = lower(def.foldOpen), lower(def.foldClose)
		def.foldMid, def.foldNoOpenAfter = lower(def.foldMid), lower(def.foldNoOpenAfter)
		words := make(map[string]Style, len(def.words))
		for k, v := range def.words {
			words[strings.ToLower(k)] = v
		}
		def.words = words
	}
	return &genericLexer{def: def}
}

func (g *genericLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	return g.lexRange(text, 0, len(text), state, out)
}

func isIdentByte(c byte) bool {
	return c == '_' || c >= 0x80 || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\v' }

func hasPrefixAt(text []byte, i int, s string) bool {
	return len(s) > 0 && len(text)-i >= len(s) && string(text[i:i+len(s)]) == s
}

func hasPrefixFoldAt(text []byte, i int, s string) bool {
	return len(s) > 0 && len(text)-i >= len(s) && strings.EqualFold(string(text[i:i+len(s)]), s)
}

// lexRange styles text[from:to], a part of a line (the whole of it, unless
// the language is embedded in another one, like JavaScript in HTML)
func (g *genericLexer) lexRange(text []byte, from, to int, state uint32, out *LineStyles) uint32 {
	d := g.def
	i := from
	kind, index, depth := gsUnpack(state)
	firstNonSpace := true
	prevSignificant := byte(0) // for the regex literals
	prevWord := ""
	var lowerBuf []byte

	// Continue what the previous line left open
	switch kind {
	case gsComment:
		if index < len(d.blockComments) {
			end, left := g.blockComment(text, i, to, index, depth, out)
			if end < 0 {
				return gsPack(gsComment, index, left)
			}
			i = end
		}
	case gsString:
		if index < len(d.strings) {
			sd := d.strings[index]
			out.Set(i, sd.style)
			end, closed := g.stringEnd(text, i, to, sd)
			i = end
			out.Set(i, StyleDefault)
			if !closed {
				return state
			}
			firstNonSpace = false
		}
	}

	for i < to {
		c := text[i]
		if isSpace(c) {
			i++
			continue
		}
		atStart := firstNonSpace
		firstNonSpace = false

		// Preprocessor directive: the character and the word after it
		if d.preprocessor != 0 && c == d.preprocessor && atStart {
			j := i + 1
			for j < to && isSpace(text[j]) {
				j++
			}
			for j < to && isIdentByte(text[j]) {
				j++
			}
			out.Span(i, j, StylePreprocessor)
			i = j
			prevSignificant = 'a'
			continue
		}

		if d.labelPrefix != 0 && c == d.labelPrefix && atStart && i+1 < to && text[i+1] != d.labelPrefix {
			out.Span(i, to, StyleLabel)
			return 0
		}

		// Comments at the start of a line
		if atStart {
			found := false
			for _, lc := range d.lineCommentsAtStart {
				if hasPrefixFoldAt(text, i, lc) && (i+len(lc) >= to || !isIdentByte(text[i+len(lc)]) || !isIdentByte(lc[len(lc)-1])) {
					found = true
					break
				}
			}
			if found {
				out.Span(i, to, StyleCommentLine)
				return 0
			}
		}

		// Block comments
		matched := false
		for bi, bc := range d.blockComments {
			if hasPrefixAt(text, i, bc.open) && !(bc.notFollowedBy != 0 && i+len(bc.open) < to && text[i+len(bc.open)] == bc.notFollowedBy) {
				out.Set(i, bc.style)
				out.Open()
				end, left := g.blockComment(text, i+len(bc.open), to, bi, 1, out)
				if end < 0 {
					// Left open: the state tells which comment
					return gsPack(gsComment, bi, left)
				}
				i = end
				matched = true
				break
			}
		}
		if matched {
			continue
		}

		// Line comments
		for _, lc := range d.docLineComments {
			if hasPrefixAt(text, i, lc) {
				out.Span(i, to, StyleCommentDoc)
				return 0
			}
		}
		for _, lc := range d.lineComments {
			if hasPrefixAt(text, i, lc) {
				if d.lineCommentAfterSpace && i > from && !isSpace(text[i-1]) {
					continue
				}
				out.Span(i, to, StyleCommentLine)
				return 0
			}
		}

		// Strings
		if si := g.stringAt(text, i); si >= 0 {
			sd := d.strings[si]
			start := i
			out.Set(i, sd.style)
			end, closed := g.stringEnd(text, i+len(sd.open), to, sd)
			i = end
			out.Set(i, StyleDefault)
			if !closed {
				if sd.multiline {
					return gsPack(gsString, si, 0)
				}
				return 0
			}
			if d.keyBeforeColon && followedByColon(text, i, to) {
				out.Span(start, i, StyleKey)
			}
			prevSignificant = '"'
			continue
		}

		// Numbers
		if isDigit(c) || (c == '.' && i+1 < to && isDigit(text[i+1])) {
			j := i + 1
			hex := c == '0' && j < to && (text[j] == 'x' || text[j] == 'X')
			for j < to {
				ch := text[j]
				if isIdentByte(ch) || ch == '.' {
					j++
					continue
				}
				if (ch == '+' || ch == '-') && !hex && (text[j-1] == 'e' || text[j-1] == 'E') {
					j++
					continue
				}
				break
			}
			out.Span(i, j, StyleNumber)
			i = j
			prevSignificant = '0'
			continue
		}

		// Variables
		if d.varPrefixes != "" && strings.IndexByte(d.varPrefixes, c) >= 0 {
			if j := g.variableEnd(text, i, to, c); j > i+1 || (j > i && c != '%') {
				out.Span(i, j, StyleVariable)
				i = j
				prevSignificant = 'a'
				continue
			}
		}

		// Commands
		if d.commandPrefix != 0 && c == d.commandPrefix && i+1 < to {
			j := i + 1
			if isIdentByte(text[j]) {
				for j < to && isIdentByte(text[j]) && text[j] < 0x80 {
					j++
				}
			} else {
				j++
			}
			out.Span(i, j, StyleKeyword)
			i = j
			prevSignificant = 'a'
			continue
		}

		// Annotations
		if d.annotation != 0 && c == d.annotation && i+1 < to && isIdentByte(text[i+1]) {
			j := i + 1
			for j < to && (isIdentByte(text[j]) || text[j] == '.') {
				j++
			}
			out.Span(i, j, StyleAnnotation)
			i = j
			prevSignificant = 'a'
			continue
		}

		// Words
		if isIdentByte(c) || (d.identStartExtra != "" && strings.IndexByte(d.identStartExtra, c) >= 0) {
			j := i + 1
			for j < to && (isIdentByte(text[j]) || (d.identExtra != "" && strings.IndexByte(d.identExtra, text[j]) >= 0)) {
				j++
			}
			word := text[i:j]
			// A string prefix: r"..."
			if d.stringPrefixes != "" && j < to && len(word) <= 3 && allIn(word, d.stringPrefixes) {
				if si := g.stringAt(text, j); si >= 0 {
					sd := d.strings[si]
					out.Set(i, sd.style)
					end, closed := g.stringEnd(text, j+len(sd.open), to, sd)
					i = end
					out.Set(i, StyleDefault)
					if !closed {
						if sd.multiline {
							return gsPack(gsString, si, 0)
						}
						return 0
					}
					prevSignificant = '"'
					continue
				}
			}
			key := word
			if d.caseInsensitive {
				lowerBuf = append(lowerBuf[:0], word...)
				lowerASCII(lowerBuf)
				key = lowerBuf
			}
			st, isWord := d.words[string(key)]
			switch {
			case isWord:
				out.Span(i, j, st)
			case d.keyBeforeColon && followedByColon(text, j, to):
				out.Span(i, j, StyleKey)
			case d.labels && atStart && j < to && text[j] == ':':
				out.Span(i, j+1, StyleLabel)
				j++
			case d.functions && followedBy(text, j, to, '('):
				out.Span(i, j, StyleFunction)
			}
			if d.foldOpen != nil || d.foldMid != nil {
				w := string(key)
				switch {
				case d.foldMid[w]:
					out.Close()
					out.Open()
				case d.foldOpen[w] && !d.foldNoOpenAfter[prevWord]:
					out.Open()
				case d.foldClose[w]:
					out.Close()
				}
				prevWord = w
			}
			i = j
			prevSignificant = 'a'
			continue
		}

		// Regular expressions
		if c == '/' && d.regexLiterals && regexAllowedAfter(prevSignificant) {
			if j := regexEnd(text, i+1, to); j > 0 {
				out.Span(i, j, StyleRegex)
				i = j
				prevSignificant = '"'
				continue
			}
		}

		// Operators and fold braces
		if k := strings.IndexByte(d.foldBraces, c); k >= 0 {
			if k%2 == 0 {
				out.Open()
			} else {
				out.Close()
			}
		}
		if c < 0x80 {
			out.Span(i, i+1, StyleOperator)
		}
		prevSignificant = c
		i++
	}
	return 0
}

func allIn(word []byte, set string) bool {
	for _, c := range word {
		if strings.IndexByte(set, c) < 0 {
			return false
		}
	}
	return true
}

func lowerASCII(b []byte) {
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
}

func followedByColon(text []byte, i, to int) bool {
	for i < to && isSpace(text[i]) {
		i++
	}
	return i < to && text[i] == ':' && (i+1 >= to || text[i+1] != ':')
}

func followedBy(text []byte, i, to int, c byte) bool {
	for i < to && isSpace(text[i]) {
		i++
	}
	return i < to && text[i] == c
}

// stringAt returns which string starts at i, -1 if none
func (g *genericLexer) stringAt(text []byte, i int) int {
	for si, sd := range g.def.strings {
		if hasPrefixAt(text, i, sd.open) {
			return si
		}
	}
	return -1
}

// stringEnd returns where the string that goes on at i ends, and whether it
// is closed on the line
func (g *genericLexer) stringEnd(text []byte, i, to int, sd stringDef) (int, bool) {
	for i < to {
		c := text[i]
		if sd.escape != 0 && c == sd.escape {
			i += 2
			continue
		}
		if hasPrefixAt(text, i, sd.close) {
			if sd.doubledClose && hasPrefixAt(text, i+len(sd.close), sd.close) {
				i += 2 * len(sd.close)
				continue
			}
			return i + len(sd.close), true
		}
		i++
	}
	return min(i, to), false
}

// blockComment styles the comment that goes on at i; returns where it ends,
// or -1 and the nesting depth if it goes on to the next line
func (g *genericLexer) blockComment(text []byte, i, to, index, depth int, out *LineStyles) (int, int) {
	bc := g.def.blockComments[index]
	out.Set(i, bc.style)
	for i < to {
		if bc.nested && hasPrefixAt(text, i, bc.open) {
			depth++
			i += len(bc.open)
			continue
		}
		if hasPrefixAt(text, i, bc.close) {
			depth--
			i += len(bc.close)
			if depth <= 0 || !bc.nested {
				out.Set(i, StyleDefault)
				out.Close()
				return i, 0
			}
			continue
		}
		i++
	}
	return -1, max(depth, 1)
}

// variableEnd returns where the variable starting at i ends
func (g *genericLexer) variableEnd(text []byte, i, to int, prefix byte) int {
	j := i + 1
	if j >= to {
		return i
	}
	if prefix == '%' {
		// %name% and %1
		if isDigit(text[j]) || text[j] == '*' {
			return j + 1
		}
		if text[j] == '~' {
			for j < to && text[j] != ' ' && !isDigit(text[j]) {
				j++
			}
			if j < to && isDigit(text[j]) {
				return j + 1
			}
			return i
		}
		k := bytes.IndexByte(text[j:to], '%')
		if k <= 0 || bytes.IndexByte(text[j:j+k], ' ') >= 0 {
			return i
		}
		return j + k + 1
	}
	switch c := text[j]; {
	case c == '{':
		k := bytes.IndexByte(text[j:to], '}')
		if k < 0 {
			return to
		}
		return j + k + 1
	case c == '(':
		if !g.def.varParens {
			return j
		}
		k := bytes.IndexByte(text[j:to], ')')
		if k < 0 {
			return to
		}
		return j + k + 1
	case isIdentByte(c):
		for j < to && isIdentByte(text[j]) {
			j++
		}
		return j
	case strings.IndexByte("#?@*!$-0123456789_", c) >= 0:
		return j + 1
	}
	return i
}

func regexAllowedAfter(prev byte) bool {
	return prev == 0 || strings.IndexByte("(,=:[!&|?{};+-*%<>~^", prev) >= 0
}

// regexEnd returns where a regular expression literal starting before i ends, or -1
func regexEnd(text []byte, i, to int) int {
	inClass := false
	for i < to {
		c := text[i]
		switch {
		case c == '\\':
			i += 2
			continue
		case c == '[':
			inClass = true
		case c == ']':
			inClass = false
		case c == '/' && !inClass:
			i++
			for i < to && (text[i] >= 'a' && text[i] <= 'z') {
				i++
			}
			return i
		}
		i++
	}
	return -1
}

// wordSet makes a set of the words in the strings
func wordSet(lists ...string) map[string]bool {
	m := make(map[string]bool)
	for _, l := range lists {
		for _, w := range strings.Fields(l) {
			m[w] = true
		}
	}
	return m
}

// wordStyles makes the styles of the words: pairs of a style and the words
func wordStyles(pairs ...any) map[string]Style {
	m := make(map[string]Style)
	for k := 0; k+1 < len(pairs); k += 2 {
		st := pairs[k].(Style)
		for _, w := range strings.Fields(pairs[k+1].(string)) {
			m[w] = st
		}
	}
	return m
}
