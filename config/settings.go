package config

import "sync"

// Settings are the options of the application (Settings > Preferences)
type Settings struct {
	// Language of the interface as a tag like "ru"; "" - the system's
	Language string `json:",omitempty"`
	// Color theme of the interface: "light"; "" - the dark one
	Theme string `json:",omitempty"`
	// Scheme is the color scheme of the text; "" - the one for the theme
	Scheme string `json:",omitempty"`

	FontSize float64
	FontFile string `json:",omitempty"`
	Zoom     int    `json:",omitempty"`

	TabSize      int
	InsertSpaces bool
	AutoIndent   bool
	AutoClose    bool
	AutoComplete bool
	SmartHome    bool
	MultiEdit    bool
	ScrollPast   bool
	CaretWidth   int
	CaretBlink   int
	EdgeColumn   int
	WordChars    string `json:",omitempty"`
	// CopyLineNoSel: Ctrl+C and Ctrl+X without a selection copy the line
	CopyLineNoSel bool

	WordWrap         bool
	ShowSpaces       bool
	ShowEOL          bool
	ShowIndentGuides bool
	ShowWrapSymbol   bool
	LineNumbers      bool
	BookmarkMargin   bool
	FoldMargin       bool
	CurrentLine      bool
	SmartHighlight   bool
	SmartMatchCase   bool
	SmartWholeWord   bool
	BraceMatch       bool
	ChangeHistory    bool

	// New documents: "crlf", "lf", "cr"; "" - the system's
	NewDocEOL string `json:",omitempty"`
	// "utf-8", "utf-8-bom" or a character set; "" - UTF-8
	NewDocEncoding string `json:",omitempty"`
	// ANSI is the character set of the files that are not Unicode; "" - by the language
	ANSI string `json:",omitempty"`
	// NewDocLanguage is the language ID of new documents; "" - the plain text
	NewDocLanguage string `json:",omitempty"`

	// RememberSession reopens the files of the last session
	RememberSession bool
	// Backup keeps the unsaved documents between the sessions, so the
	// application closes without asking to save them
	Backup bool
	// BackupSeconds is how often the unsaved changes are backed up
	BackupSeconds int
	// LargeFileMB: files above it open without highlighting, folding and wrapping
	LargeFileMB int
	// CheckFileChanges reports the files changed by other programs
	CheckFileChanges bool
	// AutoReload reloads the changed files that have no unsaved changes without asking
	AutoReload bool

	AlwaysOnTop   bool `json:",omitempty"`
	ShowToolbar   bool
	ShowStatusBar bool
	// TabCloseButtons shows a close button on each tab
	TabCloseButtons bool
	// DoubleClickCloses: a double click on a tab closes it
	DoubleClickCloses bool `json:",omitempty"`
	// SingleInstance opens the files in the running window
	SingleInstance bool
	MaxRecent      int

	// LangTabs are the tab settings of languages that differ from the general ones
	LangTabs map[string]LangTab `json:",omitempty"`
}

// LangTab is the tab setting of a language
type LangTab struct {
	TabSize      int
	InsertSpaces bool
}

var (
	settingsMtx sync.Mutex
	settings    = DefaultSettings()
)

// DefaultSettings returns the settings of the first start
func DefaultSettings() Settings {
	return Settings{
		FontSize: 14, TabSize: 4, AutoIndent: true, AutoClose: false, AutoComplete: true, SmartHome: true, MultiEdit: true,
		CaretWidth: 2, CaretBlink: 530, CopyLineNoSel: true,
		ShowWrapSymbol: true, LineNumbers: true, BookmarkMargin: true, FoldMargin: true, CurrentLine: true,
		SmartHighlight: true, SmartWholeWord: true, BraceMatch: true, ChangeHistory: true,
		RememberSession: true, Backup: true, BackupSeconds: 7, LargeFileMB: 200,
		CheckFileChanges: true, ShowToolbar: true, ShowStatusBar: true, TabCloseButtons: true,
		SingleInstance: true, MaxRecent: 15,
		LangTabs: map[string]LangTab{
			"python": {TabSize: 4, InsertSpaces: true}, "yaml": {TabSize: 2, InsertSpaces: true},
			"makefile": {TabSize: 8, InsertSpaces: false}, "haskell": {TabSize: 4, InsertSpaces: true},
			"nim": {TabSize: 2, InsertSpaces: true}, "coffeescript": {TabSize: 2, InsertSpaces: true},
			"fsharp": {TabSize: 4, InsertSpaces: true},
		},
	}
}

func settingsFile() string { return "settings.json" }

// LoadSettings reads the settings file; missing values get the defaults
func LoadSettings() {
	s := DefaultSettings()
	readJSON(settingsFile(), &s)
	if s.TabSize < 1 || s.TabSize > 16 {
		s.TabSize = 4
	}
	if s.FontSize < 6 || s.FontSize > 72 {
		s.FontSize = 14
	}
	if s.BackupSeconds < 1 {
		s.BackupSeconds = 7
	}
	if s.LargeFileMB < 1 {
		s.LargeFileMB = 200
	}
	if s.MaxRecent < 0 {
		s.MaxRecent = 15
	}
	settingsMtx.Lock()
	settings = s
	settingsMtx.Unlock()
}

// GetSettings returns the current settings; safe to call from any goroutine
func GetSettings() Settings {
	settingsMtx.Lock()
	defer settingsMtx.Unlock()
	return settings
}

// SetSettings applies and saves the settings
func SetSettings(s Settings) error {
	settingsMtx.Lock()
	settings = s
	settingsMtx.Unlock()
	return writeJSON(settingsFile(), s)
}

// TabFor returns the tab setting of the language
func (s Settings) TabFor(lang string) LangTab {
	if t, ok := s.LangTabs[lang]; ok && t.TabSize > 0 {
		return t
	}
	return LangTab{TabSize: s.TabSize, InsertSpaces: s.InsertSpaces}
}
