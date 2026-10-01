package editor

// Style is the kind of a piece of text for the syntax highlighting; the
// colors of the styles come from the color scheme (see Scheme)
type Style uint8

const (
	StyleDefault Style = iota
	StyleComment
	StyleCommentLine
	StyleCommentDoc
	StyleNumber
	StyleKeyword
	StyleType
	StyleBuiltin
	StyleString
	StyleChar
	StyleOperator
	StylePreprocessor
	StyleRegex
	StyleTag
	StyleAttribute
	StyleEntity
	StyleCDATA
	StyleHeading
	StyleEmphasis
	StyleStrong
	StyleCode
	StyleLink
	StyleQuote
	StyleVariable
	StyleAnnotation
	StyleSection
	StyleKey
	StyleValue
	StyleDiffAdded
	StyleDiffRemoved
	StyleDiffHeader
	StyleDiffPosition
	StyleError
	StyleSearchHeader
	StyleSearchFile
	StyleSearchLineNo
	StyleLabel
	StyleConstant
	StyleFunction

	NumStyles
)

// StyleNames are the names of the styles as used in color scheme files
var StyleNames = [NumStyles]string{
	"default", "comment", "commentline", "commentdoc", "number", "keyword", "type", "builtin",
	"string", "char", "operator", "preprocessor", "regex", "tag", "attribute", "entity", "cdata",
	"heading", "emphasis", "strong", "code", "link", "quote", "variable", "annotation", "section",
	"key", "value", "diffadded", "diffremoved", "diffheader", "diffposition", "error",
	"searchheader", "searchfile", "searchlineno", "label", "constant", "function",
}

// NumMarkStyles is the number of styles of Search > Mark (1-5) plus the one of Find > Mark All
const NumMarkStyles = 6

// MarkStyleFind is the mark style used by Mark All of the search
const MarkStyleFind = 5
