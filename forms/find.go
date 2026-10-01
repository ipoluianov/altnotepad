package forms

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// The modes of the search panel, like the tabs of the Find dialog of Notepad++
const (
	findModeFind = iota
	findModeReplace
	findModeFiles
	findModeMark
)

// FindPanel is the search panel at the bottom of the window: find,
// replace, find in files and mark. It stays open while the text is edited.
type FindPanel struct {
	ui.Widget

	main *MainForm
	mode int

	modeButtons [4]*ui.Button
	lblStatus   *ui.Label

	lblReplace *ui.Label
	rowFiles   *ui.Panel

	cmbFind    *ui.EditableComboBox
	cmbReplace *ui.EditableComboBox
	cmbFilters *ui.EditableComboBox
	cmbDir     *ui.EditableComboBox

	chkWholeWord *ui.Checkbox
	chkMatchCase *ui.Checkbox
	chkWrap      *ui.Checkbox
	chkInSel     *ui.Checkbox
	chkBackward  *ui.Checkbox
	chkDotAll    *ui.Checkbox
	chkSubdirs   *ui.Checkbox
	chkHidden    *ui.Checkbox
	chkBookmark  *ui.Checkbox
	chkPurge     *ui.Checkbox
	rbNormal     *ui.RadioButton
	rbExtended   *ui.RadioButton
	rbRegex      *ui.RadioButton

	findButtons    []*ui.Button
	replaceButtons []*ui.Button
	filesButtons   []*ui.Button
	markButtons    []*ui.Button

	// incStart is where the incremental search started
	incStart int
	// fits fit the widths of the check boxes to their texts
	fits []func()
}

