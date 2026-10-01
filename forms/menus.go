package forms

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/altnotepad/app"
	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// ed runs an editor command in the active view
func (c *MainForm) ed(cmd string) func() {
	return func() {
		c.view().Exec(cmd, "")
		c.status.refresh()
	}
}

// edArg runs an editor command with an argument in the active view
func (c *MainForm) edArg(cmd, arg string) func() {
	return func() {
		c.view().Exec(cmd, arg)
		c.status.refresh()
	}
}

// registerCommands makes the commands of the menus and the shortcuts
func (c *MainForm) registerCommands() {
	t := T
	// ---- File ----
	c.cmd("file.new", func() string { return t().MenuNew }, "Ctrl+N", c.newDocument)
	c.cmd("file.open", func() string { return t().MenuOpen }, "Ctrl+O", c.openDialog)
	c.cmd("file.openFolder", func() string { return t().MenuOpenContainingFolder }, "", func() {
		if p := c.curPath(); p != "" {
			app.OpenInFileManager(p)
		}
	})
	c.cmd("file.openDefault", func() string { return t().MenuOpenDefaultViewer }, "", func() {
		if p := c.curPath(); p != "" {
			app.OpenURL(p)
		}
	})
	c.cmd("file.workspace", func() string { return t().MenuOpenFolderWorkspace }, "", c.workspace.chooseFolder)
	c.cmd("file.reload", func() string { return t().MenuReload }, "", func() { c.reload(c.curDoc(), true) })
	c.cmd("file.save", func() string { return t().MenuSave }, "Ctrl+S", func() { c.save(c.curDoc(), nil) })
	c.cmd("file.saveAs", func() string { return t().MenuSaveAs }, "Ctrl+Alt+S", func() { c.saveAs(c.curDoc(), false, nil) })
	c.cmd("file.saveCopy", func() string { return t().MenuSaveCopy }, "", func() { c.saveAs(c.curDoc(), true, nil) })
	c.cmd("file.saveAll", func() string { return t().MenuSaveAll }, "Ctrl+Shift+S", c.saveAll)
	c.cmd("file.rename", func() string { return t().MenuRename }, "", func() { c.renameDoc(c.curDoc()) })
	c.cmd("file.close", func() string { return t().MenuClose }, "Ctrl+W | Ctrl+F4", func() {
		if d := c.curDoc(); d != nil {
			c.closeDoc(c.activePane(), d, nil)
		}
	})
	c.cmd("file.closeAll", func() string { return t().MenuCloseAll }, "Ctrl+Shift+W", func() { c.closeAll(nil) })
	c.cmd("file.closeOthers", func() string { return t().MenuCloseAllButActive }, "", func() {
		cur := c.curDoc()
		c.closeWhere(func(i int, d *Doc) bool { return d != cur })
	})
	c.cmd("file.closeLeft", func() string { return t().MenuCloseLeft }, "", func() {
		cur := c.activePane().cur
		c.closeWhere(func(i int, d *Doc) bool { return i < cur })
	})
	c.cmd("file.closeRight", func() string { return t().MenuCloseRight }, "", func() {
		cur := c.activePane().cur
		c.closeWhere(func(i int, d *Doc) bool { return i > cur })
	})
	c.cmd("file.closeUnchanged", func() string { return t().MenuCloseUnchanged }, "", func() {
		c.closeWhere(func(i int, d *Doc) bool { return !d.IsModified() })
	})
	c.cmd("file.deleteFile", func() string { return t().MenuDeleteFile }, "", c.deleteCurrentFile)
	c.cmd("file.loadSession", func() string { return t().MenuLoadSession }, "", c.loadSessionDialog)
	c.cmd("file.saveSession", func() string { return t().MenuSaveSession }, "", c.saveSessionDialog)
	c.cmd("file.print", func() string { return t().MenuPrint }, "Ctrl+P", c.printDoc)
	c.cmd("file.exportHTML", func() string { return t().MenuExportHTML }, "", c.exportHTML)
	c.cmd("file.restoreClosed", func() string { return t().MenuRestoreClosed }, "Ctrl+Shift+T", c.restoreClosed)
	c.cmd("file.clearRecent", func() string { return t().MenuClearRecent }, "", func() {
		config.UpdateHistory(func(h *config.History) { h.RecentFiles = nil })
	})
	c.cmd("file.exit", func() string { return t().MenuExit }, "Alt+F4", func() {
		if c.RequestExit() {
			c.Form().Close()
		}
	})

	// ---- Edit ----
	c.viewCmd("edit.undo", func() string { return t().MenuUndo }, "Ctrl+Z", c.ed("undo"))
	c.viewCmd("edit.redo", func() string { return t().MenuRedo }, "Ctrl+Y", c.ed("redo"))
	c.viewCmd("edit.cut", func() string { return t().MenuCut }, "Ctrl+X", c.ed("cut"))
	c.viewCmd("edit.copy", func() string { return t().MenuCopy }, "Ctrl+C", c.ed("copy"))
	c.viewCmd("edit.paste", func() string { return t().MenuPaste }, "Ctrl+V", c.ed("paste"))
	c.viewCmd("edit.delete", func() string { return t().MenuDelete }, "Del", c.ed("delete"))
	c.viewCmd("edit.selectAll", func() string { return t().MenuSelectAll }, "Ctrl+A", c.ed("select-all"))
	c.cmd("edit.dateShort", func() string { return t().MenuDateShort }, "", func() {
		c.view().Exec("insert-text", time.Now().Format(t().DateShortLayout))
	})
	c.cmd("edit.dateLong", func() string { return t().MenuDateLong }, "", func() {
		c.view().Exec("insert-text", time.Now().Format(t().DateLongLayout))
	})
	c.cmd("edit.dateISO", func() string { return t().MenuDateISO }, "", func() {
		c.view().Exec("insert-text", time.Now().Format("2006-01-02T15:04:05Z07:00"))
	})
	c.cmd("edit.copyPath", func() string { return t().MenuCopyFullPath }, "", func() { c.copyInfo(0) })
	c.cmd("edit.copyName", func() string { return t().MenuCopyFileName }, "", func() { c.copyInfo(1) })
	c.cmd("edit.copyDir", func() string { return t().MenuCopyDirPath }, "", func() { c.copyInfo(2) })
	c.viewCmd("edit.indent", func() string { return t().MenuIndent }, "Tab", c.ed("indent"))
	c.viewCmd("edit.unindent", func() string { return t().MenuUnindent }, "Shift+Tab", c.ed("unindent"))
	c.cmd("edit.upper", func() string { return t().MenuUpperCase }, "Ctrl+Shift+U", c.ed("case-upper"))
	c.cmd("edit.lower", func() string { return t().MenuLowerCase }, "Ctrl+U", c.ed("case-lower"))
	c.cmd("edit.proper", func() string { return t().MenuProperCase }, "Alt+U", c.ed("case-proper"))
	c.cmd("edit.properBlend", func() string { return t().MenuProperCaseBlend }, "Alt+Shift+U", c.ed("case-proper-blend"))
	c.cmd("edit.sentence", func() string { return t().MenuSentenceCase }, "Ctrl+Alt+U", c.ed("case-sentence"))
	c.cmd("edit.sentenceBlend", func() string { return t().MenuSentenceCaseBlend }, "Ctrl+Alt+Shift+U", c.ed("case-sentence-blend"))
	c.cmd("edit.invert", func() string { return t().MenuInvertCase }, "Ctrl+Alt+I", c.ed("case-invert"))
	c.cmd("edit.random", func() string { return t().MenuRandomCase }, "Ctrl+Alt+R", c.ed("case-random"))
	c.cmd("line.duplicate", func() string { return t().MenuDuplicateLine }, "Ctrl+D", c.ed("duplicate"))
	c.cmd("line.removeDups", func() string { return t().MenuRemoveDupLines }, "", c.ed("remove-dup-lines"))
	c.cmd("line.removeConsecutiveDups", func() string { return t().MenuRemoveConsecutiveDupLines }, "", c.ed("remove-consecutive-dup-lines"))
	c.cmd("line.split", func() string { return t().MenuSplitLines }, "Ctrl+I", c.ed("split-lines"))
	c.cmd("line.join", func() string { return t().MenuJoinLines }, "Ctrl+J", c.ed("join-lines"))
	c.viewCmd("line.moveUp", func() string { return t().MenuMoveLineUp }, "Ctrl+Shift+Up", c.ed("move-lines-up"))
	c.viewCmd("line.moveDown", func() string { return t().MenuMoveLineDown }, "Ctrl+Shift+Down", c.ed("move-lines-down"))
	c.cmd("line.delete", func() string { return t().MenuDeleteLine }, "Ctrl+Shift+L", c.ed("delete-lines"))
	c.cmd("line.cut", func() string { return t().MenuCutLine }, "Ctrl+L", c.ed("cut-lines"))
	c.cmd("line.transpose", func() string { return t().MenuTransposeLine }, "Ctrl+T", c.ed("transpose"))
	c.cmd("line.removeEmpty", func() string { return t().MenuRemoveEmptyLines }, "", c.ed("remove-empty-lines"))
	c.cmd("line.removeBlank", func() string { return t().MenuRemoveBlankLines }, "", c.ed("remove-blank-lines"))
	c.cmd("line.insertAbove", func() string { return t().MenuInsertLineAbove }, "Ctrl+Alt+Enter", c.ed("insert-line-above"))
	c.cmd("line.insertBelow", func() string { return t().MenuInsertLineBelow }, "Ctrl+Alt+Shift+Enter", c.ed("insert-line-below"))
	c.cmd("line.reverse", func() string { return t().MenuReverseLines }, "", c.ed("reverse-lines"))
	c.cmd("line.shuffle", func() string { return t().MenuShuffleLines }, "", c.ed("shuffle-lines"))
	sorts := []struct {
		id   string
		text func() string
	}{
		{"asc", func() string { return t().MenuSortAsc }}, {"desc", func() string { return t().MenuSortDesc }},
		{"asc-ci", func() string { return t().MenuSortAscCI }}, {"desc-ci", func() string { return t().MenuSortDescCI }},
		{"int-asc", func() string { return t().MenuSortIntAsc }}, {"int-desc", func() string { return t().MenuSortIntDesc }},
		{"dec-asc", func() string { return t().MenuSortDecAsc }}, {"dec-desc", func() string { return t().MenuSortDecDesc }},
		{"len-asc", func() string { return t().MenuSortLenAsc }}, {"len-desc", func() string { return t().MenuSortLenDesc }},
	}
	for _, s := range sorts {
		c.cmd("sort."+s.id, s.text, "", c.ed("sort-"+s.id))
	}
	c.cmd("comment.toggle", func() string { return t().MenuToggleComment }, "Ctrl+Q", c.ed("toggle-comment"))
	c.cmd("comment.line", func() string { return t().MenuLineComment }, "Ctrl+K", c.ed("comment"))
	c.cmd("comment.uncomment", func() string { return t().MenuLineUncomment }, "Ctrl+Shift+K", c.ed("uncomment"))
	c.cmd("comment.block", func() string { return t().MenuBlockComment }, "Ctrl+Shift+Q", c.ed("block-comment"))
	c.cmd("comment.blockUncomment", func() string { return t().MenuBlockUncomment }, "", c.ed("block-uncomment"))
	c.cmd("edit.complete", func() string { return t().MenuWordCompletion }, "Ctrl+Enter | Ctrl+Space", c.wordCompletion)
	for i, e := range []editor.EOL{editor.EOLCRLF, editor.EOLLF, editor.EOLCR} {
		eol := e
		texts := []func() string{func() string { return t().MenuEOLWindows }, func() string { return t().MenuEOLUnix }, func() string { return t().MenuEOLMac }}
		cm := c.cmd("eol."+strconv.Itoa(i), texts[i], "", func() { c.setEOL(eol) })
		cm.Checked = func() bool { d := c.curDoc(); return d != nil && d.doc.EOL == eol }
	}
	c.cmd("blank.trimTrailing", func() string { return t().MenuTrimTrailing }, "", c.ed("trim-trailing"))
	c.cmd("blank.trimLeading", func() string { return t().MenuTrimLeading }, "", c.ed("trim-leading"))
	c.cmd("blank.trimBoth", func() string { return t().MenuTrimBoth }, "", c.ed("trim-both"))
	c.cmd("blank.eolToSpace", func() string { return t().MenuEOLToSpace }, "", c.ed("eol-to-space"))
	c.cmd("blank.trimEOL", func() string { return t().MenuRemoveBlankEOL }, "", c.ed("trim-eol-to-space"))
	c.cmd("blank.tabToSpace", func() string { return t().MenuTabToSpace }, "", c.ed("tabs-to-spaces"))
	c.cmd("blank.spaceToTabAll", func() string { return t().MenuSpaceToTabAll }, "", c.ed("spaces-to-tabs-all"))
	c.cmd("blank.spaceToTabLeading", func() string { return t().MenuSpaceToTabLeading }, "", c.ed("spaces-to-tabs-leading"))
	c.cmd("edit.columnEditor", func() string { return t().MenuColumnEditor }, "Alt+C", c.columnEditor)
	c.cmd("multi.all", func() string { return t().MenuMultiAll }, "Ctrl+Alt+Shift+D", c.ed("multi-select-all"))
	c.cmd("multi.allCase", func() string { return t().MenuMultiAllCase }, "", c.ed("multi-select-all-case"))
	c.cmd("multi.next", func() string { return t().MenuMultiNext }, "Ctrl+Alt+D", c.ed("multi-select-next"))
	c.cmd("multi.undo", func() string { return t().MenuMultiUndo }, "", c.ed("multi-select-undo"))
	c.cmd("multi.skip", func() string { return t().MenuMultiSkip }, "", c.ed("multi-select-skip"))
	cm := c.cmd("edit.readOnly", func() string { return t().MenuReadOnly }, "", func() {
		if d := c.curDoc(); d != nil {
			d.doc.ReadOnly = !d.doc.ReadOnly
			c.activePane().refreshTabs()
			c.status.refresh()
		}
	})
	cm.Checked = func() bool { d := c.curDoc(); return d != nil && d.doc.ReadOnly }

	// ---- Search ----
	c.cmd("search.find", func() string { return t().MenuFind }, "Ctrl+F", func() { c.find.open(findModeFind) })
	c.cmd("search.findInFiles", func() string { return t().MenuFindInFiles }, "Ctrl+Shift+F", func() { c.find.open(findModeFiles) })
	c.cmd("search.findNext", func() string { return t().MenuFindNext }, "F3", func() { c.find.findNext(false) })
	c.cmd("search.findPrev", func() string { return t().MenuFindPrev }, "Shift+F3", func() { c.find.findNext(true) })
	c.cmd("search.selectFindNext", func() string { return t().MenuSelectFindNext }, "Ctrl+F3", func() { c.find.selectAndFind(false) })
	c.cmd("search.selectFindPrev", func() string { return t().MenuSelectFindPrev }, "Ctrl+Shift+F3", func() { c.find.selectAndFind(true) })
	c.cmd("search.replace", func() string { return t().MenuReplace }, "Ctrl+H", func() { c.find.open(findModeReplace) })
	c.cmd("search.incremental", func() string { return t().MenuIncremental }, "Ctrl+Alt+I", func() { c.find.open(findModeFind) })
	c.cmd("search.mark", func() string { return t().MenuMark }, "Ctrl+M", func() { c.find.open(findModeMark) })
	c.cmd("search.results", func() string { return t().MenuSearchResults }, "F7", c.toggleResults)
	c.cmd("search.nextResult", func() string { return t().MenuNextResult }, "F4", func() { c.results.jump(1) })
	c.cmd("search.prevResult", func() string { return t().MenuPrevResult }, "Shift+F4", func() { c.results.jump(-1) })
	c.cmd("search.goto", func() string { return t().MenuGoTo }, "Ctrl+G", c.gotoDialog)
	c.cmd("search.brace", func() string { return t().MenuGotoBrace }, "Ctrl+B", c.ed("goto-brace"))
	c.cmd("search.selectBrace", func() string { return t().MenuSelectBrace }, "Ctrl+Alt+B", c.ed("select-to-brace"))
	for i := 0; i < 5; i++ {
		n := i
		c.cmd(fmt.Sprintf("style.all%d", n), func() string { return t().MenuUsingStyle(n + 1) }, "", func() { c.styleToken(n, true) })
		c.cmd(fmt.Sprintf("style.one%d", n), func() string { return t().MenuUsingStyle(n + 1) }, "", func() { c.styleToken(n, false) })
		c.cmd(fmt.Sprintf("style.clear%d", n), func() string { return t().MenuClearStyle(n + 1) }, "", func() { c.clearStyle(n) })
		c.cmd(fmt.Sprintf("style.next%d", n), func() string { return t().MenuUsingStyle(n + 1) }, "", func() { c.jumpStyle(n, 1) })
		c.cmd(fmt.Sprintf("style.prev%d", n), func() string { return t().MenuUsingStyle(n + 1) }, "", func() { c.jumpStyle(n, -1) })
	}
	c.cmd("style.clearAll", func() string { return t().MenuClearAllStyles }, "", func() {
		for i := 0; i < editor.NumMarkStyles; i++ {
			c.clearStyle(i)
		}
	})
	c.cmd("change.next", func() string { return t().MenuNextChange }, "Ctrl+Alt+Down", c.ed("next-change"))
	c.cmd("change.prev", func() string { return t().MenuPrevChange }, "Ctrl+Alt+Up", c.ed("prev-change"))
	c.cmd("change.clear", func() string { return t().MenuClearChanges }, "", c.ed("clear-change-history"))
	c.cmd("bookmark.toggle", func() string { return t().MenuToggleBookmark }, "Ctrl+F2", c.ed("toggle-bookmark"))
	c.cmd("bookmark.next", func() string { return t().MenuNextBookmark }, "F2", c.ed("next-bookmark"))
	c.cmd("bookmark.prev", func() string { return t().MenuPrevBookmark }, "Shift+F2", c.ed("prev-bookmark"))
	c.cmd("bookmark.clear", func() string { return t().MenuClearBookmarks }, "", c.ed("clear-bookmarks"))
	c.cmd("bookmark.cut", func() string { return t().MenuCutBookmarked }, "", c.ed("cut-bookmarked-lines"))
	c.cmd("bookmark.copy", func() string { return t().MenuCopyBookmarked }, "", c.ed("copy-bookmarked-lines"))
	c.cmd("bookmark.remove", func() string { return t().MenuRemoveBookmarked }, "", c.ed("remove-bookmarked-lines"))
	c.cmd("bookmark.removeOthers", func() string { return t().MenuRemoveUnbookmarked }, "", c.ed("remove-unbookmarked-lines"))
	c.cmd("bookmark.inverse", func() string { return t().MenuInverseBookmarks }, "", c.ed("inverse-bookmarks"))

	// ---- View ----
	cm = c.cmd("view.onTop", func() string { return t().MenuAlwaysOnTop }, "", func() {
		c.updateSetting(func(s *config.Settings) { s.AlwaysOnTop = !s.AlwaysOnTop })
	})
	cm.Checked = func() bool { return c.settings().AlwaysOnTop }
	cm = c.cmd("view.fullscreen", func() string { return t().MenuFullScreen }, "F11", c.toggleFullScreen)
	cm.Checked = func() bool { return c.fullscreen }
	toggle := func(id string, text func() string, keys string, get func(s *config.Settings) *bool) {
		cm := c.cmd(id, text, keys, func() {
			c.updateSetting(func(s *config.Settings) { p := get(s); *p = !*p })
		})
		cm.Checked = func() bool { s := c.settings(); return *get(&s) }
	}
	toggle("view.spaces", func() string { return t().MenuShowSpaces }, "", func(s *config.Settings) *bool { return &s.ShowSpaces })
	toggle("view.eol", func() string { return t().MenuShowEOL }, "", func(s *config.Settings) *bool { return &s.ShowEOL })
	cm = c.cmd("view.allChars", func() string { return t().MenuShowAllChars }, "", func() {
		c.updateSetting(func(s *config.Settings) {
			on := !(s.ShowSpaces && s.ShowEOL)
			s.ShowSpaces, s.ShowEOL = on, on
		})
	})
	cm.Checked = func() bool { s := c.settings(); return s.ShowSpaces && s.ShowEOL }
	toggle("view.indentGuides", func() string { return t().MenuShowIndentGuides }, "", func(s *config.Settings) *bool { return &s.ShowIndentGuides })
	toggle("view.wrapSymbol", func() string { return t().MenuShowWrapSymbol }, "", func(s *config.Settings) *bool { return &s.ShowWrapSymbol })
	toggle("view.wrap", func() string { return t().MenuWordWrap }, "", func(s *config.Settings) *bool { return &s.WordWrap })
	toggle("view.lineNumbers", func() string { return t().MenuLineNumbers }, "", func(s *config.Settings) *bool { return &s.LineNumbers })
	toggle("view.toolbar", func() string { return t().MenuToolbar }, "", func(s *config.Settings) *bool { return &s.ShowToolbar })
	toggle("view.statusbar", func() string { return t().MenuStatusBar }, "", func(s *config.Settings) *bool { return &s.ShowStatusBar })
	c.cmd("view.zoomIn", func() string { return t().MenuZoomIn }, "Ctrl+Num+ | Ctrl+=", func() { c.setZoom(c.view().Zoom() + 1) })
	c.cmd("view.zoomOut", func() string { return t().MenuZoomOut }, "Ctrl+Num- | Ctrl+-", func() { c.setZoom(c.view().Zoom() - 1) })
	c.cmd("view.zoomReset", func() string { return t().MenuZoomReset }, "Ctrl+Num/ | Ctrl+0", func() { c.setZoom(0) })
	c.cmd("view.moveOther", func() string { return t().MenuMoveToOtherView }, "", func() { c.moveToOtherView(false) })
	c.cmd("view.cloneOther", func() string { return t().MenuCloneToOtherView }, "", func() { c.moveToOtherView(true) })
	c.cmd("view.focusOther", func() string { return t().MenuFocusOtherView }, "F8", c.focusOtherView)
	c.cmd("view.nextTab", func() string { return t().MenuNextTab }, "Ctrl+Tab | Ctrl+PgDn", func() { c.switchTab(1) })
	c.cmd("view.prevTab", func() string { return t().MenuPrevTab }, "Ctrl+Shift+Tab | Ctrl+PgUp", func() { c.switchTab(-1) })
	c.cmd("view.moveTabForward", func() string { return t().MenuMoveTabForward }, "Ctrl+Shift+PgDn", func() { c.moveTab(1) })
	c.cmd("view.moveTabBackward", func() string { return t().MenuMoveTabBackward }, "Ctrl+Shift+PgUp", func() { c.moveTab(-1) })
	for i := 1; i <= 9; i++ {
		n := i
		c.cmd(fmt.Sprintf("view.tab%d", n), func() string { return t().MenuTabN(n) }, fmt.Sprintf("Ctrl+%d", n), func() {
			p := c.activePane()
			if n == 9 {
				p.activate(len(p.docs) - 1)
			} else {
				p.activate(n - 1)
			}
			c.focusEditor()
		})
	}
	c.cmd("fold.all", func() string { return t().MenuFoldAll }, "Alt+0", c.ed("fold-all"))
	c.cmd("fold.none", func() string { return t().MenuUnfoldAll }, "Alt+Shift+0", c.ed("unfold-all"))
	c.cmd("fold.current", func() string { return t().MenuFoldCurrent }, "Ctrl+Alt+F", c.ed("fold-current"))
	c.cmd("fold.uncurrent", func() string { return t().MenuUnfoldCurrent }, "Ctrl+Alt+Shift+F", c.ed("unfold-current"))
	for i := 1; i <= 8; i++ {
		n := strconv.Itoa(i)
		c.cmd("fold.level"+n, func() string { return n }, "Alt+"+n, c.edArg("fold-level", n))
		c.cmd("fold.unlevel"+n, func() string { return n }, "Alt+Shift+"+n, c.edArg("unfold-level", n))
	}
	c.cmd("view.summary", func() string { return t().MenuSummary }, "", c.showSummary)
	cm = c.cmd("view.workspace", func() string { return t().MenuFolderWorkspace }, "", c.toggleWorkspace)
	cm.Checked = func() bool { return c.workspace.IsVisible() }
	cm = c.cmd("view.funcList", func() string { return t().MenuFunctionList }, "", c.toggleFuncList)
	cm.Checked = func() bool { return c.funcList.IsVisible() }
	cm = c.cmd("view.docMap", func() string { return t().MenuDocumentMap }, "", c.toggleDocMap)
	cm.Checked = func() bool { return c.docMap.IsVisible() }
	cm = c.cmd("view.monitoring", func() string { return t().MenuMonitoring }, "", c.toggleMonitoring)
	cm.Checked = func() bool { d := c.curDoc(); return d != nil && d.monitoring }

	// ---- Encoding ----
	encs := []struct {
		id   string
		enc  editor.Encoding
		text func() string
		conv func() string
	}{
		{"ansi", editor.Encoding{ID: "ansi"}, func() string { return t().MenuEncANSI }, func() string { return t().MenuConvANSI }},
		{"utf8", editor.EncodingUTF8, func() string { return t().MenuEncUTF8 }, func() string { return t().MenuConvUTF8 }},
		{"utf8bom", editor.EncodingUTF8BOM, func() string { return t().MenuEncUTF8BOM }, func() string { return t().MenuConvUTF8BOM }},
		{"utf16be", editor.EncodingUTF16BE, func() string { return t().MenuEncUTF16BE }, func() string { return t().MenuConvUTF16BE }},
		{"utf16le", editor.EncodingUTF16LE, func() string { return t().MenuEncUTF16LE }, func() string { return t().MenuConvUTF16LE }},
	}
	for _, e := range encs {
		enc := e.enc
		cm := c.cmd("enc."+e.id, e.text, "", func() { c.encodeIn(realEnc(enc)) })
		cm.Checked = func() bool { d := c.curDoc(); return d != nil && sameEnc(d.doc.Encoding, realEnc(enc)) }
		c.cmd("conv."+e.id, e.conv, "", func() { c.convertTo(realEnc(enc)) })
	}
	for _, cs := range editor.Charsets {
		if cs.Group == "" {
			continue
		}
		enc := editor.Encoding{ID: cs.ID}
		name := cs.Name
		cm := c.cmd("charset."+cs.ID, func() string { return name }, "", func() { c.encodeIn(enc) })
		cm.Checked = func() bool { d := c.curDoc(); return d != nil && d.doc.Encoding.ID == enc.ID }
	}

	// ---- Language ----
	for _, l := range append([]*editor.Language{editor.LangText}, editor.Languages()...) {
		lang := l
		text := func() string { return lang.Name }
		if lang == editor.LangText {
			text = func() string { return t().NormalText }
		}
		cm := c.cmd("lang."+lang.ID, text, "", func() { c.setLanguage(lang) })
		cm.Checked = func() bool { d := c.curDoc(); return d != nil && d.doc.Language() == lang }
	}

	// ---- Settings ----
	c.cmd("settings.preferences", func() string { return t().MenuPreferences }, "", c.ShowSettings)
	for _, sc := range editor.Schemes {
		scheme := sc
		cm := c.cmd("scheme."+scheme.ID, func() string { return scheme.Name }, "", func() {
			c.updateSetting(func(s *config.Settings) { s.Scheme = scheme.ID })
		})
		cm.Checked = func() bool { return editorScheme(c.settings()) == scheme }
	}
	c.cmd("settings.shortcuts", func() string { return t().MenuShortcuts }, "", c.showShortcuts)

	// ---- Tools ----
	for _, h := range hashNames {
		name := h
		c.cmd("hash.text."+name, func() string { return t().MenuHashGenerate(name) }, "", func() { c.hashDialog(name) })
		c.cmd("hash.files."+name, func() string { return t().MenuHashFiles(name) }, "", func() { c.hashFiles(name) })
		c.cmd("hash.sel."+name, func() string { return t().MenuHashSelection(name) }, "", func() { c.hashSelection(name) })
	}
	c.cmd("tools.base64enc", func() string { return t().MenuBase64Encode }, "", func() { c.transformSelection(base64Encode) })
	c.cmd("tools.base64dec", func() string { return t().MenuBase64Decode }, "", func() { c.transformSelection(base64Decode) })
	c.cmd("tools.urlenc", func() string { return t().MenuURLEncode }, "", func() { c.transformSelection(urlEncode) })
	c.cmd("tools.urldec", func() string { return t().MenuURLDecode }, "", func() { c.transformSelection(urlDecode) })
	c.cmd("tools.jsonPretty", func() string { return t().MenuJSONFormat }, "", func() { c.transformSelection(jsonPretty) })
	c.cmd("tools.jsonMin", func() string { return t().MenuJSONMinify }, "", func() { c.transformSelection(jsonMinify) })

	// ---- Macro ----
	c.cmd("macro.record", func() string { return t().MenuStartRecording }, "Ctrl+Shift+R", c.toggleRecording)
	c.cmd("macro.stop", func() string { return t().MenuStopRecording }, "", c.stopRecording)
	c.cmd("macro.play", func() string { return t().MenuPlayback }, "Ctrl+Shift+P", func() { c.playMacro(c.macro, 1) })
	c.cmd("macro.save", func() string { return t().MenuSaveMacro }, "", c.saveMacro)
	c.cmd("macro.runMulti", func() string { return t().MenuRunMacroMulti }, "", c.runMacroMultiple)
	c.cmd("macro.trimSave", func() string { return t().MenuTrimSave }, "", func() {
		c.view().Exec("trim-trailing", "")
		c.save(c.curDoc(), nil)
	})

	// ---- Run ----
	c.cmd("run.run", func() string { return t().MenuRun }, "F5", c.runDialog)
	c.cmd("run.browser", func() string { return t().MenuOpenInBrowser }, "", func() {
		if p := c.curPath(); p != "" {
			app.OpenURL("file://" + filepath.ToSlash(p))
		}
	})
	c.cmd("run.search", func() string { return t().MenuSearchInternet }, "", func() {
		if s := strings.TrimSpace(c.view().SelectedText()); s != "" {
			app.OpenURL("https://www.google.com/search?q=" + urlQueryEscape(s))
		}
	})
	c.cmd("run.openFile", func() string { return t().MenuOpenSelectedFile }, "", c.openSelectedFile)

	// ---- Window ----
	c.cmd("window.list", func() string { return t().MenuWindows }, "", c.windowsDialog)
	c.cmd("window.sortName", func() string { return t().MenuSortTabsByName }, "", func() { c.sortTabs(false) })
	c.cmd("window.sortPath", func() string { return t().MenuSortTabsByPath }, "", func() { c.sortTabs(true) })

	// ---- ? ----
	c.cmd("help.docs", func() string { return t().MenuHelp }, "Shift+F1", func() { openDocs(c, "help_menu") })
	c.cmd("help.site", func() string { return t().MenuHomePage }, "", func() { app.OpenSiteURL(app.Website, "help_menu") })
	c.cmd("help.about", func() string { return t().MenuAbout }, "F1", func() { c.ShowDialog(NewAboutDialog()) })
}

