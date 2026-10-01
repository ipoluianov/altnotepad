package forms

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// Files larger than this are read in the background
const asyncLoadSize = 16 << 20

// nextUntitledTitle returns the name of a new document: "new N" with the
// smallest N not in use
func (c *MainForm) nextUntitledTitle() string {
	used := map[string]bool{}
	for _, d := range c.docs {
		if d.IsUntitled() {
			used[d.title] = true
		}
	}
	for n := 1; ; n++ {
		t := fmt.Sprintf("%s %d", T().NewDocName, n)
		if !used[t] {
			return t
		}
	}
}

// makeUntitled makes a new empty document with the settings of new documents
func (c *MainForm) makeUntitled() *Doc {
	s := c.settings()
	d := newDoc()
	d.title = c.nextUntitledTitle()
	d.doc.EOL = defaultEOL(s)
	d.doc.Encoding = defaultEncoding(s)
	if s.NewDocLanguage != "" {
		d.doc.SetLanguage(editor.LanguageByID(s.NewDocLanguage))
	}
	c.docs = append(c.docs, d)
	return d
}

// newDocument opens a new document in the active view
func (c *MainForm) newDocument() {
	c.newDocumentIn(c.activePane())
	c.focusEditor()
}

func (c *MainForm) newDocumentIn(p *Pane) *Doc {
	d := c.makeUntitled()
	p.add(d)
	return d
}

// findOpenDoc returns the document of the file if it is open
func (c *MainForm) findOpenDoc(path string) *Doc {
	abs := absPath(path)
	for _, d := range c.docs {
		if d.Path() != "" && samePath(d.Path(), abs) {
			return d
		}
	}
	return nil
}

func absPath(path string) string {
	if a, err := filepath.Abs(path); err == nil {
		return a
	}
	return path
}

func samePath(a, b string) bool {
	if os.PathSeparator == '\\' {
		return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

// showDoc activates the document in the active view (or where it is open)
func (c *MainForm) showDoc(d *Doc) {
	p := c.activePane()
	if p.indexOf(d) < 0 {
		other := c.panes[1-p.index]
		if other.IsVisible() && other.indexOf(d) >= 0 {
			p = other
		}
	}
	p.add(d)
	c.setActivePane(p.index)
	p.view.Focus()
}

// OpenFiles opens the files (activates those already open); a folder opens
// as the workspace
func (c *MainForm) OpenFiles(paths []string) {
	for _, path := range paths {
		c.openFile(path, 0)
	}
}

// openFile opens a file and goes to the line (1-based; 0 - where it was)
func (c *MainForm) openFile(path string, line int) *Doc {
	path = absPath(path)
	if st, err := os.Stat(path); err == nil && st.IsDir() {
		c.workspace.openFolder(path)
		return nil
	}
	if d := c.findOpenDoc(path); d != nil {
		c.showDoc(d)
		if line > 0 {
			c.view().GotoLine(line - 1)
		}
		return d
	}
	st, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Notepad++ offers to create it
			ui.ShowQuestionMessageBoxYesNo(c, T().OpenTitle, T().CreateFileAsk(path), func() {
				if f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					f.Close()
					c.openFile(path, line)
				} else {
					c.showError(err)
				}
			}, nil)
			return nil
		}
		c.showError(err)
		return nil
	}
	d := newDoc()
	if st.Size() > asyncLoadSize {
		c.loadAsync(d, path, line)
		return d
	}
	if err := loadInto(d, path, editor.Encoding{}); err != nil {
		c.showError(err)
		return nil
	}
	c.addOpened(d, line)
	return d
}

// loadAsync reads a large file in the background, showing the progress
func (c *MainForm) loadAsync(d *Doc, path string, line int) {
	c.status.setMessage(T().Loading(filepath.Base(path), 0))
	form := c.Form()
	go func() {
		last := -1
		lf, err := editor.LoadFile(path, editor.Encoding{}, func(done, total int64) {
			pct := int(done * 100 / max(1, total))
			if pct != last {
				last = pct
				form.Invoke(func() { c.status.setMessage(T().Loading(filepath.Base(path), pct)) })
			}
		})
		form.Invoke(func() {
			c.status.setMessage("")
			if err != nil {
				c.showError(err)
				return
			}
			if existing := c.findOpenDoc(path); existing != nil {
				c.showDoc(existing)
				return
			}
			applyLoaded(d, path, lf)
			c.addOpened(d, line)
			if d.doc.IsLarge() {
				c.toast(T().LargeFileMode)
			}
		})
	}()
}

