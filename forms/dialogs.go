package forms

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ipoluianov/altnotepad/app"
	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// dialogButtons adds a row of buttons at the right of the panel
func dialogButtons(panel *ui.Panel, buttons ...*ui.Button) {
	panel.AddWidget(0, 0, ui.NewHSpacer())
	for i, b := range buttons {
		panel.AddWidget(0, i+1, b)
	}
}

// ---- Save changes ----

const (
	answerCancel = iota
	answerYes
	answerNo
)

// SaveChangesDialog asks whether to save a document: Yes, No or Cancel
type SaveChangesDialog struct {
	ui.DialogContent
	onAnswer func(answer int)
	answered bool
}

// showSaveChanges asks to save the document named name
func showSaveChanges(parent ui.Widgeter, name string, onAnswer func(answer int)) {
	var c SaveChangesDialog
	c.InitWidget()
	c.onAnswer = onAnswer
	lbl := ui.NewLabel(T().SaveChangesAsk(name))
	c.AddWidget(0, 0, lbl)
	c.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	c.AddWidget(2, 0, buttons)
	yes := ui.NewButton(T().Save)
	no := ui.NewButton(T().DontSave)
	cancel := ui.NewButton(ui.UIText().Cancel)
	answer := func(a int) func() {
		return func() {
			c.answered = true
			c.Form().Close()
			c.RunInParent(func() { onAnswer(a) })
		}
	}
	yes.SetOnClick(answer(answerYes))
	no.SetOnClick(answer(answerNo))
	cancel.SetOnClick(answer(answerCancel))
	dialogButtons(buttons, yes, no, cancel)
	c.OnDialogReject = func() bool {
		if !c.answered {
			c.answered = true
			c.RunInParent(func() { onAnswer(answerCancel) })
		}
		return true
	}
	c.OnDialogShow = func() {
		c.Form().SetTitle(T().SaveTitle)
		w, _, _ := ui.MeasureText(ui.ThemeFontFamily(), ui.ThemeFontSize(), lbl.Text())
		c.Form().SetSize(max(420, min(w+60, 800)), 150)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(yes)
		c.Form().SetCancelButton(cancel)
		yes.Focus()
	}
	parent.(interface{ ShowDialog(ui.Widgeter) }).ShowDialog(&c)
}

// ---- Input ----

// InputDialog asks for a line of text
type InputDialog struct {
	ui.DialogContent
	txt *ui.TextBox
}

func NewInputDialog(title, label, value string, onAccept func(text string)) *InputDialog {
	var c InputDialog
	c.InitWidget()
	content := ui.NewPanel()
	c.AddWidget(0, 0, content)
	c.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	c.AddWidget(2, 0, buttons)
	content.AddWidget(0, 0, ui.NewLabel(label))
	c.txt = ui.NewTextBox()
	c.txt.SetText(value)
	content.AddWidget(1, 0, c.txt)
	ok := ui.NewButton(ui.UIText().OK)
	cancel := ui.NewButton(ui.UIText().Cancel)
	ok.SetOnClick(func() {
		text := c.txt.Text()
		c.Form().Close()
		c.RunInParent(func() { onAccept(text) })
	})
	cancel.SetOnClick(func() { c.Form().Close() })
	dialogButtons(buttons, ok, cancel)
	c.OnDialogShow = func() {
		c.Form().SetTitle(title)
		c.Form().SetSize(460, 160)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(ok)
		c.Form().SetCancelButton(cancel)
		c.txt.Focus()
		c.txt.SelectAllText()
	}
	return &c
}

// ---- Go to ----

