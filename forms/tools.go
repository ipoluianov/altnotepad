package forms

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// ---- Hashes ----

var hashNames = []string{"MD5", "SHA-1", "SHA-256", "SHA-512"}

func newHash(name string) hash.Hash {
	switch name {
	case "SHA-1":
		return sha1.New()
	case "SHA-256":
		return sha256.New()
	case "SHA-512":
		return sha512.New()
	}
	return md5.New()
}

func hashOf(name string, data []byte) string {
	h := newHash(name)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// hashDialog computes the hash of a text typed in, a line at a time or of the whole
func (c *MainForm) hashDialog(name string) {
	var d ui.DialogContent
	d.InitWidget()
	in := ui.NewTextBox()
	in.SetMultiline(true)
	in.SetYExpandable(true)
	in.SetText(c.view().SelectedText())
	out := ui.NewTextBox()
	out.SetMultiline(true)
	out.SetYExpandable(true)
	out.SetReadOnly(true)
	perLine := ui.NewCheckbox(T().HashEachLine)
	update := func() {
		text := in.Text()
		if perLine.Checked() {
			var b strings.Builder
			for i, l := range strings.Split(text, "\n") {
				if i > 0 {
					b.WriteString("\n")
				}
				b.WriteString(hashOf(name, []byte(strings.TrimSuffix(l, "\r"))))
			}
			out.SetText(b.String())
			return
		}
		out.SetText(hashOf(name, []byte(text)))
	}
	in.SetOnTextChanged(update)
	perLine.SetOnStateChanged(update)
	update()
	d.AddWidget(0, 0, ui.NewLabel(T().HashInput))
	d.AddWidget(1, 0, in)
	d.AddWidget(2, 0, perLine)
	d.AddWidget(3, 0, ui.NewLabel(name+":"))
	d.AddWidget(4, 0, out)
	buttons := ui.NewPanel()
	d.AddWidget(5, 0, buttons)
	copyBtn := ui.NewButton(T().CopyToClipboard)
	copyBtn.SetOnClick(func() { ui.ClipboardSetText(out.Text()) })
	closeBtn := ui.NewButton(T().Close)
	closeBtn.SetOnClick(func() { d.Form().Close() })
	dialogButtons(buttons, copyBtn, closeBtn)
	d.OnDialogShow = func() {
		d.Form().SetTitle(T().MenuHashGenerate(name))
		d.Form().SetSize(620, 460)
		d.Form().MoveToCenterOfParent()
		d.Form().SetCancelButton(closeBtn)
		in.Focus()
	}
	c.ShowDialog(&d)
}

// hashFiles asks for files and shows their hashes in a new document
func (c *MainForm) hashFiles(name string) {
	c.Form().ShowOpenFileDialog(ui.OpenFileDialogOptions{Title: T().MenuHashFiles(name), AllowMultiple: true, DefaultDirectory: c.startDir()},
		func(paths []string, err error) {
			if err != nil || len(paths) == 0 {
				return
			}
			var b strings.Builder
			for _, p := range paths {
				f, err := os.Open(p)
				if err != nil {
					fmt.Fprintf(&b, "%s  %s\n", err, p)
					continue
				}
				h := newHash(name)
				_, err = io.Copy(h, f)
				f.Close()
				if err != nil {
					fmt.Fprintf(&b, "%s  %s\n", err, p)
					continue
				}
				fmt.Fprintf(&b, "%s  %s\n", hex.EncodeToString(h.Sum(nil)), p)
			}
			d := c.newDocumentIn(c.activePane())
			d.doc.Insert(0, []byte(b.String()))
			c.focusEditor()
		})
}

// hashSelection copies the hash of the selected text to the clipboard
func (c *MainForm) hashSelection(name string) {
	text := c.view().SelectedText()
	if text == "" {
		return
	}
	sum := hashOf(name, []byte(text))
	ui.ClipboardSetText(sum)
	c.toast(T().CopiedToClipboard(sum))
}

// ---- Conversions of the selection ----

// transformSelection replaces each selection (or the whole text without
// one) with what f makes of it
func (c *MainForm) transformSelection(f func(s string) (string, error)) {
	v := c.view()
	doc := v.Document()
	sels := v.Selections()
	var edits []editor.Edit
	whole := !v.HasSelection()
	if whole {
		sels = []editor.Sel{{Anchor: 0, Caret: doc.Len()}}
	}
	for _, s := range sels {
		if s.Start() == s.End() {
			continue
		}
		out, err := f(string(doc.Text(s.Start(), s.End())))
		if err != nil {
			c.showError(err)
			return
		}
		edits = append(edits, editor.Edit{Pos: s.Start(), Len: s.End() - s.Start(), Text: []byte(out)})
	}
	doc.BeginAction()
	doc.ApplyEdits(edits)
	doc.EndAction()
	if whole {
		v.SetCaret(0)
	}
	v.Refresh()
}

func base64Encode(s string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(s)), nil
}