// addOpened adds a loaded document to the active view; an empty new
// document there is replaced by it
func (c *MainForm) addOpened(d *Doc, line int) {
	p := c.activePane()
	var replace *Doc
	if len(p.docs) == 1 && p.docs[0].isEmptyUntitled() && !p.docs[0].inView(1-p.index) {
		replace = p.docs[0]
	}
	c.docs = append(c.docs, d)
	p.add(d)
	if replace != nil {
		p.remove(replace)
		c.forgetDoc(replace)
	}
	c.addRecent(d.Path())
	if line > 0 {
		p.view.GotoLine(line - 1)
	}
	c.setActivePane(p.index)
	p.view.Focus()
	c.status.refresh()
}

// forgetDoc drops a document closed in all the views
func (c *MainForm) forgetDoc(d *Doc) {
	for i, x := range c.docs {
		if x == d {
			c.docs = append(c.docs[:i], c.docs[i+1:]...)
			break
		}
	}
	d.release(0)
	d.release(1)
	c.removeBackup(d)
}

func (c *MainForm) addRecent(path string) {
	if path == "" {
		return
	}
	max := c.settings().MaxRecent
	config.UpdateHistory(func(h *config.History) {
		h.RecentFiles = config.PushFront(h.RecentFiles, path, max)
		h.LastDir = filepath.Dir(path)
	})
}

// fileFilters are the filters of the file dialogs
func fileFilters() []ui.FileDialogFilter {
	return []ui.FileDialogFilter{
		{DisplayName: T().AllFiles, Patterns: []string{"*"}},
		{DisplayName: T().TextFiles, Patterns: []string{"*.txt", "*.log", "*.md", "*.ini", "*.cfg", "*.conf"}},
	}
}

// startDir returns the folder the file dialogs open in
func (c *MainForm) startDir() string {
	if p := c.curPath(); p != "" {
		return filepath.Dir(p)
	}
	if d := c.workspace.root; d != "" {
		return d
	}
	if h := config.GetHistory(); h.LastDir != "" {
		return h.LastDir
	}
	home, _ := os.UserHomeDir()
	return home
}

// openDialog asks for files to open
func (c *MainForm) openDialog() {
	c.Form().ShowOpenFileDialog(ui.OpenFileDialogOptions{Title: T().OpenTitle, AllowMultiple: true,
		DefaultDirectory: c.startDir(), Filters: fileFilters()},
		func(paths []string, err error) {
			if errors.Is(err, ui.ErrNoFileDialog) {
				c.askPath(T().OpenTitle, c.startDir()+string(os.PathSeparator), func(path string) { c.openFile(path, 0) })
				return
			}
			if err != nil {
				c.showError(err)
				return
			}
			c.OpenFiles(paths)
		})
}

// askPath asks for a path typed in, when the system has no file dialog
// (Linux without zenity and kdialog)
func (c *MainForm) askPath(title, value string, onAccept func(path string)) {
	c.ShowDialog(NewInputDialog(title, T().PathLabel, value, func(path string) {
		path = strings.TrimSpace(path)
		if strings.HasPrefix(path, "~/") {
			if home, err := os.UserHomeDir(); err == nil {
				path = filepath.Join(home, path[2:])
			}
		}
		if path != "" {
			onAccept(absPath(path))
		}
	}))
}

// ---- Saving ----

// save saves the document; a new one asks for a name. done gets whether it was saved.
func (c *MainForm) save(d *Doc, done func(ok bool)) {
	if d == nil {
		return
	}
	if d.IsUntitled() {
		c.saveAs(d, false, done)
		return
	}
	c.writeDoc(d, d.Path(), false, done)
}