// gotoDialog asks for a line or an offset to go to
func (c *MainForm) gotoDialog() {
	v := c.view()
	doc := v.Document()
	var d ui.DialogContent
	d.InitWidget()
	content := ui.NewPanel()
	d.AddWidget(0, 0, content)
	d.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	d.AddWidget(2, 0, buttons)

	rbLine := ui.NewRadioButton(T().GotoLine)
	rbOffset := ui.NewRadioButton(T().GotoOffset)
	rbLine.SetChecked(true)
	modes := ui.NewPanel()
	modes.AddWidget(0, 0, rbLine)
	modes.AddWidget(0, 1, rbOffset)
	modes.AddWidget(0, 2, ui.NewHSpacer())
	content.AddWidget(0, 0, modes)
	line, _ := v.CaretLineCol()
	info := ui.NewLabel("")
	content.AddWidget(1, 0, info)
	txt := ui.NewTextBox()
	content.AddWidget(2, 0, txt)
	update := func() {
		if rbOffset.Checked() {
			info.SetText(T().GotoOffsetInfo(v.Caret(), doc.Len()))
		} else {
			info.SetText(T().GotoLineInfo(line, doc.LineCount()))
		}
	}
	rbLine.SetOnStateChanged(func(*ui.RadioButton, bool) { update() })
	rbOffset.SetOnStateChanged(func(*ui.RadioButton, bool) { update() })
	update()
	ok := ui.NewButton(T().Go)
	cancel := ui.NewButton(ui.UIText().Cancel)
	ok.SetOnClick(func() {
		n, err := strconv.Atoi(strings.TrimSpace(txt.Text()))
		offset := rbOffset.Checked()
		d.Form().Close()
		if err != nil {
			return
		}
		d.RunInParent(func() {
			if offset {
				v.GotoPos(max(0, min(n, doc.Len())))
			} else {
				v.GotoLine(n - 1)
			}
			v.Focus()
		})
	})
	cancel.SetOnClick(func() { d.Form().Close() })
	dialogButtons(buttons, ok, cancel)
	d.OnDialogShow = func() {
		d.Form().SetTitle(T().GotoTitle)
		d.Form().SetSize(400, 190)
		d.Form().MoveToCenterOfParent()
		d.Form().SetAcceptButton(ok)
		d.Form().SetCancelButton(cancel)
		txt.Focus()
	}
	c.ShowDialog(&d)
}

// ---- Column editor ----

// columnEditor inserts a text or a sequence of numbers at the caret column
// of each line of the selection (Edit > Column Editor)
func (c *MainForm) columnEditor() {
	v := c.view()
	var d ui.DialogContent
	d.InitWidget()
	content := ui.NewPanel()
	d.AddWidget(0, 0, content)
	d.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	d.AddWidget(2, 0, buttons)

	rbText := ui.NewRadioButton(T().ColumnText)
	rbNum := ui.NewRadioButton(T().ColumnNumbers)
	rbText.SetChecked(true)
	txt := ui.NewTextBox()
	numInit := ui.NewTextBox()
	numInit.SetText("1")
	numStep := ui.NewTextBox()
	numStep.SetText("1")
	numRepeat := ui.NewTextBox()
	numRepeat.SetText("1")
	cmbFormat := ui.NewComboBox()
	for _, f := range []string{"Dec", "Hex", "Oct", "Bin"} {
		cmbFormat.AddItem(f, f)
	}
	cmbFormat.SetSelectedIndex(0)
	cmbLeading := ui.NewComboBox()
	cmbLeading.AddItem(T().LeadingNone, "")
	cmbLeading.AddItem(T().LeadingZeros, "0")
	cmbLeading.AddItem(T().LeadingSpaces, " ")
	cmbLeading.SetSelectedIndex(0)

	content.AddWidget(0, 0, rbText)
	content.AddWidget(0, 1, txt)
	content.AddWidget(1, 0, rbNum)
	grid := ui.NewPanel()
	grid.AddWidget(0, 0, ui.NewLabel(T().ColumnInitial))
	grid.AddWidget(0, 1, numInit)
	grid.AddWidget(1, 0, ui.NewLabel(T().ColumnIncrease))
	grid.AddWidget(1, 1, numStep)
	grid.AddWidget(2, 0, ui.NewLabel(T().ColumnRepeat))
	grid.AddWidget(2, 1, numRepeat)
	grid.AddWidget(3, 0, ui.NewLabel(T().ColumnFormat))
	grid.AddWidget(3, 1, cmbFormat)
	grid.AddWidget(4, 0, ui.NewLabel(T().ColumnLeading))
	grid.AddWidget(4, 1, cmbLeading)
	content.AddWidget(2, 1, grid)

	ok := ui.NewButton(ui.UIText().OK)
	cancel := ui.NewButton(ui.UIText().Cancel)
	ok.SetOnClick(func() {
		isText := rbText.Checked()
		text := txt.Text()
		init, _ := strconv.ParseInt(strings.TrimSpace(numInit.Text()), 10, 64)
		step, _ := strconv.ParseInt(strings.TrimSpace(numStep.Text()), 10, 64)
		repeat, _ := strconv.Atoi(strings.TrimSpace(numRepeat.Text()))
		format, _ := cmbFormat.SelectedItemData().(string)
		leading, _ := cmbLeading.SelectedItemData().(string)
		d.Form().Close()
		d.RunInParent(func() {
			c.insertColumn(v, isText, text, init, step, max(1, repeat), format, leading)
		})
	})
	cancel.SetOnClick(func() { d.Form().Close() })
	dialogButtons(buttons, ok, cancel)
	d.OnDialogShow = func() {
		d.Form().SetTitle(T().ColumnTitle)
		d.Form().SetSize(460, 330)
		d.Form().MoveToCenterOfParent()
		d.Form().SetAcceptButton(ok)
		d.Form().SetCancelButton(cancel)
		txt.Focus()
	}
	c.ShowDialog(&d)
}

