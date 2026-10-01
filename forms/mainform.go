package forms

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ipoluianov/altnotepad/app"
	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// MainForm is the main window: the menu, the tool bar, the two views with
// their tabs, the search panel and results, and the status bar
type MainForm struct {
	ui.Widget

	toolbar    *Toolbar
	sideSplit  *ui.Splitter
	viewSplit  *ui.Splitter
	vertSplit  *ui.Splitter
	funcSplit  *ui.Splitter
	funcList   *FuncListPanel
	rightSplit *ui.Splitter
	docMap     *DocMapPanel
	panes      [2]*Pane
	active     int
	docs       []*Doc
	find       *FindPanel
	results    *ResultsPanel
	status     *StatusBar
	workspace  *Workspace
	menuBar    *ui.MenuBar
	commands   map[string]*Command
	shortcuts  map[Shortcut]*Command
	menuChecks []menuCheck

	untitled int
	// closed are the files closed last, for Restore Recent Closed File
	closed []string

	macro     *editor.Macro
	recording bool

	fullscreen bool
	secondW    int
	resultsH   int
	sideW      int
	funcW      int

	lastTitle string
	quitting  bool
}

var mainForm *MainForm

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	mainForm = &c
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	c.commands = make(map[string]*Command)
	c.shortcuts = make(map[Shortcut]*Command)

	c.panes[0] = newPane(&c, 0)
	c.panes[1] = newPane(&c, 1)
	c.panes[1].SetVisible(false)
	c.viewSplit = ui.NewHSplitter()
	c.viewSplit.SetWidgets(c.panes[0], c.panes[1])
	c.secondW = 500
	c.viewSplit.SetSecondSize(c.secondW)
	c.viewSplit.SetOnSplitChanged(func() {
		if c.panes[1].IsVisible() {
			c.secondW = c.viewSplit.SecondSize()
		}
	})

	c.workspace = NewWorkspace(&c)
	c.workspace.SetVisible(false)
	c.sideSplit = ui.NewHSplitter()
	c.sideSplit.SetWidgets(c.workspace, c.viewSplit)
	c.sideW = 260
	c.sideSplit.SetFirstSize(c.sideW)
	c.sideSplit.SetOnSplitChanged(func() {
		if c.workspace.IsVisible() {
			c.sideW = c.sideSplit.FirstSize()
		}
	})

	c.funcList = NewFuncListPanel(&c)
	c.funcList.SetVisible(false)
	c.docMap = NewDocMapPanel(&c)
	c.docMap.SetVisible(false)
	c.rightSplit = ui.NewVSplitter()
	c.rightSplit.SetWidgets(c.funcList, c.docMap)
	c.rightSplit.SetVisible(false)
	c.funcSplit = ui.NewHSplitter()
	c.funcSplit.SetWidgets(c.sideSplit, c.rightSplit)
	c.funcW = 280
	c.funcSplit.SetSecondSize(c.funcW)
	c.funcSplit.SetOnSplitChanged(func() {
		if c.rightSplit.IsVisible() {
			c.funcW = c.funcSplit.SecondSize()
		}
	})

	c.results = NewResultsPanel(&c)
	c.results.SetVisible(false)
	c.vertSplit = ui.NewVSplitter()
	c.vertSplit.SetWidgets(c.funcSplit, c.results)
	c.resultsH = 220
	c.vertSplit.SetSecondSize(c.resultsH)
	c.vertSplit.SetOnSplitChanged(func() {
		if c.results.IsVisible() {
			c.resultsH = c.vertSplit.SecondSize()
		}
	})

	c.registerCommands()
	c.buildContextMenus()
	c.toolbar = NewToolbar(&c)
	c.find = NewFindPanel(&c)
	c.find.SetVisible(false)
	c.status = NewStatusBar(&c)

	c.AddWidget(0, 0, c.toolbar)
	c.AddWidget(1, 0, c.vertSplit)
	c.AddWidget(2, 0, c.find)
	c.AddWidget(3, 0, c.status)

	c.AddTimer(1000, c.checkFiles)
	c.AddTimer(500, c.backupTimer)
	c.AddTimer(300, c.refreshTitle)

	themeListeners = append(themeListeners, c.applySettingsToViews)
	return &c
}

func (c *MainForm) settings() config.Settings { return config.GetSettings() }

// activePane returns the view with the focus
func (c *MainForm) activePane() *Pane { return c.panes[c.active] }

// view returns the editor of the view with the focus
func (c *MainForm) view() *editor.View { return c.activePane().view }

// curDoc returns the document of the view with the focus, nil if none
func (c *MainForm) curDoc() *Doc { return c.activePane().current() }