// writeDoc writes the document to the file; a copy does not change the document
func (c *MainForm) writeDoc(d *Doc, path string, copyOnly bool, done func(ok bool)) {
	finish := func(ok bool) {
		if done != nil {
			done(ok)
		}
	}
	write := func(replace bool) error {
		return editor.SaveFile(path, d.doc.Buffer(), d.doc.Encoding, d.doc.EOL, replace)
	}
	err := write(false)
	if errors.Is(err, editor.ErrUnencodable) {
		ui.ShowQuestionMessageBoxYesNo(c, T().SaveTitle, T().Unencodable(d.doc.Encoding.Name()), func() {
			if err := write(true); err != nil {
				c.showError(err)
				finish(false)
				return
			}
			c.saved(d, path, copyOnly)
			finish(true)
		}, func() { finish(false) })
		return
	}
	if err != nil {
		c.showError(fmt.Errorf("%s: %w", path, err))
		finish(false)
		return
	}
	c.saved(d, path, copyOnly)
	finish(true)
}

// saved updates the document after it was written to the file
func (c *MainForm) saved(d *Doc, path string, copyOnly bool) {
	if copyOnly {
		c.addRecent(path)
		return
	}
	renamed := !samePath(d.Path(), path)
	d.doc.Path = path
	d.doc.SetSavePoint()
	d.statFile()
	if renamed && !d.langSet {
		d.doc.SetLanguage(editor.LanguageForFile(path, d.doc.Text(0, min(d.doc.Len(), 1024))))
		for _, p := range c.panes {
			if p.current() == d {
				c.applyViewOptions(p)
				p.view.Refresh()
			}
		}
	}
	c.removeBackup(d)
	c.addRecent(path)
	for _, p := range c.panes {
		p.refreshTabs()
	}
	c.status.refresh()
	c.refreshTitle()
	c.workspace.refreshSoon()
}

// saveAs asks for a file name and saves the document; a copy leaves the document as it is
func (c *MainForm) saveAs(d *Doc, copyOnly bool, done func(ok bool)) {
	if d == nil {
		return
	}
	name := d.Name()
	if d.IsUntitled() {
		name += defaultExt(d)
	}
	title := T().SaveAsTitle
	if copyOnly {
		title = T().SaveCopyTitle
	}
	dir := c.startDir()
	if d.Path() != "" {
		dir = filepath.Dir(d.Path())
	}
	c.Form().ShowSaveFileDialog(ui.SaveFileDialogOptions{Title: title, DefaultDirectory: dir, DefaultFileName: name, Filters: fileFilters()},
		func(path string, err error) {
			if errors.Is(err, ui.ErrNoFileDialog) {
				c.askPath(title, filepath.Join(dir, name), func(path string) {
					if other := c.findOpenDoc(path); other != nil && other != d {
						c.showError(errors.New(T().AlreadyOpen(path)))
						return
					}
					c.writeDoc(d, path, copyOnly, done)
				})
				return
			}
			if err != nil {
				c.showError(err)
				if done != nil {
					done(false)
				}
				return
			}
			if path == "" {
				if done != nil {
					done(false)
				}
				return
			}
			if other := c.findOpenDoc(path); other != nil && other != d {
				c.showError(errors.New(T().AlreadyOpen(path)))
				if done != nil {
					done(false)
				}
				return
			}
			c.writeDoc(d, path, copyOnly, done)
		})
}

// defaultExt returns the extension of a new document's file by its language
func defaultExt(d *Doc) string {
	if l := d.doc.Language(); l != nil && len(l.Extensions) > 0 && l.ID != "text" {
		return "." + l.Extensions[0]
	}
	return ".txt"
}

// saveAll saves the modified documents, one after another
func (c *MainForm) saveAll() {
	var list []*Doc
	for _, d := range c.docs {
		if d.IsModified() {
			list = append(list, d)
		}
	}
	var next func(i int)
	next = func(i int) {
		if i >= len(list) {
			return
		}
		c.save(list[i], func(ok bool) { next(i + 1) })
	}
	next(0)
}

// ---- Closing ----

// closeDoc closes the document in the pane, asking to save it when it is
// closed in all the views; done gets whether it was closed
func (c *MainForm) closeDoc(p *Pane, d *Doc, done func(closed bool)) {
	finish := func(ok bool) {
		if done != nil {
			done(ok)
		}
	}
	other := c.panes[1-p.index]
	if other.indexOf(d) >= 0 || !d.IsModified() {
		c.doClose(p, d)
		finish(true)
		return
	}
	p.activate(p.indexOf(d))
	c.setActivePane(p.index)
	showSaveChanges(c, d.FullName(), func(answer int) {
		switch answer {
		case answerYes:
			c.save(d, func(ok bool) {
				if ok {
					c.doClose(p, d)
				}
				finish(ok)
			})
		case answerNo:
			c.doClose(p, d)
			finish(true)
		default:
			finish(false)
		}
	})
}

