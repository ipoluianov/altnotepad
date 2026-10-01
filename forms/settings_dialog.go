package forms

import (
	"strings"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// ShowSettings opens the preferences
func (c *MainForm) ShowSettings() {
	c.ShowDialog(NewSettingsDialog(c.settings(), c.ApplySettings))
}

// SettingsDialog edits the preferences (config.Settings)
type SettingsDialog struct {
	ui.DialogContent
	s        config.Settings
	onAccept func(config.Settings)
	apply    []func(s *config.Settings)
}

// page is a page of the settings: rows of a label and a control
type settingsPage struct {
	panel *ui.Panel
	row   int
}

func newSettingsPage() *settingsPage {
	p := &settingsPage{panel: ui.NewPanel()}
	return p
}

// add adds a row: a control with a label on its left, or a check box alone
func (p *settingsPage) add(label string, w ui.Widgeter) {
	if label != "" {
		row := ui.NewPanel()
		row.SetPanelPadding(0)
		lbl := ui.NewLabel(label)
		lbl.SetMinWidth(290)
		lbl.SetMaxWidth(290)
		row.AddWidget(0, 0, lbl)
		row.AddWidget(0, 1, w)
		p.panel.AddWidget(p.row, 0, row)
	} else {
		p.panel.AddWidget(p.row, 0, w)
	}
	p.row++
}

func (p *settingsPage) finish() *ui.Panel {
	p.panel.AddWidget(p.row, 0, ui.NewVSpacer())
	return p.panel
}

func NewSettingsDialog(s config.Settings, onAccept func(config.Settings)) *SettingsDialog {
	var c SettingsDialog
	c.InitWidget()
	c.s = s
	c.onAccept = onAccept

	check := func(p *settingsPage, text string, val *bool) {
		chk := ui.NewCheckbox(text)
		chk.SetChecked(*val)
		chk.SetXExpandable(true)
		p.add("", chk)
		c.apply = append(c.apply, func(s *config.Settings) { *val = chk.Checked() })
	}
	number := func(p *settingsPage, label string, val *int, minV, maxV int) {
		n := ui.NewNumBox()
		n.SetMin(float64(minV))
		n.SetMax(float64(maxV))
		n.SetDecimals(0)
		n.SetValue(float64(*val))
		n.SetMaxWidth(160)
		p.add(label, n)
		c.apply = append(c.apply, func(s *config.Settings) { *val = int(n.Value()) })
	}
	choice := func(p *settingsPage, label string, val *string, items [][2]string) {
		cmb := ui.NewComboBox()
		cmb.SetSelectedIndex(0)
		for i, it := range items {
			cmb.AddItem(it[1], it[0])
			if it[0] == *val {
				cmb.SetSelectedIndex(i)
			}
		}
		if *val == "" {
			cmb.SetSelectedIndex(0)
		}
		p.add(label, cmb)
		c.apply = append(c.apply, func(s *config.Settings) { *val, _ = cmb.SelectedItemData().(string) })
	}
	text := func(p *settingsPage, label string, val *string) {
		t := ui.NewTextBox()
		t.SetText(*val)
		p.add(label, t)
		c.apply = append(c.apply, func(s *config.Settings) { *val = strings.TrimSpace(t.Text()) })
	}
	cs := &c.s

	// General
	general := newSettingsPage()
	langs := [][2]string{{"", T().LanguageSystem}}
	for _, l := range languages {
		langs = append(langs, [2]string{l.tag, l.name})
	}
	choice(general, T().Language, &cs.Language, langs)
	choice(general, T().Theme, &cs.Theme, [][2]string{{themeDark, T().ThemeDark}, {themeLight, T().ThemeLight}})
	schemes := [][2]string{{"", T().SchemeByTheme}}
	for _, sc := range editor.Schemes {
		schemes = append(schemes, [2]string{sc.ID, sc.Name})
	}
	choice(general, T().ColorScheme, &cs.Scheme, schemes)
	check(general, T().ShowToolbar, &cs.ShowToolbar)
	check(general, T().ShowStatusBar, &cs.ShowStatusBar)
	check(general, T().TabCloseButtons, &cs.TabCloseButtons)
	check(general, T().DoubleClickCloses, &cs.DoubleClickCloses)
	check(general, T().AlwaysOnTop, &cs.AlwaysOnTop)
	check(general, T().SingleInstance, &cs.SingleInstance)
	check(general, T().RememberSession, &cs.RememberSession)
	check(general, T().BackupSession, &cs.Backup)
	number(general, T().BackupEvery, &cs.BackupSeconds, 1, 600)
	number(general, T().MaxRecent, &cs.MaxRecent, 0, 50)

	// Editing
	editing := newSettingsPage()
	number(editing, T().TabSize, &cs.TabSize, 1, 16)
	check(editing, T().InsertSpaces, &cs.InsertSpaces)
	check(editing, T().AutoIndent, &cs.AutoIndent)
	check(editing, T().AutoClose, &cs.AutoClose)
	check(editing, T().AutoCompletion, &cs.AutoComplete)
	check(editing, T().SmartHome, &cs.SmartHome)
	check(editing, T().MultiEdit, &cs.MultiEdit)
	check(editing, T().ScrollPast, &cs.ScrollPast)
	check(editing, T().CopyLineNoSel, &cs.CopyLineNoSel)
	number(editing, T().CaretWidth, &cs.CaretWidth, 1, 4)
	number(editing, T().CaretBlink, &cs.CaretBlink, 0, 2000)
	number(editing, T().EdgeColumn, &cs.EdgeColumn, 0, 500)
	text(editing, T().WordChars, &cs.WordChars)

	// Display
	display := newSettingsPage()
	fontSize := int(cs.FontSize)
	number(display, T().FontSize, &fontSize, 6, 72)
	c.apply = append(c.apply, func(s *config.Settings) { s.FontSize = float64(fontSize) })
	fontRow := ui.NewPanel()
	fontRow.SetPanelPadding(0)
	fontFile := ui.NewTextBox()
	fontFile.SetText(cs.FontFile)
	fontFile.SetHint(T().BuiltInFont)
	browse := ui.NewButton("...")
	browse.SetOnClick(func() {
		c.Form().ShowOpenFileDialog(ui.OpenFileDialogOptions{Title: T().FontFile,
			Filters: []ui.FileDialogFilter{{DisplayName: T().Fonts, Patterns: []string{"*.ttf", "*.otf", "*.ttc"}}}},
			func(paths []string, err error) {
				if err == nil && len(paths) > 0 {
					fontFile.SetText(paths[0])
				}
			})
	})
	fontRow.AddWidget(0, 0, fontFile)
	fontRow.AddWidget(0, 1, browse)
	display.add(T().FontFile, fontRow)
	c.apply = append(c.apply, func(s *config.Settings) { s.FontFile = strings.TrimSpace(fontFile.Text()) })
	check(display, T().LineNumbers, &cs.LineNumbers)
	check(display, T().BookmarkMargin, &cs.BookmarkMargin)
	check(display, T().FoldMargin, &cs.FoldMargin)
	check(display, T().CurrentLine, &cs.CurrentLine)
	check(display, T().SmartHighlight, &cs.SmartHighlight)
	check(display, T().SmartMatchCase, &cs.SmartMatchCase)
	check(display, T().SmartWholeWord, &cs.SmartWholeWord)
	check(display, T().BraceMatch, &cs.BraceMatch)
	check(display, T().ChangeHistory, &cs.ChangeHistory)
	check(display, T().WrapSymbol, &cs.ShowWrapSymbol)

	// New document and files
	files := newSettingsPage()
	choice(files, T().NewDocEOL, &cs.NewDocEOL, [][2]string{{"", T().SystemDefault}, {"crlf", "Windows (CR LF)"}, {"lf", "Unix (LF)"}, {"cr", "Macintosh (CR)"}})
	encs := [][2]string{{"", "UTF-8"}, {"utf-8-bom", "UTF-8-BOM"}, {"ansi", "ANSI"}}
	choice(files, T().NewDocEncoding, &cs.NewDocEncoding, encs)
	ansi := [][2]string{{"", T().ByLanguage}}
	for _, ch := range editor.Charsets {
		if ch.Group != "" && !strings.HasPrefix(ch.ID, "utf") {
			ansi = append(ansi, [2]string{ch.ID, charsetGroupName(ch.Group) + ": " + ch.Name})
		}
	}
	choice(files, T().ANSICharset, &cs.ANSI, ansi)
	langList := [][2]string{{"", T().NormalText}}
	for _, l := range editor.Languages() {
		if l.ID != "text" {
			langList = append(langList, [2]string{l.ID, l.Name})
		}
	}
	choice(files, T().NewDocLanguage, &cs.NewDocLanguage, langList)
	number(files, T().LargeFileMB, &cs.LargeFileMB, 1, 100000)
	check(files, T().CheckFileChanges, &cs.CheckFileChanges)
	check(files, T().AutoReload, &cs.AutoReload)

	tabs := ui.NewTabWidget()
	tabs.AddPage(T().PageGeneral, general.finish())
	tabs.AddPage(T().PageEditing, editing.finish())
	tabs.AddPage(T().PageDisplay, display.finish())
	tabs.AddPage(T().PageFiles, files.finish())
	c.AddWidget(0, 0, tabs)
	buttons := ui.NewPanel()
	c.AddWidget(1, 0, buttons)
	ok := ui.NewButton(ui.UIText().OK)
	cancel := ui.NewButton(ui.UIText().Cancel)
	ok.SetOnClick(c.Accept)
	cancel.SetOnClick(func() { c.Form().Close() })
	dialogButtons(buttons, ok, cancel)

	c.OnDialogShow = func() {
		c.Form().SetTitle(T().MenuPreferences)
		c.Form().SetSize(680, 660)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(ok)
		c.Form().SetCancelButton(cancel)
	}
	return &c
}

func (c *SettingsDialog) Accept() {
	s := config.GetSettings()
	// The fields not in the dialog are kept
	c.s.Zoom = s.Zoom
	c.s.LangTabs = s.LangTabs
	for _, f := range c.apply {
		f(&c.s)
	}
	result := c.s
	if c.onAccept != nil {
		c.RunInParent(func() { c.onAccept(result) })
	}
	c.Form().Close()
}
