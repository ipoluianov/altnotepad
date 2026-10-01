package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui"
	"github.com/ipoluianov/nui/ui/i18n"
)

// Strings are all the texts of the application. A misspelled field is a
// compile error; a field a language leaves empty falls back to English
// (strings_test.go checks that none is left).
type Strings struct {
	Error      string
	NewDocName string // "new" of "new 1"
	NormalText string

	// Main menu
	MenuFile, MenuEdit, MenuSearch, MenuView, MenuEncoding, MenuLanguage, MenuSettings, MenuTools, MenuMacro, MenuRunMenu, MenuWindow string

	// File
	MenuNew, MenuOpen, MenuOpenContainingFolderSub, MenuOpenContainingFolder, MenuOpenDefaultViewer, MenuOpenFolderWorkspace string
	MenuReload, MenuSave, MenuSaveAs, MenuSaveCopy, MenuSaveAll, MenuRename, MenuClose, MenuCloseAll, MenuCloseMore          string
	MenuCloseAllButActive, MenuCloseLeft, MenuCloseRight, MenuCloseUnchanged, MenuDeleteFile, MenuLoadSession                string
	MenuSaveSession, MenuRecentFiles, MenuClearRecent, MenuRestoreClosed, MenuExit                                           string
	MenuFunctionList, Filter, MenuDocumentMap                                                                                string
	MenuChangeHistory, MenuNextChange, MenuPrevChange, MenuClearChanges                                                      string
	MenuPrint, MenuExportHTML                                                                                                string
	SavedTo                                                                                                                  func(path string) string
	MenuMultiSelect, MenuMultiAll, MenuMultiAllCase, MenuMultiNext, MenuMultiUndo, MenuMultiSkip                             string

	// Edit
	MenuUndo, MenuRedo, MenuCut, MenuCopy, MenuPaste, MenuDelete, MenuSelectAll, MenuInsert, MenuDateShort, MenuDateLong, MenuDateISO string
	MenuCopyToClipboard, MenuCopyFullPath, MenuCopyFileName, MenuCopyDirPath, MenuIndentSub, MenuIndent, MenuUnindent                 string
	MenuConvertCase, MenuUpperCase, MenuLowerCase, MenuProperCase, MenuProperCaseBlend, MenuSentenceCase, MenuSentenceCaseBlend       string
	MenuInvertCase, MenuRandomCase, MenuLineOperations, MenuDuplicateLine, MenuRemoveDupLines, MenuRemoveConsecutiveDupLines          string
	MenuSplitLines, MenuJoinLines, MenuMoveLineUp, MenuMoveLineDown, MenuDeleteLine, MenuCutLine, MenuTransposeLine                   string
	MenuRemoveEmptyLines, MenuRemoveBlankLines, MenuInsertLineAbove, MenuInsertLineBelow, MenuReverseLines, MenuShuffleLines          string
	MenuSortAsc, MenuSortDesc, MenuSortAscCI, MenuSortDescCI, MenuSortIntAsc, MenuSortIntDesc, MenuSortDecAsc, MenuSortDecDesc        string
	MenuSortLenAsc, MenuSortLenDesc, MenuComment, MenuToggleComment, MenuLineComment, MenuLineUncomment, MenuBlockComment             string
	MenuBlockUncomment, MenuAutoCompletion, MenuWordCompletion, MenuEOLConversion, MenuEOLWindows, MenuEOLUnix, MenuEOLMac            string
	MenuBlankOperations, MenuTrimTrailing, MenuTrimLeading, MenuTrimBoth, MenuEOLToSpace, MenuRemoveBlankEOL, MenuTabToSpace          string
	MenuSpaceToTabAll, MenuSpaceToTabLeading, MenuColumnEditor, MenuReadOnly                                                          string
	DateShortLayout, DateLongLayout                                                                                                   string

	// Search
	MenuFind, MenuFindInFiles, MenuFindNext, MenuFindPrev, MenuSelectFindNext, MenuSelectFindPrev, MenuReplace, MenuIncremental string
	MenuMark, MenuSearchResults, MenuNextResult, MenuPrevResult, MenuGoTo, MenuGotoBrace, MenuSelectBrace, MenuStyleAll         string
	MenuStyleOne, MenuClearStyleSub, MenuClearAllStyles, MenuJumpUp, MenuJumpDown, MenuBookmark, MenuToggleBookmark             string
	MenuNextBookmark, MenuPrevBookmark, MenuClearBookmarks, MenuCutBookmarked, MenuCopyBookmarked, MenuRemoveBookmarked         string
	MenuRemoveUnbookmarked, MenuInverseBookmarks                                                                                string
	MenuUsingStyle, MenuClearStyle                                                                                              func(n int) string

	// View
	MenuAlwaysOnTop, MenuFullScreen, MenuShowSymbol, MenuShowSpaces, MenuShowEOL, MenuShowAllChars, MenuShowIndentGuides string
	MenuShowWrapSymbol, MenuZoom, MenuZoomIn, MenuZoomOut, MenuZoomReset, MenuMoveClone, MenuMoveToOtherView             string
	MenuCloneToOtherView, MenuTab, MenuNextTab, MenuPrevTab, MenuMoveTabForward, MenuMoveTabBackward, MenuWordWrap       string
	MenuLineNumbers, MenuFocusOtherView, MenuFoldAll, MenuUnfoldAll, MenuFoldCurrent, MenuUnfoldCurrent, MenuFoldLevel   string
	MenuUnfoldLevel, MenuSummary, MenuFolderWorkspace, MenuMonitoring, MenuToolbar, MenuStatusBar                        string
	MenuTabN                                                                                                             func(n int) string

	// Encoding
	MenuEncANSI, MenuEncUTF8, MenuEncUTF8BOM, MenuEncUTF16BE, MenuEncUTF16LE, MenuCharacterSets string
	MenuConvANSI, MenuConvUTF8, MenuConvUTF8BOM, MenuConvUTF16BE, MenuConvUTF16LE               string
	CharsetGroups                                                                               map[string]string

	// Settings, tools, macro, run, window, help
	MenuPreferences, MenuColorScheme, MenuShortcuts                                                     string
	MenuHashGenerate, MenuHashFiles, MenuHashSelection                                                  func(name string) string
	MenuBase64Encode, MenuBase64Decode, MenuURLEncode, MenuURLDecode, MenuJSONFormat, MenuJSONMinify    string
	MenuStartRecording, MenuStopRecording, MenuPlayback, MenuSaveMacro, MenuRunMacroMulti, MenuTrimSave string
	MenuRun, MenuOpenInBrowser, MenuSearchInternet, MenuOpenSelectedFile                                string
	MenuWindows, MenuSortTabsByName, MenuSortTabsByPath, MenuHelp, MenuHomePage, MenuAbout              string

	// Files
	OpenTitle, SaveTitle, SaveAsTitle, SaveCopyTitle, AllFiles, TextFiles, Save, DontSave, PathLabel string
	SaveChangesAsk, CreateFileAsk, AlreadyOpen, FileExists, FileNotFound                             func(path string) string
	Unencodable                                                                                      func(enc string) string
	Loading                                                                                          func(name string, percent int) string
	LargeFileMode, LargeFileMark                                                                     string
	ReloadTitle, RenameTitle, NewName, DeleteFileTitle                                               string
	ReloadAsk, DeleteFileAsk                                                                         func(path string) string
	FileChangedTitle, FileDeletedTitle                                                               string
	FileChangedAsk, FileChangedModifiedAsk, FileDeletedAsk, NotASession                              func(path string) string

	// Status bar
	StatusSize func(length, lines int) string
	StatusPos  func(line, col, pos int) string
	StatusSel  string

	// Dialogs
	GotoTitle, GotoLine, GotoOffset, Go                                                               string
	GotoLineInfo, GotoOffsetInfo                                                                      func(cur, total int) string
	ColumnTitle, ColumnText, ColumnNumbers, ColumnInitial, ColumnIncrease, ColumnRepeat, ColumnFormat string
	ColumnLeading, LeadingNone, LeadingZeros, LeadingSpaces                                           string
	FullPath, Modified, SummaryChars, SummaryWords, SummaryLines, SummaryNonBlankLines                string
	SummaryLength, SummarySelected, Bytes                                                             string
	RunTitle, RunLabel, Run                                                                           string
	Name, State, ModifiedState, Activate, CloseWindows, Command, Shortcut                             string
	HashInput, HashEachLine, CopyToClipboard                                                          string
	CopiedToClipboard                                                                                 func(s string) string
	MacroName, MacroTimes                                                                             string

	// Find
	FindTab, ReplaceTab, FindInFilesTab, MarkTab, FindWhat, ReplaceWith, Filters, Directory, CurrentFolder string
	MatchCase, WholeWord, WrapAround, Backward, InSelection, InSubfolders, InHidden, BookmarkLine          string
	PurgeEach, SearchMode, ModeNormal, ModeExtended, ModeRegex, DotAll                                     string
	FindNext, FindPrev, Count, FindAllCurrent, FindAllOpen, Replace, ReplaceAll, ReplaceAllOpen            string
	FindAll, ReplaceInFiles, MarkAll, ClearMarks, CopyMarked, Wrapped, ReadOnlyDoc                         string
	NotFound                                                                                               func(text string) string
	CountResult, ReplacedCount, MarkedCount, FoundCount, HitsCount                                         func(n int) string
	FoundInFiles, ReplacedInFiles                                                                          func(n, files int) string
	ReplaceInFilesAsk                                                                                      func(find, repl, dir string) string
	NoSuchFolder, Searching                                                                                func(path string) string
	SearchResults, Stop, Stopped, Clear, Line                                                              string
	ResultsHeader                                                                                          func(text string, hits, files, searched int) string

	// Workspace
	Workspace, Refresh, CollapseAll string

	// Settings dialog
	PageGeneral, PageEditing, PageDisplay, PageFiles                                               string
	Language, LanguageSystem, Theme, ThemeDark, ThemeLight, ColorScheme, SchemeByTheme             string
	ShowToolbar, ShowStatusBar, TabCloseButtons, DoubleClickCloses, AlwaysOnTop, SingleInstance    string
	RememberSession, BackupSession, BackupEvery, MaxRecent                                         string
	TabSize, InsertSpaces, AutoIndent, AutoClose, AutoCompletion, SmartHome, MultiEdit, ScrollPast string
	CopyLineNoSel, CaretWidth, CaretBlink, EdgeColumn, WordChars                                   string
	FontSize, FontFile, BuiltInFont, Fonts, LineNumbers, BookmarkMargin, FoldMargin, CurrentLine   string
	SmartHighlight, SmartMatchCase, SmartWholeWord, BraceMatch, WrapSymbol, ChangeHistory          string
	NewDocEOL, NewDocEncoding, ANSICharset, NewDocLanguage, SystemDefault, ByLanguage              string
	LargeFileMB, CheckFileChanges, AutoReload                                                      string

	// About
	AboutTitle                                                      func(name string) string
	AboutDescription, Version, Author, License, VisitWebsite, Close string
}

