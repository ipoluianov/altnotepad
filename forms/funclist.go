package forms

import (
	"strings"
	"time"

	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// FuncListPanel is the side panel with the functions of the active
// document (View > Function List); a double click goes to one
type FuncListPanel struct {
	ui.Widget

	main    *MainForm
	filter  *ui.TextBox
	tree    *ui.TreeView
	doc     *editor.Document
	version int64
	items   []editor.FuncItem
	changed time.Time
}

func NewFuncListPanel(main *MainForm) *FuncListPanel {
	var c FuncListPanel
	c.InitWidget()
	c.main = main
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	head := ui.NewPanel()
	head.SetPanelPadding(2)
	title := ui.NewLabel(T().MenuFunctionList)
	title.SetTextFunc(func() string { return T().MenuFunctionList })
	head.AddWidget(0, 0, title)
	head.AddWidget(0, 1, ui.NewHSpacer())
	head.AddWidget(0, 2, smallButton("sync", func() string { return T().Refresh }, func() { c.version = -1; c.refresh() }))
	head.AddWidget(0, 3, smallButton("x", func() string { return T().Close }, func() {
		c.SetVisible(false)
		c.main.updateLayout()
	}))
	c.AddWidget(0, 0, head)
	c.filter = ui.NewTextBox()
	c.filter.SetHintFunc(func() string { return T().Filter })
	c.filter.SetOnTextChanged(func() { c.fill() })
	c.AddWidget(1, 0, c.filter)
	c.tree = ui.NewTreeView()
	c.tree.SetColumnCount(1)
	c.tree.SetHeaderVisible(false)
	c.tree.SetStretchLastColumn(true)
	c.tree.SetYExpandable(true)
	c.tree.SetOnNodeActivated(func(n *ui.TreeNode) {
		if line, ok := n.Data().(int); ok {
			v := c.main.view()
			v.GotoLine(line)
			v.Focus()
		}
	})
	c.AddWidget(2, 0, c.tree)
	// Follows the edits a second after the typing stops
	c.AddTimer(500, func() {
		if c.IsVisible() && c.doc != nil && c.doc.Version() != c.version && time.Since(c.main.view().LastEdit()) > time.Second {
			c.refresh()
		}
	})
	return &c
}

// refresh reads the functions of the active document when it changed
func (c *FuncListPanel) refresh() {
	if !c.IsVisible() {
		return
	}
	d := c.main.curDoc()
	if d == nil {
		return
	}
	if c.doc == d.doc && c.version == d.doc.Version() {
		return
	}
	c.doc = d.doc
	c.version = d.doc.Version()
	c.items = editor.FunctionList(d.doc)
	c.fill()
}

// fill shows the functions matching the filter
func (c *FuncListPanel) fill() {
	filter := strings.ToLower(strings.TrimSpace(c.filter.Text()))
	c.tree.Clear()
	var parents []*ui.TreeNode // the last node of each level
	for _, it := range c.items {
		if filter != "" && !strings.Contains(strings.ToLower(it.Name), filter) {
			continue
		}
		var parent *ui.TreeNode
		if filter == "" {
			for lvl := min(it.Level, len(parents)) - 1; lvl >= 0; lvl-- {
				if parents[lvl] != nil {
					parent = parents[lvl]
					break
				}
			}
		}
		n := c.tree.AddNode(parent, it.Name)
		n.SetData(it.Line)
		n.SetReadOnly(0, true)
		for len(parents) <= it.Level {
			parents = append(parents, nil)
		}
		parents[it.Level] = n
		parents = parents[:it.Level+1]
	}
	c.tree.ExpandAll()
}
