package editor

import (
	"bytes"
	"strings"
)

// markupLexer styles XML and HTML; in HTML the scripts, styles and PHP code
// are styled by the lexers of their languages
type markupLexer struct {
	html   bool
	script subLexer
	style  subLexer
	php    subLexer
}

// Markup lexer states: the mode in the high bits, the state of the
// embedded language in the low ones
const (
	mkText = iota
	mkTag
	mkAttrDQ
	mkAttrSQ
	mkComment
	mkCDATA
	mkScript
	mkStyle
	mkPHP
	mkDecl // <!DOCTYPE ...>, <?xml ... ?>

	mkFlagScript = 1 << 27 // the tag being read is <script>
	mkFlagStyle  = 1 << 26 // the tag being read is <style>
	mkFlagOpened = 1 << 25 // the tag being read opened a fold
	mkSubMask    = 1<<24 - 1
)

func mkPack(mode int, flags uint32, sub uint32) uint32 {
	return uint32(mode)<<28 | flags | sub&mkSubMask
}

var htmlVoidElements = wordSet("area base br col embed hr img input link meta param source track wbr")

func (m *markupLexer) LexLine(text []byte, state uint32, out *LineStyles) uint32 {
	mode := int(state >> 28)
	flags := state & (mkFlagScript | mkFlagStyle | mkFlagOpened)
	sub := state & mkSubMask
	to := len(text)
	i := 0

	for i < to {
		switch mode {
		case mkComment:
			out.Set(i, StyleComment)
			k := bytes.Index(text[i:], []byte("-->"))
			if k < 0 {
				return mkPack(mkComment, 0, 0)
			}
			i += k + 3
			out.Set(i, StyleDefault)
			out.Close()
			mode = mkText

		case mkCDATA:
			out.Set(i, StyleCDATA)
			k := bytes.Index(text[i:], []byte("]]>"))
			if k < 0 {
				return mkPack(mkCDATA, 0, 0)
			}
			i += k + 3
			out.Set(i, StyleDefault)
			out.Close()
			mode = mkText

		case mkDecl:
			out.Set(i, StylePreprocessor)
			k := bytes.IndexByte(text[i:], '>')
			if k < 0 {
				return mkPack(mkDecl, 0, 0)
			}
			i += k + 1
			out.Set(i, StyleDefault)
			mode = mkText

		case mkScript, mkStyle:
			closeTag := "</script"
			lexer := m.script
			if mode == mkStyle {
				closeTag, lexer = "</style", m.style
			}
			end := indexFold(text[i:], closeTag)
			if end < 0 {
				end = to
			} else {
				end += i
			}
			// PHP inside a script
			if m.php != nil {
				if k := bytes.Index(text[i:end], []byte("<?")); k >= 0 {
					end = i + k
				}
			}
			sub = lexer.lexRange(text, i, end, sub, out)
			out.Set(end, StyleDefault)
			i = end
			if i >= to {
				return mkPack(mode, 0, sub)
			}
			if m.php != nil && hasPrefixAt(text, i, "<?") {
				// Back to the script after the PHP code: kept in the flags
				if mode == mkScript {
					flags = mkFlagScript
				} else {
					flags = mkFlagStyle
				}
				i = m.phpOpen(text, i, out)
				mode = mkPHP
				sub = 0
				continue
			}
			mode = mkText
			sub = 0

		case mkPHP:
			end := bytes.Index(text[i:], []byte("?>"))
			if end < 0 {
				end = to
			} else {
				end += i
			}
			sub = m.php.lexRange(text, i, end, sub, out)
			i = end
			if i >= to {
				return mkPack(mkPHP, flags, sub)
			}
			out.Span(i, i+2, StylePreprocessor)
			i += 2
			sub = 0
			switch {
			case flags&mkFlagScript != 0:
				mode = mkScript
			case flags&mkFlagStyle != 0:
				mode = mkStyle
			default:
				mode = mkText
			}
			flags = 0

		case mkAttrDQ, mkAttrSQ:
			q := byte('"')
			if mode == mkAttrSQ {
				q = '\''
			}
			out.Set(i, StyleString)
			k := bytes.IndexByte(text[i:], q)
			if k < 0 {
				return mkPack(mode, flags, 0)
			}
			i += k + 1
			out.Set(i, StyleDefault)
			mode = mkTag

		case mkTag:
			c := text[i]
			switch {
			case isSpace(c):
				i++
			case c == '>':
				out.Span(i, i+1, StyleTag)
				i++
				mode = mkText
				switch {
				case flags&mkFlagScript != 0:
					mode = mkScript
				case flags&mkFlagStyle != 0:
					mode = mkStyle
				}
				flags = 0
				sub = 0
			case c == '/' && i+1 < to && text[i+1] == '>':
				out.Span(i, i+2, StyleTag)
				if flags&mkFlagOpened != 0 {
					out.Close() // opened by the tag name
				}
				i += 2
				mode = mkText
				flags = 0
			case c == '"':
				mode = mkAttrDQ
				out.Set(i, StyleString)
				i++
			case c == '\'':
				mode = mkAttrSQ
				out.Set(i, StyleString)
				i++
			case c == '=':
				out.Span(i, i+1, StyleOperator)
				i++
			case m.php != nil && hasPrefixAt(text, i, "<?"):
				i = m.phpOpen(text, i, out)
				mode = mkPHP
				sub = 0
			default:
				j := i
				for j < to && !isSpace(text[j]) && text[j] != '=' && text[j] != '>' && text[j] != '/' && text[j] != '"' && text[j] != '\'' {
					j++
				}
				if j == i {
					j++
				}
				out.Span(i, j, StyleAttribute)
				i = j
			}

		default: // mkText
			k := i
			for k < to && text[k] != '<' && text[k] != '&' {
				k++
			}
			i = k
			if i >= to {
				break
			}
			if text[i] == '&' {
				j := i + 1
				for j < to && j-i < 32 && (isIdentByte(text[j]) || text[j] == '#') {
					j++
				}
				if j < to && text[j] == ';' && j > i+1 {
					out.Span(i, j+1, StyleEntity)
					i = j + 1
				} else {
					i++
				}
				continue
			}
			switch {
			case hasPrefixAt(text, i, "<!--"):
				out.Set(i, StyleComment)
				out.Open()
				i += 4
				mode = mkComment
			case hasPrefixAt(text, i, "<![CDATA["):
				out.Set(i, StyleCDATA)
				out.Open()
				i += 9
				mode = mkCDATA
			case m.php != nil && hasPrefixAt(text, i, "<?"):
				i = m.phpOpen(text, i, out)
				mode = mkPHP
				sub = 0
			case hasPrefixAt(text, i, "<!") || hasPrefixAt(text, i, "<?"):
				out.Set(i, StylePreprocessor)
				i += 2
				mode = mkDecl
			default:
				j := i + 1
				closing := j < to && text[j] == '/'
				if closing {
					j++
				}
				nameStart := j
				for j < to && (isIdentByte(text[j]) || text[j] == ':' || text[j] == '-' || text[j] == '.') {
					j++
				}
				if j == nameStart {
					// A lone '<'
					i++
					continue
				}
				out.Span(i, j, StyleTag)
				name := strings.ToLower(string(text[nameStart:j]))
				i = j
				flags = 0
				if closing {
					out.Close()
					// The rest of a closing tag up to '>'
					mode = mkTag
					continue
				}
				if !(m.html && htmlVoidElements[name]) {
					out.Open()
					flags |= mkFlagOpened
				}
				if m.html && name == "script" {
					flags |= mkFlagScript
				} else if m.html && name == "style" {
					flags |= mkFlagStyle
				}
				mode = mkTag
			}
		}
	}
	return mkPack(mode, flags, sub)
}

// phpOpen styles "<?php" at i and returns where the code starts
func (m *markupLexer) phpOpen(text []byte, i int, out *LineStyles) int {
	j := i + 2
	if hasPrefixFoldAt(text, j, "php") {
		j += 3
	} else if hasPrefixAt(text, j, "=") {
		j++
	}
	out.Span(i, j, StylePreprocessor)
	return j
}

// indexFold is bytes.Index ignoring the case of ASCII letters
func indexFold(s []byte, sub string) int {
	n := len(sub)
	for i := 0; i+n <= len(s); i++ {
		if s[i] == '<' && strings.EqualFold(string(s[i:i+n]), sub) {
			return i
		}
	}
	return -1
}