func (c *MainForm) setActivePane(i int) {
	if !c.panes[i].IsVisible() {
		return
	}
	if c.active != i {
		c.active = i
		c.docActivated(c.panes[i])
	}
	c.panes[0].tabs.SetFocusedView(i == 0)
	c.panes[1].tabs.SetFocusedView(i == 1)
}

// focusEditor gives the focus to the editor of the active view
func (c *MainForm) focusEditor() {
	c.view().Focus()
}

// docActivated is called when a view shows another document
func (c *MainForm) docActivated(p *Pane) {
	if p.index != c.active {
		return
	}
	c.status.refresh()
	c.toolbar.refresh()
	c.refreshTitle()
	c.workspace.selectFile(c.curPath())
	c.funcList.refresh()
	c.docMap.m.SetView(c.view())
}

func (c *MainForm) curPath() string {
	if d := c.curDoc(); d != nil {
		return d.Path()
	}
	return ""
}

// caretMoved is called when the caret of a view moves
func (c *MainForm) caretMoved(p *Pane) {
	if p.index == c.active {
		c.status.refreshCaret()
	}
}

// docEdited is called after an edit in a view
func (c *MainForm) docEdited(p *Pane) {
	c.status.refresh()
	c.toolbar.refresh()
}

// modifiedChanged is called when a document becomes modified or saved
func (c *MainForm) modifiedChanged(d *Doc) {
	for _, p := range c.panes {
		if p.indexOf(d) >= 0 {
			p.refreshTabs()
		}
	}
	c.refreshTitle()
	c.toolbar.refresh()
}

// refreshTitle shows the file of the active document in the window title
func (c *MainForm) refreshTitle() {
	title := app.DisplayName
	if d := c.curDoc(); d != nil {
		name := d.FullName()
		if d.IsModified() {
			name = "*" + name
		}
		title = name + " - " + app.DisplayName
	}
	if c.recording {
		title = "● " + title
	}
	if title != c.lastTitle && c.Form() != nil {
		c.lastTitle = title
		c.Form().SetTitle(title)
	}
}

// viewOptions builds the editor options from the settings for the document
func (c *MainForm) viewOptions(d *Doc) editor.Options {
	s := c.settings()
	o := editor.DefaultOptions()
	o.FontSize = s.FontSize
	o.FontFile = s.FontFile
	o.Zoom = s.Zoom
	lang := ""
	if d != nil && d.doc.Language() != nil {
		lang = d.doc.Language().ID
	}
	tab := s.TabFor(lang)
	o.TabSize = tab.TabSize
	o.InsertTabs = !tab.InsertSpaces
	o.AutoIndent = s.AutoIndent
	o.AutoClose = s.AutoClose
	o.SmartHome = s.SmartHome
	o.WordChars = s.WordChars
	o.MultiEdit = s.MultiEdit
	o.ScrollPast = s.ScrollPast
	o.CaretWidth = s.CaretWidth
	o.CaretBlink = s.CaretBlink
	o.EdgeColumn = s.EdgeColumn
	o.WordWrap = s.WordWrap
	o.LineNumbers = s.LineNumbers
	o.Bookmarks = s.BookmarkMargin
	o.FoldMargin = s.FoldMargin
	o.ShowSpaces = s.ShowSpaces
	o.ShowEOL = s.ShowEOL
	o.ShowIndentGuides = s.ShowIndentGuides
	o.ShowWrapSymbol = s.ShowWrapSymbol
	o.CurrentLine = s.CurrentLine
	o.SmartHighlight = s.SmartHighlight
	o.SmartMatchCase = s.SmartMatchCase
	o.SmartWholeWord = s.SmartWholeWord
	o.BraceMatch = s.BraceMatch
	o.CopyLineNoSel = s.CopyLineNoSel
	o.AutoComplete = s.AutoComplete
	o.ChangeHistory = s.ChangeHistory
	return o
}

// applyViewOptions applies the settings to the editor of the pane for its document
func (c *MainForm) applyViewOptions(p *Pane) {
	p.view.SetOptions(c.viewOptions(p.current()))
	p.view.SetScheme(editorScheme(c.settings()))
}

// applyDocOptions sets up a document added to a pane
func (c *MainForm) applyDocOptions(d *Doc, pane int) {
	d.doc.OnModifiedChanged = func() { c.modifiedChanged(d) }
}

// applySettingsToViews applies the settings to all the views
func (c *MainForm) applySettingsToViews() {
	for _, p := range c.panes {
		c.applyViewOptions(p)
		p.tabs.SetCloseButtons(c.settings().TabCloseButtons)
	}
	c.results.applySettings()
	c.toolbar.SetVisible(c.settings().ShowToolbar)
	c.status.SetVisible(c.settings().ShowStatusBar)
	c.toolbar.refresh()
	c.status.refresh()
	c.updateLayout()
}

