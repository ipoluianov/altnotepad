package config

// WindowState is the main window layout, restored on the next start
type WindowState struct {
	X, Y          int
	Width, Height int
	Maximized     bool

	SecondViewWidth int `json:",omitempty"`
	ResultsHeight   int `json:",omitempty"`
	SideWidth       int `json:",omitempty"`
	FuncListWidth   int `json:",omitempty"`

	// The side panels open: the folder of the workspace and the function list
	Workspace string `json:",omitempty"`
	FuncList  bool   `json:",omitempty"`
	DocMap    bool   `json:",omitempty"`
}

// LoadWindowState returns the saved window layout; ok is false when there is none
func LoadWindowState() (state WindowState, ok bool) {
	if !readJSON("window.json", &state) || state.Width <= 0 || state.Height <= 0 {
		return WindowState{}, false
	}
	return state, true
}

func SaveWindowState(state WindowState) error {
	return writeJSON("window.json", state)
}

// SessionFile is a document of the session
type SessionFile struct {
	// Path of the file; "" for a new document never saved
	Path string `json:",omitempty"`
	// Title of a new document, like "new 1"
	Title string `json:",omitempty"`
	// Backup is the file of the backup directory with the unsaved text
	Backup string `json:",omitempty"`
	// Views: which views show the document, bit 0 - the first one
	Views int

	Language string `json:",omitempty"`
	// LanguageSet: the language was chosen by hand
	LanguageSet bool   `json:",omitempty"`
	Encoding    string `json:",omitempty"`
	BOM         bool   `json:",omitempty"`
	EOL         int
	ReadOnly    bool `json:",omitempty"`

	Caret     int
	Anchor    int
	TopLine   int
	Caret2    int   `json:",omitempty"`
	Anchor2   int   `json:",omitempty"`
	TopLine2  int   `json:",omitempty"`
	Folds     []int `json:",omitempty"`
	Bookmarks []int `json:",omitempty"`
	// ModTime of the file when backed up (Unix nano), to tell whether it changed since
	ModTime int64 `json:",omitempty"`
}

// Session is the set of the documents open when the application closed
type Session struct {
	Files []SessionFile
	// Order of the tabs of each view: indexes into Files
	Tabs [2][]int
	// Active is the index of the active tab of each view (into Tabs), ActiveView the view with the focus
	Active     [2]int
	ActiveView int
}

// LoadSession returns the saved session named name ("" - the last one)
func LoadSession(name string) (Session, bool) {
	var s Session
	if name == "" {
		name = "session.json"
	}
	return s, readJSON(name, &s)
}

// SaveSession saves the session as the last one
func SaveSession(s Session) error {
	return writeJSON("session.json", s)
}
