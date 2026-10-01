package forms

import (
	"image"
	"strings"

	"github.com/ipoluianov/altnotepad/editor"
	"github.com/ipoluianov/nui/ui"
)

// Shortcut is a key with modifiers
type Shortcut struct {
	Key   ui.Key
	Ctrl  bool
	Shift bool
	Alt   bool
}

var keyNames = map[string]ui.Key{
	"F1": ui.KeyF1, "F2": ui.KeyF2, "F3": ui.KeyF3, "F4": ui.KeyF4, "F5": ui.KeyF5, "F6": ui.KeyF6,
	"F7": ui.KeyF7, "F8": ui.KeyF8, "F9": ui.KeyF9, "F10": ui.KeyF10, "F11": ui.KeyF11, "F12": ui.KeyF12,
	"Tab": ui.KeyTab, "Enter": ui.KeyEnter, "Esc": ui.KeyEsc, "Space": ui.KeySpace,
	"Up": ui.KeyArrowUp, "Down": ui.KeyArrowDown, "Left": ui.KeyArrowLeft, "Right": ui.KeyArrowRight,
	"Home": ui.KeyHome, "End": ui.KeyEnd, "PgUp": ui.KeyPageUp, "PgDn": ui.KeyPageDown,
	"Ins": ui.KeyInsert, "Del": ui.KeyDelete, "Backspace": ui.KeyBackspace,
	"Num+": ui.KeyNumpadPlus, "Num-": ui.KeyNumpadMinus, "Num/": ui.KeyNumpadSlash,
	"=": ui.KeyEqual, "-": ui.KeyMinus, "[": ui.KeyLeftBracket, "]": ui.KeyRightBracket,
	"/": ui.KeySlash, "\\": ui.KeyBackslash, ",": ui.KeyComma, ".": ui.KeyDot, ";": ui.KeySemicolon, "`": ui.KeyGrave,
	"0": ui.Key0, "1": ui.Key1, "2": ui.Key2, "3": ui.Key3, "4": ui.Key4, "5": ui.Key5, "6": ui.Key6, "7": ui.Key7, "8": ui.Key8, "9": ui.Key9,
}

func init() {
	for c := 'A'; c <= 'Z'; c++ {
		keyNames[string(c)] = ui.Key(ui.KeyA + int(c-'A'))
	}
}

// parseShortcut parses "Ctrl+Shift+S"; the zero Shortcut for ""
func parseShortcut(s string) Shortcut {
	var sc Shortcut
	if s == "" {
		return sc
	}
	parts := strings.Split(s, "+")
	// "Ctrl+Num+" has an empty last part
	if strings.HasSuffix(s, "++") || strings.HasSuffix(s, "Num+") {
		parts = append(parts[:len(parts)-2], parts[len(parts)-2]+"+")
	}
	for i, p := range parts {
		if i < len(parts)-1 {
			switch p {
			case "Ctrl":
				sc.Ctrl = true
			case "Shift":
				sc.Shift = true
			case "Alt":
				sc.Alt = true
			}
			continue
		}
		k, ok := keyNames[p]
		if !ok {
			panic("unknown key " + p)
		}
		sc.Key = k
	}
	return sc
}

var keyToName map[ui.Key]string

func (s Shortcut) String() string {
	if s.Key == 0 {
		return ""
	}
	if keyToName == nil {
		keyToName = make(map[ui.Key]string, len(keyNames))
		for n, k := range keyNames {
			keyToName[k] = n
		}
	}
	var b strings.Builder
	if s.Ctrl {
		b.WriteString("Ctrl+")
	}
	if s.Alt {
		b.WriteString("Alt+")
	}
	if s.Shift {
		b.WriteString("Shift+")
	}
	b.WriteString(keyToName[s.Key])
	return b.String()
}

// Command is an action of the menus, the tool bar and the shortcuts
type Command struct {
	ID   string
	Text func() string
	Keys []Shortcut
	Run  func()
	// Checked shows a check mark at the menu item
	Checked func() bool
	// ViewKey: the shortcut is handled by the editor itself, the menu only shows it
	ViewKey bool
}