func base64Decode(s string) (string, error) {
	s = strings.Join(strings.Fields(s), "")
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		if b2, err2 := base64.RawStdEncoding.DecodeString(s); err2 == nil {
			return string(b2), nil
		}
		if b2, err2 := base64.URLEncoding.DecodeString(s); err2 == nil {
			return string(b2), nil
		}
		return "", err
	}
	return string(b), nil
}

func urlEncode(s string) (string, error) { return url.QueryEscape(s), nil }

func urlDecode(s string) (string, error) { return url.QueryUnescape(s) }

func urlQueryEscape(s string) string { return url.QueryEscape(s) }

func jsonPretty(s string) (string, error) {
	var b bytes.Buffer
	if err := json.Indent(&b, []byte(s), "", "    "); err != nil {
		return "", err
	}
	return b.String(), nil
}

func jsonMinify(s string) (string, error) {
	var b bytes.Buffer
	if err := json.Compact(&b, []byte(s)); err != nil {
		return "", err
	}
	return b.String(), nil
}

// ---- Word completion ----

func (c *MainForm) wordCompletion() {
	c.view().Exec("complete", "")
}

func init() {
	editor.CompletionWords = func(v *editor.View) []string {
		if l := v.Document().Language(); l != nil {
			return l.Words()
		}
		return nil
	}
}

// ---- Styles of the tokens (Search > Style All Occurrences) ----

// styleToken marks the selected text (or the word at the caret) with the
// mark style: all its occurrences or just this one
func (c *MainForm) styleToken(style int, all bool) {
	v := c.view()
	doc := v.Document()
	s := v.MainSelection()
	a, b := s.Start(), s.End()
	if a == b {
		a, b = v.WordAtCaret()
	}
	if a == b {
		return
	}
	if !all {
		doc.Marks[style].Add(editor.Range{Start: a, End: b})
		v.Refresh()
		return
	}
	text := string(doc.Text(a, b))
	m, err := editor.NewMatcher(text, editor.SearchOptions{MatchCase: true, WholeWord: s.Start() == s.End()})
	if err != nil {
		return
	}
	m.FindAll(doc, 0, doc.Len(), func(x, y int) bool {
		doc.Marks[style].Add(editor.Range{Start: x, End: y})
		return true
	})
	v.Refresh()
}

func (c *MainForm) clearStyle(style int) {
	c.view().Document().Marks[style].Clear()
	c.view().Refresh()
}

// jumpStyle goes to the next (or previous) text marked with the style
func (c *MainForm) jumpStyle(style int, dir int) {
	v := c.view()
	doc := v.Document()
	ranges := doc.Marks[style].Ranges()
	if len(ranges) == 0 {
		return
	}
	pos := v.Caret()
	s := v.MainSelection()
	if dir > 0 {
		for _, r := range ranges {
			if r.Start >= s.End() && !(r.Start == s.Start() && r.End == s.End()) {
				v.SelectRange(r.Start, r.End)
				return
			}
		}
		v.SelectRange(ranges[0].Start, ranges[0].End)
		return
	}
	for i := len(ranges) - 1; i >= 0; i-- {
		if ranges[i].End <= min(pos, s.Start()) {
			v.SelectRange(ranges[i].Start, ranges[i].End)
			return
		}
	}
	last := ranges[len(ranges)-1]
	v.SelectRange(last.Start, last.End)
}

