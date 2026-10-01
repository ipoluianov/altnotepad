package forms

import (
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// Pane is one of the two views: the tabs of its documents and the editor
type Pane struct {
	ui.Widget

	index int
	main  *MainForm
	tabs  *TabBar
	view  *editor.View
	docs  []*Doc
	cur   int
}

func newPane(main *MainForm, index int) *Pane {
	var c Pane
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	c.index = index
	c.main = main
	c.cur = -1
	c.tabs = NewTabBar()
	c.view = editor.NewView()
	c.AddWidget(0, 0, c.tabs)
	c.AddWidget(1, 0, c.view)

	c.tabs.OnSelect = func(i int) {
		c.activate(i)
		c.main.setActivePane(c.index)
		c.view.Focus()
	}
	c.tabs.OnClose = func(i int) {
		if i >= 0 && i < len(c.docs) {
			c.main.closeDoc(&c, c.docs[i], nil)
		}
	}
	c.tabs.OnMove = func(from, to int) {
		d := c.docs[from]
		c.docs = append(c.docs[:from], c.docs[from+1:]...)
		c.docs = append(c.docs[:to], append([]*Doc{d}, c.docs[to:]...)...)
		c.cur = to
		c.refreshTabs()
	}
	c.tabs.OnDoubleClick = func(i int) {
		switch {
		case i < 0:
			c.main.setActivePane(c.index)
			c.main.newDocument()
		case c.main.settings().DoubleClickCloses:
			c.main.closeDoc(&c, c.docs[i], nil)
		}
	}
	c.tabs.OnContextMenu = func(i int, x, y int) {
		if i >= 0 {
			c.activate(i)
			c.main.setActivePane(c.index)
			c.main.showTabMenu(x, y)
		}
	}
	c.view.OnFocused = func() { c.main.setActivePane(c.index) }
	c.view.OnSelectionChanged = func() { c.main.caretMoved(&c) }
	c.view.OnModified = func() { c.main.docEdited(&c) }
	c.view.OnContextMenu = func(x, y int) {
		c.main.setActivePane(c.index)
		c.main.showEditorMenu(x, y)
	}
	c.view.OnZoom = func(zoom int) { c.main.zoomChanged(zoom) }
	c.view.OnOvertypeChanged = func() { c.main.status.refresh() }
	return &c
}

// current returns the active document of the pane, nil if none
func (c *Pane) current() *Doc {
	if c.cur < 0 || c.cur >= len(c.docs) {
		return nil
	}
	return c.docs[c.cur]
}

func (c *Pane) indexOf(d *Doc) int {
	for i, x := range c.docs {
		if x == d {
			return i
		}
	}
	return -1
}

// add opens the document in the pane after the active one and activates it
func (c *Pane) add(d *Doc) {
	if i := c.indexOf(d); i >= 0 {
		c.activate(i)
		return
	}
	at := len(c.docs)
	c.docs = append(c.docs, d)
	c.main.applyDocOptions(d, c.index)
	c.activate(at)
}

// remove takes the document out of the pane, activating a neighbor
func (c *Pane) remove(d *Doc) {
	i := c.indexOf(d)
	if i < 0 {
		return
	}
	c.docs = append(c.docs[:i], c.docs[i+1:]...)
	d.release(c.index)
	if len(c.docs) == 0 {
		c.cur = -1
	} else {
		if i < c.cur {
			c.cur--
		}
		c.cur = max(0, min(c.cur, len(c.docs)-1))
	}
	if cur := c.current(); cur != nil {
		c.view.SetState(cur.state(c.index))
		c.main.applyViewOptions(c)
	}
	c.refreshTabs()
}

// activate shows the document of the tab
func (c *Pane) activate(i int) {
	if i < 0 || i >= len(c.docs) {
		return
	}
	c.cur = i
	d := c.docs[i]
	c.view.SetState(d.state(c.index))
	c.main.applyViewOptions(c)
	c.refreshTabs()
	c.main.docActivated(c)
}

// refreshTabs shows the names and states of the documents in the tabs
func (c *Pane) refreshTabs() {
	items := make([]TabItem, len(c.docs))
	for i, d := range c.docs {
		items[i] = TabItem{Title: d.Name(), Tooltip: d.FullName(), Modified: d.IsModified(),
			ReadOnly: d.doc.ReadOnly || d.doc.FileReadOnly, Monitoring: d.monitoring}
	}
	c.tabs.SetItems(items, c.cur)
}