// menuText is the text of the menu item: the name and the first shortcut
func (c *Command) menuText() string {
	t := c.Text()
	if len(c.Keys) > 0 {
		t += "  (" + c.Keys[0].String() + ")"
	}
	return t
}

// cmd registers a command; keys are like "Ctrl+S", several separated by " | "
func (c *MainForm) cmd(id string, text func() string, keys string, run func()) *Command {
	cmd := &Command{ID: id, Text: text, Run: run}
	if keys != "" {
		for _, k := range strings.Split(keys, " | ") {
			sc := parseShortcut(k)
			cmd.Keys = append(cmd.Keys, sc)
			c.shortcuts[sc] = cmd
		}
	}
	c.commands[id] = cmd
	return cmd
}

// viewCmd registers a command whose shortcut the editor handles itself
func (c *MainForm) viewCmd(id string, text func() string, keys string, run func()) *Command {
	cmd := c.cmd(id, text, "", run)
	cmd.ViewKey = true
	for _, k := range strings.Split(keys, " | ") {
		if k != "" {
			cmd.Keys = append(cmd.Keys, parseShortcut(k))
		}
	}
	return cmd
}

// runCommand runs the command by ID
func (c *MainForm) runCommand(id string) {
	if cmd := c.commands[id]; cmd != nil {
		cmd.Run()
	}
}

// menuCheck updates the check mark of a menu item when its menu is shown
type menuCheck struct {
	item    *ui.ContextMenuItem
	checked func() bool
}

var checkIcon image.Image

// addItem adds the command to the menu
func (c *MainForm) addItem(menu *ui.ContextMenu, id string) *ui.ContextMenuItem {
	cmd := c.commands[id]
	if cmd == nil {
		panic("unknown command " + id)
	}
	item := menu.AddItem(cmd.menuText(), func() { cmd.Run() })
	item.SetTextFunc(cmd.menuText)
	if cmd.Checked != nil {
		c.menuChecks = append(c.menuChecks, menuCheck{item, cmd.Checked})
	}
	return item
}

// addItems adds the commands; "-" adds a separator
func (c *MainForm) addItems(menu *ui.ContextMenu, ids ...string) {
	for _, id := range ids {
		if id == "-" {
			menu.AddSeparator()
			continue
		}
		c.addItem(menu, id)
	}
}

// subMenu adds a submenu with the commands
func (c *MainForm) subMenu(menu *ui.ContextMenu, text func() string, ids ...string) *ui.ContextMenu {
	sub := ui.NewContextMenu(c)
	c.addItems(sub, ids...)
	item := menu.AddItemWithSubmenu(text(), sub)
	item.SetTextFunc(text)
	sub.SetOnShow(c.updateChecks)
	return sub
}

// updateChecks shows the check marks of the menu items
func (c *MainForm) updateChecks() {
	if checkIcon == nil {
		checkIcon = loadIcon16("check")
	}
	for _, mc := range c.menuChecks {
		if mc.checked() {
			mc.item.SetImage(checkIcon)
		} else {
			mc.item.SetImage(nil)
		}
	}
}

// OnGlobalKeyDown runs the command of a shortcut; other keys go on to the focused widget
func (c *MainForm) OnGlobalKeyDown(key ui.Key, mods ui.KeyModifiers) bool {
	editor.TrackModifiers(key, mods, true)
	sc := Shortcut{Key: key, Ctrl: mods.Ctrl || mods.Cmd, Shift: mods.Shift, Alt: mods.Alt}
	if key == ui.KeyEsc && !sc.Ctrl && !sc.Alt && !sc.Shift {
		if c.find.IsVisible() && c.find.hasFocus() {
			c.find.close()
			return true
		}
		if c.fullscreen {
			c.toggleFullScreen()
			return true
		}
		return false
	}
	cmd := c.shortcuts[sc]
	if cmd == nil {
		return false
	}
	cmd.Run()
	return true
}