// languages are offered in the settings; the names stay in their own language
var languages = []struct{ tag, name string }{
	{"en", "English"},
	{"ru", "Русский"},
	{"pl", "Polski"},
	{"sr", "Српски"},
	{"de", "Deutsch"},
	{"fr", "Français"},
	{"es", "Español"},
	{"it", "Italiano"},
	{"pt", "Português"},
	{"zh", "中文"},
	{"ja", "日本語"},
	{"ko", "한국어"},
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

var en = Strings{
	Error:      "Error",
	NewDocName: "new",
	NormalText: "Normal text",

	MenuFile: "File", MenuEdit: "Edit", MenuSearch: "Search", MenuView: "View", MenuEncoding: "Encoding", MenuLanguage: "Language",
	MenuSettings: "Settings", MenuTools: "Tools", MenuMacro: "Macro", MenuRunMenu: "Run", MenuWindow: "Window",

	MenuNew: "New", MenuOpen: "Open...", MenuOpenContainingFolderSub: "Open Containing Folder",
	MenuOpenContainingFolder: "Explorer", MenuOpenDefaultViewer: "Open in Default Viewer", MenuOpenFolderWorkspace: "Open Folder as Workspace...",
	MenuReload: "Reload from Disk", MenuSave: "Save", MenuSaveAs: "Save As...", MenuSaveCopy: "Save a Copy As...", MenuSaveAll: "Save All",
	MenuRename: "Rename...", MenuClose: "Close", MenuCloseAll: "Close All", MenuCloseMore: "Close Multiple Documents",
	MenuCloseAllButActive: "Close All but Active Document", MenuCloseLeft: "Close All to the Left", MenuCloseRight: "Close All to the Right",
	MenuCloseUnchanged: "Close All Unchanged", MenuDeleteFile: "Delete from Disk...", MenuLoadSession: "Load Session...",
	MenuSaveSession: "Save Session...", MenuRecentFiles: "Recent Files", MenuClearRecent: "Empty Recent Files List",
	MenuRestoreClosed: "Restore Recent Closed File", MenuExit: "Exit",
	MenuPrint: "Print...", MenuExportHTML: "Export to HTML...", MenuFunctionList: "Function List", Filter: "Filter", MenuDocumentMap: "Document Map",
	MenuChangeHistory: "Change History", MenuNextChange: "Go to Next Change", MenuPrevChange: "Go to Previous Change", MenuClearChanges: "Clear Change History",
	SavedTo:         func(path string) string { return "Saved to " + path },
	MenuMultiSelect: "Multi-select", MenuMultiAll: "Multi-select All", MenuMultiAllCase: "Multi-select All (Match Case)",
	MenuMultiNext: "Multi-select Next", MenuMultiUndo: "Undo the Latest Added Multi-select", MenuMultiSkip: "Skip Current & Go to Next Multi-select",

	MenuUndo: "Undo", MenuRedo: "Redo", MenuCut: "Cut", MenuCopy: "Copy", MenuPaste: "Paste", MenuDelete: "Delete", MenuSelectAll: "Select All",
	MenuInsert: "Insert", MenuDateShort: "Date Time (short)", MenuDateLong: "Date Time (long)", MenuDateISO: "Date Time (ISO 8601)",
	MenuCopyToClipboard: "Copy to Clipboard", MenuCopyFullPath: "Current Full File Path", MenuCopyFileName: "Current Filename",
	MenuCopyDirPath: "Current Directory Path", MenuIndentSub: "Indent", MenuIndent: "Increase Line Indent", MenuUnindent: "Decrease Line Indent",
	MenuConvertCase: "Convert Case to", MenuUpperCase: "UPPERCASE", MenuLowerCase: "lowercase", MenuProperCase: "Proper Case",
	MenuProperCaseBlend: "Proper Case (blend)", MenuSentenceCase: "Sentence case", MenuSentenceCaseBlend: "Sentence case (blend)",
	MenuInvertCase: "iNVERT cASE", MenuRandomCase: "ranDOm CasE", MenuLineOperations: "Line Operations",
	MenuDuplicateLine: "Duplicate Current Line", MenuRemoveDupLines: "Remove Duplicate Lines", MenuRemoveConsecutiveDupLines: "Remove Consecutive Duplicate Lines",
	MenuSplitLines: "Split Lines", MenuJoinLines: "Join Lines", MenuMoveLineUp: "Move Up Current Line", MenuMoveLineDown: "Move Down Current Line",
	MenuDeleteLine: "Delete Current Line", MenuCutLine: "Cut Current Line", MenuTransposeLine: "Swap Current Line with Previous",
	MenuRemoveEmptyLines: "Remove Empty Lines", MenuRemoveBlankLines: "Remove Empty Lines (Containing Blank characters)",
	MenuInsertLineAbove: "Insert Blank Line Above Current", MenuInsertLineBelow: "Insert Blank Line Below Current",
	MenuReverseLines: "Reverse Line Order", MenuShuffleLines: "Randomize Line Order",
	MenuSortAsc: "Sort Lines Lexicographically Ascending", MenuSortDesc: "Sort Lines Lexicographically Descending",
	MenuSortAscCI: "Sort Lines Lex. Ascending Ignoring Case", MenuSortDescCI: "Sort Lines Lex. Descending Ignoring Case",
	MenuSortIntAsc: "Sort Lines As Integers Ascending", MenuSortIntDesc: "Sort Lines As Integers Descending",
	MenuSortDecAsc: "Sort Lines As Decimals Ascending", MenuSortDecDesc: "Sort Lines As Decimals Descending",
	MenuSortLenAsc: "Sort Lines By Length Ascending", MenuSortLenDesc: "Sort Lines By Length Descending",
	MenuComment: "Comment/Uncomment", MenuToggleComment: "Toggle Single Line Comment", MenuLineComment: "Single Line Comment",
	MenuLineUncomment: "Single Line Uncomment", MenuBlockComment: "Block Comment", MenuBlockUncomment: "Block Uncomment",
	MenuAutoCompletion: "Auto-Completion", MenuWordCompletion: "Word Completion", MenuEOLConversion: "EOL Conversion",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Blank Operations", MenuTrimTrailing: "Trim Trailing Space", MenuTrimLeading: "Trim Leading Space",
	MenuTrimBoth: "Trim Leading and Trailing Space", MenuEOLToSpace: "EOL to Space", MenuRemoveBlankEOL: "Remove Unnecessary Blank and EOL",
	MenuTabToSpace: "TAB to Space", MenuSpaceToTabAll: "Space to TAB (All)", MenuSpaceToTabLeading: "Space to TAB (Leading)",
	MenuColumnEditor: "Column Editor...", MenuReadOnly: "Read-Only",
	DateShortLayout: "15:04 02.01.2006", DateLongLayout: "15:04:05 Monday, January 2, 2006",

	MenuFind: "Find...", MenuFindInFiles: "Find in Files...", MenuFindNext: "Find Next", MenuFindPrev: "Find Previous",
	MenuSelectFindNext: "Select and Find Next", MenuSelectFindPrev: "Select and Find Previous", MenuReplace: "Replace...",
	MenuIncremental: "Incremental Search", MenuMark: "Mark...", MenuSearchResults: "Search Results Window",
	MenuNextResult: "Next Search Result", MenuPrevResult: "Previous Search Result", MenuGoTo: "Go to...",
	MenuGotoBrace: "Go to Matching Brace", MenuSelectBrace: "Select All In-between {} [] or ()",
	MenuStyleAll: "Style All Occurrences of Token", MenuStyleOne: "Style One Token", MenuClearStyleSub: "Clear Style",
	MenuClearAllStyles: "Clear all Styles", MenuJumpUp: "Jump Up", MenuJumpDown: "Jump Down", MenuBookmark: "Bookmark",
	MenuToggleBookmark: "Toggle Bookmark", MenuNextBookmark: "Next Bookmark", MenuPrevBookmark: "Previous Bookmark",
	MenuClearBookmarks: "Clear All Bookmarks", MenuCutBookmarked: "Cut Bookmarked Lines", MenuCopyBookmarked: "Copy Bookmarked Lines",
	MenuRemoveBookmarked: "Remove Bookmarked Lines", MenuRemoveUnbookmarked: "Remove Non-Bookmarked Lines", MenuInverseBookmarks: "Inverse Bookmark",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Using %s Style", ordinal(n)) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Clear %s Style", ordinal(n)) },

	MenuAlwaysOnTop: "Always on Top", MenuFullScreen: "Toggle Full Screen Mode", MenuShowSymbol: "Show Symbol",
	MenuShowSpaces: "Show Space and Tab", MenuShowEOL: "Show End of Line", MenuShowAllChars: "Show All Characters",
	MenuShowIndentGuides: "Show Indent Guide", MenuShowWrapSymbol: "Show Wrap Symbol", MenuZoom: "Zoom", MenuZoomIn: "Zoom In",
	MenuZoomOut: "Zoom Out", MenuZoomReset: "Restore Default Zoom", MenuMoveClone: "Move/Clone Current Document",
	MenuMoveToOtherView: "Move to Other View", MenuCloneToOtherView: "Clone to Other View", MenuTab: "Tab",
	MenuNextTab: "Next Tab", MenuPrevTab: "Previous Tab", MenuMoveTabForward: "Move Tab Forward", MenuMoveTabBackward: "Move Tab Backward",
	MenuWordWrap: "Word Wrap", MenuLineNumbers: "Line Numbers", MenuFocusOtherView: "Focus on Another View", MenuFoldAll: "Fold All",
	MenuUnfoldAll: "Unfold All", MenuFoldCurrent: "Fold Current Level", MenuUnfoldCurrent: "Unfold Current Level",
	MenuFoldLevel: "Fold Level", MenuUnfoldLevel: "Unfold Level", MenuSummary: "Summary...", MenuFolderWorkspace: "Folder as Workspace",
	MenuMonitoring: "Monitoring (tail -f)", MenuToolbar: "Toolbar", MenuStatusBar: "Status Bar",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Last Tab"
		}
		return fmt.Sprintf("%s Tab", ordinal(n))
	},

	MenuEncANSI: "ANSI", MenuEncUTF8: "UTF-8", MenuEncUTF8BOM: "UTF-8-BOM", MenuEncUTF16BE: "UTF-16 BE BOM", MenuEncUTF16LE: "UTF-16 LE BOM",
	MenuCharacterSets: "Character Sets", MenuConvANSI: "Convert to ANSI", MenuConvUTF8: "Convert to UTF-8", MenuConvUTF8BOM: "Convert to UTF-8-BOM",
	MenuConvUTF16BE: "Convert to UTF-16 BE BOM", MenuConvUTF16LE: "Convert to UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Arabic", "Baltic": "Baltic", "Celtic": "Celtic", "Cyrillic": "Cyrillic",
		"Central European": "Central European", "Chinese": "Chinese", "Greek": "Greek", "Hebrew": "Hebrew", "Japanese": "Japanese",
		"Korean": "Korean", "North European": "North European", "Thai": "Thai", "Turkish": "Turkish", "Vietnamese": "Vietnamese",
		"Western European": "Western European"},

	MenuPreferences: "Preferences...", MenuColorScheme: "Color Scheme", MenuShortcuts: "Keyboard Shortcuts",
	MenuHashGenerate:  func(name string) string { return "Generate " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Generate " + name + " from Files..." },
	MenuHashSelection: func(name string) string { return "Copy " + name + " of Selection" },
	MenuBase64Encode:  "Base64 Encode", MenuBase64Decode: "Base64 Decode", MenuURLEncode: "URL Encode", MenuURLDecode: "URL Decode",
	MenuJSONFormat: "JSON Format", MenuJSONMinify: "JSON Minify",
	MenuStartRecording: "Start/Stop Recording", MenuStopRecording: "Stop Recording", MenuPlayback: "Playback",
	MenuSaveMacro: "Save Current Recorded Macro...", MenuRunMacroMulti: "Run a Macro Multiple Times...", MenuTrimSave: "Trim Trailing Space and Save",
	MenuRun: "Run...", MenuOpenInBrowser: "Open in Browser", MenuSearchInternet: "Search on Internet", MenuOpenSelectedFile: "Open File (Selected Name)",
	MenuWindows: "Windows...", MenuSortTabsByName: "Sort Tabs by Name", MenuSortTabsByPath: "Sort Tabs by Path",
	MenuHelp: "Online Help", MenuHomePage: "Home Page", MenuAbout: "About AltNotepad",

	OpenTitle: "Open", SaveTitle: "Save", SaveAsTitle: "Save As", SaveCopyTitle: "Save a Copy As", AllFiles: "All files", TextFiles: "Text files",
	Save: "Save", DontSave: "Don't Save", PathLabel: "Path of the file:",
	SaveChangesAsk: func(path string) string { return "Save changes to \"" + path + "\"?" },
	CreateFileAsk:  func(path string) string { return "\"" + path + "\" doesn't exist. Create it?" },
	AlreadyOpen:    func(path string) string { return "\"" + path + "\" is open in another tab" },
	FileExists:     func(path string) string { return "\"" + path + "\" already exists" },
	FileNotFound:   func(path string) string { return "File not found: " + path },
	Unencodable: func(enc string) string {
		return "Some characters cannot be saved in " + enc + ". Save them as \"?\" anyway?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Loading %s: %d%%", name, percent) },
	LargeFileMode: "Large file: syntax highlighting, folding and word wrap are off",
	LargeFileMark: "(large file)",
	ReloadTitle:   "Reload", RenameTitle: "Rename", NewName: "New name:", DeleteFileTitle: "Delete from Disk",
	ReloadAsk:        func(path string) string { return "Reload \"" + path + "\"? The unsaved changes will be lost." },
	DeleteFileAsk:    func(path string) string { return "Delete \"" + path + "\" from the disk?" },
	FileChangedTitle: "File Changed", FileDeletedTitle: "File Deleted",
	FileChangedAsk: func(path string) string {
		return "\"" + path + "\"\n\nThis file has been modified by another program.\nDo you want to reload it?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "\"" + path + "\"\n\nThis file has been modified by another program.\nReload it and lose the changes made in the editor?"
	},
	FileDeletedAsk: func(path string) string {
		return "\"" + path + "\"\n\nThis file doesn't exist anymore.\nKeep it in the editor?"
	},
	NotASession: func(name string) string { return name + " is not a session file" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("length: %d   lines: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Ln: %d   Col: %d   Pos: %d", line, col, pos) },
	StatusSel:  "Sel:",

	GotoTitle: "Go To...", GotoLine: "Line", GotoOffset: "Offset", Go: "Go",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("You are here: %d   You want to go to (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("You are here: %d   You want to go to (0 - %d):", cur, total)
	},
	ColumnTitle: "Column / Multi-Selection Editor", ColumnText: "Text to Insert", ColumnNumbers: "Number to Insert",
	ColumnInitial: "Initial number:", ColumnIncrease: "Increase by:", ColumnRepeat: "Repeat:", ColumnFormat: "Format:",
	ColumnLeading: "Leading:", LeadingNone: "None", LeadingZeros: "Zeros", LeadingSpaces: "Spaces",
	FullPath: "Full file path", Modified: "Modified", SummaryChars: "Characters (without line endings)", SummaryWords: "Words",
	SummaryLines: "Lines", SummaryNonBlankLines: "Non-blank lines", SummaryLength: "Document length", SummarySelected: "Selected characters",
	Bytes: "bytes", RunTitle: "Run", RunLabel: "The program to run:", Run: "Run",
	Name: "Name", State: "State", ModifiedState: "modified", Activate: "Activate", CloseWindows: "Close Window(s)",
	Command: "Command", Shortcut: "Shortcut", HashInput: "Text:", HashEachLine: "Treat each line as a separate string",
	CopyToClipboard:   "Copy to Clipboard",
	CopiedToClipboard: func(s string) string { return "Copied to the clipboard: " + s },
	MacroName:         "Name of the macro:", MacroTimes: "Run how many times (* - up to the end of the file):",

	FindTab: "Find", ReplaceTab: "Replace", FindInFilesTab: "Find in Files", MarkTab: "Mark", FindWhat: "Find what:",
	ReplaceWith: "Replace with:", Filters: "Filters:", Directory: "Directory:", CurrentFolder: "Current folder",
	MatchCase: "Match case", WholeWord: "Match whole word only", WrapAround: "Wrap around", Backward: "Backward direction",
	InSelection: "In selection", InSubfolders: "In all sub-folders", InHidden: "In hidden folders", BookmarkLine: "Bookmark line",
	PurgeEach: "Purge for each search", SearchMode: "Search Mode:", ModeNormal: "Normal", ModeExtended: "Extended (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Regular expression", DotAll: ". matches newline",
	FindNext: "Find Next", FindPrev: "Find Previous", Count: "Count", FindAllCurrent: "Find All in Current Document",
	FindAllOpen: "Find All in All Opened Documents", Replace: "Replace", ReplaceAll: "Replace All",
	ReplaceAllOpen: "Replace All in All Opened Documents", FindAll: "Find All", ReplaceInFiles: "Replace in Files",
	MarkAll: "Mark All", ClearMarks: "Clear all marks", CopyMarked: "Copy Marked Text", Wrapped: "Reached the end, continued from the start",
	ReadOnlyDoc: "The document is read-only",
	NotFound:    func(text string) string { return "Can't find the text \"" + text + "\"" },
	CountResult: func(n int) string { return fmt.Sprintf("Count: %d %s", n, plural(n, "match", "matches")) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d %s replaced", n, plural(n, "occurrence was", "occurrences were"))
	},
	MarkedCount: func(n int) string { return fmt.Sprintf("%d %s marked", n, plural(n, "match", "matches")) },
	FoundCount:  func(n int) string { return fmt.Sprintf("%d %s found", n, plural(n, "hit", "hits")) },
	HitsCount:   func(n int) string { return fmt.Sprintf("(%d %s)", n, plural(n, "hit", "hits")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s in %d %s", n, plural(n, "hit", "hits"), files, plural(files, "file", "files"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s replaced in %d %s", n, plural(n, "occurrence", "occurrences"), files, plural(files, "file", "files"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Replace all \"" + find + "\" with \"" + repl + "\" in all the files of\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "No such folder: " + path },
	Searching:     func(path string) string { return "Searching: " + path },
	SearchResults: "Search results", Stop: "Stop", Stopped: "(stopped)", Clear: "Clear", Line: "Line",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Search \"%s\" (%d %s in %d %s of %d searched)", text, hits, plural(hits, "hit", "hits"), files, plural(files, "file", "files"), searched)
	},

	Workspace: "Workspace", Refresh: "Refresh", CollapseAll: "Collapse All",

	PageGeneral: "General", PageEditing: "Editing", PageDisplay: "Display", PageFiles: "New Document & Files",
	Language: "Language:", LanguageSystem: "As in the system", Theme: "Theme:", ThemeDark: "Dark", ThemeLight: "Light",
	ColorScheme: "Color scheme:", SchemeByTheme: "By the theme",
	ShowToolbar: "Show the toolbar", ShowStatusBar: "Show the status bar", TabCloseButtons: "Close buttons on the tabs",
	DoubleClickCloses: "Double click closes a tab", AlwaysOnTop: "Always on top", SingleInstance: "Open files in the running window",
	RememberSession: "Remember the open files for the next session", BackupSession: "Keep unsaved changes between sessions (backup)",
	BackupEvery: "Backup every, seconds:", MaxRecent: "Recent files in the list:",
	TabSize: "Tab size:", InsertSpaces: "Replace tabs with spaces", AutoIndent: "Auto-indent", AutoClose: "Auto-close brackets and quotes",
	AutoCompletion: "Word completion while typing", SmartHome: "Home goes to the first non-blank character", MultiEdit: "Multi-editing (Ctrl+Click)",
	ScrollPast: "Scroll beyond the last line", CopyLineNoSel: "Copy/cut the line when nothing is selected", CaretWidth: "Caret width:",
	CaretBlink: "Caret blink rate, ms (0 - no blinking):", EdgeColumn: "Long line marker column (0 - none):", WordChars: "Extra word characters:",
	FontSize: "Font size:", FontFile: "Font file:", BuiltInFont: "JetBrains Mono (built-in)", Fonts: "Fonts", LineNumbers: "Line numbers",
	BookmarkMargin: "Bookmark margin", FoldMargin: "Folding margin", CurrentLine: "Highlight the current line", SmartHighlight: "Smart highlighting",
	SmartMatchCase: "Smart highlighting: match case", SmartWholeWord: "Smart highlighting: whole word only", BraceMatch: "Highlight matching braces",
	WrapSymbol: "Show the wrap symbol", ChangeHistory: "Change history in the margin",
	NewDocEOL: "Line breaks of new documents:", NewDocEncoding: "Encoding of new documents:", ANSICharset: "ANSI character set:",
	NewDocLanguage: "Language of new documents:", SystemDefault: "As in the system", ByLanguage: "By the interface language",
	LargeFileMB: "Large file restriction, MB:", CheckFileChanges: "Detect files changed by other programs", AutoReload: "Reload them silently when not modified",

	AboutTitle:       func(name string) string { return "About " + name },
	AboutDescription: "A text and source code editor", Version: "Version", Author: "Author:", License: "License:",
	VisitWebsite: "Visit Website", Close: "Close",
}