// realEnc resolves the ANSI encoding to its character set
func realEnc(e editor.Encoding) editor.Encoding {
	if e.ID == "ansi" {
		return editor.Encoding{ID: editor.DefaultANSI}
	}
	return e
}

func sameEnc(a, b editor.Encoding) bool {
	return a.ID == b.ID && (a.BOM == b.BOM || !a.IsUnicode())
}

// buildMenuBar makes the main menu
func (c *MainForm) buildMenuBar() *ui.MenuBar {
	t := T
	bar := ui.NewMenuBar()
	add := func(text func() string) *ui.ContextMenu {
		menu := ui.NewContextMenu(c)
		item := bar.AddMenuItem(text(), menu)
		item.SetTextFunc(text)
		menu.SetOnShow(c.updateChecks)
		return menu
	}

	file := add(func() string { return t().MenuFile })
	c.addItems(file, "file.new", "file.open")
	c.subMenu(file, func() string { return t().MenuOpenContainingFolderSub }, "file.openFolder", "file.openDefault", "file.workspace")
	c.addItems(file, "file.reload", "file.save", "file.saveAs", "file.saveCopy", "file.saveAll", "file.rename", "file.close", "file.closeAll")
	c.subMenu(file, func() string { return t().MenuCloseMore }, "file.closeOthers", "file.closeLeft", "file.closeRight", "file.closeUnchanged")
	c.addItems(file, "file.deleteFile", "-", "file.loadSession", "file.saveSession", "-", "file.print", "file.exportHTML", "-")
	c.recentMenu(file)
	c.addItems(file, "file.restoreClosed", "-", "file.exit")

	edit := add(func() string { return t().MenuEdit })
	c.addItems(edit, "edit.undo", "edit.redo", "-", "edit.cut", "edit.copy", "edit.paste", "edit.delete", "edit.selectAll", "-")
	c.subMenu(edit, func() string { return t().MenuInsert }, "edit.dateShort", "edit.dateLong", "edit.dateISO")
	c.subMenu(edit, func() string { return t().MenuCopyToClipboard }, "edit.copyPath", "edit.copyName", "edit.copyDir")
	c.subMenu(edit, func() string { return t().MenuIndentSub }, "edit.indent", "edit.unindent")
	c.subMenu(edit, func() string { return t().MenuConvertCase }, "edit.upper", "edit.lower", "edit.proper", "edit.properBlend", "edit.sentence", "edit.sentenceBlend", "edit.invert", "edit.random")
	c.subMenu(edit, func() string { return t().MenuLineOperations }, "line.duplicate", "line.removeDups", "line.removeConsecutiveDups", "line.split", "line.join",
		"line.moveUp", "line.moveDown", "line.delete", "line.cut", "line.transpose", "line.removeEmpty", "line.removeBlank", "line.insertAbove", "line.insertBelow", "-",
		"line.reverse", "line.shuffle", "-", "sort.asc", "sort.desc", "sort.asc-ci", "sort.desc-ci", "sort.int-asc", "sort.int-desc", "sort.dec-asc", "sort.dec-desc", "sort.len-asc", "sort.len-desc")
	c.subMenu(edit, func() string { return t().MenuMultiSelect }, "multi.all", "multi.allCase", "multi.next", "multi.undo", "multi.skip")
	c.subMenu(edit, func() string { return t().MenuComment }, "comment.toggle", "comment.line", "comment.uncomment", "comment.block", "comment.blockUncomment")
	c.subMenu(edit, func() string { return t().MenuAutoCompletion }, "edit.complete")
	c.subMenu(edit, func() string { return t().MenuEOLConversion }, "eol.0", "eol.1", "eol.2")
	c.subMenu(edit, func() string { return t().MenuBlankOperations }, "blank.trimTrailing", "blank.trimLeading", "blank.trimBoth", "blank.eolToSpace", "blank.trimEOL", "-",
		"blank.tabToSpace", "blank.spaceToTabAll", "blank.spaceToTabLeading")
	c.addItems(edit, "-", "edit.columnEditor", "edit.readOnly")

	search := add(func() string { return t().MenuSearch })
	c.addItems(search, "search.find", "search.findInFiles", "search.findNext", "search.findPrev", "search.selectFindNext", "search.selectFindPrev", "search.replace",
		"search.incremental", "search.results", "search.nextResult", "search.prevResult", "search.goto", "search.brace", "search.selectBrace", "search.mark", "-")
	c.subMenu(search, func() string { return t().MenuStyleAll }, "style.all0", "style.all1", "style.all2", "style.all3", "style.all4")
	c.subMenu(search, func() string { return t().MenuStyleOne }, "style.one0", "style.one1", "style.one2", "style.one3", "style.one4")
	c.subMenu(search, func() string { return t().MenuClearStyleSub }, "style.clear0", "style.clear1", "style.clear2", "style.clear3", "style.clear4", "style.clearAll")
	c.subMenu(search, func() string { return t().MenuJumpUp }, "style.prev0", "style.prev1", "style.prev2", "style.prev3", "style.prev4")
	c.subMenu(search, func() string { return t().MenuJumpDown }, "style.next0", "style.next1", "style.next2", "style.next3", "style.next4")
	search.AddSeparator()
	c.subMenu(search, func() string { return t().MenuChangeHistory }, "change.next", "change.prev", "change.clear")
	c.subMenu(search, func() string { return t().MenuBookmark }, "bookmark.toggle", "bookmark.next", "bookmark.prev", "bookmark.clear", "bookmark.cut",
		"bookmark.copy", "bookmark.remove", "bookmark.removeOthers", "bookmark.inverse")

	view := add(func() string { return t().MenuView })
	c.addItems(view, "view.onTop", "view.fullscreen", "-")
	c.subMenu(view, func() string { return t().MenuShowSymbol }, "view.spaces", "view.eol", "view.allChars", "-", "view.indentGuides", "view.wrapSymbol")
	c.subMenu(view, func() string { return t().MenuZoom }, "view.zoomIn", "view.zoomOut", "view.zoomReset")
	c.subMenu(view, func() string { return t().MenuMoveClone }, "view.moveOther", "view.cloneOther")
	c.subMenu(view, func() string { return t().MenuTab }, "view.tab1", "view.tab2", "view.tab3", "view.tab4", "view.tab5", "view.tab6", "view.tab7", "view.tab8", "view.tab9", "-",
		"view.nextTab", "view.prevTab", "view.moveTabForward", "view.moveTabBackward")
	c.addItems(view, "view.wrap", "view.lineNumbers", "view.focusOther", "-", "fold.all", "fold.none", "fold.current", "fold.uncurrent")
	c.subMenu(view, func() string { return t().MenuFoldLevel }, "fold.level1", "fold.level2", "fold.level3", "fold.level4", "fold.level5", "fold.level6", "fold.level7", "fold.level8")
	c.subMenu(view, func() string { return t().MenuUnfoldLevel }, "fold.unlevel1", "fold.unlevel2", "fold.unlevel3", "fold.unlevel4", "fold.unlevel5", "fold.unlevel6", "fold.unlevel7", "fold.unlevel8")
	c.addItems(view, "-", "view.summary", "-", "view.workspace", "view.funcList", "view.docMap", "-", "view.monitoring", "-", "view.toolbar", "view.statusbar")

	enc := add(func() string { return t().MenuEncoding })
	c.addItems(enc, "enc.ansi", "enc.utf8", "enc.utf8bom", "enc.utf16be", "enc.utf16le")
	sets := ui.NewContextMenu(c)
	var groups []string
	byGroup := map[string][]string{}
	for _, cs := range editor.Charsets {
		if cs.Group == "" {
			continue
		}
		if _, ok := byGroup[cs.Group]; !ok {
			groups = append(groups, cs.Group)
		}
		byGroup[cs.Group] = append(byGroup[cs.Group], "charset."+cs.ID)
	}
	for _, g := range groups {
		group := g
		c.subMenu(sets, func() string { return charsetGroupName(group) }, byGroup[g]...)
	}
	setsItem := enc.AddItemWithSubmenu(t().MenuCharacterSets, sets)
	setsItem.SetTextFunc(func() string { return t().MenuCharacterSets })
	c.addItems(enc, "-", "conv.ansi", "conv.utf8", "conv.utf8bom", "conv.utf16be", "conv.utf16le")

	langMenu := add(func() string { return t().MenuLanguage })
	c.addItems(langMenu, "lang.text")
	letters := map[string][]string{}
	var order []string
	for _, l := range editor.Languages() {
		if l.ID == "text" {
			continue
		}
		letter := strings.ToUpper(string([]rune(l.Name)[0]))
		if _, ok := letters[letter]; !ok {
			order = append(order, letter)
		}
		letters[letter] = append(letters[letter], "lang."+l.ID)
	}
	sort.Strings(order)
	for _, letter := range order {
		l := letter
		c.subMenu(langMenu, func() string { return l }, letters[l]...)
	}

	settingsMenu := add(func() string { return t().MenuSettings })
	c.addItems(settingsMenu, "settings.preferences")
	var schemeIDs []string
	for _, sc := range editor.Schemes {
		schemeIDs = append(schemeIDs, "scheme."+sc.ID)
	}
	c.subMenu(settingsMenu, func() string { return t().MenuColorScheme }, schemeIDs...)
	c.addItems(settingsMenu, "settings.shortcuts")

	tools := add(func() string { return t().MenuTools })
	for _, h := range hashNames {
		name := h
		c.subMenu(tools, func() string { return name }, "hash.text."+name, "hash.files."+name, "hash.sel."+name)
	}
	tools.AddSeparator()
	c.addItems(tools, "tools.base64enc", "tools.base64dec", "tools.urlenc", "tools.urldec", "-", "tools.jsonPretty", "tools.jsonMin")

	macro := add(func() string { return t().MenuMacro })
	c.addItems(macro, "macro.record", "macro.stop", "macro.play", "macro.save", "macro.runMulti", "-", "macro.trimSave")
	c.savedMacrosMenu(macro)

	run := add(func() string { return t().MenuRunMenu })
	c.addItems(run, "run.run", "-", "run.browser", "run.search", "run.openFile")

	window := add(func() string { return t().MenuWindow })
	c.addItems(window, "window.list", "window.sortName", "window.sortPath", "-")
	c.openDocsMenu(window)

	help := add(func() string { return "?" })
	c.addItems(help, "help.docs", "help.site", "settings.shortcuts", "-", "help.about")

	c.menuBar = bar
	return bar
}

