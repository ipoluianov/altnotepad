package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ipoluianov/altnotepad/app"
	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/forms"
	"github.com/ipoluianov/nui/ui"
)

// parseArgs returns the files of the command line and the line to go to (-n42)
func parseArgs(args []string) (req app.OpenRequest, newInstance bool) {
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "-n") && len(a) > 2:
			req.Line, _ = strconv.Atoi(a[2:])
		case a == "-multiInst" || a == "--new-window":
			newInstance = true
		case strings.HasPrefix(a, "-") && len(a) > 1:
			// Other options of Notepad++ are ignored
		default:
			if abs, err := filepath.Abs(a); err == nil {
				a = abs
			}
			req.Files = append(req.Files, a)
		}
	}
	return req, newInstance
}

func main() {
	config.Init()
	s := config.GetSettings()
	req, newInstance := parseArgs(os.Args[1:])
	if s.SingleInstance && !newInstance && app.SendToRunning(config.ConfigDirectory(), req) {
		return
	}

	forms.SetLanguage(s.Language)
	forms.ApplyTheme(s.Theme)
	forms.ApplyANSI(s)
	ui.SetAppIcon(appIcon())
	form := ui.NewForm()
	mainForm := forms.NewMainForm()
	form.Panel().SetPanelPadding(0)
	form.Panel().AddWidget(0, 0, mainForm)
	form.SetMenuBar(mainForm.BuildMenuBar())
	maximized := mainForm.RestoreWindowState(form)
	form.OnClose = mainForm.RequestExit
	form.SetOnGlobalKeyDown(mainForm.OnGlobalKeyDown)
	form.SetOnLanguageChanged(mainForm.ApplyLanguage)
	form.SetOnFilesDropped(func(files []string, x, y int) { mainForm.OpenFiles(files) })

	mainForm.RestoreSession()
	for _, f := range req.Files {
		mainForm.OpenFile(f, req.Line)
	}
	mainForm.EnsureDocument()

	stop := func() {}
	if s.SingleInstance {
		stop = app.ListenForInstances(config.ConfigDirectory(), func(r app.OpenRequest) {
			form.Invoke(func() {
				for _, f := range r.Files {
					mainForm.OpenFile(f, r.Line)
				}
				form.RequestAttention()
			})
		})
	}
	form.Show()
	form.Invoke(func() {
		if maximized {
			form.Maximize()
		}
		form.SetAlwaysOnTop(config.GetSettings().AlwaysOnTop)
		mainForm.Activate()
	})
	form.Exec()
	stop()
}
