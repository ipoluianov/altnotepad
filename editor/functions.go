package editor

import (
	"regexp"
	"strings"
)

// FuncItem is an entry of the function list: a function, a class, a heading...
type FuncItem struct {
	Name  string
	Line  int // 0-based
	Level int // nesting, by the indentation or the heading level
	Kind  string
}

// funcRule finds the entries of a language: the name is the group named
// "name" of the expression
type funcRule struct {
	re   *regexp.Regexp
	kind string
}

func rules(kind string, exprs ...string) []funcRule {
	out := make([]funcRule, len(exprs))
	for i, e := range exprs {
		out[i] = funcRule{re: regexp.MustCompile(e), kind: kind}
	}
	return out
}

var cLikeFunc = `^[\w\s\*&:<>,\[\]~.]*?\b(?P<name>~?[A-Za-z_]\w*(?:::~?\w+)*)\s*\([^;{}]*\)\s*(?:const\s*)?(?:noexcept\s*)?(?:override\s*)?(?:throws [\w, .]+)?\s*(?:\{.*)?$`

var funcRules = map[string][]funcRule{
	"c":    rules("function", cLikeFunc),
	"cpp":  append(rules("class", `^\s*(?:class|struct|namespace)\s+(?P<name>\w+)[^;]*$`), rules("function", cLikeFunc)...),
	"objc": append(rules("class", `^\s*@(?:interface|implementation)\s+(?P<name>\w+)`), rules("function", `^\s*[-+]\s*\([^)]*\)\s*(?P<name>\w+)`, cLikeFunc)...),
	"cs": append(rules("class", `^\s*(?:[\w\s]*\s)?(?:class|struct|interface|enum|record)\s+(?P<name>\w+)`),
		rules("function", `^\s*(?:(?:public|private|protected|internal|static|virtual|override|async|abstract|sealed|extern|unsafe|new|partial)\s+)+[\w<>\[\],.?\s]*?\b(?P<name>\w+)\s*(?:<[^>]*>)?\s*\([^;]*$`)...),
	"java": append(rules("class", `^\s*(?:[\w\s]*\s)?(?:class|interface|enum|record)\s+(?P<name>\w+)`),
		rules("function", `^\s*(?:(?:public|private|protected|static|final|abstract|synchronized|native|default)\s+)+[\w<>\[\],.?\s]*?\b(?P<name>\w+)\s*\([^;]*$`)...),
	"javascript": rules("function", `\bfunction\s*\*?\s*(?P<name>[\w$]+)\s*\(`, `^\s*(?:export\s+)?(?:default\s+)?class\s+(?P<name>[\w$]+)`,
		`^\s*(?:export\s+)?(?:const|let|var)\s+(?P<name>[\w$]+)\s*=\s*(?:async\s*)?(?:function|\([^)]*\)\s*=>|[\w$]+\s*=>)`,
		`^\s+(?:static\s+)?(?:async\s+)?(?:get\s+|set\s+)?(?P<name>(?:[\w$]+))\s*\([^)]*\)\s*\{`),
	"typescript": rules("function", `\bfunction\s*\*?\s*(?P<name>[\w$]+)\s*[<(]`, `^\s*(?:export\s+)?(?:default\s+)?(?:abstract\s+)?(?:class|interface|enum)\s+(?P<name>[\w$]+)`,
		`^\s*(?:export\s+)?(?:const|let|var)\s+(?P<name>[\w$]+)\s*(?::[^=]+)?=\s*(?:async\s*)?(?:function|\([^)]*\)\s*(?::[^=]+)?=>)`,
		`^\s+(?:(?:public|private|protected|static|readonly|async|abstract|override)\s+)*(?:get\s+|set\s+)?(?P<name>[\w$]+)\s*(?:<[^>]*>)?\([^)]*\)\s*(?::[^{]+)?\{`),
	"go":           rules("function", `^func\s+(?:\([^)]*\)\s*)?(?P<name>\w+)`, `^type\s+(?P<name>\w+)\s+(?:struct|interface)`),
	"rust":         rules("function", `^\s*(?:pub(?:\([^)]*\))?\s+)?(?:const\s+)?(?:async\s+)?(?:unsafe\s+)?(?:extern\s+"[^"]*"\s+)?fn\s+(?P<name>\w+)`, `^\s*(?:pub(?:\([^)]*\))?\s+)?(?:struct|enum|trait|mod|union)\s+(?P<name>\w+)`, `^\s*impl(?:<[^>]*>)?\s+(?P<name>[\w:<>]+(?:\s+for\s+[\w:<>]+)?)`),
	"python":       rules("function", `^\s*(?:async\s+)?def\s+(?P<name>\w+)`, `^\s*class\s+(?P<name>\w+)`),
	"ruby":         rules("function", `^\s*def\s+(?P<name>[\w.?!=]+)`, `^\s*(?:class|module)\s+(?P<name>[\w:]+)`),
	"perl":         rules("function", `^\s*sub\s+(?P<name>\w+)`, `^\s*package\s+(?P<name>[\w:]+)`),
	"lua":          rules("function", `\bfunction\s+(?P<name>[\w.:]+)`, `^\s*(?:local\s+)?(?P<name>[\w.]+)\s*=\s*function`),
	"bash":         rules("function", `^\s*(?:function\s+)?(?P<name>[\w-]+)\s*\(\)`, `^\s*function\s+(?P<name>[\w-]+)`),
	"powershell":   rules("function", `(?i)^\s*(?:function|filter|workflow)\s+(?P<name>[\w-]+)`, `(?i)^\s*class\s+(?P<name>\w+)`),
	"php":          rules("function", `\bfunction\s+(?P<name>\w+)\s*\(`, `^\s*(?:abstract\s+|final\s+)?(?:class|interface|trait)\s+(?P<name>\w+)`),
	"pascal":       rules("function", `(?i)^\s*(?:class\s+)?(?:procedure|function|constructor|destructor)\s+(?P<name>[\w.]+)`),
	"vb":           rules("function", `(?i)^\s*(?:(?:public|private|friend|protected|shared|overrides|overridable|static)\s+)*(?:sub|function|property)\s+(?:get\s+|let\s+|set\s+)?(?P<name>\w+)`, `(?i)^\s*(?:(?:public|private|friend)\s+)?(?:class|module|structure)\s+(?P<name>\w+)`),
	"sql":          rules("function", `(?i)\bcreate\s+(?:or\s+replace\s+)?(?:procedure|function|view|table|trigger|index|package)\s+(?:if\s+not\s+exists\s+)?(?P<name>[\w.\[\]"`+"`"+`]+)`),
	"markdown":     rules("heading", `^(?P<level>#{1,6})\s+(?P<name>.+?)\s*#*$`),
	"ini":          rules("section", `^\s*\[(?P<name>[^\]]+)\]`),
	"innosetup":    rules("section", `^\s*\[(?P<name>[^\]]+)\]`),
	"css":          rules("selector", `^\s*(?P<name>[^{}@;/\s][^{};]*?)\s*\{`, `^\s*(?P<name>@media[^{]*)\{`),
	"batch":        rules("label", `^:(?P<name>[\w.-]+)`),
	"makefile":     rules("target", `^(?P<name>[\w./%$() -]+?)\s*::?(?:[^=]|$)`),
	"yaml":         rules("key", `^(?P<name>[\w"'.-][^:#]*):(?:\s|$)`),
	"haskell":      rules("function", `^(?P<name>[a-z_]\w*'?)\s*::`, `^(?:data|newtype|type|class)\s+(?P<name>\w+)`),
	"kotlin":       rules("function", `\bfun\s+(?:<[^>]*>\s*)?(?:[\w.]+\.)?(?P<name>\w+)`, `^\s*(?:[\w\s]*\s)?(?:class|object|interface)\s+(?P<name>\w+)`),
	"swift":        rules("function", `\bfunc\s+(?P<name>\w+)`, `^\s*(?:[\w\s]*\s)?(?:class|struct|enum|protocol|extension)\s+(?P<name>\w+)`),
	"scala":        rules("function", `\bdef\s+(?P<name>\w+)`, `^\s*(?:[\w\s]*\s)?(?:class|object|trait)\s+(?P<name>\w+)`),
	"dart":         rules("function", `^\s*(?:class|mixin|enum|extension)\s+(?P<name>\w+)`, `^\s*(?:static\s+)?(?:Future<[^>]*>|[\w<>?]+)\s+(?P<name>\w+)\s*\([^;]*\)\s*(?:async\s*)?\{`),
	"d":            rules("function", `^\s*(?:class|struct|interface)\s+(?P<name>\w+)`, cLikeFunc),
	"julia":        rules("function", `^\s*function\s+(?P<name>[\w.!]+)`, `^\s*(?:mutable\s+)?struct\s+(?P<name>\w+)`, `^\s*module\s+(?P<name>\w+)`),
	"elixir":       rules("function", `^\s*defp?\s+(?P<name>[\w?!]+)`, `^\s*defmodule\s+(?P<name>[\w.]+)`),
	"erlang":       rules("function", `^(?P<name>[a-z]\w*)\s*\(.*\)\s*(?:when .*)?->`),
	"tcl":          rules("function", `^\s*proc\s+(?P<name>\S+)`),
	"r":            rules("function", `^\s*(?P<name>[\w.]+)\s*(?:<-|=)\s*function`),
	"fortran":      rules("function", `(?i)^\s*(?:[\w\s]*\s)?(?:subroutine|function|program|module)\s+(?P<name>\w+)`),
	"matlab":       rules("function", `^\s*function\s+(?:.*=\s*)?(?P<name>\w+)`),
	"nim":          rules("function", `^\s*(?:proc|func|method|iterator|macro|template)\s+(?P<name>\w+)`),
	"zig":          rules("function", `\bfn\s+(?P<name>\w+)`),
	"verilog":      rules("function", `^\s*(?:module|function|task|interface|package)\s+(?:automatic\s+)?(?P<name>\w+)`),
	"vhdl":         rules("function", `(?i)^\s*(?:entity|architecture|package|procedure|function)\s+(?P<name>\w+)`),
	"asm":          rules("label", `^(?P<name>[A-Za-z_.$][\w.$]*):`),
	"tex":          rules("heading", `\\(?P<level>part|chapter|section|subsection|subsubsection)\*?\{(?P<name>[^}]*)\}`),
	"groovy":       rules("function", `^\s*(?:[\w\s]*\s)?(?:class|interface|enum|trait)\s+(?P<name>\w+)`, `^\s*(?:(?:public|private|protected|static|def)\s+)+[\w<>\[\]]*\s*(?P<name>\w+)\s*\([^;]*$`),
	"autoit":       rules("function", `(?i)^\s*func\s+(?P<name>\w+)`),
	"nsis":         rules("function", `(?i)^\s*(?:function|section)\s+(?P<name>[^\s;]+)`),
	"protobuf":     rules("message", `^\s*(?:message|enum|service)\s+(?P<name>\w+)`, `^\s*rpc\s+(?P<name>\w+)`),
	"graphql":      rules("type", `^\s*(?:type|input|interface|enum|union|scalar|query|mutation|subscription|fragment)\s+(?P<name>\w+)`),
	"cmake":        rules("function", `(?i)^\s*(?:function|macro)\s*\(\s*(?P<name>[\w-]+)`),
	"fsharp":       rules("function", `^\s*(?:let|member)\s+(?:rec\s+|inline\s+|private\s+)*(?:\w+\.)?(?P<name>\w+)`, `^\s*type\s+(?P<name>\w+)`),
	"ocaml":        rules("function", `^\s*let\s+(?:rec\s+)?(?P<name>[a-z_]\w*)`, `^\s*(?:module|type)\s+(?P<name>\w+)`),
	"coffeescript": rules("function", `^\s*(?P<name>[\w$.]+)\s*[:=]\s*(?:\([^)]*\)\s*)?[-=]>`, `^\s*class\s+(?P<name>[\w$.]+)`),
}