// charsetGroupName returns the name of a group of character sets in the language of the application
func charsetGroupName(g string) string {
	if n, ok := T().CharsetGroups[g]; ok {
		return n
	}
	return g
}

// pooledItems adds count hidden items to the menu, shown by fill before the menu is shown
func pooledItems(menu *ui.ContextMenu, count int) []*ui.ContextMenuItem {
	items := make([]*ui.ContextMenuItem, count)
	for i := range items {
		items[i] = menu.AddItem("", nil)
		items[i].SetVisible(false)
	}
	return items
}

// recentMenu adds the submenu of the recent files
func (c *MainForm) recentMenu(file *ui.ContextMenu) {
	sub := ui.NewContextMenu(c)
	items := pooledItems(sub, 30)
	sub.AddSeparator()
	c.addItems(sub, "file.clearRecent")
	sub.SetOnShow(func() {
		recent := config.GetHistory().RecentFiles
		for i, item := range items {
			if i < len(recent) {
				path := recent[i]
				item.SetText(fmt.Sprintf("%d: %s", i+1, path))
				item.OnClick = func() { c.openFile(path, 0) }
				item.SetVisible(true)
			} else {
				item.SetVisible(false)
			}
		}
	})
	item := file.AddItemWithSubmenu(T().MenuRecentFiles, sub)
	item.SetTextFunc(func() string { return T().MenuRecentFiles })
}

