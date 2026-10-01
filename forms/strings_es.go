package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func esP(n int, one, many string) string { return i18n.Plural("es", n, one, many) }

var es = Strings{
	Error:      "Error",
	NewDocName: "nuevo",
	NormalText: "Texto normal",

	MenuFile: "Archivo", MenuEdit: "Editar", MenuSearch: "Buscar", MenuView: "Ver", MenuEncoding: "Codificación", MenuLanguage: "Lenguaje",
	MenuSettings: "Configuración", MenuTools: "Herramientas", MenuMacro: "Macro", MenuRunMenu: "Ejecutar", MenuWindow: "Ventana",

	MenuNew: "Nuevo", MenuOpen: "Abrir...", MenuOpenContainingFolderSub: "Abrir la carpeta del documento",
	MenuOpenContainingFolder: "Administrador de archivos", MenuOpenDefaultViewer: "Abrir con el programa predeterminado", MenuOpenFolderWorkspace: "Abrir carpeta como espacio de trabajo...",
	MenuReload: "Recargar desde el disco", MenuSave: "Guardar", MenuSaveAs: "Guardar como...", MenuSaveCopy: "Guardar una copia como...", MenuSaveAll: "Guardar todo",
	MenuRename: "Renombrar...", MenuClose: "Cerrar", MenuCloseAll: "Cerrar todo", MenuCloseMore: "Cerrar varios documentos",
	MenuCloseAllButActive: "Cerrar todo excepto el actual", MenuCloseLeft: "Cerrar todo a la izquierda", MenuCloseRight: "Cerrar todo a la derecha",
	MenuCloseUnchanged: "Cerrar los no modificados", MenuDeleteFile: "Eliminar del disco...", MenuLoadSession: "Cargar sesión...",
	MenuSaveSession: "Guardar sesión...", MenuRecentFiles: "Archivos recientes", MenuClearRecent: "Vaciar la lista de archivos recientes",
	MenuRestoreClosed: "Restaurar el último archivo cerrado", MenuExit: "Salir",
	MenuPrint: "Imprimir...", MenuExportHTML: "Exportar a HTML...", MenuFunctionList: "Lista de funciones", Filter: "Filtro", MenuDocumentMap: "Mapa del documento",
	MenuChangeHistory: "Historial de cambios", MenuNextChange: "Ir al cambio siguiente", MenuPrevChange: "Ir al cambio anterior", MenuClearChanges: "Borrar el historial de cambios",
	SavedTo:         func(path string) string { return "Guardado en " + path },
	MenuMultiSelect: "Selección múltiple", MenuMultiAll: "Seleccionar todas las apariciones", MenuMultiAllCase: "Seleccionar todas (coincidir mayúsculas)",
	MenuMultiNext: "Añadir la aparición siguiente", MenuMultiUndo: "Deshacer la última añadida", MenuMultiSkip: "Omitir la actual y añadir la siguiente",

	MenuUndo: "Deshacer", MenuRedo: "Rehacer", MenuCut: "Cortar", MenuCopy: "Copiar", MenuPaste: "Pegar", MenuDelete: "Eliminar", MenuSelectAll: "Seleccionar todo",
	MenuInsert: "Insertar", MenuDateShort: "Fecha y hora (corta)", MenuDateLong: "Fecha y hora (larga)", MenuDateISO: "Fecha y hora (ISO 8601)",
	MenuCopyToClipboard: "Copiar al portapapeles", MenuCopyFullPath: "Ruta completa del archivo", MenuCopyFileName: "Nombre del archivo",
	MenuCopyDirPath: "Ruta de la carpeta", MenuIndentSub: "Sangría", MenuIndent: "Aumentar sangría", MenuUnindent: "Reducir sangría",
	MenuConvertCase: "Convertir a", MenuUpperCase: "MAYÚSCULAS", MenuLowerCase: "minúsculas", MenuProperCase: "Tipo Título",
	MenuProperCaseBlend: "Tipo Título (mezclado)", MenuSentenceCase: "Tipo oración", MenuSentenceCaseBlend: "Tipo oración (mezclado)",
	MenuInvertCase: "iNVERTIR mAYÚSCULAS", MenuRandomCase: "aLeAtOrIo", MenuLineOperations: "Operaciones con líneas",
	MenuDuplicateLine: "Duplicar la línea actual", MenuRemoveDupLines: "Eliminar líneas duplicadas", MenuRemoveConsecutiveDupLines: "Eliminar duplicados consecutivos",
	MenuSplitLines: "Dividir líneas", MenuJoinLines: "Unir líneas", MenuMoveLineUp: "Mover la línea arriba", MenuMoveLineDown: "Mover la línea abajo",
	MenuDeleteLine: "Eliminar la línea actual", MenuCutLine: "Cortar la línea actual", MenuTransposeLine: "Intercambiar con la línea anterior",
	MenuRemoveEmptyLines: "Eliminar líneas vacías", MenuRemoveBlankLines: "Eliminar líneas vacías (también con espacios)",
	MenuInsertLineAbove: "Insertar línea vacía encima", MenuInsertLineBelow: "Insertar línea vacía debajo",
	MenuReverseLines: "Invertir el orden de las líneas", MenuShuffleLines: "Ordenar las líneas al azar",
	MenuSortAsc: "Ordenar líneas ascendente", MenuSortDesc: "Ordenar líneas descendente",
	MenuSortAscCI: "Ordenar ascendente sin distinguir mayúsculas", MenuSortDescCI: "Ordenar descendente sin distinguir mayúsculas",
	MenuSortIntAsc: "Ordenar como enteros ascendente", MenuSortIntDesc: "Ordenar como enteros descendente",
	MenuSortDecAsc: "Ordenar como decimales ascendente", MenuSortDecDesc: "Ordenar como decimales descendente",
	MenuSortLenAsc: "Ordenar por longitud ascendente", MenuSortLenDesc: "Ordenar por longitud descendente",
	MenuComment: "Comentarios", MenuToggleComment: "Alternar comentario de línea", MenuLineComment: "Comentar líneas",
	MenuLineUncomment: "Descomentar líneas", MenuBlockComment: "Comentario de bloque", MenuBlockUncomment: "Quitar comentario de bloque",
	MenuAutoCompletion: "Autocompletado", MenuWordCompletion: "Completar palabra", MenuEOLConversion: "Conversión de fin de línea",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Operaciones con espacios", MenuTrimTrailing: "Quitar espacios al final", MenuTrimLeading: "Quitar espacios al principio",
	MenuTrimBoth: "Quitar espacios al principio y al final", MenuEOLToSpace: "Fin de línea a espacio", MenuRemoveBlankEOL: "Quitar espacios y fines de línea innecesarios",
	MenuTabToSpace: "Tabulación a espacios", MenuSpaceToTabAll: "Espacios a tabulación (todos)", MenuSpaceToTabLeading: "Espacios a tabulación (al principio)",
	MenuColumnEditor: "Editor de columnas...", MenuReadOnly: "Solo lectura",
	DateShortLayout: "15:04 02/01/2006", DateLongLayout: "15:04:05 02/01/2006 (Monday)",

	MenuFind: "Buscar...", MenuFindInFiles: "Buscar en archivos...", MenuFindNext: "Buscar siguiente", MenuFindPrev: "Buscar anterior",
	MenuSelectFindNext: "Seleccionar y buscar siguiente", MenuSelectFindPrev: "Seleccionar y buscar anterior", MenuReplace: "Reemplazar...",
	MenuIncremental: "Búsqueda incremental", MenuMark: "Marcar...", MenuSearchResults: "Ventana de resultados",
	MenuNextResult: "Resultado siguiente", MenuPrevResult: "Resultado anterior", MenuGoTo: "Ir a...",
	MenuGotoBrace: "Ir a la llave correspondiente", MenuSelectBrace: "Seleccionar todo entre llaves",
	MenuStyleAll: "Resaltar todas las apariciones", MenuStyleOne: "Resaltar una aparición", MenuClearStyleSub: "Quitar resaltado",
	MenuClearAllStyles: "Quitar todos los resaltados", MenuJumpUp: "Saltar arriba", MenuJumpDown: "Saltar abajo", MenuBookmark: "Marcadores",
	MenuToggleBookmark: "Poner/quitar marcador", MenuNextBookmark: "Marcador siguiente", MenuPrevBookmark: "Marcador anterior",
	MenuClearBookmarks: "Quitar todos los marcadores", MenuCutBookmarked: "Cortar líneas marcadas", MenuCopyBookmarked: "Copiar líneas marcadas",
	MenuRemoveBookmarked: "Eliminar líneas marcadas", MenuRemoveUnbookmarked: "Eliminar líneas no marcadas", MenuInverseBookmarks: "Invertir marcadores",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Con el estilo %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Quitar el estilo %d", n) },

	MenuAlwaysOnTop: "Siempre visible", MenuFullScreen: "Pantalla completa", MenuShowSymbol: "Mostrar símbolos",
	MenuShowSpaces: "Espacios y tabulaciones", MenuShowEOL: "Fines de línea", MenuShowAllChars: "Todos los caracteres",
	MenuShowIndentGuides: "Guías de sangría", MenuShowWrapSymbol: "Símbolo de ajuste", MenuZoom: "Zoom", MenuZoomIn: "Acercar",
	MenuZoomOut: "Alejar", MenuZoomReset: "Zoom predeterminado", MenuMoveClone: "Mover/clonar el documento",
	MenuMoveToOtherView: "Mover a la otra vista", MenuCloneToOtherView: "Clonar en la otra vista", MenuTab: "Pestañas",
	MenuNextTab: "Pestaña siguiente", MenuPrevTab: "Pestaña anterior", MenuMoveTabForward: "Mover la pestaña a la derecha", MenuMoveTabBackward: "Mover la pestaña a la izquierda",
	MenuWordWrap: "Ajuste de línea", MenuLineNumbers: "Números de línea", MenuFocusOtherView: "Pasar a la otra vista", MenuFoldAll: "Plegar todo",
	MenuUnfoldAll: "Desplegar todo", MenuFoldCurrent: "Plegar el nivel actual", MenuUnfoldCurrent: "Desplegar el nivel actual",
	MenuFoldLevel: "Plegar nivel", MenuUnfoldLevel: "Desplegar nivel", MenuSummary: "Resumen...", MenuFolderWorkspace: "Carpeta como espacio de trabajo",
	MenuMonitoring: "Supervisión (tail -f)", MenuToolbar: "Barra de herramientas", MenuStatusBar: "Barra de estado",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Última pestaña"
		}
		return fmt.Sprintf("Pestaña %d", n)
	},

	MenuEncANSI: "Codificar en ANSI", MenuEncUTF8: "Codificar en UTF-8", MenuEncUTF8BOM: "Codificar en UTF-8-BOM", MenuEncUTF16BE: "Codificar en UTF-16 BE BOM",
	MenuEncUTF16LE: "Codificar en UTF-16 LE BOM", MenuCharacterSets: "Juegos de caracteres", MenuConvANSI: "Convertir a ANSI", MenuConvUTF8: "Convertir a UTF-8",
	MenuConvUTF8BOM: "Convertir a UTF-8-BOM", MenuConvUTF16BE: "Convertir a UTF-16 BE BOM", MenuConvUTF16LE: "Convertir a UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Árabe", "Baltic": "Báltico", "Celtic": "Celta", "Cyrillic": "Cirílico",
		"Central European": "Europa central", "Chinese": "Chino", "Greek": "Griego", "Hebrew": "Hebreo", "Japanese": "Japonés",
		"Korean": "Coreano", "North European": "Europa del norte", "Thai": "Tailandés", "Turkish": "Turco", "Vietnamese": "Vietnamita",
		"Western European": "Europa occidental"},

	MenuPreferences: "Preferencias...", MenuColorScheme: "Esquema de colores", MenuShortcuts: "Atajos de teclado",
	MenuHashGenerate:  func(name string) string { return "Generar " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Generar " + name + " de archivos..." },
	MenuHashSelection: func(name string) string { return "Copiar " + name + " de la selección" },
	MenuBase64Encode:  "Codificar en Base64", MenuBase64Decode: "Decodificar Base64", MenuURLEncode: "Codificar URL", MenuURLDecode: "Decodificar URL",
	MenuJSONFormat: "Formatear JSON", MenuJSONMinify: "Compactar JSON",
	MenuStartRecording: "Iniciar/detener la grabación", MenuStopRecording: "Detener la grabación", MenuPlayback: "Reproducir",
	MenuSaveMacro: "Guardar la macro grabada...", MenuRunMacroMulti: "Ejecutar la macro varias veces...", MenuTrimSave: "Quitar espacios al final y guardar",
	MenuRun: "Ejecutar...", MenuOpenInBrowser: "Abrir en el navegador", MenuSearchInternet: "Buscar en Internet", MenuOpenSelectedFile: "Abrir archivo (nombre seleccionado)",
	MenuWindows: "Ventanas...", MenuSortTabsByName: "Ordenar pestañas por nombre", MenuSortTabsByPath: "Ordenar pestañas por ruta",
	MenuHelp: "Ayuda en línea", MenuHomePage: "Página de inicio", MenuAbout: "Acerca de AltNotepad",

	OpenTitle: "Abrir", SaveTitle: "Guardar", SaveAsTitle: "Guardar como", SaveCopyTitle: "Guardar una copia como", AllFiles: "Todos los archivos",
	TextFiles: "Archivos de texto", Save: "Guardar", DontSave: "No guardar", PathLabel: "Ruta del archivo:",
	SaveChangesAsk: func(path string) string { return "¿Guardar los cambios de «" + path + "»?" },
	CreateFileAsk:  func(path string) string { return "«" + path + "» no existe. ¿Crearlo?" },
	AlreadyOpen:    func(path string) string { return "«" + path + "» está abierto en otra pestaña" },
	FileExists:     func(path string) string { return "«" + path + "» ya existe" },
	FileNotFound:   func(path string) string { return "Archivo no encontrado: " + path },
	Unencodable: func(enc string) string {
		return "Algunos caracteres no se pueden guardar en " + enc + ". ¿Guardarlos como «?»?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Cargando %s: %d%%", name, percent) },
	LargeFileMode: "Archivo grande: resaltado de sintaxis, plegado y ajuste de línea desactivados",
	LargeFileMark: "(archivo grande)",
	ReloadTitle:   "Recargar", RenameTitle: "Renombrar", NewName: "Nuevo nombre:", DeleteFileTitle: "Eliminar del disco",
	ReloadAsk:        func(path string) string { return "¿Recargar «" + path + "»? Se perderán los cambios no guardados." },
	DeleteFileAsk:    func(path string) string { return "¿Eliminar «" + path + "» del disco?" },
	FileChangedTitle: "Archivo modificado", FileDeletedTitle: "Archivo eliminado",
	FileChangedAsk: func(path string) string {
		return "«" + path + "»\n\nOtro programa ha modificado este archivo.\n¿Quiere recargarlo?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "«" + path + "»\n\nOtro programa ha modificado este archivo.\n¿Recargarlo y perder los cambios del editor?"
	},
	FileDeletedAsk: func(path string) string {
		return "«" + path + "»\n\nEste archivo ya no existe.\n¿Mantenerlo en el editor?"
	},
	NotASession: func(name string) string { return name + " no es un archivo de sesión" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("longitud: %d   líneas: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Lín: %d   Col: %d   Pos: %d", line, col, pos) },
	StatusSel:  "Sel:",

	GotoTitle: "Ir a...", GotoLine: "Línea", GotoOffset: "Posición", Go: "Ir",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Está aquí: %d   Ir a (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Está aquí: %d   Ir a (0 - %d):", cur, total)
	},
	ColumnTitle: "Editor de columnas y selección múltiple", ColumnText: "Texto a insertar", ColumnNumbers: "Número a insertar",
	ColumnInitial: "Número inicial:", ColumnIncrease: "Incrementar en:", ColumnRepeat: "Repetir:", ColumnFormat: "Formato:",
	ColumnLeading: "Rellenar con:", LeadingNone: "Nada", LeadingZeros: "Ceros", LeadingSpaces: "Espacios",
	FullPath: "Ruta completa", Modified: "Modificado", SummaryChars: "Caracteres (sin fines de línea)", SummaryWords: "Palabras",
	SummaryLines: "Líneas", SummaryNonBlankLines: "Líneas no vacías", SummaryLength: "Longitud del documento", SummarySelected: "Caracteres seleccionados",
	Bytes: "bytes", RunTitle: "Ejecutar", RunLabel: "Programa a ejecutar:", Run: "Ejecutar",
	Name: "Nombre", State: "Estado", ModifiedState: "modificado", Activate: "Activar", CloseWindows: "Cerrar",
	Command: "Comando", Shortcut: "Atajo", HashInput: "Texto:", HashEachLine: "Tratar cada línea por separado",
	CopyToClipboard:   "Copiar",
	CopiedToClipboard: func(s string) string { return "Copiado al portapapeles: " + s },
	MacroName:         "Nombre de la macro:", MacroTimes: "Cuántas veces ejecutar (* - hasta el final del archivo):",

	FindTab: "Buscar", ReplaceTab: "Reemplazar", FindInFilesTab: "Buscar en archivos", MarkTab: "Marcar", FindWhat: "Buscar:",
	ReplaceWith: "Reemplazar con:", Filters: "Filtros:", Directory: "Carpeta:", CurrentFolder: "Carpeta actual",
	MatchCase: "Coincidir mayúsculas", WholeWord: "Solo palabras completas", WrapAround: "Dar la vuelta", Backward: "Hacia atrás",
	InSelection: "En la selección", InSubfolders: "En todas las subcarpetas", InHidden: "En carpetas ocultas", BookmarkLine: "Marcar la línea",
	PurgeEach: "Limpiar en cada búsqueda", SearchMode: "Modo de búsqueda:", ModeNormal: "Normal", ModeExtended: "Extendido (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Expresión regular", DotAll: ". incluye saltos de línea",
	FindNext: "Buscar siguiente", FindPrev: "Buscar anterior", Count: "Contar", FindAllCurrent: "Buscar todo en el documento actual",
	FindAllOpen: "Buscar todo en los documentos abiertos", Replace: "Reemplazar", ReplaceAll: "Reemplazar todo",
	ReplaceAllOpen: "Reemplazar todo en los documentos abiertos", FindAll: "Buscar todo", ReplaceInFiles: "Reemplazar en archivos",
	MarkAll: "Marcar todo", ClearMarks: "Quitar las marcas", CopyMarked: "Copiar el texto marcado", Wrapped: "Fin alcanzado, se continuó desde el principio",
	ReadOnlyDoc: "El documento es de solo lectura",
	NotFound:    func(text string) string { return "No se encuentra «" + text + "»" },
	CountResult: func(n int) string { return fmt.Sprintf("Total: %d %s", n, esP(n, "coincidencia", "coincidencias")) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, esP(n, "aparición reemplazada", "apariciones reemplazadas"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, esP(n, "coincidencia marcada", "coincidencias marcadas"))
	},
	FoundCount: func(n int) string { return fmt.Sprintf("%d %s", n, esP(n, "resultado", "resultados")) },
	HitsCount:  func(n int) string { return fmt.Sprintf("(%d %s)", n, esP(n, "resultado", "resultados")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s en %d %s", n, esP(n, "resultado", "resultados"), files, esP(files, "archivo", "archivos"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s en %d %s", n, esP(n, "reemplazo", "reemplazos"), files, esP(files, "archivo", "archivos"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "¿Reemplazar todos los «" + find + "» por «" + repl + "» en todos los archivos de\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "No existe la carpeta: " + path },
	Searching:     func(path string) string { return "Buscando: " + path },
	SearchResults: "Resultados de la búsqueda", Stop: "Detener", Stopped: "(detenida)", Clear: "Limpiar", Line: "Línea",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Buscar «%s» (%d %s en %d %s de %d revisados)", text, hits, esP(hits, "resultado", "resultados"), files, esP(files, "archivo", "archivos"), searched)
	},

	Workspace: "Espacio de trabajo", Refresh: "Actualizar", CollapseAll: "Contraer todo",

	PageGeneral: "General", PageEditing: "Edición", PageDisplay: "Visualización", PageFiles: "Nuevo documento y archivos",
	Language: "Idioma:", LanguageSystem: "Como en el sistema", Theme: "Tema:", ThemeDark: "Oscuro", ThemeLight: "Claro",
	ColorScheme: "Esquema de colores:", SchemeByTheme: "Según el tema",
	ShowToolbar: "Mostrar la barra de herramientas", ShowStatusBar: "Mostrar la barra de estado", TabCloseButtons: "Botones de cierre en las pestañas",
	DoubleClickCloses: "El doble clic cierra la pestaña", AlwaysOnTop: "Siempre visible", SingleInstance: "Abrir los archivos en la ventana existente",
	RememberSession: "Recordar los archivos abiertos para la próxima sesión", BackupSession: "Conservar los cambios sin guardar entre sesiones",
	BackupEvery: "Copia de seguridad cada, segundos:", MaxRecent: "Archivos recientes en la lista:",
	TabSize: "Tamaño de tabulación:", InsertSpaces: "Sustituir tabulaciones por espacios", AutoIndent: "Sangría automática", AutoClose: "Cerrar paréntesis y comillas",
	AutoCompletion: "Completar palabras al escribir", SmartHome: "Inicio va al primer carácter no blanco", MultiEdit: "Edición múltiple (Ctrl+clic)",
	ScrollPast: "Desplazar más allá de la última línea", CopyLineNoSel: "Copiar/cortar la línea sin selección", CaretWidth: "Ancho del cursor:",
	CaretBlink: "Parpadeo del cursor, ms (0 - sin parpadeo):", EdgeColumn: "Marca de líneas largas, columna (0 - ninguna):", WordChars: "Caracteres de palabra adicionales:",
	FontSize: "Tamaño de fuente:", FontFile: "Archivo de fuente:", BuiltInFont: "JetBrains Mono (integrada)", Fonts: "Fuentes", LineNumbers: "Números de línea",
	BookmarkMargin: "Margen de marcadores", FoldMargin: "Margen de plegado", CurrentLine: "Resaltar la línea actual", SmartHighlight: "Resaltado inteligente",
	SmartMatchCase: "Resaltado inteligente: coincidir mayúsculas", SmartWholeWord: "Resaltado inteligente: solo palabras completas", BraceMatch: "Resaltar llaves correspondientes",
	WrapSymbol: "Mostrar el símbolo de ajuste", ChangeHistory: "Historial de cambios en el margen",
	NewDocEOL: "Fin de línea de documentos nuevos:", NewDocEncoding: "Codificación de documentos nuevos:", ANSICharset: "Juego de caracteres ANSI:",
	NewDocLanguage: "Lenguaje de documentos nuevos:", SystemDefault: "Como en el sistema", ByLanguage: "Según el idioma de la interfaz",
	LargeFileMB: "Límite de archivos grandes, MB:", CheckFileChanges: "Detectar archivos modificados por otros programas", AutoReload: "Recargarlos sin preguntar si no están modificados",

	AboutTitle:       func(name string) string { return "Acerca de " + name },
	AboutDescription: "Un editor de texto y código fuente", Version: "Versión", Author: "Autor:", License: "Licencia:",
	VisitWebsite: "Visitar el sitio web", Close: "Cerrar",
}