func (c *MainForm) updateLayout() {
	if c.Form() != nil {
		c.Form().UpdateLayout()
		c.Form().Update()
	}
}

// ApplySettings saves the settings and applies them
func (c *MainForm) ApplySettings(s config.Settings) {
	old := c.settings()
	if err := config.SetSettings(s); err != nil {
		c.showError(err)
	}
	if s.Language != old.Language {
		SetLanguage(s.Language)
	}
	applyANSI(s)
	if s.Theme != old.Theme {
		ApplyTheme(s.Theme)
	}
	c.Form().SetAlwaysOnTop(s.AlwaysOnTop)
	c.applySettingsToViews()
}

// applyANSI sets the character set of the non-Unicode files
func applyANSI(s config.Settings) {
	if s.ANSI != "" {
		editor.DefaultANSI = s.ANSI
		return
	}
	lang := s.Language
	if lang == "" {
		lang = ui.SystemLanguage()
	}
	switch lang {
	case "ru", "sr":
		editor.DefaultANSI = "windows-1251"
	case "pl":
		editor.DefaultANSI = "windows-1250"
	case "zh":
		editor.DefaultANSI = "gbk"
	case "ja":
		editor.DefaultANSI = "shift_jis"
	case "ko":
		editor.DefaultANSI = "euc-kr"
	default:
		editor.DefaultANSI = "windows-1252"
	}
}

// updateSetting changes one setting and applies it
func (c *MainForm) updateSetting(f func(s *config.Settings)) {
	s := c.settings()
	f(&s)
	c.ApplySettings(s)
}

func (c *MainForm) zoomChanged(zoom int) {
	s := c.settings()
	s.Zoom = zoom
	config.SetSettings(s)
	for _, p := range c.panes {
		if p.view.Zoom() != zoom {
			p.view.SetZoom(zoom)
		}
	}
	c.status.refresh()
}

// showError shows an error in a message box
func (c *MainForm) showError(err error) {
	ui.ShowMessageBox(c, T().Error, err.Error())
}

// toast shows a short message
func (c *MainForm) toast(text string) {
	if c.Form() != nil {
		c.Form().ShowToastFor(text, ui.ToastInfo, 2500*time.Millisecond)
	}
}

// ---- The second view ----

// showSecondView shows or hides the second view
func (c *MainForm) showSecondView(show bool) {
	if c.panes[1].IsVisible() == show {
		return
	}
	c.panes[1].SetVisible(show)
	if show {
		c.viewSplit.SetSecondSize(max(200, c.secondW))
		if w := c.viewSplit.Width(); w > 0 && c.secondW > w-200 {
			c.viewSplit.SetSecondSize(w / 2)
		}
	} else if c.active == 1 {
		c.setActivePane(0)
		c.focusEditor()
	}
	c.updateLayout()
}

// moveToOtherView moves (or clones) the active document to the other view
func (c *MainForm) moveToOtherView(clone bool) {
	src := c.activePane()
	d := src.current()
	if d == nil {
		return
	}
	dst := c.panes[1-src.index]
	if !clone && len(src.docs) == 1 && src.index == 0 && len(dst.docs) == 0 {
		// Moving the only document leaves the view empty: a new document takes its place
	}
	st := d.state(src.index)
	c.showSecondView(true)
	dstState := d.state(dst.index)
	dstState.SetPosition(st.Caret(), st.Anchor(), st.TopLine())
	dst.add(d)
	if !clone {
		src.remove(d)
		if len(src.docs) == 0 {
			if src.index == 1 {
				c.showSecondView(false)
			} else {
				c.newDocumentIn(src)
			}
		}
	}
	c.setActivePane(dst.index)
	dst.view.Focus()
}

// focusOtherView moves the focus to the other view
func (c *MainForm) focusOtherView() {
	other := 1 - c.active
	if c.panes[other].IsVisible() && len(c.panes[other].docs) > 0 {
		c.setActivePane(other)
		c.focusEditor()
	}
}

// ---- Window state ----

const (
	defaultWindowWidth  = 1200
	defaultWindowHeight = 820
)