// insertColumn inserts the text or the numbers at the column of the caret
// on the lines of the selection (from the caret line to the end without one)
func (c *MainForm) insertColumn(v *editor.View, isText bool, text string, init, step int64, repeat int, format, leading string) {
	doc := v.Document()
	sels := v.Selections()
	first, last := doc.LineCount(), 0
	_, cc := v.CaretLineCol()
	col := cc - 1
	for _, s := range sels {
		l1, l2 := doc.LineOfOffset(s.Start()), doc.LineOfOffset(s.End())
		first, last = min(first, l1), max(last, l2)
		if v.IsRectSelection() {
			col = min(col, colOfSel(v, s))
		}
	}
	if len(sels) == 1 && sels[0].Start() == sels[0].End() {
		last = doc.LineCount() - 1
	}
	numbers := make([]string, 0, last-first+1)
	width := 0
	for i := 0; i <= last-first; i++ {
		n := init + step*int64(i/repeat)
		var s string
		switch format {
		case "Hex":
			s = strings.ToUpper(strconv.FormatInt(n, 16))
		case "Oct":
			s = strconv.FormatInt(n, 8)
		case "Bin":
			s = strconv.FormatInt(n, 2)
		default:
			s = strconv.FormatInt(n, 10)
		}
		numbers = append(numbers, s)
		width = max(width, len(s))
	}
	var edits []editor.Edit
	for l := first; l <= last; l++ {
		ins := text
		if !isText {
			ins = numbers[l-first]
			if leading != "" {
				ins = strings.Repeat(leading, width-len(ins)) + ins
			}
		}
		pos, beyond := editor.PosOfColumn(doc, l, col, v.Options().TabSize)
		edits = append(edits, editor.Edit{Pos: pos, Text: []byte(strings.Repeat(" ", beyond) + ins)})
	}
	doc.BeginAction()
	doc.ApplyEdits(edits)
	doc.EndAction()
	v.Refresh()
}

func colOfSel(v *editor.View, s editor.Sel) int {
	return editor.ColumnOf(v.Document(), s.Start(), v.Options().TabSize) + min(s.AnchorVS, s.CaretVS)
}

// ---- Summary ----

func (c *MainForm) showSummary() {
	d := c.curDoc()
	if d == nil {
		return
	}
	doc := d.doc
	var b strings.Builder
	if !d.IsUntitled() {
		fmt.Fprintf(&b, "%s: %s\n", T().FullPath, d.Path())
		if st, err := os.Stat(d.Path()); err == nil {
			fmt.Fprintf(&b, "%s: %s\n", T().Modified, st.ModTime().Format("2006-01-02 15:04:05"))
		}
		b.WriteString("\n")
	}
	chars, words, nonBlank := 0, 0, 0
	inWord := false
	doc.Buffer().ForEachPiece(0, doc.Len(), func(p []byte) bool {
		for len(p) > 0 {
			r, size := utf8.DecodeRune(p)
			p = p[size:]
			if r != '\n' {
				chars++
			}
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				if !inWord {
					words++
					inWord = true
				}
			} else {
				inWord = false
			}
		}
		return true
	})
	doc.Buffer().ForEachLine(0, func(line int, text []byte) bool {
		if len(strings.TrimSpace(string(text))) > 0 {
			nonBlank++
		}
		return true
	})
	fmt.Fprintf(&b, "%s: %d\n", T().SummaryChars, chars)
	fmt.Fprintf(&b, "%s: %d\n", T().SummaryWords, words)
	fmt.Fprintf(&b, "%s: %d\n", T().SummaryLines, doc.LineCount())
	fmt.Fprintf(&b, "%s: %d\n", T().SummaryNonBlankLines, nonBlank)
	fmt.Fprintf(&b, "%s: %d %s\n", T().SummaryLength, doc.Len(), T().Bytes)
	sel := c.view().AllSelectedText()
	if sel != "" {
		fmt.Fprintf(&b, "%s: %d\n", T().SummarySelected, utf8.RuneCountInString(sel))
	}
	ui.ShowMessageBox(c, T().MenuSummary, b.String())
}