// ---- Macros ----

func (c *MainForm) toggleRecording() {
	if c.recording {
		c.stopRecording()
		return
	}
	c.macro = &editor.Macro{}
	c.recording = true
	for _, p := range c.panes {
		p.view.StartRecording(c.macro)
	}
	c.toolbar.refresh()
	c.refreshTitle()
}

func (c *MainForm) stopRecording() {
	if !c.recording {
		return
	}
	c.recording = false
	for _, p := range c.panes {
		p.view.StopRecording()
	}
	c.toolbar.refresh()
	c.refreshTitle()
}

// recordStep adds a step of the application (e.g. a search) to the macro being recorded
func (c *MainForm) recordStep(cmd, arg string) {
	if c.recording {
		c.view().RecordStep(cmd, arg)
	}
}

func (c *MainForm) playMacro(m *editor.Macro, times int) {
	if m == nil || len(m.Steps) == 0 || c.recording {
		return
	}
	c.view().Play(m, times)
	c.status.refresh()
}

func (c *MainForm) saveMacro() {
	if c.macro == nil || len(c.macro.Steps) == 0 {
		return
	}
	m := c.macro
	c.ShowDialog(NewInputDialog(T().MenuSaveMacro, T().MacroName, "", func(name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			return
		}
		steps := make([]config.MacroStep, len(m.Steps))
		for i, s := range m.Steps {
			steps[i] = config.MacroStep{Cmd: s.Cmd, Arg: s.Arg}
		}
		config.UpdateHistory(func(h *config.History) {
			if h.Macros == nil {
				h.Macros = map[string][]config.MacroStep{}
			}
			h.Macros[name] = steps
		})
	}))
}

func (c *MainForm) playSavedMacro(name string) {
	steps := config.GetHistory().Macros[name]
	m := &editor.Macro{}
	for _, s := range steps {
		m.Steps = append(m.Steps, editor.MacroStep{Cmd: s.Cmd, Arg: s.Arg})
	}
	c.playMacro(m, 1)
}

// runMacroMultiple runs the macro a number of times or up to the end of the document
func (c *MainForm) runMacroMultiple() {
	if c.macro == nil || len(c.macro.Steps) == 0 {
		return
	}
	m := c.macro
	c.ShowDialog(NewInputDialog(T().MenuRunMacroMulti, T().MacroTimes, "1", func(s string) {
		s = strings.TrimSpace(s)
		if s == "*" || strings.EqualFold(s, "end") {
			c.playMacro(m, -1)
			return
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			return
		}
		c.playMacro(m, n)
	}))
}

// ---- Sessions ----

func (c *MainForm) saveSessionDialog() {
	c.Form().ShowSaveFileDialog(ui.SaveFileDialogOptions{Title: T().MenuSaveSession, DefaultDirectory: c.startDir(), DefaultFileName: "session.json"},
		func(path string, err error) {
			if err != nil || path == "" {
				return
			}
			s := c.sessionOf(false)
			bs, err := json.MarshalIndent(s, "", "  ")
			if err == nil {
				err = os.WriteFile(path, bs, 0644)
			}
			if err != nil {
				c.showError(err)
			}
		})
}

func (c *MainForm) loadSessionDialog() {
	c.Form().ShowOpenFileDialog(ui.OpenFileDialogOptions{Title: T().MenuLoadSession, DefaultDirectory: c.startDir()},
		func(paths []string, err error) {
			if err != nil || len(paths) == 0 {
				return
			}
			bs, err := os.ReadFile(paths[0])
			if err != nil {
				c.showError(err)
				return
			}
			var s config.Session
			if err := json.Unmarshal(bs, &s); err != nil {
				c.showError(errors.New(T().NotASession(filepath.Base(paths[0]))))
				return
			}
			// The files of the session are opened in addition to the open ones
			for _, f := range s.Files {
				if f.Path != "" {
					c.openFile(f.Path, 0)
				}
			}
		})
}
