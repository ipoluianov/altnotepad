package forms

import (
	"image/color"

	"github.com/ipoluianov/altnotepad/config"
	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// The color themes offered in the settings (config.Settings.Theme)
const (
	themeDark  = ""
	themeLight = "light"
)

// themeColors is a color of the application in the dark and the light theme
type themeColors struct {
	dark, light color.RGBA
}

// get returns the color for the current theme
func (c themeColors) get() color.RGBA {
	if ui.IsDarkTheme {
		return c.dark
	}
	return c.light
}

var (
	colorLink          = themeColors{ui.ColorFromHex("#3d8bf2"), ui.ColorFromHex("#1d6fb5")}
	colorSecondaryText = themeColors{ui.ColorFromHex("#8a949a"), ui.ColorFromHex("#7a7a7a")}
	colorModified      = themeColors{ui.ColorFromHex("#ef6c57"), ui.ColorFromHex("#d9412b")}
	colorAccent        = themeColors{ui.ColorFromHex("#3d8bf2"), ui.ColorFromHex("#1d6fb5")}
	colorRecording     = themeColors{ui.ColorFromHex("#ef5350"), ui.ColorFromHex("#d32f2f")}
)

// themedWidgets are recolored when the theme changes
var themeListeners []func()

// ApplyTheme switches the application to the theme of the settings and
// repaints the open windows
func ApplyTheme(theme string) {
	if theme == themeLight {
		ui.ApplyLightTheme()
	} else {
		ui.ApplyDarkTheme()
	}
	for _, setIcon := range themedIcons {
		setIcon()
	}
	for _, f := range themeListeners {
		f()
	}
}

// editorScheme returns the color scheme of the text for the settings
func editorScheme(s config.Settings) *editor.Scheme {
	if s.Scheme != "" {
		return editor.SchemeByID(s.Scheme)
	}
	if s.Theme == themeLight {
		return editor.SchemeByID("default")
	}
	return editor.SchemeByID("obsidian")
}
