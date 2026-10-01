package forms

import (
	"image"

	"github.com/ipoluianov/nui/ui"
)

// Toolbar is the row of the buttons of the frequent commands
type Toolbar struct {
	ui.Widget

	main    *MainForm
	buttons map[string]*ui.ToolButton
	col     int
}

const toolButtonSize = 30

func NewToolbar(main *MainForm) *Toolbar {
	var c Toolbar
	c.InitWidget()
	c.main = main
	c.buttons = make(map[string]*ui.ToolButton)
	c.SetPanelPadding(2)
	c.SetCellPadding(1)
	c.SetVisible(main.settings().ShowToolbar)

	groups := [][]struct{ icon, cmd string }{
		{{"new", "file.new"}, {"open", "file.open"}, {"save", "file.save"}, {"saveall", "file.saveAll"}, {"close", "file.close"}},
		{{"cut", "edit.cut"}, {"copy", "edit.copy"}, {"paste", "edit.paste"}},
		{{"undo", "edit.undo"}, {"redo", "edit.redo"}},
		{{"find", "search.find"}, {"replace", "search.replace"}},
		{{"zoomin", "view.zoomIn"}, {"zoomout", "view.zoomOut"}},
		{{"wrap", "view.wrap"}, {"allchars", "view.allChars"}, {"indent", "view.indentGuides"}, {"bookmark", "bookmark.toggle"}},
		{{"split", "view.cloneOther"}, {"folder", "view.workspace"}, {"sync", "view.monitoring"}},
		{{"record", "macro.record"}, {"stop", "macro.stop"}, {"play", "macro.play"}},
	}
	for gi, g := range groups {
		if gi > 0 {
			space := ui.NewSpace()
			space.SetSize(8, 0)
			c.AddWidget(0, c.col, space)
			c.col++
		}
		for _, b := range g {
			c.addButton(b.icon, b.cmd)
		}
	}
	c.AddWidget(0, c.col, ui.NewHSpacer())
	return &c
}

func (c *Toolbar) addButton(icon, cmdID string) {
	cmd := c.main.commands[cmdID]
	btn := ui.NewToolButton(nil, "", func() {
		cmd.Run()
		c.refresh()
		if cmdID != "search.find" && cmdID != "search.replace" {
			c.main.focusEditor()
		}
	})
	btn.SetButtonSize(toolButtonSize, toolButtonSize)
	setIcon(icon+"-16", func(img image.Image) { btn.SetImage(img) })
	btn.SetTooltipFunc(cmd.menuText)
	c.buttons[cmdID] = btn
	c.AddWidget(0, c.col, btn)
	c.col++
}

// refresh shows the state of the commands on the buttons
func (c *Toolbar) refresh() {
	if !c.IsVisible() {
		return
	}
	m := c.main
	d := m.curDoc()
	v := m.view()
	set := func(id string, enabled bool) {
		if b := c.buttons[id]; b != nil {
			b.SetEnabled(enabled)
		}
	}
	check := func(id string) {
		if b := c.buttons[id]; b != nil {
			if cmd := m.commands[id]; cmd != nil && cmd.Checked != nil {
				b.SetChecked(cmd.Checked())
			}
		}
	}
	set("file.save", d != nil && (d.IsModified() || d.IsUntitled()))
	set("edit.undo", v.Document().CanUndo())
	set("edit.redo", v.Document().CanRedo())
	set("macro.stop", m.recording)
	set("macro.play", !m.recording && m.macro != nil && len(m.macro.Steps) > 0)
	for _, id := range []string{"view.wrap", "view.allChars", "view.indentGuides", "view.workspace", "view.monitoring"} {
		check(id)
	}
	if b := c.buttons["macro.record"]; b != nil {
		if m.recording {
			b.SetHighlight(colorRecording.get())
		} else {
			b.SetHighlight(nil)
		}
	}
}
