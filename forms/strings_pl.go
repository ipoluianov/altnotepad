package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func plP(n int, one, few, many string) string { return i18n.Plural("pl", n, one, few, many) }

var plStrings = Strings{
	Error:      "Błąd",
	NewDocName: "nowy",
	NormalText: "Zwykły tekst",

	MenuFile: "Plik", MenuEdit: "Edycja", MenuSearch: "Szukaj", MenuView: "Widok", MenuEncoding: "Kodowanie", MenuLanguage: "Składnia",
	MenuSettings: "Ustawienia", MenuTools: "Narzędzia", MenuMacro: "Makro", MenuRunMenu: "Uruchom", MenuWindow: "Okno",

	MenuNew: "Nowy", MenuOpen: "Otwórz...", MenuOpenContainingFolderSub: "Otwórz folder dokumentu",
	MenuOpenContainingFolder: "Menedżer plików", MenuOpenDefaultViewer: "Otwórz w domyślnym programie", MenuOpenFolderWorkspace: "Otwórz folder jako obszar roboczy...",
	MenuReload: "Wczytaj ponownie z dysku", MenuSave: "Zapisz", MenuSaveAs: "Zapisz jako...", MenuSaveCopy: "Zapisz kopię jako...", MenuSaveAll: "Zapisz wszystko",
	MenuRename: "Zmień nazwę...", MenuClose: "Zamknij", MenuCloseAll: "Zamknij wszystko", MenuCloseMore: "Zamknij wiele dokumentów",
	MenuCloseAllButActive: "Zamknij wszystko oprócz bieżącego", MenuCloseLeft: "Zamknij wszystko po lewej", MenuCloseRight: "Zamknij wszystko po prawej",
	MenuCloseUnchanged: "Zamknij niezmienione", MenuDeleteFile: "Usuń z dysku...", MenuLoadSession: "Wczytaj sesję...",
	MenuSaveSession: "Zapisz sesję...", MenuRecentFiles: "Ostatnie pliki", MenuClearRecent: "Wyczyść listę ostatnich plików",
	MenuRestoreClosed: "Przywróć ostatnio zamknięty plik", MenuExit: "Zakończ",
	MenuPrint: "Drukuj...", MenuExportHTML: "Eksportuj do HTML...", MenuFunctionList: "Lista funkcji", Filter: "Filtr", MenuDocumentMap: "Mapa dokumentu",
	MenuChangeHistory: "Historia zmian", MenuNextChange: "Przejdź do następnej zmiany", MenuPrevChange: "Przejdź do poprzedniej zmiany", MenuClearChanges: "Wyczyść historię zmian",
	SavedTo:         func(path string) string { return "Zapisano do " + path },
	MenuMultiSelect: "Zaznaczenie wielokrotne", MenuMultiAll: "Zaznacz wszystkie wystąpienia", MenuMultiAllCase: "Zaznacz wszystkie (uwzględnij wielkość liter)",
	MenuMultiNext: "Dodaj następne wystąpienie", MenuMultiUndo: "Cofnij ostatnio dodane", MenuMultiSkip: "Pomiń bieżące i dodaj następne",

	MenuUndo: "Cofnij", MenuRedo: "Ponów", MenuCut: "Wytnij", MenuCopy: "Kopiuj", MenuPaste: "Wklej", MenuDelete: "Usuń", MenuSelectAll: "Zaznacz wszystko",
	MenuInsert: "Wstaw", MenuDateShort: "Data i godzina (krótka)", MenuDateLong: "Data i godzina (długa)", MenuDateISO: "Data i godzina (ISO 8601)",
	MenuCopyToClipboard: "Kopiuj do schowka", MenuCopyFullPath: "Pełna ścieżka pliku", MenuCopyFileName: "Nazwa pliku",
	MenuCopyDirPath: "Ścieżka folderu", MenuIndentSub: "Wcięcie", MenuIndent: "Zwiększ wcięcie", MenuUnindent: "Zmniejsz wcięcie",
	MenuConvertCase: "Zamień na", MenuUpperCase: "WIELKIE LITERY", MenuLowerCase: "małe litery", MenuProperCase: "Jak W Tytule",
	MenuProperCaseBlend: "Jak W Tytule (mieszane)", MenuSentenceCase: "Jak w zdaniu", MenuSentenceCaseBlend: "Jak w zdaniu (mieszane)",
	MenuInvertCase: "oDWRÓĆ wIELKOŚĆ", MenuRandomCase: "lOsOwO", MenuLineOperations: "Operacje na wierszach",
	MenuDuplicateLine: "Powiel bieżący wiersz", MenuRemoveDupLines: "Usuń zduplikowane wiersze", MenuRemoveConsecutiveDupLines: "Usuń kolejne zduplikowane wiersze",
	MenuSplitLines: "Podziel wiersze", MenuJoinLines: "Połącz wiersze", MenuMoveLineUp: "Przesuń wiersz w górę", MenuMoveLineDown: "Przesuń wiersz w dół",
	MenuDeleteLine: "Usuń bieżący wiersz", MenuCutLine: "Wytnij bieżący wiersz", MenuTransposeLine: "Zamień z poprzednim wierszem",
	MenuRemoveEmptyLines: "Usuń puste wiersze", MenuRemoveBlankLines: "Usuń puste wiersze (także ze spacjami)",
	MenuInsertLineAbove: "Wstaw pusty wiersz powyżej", MenuInsertLineBelow: "Wstaw pusty wiersz poniżej",
	MenuReverseLines: "Odwróć kolejność wierszy", MenuShuffleLines: "Losowa kolejność wierszy",
	MenuSortAsc: "Sortuj wiersze rosnąco", MenuSortDesc: "Sortuj wiersze malejąco",
	MenuSortAscCI: "Sortuj rosnąco bez rozróżniania wielkości liter", MenuSortDescCI: "Sortuj malejąco bez rozróżniania wielkości liter",
	MenuSortIntAsc: "Sortuj jako liczby całkowite rosnąco", MenuSortIntDesc: "Sortuj jako liczby całkowite malejąco",
	MenuSortDecAsc: "Sortuj jako liczby dziesiętne rosnąco", MenuSortDecDesc: "Sortuj jako liczby dziesiętne malejąco",
	MenuSortLenAsc: "Sortuj według długości rosnąco", MenuSortLenDesc: "Sortuj według długości malejąco",
	MenuComment: "Komentarze", MenuToggleComment: "Przełącz komentarz wiersza", MenuLineComment: "Zakomentuj wiersze",
	MenuLineUncomment: "Odkomentuj wiersze", MenuBlockComment: "Komentarz blokowy", MenuBlockUncomment: "Usuń komentarz blokowy",
	MenuAutoCompletion: "Autouzupełnianie", MenuWordCompletion: "Uzupełnij słowo", MenuEOLConversion: "Konwersja końca wiersza",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Operacje na odstępach", MenuTrimTrailing: "Usuń końcowe odstępy", MenuTrimLeading: "Usuń początkowe odstępy",
	MenuTrimBoth: "Usuń odstępy na początku i końcu", MenuEOLToSpace: "Koniec wiersza na spację", MenuRemoveBlankEOL: "Usuń zbędne odstępy i końce wierszy",
	MenuTabToSpace: "Tabulatory na spacje", MenuSpaceToTabAll: "Spacje na tabulatory (wszystkie)", MenuSpaceToTabLeading: "Spacje na tabulatory (początkowe)",
	MenuColumnEditor: "Edytor kolumn...", MenuReadOnly: "Tylko do odczytu",
	DateShortLayout: "15:04 02.01.2006", DateLongLayout: "15:04:05 02.01.2006 (Monday)",

	MenuFind: "Znajdź...", MenuFindInFiles: "Znajdź w plikach...", MenuFindNext: "Znajdź następny", MenuFindPrev: "Znajdź poprzedni",
	MenuSelectFindNext: "Zaznacz i znajdź następny", MenuSelectFindPrev: "Zaznacz i znajdź poprzedni", MenuReplace: "Zamień...",
	MenuIncremental: "Wyszukiwanie przyrostowe", MenuMark: "Oznacz...", MenuSearchResults: "Okno wyników wyszukiwania",
	MenuNextResult: "Następny wynik", MenuPrevResult: "Poprzedni wynik", MenuGoTo: "Przejdź do...",
	MenuGotoBrace: "Przejdź do pasującego nawiasu", MenuSelectBrace: "Zaznacz wszystko między nawiasami",
	MenuStyleAll: "Wyróżnij wszystkie wystąpienia", MenuStyleOne: "Wyróżnij jedno wystąpienie", MenuClearStyleSub: "Usuń wyróżnienie",
	MenuClearAllStyles: "Usuń wszystkie wyróżnienia", MenuJumpUp: "Skocz w górę", MenuJumpDown: "Skocz w dół", MenuBookmark: "Zakładki",
	MenuToggleBookmark: "Przełącz zakładkę", MenuNextBookmark: "Następna zakładka", MenuPrevBookmark: "Poprzednia zakładka",
	MenuClearBookmarks: "Usuń wszystkie zakładki", MenuCutBookmarked: "Wytnij wiersze z zakładkami", MenuCopyBookmarked: "Kopiuj wiersze z zakładkami",
	MenuRemoveBookmarked: "Usuń wiersze z zakładkami", MenuRemoveUnbookmarked: "Usuń wiersze bez zakładek", MenuInverseBookmarks: "Odwróć zakładki",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Styl %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Usuń styl %d", n) },

	MenuAlwaysOnTop: "Zawsze na wierzchu", MenuFullScreen: "Pełny ekran", MenuShowSymbol: "Pokaż symbole",
	MenuShowSpaces: "Spacje i tabulatory", MenuShowEOL: "Końce wierszy", MenuShowAllChars: "Wszystkie znaki",
	MenuShowIndentGuides: "Prowadnice wcięć", MenuShowWrapSymbol: "Symbol zawijania", MenuZoom: "Powiększenie", MenuZoomIn: "Powiększ",
	MenuZoomOut: "Pomniejsz", MenuZoomReset: "Domyślne powiększenie", MenuMoveClone: "Przenieś/sklonuj dokument",
	MenuMoveToOtherView: "Przenieś do drugiego widoku", MenuCloneToOtherView: "Sklonuj do drugiego widoku", MenuTab: "Karty",
	MenuNextTab: "Następna karta", MenuPrevTab: "Poprzednia karta", MenuMoveTabForward: "Przesuń kartę w prawo", MenuMoveTabBackward: "Przesuń kartę w lewo",
	MenuWordWrap: "Zawijanie wierszy", MenuLineNumbers: "Numery wierszy", MenuFocusOtherView: "Przejdź do drugiego widoku", MenuFoldAll: "Zwiń wszystko",
	MenuUnfoldAll: "Rozwiń wszystko", MenuFoldCurrent: "Zwiń bieżący poziom", MenuUnfoldCurrent: "Rozwiń bieżący poziom",
	MenuFoldLevel: "Zwiń poziom", MenuUnfoldLevel: "Rozwiń poziom", MenuSummary: "Podsumowanie...", MenuFolderWorkspace: "Folder jako obszar roboczy",
	MenuMonitoring: "Monitorowanie (tail -f)", MenuToolbar: "Pasek narzędzi", MenuStatusBar: "Pasek stanu",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Ostatnia karta"
		}
		return fmt.Sprintf("Karta %d", n)
	},

	MenuEncANSI: "Koduj w ANSI", MenuEncUTF8: "Koduj w UTF-8", MenuEncUTF8BOM: "Koduj w UTF-8-BOM", MenuEncUTF16BE: "Koduj w UTF-16 BE BOM",
	MenuEncUTF16LE: "Koduj w UTF-16 LE BOM", MenuCharacterSets: "Zestawy znaków", MenuConvANSI: "Konwertuj na ANSI", MenuConvUTF8: "Konwertuj na UTF-8",
	MenuConvUTF8BOM: "Konwertuj na UTF-8-BOM", MenuConvUTF16BE: "Konwertuj na UTF-16 BE BOM", MenuConvUTF16LE: "Konwertuj na UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Arabskie", "Baltic": "Bałtyckie", "Celtic": "Celtyckie", "Cyrillic": "Cyrylica",
		"Central European": "Środkowoeuropejskie", "Chinese": "Chińskie", "Greek": "Greckie", "Hebrew": "Hebrajskie", "Japanese": "Japońskie",
		"Korean": "Koreańskie", "North European": "Północnoeuropejskie", "Thai": "Tajskie", "Turkish": "Tureckie", "Vietnamese": "Wietnamskie",
		"Western European": "Zachodnioeuropejskie"},

	MenuPreferences: "Preferencje...", MenuColorScheme: "Schemat kolorów", MenuShortcuts: "Skróty klawiszowe",
	MenuHashGenerate:  func(name string) string { return "Generuj " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Generuj " + name + " plików..." },
	MenuHashSelection: func(name string) string { return "Kopiuj " + name + " zaznaczenia" },
	MenuBase64Encode:  "Koduj Base64", MenuBase64Decode: "Dekoduj Base64", MenuURLEncode: "Koduj URL", MenuURLDecode: "Dekoduj URL",
	MenuJSONFormat: "Formatuj JSON", MenuJSONMinify: "Kompaktuj JSON",
	MenuStartRecording: "Rozpocznij/zatrzymaj nagrywanie", MenuStopRecording: "Zatrzymaj nagrywanie", MenuPlayback: "Odtwórz",
	MenuSaveMacro: "Zapisz nagrane makro...", MenuRunMacroMulti: "Uruchom makro wielokrotnie...", MenuTrimSave: "Usuń końcowe odstępy i zapisz",
	MenuRun: "Uruchom...", MenuOpenInBrowser: "Otwórz w przeglądarce", MenuSearchInternet: "Szukaj w Internecie", MenuOpenSelectedFile: "Otwórz plik (zaznaczona nazwa)",
	MenuWindows: "Okna...", MenuSortTabsByName: "Sortuj karty według nazwy", MenuSortTabsByPath: "Sortuj karty według ścieżki",
	MenuHelp: "Pomoc online", MenuHomePage: "Strona domowa", MenuAbout: "O programie AltNotepad",

	OpenTitle: "Otwórz", SaveTitle: "Zapisz", SaveAsTitle: "Zapisz jako", SaveCopyTitle: "Zapisz kopię jako", AllFiles: "Wszystkie pliki",
	TextFiles: "Pliki tekstowe", Save: "Zapisz", DontSave: "Nie zapisuj", PathLabel: "Ścieżka pliku:",
	SaveChangesAsk: func(path string) string { return "Zapisać zmiany w „" + path + "”?" },
	CreateFileAsk:  func(path string) string { return "„" + path + "” nie istnieje. Utworzyć?" },
	AlreadyOpen:    func(path string) string { return "„" + path + "” jest otwarty w innej karcie" },
	FileExists:     func(path string) string { return "„" + path + "” już istnieje" },
	FileNotFound:   func(path string) string { return "Nie znaleziono pliku: " + path },
	Unencodable: func(enc string) string {
		return "Niektórych znaków nie można zapisać w " + enc + ". Zapisać je jako „?”?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Wczytywanie %s: %d%%", name, percent) },
	LargeFileMode: "Duży plik: podświetlanie składni, zwijanie i zawijanie wierszy wyłączone",
	LargeFileMark: "(duży plik)",
	ReloadTitle:   "Wczytaj ponownie", RenameTitle: "Zmień nazwę", NewName: "Nowa nazwa:", DeleteFileTitle: "Usuń z dysku",
	ReloadAsk: func(path string) string {
		return "Wczytać ponownie „" + path + "”? Niezapisane zmiany zostaną utracone."
	},
	DeleteFileAsk:    func(path string) string { return "Usunąć „" + path + "” z dysku?" },
	FileChangedTitle: "Plik zmieniony", FileDeletedTitle: "Plik usunięty",
	FileChangedAsk: func(path string) string {
		return "„" + path + "”\n\nTen plik został zmieniony przez inny program.\nWczytać go ponownie?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "„" + path + "”\n\nTen plik został zmieniony przez inny program.\nWczytać go ponownie i utracić zmiany z edytora?"
	},
	FileDeletedAsk: func(path string) string {
		return "„" + path + "”\n\nTen plik już nie istnieje.\nZachować go w edytorze?"
	},
	NotASession: func(name string) string { return name + " nie jest plikiem sesji" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("długość: %d   wiersze: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Wrs: %d   Kol: %d   Poz: %d", line, col, pos) },
	StatusSel:  "Zazn:",

	GotoTitle: "Przejdź do...", GotoLine: "Wiersz", GotoOffset: "Pozycja", Go: "Przejdź",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Jesteś tutaj: %d   Przejdź do (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Jesteś tutaj: %d   Przejdź do (0 - %d):", cur, total)
	},
	ColumnTitle: "Edytor kolumn i zaznaczenie wielokrotne", ColumnText: "Tekst do wstawienia", ColumnNumbers: "Liczba do wstawienia",
	ColumnInitial: "Liczba początkowa:", ColumnIncrease: "Zwiększaj o:", ColumnRepeat: "Powtórz:", ColumnFormat: "Format:",
	ColumnLeading: "Wypełnienie:", LeadingNone: "Brak", LeadingZeros: "Zera", LeadingSpaces: "Spacje",
	FullPath: "Pełna ścieżka", Modified: "Zmodyfikowano", SummaryChars: "Znaki (bez końców wierszy)", SummaryWords: "Słowa",
	SummaryLines: "Wiersze", SummaryNonBlankLines: "Niepuste wiersze", SummaryLength: "Długość dokumentu", SummarySelected: "Zaznaczone znaki",
	Bytes: "bajtów", RunTitle: "Uruchom", RunLabel: "Program do uruchomienia:", Run: "Uruchom",
	Name: "Nazwa", State: "Stan", ModifiedState: "zmieniony", Activate: "Aktywuj", CloseWindows: "Zamknij",
	Command: "Polecenie", Shortcut: "Skrót", HashInput: "Tekst:", HashEachLine: "Traktuj każdy wiersz osobno",
	CopyToClipboard:   "Kopiuj",
	CopiedToClipboard: func(s string) string { return "Skopiowano do schowka: " + s },
	MacroName:         "Nazwa makra:", MacroTimes: "Ile razy uruchomić (* - do końca pliku):",

	FindTab: "Znajdź", ReplaceTab: "Zamień", FindInFilesTab: "Znajdź w plikach", MarkTab: "Oznacz", FindWhat: "Znajdź:",
	ReplaceWith: "Zamień na:", Filters: "Filtry:", Directory: "Folder:", CurrentFolder: "Bieżący folder",
	MatchCase: "Uwzględnij wielkość liter", WholeWord: "Tylko całe słowa", WrapAround: "Zawijaj wyszukiwanie", Backward: "Wstecz",
	InSelection: "W zaznaczeniu", InSubfolders: "We wszystkich podfolderach", InHidden: "W ukrytych folderach", BookmarkLine: "Dodaj zakładkę",
	PurgeEach: "Czyść przy każdym wyszukiwaniu", SearchMode: "Tryb wyszukiwania:", ModeNormal: "Zwykły", ModeExtended: "Rozszerzony (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Wyrażenie regularne", DotAll: ". obejmuje nowe wiersze",
	FindNext: "Znajdź następny", FindPrev: "Znajdź poprzedni", Count: "Policz", FindAllCurrent: "Znajdź wszystko w bieżącym dokumencie",
	FindAllOpen: "Znajdź wszystko w otwartych dokumentach", Replace: "Zamień", ReplaceAll: "Zamień wszystko",
	ReplaceAllOpen: "Zamień wszystko w otwartych dokumentach", FindAll: "Znajdź wszystko", ReplaceInFiles: "Zamień w plikach",
	MarkAll: "Oznacz wszystko", ClearMarks: "Usuń oznaczenia", CopyMarked: "Kopiuj oznaczony tekst", Wrapped: "Osiągnięto koniec, kontynuowano od początku",
	ReadOnlyDoc: "Dokument jest tylko do odczytu",
	NotFound:    func(text string) string { return "Nie znaleziono „" + text + "”" },
	CountResult: func(n int) string {
		return fmt.Sprintf("Razem: %d %s", n, plP(n, "dopasowanie", "dopasowania", "dopasowań"))
	},
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("Zamieniono: %d %s", n, plP(n, "wystąpienie", "wystąpienia", "wystąpień"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("Oznaczono: %d %s", n, plP(n, "dopasowanie", "dopasowania", "dopasowań"))
	},
	FoundCount: func(n int) string { return fmt.Sprintf("%d %s", n, plP(n, "wynik", "wyniki", "wyników")) },
	HitsCount:  func(n int) string { return fmt.Sprintf("(%d %s)", n, plP(n, "wynik", "wyniki", "wyników")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s w %d %s", n, plP(n, "wynik", "wyniki", "wyników"), files, plP(files, "pliku", "plikach", "plikach"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s w %d %s", n, plP(n, "zamiana", "zamiany", "zamian"), files, plP(files, "pliku", "plikach", "plikach"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Zamienić wszystkie „" + find + "” na „" + repl + "” we wszystkich plikach w\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "Folder nie istnieje: " + path },
	Searching:     func(path string) string { return "Wyszukiwanie: " + path },
	SearchResults: "Wyniki wyszukiwania", Stop: "Zatrzymaj", Stopped: "(zatrzymano)", Clear: "Wyczyść", Line: "Wiersz",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Szukaj „%s” (%d %s w %d z %d przeszukanych plików)", text, hits, plP(hits, "wynik", "wyniki", "wyników"), files, searched)
	},

	Workspace: "Obszar roboczy", Refresh: "Odśwież", CollapseAll: "Zwiń wszystko",

	PageGeneral: "Ogólne", PageEditing: "Edycja", PageDisplay: "Wyświetlanie", PageFiles: "Nowy dokument i pliki",
	Language: "Język:", LanguageSystem: "Jak w systemie", Theme: "Motyw:", ThemeDark: "Ciemny", ThemeLight: "Jasny",
	ColorScheme: "Schemat kolorów:", SchemeByTheme: "Według motywu",
	ShowToolbar: "Pokaż pasek narzędzi", ShowStatusBar: "Pokaż pasek stanu", TabCloseButtons: "Przyciski zamykania na kartach",
	DoubleClickCloses: "Dwuklik zamyka kartę", AlwaysOnTop: "Zawsze na wierzchu", SingleInstance: "Otwieraj pliki w już otwartym oknie",
	RememberSession: "Zapamiętaj otwarte pliki na następną sesję", BackupSession: "Zachowuj niezapisane zmiany między sesjami",
	BackupEvery: "Kopia zapasowa co, sekund:", MaxRecent: "Ostatnich plików na liście:",
	TabSize: "Rozmiar tabulatora:", InsertSpaces: "Zamieniaj tabulatory na spacje", AutoIndent: "Automatyczne wcięcie", AutoClose: "Zamykaj nawiasy i cudzysłowy",
	AutoCompletion: "Uzupełniaj słowa podczas pisania", SmartHome: "Home przechodzi do pierwszego niebiałego znaku", MultiEdit: "Edycja wielokrotna (Ctrl+klik)",
	ScrollPast: "Przewijaj za ostatni wiersz", CopyLineNoSel: "Kopiuj/wycinaj wiersz bez zaznaczenia", CaretWidth: "Szerokość kursora:",
	CaretBlink: "Miganie kursora, ms (0 - bez migania):", EdgeColumn: "Znacznik długich wierszy, kolumna (0 - brak):", WordChars: "Dodatkowe znaki słów:",
	FontSize: "Rozmiar czcionki:", FontFile: "Plik czcionki:", BuiltInFont: "JetBrains Mono (wbudowana)", Fonts: "Czcionki", LineNumbers: "Numery wierszy",
	BookmarkMargin: "Margines zakładek", FoldMargin: "Margines zwijania", CurrentLine: "Wyróżniaj bieżący wiersz", SmartHighlight: "Inteligentne wyróżnianie",
	SmartMatchCase: "Inteligentne wyróżnianie: wielkość liter", SmartWholeWord: "Inteligentne wyróżnianie: tylko całe słowa", BraceMatch: "Wyróżniaj pasujące nawiasy",
	WrapSymbol: "Pokaż symbol zawijania", ChangeHistory: "Historia zmian na marginesie",
	NewDocEOL: "Koniec wiersza nowych dokumentów:", NewDocEncoding: "Kodowanie nowych dokumentów:", ANSICharset: "Zestaw znaków ANSI:",
	NewDocLanguage: "Składnia nowych dokumentów:", SystemDefault: "Jak w systemie", ByLanguage: "Według języka interfejsu",
	LargeFileMB: "Próg dużego pliku, MB:", CheckFileChanges: "Wykrywaj pliki zmienione przez inne programy", AutoReload: "Wczytuj je ponownie bez pytania, jeśli niezmienione",

	AboutTitle:       func(name string) string { return "O programie " + name },
	AboutDescription: "Edytor tekstu i kodu źródłowego", Version: "Wersja", Author: "Autor:", License: "Licencja:",
	VisitWebsite: "Odwiedź stronę", Close: "Zamknij",
}
