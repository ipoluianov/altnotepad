package editor

import (
	"path/filepath"
	"sort"
	"strings"
)

// Language is a language of the Language menu: how its files are named,
// highlighted, commented and folded
type Language struct {
	ID   string
	Name string

	Extensions []string // lower case, without the dot
	FileNames  []string // whole names, lower case: "makefile"
	Shebangs   []string // interpreters of #! lines: "python", "bash"

	LineComment  string
	BlockComment [2]string

	// FoldIndent: blocks fold by indentation, as in Python
	FoldIndent bool
	// NoFold: the language has no folding
	NoFold bool
	// Hidden: not offered in the menu, e.g. the search results
	Hidden bool

	lexer Lexer
	words []string
}

// Words returns the keywords of the language, for the auto-completion
func (l *Language) Words() []string {
	return l.words
}

var (
	languages     []*Language
	languagesByID = map[string]*Language{}
)

// LangText is the plain text
var LangText = &Language{ID: "text", Name: "Normal text", Extensions: []string{"txt", "log"}, NoFold: true, lexer: plainLexer{}}

// Languages returns the languages of the menu, sorted by name
func Languages() []*Language {
	out := make([]*Language, 0, len(languages))
	for _, l := range languages {
		if !l.Hidden {
			out = append(out, l)
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out
}

// LanguageByID returns the language; the plain text for an unknown ID
func LanguageByID(id string) *Language {
	if l, ok := languagesByID[id]; ok {
		return l
	}
	return LangText
}

// LanguageForFile picks the language by the name of the file and, when
// that tells nothing, by its first line (#!/usr/bin/env python, <?xml)
func LanguageForFile(path string, firstLine []byte) *Language {
	name := strings.ToLower(filepath.Base(path))
	for _, l := range languages {
		for _, fn := range l.FileNames {
			if name == fn {
				return l
			}
		}
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
	if ext != "" {
		for _, l := range languages {
			for _, e := range l.Extensions {
				if e == ext {
					return l
				}
			}
		}
	}
	first := strings.TrimSpace(string(firstLine))
	if strings.HasPrefix(first, "#!") {
		fields := strings.Fields(strings.TrimPrefix(first, "#!"))
		interp := ""
		if len(fields) > 0 {
			interp = filepath.Base(fields[0])
			if interp == "env" && len(fields) > 1 {
				interp = fields[1]
			}
		}
		for _, l := range languages {
			for _, s := range l.Shebangs {
				if strings.HasPrefix(interp, s) {
					return l
				}
			}
		}
	}
	if strings.HasPrefix(first, "<?xml") {
		return LanguageByID("xml")
	}
	if strings.HasPrefix(strings.ToLower(first), "<!doctype html") || strings.HasPrefix(strings.ToLower(first), "<html") {
		return LanguageByID("html")
	}
	return LangText
}

func addLanguage(l *Language, def *genericDef) *Language {
	if def != nil {
		l.lexer = newGeneric(def)
		for w := range def.words {
			l.words = append(l.words, w)
		}
		sort.Strings(l.words)
	}
	languages = append(languages, l)
	languagesByID[l.ID] = l
	return l
}

// kw makes the styles of the words: keywords, types, built-ins and constants
func kw(keywords, types, builtins, constants string) map[string]Style {
	return wordStyles(StyleKeyword, keywords, StyleType, types, StyleBuiltin, builtins, StyleConstant, constants)
}

var (
	strDQ     = stringDef{open: `"`, close: `"`, escape: '\\', style: StyleString}
	strSQ     = stringDef{open: `'`, close: `'`, escape: '\\', style: StyleString}
	charSQ    = stringDef{open: `'`, close: `'`, escape: '\\', style: StyleChar}
	strBT     = stringDef{open: "`", close: "`", escape: '\\', multiline: true, style: StyleString}
	strRawBT  = stringDef{open: "`", close: "`", multiline: true, style: StyleString}
	strTDQ    = stringDef{open: `"""`, close: `"""`, escape: '\\', multiline: true, style: StyleString}
	strTSQ    = stringDef{open: `'''`, close: `'''`, escape: '\\', multiline: true, style: StyleString}
	cComment  = blockDef{open: "/*", close: "*/", style: StyleComment}
	docCmt    = blockDef{open: "/**", close: "*/", style: StyleCommentDoc, notFollowedBy: '/'}
	sqlString = stringDef{open: `'`, close: `'`, style: StyleString, doubledClose: true}
)

// cLike returns the definition of a language with C comments and strings
func cLike(words map[string]Style) *genericDef {
	return &genericDef{
		lineComments:  []string{"//"},
		blockComments: []blockDef{docCmt, cComment},
		strings:       []stringDef{strDQ, charSQ},
		words:         words,
		foldBraces:    "{}",
	}
}

const (
	cKeywords   = "auto break case const continue default do else enum extern for goto if inline register restrict return sizeof static struct switch typedef union volatile while _Alignas _Alignof _Atomic _Generic _Noreturn _Static_assert _Thread_local"
	cTypes      = "char double float int long short signed unsigned void bool _Bool _Complex size_t ssize_t int8_t int16_t int32_t int64_t uint8_t uint16_t uint32_t uint64_t intptr_t uintptr_t ptrdiff_t wchar_t FILE va_list"
	cppKeywords = cKeywords + " alignas alignof and and_eq asm bitand bitor catch class compl concept consteval constexpr constinit const_cast co_await co_return co_yield decltype delete dynamic_cast explicit export friend mutable namespace new noexcept not not_eq operator or or_eq private protected public reinterpret_cast requires static_assert static_cast template this thread_local throw try typeid typename using virtual xor xor_eq override final"
	cppTypes    = cTypes + " char8_t char16_t char32_t nullptr_t string wstring vector map set unordered_map unordered_set shared_ptr unique_ptr weak_ptr"
	jsKeywords  = "break case catch class const continue debugger default delete do else export extends finally for function if import in instanceof let new return super switch this throw try typeof var void while with yield async await of static get set from as"
	jsBuiltins  = "console window document Math JSON Object Array String Number Boolean Promise Map Set WeakMap WeakSet Symbol Date RegExp Error parseInt parseFloat isNaN require module exports process globalThis setTimeout setInterval clearTimeout clearInterval fetch"
	jsConstants = "true false null undefined NaN Infinity"
	tsKeywords  = jsKeywords + " abstract declare enum implements interface keyof namespace private protected public readonly type is infer module override satisfies"
	tsTypes     = "any boolean number string symbol unknown never object void bigint"
)

func init() {
	languages = append(languages, LangText)
	languagesByID[LangText.ID] = LangText

	// C and its family
	addLanguage(&Language{ID: "c", Name: "C", Extensions: []string{"c", "h"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		withPre(cLike(kw(cKeywords, cTypes, "printf malloc free memcpy memset strlen strcpy strcmp sizeof", "NULL true false EOF"))))
	addLanguage(&Language{ID: "cpp", Name: "C++", Extensions: []string{"cpp", "cxx", "cc", "c++", "hpp", "hxx", "hh", "h++", "ino", "inl", "ipp", "tcc"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		withPre(cLike(kw(cppKeywords, cppTypes, "std cout cin cerr endl printf malloc free", "nullptr NULL true false"))))
	addLanguage(&Language{ID: "objc", Name: "Objective-C", Extensions: []string{"m", "mm"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		withAnnotation(withPre(cLike(kw(cppKeywords+" self super in out inout bycopy byref oneway", cTypes+" id Class SEL IMP BOOL NSString NSObject NSInteger NSUInteger", "NSLog", "nil Nil YES NO NULL true false")))))
	addLanguage(&Language{ID: "cs", Name: "C#", Extensions: []string{"cs", "csx"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withPre(cLike(kw("abstract as base break case catch checked class const continue default delegate do else enum event explicit extern finally fixed for foreach goto if implicit in interface internal is lock namespace new operator out override params private protected public readonly ref return sealed sizeof stackalloc static struct switch this throw try typeof unchecked unsafe using virtual volatile while add alias ascending async await by descending dynamic equals from get global group into join let nameof on orderby partial record remove select set value var when where yield init required",
				"bool byte char decimal double float int long object sbyte short string uint ulong ushort void nint nuint", "Console Math String List Dictionary Task", "true false null")))
			d.docLineComments = []string{"///"}
			d.strings = []stringDef{{open: `@"`, close: `"`, multiline: true, style: StyleString, doubledClose: true}, {open: `$"`, close: `"`, escape: '\\', style: StyleString}, strDQ, charSQ}
			return d
		}())
	addLanguage(&Language{ID: "java", Name: "Java", Extensions: []string{"java", "jav"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("abstract assert break case catch class const continue default do else enum extends final finally for goto if implements import instanceof interface native new package private protected public return static strictfp super switch synchronized this throw throws transient try volatile while var record yield sealed permits",
				"boolean byte char double float int long short void String Object Integer Long Double Boolean Character List Map Set", "System Math", "true false null")))
			d.strings = []stringDef{strTDQ, strDQ, charSQ}
			return d
		}())
	addLanguage(&Language{ID: "javascript", Name: "JavaScript", Extensions: []string{"js", "mjs", "cjs", "jsx", "jsm"}, Shebangs: []string{"node"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		jsDef(kw(jsKeywords, "", jsBuiltins, jsConstants)))
	addLanguage(&Language{ID: "typescript", Name: "TypeScript", Extensions: []string{"ts", "tsx", "mts", "cts"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		withAnnotation(jsDef(kw(tsKeywords, tsTypes, jsBuiltins, jsConstants))))
	addLanguage(&Language{ID: "json", Name: "JSON", Extensions: []string{"json", "jsonc", "json5", "geojson", "webmanifest", "babelrc", "eslintrc"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		&genericDef{lineComments: []string{"//"}, blockComments: []blockDef{cComment}, strings: []stringDef{strDQ, strSQ},
			words: kw("", "", "", "true false null"), foldBraces: "{}[]", keyBeforeColon: true})
	addLanguage(&Language{ID: "go", Name: "Go", Extensions: []string{"go"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := cLike(kw("break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var",
				"bool byte complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr any comparable",
				"append cap clear close complex copy delete imag len make max min new panic print println real recover", "true false nil iota"))
			d.blockComments = []blockDef{cComment}
			d.strings = []stringDef{strDQ, strRawBT, charSQ}
			return d
		}())
	addLanguage(&Language{ID: "rust", Name: "Rust", Extensions: []string{"rs"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := cLike(kw("as async await break const continue crate dyn else enum extern fn for if impl in let loop match mod move mut pub ref return self Self static struct super trait type unsafe use where while macro_rules union",
				"bool char str u8 u16 u32 u64 u128 usize i8 i16 i32 i64 i128 isize f32 f64 String Vec Option Result Box Rc Arc HashMap HashSet",
				"println print format vec panic assert assert_eq debug_assert write writeln", "true false None Some Ok Err"))
			d.blockComments = []blockDef{{open: "/*", close: "*/", style: StyleComment, nested: true}}
			d.docLineComments = []string{"///", "//!"}
			d.strings = []stringDef{strDQ}
			d.annotation = '#'
			return d
		}())
	addLanguage(&Language{ID: "swift", Name: "Swift", Extensions: []string{"swift"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("associatedtype class deinit enum extension fileprivate func import init inout internal let open operator private protocol public rethrows static struct subscript typealias var break case continue default defer do else fallthrough for guard if in repeat return switch where while as catch is super self Self throw throws try async await some any",
				"Int Int8 Int16 Int32 Int64 UInt Double Float String Bool Character Array Dictionary Set Optional AnyObject Void", "print", "true false nil")))
			d.blockComments = []blockDef{{open: "/*", close: "*/", style: StyleComment, nested: true}}
			d.strings = []stringDef{strTDQ, strDQ}
			return d
		}())
	addLanguage(&Language{ID: "kotlin", Name: "Kotlin", Extensions: []string{"kt", "kts"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("package import class interface fun val var if else when for while do return break continue object companion data sealed open override private public protected internal abstract final enum try catch finally throw is as in out by constructor init this super typealias suspend inline reified lateinit where",
				"Int Long Short Byte Double Float Boolean Char String Unit Any Nothing List Map Set Array MutableList MutableMap", "println print listOf mapOf setOf arrayOf", "true false null")))
			d.strings = []stringDef{strTDQ, strDQ, charSQ}
			return d
		}())
	addLanguage(&Language{ID: "scala", Name: "Scala", Extensions: []string{"scala", "sc", "sbt"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("abstract case catch class def do else extends final finally for forSome if implicit import lazy match new object override package private protected return sealed super this throw trait try type val var while with yield given using enum then export",
				"Int Long Short Byte Double Float Boolean Char String Unit Any AnyRef Nothing List Map Set Seq Option", "println", "true false null None Some")))
			d.strings = []stringDef{strTDQ, strDQ, charSQ}
			return d
		}())
	addLanguage(&Language{ID: "groovy", Name: "Groovy", Extensions: []string{"groovy", "gradle", "gvy"}, Shebangs: []string{"groovy"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("abstract as assert break case catch class const continue def default do else enum extends final finally for goto if implements import in instanceof interface new package private protected public return static super switch synchronized this throw throws trait try var while",
				"boolean byte char double float int long short void String Object", "println", "true false null")))
			d.strings = []stringDef{strTDQ, strTSQ, strDQ, strSQ}
			return d
		}())
	addLanguage(&Language{ID: "dart", Name: "Dart", Extensions: []string{"dart"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("abstract as assert async await break case catch class const continue covariant default deferred do dynamic else enum export extends extension external factory final finally for get hide if implements import in interface is late library mixin new on operator part required rethrow return set show static super switch sync this throw try typedef var void while with yield",
				"int double num String bool List Map Set Future Stream Object Iterable", "print", "true false null")))
			d.docLineComments = []string{"///"}
			d.strings = []stringDef{strTDQ, strTSQ, strDQ, strSQ}
			return d
		}())
	addLanguage(&Language{ID: "d", Name: "D", Extensions: []string{"d", "di"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := withAnnotation(cLike(kw("abstract alias align asm assert auto body break case cast catch class const continue debug default delegate delete deprecated do else enum export extern final finally for foreach foreach_reverse function goto if immutable import in inout interface invariant is lazy mixin module new nothrow out override package pragma private protected public pure ref return scope shared static struct super switch synchronized template this throw try typeid typeof union unittest version while with",
				"bool byte ubyte short ushort int uint long ulong char wchar dchar float double real void string size_t", "writeln writefln", "true false null")))
			d.blockComments = append(d.blockComments, blockDef{open: "/+", close: "+/", style: StyleComment, nested: true})
			d.strings = []stringDef{strDQ, strRawBT, charSQ}
			return d
		}())
	addLanguage(&Language{ID: "zig", Name: "Zig", Extensions: []string{"zig"}, LineComment: "//"},
		func() *genericDef {
			d := cLike(kw("addrspace align allowzero and anyframe anytype asm async await break callconv catch comptime const continue defer else enum errdefer error export extern fn for if inline linksection noalias noinline nosuspend opaque or orelse packed pub resume return struct suspend switch test threadlocal try union unreachable usingnamespace var volatile while",
				"i8 u8 i16 u16 i32 u32 i64 u64 i128 u128 isize usize f16 f32 f64 f128 bool void noreturn type anyerror comptime_int comptime_float", "", "true false null undefined"))
			d.blockComments = nil
			d.docLineComments = []string{"///", "//!"}
			d.annotation = '@'
			return d
		}())

	// Web
	jsLexer := newGeneric(jsDef(kw(jsKeywords, "", jsBuiltins, jsConstants)))
	phpDef := &genericDef{lineComments: []string{"//", "#"}, blockComments: []blockDef{docCmt, cComment},
		strings: []stringDef{strDQ, strSQ, strBT}, varPrefixes: "$", foldBraces: "{}",
		words: kw("abstract and array as break callable case catch class clone const continue declare default do echo else elseif empty enddeclare endfor endforeach endif endswitch endwhile eval exit extends final finally fn for foreach function global goto if implements include include_once instanceof insteadof interface isset list match namespace new or print private protected public readonly require require_once return static switch throw trait try unset use var while xor yield enum",
			"int float bool string void mixed object iterable never self parent", "strlen count isset print_r var_dump die", "true false null TRUE FALSE NULL"),
		caseInsensitive: false}
	phpLexer := newGeneric(phpDef)
	addLanguage(&Language{ID: "html", Name: "HTML", Extensions: []string{"html", "htm", "shtml", "shtm", "xhtml", "xht", "hta", "vue", "svelte"}, BlockComment: [2]string{"<!--", "-->"}},
		nil).lexer = &markupLexer{html: true, script: jsLexer, style: cssLexer{}}
	addLanguage(&Language{ID: "xml", Name: "XML", Extensions: []string{"xml", "xsl", "xslt", "xsd", "xaml", "svg", "rss", "atom", "plist", "csproj", "vcxproj", "props", "targets", "config", "resx", "wsdl", "kml", "gpx", "xul", "fxml", "iml", "nuspec", "storyboard", "xib", "tld", "dtd", "ui"}, BlockComment: [2]string{"<!--", "-->"}},
		nil).lexer = &markupLexer{}
	addLanguage(&Language{ID: "php", Name: "PHP", Extensions: []string{"php", "php3", "php4", "php5", "phtml", "phps", "inc"}, Shebangs: []string{"php"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		nil).lexer = &markupLexer{html: true, script: jsLexer, style: cssLexer{}, php: phpLexer}
	languagesByID["php"].words = sortedWords(phpDef.words)
	addLanguage(&Language{ID: "css", Name: "CSS", Extensions: []string{"css", "scss", "less", "sass"}, BlockComment: [2]string{"/*", "*/"}}, nil).lexer = cssLexer{}
	addLanguage(&Language{ID: "markdown", Name: "Markdown", Extensions: []string{"md", "markdown", "mdown", "mkd", "mkdn", "mdwn", "rmd"}, BlockComment: [2]string{"<!--", "-->"}},
		nil).lexer = markdownLexer{}

	// Scripting
	addLanguage(&Language{ID: "python", Name: "Python", Extensions: []string{"py", "pyw", "pyi", "pyx", "pxd", "sconstruct", "sconscript", "gyp", "gypi", "bzl", "wscript"}, FileNames: []string{"sconstruct", "sconscript", "build", "workspace"}, Shebangs: []string{"python"}, LineComment: "#", FoldIndent: true},
		&genericDef{lineComments: []string{"#"}, strings: []stringDef{strTDQ, strTSQ, strDQ, strSQ}, stringPrefixes: "rRbBuUfF", annotation: '@',
			words: kw("and as assert async await break class continue def del elif else except finally for from global if import in is lambda nonlocal not or pass raise return try while with yield match case",
				"int float str bytes bool list dict set tuple object complex frozenset bytearray", "print len range open type isinstance issubclass super enumerate zip map filter sorted reversed min max sum abs any all iter next hasattr getattr setattr delattr repr format input round divmod id hash callable vars dir globals locals self cls __init__ __name__ __main__",
				"True False None NotImplemented Ellipsis")})
	addLanguage(&Language{ID: "ruby", Name: "Ruby", Extensions: []string{"rb", "rbw", "rake", "gemspec", "ru", "podspec"}, FileNames: []string{"rakefile", "gemfile", "podfile", "vagrantfile", "guardfile"}, Shebangs: []string{"ruby"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, lineCommentsAtStart: nil, blockComments: []blockDef{{open: "=begin", close: "=end", style: StyleComment}},
			strings: []stringDef{strDQ, strSQ, strBT}, varPrefixes: "@$", regexLiterals: true,
			words: kw("BEGIN END alias and begin break case class def defined? do else elsif end ensure for if in module next not or redo rescue retry return self super then undef unless until when while yield",
				"", "require require_relative puts print p attr_accessor attr_reader attr_writer include extend raise lambda proc loop", "true false nil"),
			foldOpen: wordSet("def class module do begin case while until unless"), foldClose: wordSet("end"), foldNoOpenAfter: wordSet("end")})
	addLanguage(&Language{ID: "perl", Name: "Perl", Extensions: []string{"pl", "pm", "plx", "t", "pod"}, Shebangs: []string{"perl"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, blockComments: []blockDef{{open: "=pod", close: "=cut", style: StyleComment}, {open: "=head", close: "=cut", style: StyleComment}},
			strings: []stringDef{strDQ, strSQ, strBT}, varPrefixes: "$@", regexLiterals: true, foldBraces: "{}",
			words: kw("my our local sub if elsif else unless while until for foreach do last next redo return use require package BEGIN END no and or not eq ne lt gt le ge cmp",
				"", "print printf die warn eval defined undef shift push pop split join keys values exists delete open close chomp scalar ref bless", "")})
	addLanguage(&Language{ID: "lua", Name: "Lua", Extensions: []string{"lua", "wlua", "rockspec"}, Shebangs: []string{"lua"}, LineComment: "--", BlockComment: [2]string{"--[[", "]]"}},
		&genericDef{lineComments: []string{"--"}, blockComments: []blockDef{{open: "--[[", close: "]]", style: StyleComment}},
			strings: []stringDef{{open: "[[", close: "]]", multiline: true, style: StyleString}, strDQ, strSQ}, foldBraces: "{}",
			words: kw("and break do else elseif end for function goto if in local not or repeat return then until while",
				"", "print pairs ipairs type tostring tonumber require string table math io os error pcall xpcall select setmetatable getmetatable rawget rawset assert unpack next", "true false nil"),
			foldOpen: wordSet("function do if repeat"), foldClose: wordSet("end until"), foldMid: wordSet("else elseif")})
	addLanguage(&Language{ID: "bash", Name: "Shell", Extensions: []string{"sh", "bash", "zsh", "ksh", "csh", "fish", "bashrc", "zshrc", "profile", "bash_profile", "command", "ebuild", "eclass"}, FileNames: []string{".bashrc", ".zshrc", ".profile", ".bash_profile", ".bash_aliases", "pkgbuild", "apkbuild"}, Shebangs: []string{"sh", "bash", "zsh", "ksh", "dash", "ash", "fish"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, lineCommentAfterSpace: true,
			strings:     []stringDef{{open: `"`, close: `"`, escape: '\\', multiline: true, style: StyleString}, {open: `'`, close: `'`, multiline: true, style: StyleString}, strBT},
			varPrefixes: "$", foldBraces: "{}", identExtra: "-",
			words: kw("if then else elif fi case esac for select while until do done in function time coproc return exit break continue local export declare readonly typeset unset shift source alias eval exec set trap let",
				"", "echo printf read cd pwd test true false cat grep sed awk ls rm cp mv mkdir chmod chown find xargs sudo", ""),
			foldOpen: wordSet("if do case"), foldClose: wordSet("fi done esac"), foldMid: wordSet("else elif")})
	addLanguage(&Language{ID: "powershell", Name: "PowerShell", Extensions: []string{"ps1", "psm1", "psd1", "ps1xml", "pssc", "psrc", "cdxml"}, Shebangs: []string{"pwsh", "powershell"}, LineComment: "#", BlockComment: [2]string{"<#", "#>"}},
		&genericDef{lineComments: []string{"#"}, blockComments: []blockDef{{open: "<#", close: "#>", style: StyleComment}},
			strings:     []stringDef{{open: `@"`, close: `"@`, multiline: true, style: StyleString}, {open: `@'`, close: `'@`, multiline: true, style: StyleString}, {open: `"`, close: `"`, escape: '`', style: StyleString}, {open: `'`, close: `'`, style: StyleString, doubledClose: true}},
			varPrefixes: "$", foldBraces: "{}", caseInsensitive: true, identExtra: "-",
			words: kw("begin break catch class continue data define do dynamicparam else elseif end enum exit filter finally for foreach from function hidden if in param process return static switch throw trap try until using var while workflow -eq -ne -gt -ge -lt -le -like -notlike -match -notmatch -contains -notcontains -in -notin -and -or -not -xor -is -isnot -replace",
				"", "Write-Host Write-Output Get-Item Set-Item Get-ChildItem Get-Content Set-Content Where-Object ForEach-Object Select-Object", "$true $false $null")})
	addLanguage(&Language{ID: "batch", Name: "Batch", Extensions: []string{"bat", "cmd", "nt"}, LineComment: "REM "},
		&genericDef{lineCommentsAtStart: []string{"rem", "::", "@rem"}, labelPrefix: ':', varPrefixes: "%", caseInsensitive: true,
			strings: []stringDef{{open: `"`, close: `"`, style: StyleString}}, foldBraces: "()",
			words: kw("if else goto call exit set setlocal endlocal for in do echo not exist defined errorlevel shift pause cls cd chdir copy xcopy del erase dir move ren rename md rd mkdir rmdir start title type choice pushd popd equ neq lss leq gtr geq nul con enabledelayedexpansion enableextensions verify ver vol path prompt color mode break",
				"", "", "on off")})
	addLanguage(&Language{ID: "tcl", Name: "TCL", Extensions: []string{"tcl", "tk", "itcl", "exp"}, Shebangs: []string{"tclsh", "wish", "expect"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, strings: []stringDef{strDQ}, varPrefixes: "$", foldBraces: "{}",
			words: kw("after append array break catch cd close concat continue else elseif eof error eval exec exit expr for foreach format gets glob global if incr info join lappend lindex linsert list llength lrange lreplace lsearch lsort namespace open package proc puts read regexp regsub rename return scan seek set source split string switch then trace unset update uplevel upvar variable vwait while", "", "", "")})
	addLanguage(&Language{ID: "r", Name: "R", Extensions: []string{"r", "rprofile"}, Shebangs: []string{"rscript"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, strings: []stringDef{strDQ, strSQ, strRawBT}, foldBraces: "{}", identExtra: ".",
			words: kw("if else repeat while function for in next break return switch", "", "print cat paste library require c list data.frame matrix vector length sum mean", "TRUE FALSE NULL Inf NaN NA NA_integer_ NA_real_ NA_character_ T F")})
	addLanguage(&Language{ID: "julia", Name: "Julia", Extensions: []string{"jl"}, Shebangs: []string{"julia"}, LineComment: "#", BlockComment: [2]string{"#=", "=#"}},
		&genericDef{lineComments: []string{"#"}, blockComments: []blockDef{{open: "#=", close: "=#", style: StyleComment, nested: true}}, strings: []stringDef{strTDQ, strDQ, strRawBT}, annotation: '@',
			words: kw("abstract baremodule begin break catch const continue do else elseif end export finally for function global if import in let local macro module mutable primitive quote return struct try type using where while",
				"Int Int8 Int16 Int32 Int64 UInt8 Float32 Float64 Bool Char String Vector Matrix Array Dict Any Nothing", "println print", "true false nothing missing"),
			foldOpen: wordSet("function if for while begin let module baremodule struct macro quote do try"), foldClose: wordSet("end"), foldNoOpenAfter: wordSet("end")})
	addLanguage(&Language{ID: "coffeescript", Name: "CoffeeScript", Extensions: []string{"coffee", "litcoffee", "cson"}, LineComment: "#", BlockComment: [2]string{"###", "###"}, FoldIndent: true},
		&genericDef{lineComments: []string{"#"}, blockComments: []blockDef{{open: "###", close: "###", style: StyleComment}}, strings: []stringDef{strTDQ, strTSQ, strDQ, strSQ}, regexLiterals: true, annotation: '@',
			words: kw("and break by catch class continue debugger delete do else extends finally for if in instanceof is isnt loop new not of or return super switch then this throw try typeof unless until when while yield", "", "console require", "true false null undefined yes no on off")})

	// Functional
	addLanguage(&Language{ID: "haskell", Name: "Haskell", Extensions: []string{"hs", "lhs"}, LineComment: "--", BlockComment: [2]string{"{-", "-}"}, FoldIndent: true},
		&genericDef{lineComments: []string{"--"}, blockComments: []blockDef{{open: "{-", close: "-}", style: StyleComment, nested: true}}, strings: []stringDef{strDQ, charSQ}, identExtra: "'",
			words: kw("case class data default deriving do else foreign if import in infix infixl infixr instance let module newtype of then type where qualified as hiding forall", "Int Integer Float Double Char String Bool Maybe Either IO", "putStrLn print show map filter foldr foldl", "True False Nothing Just Left Right")})
	addLanguage(&Language{ID: "fsharp", Name: "F#", Extensions: []string{"fs", "fsi", "fsx", "fsscript"}, LineComment: "//", BlockComment: [2]string{"(*", "*)"}, FoldIndent: true},
		&genericDef{lineComments: []string{"//"}, blockComments: []blockDef{{open: "(*", close: "*)", style: StyleComment, nested: true, notFollowedBy: ')'}}, strings: []stringDef{strTDQ, strDQ}, annotation: 0,
			words: kw("abstract and as assert base begin class default delegate do done downcast downto elif else end exception extern false finally fixed for fun function global if in inherit inline interface internal lazy let match member module mutable namespace new not null of open or override private public rec return select sig static struct then to true try type upcast use val void when while with yield async task",
				"int float string bool char unit list option seq array decimal byte", "printfn printf sprintf failwith", "true false null None Some")})
	addLanguage(&Language{ID: "ocaml", Name: "OCaml", Extensions: []string{"ml", "mli", "mll", "mly"}, BlockComment: [2]string{"(*", "*)"}},
		&genericDef{blockComments: []blockDef{{open: "(*", close: "*)", style: StyleComment, nested: true, notFollowedBy: ')'}}, strings: []stringDef{strDQ},
			words:    kw("and as assert begin class constraint do done downto else end exception external for fun function functor if in include inherit initializer lazy let match method module mutable new object of open or private rec sig struct then to try type val virtual when while with", "int float string bool char unit list option array", "print_endline print_string printf", "true false None Some"),
			foldOpen: wordSet("begin struct sig object"), foldClose: wordSet("end")})
	addLanguage(&Language{ID: "erlang", Name: "Erlang", Extensions: []string{"erl", "hrl", "escript"}, Shebangs: []string{"escript"}, LineComment: "%"},
		&genericDef{lineComments: []string{"%"}, strings: []stringDef{strDQ, {open: "'", close: "'", escape: '\\', style: StyleConstant}}, preprocessor: '-',
			words:    kw("after and andalso band begin bnot bor bsl bsr bxor case catch cond div end fun if let not of or orelse receive rem try when xor maybe else", "", "io format spawn self length hd tl element", "true false undefined ok error"),
			foldOpen: wordSet("case if receive try fun begin"), foldClose: wordSet("end")})
	addLanguage(&Language{ID: "elixir", Name: "Elixir", Extensions: []string{"ex", "exs", "eex", "heex", "leex"}, Shebangs: []string{"elixir"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, strings: []stringDef{strTDQ, strDQ, strSQ}, annotation: '@',
			words:    kw("def defp defmodule defmacro defmacrop defstruct defprotocol defimpl defdelegate defexception defguard do end if else unless case cond fn when import alias require use quote unquote receive try catch rescue after raise in and or not with for", "", "IO Enum List Map String Kernel", "true false nil"),
			foldOpen: wordSet("do fn"), foldClose: wordSet("end")})
	lispWords := kw("def defn defn- defmacro defmulti defmethod defprotocol defrecord deftype let letfn fn if if-not when when-not do cond case loop recur ns require import use try catch finally throw quote defun defvar defparameter defconstant setq setf lambda progn prog1 let* flet labels unless dolist dotimes define define-syntax set! begin",
		"", "println prn str map filter reduce first rest cons list vector hash-map car cdr format apply funcall", "nil true false t")
	addLanguage(&Language{ID: "lisp", Name: "Lisp", Extensions: []string{"lisp", "lsp", "cl", "el", "scm", "ss", "rkt", "clj", "cljs", "cljc", "edn", "fnl"}, LineComment: ";", BlockComment: [2]string{"#|", "|#"}},
		&genericDef{lineComments: []string{";"}, blockComments: []blockDef{{open: "#|", close: "|#", style: StyleComment, nested: true}}, strings: []stringDef{strDQ},
			identExtra: "-*+!?<>=/.:&%", words: lispWords, foldBraces: "()[]{}"})

	// Data and config
	addLanguage(&Language{ID: "yaml", Name: "YAML", Extensions: []string{"yaml", "yml", "clang-format", "clang-tidy"}, FileNames: []string{".clang-format", ".clang-tidy"}, LineComment: "#", FoldIndent: true}, nil).lexer = yamlLexer{}
	addLanguage(&Language{ID: "ini", Name: "INI", Extensions: []string{"ini", "inf", "cfg", "reg", "url", "desktop", "toml", "service", "editorconfig", "gitconfig", "npmrc", "pylintrc", "flake8", "wixproj"}, FileNames: []string{".editorconfig", ".gitconfig", ".npmrc", "setup.cfg", "tox.ini", ".pylintrc"}, LineComment: ";"}, nil).lexer = iniLexer{}
	addLanguage(&Language{ID: "properties", Name: "Properties", Extensions: []string{"properties", "conf", "env", "gitattributes", "gitignore", "dockerignore", "hgignore", "npmignore", "gitmodules"}, FileNames: []string{".env", ".gitignore", ".gitattributes", ".dockerignore", ".gitmodules", ".hgignore", ".npmignore"}, LineComment: "#"}, nil).lexer = iniLexer{props: true}
	addLanguage(&Language{ID: "diff", Name: "Diff", Extensions: []string{"diff", "patch", "rej"}}, nil).lexer = diffLexer{}
	addLanguage(&Language{ID: "sql", Name: "SQL", Extensions: []string{"sql", "ddl", "dml", "pgsql", "mysql", "psql", "plsql", "tsql", "hql", "cql"}, LineComment: "--", BlockComment: [2]string{"/*", "*/"}},
		&genericDef{lineComments: []string{"--", "#"}, blockComments: []blockDef{cComment}, strings: []stringDef{sqlString, {open: `"`, close: `"`, style: StyleVariable, doubledClose: true}, {open: "`", close: "`", style: StyleVariable}},
			caseInsensitive: true, varPrefixes: "@", foldBraces: "()",
			words: kw("select from where insert into values update set delete create table alter drop index view trigger procedure function begin end if else then case when as and or not null is in like between join inner left right outer full cross on group by order having limit offset union all distinct exists primary key foreign references unique default check constraint database schema grant revoke commit rollback transaction declare return returns asc desc with recursive cast top fetch next rows only over partition row_number rank truncate replace merge using matched execute exec go use show describe explain analyze vacuum add column rename to modify cascade restrict temporary temp if_exists auto_increment identity sequence",
				"int integer bigint smallint tinyint decimal numeric float real double precision varchar nvarchar char nchar text ntext date time timestamp datetime datetime2 interval boolean bool blob clob bytea serial bigserial uuid json jsonb xml money bit varbinary binary image",
				"count sum avg min max coalesce nullif substring substr upper lower trim ltrim rtrim length len now getdate current_date current_timestamp round floor ceil abs isnull ifnull concat convert extract date_trunc string_agg group_concat",
				"true false null")})
	addLanguage(&Language{ID: "protobuf", Name: "Protocol Buffers", Extensions: []string{"proto"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		cLike(kw("syntax edition package import option message enum service rpc returns repeated optional required oneof map reserved extensions extend stream to max weak public",
			"double float int32 int64 uint32 uint64 sint32 sint64 fixed32 fixed64 sfixed32 sfixed64 bool string bytes", "", "true false")))
	addLanguage(&Language{ID: "graphql", Name: "GraphQL", Extensions: []string{"graphql", "gql", "graphqls"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, strings: []stringDef{strTDQ, strDQ}, foldBraces: "{}", varPrefixes: "$", annotation: '@',
			words: kw("query mutation subscription fragment on type interface union enum input scalar schema extend implements directive repeatable", "Int Float String Boolean ID", "", "true false null")})

	// Build
	addLanguage(&Language{ID: "makefile", Name: "Makefile", Extensions: []string{"mak", "mk", "make", "mak.in"}, FileNames: []string{"makefile", "gnumakefile", "makefile.in", "makefile.am", "bsdmakefile"}, Shebangs: []string{"make"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, strings: []stringDef{strDQ, strSQ}, varPrefixes: "$", varParens: true, labels: true, preprocessor: 0,
			words:    kw("ifeq ifneq ifdef ifndef else endif include sinclude -include define endef export unexport override vpath", "", "subst patsubst strip findstring filter filter-out sort word wordlist words firstword lastword dir notdir suffix basename addsuffix addprefix join wildcard realpath abspath error warning shell origin flavor foreach call eval value info", ".PHONY .SUFFIXES .DEFAULT .PRECIOUS .INTERMEDIATE .SECONDARY .DELETE_ON_ERROR .IGNORE .SILENT"),
			foldOpen: wordSet("ifeq ifneq ifdef ifndef define"), foldClose: wordSet("endif endef"), foldMid: wordSet("else")})
	addLanguage(&Language{ID: "cmake", Name: "CMake", Extensions: []string{"cmake"}, FileNames: []string{"cmakelists.txt"}, LineComment: "#"},
		&genericDef{lineComments: []string{"#"}, blockComments: []blockDef{{open: "#[[", close: "]]", style: StyleComment}}, strings: []stringDef{strDQ}, varPrefixes: "$", caseInsensitive: true,
			words: kw("if else elseif endif foreach endforeach while endwhile function endfunction macro endmacro return break continue block endblock",
				"", "set unset option project cmake_minimum_required add_executable add_library target_link_libraries include_directories target_include_directories target_compile_definitions target_compile_options find_package find_library find_path message install add_subdirectory list string file include configure_file add_custom_command add_custom_target add_definitions add_dependencies set_target_properties get_filename_component execute_process",
				"on off true false yes no private public interface required"),
			foldOpen: wordSet("if foreach while function macro block"), foldClose: wordSet("endif endforeach endwhile endfunction endmacro endblock"), foldMid: wordSet("else elseif")})
	addLanguage(&Language{ID: "dockerfile", Name: "Dockerfile", Extensions: []string{"dockerfile", "containerfile"}, FileNames: []string{"dockerfile", "containerfile"}, LineComment: "#"},
		&genericDef{lineCommentsAtStart: []string{"#"}, strings: []stringDef{strDQ, strSQ}, varPrefixes: "$", caseInsensitive: true,
			words: kw("from run cmd label expose env add copy entrypoint volume user workdir arg onbuild stopsignal healthcheck shell as maintainer", "", "", "")})

	// Other programming languages
	addLanguage(&Language{ID: "pascal", Name: "Pascal", Extensions: []string{"pas", "pp", "dpr", "dpk", "lpr", "lpk", "dfm", "inc"}, LineComment: "//", BlockComment: [2]string{"{", "}"}},
		&genericDef{lineComments: []string{"//"}, blockComments: []blockDef{{open: "{$", close: "}", style: StylePreprocessor}, {open: "{", close: "}", style: StyleComment}, {open: "(*", close: "*)", style: StyleComment}},
			strings: []stringDef{sqlString}, caseInsensitive: true, varPrefixes: "",
			words: kw("and array as asm begin case class const constructor destructor div do downto else end except exports file finalization finally for function goto if implementation in inherited initialization inline interface is label library mod nil not object of or out packed procedure program property raise record repeat resourcestring set shl shr string then threadvar to try type unit until uses var while with xor private protected public published override virtual abstract overload reintroduce strict",
				"integer cardinal shortint smallint longint int64 byte word longword real single double extended currency boolean char widechar ansichar pchar pointer variant", "writeln write readln inc dec length setlength high low ord chr", "true false nil self result"),
			foldOpen: wordSet("begin case record try repeat class object interface"), foldClose: wordSet("end until"), foldNoOpenAfter: wordSet("of")})
	addLanguage(&Language{ID: "vb", Name: "Visual Basic", Extensions: []string{"vb", "vbs", "bas", "cls", "frm", "vba", "wsf"}, LineComment: "'"},
		&genericDef{lineComments: []string{"'"}, lineCommentsAtStart: []string{"rem"}, strings: []stringDef{{open: `"`, close: `"`, style: StyleString, doubledClose: true}}, caseInsensitive: true, preprocessor: '#',
			words: kw("addhandler addressof alias and andalso as byref byval call case catch class const continue declare default delegate dim do each else elseif end enum erase error event exit explicit finally for friend function get global gosub goto handles if implements imports in inherits interface is isnot let lib like loop me mod module mustinherit mustoverride mybase myclass namespace narrowing new next not nothing notinheritable notoverridable of on operator option optional or orelse overloads overridable overrides paramarray partial private property protected public raiseevent readonly redim removehandler resume return select set shadows shared static step stop structure sub synclock then throw to try typeof until using when while widening with withevents writeonly xor wend",
				"boolean byte char date decimal double integer long object sbyte short single string uinteger ulong ushort variant currency", "msgbox inputbox createobject len left right mid trim cstr cint clng cdbl", "true false nothing null empty"),
			foldOpen: wordSet("sub function for do while select with class module property structure enum try namespace interface"), foldClose: wordSet("end next loop wend"), foldNoOpenAfter: wordSet("end exit declare")})
	addLanguage(&Language{ID: "fortran", Name: "Fortran", Extensions: []string{"f", "for", "f90", "f95", "f03", "f08", "f77", "fpp"}, LineComment: "!"},
		&genericDef{lineComments: []string{"!"}, strings: []stringDef{strDQ, sqlString}, caseInsensitive: true,
			words: kw("program end subroutine function module use implicit none if then else elseif endif do enddo while call return stop contains type allocate deallocate print write read open close select case default where forall interface procedure intent in out inout parameter dimension allocatable pointer target save data common equivalence external intrinsic include recursive pure elemental result exit cycle goto continue format",
				"integer real double precision character logical complex", "size shape abs sqrt sin cos exp log max min mod sum", ".true. .false."),
			foldOpen: wordSet("program subroutine function module do if select interface type where"), foldClose: wordSet("end enddo endif"), foldNoOpenAfter: wordSet("end")})
	addLanguage(&Language{ID: "matlab", Name: "MATLAB", Extensions: []string{"matlab", "octave"}, LineComment: "%", BlockComment: [2]string{"%{", "%}"}},
		&genericDef{lineComments: []string{"%"}, blockComments: []blockDef{{open: "%{", close: "%}", style: StyleComment}}, strings: []stringDef{strDQ, sqlString},
			words:    kw("break case catch classdef continue else elseif end for function global if otherwise parfor persistent return spmd switch try while methods properties events enumeration", "", "disp fprintf sprintf zeros ones size length numel plot figure", "true false pi inf nan eps"),
			foldOpen: wordSet("function if for while switch try parfor classdef methods properties"), foldClose: wordSet("end"), foldMid: wordSet("else elseif catch")})
	addLanguage(&Language{ID: "asm", Name: "Assembly", Extensions: []string{"asm", "s", "nasm", "masm", "inc"}, LineComment: ";"},
		&genericDef{lineComments: []string{";", "#"}, lineCommentAfterSpace: false, strings: []stringDef{strDQ, strSQ}, caseInsensitive: true, labels: true, preprocessor: '%', identStartExtra: ".",
			words: kw("mov movzx movsx push pop call ret jmp je jne jz jnz jg jge jl jle ja jae jb jbe jc jnc js jns loop add sub mul imul div idiv inc dec neg cmp test and or xor not shl shr sal sar rol ror lea int nop hlt cli sti syscall sysenter leave enter xchg cmov cmove cmovne rep movsb stosb lodsb",
				"eax ebx ecx edx esi edi esp ebp rax rbx rcx rdx rsi rdi rsp rbp r8 r9 r10 r11 r12 r13 r14 r15 ax bx cx dx si di sp bp al ah bl bh cl ch dl dh cs ds es fs gs ss",
				"db dw dd dq dt resb resw resd resq times equ section segment global extern org bits use16 use32 use64 align byte word dword qword ptr .text .data .bss .globl .section .align .ascii .asciz .byte .word .long .quad", "")})
	addLanguage(&Language{ID: "vhdl", Name: "VHDL", Extensions: []string{"vhd", "vhdl", "vho"}, LineComment: "--"},
		&genericDef{lineComments: []string{"--"}, strings: []stringDef{strDQ}, caseInsensitive: true,
			words: kw("abs access after alias all and architecture array assert attribute begin block body buffer bus case component configuration constant disconnect downto else elsif end entity exit file for function generate generic group guarded if impure in inertial inout is label library linkage literal loop map mod nand new next nor not null of on open or others out package port postponed procedure process pure range record register reject rem report return rol ror select severity signal shared sla sll sra srl subtype then to transport type unaffected units until use variable wait when while with xnor xor",
				"std_logic std_logic_vector std_ulogic signed unsigned integer natural positive boolean bit bit_vector string real time", "rising_edge falling_edge to_integer to_unsigned resize", "true false"),
			foldOpen: wordSet("process begin case loop generate record"), foldClose: wordSet("end"), foldNoOpenAfter: wordSet("end")})
	addLanguage(&Language{ID: "verilog", Name: "Verilog", Extensions: []string{"v", "sv", "svh", "vh"}, LineComment: "//", BlockComment: [2]string{"/*", "*/"}},
		func() *genericDef {
			d := cLike(kw("always and assign automatic begin case casex casez cell config deassign default defparam design disable edge else end endcase endconfig endfunction endgenerate endmodule endprimitive endspecify endtable endtask event for force forever fork function generate genvar if ifnone incdir include initial inout input instance join large liblist library localparam macromodule medium module nand negedge nmos nor noshowcancelled not notif0 notif1 or output parameter pmos posedge primitive pull0 pull1 pulldown pullup rcmos real realtime reg release repeat rnmos rpmos rtran rtranif0 rtranif1 scalared signed small specify specparam strong0 strong1 supply0 supply1 table task time tran tranif0 tranif1 tri tri0 tri1 triand trior trireg unsigned use vectored wait wand weak0 weak1 while wire wor xnor xor logic bit byte int always_ff always_comb always_latch interface endinterface package endpackage class endclass",
				"", "$display $monitor $finish $time $stop $write", ""))
			d.foldBraces = ""
			d.preprocessor = '`'
			d.foldOpen = wordSet("begin module case casex casez function task fork generate class interface package")
			d.foldClose = wordSet("end endmodule endcase endfunction endtask join endgenerate endclass endinterface endpackage")
			return d
		}())
	addLanguage(&Language{ID: "nim", Name: "Nim", Extensions: []string{"nim", "nims", "nimble"}, LineComment: "#", BlockComment: [2]string{"#[", "]#"}, FoldIndent: true},
		&genericDef{lineComments: []string{"#"}, blockComments: []blockDef{{open: "#[", close: "]#", style: StyleComment, nested: true}}, strings: []stringDef{strTDQ, strDQ, charSQ},
			words: kw("addr and as asm bind block break case cast concept const continue converter defer discard distinct div do elif else end enum except export finally for from func if import in include interface is isnot iterator let macro method mixin mod nil not notin object of or out proc ptr raise ref return shl shr static template try tuple type using var when while xor yield",
				"int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 float float32 float64 bool char string seq array set", "echo len", "true false nil")})
	addLanguage(&Language{ID: "ada", Name: "Ada", Extensions: []string{"ada", "adb", "ads"}, LineComment: "--"},
		&genericDef{lineComments: []string{"--"}, strings: []stringDef{{open: `"`, close: `"`, style: StyleString, doubledClose: true}}, caseInsensitive: true,
			words: kw("abort abs abstract accept access aliased all and array at begin body case constant declare delay delta digits do else elsif end entry exception exit for function generic goto if in interface is limited loop mod new not null of or others out overriding package pragma private procedure protected raise range record rem renames requeue return reverse select separate some subtype synchronized tagged task terminate then type until use when while with xor",
				"integer natural positive float boolean character string duration", "", "true false"),
			foldOpen: wordSet("begin loop record case select"), foldClose: wordSet("end"), foldNoOpenAfter: wordSet("end")})
	addLanguage(&Language{ID: "tex", Name: "TeX", Extensions: []string{"tex", "sty", "cls", "bib", "ltx", "dtx", "ins"}, LineComment: "%"},
		&genericDef{lineComments: []string{"%"}, commandPrefix: '\\', foldBraces: "{}",
			words: kw("", "", "", "")})
	addLanguage(&Language{ID: "autoit", Name: "AutoIt", Extensions: []string{"au3"}, LineComment: ";", BlockComment: [2]string{"#cs", "#ce"}},
		&genericDef{lineComments: []string{";"}, blockComments: []blockDef{{open: "#cs", close: "#ce", style: StyleComment}, {open: "#comments-start", close: "#comments-end", style: StyleComment}},
			strings: []stringDef{{open: `"`, close: `"`, style: StyleString, doubledClose: true}, sqlString}, caseInsensitive: true, varPrefixes: "$", preprocessor: '#',
			words:    kw("and byref case const continuecase continueloop default dim do else elseif endfunc endif endselect endswitch endwith enum exit exitloop false for func global if in local next not or redim return select static step switch then to true until wend while with", "", "msgbox consolewrite run sleep send", "true false"),
			foldOpen: wordSet("func if for while do select switch with"), foldClose: wordSet("endfunc endif next wend until endselect endswitch endwith"), foldMid: wordSet("else elseif")})
	addLanguage(&Language{ID: "nsis", Name: "NSIS", Extensions: []string{"nsi", "nsh"}, LineComment: ";", BlockComment: [2]string{"/*", "*/"}},
		&genericDef{lineComments: []string{";", "#"}, blockComments: []blockDef{cComment}, strings: []stringDef{strDQ, strSQ, strRawBT}, caseInsensitive: true, varPrefixes: "$", preprocessor: '!',
			words:    kw("function functionend section sectionend sectiongroup sectiongroupend page pageex pageexend var name outfile installdir requestexecutionlevel setoutpath file delete rmdir createdirectory createshortcut writeregstr readregstr deleteregkey messagebox exec execwait call return goto strcmp intcmp", "", "", "true false"),
			foldOpen: wordSet("function section sectiongroup pageex"), foldClose: wordSet("functionend sectionend sectiongroupend pageexend")})
	addLanguage(&Language{ID: "innosetup", Name: "Inno Setup", Extensions: []string{"iss", "isl"}, LineComment: ";"}, nil).lexer = iniLexer{}

	// Internal
	addLanguage(&Language{ID: "searchresults", Name: "Search results", Hidden: true}, nil).lexer = searchResultsLexer{}
}

func sortedWords(m map[string]Style) []string {
	out := make([]string, 0, len(m))
	for w := range m {
		out = append(out, w)
	}
	sort.Strings(out)
	return out
}

func withPre(d *genericDef) *genericDef {
	d.preprocessor = '#'
	return d
}

func withAnnotation(d *genericDef) *genericDef {
	d.annotation = '@'
	return d
}

func jsDef(words map[string]Style) *genericDef {
	d := cLike(words)
	d.strings = []stringDef{strDQ, strSQ, strBT}
	d.identStartExtra = "$"
	d.identExtra = "$"
	d.regexLiterals = true
	d.foldBraces = "{}[]"
	return d
}
