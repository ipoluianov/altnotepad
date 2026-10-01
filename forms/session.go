package forms

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
)

// ---- Session ----

// sessionOf describes the open documents; with backups the unsaved texts are
// written to the backup directory
func (c *MainForm) sessionOf(withBackups bool) config.Session {
	var s config.Session
	index := map[*Doc]int{}
	for _, d := range c.docs {
		if d.IsUntitled() && (!withBackups || (d.doc.Len() == 0 && !d.IsModified())) {
			continue
		}
		if !d.IsUntitled() && d.IsModified() && !withBackups {
			// Saved as it is on the disk
		}
		sf := config.SessionFile{Path: d.Path(), Title: d.title, EOL: int(d.doc.EOL),
			Encoding: d.doc.Encoding.ID, BOM: d.doc.Encoding.BOM, ReadOnly: d.doc.ReadOnly,
			LanguageSet: d.langSet, Bookmarks: d.doc.Bookmarks.Lines()}
		if l := d.doc.Language(); l != nil {
			sf.Language = l.ID
		}
		if !d.modTime.IsZero() {
			sf.ModTime = d.modTime.UnixNano()
		}
		for v := 0; v < 2; v++ {
			st := d.states[v]
			if st == nil {
				continue
			}
			sf.Views |= 1 << v
			if v == 0 || sf.Views == 2 {
				sf.Caret, sf.Anchor, sf.TopLine = st.Caret(), st.Anchor(), st.TopLine()
				sf.Folds = st.FoldedLines()
			} else {
				sf.Caret2, sf.Anchor2, sf.TopLine2 = st.Caret(), st.Anchor(), st.TopLine()
			}
		}
		if withBackups && (d.IsModified() || d.IsUntitled()) {
			if name, err := c.writeBackupNow(d); err == nil {
				sf.Backup = name
			}
		}
		index[d] = len(s.Files)
		s.Files = append(s.Files, sf)
	}
	for v, p := range c.panes {
		for i, d := range p.docs {
			if k, ok := index[d]; ok {
				s.Tabs[v] = append(s.Tabs[v], k)
				if i == p.cur {
					s.Active[v] = len(s.Tabs[v]) - 1
				}
			}
		}
	}
	s.ActiveView = c.active
	return s
}

// SaveSession remembers the open documents for the next start
func (c *MainForm) SaveSession() {
	st := c.settings()
	if !st.RememberSession {
		config.SaveSession(config.Session{})
		c.cleanBackups(nil)
		return
	}
	s := c.sessionOf(st.Backup)
	config.SaveSession(s)
	used := map[string]bool{}
	for _, f := range s.Files {
		if f.Backup != "" {
			used[f.Backup] = true
		}
	}
	c.cleanBackups(used)
}

// RestoreSession opens the documents of the last session; returns false if there were none
func (c *MainForm) RestoreSession() bool {
	if !c.settings().RememberSession {
		return false
	}
	s, ok := config.LoadSession("")
	if !ok {
		return false
	}
	return c.applySession(s)
}

// applySession opens the documents of a session
func (c *MainForm) applySession(s config.Session) bool {
	docs := make([]*Doc, len(s.Files))
	for i, f := range s.Files {
		docs[i] = c.restoreFile(f)
	}
	opened := false
	for v := 1; v >= 0; v-- {
		p := c.panes[v]
		for k, idx := range s.Tabs[v] {
			if idx < 0 || idx >= len(docs) || docs[idx] == nil {
				continue
			}
			d := docs[idx]
			f := s.Files[idx]
			if v == 1 {
				c.showSecondView(true)
			}
			st := d.state(v)
			if v == 0 || f.Views == 2 {
				st.SetPosition(f.Caret, f.Anchor, f.TopLine)
				st.SetFoldedLines(f.Folds)
			} else {
				st.SetPosition(f.Caret2, f.Anchor2, f.TopLine2)
			}
			p.docs = append(p.docs, d)
			c.applyDocOptions(d, v)
			if k == s.Active[v] {
				p.cur = len(p.docs) - 1
			}
			opened = true
		}
		if len(p.docs) > 0 {
			if p.cur < 0 {
				p.cur = 0
			}
			p.activate(p.cur)
		}
	}
	for _, d := range docs {
		if d != nil && !d.inView(0) && !d.inView(1) {
			c.forgetDoc(d)
		}
	}
	if len(c.panes[1].docs) == 0 {
		c.showSecondView(false)
	}
	if s.ActiveView == 1 && len(c.panes[1].docs) > 0 {
		c.setActivePane(1)
	}
	return opened
}

