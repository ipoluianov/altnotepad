package editor

import (
	"image/color"

	"github.com/ipoluianov/nui/ui"
)

// StyleColors are the colors of a style; a zero Back keeps the background
type StyleColors struct {
	Fore   color.RGBA
	Back   color.RGBA
	Bold   bool
	Italic bool
}

// Scheme is a color scheme of the editor (Settings > Style)
type Scheme struct {
	ID   string
	Name string
	Dark bool

	Styles [NumStyles]StyleColors

	Background  color.RGBA
	Foreground  color.RGBA
	CurrentLine color.RGBA
	Selection   color.RGBA
	// SelectionText, when set, is the text color of the selection
	SelectionText   color.RGBA
	Caret           color.RGBA
	GutterBack      color.RGBA
	LineNumber      color.RGBA
	LineNumberCur   color.RGBA
	FoldMarker      color.RGBA
	Whitespace      color.RGBA
	IndentGuide     color.RGBA
	Edge            color.RGBA
	BraceMatch      color.RGBA // background of the matching braces
	BraceBad        color.RGBA
	SmartHighlight  color.RGBA
	Marks           [NumMarkStyles]color.RGBA
	Bookmark        color.RGBA
	ControlCharBack color.RGBA
	FoldedLine      color.RGBA
	ScrollThumb     color.RGBA
	// The change history: the lines changed and not saved, and those saved
	ChangedLine color.RGBA
	SavedLine   color.RGBA
}

func hex(s string) color.RGBA { return ui.ColorFromHex(s) }

func withA(c color.RGBA, a uint8) color.RGBA {
	c.A = a
	return c
}

// schemeSpec is the short form of a scheme: the main colors and the syntax
type schemeSpec struct {
	id, name string
	dark     bool
	bg, fg   string
	curLine  string
	sel      string
	caret    string
	gutter   string
	lineNo   string

	comment, keyword, typ, builtin, str, number, operator, preproc, tag, attr, variable, constant, function string
}

func makeScheme(s schemeSpec) *Scheme {
	sc := &Scheme{ID: s.id, Name: s.name, Dark: s.dark}
	sc.Background = hex(s.bg)
	sc.Foreground = hex(s.fg)
	sc.CurrentLine = hex(s.curLine)
	sc.Selection = hex(s.sel)
	sc.Caret = hex(s.caret)
	sc.GutterBack = hex(s.gutter)
	sc.LineNumber = hex(s.lineNo)
	sc.LineNumberCur = sc.Foreground
	sc.FoldMarker = ui.MixColors(sc.LineNumber, sc.GutterBack, 0.2)
	sc.Whitespace = ui.MixColors(sc.Foreground, sc.Background, 0.72)
	sc.IndentGuide = ui.MixColors(sc.Foreground, sc.Background, 0.85)
	sc.Edge = ui.MixColors(sc.Foreground, sc.Background, 0.82)
	sc.ControlCharBack = ui.MixColors(sc.Foreground, sc.Background, 0.35)
	sc.FoldedLine = ui.MixColors(sc.Foreground, sc.Background, 0.6)
	sc.ScrollThumb = withA(ui.MixColors(sc.Foreground, sc.Background, 0.5), 110)
	sc.ChangedLine = hex("#f0a030")
	sc.SavedLine = hex("#3fae4f")
	if s.dark {
		sc.BraceMatch = hex("#3f6b3a")
		sc.BraceBad = hex("#8a2a2a")
		sc.SmartHighlight = withA(hex("#3a8f3a"), 110)
		sc.Marks = [NumMarkStyles]color.RGBA{withA(hex("#1e90ff"), 120), withA(hex("#ff8c00"), 120), withA(hex("#ffd700"), 110), withA(hex("#9932cc"), 130), withA(hex("#2e8b57"), 130), withA(hex("#d0d000"), 100)}
		sc.Bookmark = hex("#4d9ef5")
	} else {
		sc.BraceMatch = hex("#b5e3a6")
		sc.BraceBad = hex("#ff9a9a")
		sc.SmartHighlight = withA(hex("#00c800"), 70)
		sc.Marks = [NumMarkStyles]color.RGBA{withA(hex("#00ffff"), 110), withA(hex("#ffb000"), 110), withA(hex("#ffff00"), 130), withA(hex("#8000ff"), 70), withA(hex("#008000"), 70), withA(hex("#ff0000"), 60)}
		sc.Bookmark = hex("#1f6fd6")
	}
	st := &sc.Styles
	set := func(style Style, c string, bold, italic bool) {
		if c == "" {
			return
		}
		st[style] = StyleColors{Fore: hex(c), Bold: bold, Italic: italic}
	}
	for i := range st {
		st[i] = StyleColors{Fore: sc.Foreground}
	}
	set(StyleComment, s.comment, false, false)
	set(StyleCommentLine, s.comment, false, false)
	set(StyleCommentDoc, s.comment, false, false)
	set(StyleNumber, s.number, false, false)
	set(StyleKeyword, s.keyword, true, false)
	set(StyleType, s.typ, false, false)
	set(StyleBuiltin, s.builtin, false, false)
	set(StyleString, s.str, false, false)
	set(StyleChar, s.str, false, false)
	set(StyleOperator, s.operator, !s.dark, false)
	set(StylePreprocessor, s.preproc, false, false)
	set(StyleRegex, s.str, false, false)
	set(StyleTag, s.tag, false, false)
	set(StyleAttribute, s.attr, false, false)
	set(StyleEntity, s.constant, false, false)
	set(StyleCDATA, s.preproc, false, false)
	set(StyleHeading, s.keyword, true, false)
	set(StyleEmphasis, s.variable, false, true)
	set(StyleStrong, s.fg, true, false)
	set(StyleCode, s.str, false, false)
	set(StyleLink, s.tag, false, false)
	set(StyleQuote, s.comment, false, false)
	set(StyleVariable, s.variable, false, false)
	set(StyleAnnotation, s.preproc, false, false)
	set(StyleSection, s.keyword, true, false)
	set(StyleKey, s.attr, false, false)
	set(StyleValue, s.str, false, false)
	set(StyleLabel, s.constant, true, false)
	set(StyleConstant, s.constant, false, false)
	set(StyleFunction, s.function, false, false)
	set(StyleError, "#ff4040", false, false)
	if s.dark {
		st[StyleDiffAdded] = StyleColors{Fore: hex("#9fdf9f"), Back: hex("#1f3a1f")}
		st[StyleDiffRemoved] = StyleColors{Fore: hex("#f0a0a0"), Back: hex("#44201f")}
		st[StyleSearchHeader] = StyleColors{Fore: hex("#e0e0e0"), Back: hex("#2a3a52"), Bold: true}
		st[StyleSearchFile] = StyleColors{Fore: hex("#8fd18f"), Back: hex("#223022")}
	} else {
		st[StyleDiffAdded] = StyleColors{Fore: hex("#006000"), Back: hex("#ddffdd")}
		st[StyleDiffRemoved] = StyleColors{Fore: hex("#a00000"), Back: hex("#ffe0e0")}
		st[StyleSearchHeader] = StyleColors{Fore: hex("#000080"), Back: hex("#bbbbff"), Bold: true}
		st[StyleSearchFile] = StyleColors{Fore: hex("#008000"), Back: hex("#ffffbf")}
	}
	set(StyleDiffHeader, s.keyword, true, false)
	set(StyleDiffPosition, s.preproc, false, false)
	set(StyleSearchLineNo, s.number, false, false)
	return sc
}