// doClose closes the document in the pane without asking
func (c *MainForm) doClose(p *Pane, d *Doc) {
	other := c.panes[1-p.index]
	p.remove(d)
	if other.indexOf(d) < 0 {
		if d.Path() != "" {
			c.closed = append(c.closed, d.Path())
			if len(c.closed) > 30 {
				c.closed = c.closed[1:]
			}
		}
		c.forgetDoc(d)
	}
	if len(p.docs) == 0 {
		if p.index == 1 {
			c.showSecondView(false)
		} else if len(other.docs) > 0 && other.IsVisible() && !c.quitting {
			// The documents of the second view move to the first one
			for _, x := range append([]*Doc(nil), other.docs...) {
				st := x.state(other.index)
				x.state(0).SetPosition(st.Caret(), st.Anchor(), st.TopLine())
				p.add(x)
				other.remove(x)
			}
			c.showSecondView(false)
		} else if !c.quitting {
			c.newDocumentIn(p)
		}
	}
	c.setActivePane(c.active)
	if !c.quitting {
		c.focusEditor()
	}
	c.status.refresh()
	c.refreshTitle()
}

// closeDocs closes the documents of the pane one after another, then calls done
func (c *MainForm) closeDocs(p *Pane, list []*Doc, done func(all bool)) {
	var next func(i int)
	next = func(i int) {
		if i >= len(list) {
			if done != nil {
				done(true)
			}
			return
		}
		c.closeDoc(p, list[i], func(ok bool) {
			if !ok {
				if done != nil {
					done(false)
				}
				return
			}
			next(i + 1)
		})
	}
	next(0)
}

// closeWhere closes the documents of the active pane the filter selects
func (c *MainForm) closeWhere(filter func(i int, d *Doc) bool) {
	p := c.activePane()
	var list []*Doc
	for i, d := range p.docs {
		if filter(i, d) {
			list = append(list, d)
		}
	}
	c.closeDocs(p, list, nil)
}

// closeAll closes all the documents of both views
func (c *MainForm) closeAll(done func(all bool)) {
	c.closeDocs(c.panes[1], append([]*Doc(nil), c.panes[1].docs...), func(all bool) {
		if !all {
			if done != nil {
				done(false)
			}
			return
		}
		c.closeDocs(c.panes[0], append([]*Doc(nil), c.panes[0].docs...), done)
	})
}

// restoreClosed reopens the file closed last
func (c *MainForm) restoreClosed() {
	for len(c.closed) > 0 {
		path := c.closed[len(c.closed)-1]
		c.closed = c.closed[:len(c.closed)-1]
		if c.findOpenDoc(path) == nil {
			if _, err := os.Stat(path); err == nil {
				c.openFile(path, 0)
				return
			}
		}
	}
}

// ---- Reload, rename ----

// reload reads the file of the document again, asking when it has changes
func (c *MainForm) reload(d *Doc, ask bool) {
	if d == nil || d.IsUntitled() {
		return
	}
	do := func() {
		keep := [2][3]int{}
		for i := range d.states {
			if st := d.states[i]; st != nil {
				keep[i] = [3]int{st.Caret(), st.Anchor(), st.TopLine()}
			}
		}
		lang := d.doc.Language()
		if err := loadInto(d, d.Path(), d.doc.Encoding); err != nil {
			c.showError(err)
			return
		}
		if d.langSet {
			d.doc.SetLanguage(lang)
		}
		for i := range d.states {
			if st := d.states[i]; st != nil {
				st.SetPosition(keep[i][0], keep[i][1], keep[i][2])
			}
		}
		for _, p := range c.panes {
			if p.current() == d {
				p.view.Refresh()
				if d.monitoring {
					p.view.GotoPos(d.doc.Len())
				}
			}
			p.refreshTabs()
		}
		c.status.refresh()
	}
	if ask && d.IsModified() {
		ui.ShowQuestionMessageBoxYesNo(c, T().ReloadTitle, T().ReloadAsk(d.Path()), do, nil)
		return
	}
	do()
}

