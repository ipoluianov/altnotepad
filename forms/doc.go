package forms

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
)

// Doc is an open document: its text, file and views
type Doc struct {
	doc *editor.Document
	// title of a new document never saved, like "new 1"
	title string
	// states of the document in the two views, nil where it is not open
	states [2]*editor.ViewState
	// langSet: the language was chosen by hand, a rename keeps it
	langSet bool

	// The file as last read or written, to tell when another program changes it
	modTime  time.Time
	fileSize int64
	// missing: the file was deleted or moved by another program
	missing bool
	// asking: a question about the file change is shown
	asking bool

	monitoring bool

	backupName    string
	backupVersion int64
}

func newDoc() *Doc {
	d := &Doc{doc: editor.NewDocument()}
	return d
}

// Path returns the file of the document, "" for a new one
func (d *Doc) Path() string { return d.doc.Path }

// IsUntitled reports whether the document was never saved
func (d *Doc) IsUntitled() bool { return d.doc.Path == "" }

// Name returns the name shown in the tab
func (d *Doc) Name() string {
	if d.doc.Path == "" {
		return d.title
	}
	return filepath.Base(d.doc.Path)
}

// FullName returns the path, or the name of a new document
func (d *Doc) FullName() string {
	if d.doc.Path == "" {
		return d.title
	}
	return d.doc.Path
}

// IsModified reports whether the document has unsaved changes
func (d *Doc) IsModified() bool { return d.doc.IsModified() }

// isEmptyUntitled: a new document nothing was typed in, replaced by the first file opened
func (d *Doc) isEmptyUntitled() bool {
	return d.IsUntitled() && d.doc.Len() == 0 && !d.doc.IsModified()
}

// state returns the state of the document in the view, making it
func (d *Doc) state(view int) *editor.ViewState {
	if d.states[view] == nil {
		d.states[view] = editor.NewViewState(d.doc)
	}
	return d.states[view]
}

func (d *Doc) inView(view int) bool { return d.states[view] != nil }

// release forgets the state of the view
func (d *Doc) release(view int) {
	if d.states[view] != nil {
		d.states[view].Close()
		d.states[view] = nil
	}
}

// statFile remembers the time and size of the file on the disk
func (d *Doc) statFile() {
	if d.doc.Path == "" {
		return
	}
	if st, err := os.Stat(d.doc.Path); err == nil {
		d.modTime = st.ModTime()
		d.fileSize = st.Size()
		d.missing = false
		d.doc.FileReadOnly = st.Mode().Perm()&0200 == 0
	}
}

// loadInto reads the file into the document, replacing its text
func loadInto(d *Doc, path string, enc editor.Encoding) error {
	lf, err := editor.LoadFile(path, enc, nil)
	if err != nil {
		return err
	}
	applyLoaded(d, path, lf)
	return nil
}

// applyLoaded puts a loaded file into the document
func applyLoaded(d *Doc, path string, lf *editor.LoadedFile) {
	s := config.GetSettings()
	large := int64(lf.Buffer.Len()) > int64(s.LargeFileMB)<<20
	d.doc.Path = path
	d.doc.SetLarge(large)
	d.doc.SetBuffer(lf.Buffer)
	d.doc.Encoding = lf.Encoding
	d.doc.EOL = lf.EOL
	if !d.langSet {
		first := d.doc.Text(0, min(d.doc.Len(), d.doc.LineEnd(0), 1024))
		d.doc.SetLanguage(editor.LanguageForFile(path, first))
	}
	d.statFile()
}

// defaultEncoding returns the encoding of new documents from the settings
func defaultEncoding(s config.Settings) editor.Encoding {
	switch s.NewDocEncoding {
	case "", "utf-8":
		return editor.EncodingUTF8
	case "utf-8-bom":
		return editor.EncodingUTF8BOM
	case "ansi":
		return editor.Encoding{ID: editor.DefaultANSI}
	}
	if _, ok := editor.CharsetByID(s.NewDocEncoding); ok {
		return editor.Encoding{ID: s.NewDocEncoding}
	}
	return editor.EncodingUTF8
}

// defaultEOL returns the line break of new documents from the settings
func defaultEOL(s config.Settings) editor.EOL {
	switch s.NewDocEOL {
	case "crlf":
		return editor.EOLCRLF
	case "lf":
		return editor.EOLLF
	case "cr":
		return editor.EOLCR
	}
	return editor.DefaultEOL()
}
