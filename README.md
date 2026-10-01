# altnotepad

A text and source code editor in the spirit of Notepad++, for Linux, Windows
and macOS. Part of the [altbins](https://altbins.pro/) utilities, built on
[nui](https://github.com/ipoluianov/nui).

## Features

- Tabs, two views side by side (move or clone a document to the other view)
- Syntax highlighting and folding for about 70 languages
- Huge files: gigabytes and lines of hundreds of megabytes; files above the
  large file limit (200 MB by default) open without highlighting and folding
- Multi-editing (Ctrl+Click), rectangular selection (Alt+Drag, Alt+Shift+Arrows),
  multi-select of all the occurrences
- Find, Replace, Find in Files, Mark with normal, extended (`\n`, `\t`...) and
  regular expression modes; Find All results panel
- Encodings: UTF-8, UTF-8-BOM, UTF-16 LE/BE, and 40 character sets; EOL conversion
- Sessions: the open files and the unsaved changes are kept between the starts
- Bookmarks, change history, document map, function list, folder as workspace
- Line operations, case conversion, comments, blank operations, column editor
- Macros, word completion, file monitoring (tail -f), hashes, Base64, JSON format
- Color schemes, dark and light themes, English and Russian

## Build

```
go build -o altnotepad .
```

On Linux `libX11` is needed at run time, and `zenity` or `kdialog` for the
file dialogs. Windows: `go build -ldflags="-H=windowsgui"`.
