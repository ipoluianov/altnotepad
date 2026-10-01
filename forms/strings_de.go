package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func deP(n int, one, many string) string { return i18n.Plural("de", n, one, many) }

var de = Strings{
	Error:      "Fehler",
	NewDocName: "neu",
	NormalText: "Normaler Text",

	MenuFile: "Datei", MenuEdit: "Bearbeiten", MenuSearch: "Suchen", MenuView: "Ansicht", MenuEncoding: "Kodierung", MenuLanguage: "Sprache",
	MenuSettings: "Einstellungen", MenuTools: "Werkzeuge", MenuMacro: "Makro", MenuRunMenu: "Ausführen", MenuWindow: "Fenster",

	MenuNew: "Neu", MenuOpen: "Öffnen...", MenuOpenContainingFolderSub: "Enthaltenden Ordner öffnen",
	MenuOpenContainingFolder: "Dateimanager", MenuOpenDefaultViewer: "Im Standardprogramm öffnen", MenuOpenFolderWorkspace: "Ordner als Arbeitsbereich öffnen...",
	MenuReload: "Neu laden", MenuSave: "Speichern", MenuSaveAs: "Speichern unter...", MenuSaveCopy: "Kopie speichern unter...", MenuSaveAll: "Alle speichern",
	MenuRename: "Umbenennen...", MenuClose: "Schließen", MenuCloseAll: "Alle schließen", MenuCloseMore: "Mehrere Dokumente schließen",
	MenuCloseAllButActive: "Alle außer aktuellem schließen", MenuCloseLeft: "Alle links davon schließen", MenuCloseRight: "Alle rechts davon schließen",
	MenuCloseUnchanged: "Alle unveränderten schließen", MenuDeleteFile: "Vom Datenträger löschen...", MenuLoadSession: "Sitzung laden...",
	MenuSaveSession: "Sitzung speichern...", MenuRecentFiles: "Zuletzt geöffnet", MenuClearRecent: "Liste leeren",
	MenuRestoreClosed: "Zuletzt geschlossene Datei wiederherstellen", MenuExit: "Beenden",
	MenuPrint: "Drucken...", MenuExportHTML: "Als HTML exportieren...", MenuFunctionList: "Funktionsliste", Filter: "Filter", MenuDocumentMap: "Dokumentübersicht",
	MenuChangeHistory: "Änderungsverlauf", MenuNextChange: "Zur nächsten Änderung", MenuPrevChange: "Zur vorherigen Änderung", MenuClearChanges: "Änderungsverlauf löschen",
	SavedTo:         func(path string) string { return "Gespeichert in " + path },
	MenuMultiSelect: "Mehrfachauswahl", MenuMultiAll: "Alle Vorkommen auswählen", MenuMultiAllCase: "Alle Vorkommen auswählen (Groß-/Kleinschreibung)",
	MenuMultiNext: "Nächstes Vorkommen hinzufügen", MenuMultiUndo: "Letzte Auswahl entfernen", MenuMultiSkip: "Überspringen und nächstes hinzufügen",

	MenuUndo: "Rückgängig", MenuRedo: "Wiederholen", MenuCut: "Ausschneiden", MenuCopy: "Kopieren", MenuPaste: "Einfügen", MenuDelete: "Löschen", MenuSelectAll: "Alles auswählen",
	MenuInsert: "Einfügen", MenuDateShort: "Datum und Uhrzeit (kurz)", MenuDateLong: "Datum und Uhrzeit (lang)", MenuDateISO: "Datum und Uhrzeit (ISO 8601)",
	MenuCopyToClipboard: "In die Zwischenablage kopieren", MenuCopyFullPath: "Vollständiger Dateipfad", MenuCopyFileName: "Dateiname",
	MenuCopyDirPath: "Ordnerpfad", MenuIndentSub: "Einrückung", MenuIndent: "Einrückung vergrößern", MenuUnindent: "Einrückung verkleinern",
	MenuConvertCase: "Schreibweise umwandeln", MenuUpperCase: "GROSSBUCHSTABEN", MenuLowerCase: "kleinbuchstaben", MenuProperCase: "Erster Buchstabe Groß",
	MenuProperCaseBlend: "Erster Buchstabe Groß (gemischt)", MenuSentenceCase: "Satzanfang groß", MenuSentenceCaseBlend: "Satzanfang groß (gemischt)",
	MenuInvertCase: "uMKEHREN", MenuRandomCase: "zUfÄlLiG", MenuLineOperations: "Zeilenoperationen",
	MenuDuplicateLine: "Aktuelle Zeile duplizieren", MenuRemoveDupLines: "Doppelte Zeilen entfernen", MenuRemoveConsecutiveDupLines: "Aufeinanderfolgende doppelte Zeilen entfernen",
	MenuSplitLines: "Zeilen umbrechen", MenuJoinLines: "Zeilen verbinden", MenuMoveLineUp: "Zeile nach oben verschieben", MenuMoveLineDown: "Zeile nach unten verschieben",
	MenuDeleteLine: "Aktuelle Zeile löschen", MenuCutLine: "Aktuelle Zeile ausschneiden", MenuTransposeLine: "Mit vorheriger Zeile tauschen",
	MenuRemoveEmptyLines: "Leere Zeilen entfernen", MenuRemoveBlankLines: "Leere Zeilen entfernen (auch mit Leerzeichen)",
	MenuInsertLineAbove: "Leere Zeile darüber einfügen", MenuInsertLineBelow: "Leere Zeile darunter einfügen",
	MenuReverseLines: "Zeilenreihenfolge umkehren", MenuShuffleLines: "Zeilen zufällig anordnen",
	MenuSortAsc: "Zeilen aufsteigend sortieren", MenuSortDesc: "Zeilen absteigend sortieren",
	MenuSortAscCI: "Aufsteigend sortieren (ohne Groß-/Kleinschreibung)", MenuSortDescCI: "Absteigend sortieren (ohne Groß-/Kleinschreibung)",
	MenuSortIntAsc: "Als Ganzzahlen aufsteigend sortieren", MenuSortIntDesc: "Als Ganzzahlen absteigend sortieren",
	MenuSortDecAsc: "Als Dezimalzahlen aufsteigend sortieren", MenuSortDecDesc: "Als Dezimalzahlen absteigend sortieren",
	MenuSortLenAsc: "Nach Länge aufsteigend sortieren", MenuSortLenDesc: "Nach Länge absteigend sortieren",
	MenuComment: "Kommentieren", MenuToggleComment: "Zeilenkommentar umschalten", MenuLineComment: "Zeilen auskommentieren",
	MenuLineUncomment: "Zeilenkommentar entfernen", MenuBlockComment: "Blockkommentar", MenuBlockUncomment: "Blockkommentar entfernen",
	MenuAutoCompletion: "Autovervollständigung", MenuWordCompletion: "Wort vervollständigen", MenuEOLConversion: "Zeilenende umwandeln",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Leerzeichen-Operationen", MenuTrimTrailing: "Leerzeichen am Zeilenende entfernen", MenuTrimLeading: "Leerzeichen am Zeilenanfang entfernen",
	MenuTrimBoth: "Leerzeichen am Anfang und Ende entfernen", MenuEOLToSpace: "Zeilenenden in Leerzeichen", MenuRemoveBlankEOL: "Überflüssige Leerzeichen und Zeilenenden entfernen",
	MenuTabToSpace: "Tabulatoren in Leerzeichen", MenuSpaceToTabAll: "Leerzeichen in Tabulatoren (alle)", MenuSpaceToTabLeading: "Leerzeichen in Tabulatoren (am Anfang)",
	MenuColumnEditor: "Spalteneditor...", MenuReadOnly: "Schreibgeschützt",
	DateShortLayout: "15:04 02.01.2006", DateLongLayout: "15:04:05 02.01.2006 (Monday)",

	MenuFind: "Suchen...", MenuFindInFiles: "In Dateien suchen...", MenuFindNext: "Weitersuchen", MenuFindPrev: "Rückwärts suchen",
	MenuSelectFindNext: "Auswählen und weitersuchen", MenuSelectFindPrev: "Auswählen und rückwärts suchen", MenuReplace: "Ersetzen...",
	MenuIncremental: "Inkrementelle Suche", MenuMark: "Markieren...", MenuSearchResults: "Suchergebnisfenster",
	MenuNextResult: "Nächstes Suchergebnis", MenuPrevResult: "Vorheriges Suchergebnis", MenuGoTo: "Gehe zu...",
	MenuGotoBrace: "Zur passenden Klammer", MenuSelectBrace: "Alles zwischen Klammern auswählen",
	MenuStyleAll: "Alle Vorkommen hervorheben", MenuStyleOne: "Ein Vorkommen hervorheben", MenuClearStyleSub: "Hervorhebung entfernen",
	MenuClearAllStyles: "Alle Hervorhebungen entfernen", MenuJumpUp: "Nach oben springen", MenuJumpDown: "Nach unten springen", MenuBookmark: "Lesezeichen",
	MenuToggleBookmark: "Lesezeichen umschalten", MenuNextBookmark: "Nächstes Lesezeichen", MenuPrevBookmark: "Vorheriges Lesezeichen",
	MenuClearBookmarks: "Alle Lesezeichen entfernen", MenuCutBookmarked: "Zeilen mit Lesezeichen ausschneiden", MenuCopyBookmarked: "Zeilen mit Lesezeichen kopieren",
	MenuRemoveBookmarked: "Zeilen mit Lesezeichen entfernen", MenuRemoveUnbookmarked: "Zeilen ohne Lesezeichen entfernen", MenuInverseBookmarks: "Lesezeichen umkehren",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Mit Stil %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Stil %d entfernen", n) },

	MenuAlwaysOnTop: "Immer im Vordergrund", MenuFullScreen: "Vollbild", MenuShowSymbol: "Symbole anzeigen",
	MenuShowSpaces: "Leerzeichen und Tabulatoren", MenuShowEOL: "Zeilenenden", MenuShowAllChars: "Alle Zeichen",
	MenuShowIndentGuides: "Einrückungslinien", MenuShowWrapSymbol: "Umbruchsymbol", MenuZoom: "Zoom", MenuZoomIn: "Vergrößern",
	MenuZoomOut: "Verkleinern", MenuZoomReset: "Standardgröße", MenuMoveClone: "Dokument verschieben/klonen",
	MenuMoveToOtherView: "In andere Ansicht verschieben", MenuCloneToOtherView: "In andere Ansicht klonen", MenuTab: "Registerkarte",
	MenuNextTab: "Nächste Registerkarte", MenuPrevTab: "Vorherige Registerkarte", MenuMoveTabForward: "Registerkarte nach rechts", MenuMoveTabBackward: "Registerkarte nach links",
	MenuWordWrap: "Zeilenumbruch", MenuLineNumbers: "Zeilennummern", MenuFocusOtherView: "Zur anderen Ansicht wechseln", MenuFoldAll: "Alles einklappen",
	MenuUnfoldAll: "Alles ausklappen", MenuFoldCurrent: "Aktuelle Ebene einklappen", MenuUnfoldCurrent: "Aktuelle Ebene ausklappen",
	MenuFoldLevel: "Ebene einklappen", MenuUnfoldLevel: "Ebene ausklappen", MenuSummary: "Zusammenfassung...", MenuFolderWorkspace: "Ordner als Arbeitsbereich",
	MenuMonitoring: "Überwachung (tail -f)", MenuToolbar: "Werkzeugleiste", MenuStatusBar: "Statusleiste",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Letzte Registerkarte"
		}
		return fmt.Sprintf("Registerkarte %d", n)
	},

	MenuEncANSI: "In ANSI kodieren", MenuEncUTF8: "In UTF-8 kodieren", MenuEncUTF8BOM: "In UTF-8-BOM kodieren", MenuEncUTF16BE: "In UTF-16 BE BOM kodieren",
	MenuEncUTF16LE: "In UTF-16 LE BOM kodieren", MenuCharacterSets: "Zeichensätze", MenuConvANSI: "In ANSI konvertieren", MenuConvUTF8: "In UTF-8 konvertieren",
	MenuConvUTF8BOM: "In UTF-8-BOM konvertieren", MenuConvUTF16BE: "In UTF-16 BE BOM konvertieren", MenuConvUTF16LE: "In UTF-16 LE BOM konvertieren",
	CharsetGroups: map[string]string{"Arabic": "Arabisch", "Baltic": "Baltisch", "Celtic": "Keltisch", "Cyrillic": "Kyrillisch",
		"Central European": "Mitteleuropäisch", "Chinese": "Chinesisch", "Greek": "Griechisch", "Hebrew": "Hebräisch", "Japanese": "Japanisch",
		"Korean": "Koreanisch", "North European": "Nordeuropäisch", "Thai": "Thailändisch", "Turkish": "Türkisch", "Vietnamese": "Vietnamesisch",
		"Western European": "Westeuropäisch"},

	MenuPreferences: "Einstellungen...", MenuColorScheme: "Farbschema", MenuShortcuts: "Tastenkürzel",
	MenuHashGenerate:  func(name string) string { return name + " erzeugen..." },
	MenuHashFiles:     func(name string) string { return name + " von Dateien erzeugen..." },
	MenuHashSelection: func(name string) string { return name + " der Auswahl kopieren" },
	MenuBase64Encode:  "Base64 kodieren", MenuBase64Decode: "Base64 dekodieren", MenuURLEncode: "URL kodieren", MenuURLDecode: "URL dekodieren",
	MenuJSONFormat: "JSON formatieren", MenuJSONMinify: "JSON komprimieren",
	MenuStartRecording: "Aufzeichnung starten/beenden", MenuStopRecording: "Aufzeichnung beenden", MenuPlayback: "Abspielen",
	MenuSaveMacro: "Aufgezeichnetes Makro speichern...", MenuRunMacroMulti: "Makro mehrfach ausführen...", MenuTrimSave: "Leerzeichen am Zeilenende entfernen und speichern",
	MenuRun: "Ausführen...", MenuOpenInBrowser: "Im Browser öffnen", MenuSearchInternet: "Im Internet suchen", MenuOpenSelectedFile: "Datei öffnen (ausgewählter Name)",
	MenuWindows: "Fenster...", MenuSortTabsByName: "Registerkarten nach Name sortieren", MenuSortTabsByPath: "Registerkarten nach Pfad sortieren",
	MenuHelp: "Online-Hilfe", MenuHomePage: "Startseite", MenuAbout: "Über AltNotepad",

	OpenTitle: "Öffnen", SaveTitle: "Speichern", SaveAsTitle: "Speichern unter", SaveCopyTitle: "Kopie speichern unter", AllFiles: "Alle Dateien", TextFiles: "Textdateien",
	Save: "Speichern", DontSave: "Nicht speichern", PathLabel: "Pfad der Datei:",
	SaveChangesAsk: func(path string) string { return "Änderungen an „" + path + "“ speichern?" },
	CreateFileAsk:  func(path string) string { return "„" + path + "“ existiert nicht. Erstellen?" },
	AlreadyOpen:    func(path string) string { return "„" + path + "“ ist in einer anderen Registerkarte geöffnet" },
	FileExists:     func(path string) string { return "„" + path + "“ existiert bereits" },
	FileNotFound:   func(path string) string { return "Datei nicht gefunden: " + path },
	Unencodable: func(enc string) string {
		return "Einige Zeichen können nicht in " + enc + " gespeichert werden. Trotzdem als „?“ speichern?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("%s wird geladen: %d%%", name, percent) },
	LargeFileMode: "Große Datei: Syntaxhervorhebung, Einklappen und Zeilenumbruch sind aus",
	LargeFileMark: "(große Datei)",
	ReloadTitle:   "Neu laden", RenameTitle: "Umbenennen", NewName: "Neuer Name:", DeleteFileTitle: "Vom Datenträger löschen",
	ReloadAsk: func(path string) string {
		return "„" + path + "“ neu laden? Ungespeicherte Änderungen gehen verloren."
	},
	DeleteFileAsk:    func(path string) string { return "„" + path + "“ vom Datenträger löschen?" },
	FileChangedTitle: "Datei geändert", FileDeletedTitle: "Datei gelöscht",
	FileChangedAsk: func(path string) string {
		return "„" + path + "“\n\nDiese Datei wurde von einem anderen Programm geändert.\nNeu laden?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "„" + path + "“\n\nDiese Datei wurde von einem anderen Programm geändert.\nNeu laden und die Änderungen im Editor verwerfen?"
	},
	FileDeletedAsk: func(path string) string {
		return "„" + path + "“\n\nDiese Datei existiert nicht mehr.\nIm Editor behalten?"
	},
	NotASession: func(name string) string { return name + " ist keine Sitzungsdatei" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("Länge: %d   Zeilen: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Z: %d   Sp: %d   Pos: %d", line, col, pos) },
	StatusSel:  "Ausw:",

	GotoTitle: "Gehe zu...", GotoLine: "Zeile", GotoOffset: "Position", Go: "Los",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Sie sind hier: %d   Gehe zu (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Sie sind hier: %d   Gehe zu (0 - %d):", cur, total)
	},
	ColumnTitle: "Spalten- und Mehrfachauswahl-Editor", ColumnText: "Einzufügender Text", ColumnNumbers: "Einzufügende Zahl",
	ColumnInitial: "Startwert:", ColumnIncrease: "Erhöhen um:", ColumnRepeat: "Wiederholen:", ColumnFormat: "Format:",
	ColumnLeading: "Auffüllen:", LeadingNone: "Nichts", LeadingZeros: "Nullen", LeadingSpaces: "Leerzeichen",
	FullPath: "Vollständiger Dateipfad", Modified: "Geändert", SummaryChars: "Zeichen (ohne Zeilenenden)", SummaryWords: "Wörter",
	SummaryLines: "Zeilen", SummaryNonBlankLines: "Nicht leere Zeilen", SummaryLength: "Dokumentlänge", SummarySelected: "Ausgewählte Zeichen",
	Bytes: "Bytes", RunTitle: "Ausführen", RunLabel: "Auszuführendes Programm:", Run: "Ausführen",
	Name: "Name", State: "Status", ModifiedState: "geändert", Activate: "Aktivieren", CloseWindows: "Schließen",
	Command: "Befehl", Shortcut: "Tastenkürzel", HashInput: "Text:", HashEachLine: "Jede Zeile einzeln behandeln",
	CopyToClipboard:   "Kopieren",
	CopiedToClipboard: func(s string) string { return "In die Zwischenablage kopiert: " + s },
	MacroName:         "Name des Makros:", MacroTimes: "Wie oft ausführen (* - bis zum Dateiende):",

	FindTab: "Suchen", ReplaceTab: "Ersetzen", FindInFilesTab: "In Dateien suchen", MarkTab: "Markieren", FindWhat: "Suchen nach:",
	ReplaceWith: "Ersetzen durch:", Filters: "Filter:", Directory: "Ordner:", CurrentFolder: "Aktueller Ordner",
	MatchCase: "Groß-/Kleinschreibung", WholeWord: "Nur ganzes Wort", WrapAround: "Am Ende von vorn beginnen", Backward: "Rückwärts",
	InSelection: "In der Auswahl", InSubfolders: "In allen Unterordnern", InHidden: "In versteckten Ordnern", BookmarkLine: "Lesezeichen setzen",
	PurgeEach: "Bei jeder Suche zurücksetzen", SearchMode: "Suchmodus:", ModeNormal: "Normal", ModeExtended: "Erweitert (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Regulärer Ausdruck", DotAll: ". findet Zeilenumbruch",
	FindNext: "Weitersuchen", FindPrev: "Rückwärts", Count: "Zählen", FindAllCurrent: "Alle im aktuellen Dokument",
	FindAllOpen: "Alle in allen offenen Dokumenten", Replace: "Ersetzen", ReplaceAll: "Alle ersetzen",
	ReplaceAllOpen: "Alle in allen offenen Dokumenten ersetzen", FindAll: "Alle suchen", ReplaceInFiles: "In Dateien ersetzen",
	MarkAll: "Alle markieren", ClearMarks: "Markierungen entfernen", CopyMarked: "Markierten Text kopieren", Wrapped: "Ende erreicht, von vorn fortgesetzt",
	ReadOnlyDoc: "Das Dokument ist schreibgeschützt",
	NotFound:    func(text string) string { return "„" + text + "“ nicht gefunden" },
	CountResult: func(n int) string { return fmt.Sprintf("Anzahl: %d %s", n, deP(n, "Treffer", "Treffer")) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d %s ersetzt", n, deP(n, "Vorkommen", "Vorkommen"))
	},
	MarkedCount: func(n int) string { return fmt.Sprintf("%d %s markiert", n, deP(n, "Treffer", "Treffer")) },
	FoundCount:  func(n int) string { return fmt.Sprintf("%d %s gefunden", n, deP(n, "Treffer", "Treffer")) },
	HitsCount:   func(n int) string { return fmt.Sprintf("(%d %s)", n, deP(n, "Treffer", "Treffer")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d Treffer in %d %s", n, files, deP(files, "Datei", "Dateien"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d Vorkommen in %d %s ersetzt", n, files, deP(files, "Datei", "Dateien"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Alle „" + find + "“ durch „" + repl + "“ in allen Dateien ersetzen in\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "Ordner nicht gefunden: " + path },
	Searching:     func(path string) string { return "Suche: " + path },
	SearchResults: "Suchergebnisse", Stop: "Stopp", Stopped: "(gestoppt)", Clear: "Leeren", Line: "Zeile",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Suche „%s“ (%d Treffer in %d %s von %d durchsuchten)", text, hits, files, deP(files, "Datei", "Dateien"), searched)
	},

	Workspace: "Arbeitsbereich", Refresh: "Aktualisieren", CollapseAll: "Alle einklappen",

	PageGeneral: "Allgemein", PageEditing: "Bearbeiten", PageDisplay: "Anzeige", PageFiles: "Neues Dokument und Dateien",
	Language: "Sprache:", LanguageSystem: "Wie im System", Theme: "Design:", ThemeDark: "Dunkel", ThemeLight: "Hell",
	ColorScheme: "Farbschema:", SchemeByTheme: "Nach Design",
	ShowToolbar: "Werkzeugleiste anzeigen", ShowStatusBar: "Statusleiste anzeigen", TabCloseButtons: "Schließen-Schaltflächen auf den Registerkarten",
	DoubleClickCloses: "Doppelklick schließt Registerkarte", AlwaysOnTop: "Immer im Vordergrund", SingleInstance: "Dateien im laufenden Fenster öffnen",
	RememberSession: "Geöffnete Dateien für die nächste Sitzung merken", BackupSession: "Ungespeicherte Änderungen zwischen Sitzungen behalten",
	BackupEvery: "Sicherung alle, Sekunden:", MaxRecent: "Zuletzt geöffnete Dateien:",
	TabSize: "Tabulatorbreite:", InsertSpaces: "Tabulatoren durch Leerzeichen ersetzen", AutoIndent: "Automatische Einrückung", AutoClose: "Klammern und Anführungszeichen schließen",
	AutoCompletion: "Wortvervollständigung beim Tippen", SmartHome: "Pos1 springt zum ersten Zeichen", MultiEdit: "Mehrfachbearbeitung (Strg+Klick)",
	ScrollPast: "Über die letzte Zeile hinaus scrollen", CopyLineNoSel: "Zeile kopieren/ausschneiden ohne Auswahl", CaretWidth: "Cursorbreite:",
	CaretBlink: "Cursor-Blinkrate, ms (0 - kein Blinken):", EdgeColumn: "Markierung langer Zeilen, Spalte (0 - keine):", WordChars: "Zusätzliche Wortzeichen:",
	FontSize: "Schriftgröße:", FontFile: "Schriftdatei:", BuiltInFont: "JetBrains Mono (integriert)", Fonts: "Schriften", LineNumbers: "Zeilennummern",
	BookmarkMargin: "Lesezeichenrand", FoldMargin: "Einklapprand", CurrentLine: "Aktuelle Zeile hervorheben", SmartHighlight: "Intelligente Hervorhebung",
	SmartMatchCase: "Intelligente Hervorhebung: Groß-/Kleinschreibung", SmartWholeWord: "Intelligente Hervorhebung: nur ganze Wörter", BraceMatch: "Passende Klammern hervorheben",
	WrapSymbol: "Umbruchsymbol anzeigen", ChangeHistory: "Änderungsverlauf am Rand",
	NewDocEOL: "Zeilenenden neuer Dokumente:", NewDocEncoding: "Kodierung neuer Dokumente:", ANSICharset: "ANSI-Zeichensatz:",
	NewDocLanguage: "Sprache neuer Dokumente:", SystemDefault: "Wie im System", ByLanguage: "Nach Oberflächensprache",
	LargeFileMB: "Grenze für große Dateien, MB:", CheckFileChanges: "Von anderen Programmen geänderte Dateien erkennen", AutoReload: "Unveränderte ohne Nachfrage neu laden",

	AboutTitle:       func(name string) string { return "Über " + name },
	AboutDescription: "Ein Editor für Text und Quellcode", Version: "Version", Author: "Autor:", License: "Lizenz:",
	VisitWebsite: "Webseite besuchen", Close: "Schließen",
}