// restoreFile opens a document of the session: its backup, or its file
func (c *MainForm) restoreFile(f config.SessionFile) *Doc {
	d := newDoc()
	d.langSet = f.LanguageSet
	enc := editor.Encoding{ID: f.Encoding, BOM: f.BOM}
	if f.Backup != "" {
		lf, err := editor.LoadFile(filepath.Join(config.BackupDirectory(), f.Backup), editor.EncodingUTF8, nil)
		if err == nil {
			d.doc.SetLarge(int64(lf.Buffer.Len()) > int64(c.settings().LargeFileMB)<<20)
			d.doc.SetBuffer(lf.Buffer)
			d.doc.Path = f.Path
			d.title = f.Title
			if f.Path == "" && d.title == "" {
				d.title = c.nextUntitledTitle()
			}
			if f.Encoding != "" {
				d.doc.Encoding = enc
			}
			d.doc.EOL = editor.EOL(f.EOL)
			d.doc.SetModified()
			d.backupName = f.Backup
			if f.ModTime != 0 {
				d.modTime = time.Unix(0, f.ModTime)
				if st, err := os.Stat(f.Path); err == nil {
					d.fileSize = st.Size()
					if !st.ModTime().Equal(d.modTime) {
						d.fileSize = -1 // changed since: asked about
					}
				}
			}
		} else if f.Path == "" {
			return nil
		}
	}
	if d.backupName == "" {
		if f.Path == "" {
			return nil
		}
		if _, err := os.Stat(f.Path); err != nil {
			return nil
		}
		if f.Encoding == "" {
			enc = editor.Encoding{}
		}
		if err := loadInto(d, f.Path, enc); err != nil {
			return nil
		}
	}
	lang := editor.LanguageByID(f.Language)
	if f.Language == "" {
		lang = editor.LanguageForFile(d.FullName(), d.doc.Text(0, min(d.doc.Len(), 1024)))
	}
	d.doc.SetLanguage(lang)
	d.doc.ReadOnly = f.ReadOnly
	d.doc.Bookmarks.Set(f.Bookmarks)
	c.docs = append(c.docs, d)
	return d
}

// ---- Backups ----

func newBackupName() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b) + ".txt"
}

// writeBackupNow writes the text of the document to its backup file
func (c *MainForm) writeBackupNow(d *Doc) (string, error) {
	if d.backupName == "" {
		d.backupName = newBackupName()
	}
	if d.backupVersion == d.doc.Version() {
		if _, err := os.Stat(filepath.Join(config.BackupDirectory(), d.backupName)); err == nil {
			return d.backupName, nil
		}
	}
	if err := os.MkdirAll(config.BackupDirectory(), 0700); err != nil {
		return "", err
	}
	path := filepath.Join(config.BackupDirectory(), d.backupName)
	if err := editor.SaveFile(path, d.doc.Buffer(), editor.EncodingUTF8, editor.EOLLF, true); err != nil {
		return "", err
	}
	d.backupVersion = d.doc.Version()
	return d.backupName, nil
}

var lastBackup time.Time

// backupTimer backs up the unsaved changes every few seconds, in the
// background, so they survive a crash
func (c *MainForm) backupTimer() {
	s := c.settings()
	if !s.Backup || !s.RememberSession || time.Since(lastBackup) < time.Duration(s.BackupSeconds)*time.Second {
		return
	}
	lastBackup = time.Now()
	changed := false
	for _, d := range c.docs {
		if !(d.IsModified() || (d.IsUntitled() && d.doc.Len() > 0)) || d.backupVersion == d.doc.Version() {
			continue
		}
		if d.backupName == "" {
			d.backupName = newBackupName()
		}
		// A snapshot of the text is written on another goroutine
		buf := d.doc.Buffer().Clone()
		name := d.backupName
		d.backupVersion = d.doc.Version()
		changed = true
		go func() {
			os.MkdirAll(config.BackupDirectory(), 0700)
			editor.SaveFile(filepath.Join(config.BackupDirectory(), name), buf, editor.EncodingUTF8, editor.EOLLF, true)
		}()
	}
	if changed {
		// The session refers to the backups
		config.SaveSession(c.sessionOfNoWrite())
	}
}

