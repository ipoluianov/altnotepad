package forms

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// hit is a line with matches
type hit struct {
	line   int    // 0-based
	text   string // the line, shortened when long
	ranges [][2]int
	pos    int // the document position of the first match
	end    int
}

// fileHits are the matches in a file
type fileHits struct {
	path string
	doc  *Doc // the open document, nil for a file on the disk
	hits []hit
}

// searchResult is the result of a Find All
type searchResult struct {
	text     string
	opts     editor.SearchOptions
	files    []fileHits
	searched int
	stopped  bool
}

func (r *searchResult) total() int {
	n := 0
	for _, f := range r.files {
		for _, h := range f.hits {
			n += len(h.ranges)
		}
	}
	return n
}

const (
	maxHitLine    = 400     // longer lines are shortened in the results
	maxHitsPerDoc = 200_000 // lines with matches listed per document
)

// collectMatches finds the matches of the document between from and to, by lines
func collectMatches(doc *editor.Document, m *editor.Matcher, from, to int) fileHits {
	var fh fileHits
	last := -1
	m.FindAll(doc, from, to, func(a, b int) bool {
		line := doc.LineOfOffset(a)
		ls := doc.LineStart(line)
		if line != last {
			if len(fh.hits) >= maxHitsPerDoc {
				return false
			}
			le := doc.LineEnd(line)
			text := doc.Text(ls, min(le, ls+maxHitLine*4))
			s := string(bytes.ToValidUTF8(text, []byte("?")))
			if r := []rune(s); len(r) > maxHitLine {
				s = string(r[:maxHitLine]) + "…"
			}
			fh.hits = append(fh.hits, hit{line: line, text: s, pos: a, end: b})
			last = line
		}
		h := &fh.hits[len(fh.hits)-1]
		h.ranges = append(h.ranges, [2]int{a - ls, b - ls})
		return true
	})
	return fh
}

// resultEntry is what a line of the results refers to
type resultEntry struct {
	path string
	doc  *Doc
	hit  *hit // nil for a header line
}

// ResultsPanel shows the results of Find All, grouped by search and by file;
// a double click (or F4) opens the match
type ResultsPanel struct {
	ui.Widget

	main    *MainForm
	title   *ui.Label
	stopBtn *ui.Button
	view    *editor.View
	doc     *editor.Document
	results []*searchResult
	entries []resultEntry
	current int
	cancel  *atomic.Bool
}

func NewResultsPanel(main *MainForm) *ResultsPanel {
	var c ResultsPanel
	c.InitWidget()
	c.main = main
	c.current = -1
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	head := ui.NewPanel()
	head.SetPanelPadding(2)
	c.title = ui.NewLabel(T().SearchResults)
	c.title.SetTextFunc(func() string { return T().SearchResults })
	head.AddWidget(0, 0, c.title)
	head.AddWidget(0, 1, ui.NewHSpacer())
	c.stopBtn = ui.NewButton(T().Stop)
	c.stopBtn.SetTextFunc(func() string { return T().Stop })
	c.stopBtn.SetOnClick(func() {
		if c.cancel != nil {
			c.cancel.Store(true)
		}
	})
	c.stopBtn.SetVisible(false)
	head.AddWidget(0, 2, c.stopBtn)
	clear := ui.NewButton(T().Clear)
	clear.SetTextFunc(func() string { return T().Clear })
	clear.SetOnClick(c.clear)
	head.AddWidget(0, 3, clear)
	head.AddWidget(0, 4, smallButton("x", func() string { return T().Close }, func() {
		c.SetVisible(false)
		c.main.updateLayout()
		c.main.focusEditor()
	}))
	c.AddWidget(0, 0, head)

	c.view = editor.NewView()
	c.doc = editor.NewDocument()
	c.doc.SetLanguage(editor.LanguageByID("searchresults"))
	c.doc.ReadOnly = true
	c.view.SetState(editor.NewViewState(c.doc))
	c.view.SetReadOnlyView(true)
	c.view.OnDoubleClickLine = func(line int) bool {
		if line >= 0 && line < len(c.entries) && c.entries[line].hit != nil {
			c.current = line
			c.open(line)
			return true
		}
		return false
	}
	c.AddWidget(1, 0, c.view)
	c.applySettings()
	return &c
}

func (c *ResultsPanel) applyLanguage() {}

// applySettings shows the results with the colors and the font of the editors
func (c *ResultsPanel) applySettings() {
	o := c.main.viewOptions(nil)
	o.LineNumbers = false
	o.Bookmarks = false
	o.WordWrap = false
	o.CurrentLine = true
	o.SmartHighlight = false
	o.ShowSpaces, o.ShowEOL, o.ShowIndentGuides = false, false, false
	o.TabSize = 4
	c.view.SetOptions(o)
	c.view.SetScheme(editorScheme(c.main.settings()))
}

func (c *ResultsPanel) clear() {
	c.results = nil
	c.rebuild()
}

