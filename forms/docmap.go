package forms

import (
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// DocMapPanel is the side panel with the document map of the active view
type DocMapPanel struct {
	ui.Widget
	main *MainForm
	m    *editor.DocMap
}

func NewDocMapPanel(main *MainForm) *DocMapPanel {
	var c DocMapPanel
	c.InitWidget()
	c.main = main
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	head := ui.NewPanel()
	head.SetPanelPadding(2)
	title := ui.NewLabel(T().MenuDocumentMap)
	title.SetTextFunc(func() string { return T().MenuDocumentMap })
	head.AddWidget(0, 0, title)
	head.AddWidget(0, 1, ui.NewHSpacer())
	head.AddWidget(0, 2, smallButton("x", func() string { return T().Close }, func() { main.toggleDocMap() }))
	c.AddWidget(0, 0, head)
	c.m = editor.NewDocMap()
	c.AddWidget(1, 0, c.m)
	return &c
}