func NewFindPanel(main *MainForm) *FindPanel {
	var c FindPanel
	c.InitWidget()
	c.main = main
	c.SetPanelPadding(4)
	c.SetCellPadding(3)
	h := config.GetHistory()

	// The modes and the status
	top := ui.NewPanel()
	top.SetPanelPadding(0)
	names := []func() string{func() string { return T().FindTab }, func() string { return T().ReplaceTab },
		func() string { return T().FindInFilesTab }, func() string { return T().MarkTab }}
	for i := range c.modeButtons {
		mode := i
		b := ui.NewButton(names[i]())
		b.SetTextFunc(names[i])
		b.SetOnClick(func() { c.setMode(mode); c.cmbFind.Focus() })
		c.modeButtons[i] = b
		top.AddWidget(0, i, b)
	}
	c.lblStatus = ui.NewLabel("")
	top.AddWidget(0, 4, c.lblStatus)
	top.AddWidget(0, 5, ui.NewHSpacer())
	top.AddWidget(0, 6, smallButton("x", func() string { return T().Close + " (Esc)" }, c.close))
	c.AddWidget(0, 0, top)

	// Find what, replace with
	grid := ui.NewPanel()
	grid.SetPanelPadding(0)
	c.cmbFind = ui.NewEditableComboBox()
	c.cmbFind.SetItems(h.FindTexts)
	c.cmbFind.SetOnAccept(func(string) { c.defaultAction() })
	c.cmbFind.SetOnTextChanged(c.incremental)
	c.cmbFind.TextBox().SetOnTextBoxKeyDown(func() {
		ev := ui.CurrentEvent().Parameter.(*ui.EventTextboxKeyDown)
		if ev.Key == ui.KeyEnter && ev.Mods.Shift {
			c.findNext(true)
			ev.Processed = true
		}
	})
	lblFind := c.label(func() string { return T().FindWhat })
	grid.AddWidget(0, 0, lblFind)
	grid.AddWidget(0, 1, c.cmbFind)
	c.cmbReplace = ui.NewEditableComboBox()
	c.cmbReplace.SetItems(h.ReplaceTexts)
	c.cmbReplace.SetOnAccept(func(string) { c.replaceOne() })
	c.lblReplace = c.label(func() string { return T().ReplaceWith })
	grid.AddWidget(1, 0, c.lblReplace)
	grid.AddWidget(1, 1, c.cmbReplace)
	c.AddWidget(1, 0, grid)

	// Find in files: the filters and the folder
	c.rowFiles = ui.NewPanel()
	c.rowFiles.SetPanelPadding(0)
	c.cmbFilters = ui.NewEditableComboBox()
	c.cmbFilters.SetItems(h.FindFilters)
	c.cmbFilters.SetText("*")
	if len(h.FindFilters) > 0 {
		c.cmbFilters.SetText(h.FindFilters[0])
	}
	c.cmbFilters.SetMaxWidth(260)
	c.cmbDir = ui.NewEditableComboBox()
	c.cmbDir.SetItems(h.FindDirs)
	if len(h.FindDirs) > 0 {
		c.cmbDir.SetText(h.FindDirs[0])
	}
	browse := ui.NewButton("...")
	browse.SetOnClick(c.browseDir)
	useCur := ui.NewButton("")
	useCur.SetTextFunc(func() string { return T().CurrentFolder })
	useCur.SetOnClick(func() {
		if p := c.main.curPath(); p != "" {
			c.cmbDir.SetText(dirOf(p))
		}
	})
	c.rowFiles.AddWidget(0, 0, c.label(func() string { return T().Filters }))
	c.rowFiles.AddWidget(0, 1, c.cmbFilters)
	c.rowFiles.AddWidget(0, 2, c.label(func() string { return T().Directory }))
	c.rowFiles.AddWidget(0, 3, c.cmbDir)
	c.rowFiles.AddWidget(0, 4, browse)
	c.rowFiles.AddWidget(0, 5, useCur)
	c.AddWidget(2, 0, c.rowFiles)

	// Options
	opts := ui.NewPanel()
	opts.SetPanelPadding(0)
	col := 0
	check := func(text func() string, checked bool) *ui.Checkbox {
		chk := ui.NewCheckbox(text())
		chk.SetTextFunc(text)
		chk.SetChecked(checked)
		c.fits = append(c.fits, func() { fitWidth(chk, text()) })
		opts.AddWidget(0, col, chk)
		col++
		return chk
	}
	c.chkMatchCase = check(func() string { return T().MatchCase }, h.FindMatchCase)
	c.chkWholeWord = check(func() string { return T().WholeWord }, h.FindWholeWord)
	c.chkWrap = check(func() string { return T().WrapAround }, h.FindWrap)
	c.chkBackward = check(func() string { return T().Backward }, false)
	c.chkInSel = check(func() string { return T().InSelection }, false)
	c.chkSubdirs = check(func() string { return T().InSubfolders }, h.FindSubdirs)
	c.chkHidden = check(func() string { return T().InHidden }, h.FindHidden)
	c.chkBookmark = check(func() string { return T().BookmarkLine }, false)
	c.chkPurge = check(func() string { return T().PurgeEach }, true)
	opts.AddWidget(0, col, ui.NewHSpacer())
	c.AddWidget(3, 0, opts)

	modes := ui.NewPanel()
	modes.SetPanelPadding(0)
	modes.AddWidget(0, 0, c.label(func() string { return T().SearchMode }))
	radio := func(i int, text func() string) *ui.RadioButton {
		rb := ui.NewRadioButton(text())
		rb.SetTextFunc(text)
		c.fits = append(c.fits, func() { fitWidth(rb, text()) })
		modes.AddWidget(0, i, rb)
		return rb
	}
	c.rbNormal = radio(1, func() string { return T().ModeNormal })
	c.rbExtended = radio(2, func() string { return T().ModeExtended })
	c.rbRegex = radio(3, func() string { return T().ModeRegex })
	c.chkDotAll = ui.NewCheckbox(T().DotAll)
	c.chkDotAll.SetTextFunc(func() string { return T().DotAll })
	c.chkDotAll.SetChecked(h.FindDotAll)
	c.fits = append(c.fits, func() { fitWidth(c.chkDotAll, T().DotAll) })
	modes.AddWidget(0, 4, c.chkDotAll)
	modes.AddWidget(0, 5, ui.NewHSpacer())
	switch editor.SearchMode(h.FindMode) {
	case editor.SearchExtended:
		c.rbExtended.SetChecked(true)
	case editor.SearchRegex:
		c.rbRegex.SetChecked(true)
	default:
		c.rbNormal.SetChecked(true)
	}
	c.AddWidget(4, 0, modes)

	// The actions
	actions := ui.NewPanel()
	actions.SetPanelPadding(0)
	acol := 0
	button := func(text func() string, run func()) *ui.Button {
		b := ui.NewButton(text())
		b.SetTextFunc(text)
		b.SetOnClick(run)
		actions.AddWidget(0, acol, b)
		acol++
		return b
	}
	c.findButtons = []*ui.Button{
		button(func() string { return T().FindNext }, func() { c.findNext(false) }),
		button(func() string { return T().FindPrev }, func() { c.findNext(true) }),
		button(func() string { return T().Count }, c.count),
		button(func() string { return T().FindAllCurrent }, func() { c.findAllIn(false) }),
		button(func() string { return T().FindAllOpen }, func() { c.findAllIn(true) }),
	}
	c.replaceButtons = []*ui.Button{
		button(func() string { return T().Replace }, c.replaceOne),
		button(func() string { return T().ReplaceAll }, func() { c.replaceAll(false) }),
		button(func() string { return T().ReplaceAllOpen }, func() { c.replaceAll(true) }),
	}
	c.filesButtons = []*ui.Button{
		button(func() string { return T().FindAll }, func() { c.findInFiles(false) }),
		button(func() string { return T().ReplaceInFiles }, func() { c.findInFiles(true) }),
	}
	c.markButtons = []*ui.Button{
		button(func() string { return T().MarkAll }, c.markAll),
		button(func() string { return T().ClearMarks }, c.clearMarks),
		button(func() string { return T().CopyMarked }, c.copyMarked),
	}
	actions.AddWidget(0, acol, ui.NewHSpacer())
	c.AddWidget(5, 0, actions)

	c.applyLanguage()
	c.setMode(findModeFind)
	return &c
}