// ---- Run ----

// runDialog runs a command with the variables of the current file
func (c *MainForm) runDialog() {
	hist := config.GetHistory().RunCommands
	var d ui.DialogContent
	d.InitWidget()
	content := ui.NewPanel()
	d.AddWidget(0, 0, content)
	d.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	d.AddWidget(2, 0, buttons)
	content.AddWidget(0, 0, ui.NewLabel(T().RunLabel))
	cmb := ui.NewEditableComboBox()
	cmb.SetItems(hist)
	if len(hist) > 0 {
		cmb.SetText(hist[0])
	}
	content.AddWidget(1, 0, cmb)
	vars := ui.NewLabel("$(FULL_CURRENT_PATH) $(CURRENT_DIRECTORY) $(FILE_NAME) $(NAME_PART) $(EXT_PART) $(CURRENT_WORD) $(CURRENT_LINE)")
	vars.SetForegroundColor(colorSecondaryText.get())
	content.AddWidget(2, 0, vars)
	ok := ui.NewButton(T().Run)
	cancel := ui.NewButton(ui.UIText().Cancel)
	ok.SetOnClick(func() {
		cmdline := strings.TrimSpace(cmb.Text())
		d.Form().Close()
		if cmdline == "" {
			return
		}
		d.RunInParent(func() {
			config.UpdateHistory(func(h *config.History) { h.RunCommands = config.PushFront(h.RunCommands, cmdline, 20) })
			if err := app.RunCommand(c.expandVars(cmdline)); err != nil {
				c.showError(err)
			}
		})
	})
	cancel.SetOnClick(func() { d.Form().Close() })
	dialogButtons(buttons, ok, cancel)
	d.OnDialogShow = func() {
		d.Form().SetTitle(T().RunTitle)
		d.Form().SetSize(560, 180)
		d.Form().MoveToCenterOfParent()
		d.Form().SetAcceptButton(ok)
		d.Form().SetCancelButton(cancel)
		cmb.Focus()
	}
	c.ShowDialog(&d)
}

// expandVars replaces the variables of the Run command
func (c *MainForm) expandVars(s string) string {
	path := c.curPath()
	v := c.view()
	word := strings.TrimSpace(v.SelectedText())
	line, _ := v.CaretLineCol()
	ext := filepath.Ext(path)
	r := strings.NewReplacer(
		"$(FULL_CURRENT_PATH)", path,
		"$(CURRENT_DIRECTORY)", filepath.Dir(path),
		"$(FILE_NAME)", filepath.Base(path),
		"$(NAME_PART)", strings.TrimSuffix(filepath.Base(path), ext),
		"$(EXT_PART)", ext,
		"$(CURRENT_WORD)", word,
		"$(CURRENT_LINE)", strconv.Itoa(line),
		"$(NPP_DIRECTORY)", filepath.Dir(os.Args[0]),
	)
	return r.Replace(s)
}

// ---- Windows ----

// windowsDialog lists the open documents to switch to or close
func (c *MainForm) windowsDialog() {
	var d ui.DialogContent
	d.InitWidget()
	table := ui.NewTreeView()
	table.SetColumnCount(3)
	table.SetColumnName(0, T().Name)
	table.SetColumnName(1, T().FullPath)
	table.SetColumnName(2, T().State)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 420)
	table.SetMultiselect(true)
	for _, doc := range c.docs {
		n := table.AddNode(nil, doc.Name())
		n.SetText(1, doc.FullName())
		if doc.IsModified() {
			n.SetText(2, T().ModifiedState)
		}
		n.SetData(doc)
		n.SetReadOnly(0, true)
		n.SetReadOnly(1, true)
		n.SetReadOnly(2, true)
		if doc == c.curDoc() {
			table.SetCurrentNode(n)
		}
	}
	d.AddWidget(0, 0, table)
	buttons := ui.NewPanel()
	d.AddWidget(1, 0, buttons)
	activate := ui.NewButton(T().Activate)
	closeBtn := ui.NewButton(T().CloseWindows)
	cancel := ui.NewButton(T().Close)
	activate.SetOnClick(func() {
		n := table.CurrentNode()
		d.Form().Close()
		if n != nil {
			doc := n.Data().(*Doc)
			d.RunInParent(func() { c.showDoc(doc) })
		}
	})
	table.SetOnNodeActivated(func(n *ui.TreeNode) { activate.Push() })
	closeBtn.SetOnClick(func() {
		var list []*Doc
		for _, n := range table.SelectedNodes() {
			list = append(list, n.Data().(*Doc))
		}
		d.Form().Close()
		d.RunInParent(func() {
			for _, doc := range list {
				for _, p := range c.panes {
					if p.indexOf(doc) >= 0 {
						c.closeDoc(p, doc, nil)
					}
				}
			}
		})
	})
	cancel.SetOnClick(func() { d.Form().Close() })
	dialogButtons(buttons, activate, closeBtn, cancel)
	d.OnDialogShow = func() {
		d.Form().SetTitle(T().MenuWindows)
		d.Form().SetSize(760, 440)
		d.Form().MoveToCenterOfParent()
		d.Form().SetCancelButton(cancel)
		table.Focus()
	}
	c.ShowDialog(&d)
}