// controlWords are not function names though the C-like rule matches them
var controlWords = wordSet("if for while switch return sizeof catch else do case new delete throw using typeof foreach lock when elif unless until")

// maxFuncListSize: larger documents get no function list
const maxFuncListSize = 32 << 20

// FunctionList returns the functions, classes, headings... of the document
func FunctionList(doc *Document) []FuncItem {
	lang := doc.Language()
	if lang == nil || doc.Len() > maxFuncListSize {
		return nil
	}
	rs := funcRules[lang.ID]
	if len(rs) == 0 {
		return nil
	}
	var out []FuncItem
	var ls LineStyles
	styled := lang.ID != "markdown" && lang.ID != "yaml" && lang.ID != "ini" && lang.ID != "css" && lang.ID != "tex"
	doc.Buffer().ForEachLine(0, func(line int, text []byte) bool {
		if len(text) == 0 || len(text) > 4096 {
			return true
		}
		s := string(text)
		for _, r := range rs {
			m := r.re.FindStringSubmatchIndex(s)
			if m == nil {
				continue
			}
			ni := r.re.SubexpIndex("name")
			if ni < 0 || m[2*ni] < 0 {
				continue
			}
			name := strings.TrimSpace(s[m[2*ni]:m[2*ni+1]])
			if name == "" || controlWords[name] {
				continue
			}
			if styled {
				doc.LineStyles(line, text, &ls)
				switch ls.StyleAt(m[2*ni]) {
				case StyleComment, StyleCommentLine, StyleCommentDoc, StyleString, StyleChar:
					continue
				}
			}
			level := 0
			if li := r.re.SubexpIndex("level"); li >= 0 && m[2*li] >= 0 {
				level = headingLevel(s[m[2*li]:m[2*li+1]])
			} else {
				level = len(s) - len(strings.TrimLeft(s, " \t"))
				if level > 0 {
					level = 1 + level/8
				}
			}
			out = append(out, FuncItem{Name: name, Line: line, Level: level, Kind: r.kind})
			break
		}
		return len(out) < 20000
	})
	return out
}

func headingLevel(s string) int {
	switch s {
	case "part":
		return 0
	case "chapter":
		return 0
	case "section":
		return 1
	case "subsection":
		return 2
	case "subsubsection":
		return 3
	}
	if strings.Trim(s, "#") == "" {
		return len(s) - 1
	}
	return 0
}
