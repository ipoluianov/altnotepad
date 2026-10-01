package forms

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/nui/ui"
)

// Workspace is the side panel with the tree of a folder (Folder as Workspace):
// a double click opens a file
type Workspace struct {
	ui.Widget

	main  *MainForm
	root  string
	title *ui.Label
	tree  *ui.TreeView

	folderIcon, fileIcon image.Image
	refreshAt            time.Time
}

func NewWorkspace(main *MainForm) *Workspace {
	var c Workspace
	c.InitWidget()
	c.main = main
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	head := ui.NewPanel()
	head.SetPanelPadding(2)
	c.title = ui.NewLabel(T().Workspace)
	head.AddWidget(0, 0, c.title)
	head.AddWidget(0, 1, ui.NewHSpacer())
	head.AddWidget(0, 2, smallButton("sync", func() string { return T().Refresh }, c.refresh))
	head.AddWidget(0, 3, smallButton("collapse", func() string { return T().CollapseAll }, func() { c.tree.CollapseAll() }))
	head.AddWidget(0, 4, smallButton("x", func() string { return T().Close }, func() {
		c.SetVisible(false)
		c.main.updateLayout()
	}))
	c.AddWidget(0, 0, head)

	c.tree = ui.NewTreeView()
	c.tree.SetColumnCount(1)
	c.tree.SetHeaderVisible(false)
	c.tree.SetStretchLastColumn(true)
	c.tree.SetOnExpand(c.expand)
	c.tree.SetOnNodeActivated(func(n *ui.TreeNode) {
		path, _ := n.Data().(string)
		if path == "" {
			return
		}
		if st, err := os.Stat(path); err == nil && !st.IsDir() {
			c.main.openFile(path, 0)
		} else {
			n.Toggle()
		}
	})
	c.tree.SetYExpandable(true)
	c.AddWidget(1, 0, c.tree)
	setIcon("folder-16", func(img image.Image) { c.folderIcon = img })
	setIcon("file-16", func(img image.Image) { c.fileIcon = img })
	c.AddTimer(500, func() {
		if !c.refreshAt.IsZero() && time.Now().After(c.refreshAt) {
			c.refreshAt = time.Time{}
			c.refresh()
		}
	})
	return &c
}

func (c *Workspace) applyLanguage() {
	if c.root == "" {
		c.title.SetText(T().Workspace)
	}
}

// chooseFolder asks for the folder to show
func (c *Workspace) chooseFolder() {
	c.main.Form().ShowSelectDirectoryDialog(ui.SelectDirectoryDialogOptions{Title: T().MenuOpenFolderWorkspace, DefaultDirectory: c.main.startDir()},
		func(path string, err error) {
			if errors.Is(err, ui.ErrNoFileDialog) {
				c.main.askPath(T().MenuOpenFolderWorkspace, c.main.startDir(), c.openFolder)
				return
			}
			if err == nil && path != "" {
				c.openFolder(path)
			}
		})
}

// openFolder shows the folder in the panel
func (c *Workspace) openFolder(path string) {
	c.root = path
	c.title.SetText(filepath.Base(path))
	c.title.SetTooltip(path)
	c.refresh()
	if !c.IsVisible() {
		c.SetVisible(true)
		c.main.sideSplit.SetFirstSize(max(150, c.main.sideW))
		c.main.updateLayout()
	}
	config.UpdateHistory(func(h *config.History) { h.LastDir = path })
}

// refresh reads the folder again, keeping the expanded folders
func (c *Workspace) refresh() {
	if c.root == "" {
		return
	}
	expanded := map[string]bool{}
	var walk func(nodes []*ui.TreeNode)
	walk = func(nodes []*ui.TreeNode) {
		for _, n := range nodes {
			if n.IsExpanded() {
				if p, ok := n.Data().(string); ok {
					expanded[p] = true
				}
				walk(n.Children())
			}
		}
	}
	walk(c.tree.Nodes())
	c.tree.Clear()
	c.fill(nil, c.root, expanded)
}

// fill adds the entries of the folder under the node
func (c *Workspace) fill(parent *ui.TreeNode, dir string, expanded map[string]bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.IsDir() != b.IsDir() {
			return a.IsDir()
		}
		return strings.ToLower(a.Name()) < strings.ToLower(b.Name())
	})
	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		n := c.tree.AddNode(parent, e.Name())
		n.SetData(path)
		n.SetReadOnly(0, true)
		if e.IsDir() {
			n.SetImage(c.folderIcon)
			n.SetHasChildren(true)
			if expanded[path] {
				c.fill(n, path, expanded)
				n.SetExpanded(true)
			}
		} else {
			n.SetImage(c.fileIcon)
		}
	}
}

// expand reads a folder when it is expanded the first time
func (c *Workspace) expand(n *ui.TreeNode) {
	if n.ChildCount() > 0 {
		return
	}
	if path, ok := n.Data().(string); ok {
		c.fill(n, path, nil)
	}
}

// refreshSoon reads the folder again a bit later, after a file changed
func (c *Workspace) refreshSoon() {
	if c.root != "" && c.IsVisible() {
		c.refreshAt = time.Now().Add(700 * time.Millisecond)
	}
}

// selectFile shows the file of the active document in the tree when it is there
func (c *Workspace) selectFile(path string) {
	if path == "" || c.root == "" || !c.IsVisible() || !strings.HasPrefix(path, c.root) {
		return
	}
	var find func(nodes []*ui.TreeNode) *ui.TreeNode
	find = func(nodes []*ui.TreeNode) *ui.TreeNode {
		for _, n := range nodes {
			p, _ := n.Data().(string)
			if p == path {
				return n
			}
			if strings.HasPrefix(path, p+string(filepath.Separator)) && n.IsExpanded() {
				return find(n.Children())
			}
		}
		return nil
	}
	if n := find(c.tree.Nodes()); n != nil {
		c.tree.SetCurrentNode(n)
		c.tree.ScrollToNode(n)
	}
}