// RestoreWindowState applies the saved window layout before the form is shown.
// Returns whether the window should be maximized once shown.
func (c *MainForm) RestoreWindowState(form *ui.Form) (maximized bool) {
	form.SetSize(defaultWindowWidth, defaultWindowHeight)
	state, ok := config.LoadWindowState()
	if !ok {
		return false
	}
	form.SetSize(state.Width, state.Height)
	if state.X != 0 || state.Y != 0 {
		form.Move(state.X, state.Y)
	}
	if state.SecondViewWidth > 0 {
		c.secondW = state.SecondViewWidth
	}
	if state.ResultsHeight > 0 {
		c.resultsH = state.ResultsHeight
		c.vertSplit.SetSecondSize(c.resultsH)
	}
	if state.SideWidth > 0 {
		c.sideW = state.SideWidth
		c.sideSplit.SetFirstSize(c.sideW)
	}
	if state.FuncListWidth > 0 {
		c.funcW = state.FuncListWidth
		c.funcSplit.SetSecondSize(c.funcW)
	}
	if state.Workspace != "" {
		if st, err := os.Stat(state.Workspace); err == nil && st.IsDir() {
			c.workspace.openFolder(state.Workspace)
		}
	}
	if state.FuncList {
		c.funcList.SetVisible(true)
	}
	if state.DocMap {
		c.docMap.SetVisible(true)
	}
	c.updateRightPanel()
	return state.Maximized
}

// SaveWindowState remembers the window layout for the next start
func (c *MainForm) SaveWindowState() {
	form := c.Form()
	state, _ := config.LoadWindowState()
	state.Maximized = form.IsMaximized()
	if !state.Maximized && !c.fullscreen {
		state.X, state.Y = form.Position()
		state.Width, state.Height = form.Size()
	}
	state.SecondViewWidth = c.secondW
	state.ResultsHeight = c.resultsH
	state.SideWidth = c.sideW
	state.FuncListWidth = c.funcW
	state.Workspace = ""
	if c.workspace.IsVisible() {
		state.Workspace = c.workspace.root
	}
	state.FuncList = c.funcList.IsVisible()
	state.DocMap = c.docMap.IsVisible()
	config.SaveWindowState(state)
}

// toggleFullScreen maximizes the window without the tool bar and the status bar
func (c *MainForm) toggleFullScreen() {
	c.fullscreen = !c.fullscreen
	form := c.Form()
	if c.fullscreen {
		form.Maximize()
		c.toolbar.SetVisible(false)
		c.status.SetVisible(false)
	} else {
		form.Restore()
		c.toolbar.SetVisible(c.settings().ShowToolbar)
		c.status.SetVisible(c.settings().ShowStatusBar)
	}
	c.updateLayout()
}

// Activate gives the focus to the editor once the window is shown
func (c *MainForm) Activate() {
	c.setActivePane(c.active)
	c.focusEditor()
}

// ApplyLanguage updates the texts that do not follow the language by themselves
func (c *MainForm) ApplyLanguage() {
	c.status.refresh()
	c.find.applyLanguage()
	c.results.applyLanguage()
	c.workspace.applyLanguage()
	c.refreshTitle()
}

// toggleFuncList shows or hides the function list
func (c *MainForm) toggleFuncList() {
	c.funcList.SetVisible(!c.funcList.IsVisible())
	if c.funcList.IsVisible() {
		c.funcList.version = -1
		c.funcList.refresh()
	}
	c.updateRightPanel()
}

// toggleDocMap shows or hides the document map
func (c *MainForm) toggleDocMap() {
	c.docMap.SetVisible(!c.docMap.IsVisible())
	c.docMap.m.SetView(c.view())
	c.updateRightPanel()
}

// updateRightPanel shows the right side panel when the function list or the map is shown
func (c *MainForm) updateRightPanel() {
	show := c.funcList.IsVisible() || c.docMap.IsVisible()
	if show && !c.rightSplit.IsVisible() {
		c.funcSplit.SetSecondSize(max(150, c.funcW))
	}
	c.rightSplit.SetVisible(show)
	c.updateLayout()
}

func (c *MainForm) formatSize(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	}
	return fmt.Sprintf("%d B", n)
}

func baseName(path string) string { return filepath.Base(path) }

// BuildMenuBar makes the main menu of the window
func (c *MainForm) BuildMenuBar() *ui.MenuBar { return c.buildMenuBar() }

// OpenFile opens a file of the command line and goes to the line (1-based, 0 - none)
func (c *MainForm) OpenFile(path string, line int) { c.openFile(path, line) }

// EnsureDocument opens a new document when none is open
func (c *MainForm) EnsureDocument() {
	if len(c.panes[0].docs) == 0 {
		c.newDocumentIn(c.panes[0])
	}
	c.status.refresh()
	c.toolbar.refresh()
}

// ApplyANSI sets the character set of the non-Unicode files from the settings
func ApplyANSI(s config.Settings) { applyANSI(s) }
