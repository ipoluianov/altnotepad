package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func srP(n int, one, few, many string) string { return i18n.Plural("sr", n, one, few, many) }

var sr = Strings{
	Error:      "Грешка",
	NewDocName: "нови",
	NormalText: "Обичан текст",

	MenuFile: "Датотека", MenuEdit: "Уређивање", MenuSearch: "Претрага", MenuView: "Приказ", MenuEncoding: "Кодирање", MenuLanguage: "Синтакса",
	MenuSettings: "Подешавања", MenuTools: "Алатке", MenuMacro: "Макро", MenuRunMenu: "Покрени", MenuWindow: "Прозор",

	MenuNew: "Нова", MenuOpen: "Отвори...", MenuOpenContainingFolderSub: "Отвори фасциклу документа",
	MenuOpenContainingFolder: "Менаџер датотека", MenuOpenDefaultViewer: "Отвори подразумеваним програмом", MenuOpenFolderWorkspace: "Отвори фасциклу као радни простор...",
	MenuReload: "Поново учитај са диска", MenuSave: "Сачувај", MenuSaveAs: "Сачувај као...", MenuSaveCopy: "Сачувај копију као...", MenuSaveAll: "Сачувај све",
	MenuRename: "Преименуј...", MenuClose: "Затвори", MenuCloseAll: "Затвори све", MenuCloseMore: "Затвори више докумената",
	MenuCloseAllButActive: "Затвори све осим тренутног", MenuCloseLeft: "Затвори све лево", MenuCloseRight: "Затвори све десно",
	MenuCloseUnchanged: "Затвори непромењене", MenuDeleteFile: "Обриши са диска...", MenuLoadSession: "Учитај сесију...",
	MenuSaveSession: "Сачувај сесију...", MenuRecentFiles: "Недавне датотеке", MenuClearRecent: "Очисти листу недавних датотека",
	MenuRestoreClosed: "Врати последњу затворену датотеку", MenuExit: "Излаз",
	MenuPrint: "Штампај...", MenuExportHTML: "Извези у HTML...", MenuFunctionList: "Листа функција", Filter: "Филтер", MenuDocumentMap: "Мапа документа",
	MenuChangeHistory: "Историја измена", MenuNextChange: "Иди на следећу измену", MenuPrevChange: "Иди на претходну измену", MenuClearChanges: "Очисти историју измена",
	SavedTo:         func(path string) string { return "Сачувано у " + path },
	MenuMultiSelect: "Вишеструки избор", MenuMultiAll: "Изабери сва појављивања", MenuMultiAllCase: "Изабери сва (разликуј велика слова)",
	MenuMultiNext: "Додај следеће појављивање", MenuMultiUndo: "Поништи последње додато", MenuMultiSkip: "Прескочи тренутно и додај следеће",

	MenuUndo: "Опозови", MenuRedo: "Понови", MenuCut: "Исеци", MenuCopy: "Копирај", MenuPaste: "Налепи", MenuDelete: "Обриши", MenuSelectAll: "Изабери све",
	MenuInsert: "Уметни", MenuDateShort: "Датум и време (кратко)", MenuDateLong: "Датум и време (дуго)", MenuDateISO: "Датум и време (ISO 8601)",
	MenuCopyToClipboard: "Копирај у оставу", MenuCopyFullPath: "Пуна путања датотеке", MenuCopyFileName: "Име датотеке",
	MenuCopyDirPath: "Путања фасцикле", MenuIndentSub: "Увлачење", MenuIndent: "Повећај увлачење", MenuUnindent: "Смањи увлачење",
	MenuConvertCase: "Претвори у", MenuUpperCase: "ВЕЛИКА СЛОВА", MenuLowerCase: "мала слова", MenuProperCase: "Свака Реч Великим",
	MenuProperCaseBlend: "Свака Реч Великим (мешано)", MenuSentenceCase: "Као у реченици", MenuSentenceCaseBlend: "Као у реченици (мешано)",
	MenuInvertCase: "оБРНИ вЕЛИЧИНУ", MenuRandomCase: "нАсУмИчНо", MenuLineOperations: "Операције са редовима",
	MenuDuplicateLine: "Удвостручи тренутни ред", MenuRemoveDupLines: "Уклони дупле редове", MenuRemoveConsecutiveDupLines: "Уклони узастопне дупле редове",
	MenuSplitLines: "Подели редове", MenuJoinLines: "Спој редове", MenuMoveLineUp: "Помери ред горе", MenuMoveLineDown: "Помери ред доле",
	MenuDeleteLine: "Обриши тренутни ред", MenuCutLine: "Исеци тренутни ред", MenuTransposeLine: "Замени са претходним редом",
	MenuRemoveEmptyLines: "Уклони празне редове", MenuRemoveBlankLines: "Уклони празне редове (и са размацима)",
	MenuInsertLineAbove: "Уметни празан ред изнад", MenuInsertLineBelow: "Уметни празан ред испод",
	MenuReverseLines: "Обрни редослед редова", MenuShuffleLines: "Насумично измешај редове",
	MenuSortAsc: "Сортирај редове растуће", MenuSortDesc: "Сортирај редове опадајуће",
	MenuSortAscCI: "Сортирај растуће без обзира на велика слова", MenuSortDescCI: "Сортирај опадајуће без обзира на велика слова",
	MenuSortIntAsc: "Сортирај као целе бројеве растуће", MenuSortIntDesc: "Сортирај као целе бројеве опадајуће",
	MenuSortDecAsc: "Сортирај као децималне бројеве растуће", MenuSortDecDesc: "Сортирај као децималне бројеве опадајуће",
	MenuSortLenAsc: "Сортирај по дужини растуће", MenuSortLenDesc: "Сортирај по дужини опадајуће",
	MenuComment: "Коментари", MenuToggleComment: "Укључи/искључи коментар реда", MenuLineComment: "Коментариши редове",
	MenuLineUncomment: "Уклони коментар редова", MenuBlockComment: "Блок коментар", MenuBlockUncomment: "Уклони блок коментар",
	MenuAutoCompletion: "Аутоматско довршавање", MenuWordCompletion: "Доврши реч", MenuEOLConversion: "Претварање краја реда",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Операције са размацима", MenuTrimTrailing: "Уклони размаке на крају", MenuTrimLeading: "Уклони размаке на почетку",
	MenuTrimBoth: "Уклони размаке на почетку и крају", MenuEOLToSpace: "Крај реда у размак", MenuRemoveBlankEOL: "Уклони сувишне размаке и крајеве редова",
	MenuTabToSpace: "Табулатори у размаке", MenuSpaceToTabAll: "Размаци у табулаторе (сви)", MenuSpaceToTabLeading: "Размаци у табулаторе (почетни)",
	MenuColumnEditor: "Уређивач колона...", MenuReadOnly: "Само за читање",
	DateShortLayout: "15:04 02.01.2006.", DateLongLayout: "15:04:05 02.01.2006. (Monday)",

	MenuFind: "Пронађи...", MenuFindInFiles: "Пронађи у датотекама...", MenuFindNext: "Пронађи следеће", MenuFindPrev: "Пронађи претходно",
	MenuSelectFindNext: "Изабери и пронађи следеће", MenuSelectFindPrev: "Изабери и пронађи претходно", MenuReplace: "Замени...",
	MenuIncremental: "Постепена претрага", MenuMark: "Означи...", MenuSearchResults: "Прозор резултата претраге",
	MenuNextResult: "Следећи резултат", MenuPrevResult: "Претходни резултат", MenuGoTo: "Иди на...",
	MenuGotoBrace: "Иди на одговарајућу заграду", MenuSelectBrace: "Изабери све између заграда",
	MenuStyleAll: "Истакни сва појављивања", MenuStyleOne: "Истакни једно појављивање", MenuClearStyleSub: "Уклони истицање",
	MenuClearAllStyles: "Уклони сва истицања", MenuJumpUp: "Скочи горе", MenuJumpDown: "Скочи доле", MenuBookmark: "Обележивачи",
	MenuToggleBookmark: "Постави/уклони обележивач", MenuNextBookmark: "Следећи обележивач", MenuPrevBookmark: "Претходни обележивач",
	MenuClearBookmarks: "Уклони све обележиваче", MenuCutBookmarked: "Исеци обележене редове", MenuCopyBookmarked: "Копирај обележене редове",
	MenuRemoveBookmarked: "Уклони обележене редове", MenuRemoveUnbookmarked: "Уклони необележене редове", MenuInverseBookmarks: "Обрни обележиваче",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Стил %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Уклони стил %d", n) },

	MenuAlwaysOnTop: "Увек на врху", MenuFullScreen: "Цео екран", MenuShowSymbol: "Прикажи симболе",
	MenuShowSpaces: "Размаци и табулатори", MenuShowEOL: "Крајеви редова", MenuShowAllChars: "Сви знакови",
	MenuShowIndentGuides: "Водичи увлачења", MenuShowWrapSymbol: "Симбол преламања", MenuZoom: "Увећање", MenuZoomIn: "Увећај",
	MenuZoomOut: "Умањи", MenuZoomReset: "Подразумевано увећање", MenuMoveClone: "Премести/клонирај документ",
	MenuMoveToOtherView: "Премести у други приказ", MenuCloneToOtherView: "Клонирај у други приказ", MenuTab: "Картице",
	MenuNextTab: "Следећа картица", MenuPrevTab: "Претходна картица", MenuMoveTabForward: "Помери картицу десно", MenuMoveTabBackward: "Помери картицу лево",
	MenuWordWrap: "Преламање редова", MenuLineNumbers: "Бројеви редова", MenuFocusOtherView: "Пређи у други приказ", MenuFoldAll: "Скупи све",
	MenuUnfoldAll: "Рашири све", MenuFoldCurrent: "Скупи тренутни ниво", MenuUnfoldCurrent: "Рашири тренутни ниво",
	MenuFoldLevel: "Скупи ниво", MenuUnfoldLevel: "Рашири ниво", MenuSummary: "Резиме...", MenuFolderWorkspace: "Фасцикла као радни простор",
	MenuMonitoring: "Праћење (tail -f)", MenuToolbar: "Трака са алаткама", MenuStatusBar: "Статусна трака",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Последња картица"
		}
		return fmt.Sprintf("Картица %d", n)
	},

	MenuEncANSI: "Кодирај у ANSI", MenuEncUTF8: "Кодирај у UTF-8", MenuEncUTF8BOM: "Кодирај у UTF-8-BOM", MenuEncUTF16BE: "Кодирај у UTF-16 BE BOM",
	MenuEncUTF16LE: "Кодирај у UTF-16 LE BOM", MenuCharacterSets: "Скупови знакова", MenuConvANSI: "Претвори у ANSI", MenuConvUTF8: "Претвори у UTF-8",
	MenuConvUTF8BOM: "Претвори у UTF-8-BOM", MenuConvUTF16BE: "Претвори у UTF-16 BE BOM", MenuConvUTF16LE: "Претвори у UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Арапски", "Baltic": "Балтички", "Celtic": "Келтски", "Cyrillic": "Ћирилица",
		"Central European": "Средњоевропски", "Chinese": "Кинески", "Greek": "Грчки", "Hebrew": "Хебрејски", "Japanese": "Јапански",
		"Korean": "Корејски", "North European": "Северноевропски", "Thai": "Тајландски", "Turkish": "Турски", "Vietnamese": "Вијетнамски",
		"Western European": "Западноевропски"},

	MenuPreferences: "Подешавања...", MenuColorScheme: "Шема боја", MenuShortcuts: "Пречице на тастатури",
	MenuHashGenerate:  func(name string) string { return "Генериши " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Генериши " + name + " датотека..." },
	MenuHashSelection: func(name string) string { return "Копирај " + name + " избора" },
	MenuBase64Encode:  "Кодирај Base64", MenuBase64Decode: "Декодирај Base64", MenuURLEncode: "Кодирај URL", MenuURLDecode: "Декодирај URL",
	MenuJSONFormat: "Форматирај JSON", MenuJSONMinify: "Сажми JSON",
	MenuStartRecording: "Покрени/заустави снимање", MenuStopRecording: "Заустави снимање", MenuPlayback: "Репродукуј",
	MenuSaveMacro: "Сачувај снимљени макро...", MenuRunMacroMulti: "Покрени макро више пута...", MenuTrimSave: "Уклони размаке на крају и сачувај",
	MenuRun: "Покрени...", MenuOpenInBrowser: "Отвори у прегледачу", MenuSearchInternet: "Претражи интернет", MenuOpenSelectedFile: "Отвори датотеку (изабрано име)",
	MenuWindows: "Прозори...", MenuSortTabsByName: "Сортирај картице по имену", MenuSortTabsByPath: "Сортирај картице по путањи",
	MenuHelp: "Помоћ на мрежи", MenuHomePage: "Почетна страница", MenuAbout: "О програму AltNotepad",

	OpenTitle: "Отвори", SaveTitle: "Сачувај", SaveAsTitle: "Сачувај као", SaveCopyTitle: "Сачувај копију као", AllFiles: "Све датотеке",
	TextFiles: "Текстуалне датотеке", Save: "Сачувај", DontSave: "Не чувај", PathLabel: "Путања датотеке:",
	SaveChangesAsk: func(path string) string { return "Сачувати измене у „" + path + "“?" },
	CreateFileAsk:  func(path string) string { return "„" + path + "“ не постоји. Направити?" },
	AlreadyOpen: func(path string) string {
		return "„" + path + "“ је отворена у другој картици"
	},
	FileExists:   func(path string) string { return "„" + path + "“ већ постоји" },
	FileNotFound: func(path string) string { return "Датотека није пронађена: " + path },
	Unencodable: func(enc string) string {
		return "Неке знакове није могуће сачувати у " + enc + ". Сачувати их као „?“?"
	},
	Loading: func(name string, percent int) string {
		return fmt.Sprintf("Учитавање %s: %d%%", name, percent)
	},
	LargeFileMode: "Велика датотека: истицање синтаксе, скупљање и преламање редова су искључени",
	LargeFileMark: "(велика датотека)",
	ReloadTitle:   "Поново учитај", RenameTitle: "Преименуј", NewName: "Ново име:", DeleteFileTitle: "Обриши са диска",
	ReloadAsk: func(path string) string {
		return "Поново учитати „" + path + "“? Несачуване измене ће бити изгубљене."
	},
	DeleteFileAsk:    func(path string) string { return "Обрисати „" + path + "“ са диска?" },
	FileChangedTitle: "Датотека је измењена", FileDeletedTitle: "Датотека је обрисана",
	FileChangedAsk: func(path string) string {
		return "„" + path + "“\n\nОву датотеку је изменио други програм.\nПоново је учитати?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "„" + path + "“\n\nОву датотеку је изменио други програм.\nПоново је учитати и изгубити измене у уређивачу?"
	},
	FileDeletedAsk: func(path string) string {
		return "„" + path + "“\n\nОва датотека више не постоји.\nЗадржати је у уређивачу?"
	},
	NotASession: func(name string) string { return name + " није датотека сесије" },

	StatusSize: func(length, lines int) string {
		return fmt.Sprintf("дужина: %d   редова: %d", length, lines)
	},
	StatusPos: func(line, col, pos int) string {
		return fmt.Sprintf("Ред: %d   Кол: %d   Поз: %d", line, col, pos)
	},
	StatusSel: "Изб:",

	GotoTitle: "Иди на...", GotoLine: "Ред", GotoOffset: "Позиција", Go: "Иди",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Налазите се: %d   Иди на (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Налазите се: %d   Иди на (0 - %d):", cur, total)
	},
	ColumnTitle: "Уређивач колона и вишеструки избор", ColumnText: "Текст за уметање", ColumnNumbers: "Број за уметање",
	ColumnInitial: "Почетни број:", ColumnIncrease: "Повећавај за:", ColumnRepeat: "Понови:", ColumnFormat: "Формат:",
	ColumnLeading: "Попуна:", LeadingNone: "Без", LeadingZeros: "Нуле", LeadingSpaces: "Размаци",
	FullPath: "Пуна путања", Modified: "Измењено", SummaryChars: "Знакови (без крајева редова)", SummaryWords: "Речи",
	SummaryLines: "Редови", SummaryNonBlankLines: "Непразни редови", SummaryLength: "Дужина документа", SummarySelected: "Изабрани знакови",
	Bytes: "бајтова", RunTitle: "Покрени", RunLabel: "Програм за покретање:", Run: "Покрени",
	Name: "Име", State: "Стање", ModifiedState: "измењено", Activate: "Активирај", CloseWindows: "Затвори",
	Command: "Команда", Shortcut: "Пречица", HashInput: "Текст:", HashEachLine: "Обради сваки ред засебно",
	CopyToClipboard:   "Копирај",
	CopiedToClipboard: func(s string) string { return "Копирано у оставу: " + s },
	MacroName:         "Име макроа:", MacroTimes: "Колико пута покренути (* - до краја датотеке):",

	FindTab: "Пронађи", ReplaceTab: "Замени", FindInFilesTab: "Пронађи у датотекама", MarkTab: "Означи", FindWhat: "Пронађи:",
	ReplaceWith: "Замени са:", Filters: "Филтери:", Directory: "Фасцикла:", CurrentFolder: "Тренутна фасцикла",
	MatchCase: "Разликуј велика слова", WholeWord: "Само целе речи", WrapAround: "Настави од почетка", Backward: "Уназад",
	InSelection: "У избору", InSubfolders: "У свим потфасциклама", InHidden: "У скривеним фасциклама", BookmarkLine: "Обележи ред",
	PurgeEach: "Очисти при свакој претрази", SearchMode: "Режим претраге:", ModeNormal: "Обичан", ModeExtended: "Проширен (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Регуларни израз", DotAll: ". обухвата нове редове",
	FindNext: "Пронађи следеће", FindPrev: "Пронађи претходно", Count: "Преброј", FindAllCurrent: "Пронађи све у тренутном документу",
	FindAllOpen: "Пронађи све у отвореним документима", Replace: "Замени", ReplaceAll: "Замени све",
	ReplaceAllOpen: "Замени све у отвореним документима", FindAll: "Пронађи све", ReplaceInFiles: "Замени у датотекама",
	MarkAll: "Означи све", ClearMarks: "Уклони ознаке", CopyMarked: "Копирај означени текст", Wrapped: "Достигнут крај, настављено од почетка",
	ReadOnlyDoc: "Документ је само за читање",
	NotFound:    func(text string) string { return "Није пронађено „" + text + "“" },
	CountResult: func(n int) string {
		return fmt.Sprintf("Укупно: %d %s", n, srP(n, "поклапање", "поклапања", "поклапања"))
	},
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("Замењено: %d %s", n, srP(n, "појављивање", "појављивања", "појављивања"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("Означено: %d %s", n, srP(n, "поклапање", "поклапања", "поклапања"))
	},
	FoundCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, srP(n, "резултат", "резултата", "резултата"))
	},
	HitsCount: func(n int) string {
		return fmt.Sprintf("(%d %s)", n, srP(n, "резултат", "резултата", "резултата"))
	},
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s у %d %s", n, srP(n, "резултат", "резултата", "резултата"), files, srP(files, "датотеци", "датотеке", "датотека"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s у %d %s", n, srP(n, "замена", "замене", "замена"), files, srP(files, "датотеци", "датотеке", "датотека"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Заменити сва „" + find + "“ са „" + repl + "“ у свим датотекама у\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "Фасцикла не постоји: " + path },
	Searching:     func(path string) string { return "Претрага: " + path },
	SearchResults: "Резултати претраге", Stop: "Заустави", Stopped: "(заустављено)", Clear: "Очисти", Line: "Ред",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Претрага „%s“ (%d %s у %d од %d претражених датотека)", text, hits, srP(hits, "резултат", "резултата", "резултата"), files, searched)
	},

	Workspace: "Радни простор", Refresh: "Освежи", CollapseAll: "Скупи све",

	PageGeneral: "Опште", PageEditing: "Уређивање", PageDisplay: "Приказ", PageFiles: "Нови документ и датотеке",
	Language: "Језик:", LanguageSystem: "Као у систему", Theme: "Тема:", ThemeDark: "Тамна", ThemeLight: "Светла",
	ColorScheme: "Шема боја:", SchemeByTheme: "Према теми",
	ShowToolbar: "Прикажи траку са алаткама", ShowStatusBar: "Прикажи статусну траку", TabCloseButtons: "Дугмад за затварање на картицама",
	DoubleClickCloses: "Двоклик затвара картицу", AlwaysOnTop: "Увек на врху", SingleInstance: "Отварај датотеке у већ отвореном прозору",
	RememberSession: "Запамти отворене датотеке за следећу сесију", BackupSession: "Чувај несачуване измене између сесија",
	BackupEvery: "Резервна копија сваких, секунди:", MaxRecent: "Недавних датотека у листи:",
	TabSize: "Величина табулатора:", InsertSpaces: "Замењуј табулаторе размацима", AutoIndent: "Аутоматско увлачење", AutoClose: "Затварај заграде и наводнике",
	AutoCompletion: "Довршавај речи током куцања", SmartHome: "Home иде на први знак који није размак", MultiEdit: "Вишеструко уређивање (Ctrl+клик)",
	ScrollPast: "Померај иза последњег реда", CopyLineNoSel: "Копирај/исеци ред без избора", CaretWidth: "Ширина курсора:",
	CaretBlink: "Трептање курсора, ms (0 - без трептања):", EdgeColumn: "Ознака дугих редова, колона (0 - без):", WordChars: "Додатни знакови речи:",
	FontSize: "Величина фонта:", FontFile: "Датотека фонта:", BuiltInFont: "JetBrains Mono (уграђени)", Fonts: "Фонтови", LineNumbers: "Бројеви редова",
	BookmarkMargin: "Маргина обележивача", FoldMargin: "Маргина скупљања", CurrentLine: "Истакни тренутни ред", SmartHighlight: "Паметно истицање",
	SmartMatchCase: "Паметно истицање: разликуј велика слова", SmartWholeWord: "Паметно истицање: само целе речи", BraceMatch: "Истакни одговарајуће заграде",
	WrapSymbol: "Прикажи симбол преламања", ChangeHistory: "Историја измена на маргини",
	NewDocEOL: "Крај реда нових докумената:", NewDocEncoding: "Кодирање нових докумената:", ANSICharset: "ANSI скуп знакова:",
	NewDocLanguage: "Синтакса нових докумената:", SystemDefault: "Као у систему", ByLanguage: "Према језику интерфејса",
	LargeFileMB: "Праг велике датотеке, MB:", CheckFileChanges: "Откривај датотеке које су изменили други програми", AutoReload: "Поново их учитај без питања ако нису измењене",

	AboutTitle:       func(name string) string { return "О програму " + name },
	AboutDescription: "Уређивач текста и изворног кода", Version: "Верзија", Author: "Аутор:", License: "Лиценца:",
	VisitWebsite: "Посети сајт", Close: "Затвори",
}