// add shows a new result above the old ones and opens the panel
func (c *ResultsPanel) add(r *searchResult) {
	c.results = append([]*searchResult{r}, c.results...)
	if len(c.results) > 20 {
		c.results = c.results[:20]
	}
	c.current = -1
	c.rebuild()
	c.view.SetCaret(0)
	c.view.ScrollToLine(0)
	if !c.IsVisible() {
		c.SetVisible(true)
		c.main.vertSplit.SetSecondSize(max(120, c.main.resultsH))
		c.main.updateLayout()
	}
}

// rebuild makes the text of the results
func (c *ResultsPanel) rebuild() {
	var b strings.Builder
	c.entries = c.entries[:0]
	type mark struct{ a, b int }
	var marks []mark
	for ri, r := range c.results {
		files := len(r.files)
		header := T().ResultsHeader(r.text, r.total(), files, r.searched)
		if r.stopped {
			header += " " + T().Stopped
		}
		b.WriteString(header)
		b.WriteString("\n")
		c.entries = append(c.entries, resultEntry{})
		for fi := range r.files {
			f := &r.files[fi]
			n := 0
			for _, h := range f.hits {
				n += len(h.ranges)
			}
			fmt.Fprintf(&b, "  %s %s\n", f.path, T().HitsCount(n))
			c.entries = append(c.entries, resultEntry{path: f.path, doc: f.doc})
			for hi := range f.hits {
				h := &f.hits[hi]
				prefix := fmt.Sprintf("\t%s %d: ", T().Line, h.line+1)
				start := b.Len() + len(prefix)
				b.WriteString(prefix)
				b.WriteString(h.text)
				b.WriteString("\n")
				for _, rg := range h.ranges {
					if rg[1] <= len(h.text) {
						marks = append(marks, mark{start + rg[0], start + rg[1]})
					}
				}
				c.entries = append(c.entries, resultEntry{path: f.path, doc: f.doc, hit: h})
			}
		}
		_ = ri
	}
	c.doc.ReadOnly = false
	c.doc.SetText([]byte(b.String()))
	c.doc.ReadOnly = true
	for _, m := range marks {
		c.doc.Marks[editor.MarkStyleFind].Add(editor.Range{Start: m.a, End: m.b})
	}
	// The older searches are folded
	st := c.view.State()
	var folds []int
	for i, e := range c.entries {
		if i > 0 && e.path == "" && e.hit == nil {
			folds = append(folds, i)
		}
	}
	st.SetFoldedLines(folds)
	c.view.Refresh()
}

// open opens the match of the line of the results
func (c *ResultsPanel) open(line int) {
	e := c.entries[line]
	c.view.SetCaret(c.doc.LineStart(line))
	c.view.EnsureCaretVisible(true)
	var d *Doc
	if e.doc != nil && c.main.isOpen(e.doc) {
		d = e.doc
		c.main.showDoc(d)
	} else if e.path != "" {
		d = c.main.openFile(e.path, 0)
	}
	if d == nil {
		return
	}
	v := c.main.view()
	doc := d.doc
	if e.hit.line >= doc.LineCount() {
		return
	}
	ls := doc.LineStart(e.hit.line)
	a := ls + e.hit.ranges[0][0]
	b := ls + e.hit.ranges[0][1]
	if b <= doc.LineEnd(e.hit.line)+1 {
		v.SelectRange(a, b)
	} else {
		v.GotoLine(e.hit.line)
	}
	v.Focus()
}

// jump opens the next (dir 1) or previous match
func (c *ResultsPanel) jump(dir int) {
	if len(c.entries) == 0 {
		return
	}
	i := c.current
	for k := 0; k < len(c.entries); k++ {
		i += dir
		if i < 0 {
			i = len(c.entries) - 1
		}
		if i >= len(c.entries) {
			i = 0
		}
		if c.entries[i].hit != nil {
			c.current = i
			if !c.IsVisible() {
				c.SetVisible(true)
				c.main.updateLayout()
			}
			c.open(i)
			return
		}
	}
}

// isOpen reports whether the document is still open
func (c *MainForm) isOpen(d *Doc) bool {
	for _, x := range c.docs {
		if x == d {
			return true
		}
	}
	return false
}

// ---- Find in files ----

// fileFilter tells which files to search: "*.go *.txt;!*_test.go"
type fileFilter struct {
	include []string
	exclude []string
}

func parseFilter(s string) fileFilter {
	var f fileFilter
	for _, p := range strings.FieldsFunc(s, func(r rune) bool { return r == ';' || r == ' ' || r == ',' }) {
		if strings.HasPrefix(p, "!") {
			f.exclude = append(f.exclude, p[1:])
		} else {
			f.include = append(f.include, p)
		}
	}
	if len(f.include) == 0 {
		f.include = []string{"*"}
	}
	return f
}

func (f fileFilter) match(name string) bool {
	lname := strings.ToLower(name)
	for _, p := range f.exclude {
		if ok, _ := filepath.Match(strings.ToLower(p), lname); ok {
			return false
		}
	}
	for _, p := range f.include {
		if p == "*" || p == "*.*" {
			return true
		}
		if ok, _ := filepath.Match(strings.ToLower(p), lname); ok {
			return true
		}
	}
	return false
}