// ---- Shortcuts ----

// showShortcuts lists the keyboard shortcuts
func (c *MainForm) showShortcuts() {
	var d ui.DialogContent
	d.InitWidget()
	table := ui.NewTreeView()
	table.SetColumnCount(2)
	table.SetColumnName(0, T().Command)
	table.SetColumnName(1, T().Shortcut)
	table.SetColumnWidth(0, 420)
	var list []*Command
	for _, cmd := range c.commands {
		if len(cmd.Keys) > 0 {
			list = append(list, cmd)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	for _, cmd := range list {
		keys := make([]string, len(cmd.Keys))
		for i, k := range cmd.Keys {
			keys[i] = k.String()
		}
		n := table.AddNode(nil, cmd.Text())
		n.SetText(1, strings.Join(keys, ", "))
		n.SetReadOnly(0, true)
		n.SetReadOnly(1, true)
	}
	d.AddWidget(0, 0, table)
	buttons := ui.NewPanel()
	d.AddWidget(1, 0, buttons)
	closeBtn := ui.NewButton(T().Close)
	closeBtn.SetOnClick(func() { d.Form().Close() })
	dialogButtons(buttons, closeBtn)
	d.OnDialogShow = func() {
		d.Form().SetTitle(T().MenuShortcuts)
		d.Form().SetSize(640, 560)
		d.Form().MoveToCenterOfParent()
		d.Form().SetCancelButton(closeBtn)
		d.Form().SetAcceptButton(closeBtn)
	}
	c.ShowDialog(&d)
}

// ---- About ----

type AboutDialog struct {
	ui.DialogContent
}

func NewAboutDialog() *AboutDialog {
	var c AboutDialog
	c.InitWidget()
	content := ui.NewPanel()
	c.AddWidget(0, 0, content)
	c.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	c.AddWidget(2, 0, buttons)

	lblName := ui.NewLabel(app.DisplayName)
	lblName.SetFontSize(24)
	lines := []*ui.Label{lblName,
		ui.NewLabel(T().Version + " " + app.Version),
		ui.NewLabel(T().AboutDescription),
		ui.NewLabel(T().Author + " " + app.Author),
		ui.NewLabel(app.Copyright()),
		ui.NewLabel(T().License + " " + app.License),
		ui.NewLabel(app.Website),
	}
	for i, l := range lines {
		l.SetTextAlign(ui.HAlignCenter)
		content.AddWidget(i, 0, l)
	}
	website := ui.NewButton(T().VisitWebsite)
	website.SetOnClick(func() {
		if err := app.OpenSiteURL(app.Website, "about_dialog"); err != nil {
			ui.ShowMessageBox(&c, T().Error, err.Error())
		}
	})
	closeBtn := ui.NewButton(T().Close)
	closeBtn.SetOnClick(func() { c.Form().Close() })
	dialogButtons(buttons, website, closeBtn)
	c.OnDialogShow = func() {
		c.Form().SetTitle(T().AboutTitle(app.DisplayName))
		c.Form().SetSize(440, 340)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(closeBtn)
		c.Form().SetCancelButton(closeBtn)
	}
	return &c
}

// openDocs opens the docs on the site; campaign tells which place in the app the visit came from
func openDocs(parent ui.Widgeter, campaign string) {
	if err := app.OpenSiteURL(app.DocsURL, campaign); err != nil {
		ui.ShowMessageBox(parent, T().Error, err.Error())
	}
}