func ordinal(n int) string {
	switch n {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	}
	return fmt.Sprintf("%dth", n)
}

var ru = Strings{
	Error:      "Ошибка",
	NewDocName: "новый",
	NormalText: "Обычный текст",

	MenuFile: "Файл", MenuEdit: "Правка", MenuSearch: "Поиск", MenuView: "Вид", MenuEncoding: "Кодировки", MenuLanguage: "Синтаксис",
	MenuSettings: "Опции", MenuTools: "Инструменты", MenuMacro: "Макросы", MenuRunMenu: "Запуск", MenuWindow: "Окна",

	MenuNew: "Новый", MenuOpen: "Открыть...", MenuOpenContainingFolderSub: "Открыть папку документа",
	MenuOpenContainingFolder: "В проводнике", MenuOpenDefaultViewer: "Открыть в программе по умолчанию", MenuOpenFolderWorkspace: "Открыть папку как рабочее пространство...",
	MenuReload: "Перезагрузить с диска", MenuSave: "Сохранить", MenuSaveAs: "Сохранить как...", MenuSaveCopy: "Сохранить копию как...", MenuSaveAll: "Сохранить все",
	MenuRename: "Переименовать...", MenuClose: "Закрыть", MenuCloseAll: "Закрыть все", MenuCloseMore: "Закрыть несколько документов",
	MenuCloseAllButActive: "Закрыть все, кроме текущего", MenuCloseLeft: "Закрыть все слева", MenuCloseRight: "Закрыть все справа",
	MenuCloseUnchanged: "Закрыть все неизменённые", MenuDeleteFile: "Удалить с диска...", MenuLoadSession: "Загрузить сессию...",
	MenuSaveSession: "Сохранить сессию...", MenuRecentFiles: "Недавние файлы", MenuClearRecent: "Очистить список недавних файлов",
	MenuRestoreClosed: "Восстановить последний закрытый файл", MenuExit: "Выход",
	MenuPrint: "Печать...", MenuExportHTML: "Экспорт в HTML...", MenuFunctionList: "Список функций", Filter: "Фильтр", MenuDocumentMap: "Карта документа",
	MenuChangeHistory: "История изменений", MenuNextChange: "К следующему изменению", MenuPrevChange: "К предыдущему изменению", MenuClearChanges: "Очистить историю изменений",
	SavedTo:         func(path string) string { return "Сохранено: " + path },
	MenuMultiSelect: "Мультивыделение", MenuMultiAll: "Выделить все вхождения", MenuMultiAllCase: "Выделить все вхождения (с учётом регистра)",
	MenuMultiNext: "Добавить следующее вхождение", MenuMultiUndo: "Отменить последнее добавленное", MenuMultiSkip: "Пропустить текущее и добавить следующее",

	MenuUndo: "Отменить", MenuRedo: "Повторить", MenuCut: "Вырезать", MenuCopy: "Копировать", MenuPaste: "Вставить", MenuDelete: "Удалить", MenuSelectAll: "Выделить всё",
	MenuInsert: "Вставить", MenuDateShort: "Дату и время (кратко)", MenuDateLong: "Дату и время (полностью)", MenuDateISO: "Дату и время (ISO 8601)",
	MenuCopyToClipboard: "Копировать в буфер обмена", MenuCopyFullPath: "Полный путь к файлу", MenuCopyFileName: "Имя файла",
	MenuCopyDirPath: "Путь к папке", MenuIndentSub: "Отступ", MenuIndent: "Увеличить отступ", MenuUnindent: "Уменьшить отступ",
	MenuConvertCase: "Преобразовать регистр", MenuUpperCase: "ВСЕ ПРОПИСНЫЕ", MenuLowerCase: "все строчные", MenuProperCase: "Каждое Слово С Прописной",
	MenuProperCaseBlend: "Каждое Слово С Прописной (смешанно)", MenuSentenceCase: "Как в предложениях", MenuSentenceCaseBlend: "Как в предложениях (смешанно)",
	MenuInvertCase: "иНВЕРТИРОВАТЬ рЕГИСТР", MenuRandomCase: "сЛуЧаЙнЫй РеГиСтР", MenuLineOperations: "Операции со строками",
	MenuDuplicateLine: "Дублировать текущую строку", MenuRemoveDupLines: "Удалить повторяющиеся строки", MenuRemoveConsecutiveDupLines: "Удалить идущие подряд повторы строк",
	MenuSplitLines: "Разбить строки", MenuJoinLines: "Объединить строки", MenuMoveLineUp: "Переместить строку вверх", MenuMoveLineDown: "Переместить строку вниз",
	MenuDeleteLine: "Удалить текущую строку", MenuCutLine: "Вырезать текущую строку", MenuTransposeLine: "Поменять строку с предыдущей",
	MenuRemoveEmptyLines: "Удалить пустые строки", MenuRemoveBlankLines: "Удалить пустые строки (и из пробелов)",
	MenuInsertLineAbove: "Вставить пустую строку выше", MenuInsertLineBelow: "Вставить пустую строку ниже",
	MenuReverseLines: "Обратный порядок строк", MenuShuffleLines: "Перемешать строки",
	MenuSortAsc: "Сортировать строки по возрастанию", MenuSortDesc: "Сортировать строки по убыванию",
	MenuSortAscCI: "Сортировать по возрастанию без учёта регистра", MenuSortDescCI: "Сортировать по убыванию без учёта регистра",
	MenuSortIntAsc: "Сортировать как целые числа по возрастанию", MenuSortIntDesc: "Сортировать как целые числа по убыванию",
	MenuSortDecAsc: "Сортировать как дробные числа по возрастанию", MenuSortDecDesc: "Сортировать как дробные числа по убыванию",
	MenuSortLenAsc: "Сортировать по длине по возрастанию", MenuSortLenDesc: "Сортировать по длине по убыванию",
	MenuComment: "Комментирование", MenuToggleComment: "Закомментировать/раскомментировать строки", MenuLineComment: "Закомментировать строки",
	MenuLineUncomment: "Раскомментировать строки", MenuBlockComment: "Закомментировать блок", MenuBlockUncomment: "Раскомментировать блок",
	MenuAutoCompletion: "Автодополнение", MenuWordCompletion: "Завершение слова", MenuEOLConversion: "Формат конца строк",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Операции с пробелами", MenuTrimTrailing: "Удалить пробелы в конце строк", MenuTrimLeading: "Удалить пробелы в начале строк",
	MenuTrimBoth: "Удалить пробелы в начале и конце строк", MenuEOLToSpace: "Концы строк в пробелы", MenuRemoveBlankEOL: "Удалить лишние пробелы и концы строк",
	MenuTabToSpace: "Табуляции в пробелы", MenuSpaceToTabAll: "Пробелы в табуляции (все)", MenuSpaceToTabLeading: "Пробелы в табуляции (в начале строк)",
	MenuColumnEditor: "Редактор столбцов...", MenuReadOnly: "Только чтение",
	DateShortLayout: "15:04 02.01.2006", DateLongLayout: "15:04:05 02.01.2006 (Monday)",

	MenuFind: "Найти...", MenuFindInFiles: "Найти в файлах...", MenuFindNext: "Искать далее", MenuFindPrev: "Искать ранее",
	MenuSelectFindNext: "Выделить и искать далее", MenuSelectFindPrev: "Выделить и искать ранее", MenuReplace: "Заменить...",
	MenuIncremental: "Поиск по мере ввода", MenuMark: "Пометить...", MenuSearchResults: "Окно результатов поиска",
	MenuNextResult: "Следующий результат", MenuPrevResult: "Предыдущий результат", MenuGoTo: "Перейти...",
	MenuGotoBrace: "К парной скобке", MenuSelectBrace: "Выделить всё между скобками",
	MenuStyleAll: "Пометить все вхождения", MenuStyleOne: "Пометить одно вхождение", MenuClearStyleSub: "Снять пометку",
	MenuClearAllStyles: "Снять все пометки", MenuJumpUp: "К предыдущей пометке", MenuJumpDown: "К следующей пометке", MenuBookmark: "Закладки",
	MenuToggleBookmark: "Поставить/убрать закладку", MenuNextBookmark: "Следующая закладка", MenuPrevBookmark: "Предыдущая закладка",
	MenuClearBookmarks: "Убрать все закладки", MenuCutBookmarked: "Вырезать строки с закладками", MenuCopyBookmarked: "Копировать строки с закладками",
	MenuRemoveBookmarked: "Удалить строки с закладками", MenuRemoveUnbookmarked: "Удалить строки без закладок", MenuInverseBookmarks: "Инвертировать закладки",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Стилем %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Снять стиль %d", n) },

	MenuAlwaysOnTop: "Поверх всех окон", MenuFullScreen: "Полноэкранный режим", MenuShowSymbol: "Отображение символов",
	MenuShowSpaces: "Пробелы и табуляции", MenuShowEOL: "Концы строк", MenuShowAllChars: "Все символы",
	MenuShowIndentGuides: "Направляющие отступов", MenuShowWrapSymbol: "Символ переноса", MenuZoom: "Масштаб", MenuZoomIn: "Увеличить",
	MenuZoomOut: "Уменьшить", MenuZoomReset: "Восстановить масштаб", MenuMoveClone: "Переместить/клонировать документ",
	MenuMoveToOtherView: "Переместить в другое окно", MenuCloneToOtherView: "Клонировать в другое окно", MenuTab: "Вкладки",
	MenuNextTab: "Следующая вкладка", MenuPrevTab: "Предыдущая вкладка", MenuMoveTabForward: "Переместить вкладку вперёд", MenuMoveTabBackward: "Переместить вкладку назад",
	MenuWordWrap: "Перенос строк", MenuLineNumbers: "Номера строк", MenuFocusOtherView: "Перейти в другое окно", MenuFoldAll: "Свернуть все блоки",
	MenuUnfoldAll: "Развернуть все блоки", MenuFoldCurrent: "Свернуть текущий блок", MenuUnfoldCurrent: "Развернуть текущий блок",
	MenuFoldLevel: "Свернуть уровень", MenuUnfoldLevel: "Развернуть уровень", MenuSummary: "Сводка...", MenuFolderWorkspace: "Папка как рабочее пространство",
	MenuMonitoring: "Слежение за файлом (tail -f)", MenuToolbar: "Панель инструментов", MenuStatusBar: "Строка состояния",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Последняя вкладка"
		}
		return fmt.Sprintf("Вкладка %d", n)
	},

	MenuEncANSI: "Кодировать в ANSI", MenuEncUTF8: "Кодировать в UTF-8", MenuEncUTF8BOM: "Кодировать в UTF-8-BOM", MenuEncUTF16BE: "Кодировать в UTF-16 BE BOM",
	MenuEncUTF16LE: "Кодировать в UTF-16 LE BOM", MenuCharacterSets: "Кодировки", MenuConvANSI: "Преобразовать в ANSI", MenuConvUTF8: "Преобразовать в UTF-8",
	MenuConvUTF8BOM: "Преобразовать в UTF-8-BOM", MenuConvUTF16BE: "Преобразовать в UTF-16 BE BOM", MenuConvUTF16LE: "Преобразовать в UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Арабские", "Baltic": "Балтийские", "Celtic": "Кельтские", "Cyrillic": "Кириллица",
		"Central European": "Центральноевропейские", "Chinese": "Китайские", "Greek": "Греческие", "Hebrew": "Иврит", "Japanese": "Японские",
		"Korean": "Корейские", "North European": "Северноевропейские", "Thai": "Тайские", "Turkish": "Турецкие", "Vietnamese": "Вьетнамские",
		"Western European": "Западноевропейские"},

	MenuPreferences: "Настройки...", MenuColorScheme: "Цветовая схема", MenuShortcuts: "Горячие клавиши",
	MenuHashGenerate:  func(name string) string { return "Вычислить " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Вычислить " + name + " файлов..." },
	MenuHashSelection: func(name string) string { return "Копировать " + name + " выделения" },
	MenuBase64Encode:  "Кодировать в Base64", MenuBase64Decode: "Декодировать из Base64", MenuURLEncode: "Кодировать URL", MenuURLDecode: "Декодировать URL",
	MenuJSONFormat: "Форматировать JSON", MenuJSONMinify: "Сжать JSON",
	MenuStartRecording: "Начать/остановить запись", MenuStopRecording: "Остановить запись", MenuPlayback: "Воспроизвести",
	MenuSaveMacro: "Сохранить записанный макрос...", MenuRunMacroMulti: "Воспроизвести несколько раз...", MenuTrimSave: "Удалить пробелы в конце строк и сохранить",
	MenuRun: "Запустить...", MenuOpenInBrowser: "Открыть в браузере", MenuSearchInternet: "Искать в интернете", MenuOpenSelectedFile: "Открыть файл (выделенное имя)",
	MenuWindows: "Окна...", MenuSortTabsByName: "Сортировать вкладки по имени", MenuSortTabsByPath: "Сортировать вкладки по пути",
	MenuHelp: "Справка", MenuHomePage: "Домашняя страница", MenuAbout: "О программе",

	OpenTitle: "Открыть", SaveTitle: "Сохранение", SaveAsTitle: "Сохранить как", SaveCopyTitle: "Сохранить копию как", AllFiles: "Все файлы",
	TextFiles: "Текстовые файлы", Save: "Сохранить", DontSave: "Не сохранять", PathLabel: "Путь к файлу:",
	SaveChangesAsk: func(path string) string { return "Сохранить изменения в «" + path + "»?" },
	CreateFileAsk:  func(path string) string { return "Файла «" + path + "» нет. Создать его?" },
	AlreadyOpen:    func(path string) string { return "«" + path + "» уже открыт в другой вкладке" },
	FileExists:     func(path string) string { return "«" + path + "» уже существует" },
	FileNotFound:   func(path string) string { return "Файл не найден: " + path },
	Unencodable: func(enc string) string {
		return "Некоторые символы нельзя сохранить в кодировке " + enc + ". Сохранить их как «?»?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Загрузка %s: %d%%", name, percent) },
	LargeFileMode: "Большой файл: подсветка синтаксиса, сворачивание и перенос строк выключены",
	LargeFileMark: "(большой файл)",
	ReloadTitle:   "Перезагрузка", RenameTitle: "Переименование", NewName: "Новое имя:", DeleteFileTitle: "Удаление с диска",
	ReloadAsk: func(path string) string {
		return "Перезагрузить «" + path + "»? Несохранённые изменения будут потеряны."
	},
	DeleteFileAsk:    func(path string) string { return "Удалить «" + path + "» с диска?" },
	FileChangedTitle: "Файл изменён", FileDeletedTitle: "Файл удалён",
	FileChangedAsk: func(path string) string {
		return "«" + path + "»\n\nЭтот файл был изменён другой программой.\nПерезагрузить его?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "«" + path + "»\n\nЭтот файл был изменён другой программой.\nПерезагрузить его и потерять изменения, сделанные в редакторе?"
	},
	FileDeletedAsk: func(path string) string {
		return "«" + path + "»\n\nЭтого файла больше нет.\nОставить его в редакторе?"
	},
	NotASession: func(name string) string { return name + " — не файл сессии" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("длина: %d   строк: %d", length, lines) },
	StatusPos: func(line, col, pos int) string {
		return fmt.Sprintf("Стр: %d   Стлб: %d   Поз: %d", line, col, pos)
	},
	StatusSel: "Выд:",

	GotoTitle: "Перейти...", GotoLine: "Строка", GotoOffset: "Смещение", Go: "Перейти",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Вы здесь: %d   Перейти к (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Вы здесь: %d   Перейти к (0 - %d):", cur, total)
	},
	ColumnTitle: "Редактор столбцов и мультивыделения", ColumnText: "Вставить текст", ColumnNumbers: "Вставить числа",
	ColumnInitial: "Начальное число:", ColumnIncrease: "Увеличивать на:", ColumnRepeat: "Повторять:", ColumnFormat: "Формат:",
	ColumnLeading: "Дополнять:", LeadingNone: "Ничем", LeadingZeros: "Нулями", LeadingSpaces: "Пробелами",
	FullPath: "Полный путь к файлу", Modified: "Изменён", SummaryChars: "Символов (без концов строк)", SummaryWords: "Слов",
	SummaryLines: "Строк", SummaryNonBlankLines: "Непустых строк", SummaryLength: "Длина документа", SummarySelected: "Выделено символов",
	Bytes: "байт", RunTitle: "Запуск", RunLabel: "Программа для запуска:", Run: "Запустить",
	Name: "Имя", State: "Состояние", ModifiedState: "изменён", Activate: "Перейти", CloseWindows: "Закрыть",
	Command: "Команда", Shortcut: "Клавиши", HashInput: "Текст:", HashEachLine: "Считать каждую строку отдельно",
	CopyToClipboard:   "Копировать",
	CopiedToClipboard: func(s string) string { return "Скопировано в буфер обмена: " + s },
	MacroName:         "Название макроса:", MacroTimes: "Сколько раз выполнить (* - до конца файла):",

	FindTab: "Найти", ReplaceTab: "Заменить", FindInFilesTab: "Найти в файлах", MarkTab: "Пометки", FindWhat: "Найти:",
	ReplaceWith: "Заменить на:", Filters: "Фильтры:", Directory: "Папка:", CurrentFolder: "Папка документа",
	MatchCase: "Учитывать регистр", WholeWord: "Только целые слова", WrapAround: "Зациклить поиск", Backward: "Искать назад",
	InSelection: "В выделенном", InSubfolders: "Во всех подпапках", InHidden: "В скрытых папках", BookmarkLine: "Ставить закладки",
	PurgeEach: "Очищать при каждом поиске", SearchMode: "Режим поиска:", ModeNormal: "Обычный", ModeExtended: "Расширенный (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Регулярные выражения", DotAll: ". — любой символ, и конец строки",
	FindNext: "Искать далее", FindPrev: "Искать ранее", Count: "Подсчитать", FindAllCurrent: "Найти всё в текущем документе",
	FindAllOpen: "Найти всё во всех открытых", Replace: "Заменить", ReplaceAll: "Заменить всё",
	ReplaceAllOpen: "Заменить всё во всех открытых", FindAll: "Найти всё", ReplaceInFiles: "Заменить в файлах",
	MarkAll: "Пометить всё", ClearMarks: "Снять все пометки", CopyMarked: "Копировать помеченное", Wrapped: "Достигнут конец, поиск продолжен с начала",
	ReadOnlyDoc: "Документ только для чтения",
	NotFound:    func(text string) string { return "Не найдено: «" + text + "»" },
	CountResult: func(n int) string {
		return fmt.Sprintf("Найдено: %d %s", n, i18n.Plural("ru", n, "совпадение", "совпадения", "совпадений"))
	},
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("Заменено: %d %s", n, i18n.Plural("ru", n, "вхождение", "вхождения", "вхождений"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("Помечено: %d %s", n, i18n.Plural("ru", n, "совпадение", "совпадения", "совпадений"))
	},
	FoundCount: func(n int) string {
		return fmt.Sprintf("Найдено: %d %s", n, i18n.Plural("ru", n, "совпадение", "совпадения", "совпадений"))
	},
	HitsCount: func(n int) string {
		return fmt.Sprintf("(%d %s)", n, i18n.Plural("ru", n, "совпадение", "совпадения", "совпадений"))
	},
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s в %d %s", n, i18n.Plural("ru", n, "совпадение", "совпадения", "совпадений"), files, i18n.Plural("ru", files, "файле", "файлах", "файлах"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("Заменено %d %s в %d %s", n, i18n.Plural("ru", n, "вхождение", "вхождения", "вхождений"), files, i18n.Plural("ru", files, "файле", "файлах", "файлах"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Заменить все «" + find + "» на «" + repl + "» во всех файлах папки\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "Нет такой папки: " + path },
	Searching:     func(path string) string { return "Поиск: " + path },
	SearchResults: "Результаты поиска", Stop: "Остановить", Stopped: "(остановлен)", Clear: "Очистить", Line: "Строка",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Поиск «%s» (%d %s в %d %s из %d просмотренных)", text, hits, i18n.Plural("ru", hits, "совпадение", "совпадения", "совпадений"),
			files, i18n.Plural("ru", files, "файле", "файлах", "файлах"), searched)
	},

	Workspace: "Рабочее пространство", Refresh: "Обновить", CollapseAll: "Свернуть всё",

	PageGeneral: "Общие", PageEditing: "Правка", PageDisplay: "Вид", PageFiles: "Новый документ и файлы",
	Language: "Язык:", LanguageSystem: "Как в системе", Theme: "Тема:", ThemeDark: "Тёмная", ThemeLight: "Светлая",
	ColorScheme: "Цветовая схема:", SchemeByTheme: "По теме",
	ShowToolbar: "Показывать панель инструментов", ShowStatusBar: "Показывать строку состояния", TabCloseButtons: "Кнопки закрытия на вкладках",
	DoubleClickCloses: "Закрывать вкладку двойным щелчком", AlwaysOnTop: "Поверх всех окон", SingleInstance: "Открывать файлы в уже запущенном окне",
	RememberSession: "Запоминать открытые файлы до следующего запуска", BackupSession: "Сохранять несохранённые изменения между запусками",
	BackupEvery: "Резервное копирование каждые, секунд:", MaxRecent: "Недавних файлов в списке:",
	TabSize: "Размер табуляции:", InsertSpaces: "Заменять табуляцию пробелами", AutoIndent: "Автоотступ", AutoClose: "Автозакрытие скобок и кавычек",
	AutoCompletion: "Завершение слов при вводе", SmartHome: "Home — к первому непробельному символу", MultiEdit: "Мультиредактирование (Ctrl+щелчок)",
	ScrollPast: "Прокрутка за последнюю строку", CopyLineNoSel: "Копировать/вырезать строку, если ничего не выделено", CaretWidth: "Ширина курсора:",
	CaretBlink: "Мигание курсора, мс (0 - без мигания):", EdgeColumn: "Граница длинных строк, столбец (0 - нет):", WordChars: "Дополнительные символы слов:",
	FontSize: "Размер шрифта:", FontFile: "Файл шрифта:", BuiltInFont: "JetBrains Mono (встроенный)", Fonts: "Шрифты", LineNumbers: "Номера строк",
	BookmarkMargin: "Поле закладок", FoldMargin: "Поле сворачивания", CurrentLine: "Подсвечивать текущую строку", SmartHighlight: "Умная подсветка",
	SmartMatchCase: "Умная подсветка: с учётом регистра", SmartWholeWord: "Умная подсветка: только целые слова", BraceMatch: "Подсвечивать парные скобки",
	WrapSymbol: "Показывать символ переноса", ChangeHistory: "История изменений на поле",
	NewDocEOL: "Концы строк новых документов:", NewDocEncoding: "Кодировка новых документов:", ANSICharset: "Кодировка ANSI:",
	NewDocLanguage: "Синтаксис новых документов:", SystemDefault: "Как в системе", ByLanguage: "По языку интерфейса",
	LargeFileMB: "Ограничение для больших файлов, МБ:", CheckFileChanges: "Следить за изменением файлов другими программами",
	AutoReload: "Перезагружать их без вопросов, если не изменены",

	AboutTitle:       func(name string) string { return "О программе " + name },
	AboutDescription: "Редактор текста и исходного кода", Version: "Версия", Author: "Автор:", License: "Лицензия:",
	VisitWebsite: "Открыть сайт", Close: "Закрыть",
}

