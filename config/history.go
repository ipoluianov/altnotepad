package config

import "sync"

// History is what the application remembers between the starts: the
// recent files and the texts searched
type History struct {
	RecentFiles  []string `json:",omitempty"`
	FindTexts    []string `json:",omitempty"`
	ReplaceTexts []string `json:",omitempty"`
	FindDirs     []string `json:",omitempty"`
	FindFilters  []string `json:",omitempty"`
	LastDir      string   `json:",omitempty"`
	// Macros are the saved macros by name
	Macros map[string][]MacroStep `json:",omitempty"`
	// RunCommands are the commands of Run
	RunCommands []string `json:",omitempty"`

	FindMatchCase bool `json:",omitempty"`
	FindWholeWord bool `json:",omitempty"`
	FindWrap      bool
	FindMode      int  `json:",omitempty"`
	FindDotAll    bool `json:",omitempty"`
	FindSubdirs   bool
	FindHidden    bool `json:",omitempty"`
}

// MacroStep is a step of a saved macro
type MacroStep struct {
	Cmd string
	Arg string `json:",omitempty"`
}

var (
	historyMtx sync.Mutex
	history    = History{FindWrap: true, FindSubdirs: true}
)

const maxHistory = 20

func loadHistory() {
	h := History{FindWrap: true, FindSubdirs: true}
	readJSON("history.json", &h)
	historyMtx.Lock()
	history = h
	historyMtx.Unlock()
}

// GetHistory returns a copy of the history
func GetHistory() History {
	historyMtx.Lock()
	defer historyMtx.Unlock()
	h := history
	h.RecentFiles = append([]string(nil), history.RecentFiles...)
	return h
}

// UpdateHistory changes the history with f and saves it
func UpdateHistory(f func(h *History)) {
	historyMtx.Lock()
	f(&history)
	h := history
	historyMtx.Unlock()
	_ = writeJSON("history.json", h)
}

// PushFront puts s first in the list, without duplicates, keeping at most max items
func PushFront(list []string, s string, maxItems int) []string {
	if maxItems <= 0 {
		maxItems = maxHistory
	}
	out := []string{s}
	for _, x := range list {
		if x != s && len(out) < maxItems {
			out = append(out, x)
		}
	}
	return out
}