// openDocsMenu adds the open documents to the Window menu
func (c *MainForm) openDocsMenu(window *ui.ContextMenu) {
	items := pooledItems(window, 40)
	prev := window
	_ = prev
	onShow := func() {
		c.updateChecks()
		var list []*Doc
		for _, p := range c.panes {
			for _, d := range p.docs {
				dup := false
				for _, x := range list {
					if x == d {
						dup = true
					}
				}
				if !dup {
					list = append(list, d)
				}
			}
		}
		for i, item := range items {
			if i < len(list) {
				d := list[i]
				name := d.FullName()
				if d.IsModified() {
					name = "*" + name
				}
				item.SetText(fmt.Sprintf("%d: %s", i+1, name))
				item.OnClick = func() { c.showDoc(d) }
				item.SetVisible(true)
				if d == c.curDoc() {
					item.SetImage(loadIcon16("check"))
				} else {
					item.SetImage(nil)
				}
			} else {
				item.SetVisible(false)
			}
		}
	}
	window.SetOnShow(onShow)
}

// savedMacrosMenu adds the saved macros to the Macro menu
func (c *MainForm) savedMacrosMenu(macro *ui.ContextMenu) {
	items := pooledItems(macro, 20)
	macro.SetOnShow(func() {
		c.updateChecks()
		macros := config.GetHistory().Macros
		names := make([]string, 0, len(macros))
		for n := range macros {
			names = append(names, n)
		}
		sort.Strings(names)
		for i, item := range items {
			if i < len(names) {
				name := names[i]
				item.SetText(name)
				item.OnClick = func() { c.playSavedMacro(name) }
				item.SetVisible(true)
			} else {
				item.SetVisible(false)
			}
		}
	})
}