// isBinary reports whether the file looks binary: zero bytes in its start
func isBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()
	buf := make([]byte, 8000)
	n, _ := io.ReadFull(f, buf)
	buf = buf[:n]
	if bytes.HasPrefix(buf, []byte{0xFF, 0xFE}) || bytes.HasPrefix(buf, []byte{0xFE, 0xFF}) {
		return false // UTF-16
	}
	return bytes.IndexByte(buf, 0) >= 0
}

// findInFiles searches (and replaces) in the files of the folder, in the background
func (c *FindPanel) findInFiles(replace bool) {
	s := c.current()
	m, err := s.matcher()
	if err != nil {
		c.setStatus(err.Error(), true)
		return
	}
	dir := strings.TrimSpace(c.cmbDir.Text())
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		c.setStatus(T().NoSuchFolder(dir), true)
		return
	}
	filterText := strings.TrimSpace(c.cmbFilters.Text())
	c.remember(replace)
	config.UpdateHistory(func(h *config.History) {
		h.FindDirs = config.PushFront(h.FindDirs, dir, 0)
		if filterText != "" {
			h.FindFilters = config.PushFront(h.FindFilters, filterText, 0)
		}
	})
	start := func() {
		c.runFindInFiles(s, m, dir, parseFilter(filterText), c.chkSubdirs.Checked(), c.chkHidden.Checked(), replace)
	}
	if replace {
		ui.ShowQuestionMessageBoxYesNo(c.main, T().ReplaceInFiles, T().ReplaceInFilesAsk(s.Text, s.Replace, dir), start, nil)
		return
	}
	start()
}

func (c *FindPanel) runFindInFiles(s search, m *editor.Matcher, dir string, filter fileFilter, subdirs, hidden, replace bool) {
	main := c.main
	res := &searchResult{text: s.Text, opts: s.Opts}
	// The open files are searched in the editor, on this goroutine
	open := map[string]*Doc{}
	for _, d := range main.docs {
		if d.Path() != "" {
			open[filepath.Clean(d.Path())] = d
		}
	}
	cancel := &atomic.Bool{}
	main.results.cancel = cancel
	main.results.stopBtn.SetVisible(true)
	main.updateLayout()
	form := main.Form()
	var openToSearch []*Doc
	var replaced atomic.Int64
	largeLimit := int64(main.settings().LargeFileMB) << 20
	go func() {
		lastProgress := time.Now()
		var found []fileHits
		filepath.WalkDir(dir, func(path string, de fs.DirEntry, err error) error {
			if cancel.Load() {
				return errors.New("stopped")
			}
			if err != nil {
				return nil
			}
			name := de.Name()
			if de.IsDir() {
				if path != dir && (!subdirs || (!hidden && strings.HasPrefix(name, "."))) {
					return filepath.SkipDir
				}
				return nil
			}
			if !de.Type().IsRegular() || !filter.match(name) {
				return nil
			}
			if !hidden && strings.HasPrefix(name, ".") {
				return nil
			}
			if d, ok := open[filepath.Clean(path)]; ok {
				openToSearch = append(openToSearch, d)
				return nil
			}
			if time.Since(lastProgress) > 150*time.Millisecond {
				lastProgress = time.Now()
				p := path
				form.Invoke(func() { main.status.setMessage(T().Searching(p)) })
			}
			if info, err := de.Info(); err != nil || info.Size() > largeLimit*4 || isBinary(path) {
				return nil
			}
			res.searched++
			lf, err := editor.LoadFile(path, editor.Encoding{}, nil)
			if err != nil {
				return nil
			}
			doc := editor.NewDocument()
			doc.SetBuffer(lf.Buffer)
			fh := collectMatches(doc, m, 0, doc.Len())
			if len(fh.hits) == 0 {
				return nil
			}
			fh.path = path
			found = append(found, fh)
			if replace {
				n := replaceInDoc(doc, s, 0, doc.Len())
				if n > 0 && editor.SaveFile(path, doc.Buffer(), lf.Encoding, lf.EOL, true) == nil {
					replaced.Add(int64(n))
				}
			}
			return nil
		})
		form.Invoke(func() {
			for _, d := range openToSearch {
				res.searched++
				fh := collectMatches(d.doc, m, 0, d.doc.Len())
				if len(fh.hits) == 0 {
					continue
				}
				fh.path = d.Path()
				fh.doc = d
				found = append(found, fh)
				if replace && !d.doc.ReadOnly {
					replaced.Add(int64(replaceInDoc(d.doc, s, 0, d.doc.Len())))
				}
			}
			res.files = found
			res.stopped = cancel.Load()
			main.results.stopBtn.SetVisible(false)
			main.status.setMessage("")
			main.results.add(res)
			if replace {
				c.setStatus(T().ReplacedInFiles(int(replaced.Load()), len(found)), replaced.Load() == 0)
				main.refreshAllViews()
			} else {
				c.setStatus(T().FoundInFiles(res.total(), len(found)), res.total() == 0)
			}
		})
	}()
}