func (c *FindPanel) label(text func() string) *ui.Label {
	l := ui.NewLabel(text())
	l.SetTextFunc(text)
	return l
}

func (c *FindPanel) applyLanguage() {
	for _, f := range c.fits {
		f()
	}
}

func dirOf(path string) string {
	i := strings.LastIndexAny(path, `/\`)
	if i < 0 {
		return "."
	}
	return path[:i]
}

// setMode shows the rows and the buttons of the mode
func (c *FindPanel) setMode(mode int) {
	c.mode = mode
	for i, b := range c.modeButtons {
		if i == mode {
			b.SetRole("primary")
		} else {
			b.SetRole("")
		}
	}
	replace := mode == findModeReplace || mode == findModeFiles
	c.cmbReplace.SetVisible(replace)
	c.lblReplace.SetVisible(replace)
	c.rowFiles.SetVisible(mode == findModeFiles)
	show := func(list []*ui.Button, on bool) {
		for _, b := range list {
			b.SetVisible(on)
		}
	}
	show(c.findButtons, mode == findModeFind || mode == findModeReplace)
	show(c.replaceButtons, mode == findModeReplace)
	show(c.filesButtons, mode == findModeFiles)
	show(c.markButtons, mode == findModeMark)
	c.chkBackward.SetVisible(mode == findModeFind || mode == findModeReplace)
	c.chkWrap.SetVisible(mode != findModeFiles)
	c.chkInSel.SetVisible(mode != findModeFiles)
	c.chkSubdirs.SetVisible(mode == findModeFiles)
	c.chkHidden.SetVisible(mode == findModeFiles)
	c.chkBookmark.SetVisible(mode == findModeMark)
	c.chkPurge.SetVisible(mode == findModeMark)
	c.main.updateLayout()
}

// open shows the panel in the mode with the selected text (or the word at the caret) to find
func (c *FindPanel) open(mode int) {
	c.SetVisible(true)
	c.setMode(mode)
	v := c.main.view()
	sel := v.SelectedText()
	if sel != "" && !strings.Contains(sel, "\n") && len(sel) < 500 {
		c.cmbFind.SetText(sel)
	} else if sel == "" {
		a, b := v.WordAtCaret()
		if b > a && b-a < 200 {
			c.cmbFind.SetText(string(v.Document().Text(a, b)))
		}
	}
	if sel != "" && strings.Contains(sel, "\n") {
		c.chkInSel.SetChecked(true)
	} else {
		c.chkInSel.SetChecked(false)
	}
	if mode == findModeFiles && strings.TrimSpace(c.cmbDir.Text()) == "" {
		if p := c.main.curPath(); p != "" {
			c.cmbDir.SetText(dirOf(p))
		} else if c.main.workspace.root != "" {
			c.cmbDir.SetText(c.main.workspace.root)
		}
	}
	c.incStart = v.MainSelection().Start()
	c.lblStatus.SetText("")
	c.main.updateLayout()
	c.cmbFind.Focus()
	c.cmbFind.TextBox().SelectAllText()
}

// close hides the panel and returns to the editor
func (c *FindPanel) close() {
	c.SetVisible(false)
	c.main.updateLayout()
	c.main.focusEditor()
}

// hasFocus reports whether a field of the panel has the focus
func (c *FindPanel) hasFocus() bool {
	f := c.main.Form().FocusedWidget()
	if f == nil {
		return false
	}
	for _, w := range c.AllChildren() {
		if w.Id() == f.Id() {
			return true
		}
	}
	return false
}

func (c *FindPanel) defaultAction() {
	switch c.mode {
	case findModeFiles:
		c.findInFiles(false)
	case findModeMark:
		c.markAll()
	default:
		c.findNext(c.chkBackward.Checked())
	}
}

func (c *FindPanel) setStatus(text string, bad bool) {
	c.lblStatus.SetText(text)
	if bad {
		c.lblStatus.SetForegroundColor(ui.CurrentPalette().Error)
	} else {
		c.lblStatus.SetForegroundColor(colorAccent.get())
	}
}

// searchOptions returns the options chosen in the panel
func (c *FindPanel) searchOptions() editor.SearchOptions {
	mode := editor.SearchNormal
	switch {
	case c.rbExtended.Checked():
		mode = editor.SearchExtended
	case c.rbRegex.Checked():
		mode = editor.SearchRegex
	}
	return editor.SearchOptions{Mode: mode, MatchCase: c.chkMatchCase.Checked(), WholeWord: c.chkWholeWord.Checked(),
		DotAll: c.chkDotAll.Checked(), WordChars: c.main.settings().WordChars}
}

// remember saves the search texts and the options in the history
func (c *FindPanel) remember(withReplace bool) {
	find := c.cmbFind.Text()
	repl := c.cmbReplace.Text()
	o := c.searchOptions()
	config.UpdateHistory(func(h *config.History) {
		if find != "" {
			h.FindTexts = config.PushFront(h.FindTexts, find, 0)
		}
		if withReplace {
			h.ReplaceTexts = config.PushFront(h.ReplaceTexts, repl, 0)
		}
		h.FindMatchCase, h.FindWholeWord, h.FindWrap = o.MatchCase, o.WholeWord, c.chkWrap.Checked()
		h.FindMode, h.FindDotAll = int(o.Mode), o.DotAll
		h.FindSubdirs, h.FindHidden = c.chkSubdirs.Checked(), c.chkHidden.Checked()
	})
	hist := config.GetHistory()
	c.cmbFind.SetItems(hist.FindTexts)
	c.cmbReplace.SetItems(hist.ReplaceTexts)
}

// search is a search to make: what the panel says, kept for the macros
type search struct {
	Text      string
	Replace   string
	Opts      editor.SearchOptions
	Wrap      bool
	InSel     bool
	Backward  bool
	Bookmarks bool
	Purge     bool
}

func (c *FindPanel) current() search {
	return search{Text: c.cmbFind.Text(), Replace: c.cmbReplace.Text(), Opts: c.searchOptions(), Wrap: c.chkWrap.Checked(),
		InSel: c.chkInSel.Checked(), Backward: c.chkBackward.Checked(), Bookmarks: c.chkBookmark.Checked(), Purge: c.chkPurge.Checked()}
}

func (s search) matcher() (*editor.Matcher, error) {
	return editor.NewMatcher(s.Text, s.Opts)
}

func (s search) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func parseSearch(arg string) (search, bool) {
	var s search
	return s, json.Unmarshal([]byte(arg), &s) == nil
}

// ---- Find ----

// findNext selects the next match after the selection (before it when backward)
func (c *FindPanel) findNext(backward bool) {
	s := c.current()
	if s.Text == "" {
		c.open(findModeFind)
		return
	}
	c.remember(false)
	s.Backward = backward
	c.main.recordStep("app.find-next", s.String())
	c.doFindNext(c.main.view(), s)
}

func (c *FindPanel) doFindNext(v *editor.View, s search) bool {
	m, err := s.matcher()
	if err != nil {
		c.setStatus(err.Error(), true)
		return false
	}
	doc := v.Document()
	sel := v.MainSelection()
	from := sel.End()
	skip := -1
	if s.Backward {
		from = sel.Start()
	}
	if sel.Start() == sel.End() {
		skip = from
	}
	a, b, ok := m.FindNext(doc, from, s.Backward, 0, doc.Len(), skip)
	wrapped := false
	if !ok && s.Wrap {
		start := 0
		if s.Backward {
			start = doc.Len()
		}
		a, b, ok = m.FindNext(doc, start, s.Backward, 0, doc.Len(), -1)
		wrapped = ok
	}
	if !ok {
		c.setStatus(T().NotFound(s.Text), true)
		return false
	}
	v.SelectRange(a, b)
	if wrapped {
		c.setStatus(T().Wrapped, false)
	} else {
		c.setStatus("", false)
	}
	return true
}

// selectAndFind searches the selected text (or the word at the caret) (Ctrl+F3)
func (c *FindPanel) selectAndFind(backward bool) {
	v := c.main.view()
	text := v.SelectedText()
	if text == "" {
		a, b := v.WordAtCaret()
		if a == b {
			return
		}
		v.SetSelection(a, b)
		text = string(v.Document().Text(a, b))
	}
	c.cmbFind.SetText(text)
	c.rbNormal.SetChecked(true)
	s := c.current()
	s.Opts.Mode = editor.SearchNormal
	s.Backward = backward
	c.doFindNext(v, s)
}

// incremental selects the first match of the text typed, from where the search started
func (c *FindPanel) incremental() {
	if c.mode == findModeFiles || c.mode == findModeMark {
		return
	}
	s := c.current()
	v := c.main.view()
	if s.Text == "" {
		v.SetCaret(c.incStart)
		c.setStatus("", false)
		return
	}
	m, err := s.matcher()
	if err != nil {
		return
	}
	doc := v.Document()
	a, b, ok := m.FindNext(doc, c.incStart, false, 0, doc.Len(), -1)
	if !ok {
		a, b, ok = m.FindNext(doc, 0, false, 0, doc.Len(), -1)
	}
	if ok {
		v.SelectRange(a, b)
		c.setStatus("", false)
	} else {
		c.setStatus(T().NotFound(s.Text), true)
	}
}

// scope returns the part of the document to search: the selection or all of it
func scope(v *editor.View, inSel bool) (int, int) {
	if inSel && v.HasSelection() {
		s := v.MainSelection()
		return s.Start(), s.End()
	}
	return 0, v.Document().Len()
}

// count counts the matches in the document
func (c *FindPanel) count() {
	s := c.current()
	m, err := s.matcher()
	if err != nil {
		c.setStatus(err.Error(), true)
		return
	}
	c.remember(false)
	v := c.main.view()
	from, to := scope(v, s.InSel)
	n := 0
	m.FindAll(v.Document(), from, to, func(a, b int) bool { n++; return true })
	c.setStatus(T().CountResult(n), n == 0)
}

// ---- Replace ----

// replaceOne replaces the selected match and selects the next one
func (c *FindPanel) replaceOne() {
	s := c.current()
	c.remember(true)
	c.main.recordStep("app.replace", s.String())
	c.doReplaceOne(c.main.view(), s)
}

func (c *FindPanel) doReplaceOne(v *editor.View, s search) {
	m, err := s.matcher()
	if err != nil {
		c.setStatus(err.Error(), true)
		return
	}
	doc := v.Document()
	if doc.ReadOnly {
		c.setStatus(T().ReadOnlyDoc, true)
		return
	}
	sel := v.MainSelection()
	if sel.Start() != sel.End() {
		a, b, ok := m.FindNext(doc, sel.Start(), false, 0, doc.Len(), -1)
		if ok && a == sel.Start() && b == sel.End() {
			repl := editor.NewReplacement(s.Replace, s.Opts.Mode).Expand(m, doc, a, b)
			doc.BeginAction()
			doc.Replace(a, b-a, repl)
			doc.EndAction()
			v.SetCaret(a + len(repl))
		}
	}
	c.doFindNext(v, s)
}

// replaceAll replaces all the matches in the document (in all the open ones)
func (c *FindPanel) replaceAll(allOpen bool) {
	s := c.current()
	c.remember(true)
	if _, err := s.matcher(); err != nil {
		c.setStatus(err.Error(), true)
		return
	}
	if !allOpen {
		c.main.recordStep("app.replace-all", s.String())
		n := c.doReplaceAll(c.main.view(), s)
		c.setStatus(T().ReplacedCount(n), n == 0)
		return
	}
	total := 0
	for _, d := range c.main.docs {
		if d.doc.ReadOnly {
			continue
		}
		total += replaceInDoc(d.doc, s, 0, d.doc.Len())
	}
	c.main.refreshAllViews()
	c.setStatus(T().ReplacedCount(total), total == 0)
}

func (c *FindPanel) doReplaceAll(v *editor.View, s search) int {
	doc := v.Document()
	if doc.ReadOnly {
		c.setStatus(T().ReadOnlyDoc, true)
		return 0
	}
	from, to := scope(v, s.InSel)
	n := replaceInDoc(doc, s, from, to)
	v.Refresh()
	return n
}

// replaceInDoc replaces the matches between from and to as one undo step; returns their number
func replaceInDoc(doc *editor.Document, s search, from, to int) int {
	m, err := s.matcher()
	if err != nil {
		return 0
	}
	repl := editor.NewReplacement(s.Replace, s.Opts.Mode)
	var edits []editor.Edit
	m.FindAll(doc, from, to, func(a, b int) bool {
		edits = append(edits, editor.Edit{Pos: a, Len: b - a, Text: repl.Expand(m, doc, a, b)})
		return true
	})
	if len(edits) == 0 {
		return 0
	}
	doc.BeginAction()
	doc.ApplyEdits(edits)
	doc.EndAction()
	return len(edits)
}

// ---- Mark ----

// markAll marks the matches in the document (and bookmarks their lines)
func (c *FindPanel) markAll() {
	s := c.current()
	c.remember(false)
	c.main.recordStep("app.mark", s.String())
	n := c.doMark(c.main.view(), s)
	c.setStatus(T().MarkedCount(n), n == 0)
}

func (c *FindPanel) doMark(v *editor.View, s search) int {
	m, err := s.matcher()
	if err != nil {
		c.setStatus(err.Error(), true)
		return 0
	}
	doc := v.Document()
	if s.Purge {
		doc.Marks[editor.MarkStyleFind].Clear()
		if s.Bookmarks {
			doc.Bookmarks.Clear()
		}
	}
	from, to := scope(v, s.InSel)
	n := 0
	lastLine := -1
	m.FindAll(doc, from, to, func(a, b int) bool {
		n++
		if b > a {
			doc.Marks[editor.MarkStyleFind].Add(editor.Range{Start: a, End: b})
		}
		if s.Bookmarks {
			if l := doc.LineOfOffset(a); l != lastLine {
				doc.Bookmarks.Add(l)
				lastLine = l
			}
		}
		return true
	})
	v.Refresh()
	return n
}

func (c *FindPanel) clearMarks() {
	v := c.main.view()
	v.Document().Marks[editor.MarkStyleFind].Clear()
	v.Refresh()
	c.setStatus("", false)
}

// copyMarked copies the marked texts, one per line
func (c *FindPanel) copyMarked() {
	doc := c.main.view().Document()
	var b strings.Builder
	for _, r := range doc.Marks[editor.MarkStyleFind].Ranges() {
		b.Write(doc.Text(r.Start, r.End))
		b.WriteString("\n")
	}
	ui.ClipboardSetText(b.String())
}

// ---- Find All ----

// findAllIn lists the matches of the current document (of all the open ones) in the results
func (c *FindPanel) findAllIn(allOpen bool) {
	s := c.current()
	m, err := s.matcher()
	if err != nil {
		c.setStatus(err.Error(), true)
		return
	}
	c.remember(false)
	var docs []*Doc
	if allOpen {
		docs = c.main.docs
	} else if d := c.main.curDoc(); d != nil {
		docs = []*Doc{d}
	}
	res := &searchResult{text: s.Text, opts: s.Opts}
	for _, d := range docs {
		from, to := 0, d.doc.Len()
		if !allOpen && s.InSel {
			from, to = scope(c.main.view(), true)
		}
		fr := collectMatches(d.doc, m, from, to)
		if len(fr.hits) > 0 {
			fr.path = d.FullName()
			fr.doc = d
			res.files = append(res.files, fr)
		}
		res.searched++
	}
	c.main.results.add(res)
	c.setStatus(T().FoundCount(res.total()), res.total() == 0)
}

func (c *FindPanel) browseDir() {
	c.main.Form().ShowSelectDirectoryDialog(ui.SelectDirectoryDialogOptions{Title: T().Directory, DefaultDirectory: c.cmbDir.Text()},
		func(path string, err error) {
			if err == nil && path != "" {
				c.cmbDir.SetText(path)
			}
		})
}

// runMacroSearch runs the search steps of a macro
func (c *FindPanel) runMacroSearch(v *editor.View, cmd, arg string) bool {
	s, ok := parseSearch(arg)
	if !ok {
		return false
	}
	switch cmd {
	case "app.find-next":
		c.doFindNext(v, s)
	case "app.replace":
		c.doReplaceOne(v, s)
	case "app.replace-all":
		c.doReplaceAll(v, s)
	case "app.mark":
		c.doMark(v, s)
	default:
		return false
	}
	return true
}

func init() {
	editor.OnMacroCommand = func(v *editor.View, cmd, arg string) bool {
		if mainForm == nil {
			return false
		}
		return mainForm.find.runMacroSearch(v, cmd, arg)
	}
}

// refreshAllViews repaints the views, e.g. after their documents changed
func (c *MainForm) refreshAllViews() {
	for _, p := range c.panes {
		p.view.Refresh()
		p.refreshTabs()
	}
	c.status.refresh()
}

var _ = fmt.Sprint