// ---- Context menus ----

var editorMenu, tabMenu *ui.ContextMenu

// buildContextMenus makes the context menus of the editors and the tabs
func (c *MainForm) buildContextMenus() {
	editorMenu = ui.NewContextMenu(c)
	c.addItems(editorMenu, "edit.cut", "edit.copy", "edit.paste", "edit.delete", "edit.selectAll", "-")
	c.subMenu(editorMenu, func() string { return T().MenuStyleAll }, "style.all0", "style.all1", "style.all2", "style.all3", "style.all4")
	c.subMenu(editorMenu, func() string { return T().MenuClearStyleSub }, "style.clear0", "style.clear1", "style.clear2", "style.clear3", "style.clear4", "style.clearAll")
	c.addItems(editorMenu, "-", "edit.upper", "edit.lower", "-", "run.openFile", "run.search", "-", "comment.toggle", "comment.block", "-", "bookmark.toggle")
	editorMenu.SetOnShow(c.updateChecks)

	tabMenu = ui.NewContextMenu(c)
	c.addItems(tabMenu, "file.close", "file.closeOthers", "file.closeLeft", "file.closeRight", "file.closeUnchanged", "-",
		"file.save", "file.saveAs", "file.rename", "file.deleteFile", "file.reload", "-",
		"file.openFolder", "edit.copyPath", "edit.copyName", "edit.copyDir", "-", "edit.readOnly", "-", "view.moveOther", "view.cloneOther")
	tabMenu.SetOnShow(c.updateChecks)

	// Attached to the widgets that show them, which handle the right click themselves
	for _, p := range c.panes {
		p.view.SetContextMenu(editorMenu)
		p.tabs.SetContextMenu(tabMenu)
	}
}

