package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func itP(n int, one, many string) string { return i18n.Plural("it", n, one, many) }

var it = Strings{
	Error:      "Errore",
	NewDocName: "nuovo",
	NormalText: "Testo normale",

	MenuFile: "File", MenuEdit: "Modifica", MenuSearch: "Cerca", MenuView: "Visualizza", MenuEncoding: "Codifica", MenuLanguage: "Linguaggio",
	MenuSettings: "Impostazioni", MenuTools: "Strumenti", MenuMacro: "Macro", MenuRunMenu: "Esegui", MenuWindow: "Finestra",

	MenuNew: "Nuovo", MenuOpen: "Apri...", MenuOpenContainingFolderSub: "Apri la cartella del documento",
	MenuOpenContainingFolder: "Gestore file", MenuOpenDefaultViewer: "Apri con il programma predefinito", MenuOpenFolderWorkspace: "Apri cartella come area di lavoro...",
	MenuReload: "Ricarica dal disco", MenuSave: "Salva", MenuSaveAs: "Salva con nome...", MenuSaveCopy: "Salva una copia con nome...", MenuSaveAll: "Salva tutto",
	MenuRename: "Rinomina...", MenuClose: "Chiudi", MenuCloseAll: "Chiudi tutto", MenuCloseMore: "Chiudi più documenti",
	MenuCloseAllButActive: "Chiudi tutto tranne il documento attivo", MenuCloseLeft: "Chiudi tutto a sinistra", MenuCloseRight: "Chiudi tutto a destra",
	MenuCloseUnchanged: "Chiudi i documenti non modificati", MenuDeleteFile: "Elimina dal disco...", MenuLoadSession: "Carica sessione...",
	MenuSaveSession: "Salva sessione...", MenuRecentFiles: "File recenti", MenuClearRecent: "Svuota l'elenco dei file recenti",
	MenuRestoreClosed: "Riapri l'ultimo file chiuso", MenuExit: "Esci",
	MenuPrint: "Stampa...", MenuExportHTML: "Esporta in HTML...", MenuFunctionList: "Elenco funzioni", Filter: "Filtro", MenuDocumentMap: "Mappa del documento",
	MenuChangeHistory: "Cronologia modifiche", MenuNextChange: "Vai alla modifica successiva", MenuPrevChange: "Vai alla modifica precedente", MenuClearChanges: "Cancella la cronologia modifiche",
	SavedTo:         func(path string) string { return "Salvato in " + path },
	MenuMultiSelect: "Selezione multipla", MenuMultiAll: "Seleziona tutte le occorrenze", MenuMultiAllCase: "Seleziona tutte (maiuscole/minuscole)",
	MenuMultiNext: "Aggiungi l'occorrenza successiva", MenuMultiUndo: "Annulla l'ultima aggiunta", MenuMultiSkip: "Salta e aggiungi la successiva",

	MenuUndo: "Annulla", MenuRedo: "Ripeti", MenuCut: "Taglia", MenuCopy: "Copia", MenuPaste: "Incolla", MenuDelete: "Elimina", MenuSelectAll: "Seleziona tutto",
	MenuInsert: "Inserisci", MenuDateShort: "Data e ora (breve)", MenuDateLong: "Data e ora (estesa)", MenuDateISO: "Data e ora (ISO 8601)",
	MenuCopyToClipboard: "Copia negli appunti", MenuCopyFullPath: "Percorso completo del file", MenuCopyFileName: "Nome del file",
	MenuCopyDirPath: "Percorso della cartella", MenuIndentSub: "Rientro", MenuIndent: "Aumenta rientro", MenuUnindent: "Riduci rientro",
	MenuConvertCase: "Converti in", MenuUpperCase: "MAIUSCOLO", MenuLowerCase: "minuscolo", MenuProperCase: "Iniziali Maiuscole",
	MenuProperCaseBlend: "Iniziali Maiuscole (misto)", MenuSentenceCase: "Maiuscola a inizio frase", MenuSentenceCaseBlend: "Maiuscola a inizio frase (misto)",
	MenuInvertCase: "iNVERTI mAIUSCOLE", MenuRandomCase: "cAsUaLe", MenuLineOperations: "Operazioni sulle righe",
	MenuDuplicateLine: "Duplica la riga", MenuRemoveDupLines: "Rimuovi righe duplicate", MenuRemoveConsecutiveDupLines: "Rimuovi duplicati consecutivi",
	MenuSplitLines: "Dividi le righe", MenuJoinLines: "Unisci le righe", MenuMoveLineUp: "Sposta la riga in alto", MenuMoveLineDown: "Sposta la riga in basso",
	MenuDeleteLine: "Elimina la riga", MenuCutLine: "Taglia la riga", MenuTransposeLine: "Scambia con la riga precedente",
	MenuRemoveEmptyLines: "Rimuovi righe vuote", MenuRemoveBlankLines: "Rimuovi righe vuote (anche con spazi)",
	MenuInsertLineAbove: "Inserisci riga vuota sopra", MenuInsertLineBelow: "Inserisci riga vuota sotto",
	MenuReverseLines: "Inverti l'ordine delle righe", MenuShuffleLines: "Mescola le righe",
	MenuSortAsc: "Ordina le righe in modo crescente", MenuSortDesc: "Ordina le righe in modo decrescente",
	MenuSortAscCI: "Ordine crescente ignorando maiuscole", MenuSortDescCI: "Ordine decrescente ignorando maiuscole",
	MenuSortIntAsc: "Ordina come interi (crescente)", MenuSortIntDesc: "Ordina come interi (decrescente)",
	MenuSortDecAsc: "Ordina come decimali (crescente)", MenuSortDecDesc: "Ordina come decimali (decrescente)",
	MenuSortLenAsc: "Ordina per lunghezza (crescente)", MenuSortLenDesc: "Ordina per lunghezza (decrescente)",
	MenuComment: "Commenti", MenuToggleComment: "Commenta/decommenta le righe", MenuLineComment: "Commenta le righe",
	MenuLineUncomment: "Decommenta le righe", MenuBlockComment: "Commento a blocco", MenuBlockUncomment: "Rimuovi commento a blocco",
	MenuAutoCompletion: "Completamento automatico", MenuWordCompletion: "Completa la parola", MenuEOLConversion: "Conversione fine riga",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Operazioni sugli spazi", MenuTrimTrailing: "Rimuovi spazi a fine riga", MenuTrimLeading: "Rimuovi spazi a inizio riga",
	MenuTrimBoth: "Rimuovi spazi a inizio e fine riga", MenuEOLToSpace: "Fine riga in spazi", MenuRemoveBlankEOL: "Rimuovi spazi e fine riga superflui",
	MenuTabToSpace: "Tabulazioni in spazi", MenuSpaceToTabAll: "Spazi in tabulazioni (tutti)", MenuSpaceToTabLeading: "Spazi in tabulazioni (a inizio riga)",
	MenuColumnEditor: "Editor di colonna...", MenuReadOnly: "Sola lettura",
	DateShortLayout: "15:04 02/01/2006", DateLongLayout: "15:04:05 02/01/2006 (Monday)",

	MenuFind: "Trova...", MenuFindInFiles: "Trova nei file...", MenuFindNext: "Trova successivo", MenuFindPrev: "Trova precedente",
	MenuSelectFindNext: "Seleziona e trova successivo", MenuSelectFindPrev: "Seleziona e trova precedente", MenuReplace: "Sostituisci...",
	MenuIncremental: "Ricerca incrementale", MenuMark: "Evidenzia...", MenuSearchResults: "Finestra dei risultati",
	MenuNextResult: "Risultato successivo", MenuPrevResult: "Risultato precedente", MenuGoTo: "Vai a...",
	MenuGotoBrace: "Vai alla parentesi corrispondente", MenuSelectBrace: "Seleziona tutto tra le parentesi",
	MenuStyleAll: "Evidenzia tutte le occorrenze", MenuStyleOne: "Evidenzia un'occorrenza", MenuClearStyleSub: "Rimuovi evidenziazione",
	MenuClearAllStyles: "Rimuovi tutte le evidenziazioni", MenuJumpUp: "Salta su", MenuJumpDown: "Salta giù", MenuBookmark: "Segnalibri",
	MenuToggleBookmark: "Attiva/disattiva segnalibro", MenuNextBookmark: "Segnalibro successivo", MenuPrevBookmark: "Segnalibro precedente",
	MenuClearBookmarks: "Rimuovi tutti i segnalibri", MenuCutBookmarked: "Taglia le righe con segnalibro", MenuCopyBookmarked: "Copia le righe con segnalibro",
	MenuRemoveBookmarked: "Rimuovi le righe con segnalibro", MenuRemoveUnbookmarked: "Rimuovi le righe senza segnalibro", MenuInverseBookmarks: "Inverti i segnalibri",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Con lo stile %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Rimuovi lo stile %d", n) },

	MenuAlwaysOnTop: "Sempre in primo piano", MenuFullScreen: "Schermo intero", MenuShowSymbol: "Mostra simboli",
	MenuShowSpaces: "Spazi e tabulazioni", MenuShowEOL: "Fine riga", MenuShowAllChars: "Tutti i caratteri",
	MenuShowIndentGuides: "Guide di rientro", MenuShowWrapSymbol: "Simbolo di a capo", MenuZoom: "Zoom", MenuZoomIn: "Ingrandisci",
	MenuZoomOut: "Riduci", MenuZoomReset: "Zoom predefinito", MenuMoveClone: "Sposta/clona il documento",
	MenuMoveToOtherView: "Sposta nell'altra vista", MenuCloneToOtherView: "Clona nell'altra vista", MenuTab: "Schede",
	MenuNextTab: "Scheda successiva", MenuPrevTab: "Scheda precedente", MenuMoveTabForward: "Sposta la scheda a destra", MenuMoveTabBackward: "Sposta la scheda a sinistra",
	MenuWordWrap: "A capo automatico", MenuLineNumbers: "Numeri di riga", MenuFocusOtherView: "Passa all'altra vista", MenuFoldAll: "Comprimi tutto",
	MenuUnfoldAll: "Espandi tutto", MenuFoldCurrent: "Comprimi il livello attuale", MenuUnfoldCurrent: "Espandi il livello attuale",
	MenuFoldLevel: "Comprimi livello", MenuUnfoldLevel: "Espandi livello", MenuSummary: "Riepilogo...", MenuFolderWorkspace: "Cartella come area di lavoro",
	MenuMonitoring: "Monitoraggio (tail -f)", MenuToolbar: "Barra degli strumenti", MenuStatusBar: "Barra di stato",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Ultima scheda"
		}
		return fmt.Sprintf("Scheda %d", n)
	},

	MenuEncANSI: "Codifica in ANSI", MenuEncUTF8: "Codifica in UTF-8", MenuEncUTF8BOM: "Codifica in UTF-8-BOM", MenuEncUTF16BE: "Codifica in UTF-16 BE BOM",
	MenuEncUTF16LE: "Codifica in UTF-16 LE BOM", MenuCharacterSets: "Set di caratteri", MenuConvANSI: "Converti in ANSI", MenuConvUTF8: "Converti in UTF-8",
	MenuConvUTF8BOM: "Converti in UTF-8-BOM", MenuConvUTF16BE: "Converti in UTF-16 BE BOM", MenuConvUTF16LE: "Converti in UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Arabo", "Baltic": "Baltico", "Celtic": "Celtico", "Cyrillic": "Cirillico",
		"Central European": "Europa centrale", "Chinese": "Cinese", "Greek": "Greco", "Hebrew": "Ebraico", "Japanese": "Giapponese",
		"Korean": "Coreano", "North European": "Europa del nord", "Thai": "Thailandese", "Turkish": "Turco", "Vietnamese": "Vietnamita",
		"Western European": "Europa occidentale"},

	MenuPreferences: "Preferenze...", MenuColorScheme: "Schema di colori", MenuShortcuts: "Scorciatoie da tastiera",
	MenuHashGenerate:  func(name string) string { return "Genera " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Genera " + name + " dei file..." },
	MenuHashSelection: func(name string) string { return "Copia " + name + " della selezione" },
	MenuBase64Encode:  "Codifica Base64", MenuBase64Decode: "Decodifica Base64", MenuURLEncode: "Codifica URL", MenuURLDecode: "Decodifica URL",
	MenuJSONFormat: "Formatta JSON", MenuJSONMinify: "Compatta JSON",
	MenuStartRecording: "Avvia/ferma la registrazione", MenuStopRecording: "Ferma la registrazione", MenuPlayback: "Riproduci",
	MenuSaveMacro: "Salva la macro registrata...", MenuRunMacroMulti: "Esegui la macro più volte...", MenuTrimSave: "Rimuovi spazi a fine riga e salva",
	MenuRun: "Esegui...", MenuOpenInBrowser: "Apri nel browser", MenuSearchInternet: "Cerca su Internet", MenuOpenSelectedFile: "Apri file (nome selezionato)",
	MenuWindows: "Finestre...", MenuSortTabsByName: "Ordina le schede per nome", MenuSortTabsByPath: "Ordina le schede per percorso",
	MenuHelp: "Guida online", MenuHomePage: "Pagina iniziale", MenuAbout: "Informazioni su AltNotepad",

	OpenTitle: "Apri", SaveTitle: "Salva", SaveAsTitle: "Salva con nome", SaveCopyTitle: "Salva una copia con nome", AllFiles: "Tutti i file",
	TextFiles: "File di testo", Save: "Salva", DontSave: "Non salvare", PathLabel: "Percorso del file:",
	SaveChangesAsk: func(path string) string { return "Salvare le modifiche a «" + path + "»?" },
	CreateFileAsk:  func(path string) string { return "«" + path + "» non esiste. Crearlo?" },
	AlreadyOpen:    func(path string) string { return "«" + path + "» è aperto in un'altra scheda" },
	FileExists:     func(path string) string { return "«" + path + "» esiste già" },
	FileNotFound:   func(path string) string { return "File non trovato: " + path },
	Unencodable: func(enc string) string {
		return "Alcuni caratteri non possono essere salvati in " + enc + ". Salvarli come «?»?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Caricamento di %s: %d%%", name, percent) },
	LargeFileMode: "File grande: evidenziazione della sintassi, compressione e a capo disattivati",
	LargeFileMark: "(file grande)",
	ReloadTitle:   "Ricarica", RenameTitle: "Rinomina", NewName: "Nuovo nome:", DeleteFileTitle: "Elimina dal disco",
	ReloadAsk: func(path string) string {
		return "Ricaricare «" + path + "»? Le modifiche non salvate andranno perse."
	},
	DeleteFileAsk:    func(path string) string { return "Eliminare «" + path + "» dal disco?" },
	FileChangedTitle: "File modificato", FileDeletedTitle: "File eliminato",
	FileChangedAsk: func(path string) string {
		return "«" + path + "»\n\nQuesto file è stato modificato da un altro programma.\nRicaricarlo?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "«" + path + "»\n\nQuesto file è stato modificato da un altro programma.\nRicaricarlo perdendo le modifiche dell'editor?"
	},
	FileDeletedAsk: func(path string) string {
		return "«" + path + "»\n\nQuesto file non esiste più.\nMantenerlo nell'editor?"
	},
	NotASession: func(name string) string { return name + " non è un file di sessione" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("lunghezza: %d   righe: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Riga: %d   Col: %d   Pos: %d", line, col, pos) },
	StatusSel:  "Sel:",

	GotoTitle: "Vai a...", GotoLine: "Riga", GotoOffset: "Posizione", Go: "Vai",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Sei qui: %d   Vai a (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Sei qui: %d   Vai a (0 - %d):", cur, total)
	},
	ColumnTitle: "Editor di colonna e selezione multipla", ColumnText: "Testo da inserire", ColumnNumbers: "Numero da inserire",
	ColumnInitial: "Numero iniziale:", ColumnIncrease: "Incremento:", ColumnRepeat: "Ripeti:", ColumnFormat: "Formato:",
	ColumnLeading: "Riempimento:", LeadingNone: "Nessuno", LeadingZeros: "Zeri", LeadingSpaces: "Spazi",
	FullPath: "Percorso completo", Modified: "Modificato", SummaryChars: "Caratteri (senza fine riga)", SummaryWords: "Parole",
	SummaryLines: "Righe", SummaryNonBlankLines: "Righe non vuote", SummaryLength: "Lunghezza del documento", SummarySelected: "Caratteri selezionati",
	Bytes: "byte", RunTitle: "Esegui", RunLabel: "Programma da eseguire:", Run: "Esegui",
	Name: "Nome", State: "Stato", ModifiedState: "modificato", Activate: "Attiva", CloseWindows: "Chiudi",
	Command: "Comando", Shortcut: "Scorciatoia", HashInput: "Testo:", HashEachLine: "Tratta ogni riga separatamente",
	CopyToClipboard:   "Copia",
	CopiedToClipboard: func(s string) string { return "Copiato negli appunti: " + s },
	MacroName:         "Nome della macro:", MacroTimes: "Quante volte eseguire (* - fino alla fine del file):",

	FindTab: "Trova", ReplaceTab: "Sostituisci", FindInFilesTab: "Trova nei file", MarkTab: "Evidenzia", FindWhat: "Trova:",
	ReplaceWith: "Sostituisci con:", Filters: "Filtri:", Directory: "Cartella:", CurrentFolder: "Cartella attuale",
	MatchCase: "Maiuscole/minuscole", WholeWord: "Solo parole intere", WrapAround: "Ricomincia dall'inizio", Backward: "All'indietro",
	InSelection: "Nella selezione", InSubfolders: "In tutte le sottocartelle", InHidden: "Nelle cartelle nascoste", BookmarkLine: "Segna la riga",
	PurgeEach: "Azzera a ogni ricerca", SearchMode: "Modalità di ricerca:", ModeNormal: "Normale", ModeExtended: "Estesa (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Espressione regolare", DotAll: ". include gli a capo",
	FindNext: "Trova successivo", FindPrev: "Trova precedente", Count: "Conta", FindAllCurrent: "Trova tutto nel documento",
	FindAllOpen: "Trova tutto nei documenti aperti", Replace: "Sostituisci", ReplaceAll: "Sostituisci tutto",
	ReplaceAllOpen: "Sostituisci tutto nei documenti aperti", FindAll: "Trova tutto", ReplaceInFiles: "Sostituisci nei file",
	MarkAll: "Evidenzia tutto", ClearMarks: "Rimuovi evidenziazioni", CopyMarked: "Copia il testo evidenziato", Wrapped: "Raggiunta la fine, ripreso dall'inizio",
	ReadOnlyDoc: "Il documento è in sola lettura",
	NotFound:    func(text string) string { return "Impossibile trovare «" + text + "»" },
	CountResult: func(n int) string { return fmt.Sprintf("Totale: %d %s", n, itP(n, "corrispondenza", "corrispondenze")) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, itP(n, "occorrenza sostituita", "occorrenze sostituite"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, itP(n, "corrispondenza evidenziata", "corrispondenze evidenziate"))
	},
	FoundCount: func(n int) string { return fmt.Sprintf("%d %s", n, itP(n, "risultato", "risultati")) },
	HitsCount:  func(n int) string { return fmt.Sprintf("(%d %s)", n, itP(n, "risultato", "risultati")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s in %d file", n, itP(n, "risultato", "risultati"), files)
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s in %d file", n, itP(n, "sostituzione", "sostituzioni"), files)
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Sostituire tutti i «" + find + "» con «" + repl + "» in tutti i file di\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "Cartella inesistente: " + path },
	Searching:     func(path string) string { return "Ricerca: " + path },
	SearchResults: "Risultati della ricerca", Stop: "Ferma", Stopped: "(fermata)", Clear: "Pulisci", Line: "Riga",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Cerca «%s» (%d %s in %d file su %d esaminati)", text, hits, itP(hits, "risultato", "risultati"), files, searched)
	},

	Workspace: "Area di lavoro", Refresh: "Aggiorna", CollapseAll: "Comprimi tutto",

	PageGeneral: "Generale", PageEditing: "Modifica", PageDisplay: "Visualizzazione", PageFiles: "Nuovo documento e file",
	Language: "Lingua:", LanguageSystem: "Come nel sistema", Theme: "Tema:", ThemeDark: "Scuro", ThemeLight: "Chiaro",
	ColorScheme: "Schema di colori:", SchemeByTheme: "In base al tema",
	ShowToolbar: "Mostra la barra degli strumenti", ShowStatusBar: "Mostra la barra di stato", TabCloseButtons: "Pulsanti di chiusura sulle schede",
	DoubleClickCloses: "Il doppio clic chiude la scheda", AlwaysOnTop: "Sempre in primo piano", SingleInstance: "Apri i file nella finestra già aperta",
	RememberSession: "Ricorda i file aperti per la sessione successiva", BackupSession: "Conserva le modifiche non salvate tra le sessioni",
	BackupEvery: "Backup ogni, secondi:", MaxRecent: "File recenti nell'elenco:",
	TabSize: "Dimensione tabulazione:", InsertSpaces: "Sostituisci le tabulazioni con spazi", AutoIndent: "Rientro automatico", AutoClose: "Chiudi parentesi e virgolette",
	AutoCompletion: "Completamento delle parole durante la digitazione", SmartHome: "Home va al primo carattere non vuoto", MultiEdit: "Modifica multipla (Ctrl+clic)",
	ScrollPast: "Scorri oltre l'ultima riga", CopyLineNoSel: "Copia/taglia la riga senza selezione", CaretWidth: "Larghezza del cursore:",
	CaretBlink: "Lampeggio del cursore, ms (0 - nessuno):", EdgeColumn: "Indicatore righe lunghe, colonna (0 - nessuno):", WordChars: "Caratteri di parola aggiuntivi:",
	FontSize: "Dimensione carattere:", FontFile: "File del carattere:", BuiltInFont: "JetBrains Mono (integrato)", Fonts: "Caratteri", LineNumbers: "Numeri di riga",
	BookmarkMargin: "Margine segnalibri", FoldMargin: "Margine di compressione", CurrentLine: "Evidenzia la riga corrente", SmartHighlight: "Evidenziazione intelligente",
	SmartMatchCase: "Evidenziazione intelligente: maiuscole/minuscole", SmartWholeWord: "Evidenziazione intelligente: solo parole intere", BraceMatch: "Evidenzia le parentesi corrispondenti",
	WrapSymbol: "Mostra il simbolo di a capo", ChangeHistory: "Cronologia modifiche nel margine",
	NewDocEOL: "Fine riga dei nuovi documenti:", NewDocEncoding: "Codifica dei nuovi documenti:", ANSICharset: "Set di caratteri ANSI:",
	NewDocLanguage: "Linguaggio dei nuovi documenti:", SystemDefault: "Come nel sistema", ByLanguage: "In base alla lingua dell'interfaccia",
	LargeFileMB: "Limite dei file grandi, MB:", CheckFileChanges: "Rileva i file modificati da altri programmi", AutoReload: "Ricaricali senza chiedere se non modificati",

	AboutTitle:       func(name string) string { return "Informazioni su " + name },
	AboutDescription: "Un editor di testo e codice sorgente", Version: "Versione", Author: "Autore:", License: "Licenza:",
	VisitWebsite: "Visita il sito", Close: "Chiudi",
}
