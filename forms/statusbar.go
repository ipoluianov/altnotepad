package forms

import (
	"fmt"
	"unicode/utf8"

	"github.com/ipoluianov/nui/ui"
)

// StatusBar shows the language, the size and the caret of the document,
// its line breaks and encoding. A click on the line breaks or the encoding
// opens the menu to change them.
type StatusBar struct {
	ui.Widget

	main    *MainForm
	lblMsg  *ui.Label
	lblLang *ui.Label
	lblSize *ui.Label
	lblPos  *ui.Label
	lblEOL  *ui.Label
	lblEnc  *ui.Label
	lblIns  *ui.Label
	message string

	eolMenu  *ui.ContextMenu
	encMenu  *ui.ContextMenu
	langMenu *ui.ContextMenu
}

func NewStatusBar(main *MainForm) *StatusBar {
	var c StatusBar
	c.InitWidget()
	c.main = main
	c.SetPanelPadding(3)
	c.SetCellPadding(0)
	c.SetVisible(main.settings().ShowStatusBar)

	c.lblLang = ui.NewLabel("")
	c.lblSize = ui.NewLabel("")
	c.lblPos = ui.NewLabel("")
	c.lblEOL = ui.NewLabel("")
	c.lblEnc = ui.NewLabel("")
	c.lblIns = ui.NewLabel("")
	c.lblMsg = ui.NewLabel("")
	c.lblMsg.SetForegroundColor(colorAccent.get())
	themeListeners = append(themeListeners, func() { c.lblMsg.SetForegroundColor(colorAccent.get()) })

	col := 0
	add := func(w ui.Widgeter, minW int) {
		if col > 0 {
			sep := ui.NewSpace()
			sep.SetSize(14, 0)
			c.AddWidget(0, col, sep)
			col++
		}
		if l, ok := w.(*ui.Label); ok && minW > 0 {
			l.SetMinWidth(minW)
		}
		c.AddWidget(0, col, w)
		col++
	}
	add(c.lblLang, 120)
	add(c.lblSize, 200)
	add(c.lblPos, 260)
	add(c.lblMsg, 0)
	c.AddWidget(0, col, ui.NewHSpacer())
	col++
	add(c.lblEOL, 110)
	add(c.lblEnc, 110)
	add(c.lblIns, 40)

	c.eolMenu = ui.NewContextMenu(&c)
	main.addItems(c.eolMenu, "eol.0", "eol.1", "eol.2")
	c.eolMenu.SetOnShow(main.updateChecks)
	c.lblEOL.SetContextMenu(c.eolMenu)
	c.encMenu = ui.NewContextMenu(&c)
	main.addItems(c.encMenu, "enc.ansi", "enc.utf8", "enc.utf8bom", "enc.utf16be", "enc.utf16le", "-",
		"conv.ansi", "conv.utf8", "conv.utf8bom", "conv.utf16be", "conv.utf16le")
	c.encMenu.SetOnShow(main.updateChecks)
	c.lblEnc.SetContextMenu(c.encMenu)

	click := func(lbl *ui.Label, menu *ui.ContextMenu) {
		lbl.SetMouseCursor(ui.MouseCursorPointer)
		lbl.SetOnMouseDown(func(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
			wx, wy := lbl.RectClientAreaOnWindow()
			menu.ShowMenu(wx+x, wy+y)
			return true
		})
	}
	click(c.lblEOL, c.eolMenu)
	click(c.lblEnc, c.encMenu)
	c.lblIns.SetMouseCursor(ui.MouseCursorPointer)
	c.lblIns.SetOnMouseDown(func(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
		main.view().Exec("toggle-overtype", "")
		return true
	})
	c.lblPos.SetMouseCursor(ui.MouseCursorPointer)
	c.lblPos.SetOnMouseDown(func(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
		main.gotoDialog()
		return true
	})
	return &c
}

// setMessage shows a message, e.g. the progress of loading a file
func (c *StatusBar) setMessage(msg string) {
	c.message = msg
	c.lblMsg.SetText(msg)
}

// refresh shows the state of the active document
func (c *StatusBar) refresh() {
	d := c.main.curDoc()
	if d == nil {
		return
	}
	lang := T().NormalText
	if l := d.doc.Language(); l != nil && l.ID != "text" {
		lang = l.Name
	}
	if d.doc.IsLarge() {
		lang += " " + T().LargeFileMark
	}
	c.lblLang.SetText(lang)
	c.lblSize.SetText(T().StatusSize(d.doc.Len(), d.doc.LineCount()))
	c.lblEOL.SetText(eolName(d))
	c.lblEnc.SetText(d.doc.Encoding.Name())
	if c.main.view().IsOvertype() {
		c.lblIns.SetText("OVR")
	} else {
		c.lblIns.SetText("INS")
	}
	c.refreshCaret()
}

func eolName(d *Doc) string {
	switch d.doc.EOL {
	case 0:
		return "Windows (CR LF)"
	case 2:
		return "Macintosh (CR)"
	}
	return "Unix (LF)"
}

// refreshCaret shows the position of the caret and the selection
func (c *StatusBar) refreshCaret() {
	v := c.main.view()
	line, col := v.CaretLineCol()
	sels := v.Selections()
	selChars, selLines := 0, 0
	doc := v.Document()
	for _, s := range sels {
		if s.Start() == s.End() {
			continue
		}
		// Characters for the usual selections, bytes for the huge ones
		if n := s.End() - s.Start(); n <= 4<<20 {
			selChars += utf8.RuneCount(doc.Buffer().View(s.Start(), s.End()))
		} else {
			selChars += n
		}
		selLines += doc.LineOfOffset(s.End()) - doc.LineOfOffset(s.Start()) + 1
	}
	text := T().StatusPos(line, col, v.Caret()+1)
	if selChars > 0 {
		if len(sels) > 1 {
			text += fmt.Sprintf("  %s %d | %d", T().StatusSel, selChars, len(sels))
		} else {
			text += fmt.Sprintf("  %s %d | %d", T().StatusSel, selChars, selLines)
		}
	}
	c.lblPos.SetText(text)
}