// showEditorMenu shows the context menu of the editor at (x, y) of the form
func (c *MainForm) showEditorMenu(x, y int) {
	editorMenu.ShowMenu(x, y)
}

// showTabMenu shows the context menu of a tab at (x, y) of the form
func (c *MainForm) showTabMenu(x, y int) {
	tabMenu.ShowMenu(x, y)
}

// ---- Helpers of the commands ----

func (c *MainForm) setZoom(z int) {
	for _, p := range c.panes {
		p.view.SetZoom(z)
	}
	c.zoomChanged(c.view().Zoom())
}

func (c *MainForm) switchTab(dir int) {
	p := c.activePane()
	if len(p.docs) == 0 {
		return
	}
	p.activate((p.cur + dir + len(p.docs)) % len(p.docs))
	c.focusEditor()
}

func (c *MainForm) moveTab(dir int) {
	p := c.activePane()
	to := p.cur + dir
	if to < 0 || to >= len(p.docs) {
		return
	}
	p.docs[p.cur], p.docs[to] = p.docs[to], p.docs[p.cur]
	p.cur = to
	p.refreshTabs()
}

func (c *MainForm) sortTabs(byPath bool) {
	p := c.activePane()
	cur := p.current()
	sort.SliceStable(p.docs, func(i, j int) bool {
		a, b := p.docs[i].Name(), p.docs[j].Name()
		if byPath {
			a, b = p.docs[i].FullName(), p.docs[j].FullName()
		}
		return strings.ToLower(a) < strings.ToLower(b)
	})
	p.cur = p.indexOf(cur)
	p.refreshTabs()
}

