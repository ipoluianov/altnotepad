package forms

import (
	"bufio"
	"fmt"
	"html"
	"image/color"
	"io"
	"os"
	"path/filepath"

	"github.com/ipoluianov/altnotepad/app"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

func cssColor(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// writeHTML writes the document as HTML with the colors of the syntax
// highlighting; print adds a script that opens the print dialog of the browser
func writeHTML(w io.Writer, d *Doc, scheme *editor.Scheme, tabSize int, print bool) error {
	bw := bufio.NewWriterSize(w, 1<<20)
	doc := d.doc
	fmt.Fprintf(bw, "<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>%s</title>\n<style>\n", html.EscapeString(d.Name()))
	fmt.Fprintf(bw, "body { background: %s; color: %s; margin: 0; }\n", cssColor(scheme.Background), cssColor(scheme.Foreground))
	fmt.Fprintf(bw, "pre { font-family: 'JetBrains Mono', Consolas, 'DejaVu Sans Mono', monospace; font-size: 10pt; margin: 12px; tab-size: %d; white-space: pre-wrap; }\n", tabSize)
	for i, st := range scheme.Styles {
		fmt.Fprintf(bw, ".s%d { color: %s;", i, cssColor(st.Fore))
		if st.Back.A > 0 {
			fmt.Fprintf(bw, " background: %s;", cssColor(st.Back))
		}
		if st.Bold {
			bw.WriteString(" font-weight: bold;")
		}
		if st.Italic {
			bw.WriteString(" font-style: italic;")
		}
		bw.WriteString(" }\n")
	}
	bw.WriteString("@media print { body { background: #fff; } }\n</style>\n")
	if print {
		bw.WriteString("<script>window.onload = function() { window.print(); };</script>\n")
	}
	bw.WriteString("</head>\n<body>\n<pre>")
	var ls editor.LineStyles
	doc.Buffer().ForEachLine(0, func(line int, text []byte) bool {
		if len(text) <= 64<<10 {
			doc.LineStyles(line, text, &ls)
		} else {
			ls = editor.LineStyles{}
		}
		runs := ls.Runs
		pos := 0
		for k, r := range runs {
			start := int(r.Start)
			if start > pos {
				bw.WriteString(html.EscapeString(string(text[pos:start])))
			}
			end := len(text)
			if k+1 < len(runs) {
				end = int(runs[k+1].Start)
			}
			if end > start {
				fmt.Fprintf(bw, "<span class=\"s%d\">%s</span>", r.Style, html.EscapeString(string(text[start:end])))
			}
			pos = end
		}
		if pos < len(text) {
			bw.WriteString(html.EscapeString(string(text[pos:])))
		}
		bw.WriteByte('\n')
		return true
	})
	bw.WriteString("</pre>\n</body>\n</html>\n")
	return bw.Flush()
}

// exportHTML saves the document as HTML with the syntax highlighting
func (c *MainForm) exportHTML() {
	d := c.curDoc()
	if d == nil {
		return
	}
	name := d.Name() + ".html"
	write := func(path string) {
		f, err := os.Create(path)
		if err != nil {
			c.showError(err)
			return
		}
		err = writeHTML(f, d, c.view().Scheme(), c.view().Options().TabSize, false)
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			c.showError(err)
			return
		}
		c.toast(T().SavedTo(path))
	}
	c.Form().ShowSaveFileDialog(ui.SaveFileDialogOptions{Title: T().MenuExportHTML, DefaultDirectory: c.startDir(), DefaultFileName: name,
		Filters: []ui.FileDialogFilter{{DisplayName: "HTML", Patterns: []string{"*.html", "*.htm"}}}},
		func(path string, err error) {
			if err != nil && err == ui.ErrNoFileDialog {
				c.askPath(T().MenuExportHTML, filepath.Join(c.startDir(), name), write)
				return
			}
			if err == nil && path != "" {
				write(path)
			}
		})
}

// printDoc prints the document through the browser: it is written as HTML
// in the light colors and opened with the print dialog
func (c *MainForm) printDoc() {
	d := c.curDoc()
	if d == nil {
		return
	}
	f, err := os.CreateTemp("", "altnotepad-print-*.html")
	if err != nil {
		c.showError(err)
		return
	}
	err = writeHTML(f, d, editor.SchemeByID("default"), c.view().Options().TabSize, true)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = app.OpenURL("file://" + filepath.ToSlash(f.Name()))
	}
	if err != nil {
		c.showError(err)
	}
}