var catalog = i18n.NewCatalog(en, map[string]Strings{
	"ru": ru, "pl": plStrings, "sr": sr, "de": de, "fr": fr, "es": es,
	"it": it, "pt": pt, "zh": zh, "ja": ja, "ko": ko,
})

// The library translates its own buttons (OK, Cancel...) to Russian and Chinese only
func init() {
	for lang, s := range map[string]ui.UIStrings{
		"pl": {OK: "OK", Cancel: "Anuluj", Yes: "Tak", No: "Nie"},
		"sr": {OK: "У реду", Cancel: "Откажи", Yes: "Да", No: "Не"},
		"de": {OK: "OK", Cancel: "Abbrechen", Yes: "Ja", No: "Nein"},
		"fr": {OK: "OK", Cancel: "Annuler", Yes: "Oui", No: "Non"},
		"es": {OK: "Aceptar", Cancel: "Cancelar", Yes: "Sí", No: "No"},
		"it": {OK: "OK", Cancel: "Annulla", Yes: "Sì", No: "No"},
		"pt": {OK: "OK", Cancel: "Cancelar", Yes: "Sim", No: "Não"},
		"ja": {OK: "OK", Cancel: "キャンセル", Yes: "はい", No: "いいえ"},
		"ko": {OK: "확인", Cancel: "취소", Yes: "예", No: "아니요"},
	} {
		ui.RegisterUIStrings(lang, s)
	}
}

// T returns the texts in the language of the application
func T() *Strings {
	return catalog.Get(ui.Language())
}

// SetLanguage switches the application to the language of the settings, "" - the system's
func SetLanguage(lang string) {
	if lang == "" {
		lang = ui.SystemLanguage()
	}
	ui.SetLanguage(lang)
}