func (c *MainForm) copyInfo(kind int) {
	d := c.curDoc()
	if d == nil {
		return
	}
	text := d.FullName()
	switch kind {
	case 1:
		text = d.Name()
	case 2:
		if d.IsUntitled() {
			return
		}
		text = filepath.Dir(d.Path())
	}
	ui.ClipboardSetText(text)
}

func (c *MainForm) setEOL(eol editor.EOL) {
	d := c.curDoc()
	if d == nil || d.doc.EOL == eol {
		return
	}
	d.doc.EOL = eol
	d.doc.SetModified()
	c.view().Refresh()
	c.status.refresh()
}

// setLanguage sets the language of the active document
func (c *MainForm) setLanguage(l *editor.Language) {
	d := c.curDoc()
	if d == nil {
		return
	}
	d.langSet = true
	d.doc.SetLanguage(l)
	for _, p := range c.panes {
		if p.current() == d {
			c.applyViewOptions(p)
			p.view.Refresh()
		}
	}
	c.status.refresh()
}

// encodeIn reads the text again in the encoding (Encoding > Encode in)
func (c *MainForm) encodeIn(enc editor.Encoding) {
	d := c.curDoc()
	if d == nil {
		return
	}
	old := d.doc.Encoding
	if old.ID == enc.ID && old.ID == "utf-8" {
		if old.BOM != enc.BOM {
			d.doc.Encoding = enc
			d.doc.SetModified()
			c.status.refresh()
		}
		return
	}
	apply := func(lf *editor.LoadedFile, fromDisk bool) {
		st := c.view().State()
		caret, top := st.Caret(), st.TopLine()
		wasModified := d.IsModified()
		d.doc.SetBuffer(lf.Buffer)
		d.doc.Encoding = enc
		if fromDisk {
			d.statFile()
		}
		if wasModified || !fromDisk {
			d.doc.SetModified()
		}
		st.SetPosition(min(caret, d.doc.Len()), min(caret, d.doc.Len()), top)
		c.view().Refresh()
		c.activePane().refreshTabs()
		c.status.refresh()
	}
	if !d.IsUntitled() && !d.IsModified() {
		lf, err := editor.LoadFile(d.Path(), enc, nil)
		if err != nil {
			c.showError(err)
			return
		}
		apply(lf, true)
		return
	}
	lf, err := editor.Reinterpret(d.doc.Buffer(), old, enc)
	if err != nil {
		c.showError(err)
		return
	}
	apply(lf, false)
}