// renameDoc renames the file of the document (a new document just gets another name)
func (c *MainForm) renameDoc(d *Doc) {
	if d == nil {
		return
	}
	c.ShowDialog(NewInputDialog(T().RenameTitle, T().NewName, d.Name(), func(name string) {
		name = strings.TrimSpace(name)
		if name == "" || strings.ContainsAny(name, `/\`) {
			return
		}
		if d.IsUntitled() {
			d.title = name
			if !d.langSet {
				d.doc.SetLanguage(editor.LanguageForFile(name, nil))
			}
		} else {
			newPath := filepath.Join(filepath.Dir(d.Path()), name)
			if _, err := os.Stat(newPath); err == nil {
				c.showError(errors.New(T().FileExists(newPath)))
				return
			}
			if err := os.Rename(d.Path(), newPath); err != nil {
				c.showError(err)
				return
			}
			d.doc.Path = newPath
			d.statFile()
			if !d.langSet {
				d.doc.SetLanguage(editor.LanguageForFile(newPath, d.doc.Text(0, min(d.doc.Len(), 1024))))
			}
			c.addRecent(newPath)
			c.workspace.refreshSoon()
		}
		for _, p := range c.panes {
			p.refreshTabs()
			if p.current() == d {
				c.applyViewOptions(p)
				p.view.Refresh()
			}
		}
		c.status.refresh()
		c.refreshTitle()
	}))
}

// ---- Files changed by other programs ----

// checkFiles looks for the files changed or deleted by other programs, and
// follows the monitored ones
func (c *MainForm) checkFiles() {
	s := c.settings()
	for _, d := range append([]*Doc(nil), c.docs...) {
		if d.IsUntitled() || d.asking {
			continue
		}
		st, err := os.Stat(d.Path())
		if err != nil {
			if errors.Is(err, os.ErrNotExist) && !d.missing && s.CheckFileChanges {
				d.missing = true
				d.asking = true
				c.showDoc(d)
				ui.ShowQuestionMessageBoxYesNo(c, T().FileDeletedTitle, T().FileDeletedAsk(d.Path()), func() {
					d.asking = false
					d.doc.SetModified()
				}, func() {
					d.asking = false
					for _, p := range c.panes {
						if p.indexOf(d) >= 0 {
							c.doClose(p, d)
						}
					}
				})
			}
			continue
		}
		if d.missing {
			d.missing = false
			d.statFile()
			continue
		}
		if st.ModTime().Equal(d.modTime) && st.Size() == d.fileSize {
			continue
		}
		if d.monitoring {
			c.followFile(d, st.Size())
			continue
		}
		if !s.CheckFileChanges {
			d.statFile()
			continue
		}
		if s.AutoReload && !d.IsModified() {
			c.reload(d, false)
			continue
		}
		d.asking = true
		c.showDoc(d)
		msg := T().FileChangedAsk(d.Path())
		if d.IsModified() {
			msg = T().FileChangedModifiedAsk(d.Path())
		}
		ui.ShowQuestionMessageBoxYesNo(c, T().FileChangedTitle, msg, func() {
			d.asking = false
			c.reload(d, false)
		}, func() {
			d.asking = false
			d.statFile()
		})
	}
}

// followFile adds the end appended to a monitored file (tail -f); a file
// that changed otherwise is read again
func (c *MainForm) followFile(d *Doc, size int64) {
	if size < d.fileSize || d.IsModified() || !d.doc.Encoding.IsUnicode() || d.doc.Encoding.ID != "utf-8" {
		c.reload(d, false)
		return
	}
	f, err := os.Open(d.Path())
	if err != nil {
		return
	}
	defer f.Close()
	if _, err := f.Seek(d.fileSize, io.SeekStart); err != nil {
		return
	}
	data, err := io.ReadAll(io.LimitReader(f, size-d.fileSize))
	if err != nil {
		return
	}
	lf, err := editor.LoadBytes(data, editor.EncodingUTF8)
	if err != nil {
		return
	}
	d.doc.AppendNoUndo(lf.Buffer.Bytes(0, lf.Buffer.Len()))
	d.statFile()
	for _, p := range c.panes {
		if p.current() == d {
			p.view.GotoPos(d.doc.Len())
		}
	}
}