// sessionOfNoWrite describes the session with the backups already written
func (c *MainForm) sessionOfNoWrite() config.Session {
	s := c.sessionOf(false)
	byPath := map[string]*Doc{}
	var untitled []*Doc
	for _, d := range c.docs {
		if d.IsUntitled() {
			untitled = append(untitled, d)
		} else {
			byPath[d.Path()] = d
		}
	}
	// Add the new documents with backups, and the backups of the modified files
	index := map[*Doc]int{}
	for i, f := range s.Files {
		if d := byPath[f.Path]; d != nil {
			index[d] = i
			if d.IsModified() && d.backupName != "" {
				s.Files[i].Backup = d.backupName
			}
		}
	}
	for _, d := range untitled {
		if d.backupName == "" {
			continue
		}
		sf := config.SessionFile{Title: d.title, Backup: d.backupName, EOL: int(d.doc.EOL), Encoding: d.doc.Encoding.ID, BOM: d.doc.Encoding.BOM}
		if l := d.doc.Language(); l != nil {
			sf.Language = l.ID
		}
		for v := 0; v < 2; v++ {
			if st := d.states[v]; st != nil {
				sf.Views |= 1 << v
				sf.Caret, sf.Anchor, sf.TopLine = st.Caret(), st.Anchor(), st.TopLine()
			}
		}
		index[d] = len(s.Files)
		s.Files = append(s.Files, sf)
	}
	s.Tabs = [2][]int{}
	for v, p := range c.panes {
		for i, d := range p.docs {
			if k, ok := index[d]; ok {
				s.Tabs[v] = append(s.Tabs[v], k)
				if i == p.cur {
					s.Active[v] = len(s.Tabs[v]) - 1
				}
			}
		}
	}
	return s
}

// removeBackup deletes the backup of a document that was saved or closed
func (c *MainForm) removeBackup(d *Doc) {
	if d.backupName == "" {
		return
	}
	os.Remove(filepath.Join(config.BackupDirectory(), d.backupName))
	d.backupName = ""
	d.backupVersion = 0
}

// cleanBackups deletes the backup files not in use
func (c *MainForm) cleanBackups(used map[string]bool) {
	entries, err := os.ReadDir(config.BackupDirectory())
	if err != nil {
		return
	}
	for _, e := range entries {
		if !used[e.Name()] {
			os.Remove(filepath.Join(config.BackupDirectory(), e.Name()))
		}
	}
}

// ---- Exit ----

// RequestExit is called when the window is closed: with the backups on the
// unsaved documents are kept for the next start, otherwise it asks to save
// them. Returns whether the window can close now.
func (c *MainForm) RequestExit() bool {
	s := c.settings()
	c.SaveWindowState()
	if s.Backup && s.RememberSession {
		c.quitting = true
		c.SaveSession()
		return true
	}
	var modified []*Doc
	for _, d := range c.docs {
		if d.IsModified() {
			modified = append(modified, d)
		}
	}
	if len(modified) == 0 {
		c.quitting = true
		c.SaveSession()
		return true
	}
	var next func(i int)
	next = func(i int) {
		if i >= len(modified) {
			c.quitting = true
			c.SaveSession()
			c.Form().Close()
			return
		}
		d := modified[i]
		c.showDoc(d)
		showSaveChanges(c, d.FullName(), func(answer int) {
			switch answer {
			case answerYes:
				c.save(d, func(ok bool) {
					if ok {
						next(i + 1)
					}
				})
			case answerNo:
				// Not saved: the file stays as it is on the disk
				d.doc.SetSavePoint()
				next(i + 1)
			}
		})
	}
	next(0)
	return false
}