// convertTo changes the encoding the document is saved in (Encoding > Convert to)
func (c *MainForm) convertTo(enc editor.Encoding) {
	d := c.curDoc()
	if d == nil || sameEnc(d.doc.Encoding, enc) {
		return
	}
	d.doc.Encoding = enc
	d.doc.SetModified()
	c.status.refresh()
}

func (c *MainForm) toggleResults() {
	c.results.SetVisible(!c.results.IsVisible())
	if c.results.IsVisible() {
		c.vertSplit.SetSecondSize(max(120, c.resultsH))
	}
	c.updateLayout()
}

func (c *MainForm) toggleWorkspace() {
	if c.workspace.root == "" && !c.workspace.IsVisible() {
		c.workspace.chooseFolder()
		return
	}
	c.workspace.SetVisible(!c.workspace.IsVisible())
	if c.workspace.IsVisible() {
		c.sideSplit.SetFirstSize(max(150, c.sideW))
	}
	c.updateLayout()
}

func (c *MainForm) toggleMonitoring() {
	d := c.curDoc()
	if d == nil || d.IsUntitled() {
		return
	}
	d.monitoring = !d.monitoring
	d.doc.ReadOnly = d.monitoring
	if d.monitoring {
		c.reload(d, false)
		c.view().GotoPos(d.doc.Len())
	}
	c.activePane().refreshTabs()
	c.status.refresh()
}

// deleteCurrentFile deletes the file of the active document and closes it
func (c *MainForm) deleteCurrentFile() {
	d := c.curDoc()
	if d == nil || d.IsUntitled() {
		return
	}
	ui.ShowQuestionMessageBoxYesNo(c, T().DeleteFileTitle, T().DeleteFileAsk(d.Path()), func() {
		if err := os.Remove(d.Path()); err != nil {
			c.showError(err)
			return
		}
		d.doc.SetSavePoint()
		for _, p := range c.panes {
			if p.indexOf(d) >= 0 {
				c.doClose(p, d)
			}
		}
		c.workspace.refreshSoon()
	}, nil)
}

// openSelectedFile opens the file whose name is selected (or under the caret)
func (c *MainForm) openSelectedFile() {
	v := c.view()
	name := strings.TrimSpace(v.SelectedText())
	if name == "" {
		line := v.Document().LineOfOffset(v.Caret())
		name = strings.TrimSpace(string(v.Document().LineText(line)))
	}
	name = strings.Trim(name, "\"'<>")
	if name == "" {
		return
	}
	if !filepath.IsAbs(name) && c.curPath() != "" {
		name = filepath.Join(filepath.Dir(c.curPath()), name)
	}
	if _, err := os.Stat(name); err != nil {
		c.toast(T().FileNotFound(name))
		return
	}
	c.openFile(name, 0)
}