// Schemes are the built-in color schemes
var Schemes = []*Scheme{
	makeScheme(schemeSpec{id: "default", name: "Default", bg: "#ffffff", fg: "#000000", curLine: "#e8e8ff", sel: "#c0c0c0", caret: "#000000", gutter: "#e4e4e4", lineNo: "#808080",
		comment: "#008000", keyword: "#0000ff", typ: "#8000ff", builtin: "#0080c0", str: "#808080", number: "#ff8000", operator: "#000080", preproc: "#804000",
		tag: "#0000ff", attr: "#ff0000", variable: "#a000a0", constant: "#8000ff", function: "#000000"}),
	makeScheme(schemeSpec{id: "obsidian", name: "Obsidian", dark: true, bg: "#293134", fg: "#e0e2e4", curLine: "#2f393c", sel: "#404e51", caret: "#e0e2e4", gutter: "#232829", lineNo: "#81969a",
		comment: "#66747b", keyword: "#93c763", typ: "#678cb1", builtin: "#a082bd", str: "#ec7600", number: "#ffcd22", operator: "#e8e2b7", preproc: "#a082bd",
		tag: "#8cbbad", attr: "#c3b77e", variable: "#d39745", constant: "#ffcd22", function: "#e0e2e4"}),
	makeScheme(schemeSpec{id: "dark", name: "Deep Dark", dark: true, bg: "#1e1e1e", fg: "#d4d4d4", curLine: "#282828", sel: "#264f78", caret: "#aeafad", gutter: "#1e1e1e", lineNo: "#858585",
		comment: "#6a9955", keyword: "#569cd6", typ: "#4ec9b0", builtin: "#dcdcaa", str: "#ce9178", number: "#b5cea8", operator: "#d4d4d4", preproc: "#c586c0",
		tag: "#569cd6", attr: "#9cdcfe", variable: "#9cdcfe", constant: "#569cd6", function: "#dcdcaa"}),
	makeScheme(schemeSpec{id: "monokai", name: "Monokai", dark: true, bg: "#272822", fg: "#f8f8f2", curLine: "#3e3d32", sel: "#49483e", caret: "#f8f8f0", gutter: "#2f3029", lineNo: "#90908a",
		comment: "#75715e", keyword: "#f92672", typ: "#66d9ef", builtin: "#a6e22e", str: "#e6db74", number: "#ae81ff", operator: "#f92672", preproc: "#f92672",
		tag: "#f92672", attr: "#a6e22e", variable: "#fd971f", constant: "#ae81ff", function: "#a6e22e"}),
	makeScheme(schemeSpec{id: "dracula", name: "Dracula", dark: true, bg: "#282a36", fg: "#f8f8f2", curLine: "#343746", sel: "#44475a", caret: "#f8f8f0", gutter: "#21222c", lineNo: "#6272a4",
		comment: "#6272a4", keyword: "#ff79c6", typ: "#8be9fd", builtin: "#50fa7b", str: "#f1fa8c", number: "#bd93f9", operator: "#ff79c6", preproc: "#ff79c6",
		tag: "#ff79c6", attr: "#50fa7b", variable: "#ffb86c", constant: "#bd93f9", function: "#50fa7b"}),
	makeScheme(schemeSpec{id: "solarized-dark", name: "Solarized Dark", dark: true, bg: "#002b36", fg: "#93a1a1", curLine: "#073642", sel: "#274642", caret: "#93a1a1", gutter: "#073642", lineNo: "#586e75",
		comment: "#586e75", keyword: "#859900", typ: "#b58900", builtin: "#268bd2", str: "#2aa198", number: "#d33682", operator: "#93a1a1", preproc: "#cb4b16",
		tag: "#268bd2", attr: "#93a1a1", variable: "#268bd2", constant: "#cb4b16", function: "#268bd2"}),
	makeScheme(schemeSpec{id: "solarized-light", name: "Solarized Light", bg: "#fdf6e3", fg: "#586e75", curLine: "#eee8d5", sel: "#e6dfc8", caret: "#586e75", gutter: "#eee8d5", lineNo: "#93a1a1",
		comment: "#93a1a1", keyword: "#859900", typ: "#b58900", builtin: "#268bd2", str: "#2aa198", number: "#d33682", operator: "#586e75", preproc: "#cb4b16",
		tag: "#268bd2", attr: "#586e75", variable: "#268bd2", constant: "#cb4b16", function: "#268bd2"}),
	makeScheme(schemeSpec{id: "zenburn", name: "Zenburn", dark: true, bg: "#3f3f3f", fg: "#dcdccc", curLine: "#474747", sel: "#5f5f5f", caret: "#ffffef", gutter: "#383838", lineNo: "#9fafaf",
		comment: "#7f9f7f", keyword: "#f0dfaf", typ: "#dfdfbf", builtin: "#efef8f", str: "#cc9393", number: "#8cd0d3", operator: "#f0efd0", preproc: "#ffcfaf",
		tag: "#e3ceab", attr: "#dfdfbf", variable: "#dcdccc", constant: "#dca3a3", function: "#efef8f"}),
	makeScheme(schemeSpec{id: "nord", name: "Nord", dark: true, bg: "#2e3440", fg: "#d8dee9", curLine: "#3b4252", sel: "#434c5e", caret: "#d8dee9", gutter: "#2e3440", lineNo: "#4c566a",
		comment: "#616e88", keyword: "#81a1c1", typ: "#8fbcbb", builtin: "#88c0d0", str: "#a3be8c", number: "#b48ead", operator: "#81a1c1", preproc: "#5e81ac",
		tag: "#81a1c1", attr: "#8fbcbb", variable: "#d8dee9", constant: "#b48ead", function: "#88c0d0"}),
	makeScheme(schemeSpec{id: "github", name: "GitHub Light", bg: "#ffffff", fg: "#24292f", curLine: "#f6f8fa", sel: "#c8e1ff", caret: "#24292f", gutter: "#ffffff", lineNo: "#8c959f",
		comment: "#6e7781", keyword: "#cf222e", typ: "#953800", builtin: "#8250df", str: "#0a3069", number: "#0550ae", operator: "#24292f", preproc: "#cf222e",
		tag: "#116329", attr: "#0550ae", variable: "#953800", constant: "#0550ae", function: "#8250df"}),
	makeScheme(schemeSpec{id: "onedark", name: "One Dark", dark: true, bg: "#282c34", fg: "#abb2bf", curLine: "#2c313c", sel: "#3e4451", caret: "#528bff", gutter: "#282c34", lineNo: "#636d83",
		comment: "#5c6370", keyword: "#c678dd", typ: "#e5c07b", builtin: "#56b6c2", str: "#98c379", number: "#d19a66", operator: "#56b6c2", preproc: "#c678dd",
		tag: "#e06c75", attr: "#d19a66", variable: "#e06c75", constant: "#d19a66", function: "#61afef"}),
}

// SchemeByID returns the scheme; the default one for an unknown ID
func SchemeByID(id string) *Scheme {
	for _, s := range Schemes {
		if s.ID == id {
			return s
		}
	}
	return Schemes[0]
}
