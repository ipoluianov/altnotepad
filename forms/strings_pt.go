package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func ptP(n int, one, many string) string { return i18n.Plural("pt", n, one, many) }

var pt = Strings{
	Error:      "Erro",
	NewDocName: "novo",
	NormalText: "Texto normal",

	MenuFile: "Arquivo", MenuEdit: "Editar", MenuSearch: "Pesquisar", MenuView: "Exibir", MenuEncoding: "Codificação", MenuLanguage: "Linguagem",
	MenuSettings: "Configurações", MenuTools: "Ferramentas", MenuMacro: "Macro", MenuRunMenu: "Executar", MenuWindow: "Janela",

	MenuNew: "Novo", MenuOpen: "Abrir...", MenuOpenContainingFolderSub: "Abrir a pasta do documento",
	MenuOpenContainingFolder: "Gerenciador de arquivos", MenuOpenDefaultViewer: "Abrir com o programa padrão", MenuOpenFolderWorkspace: "Abrir pasta como área de trabalho...",
	MenuReload: "Recarregar do disco", MenuSave: "Salvar", MenuSaveAs: "Salvar como...", MenuSaveCopy: "Salvar uma cópia como...", MenuSaveAll: "Salvar tudo",
	MenuRename: "Renomear...", MenuClose: "Fechar", MenuCloseAll: "Fechar tudo", MenuCloseMore: "Fechar vários documentos",
	MenuCloseAllButActive: "Fechar tudo exceto o atual", MenuCloseLeft: "Fechar tudo à esquerda", MenuCloseRight: "Fechar tudo à direita",
	MenuCloseUnchanged: "Fechar os não modificados", MenuDeleteFile: "Excluir do disco...", MenuLoadSession: "Carregar sessão...",
	MenuSaveSession: "Salvar sessão...", MenuRecentFiles: "Arquivos recentes", MenuClearRecent: "Limpar a lista de arquivos recentes",
	MenuRestoreClosed: "Reabrir o último arquivo fechado", MenuExit: "Sair",
	MenuPrint: "Imprimir...", MenuExportHTML: "Exportar para HTML...", MenuFunctionList: "Lista de funções", Filter: "Filtro", MenuDocumentMap: "Mapa do documento",
	MenuChangeHistory: "Histórico de alterações", MenuNextChange: "Ir para a próxima alteração", MenuPrevChange: "Ir para a alteração anterior", MenuClearChanges: "Limpar o histórico de alterações",
	SavedTo:         func(path string) string { return "Salvo em " + path },
	MenuMultiSelect: "Seleção múltipla", MenuMultiAll: "Selecionar todas as ocorrências", MenuMultiAllCase: "Selecionar todas (diferenciar maiúsculas)",
	MenuMultiNext: "Adicionar a próxima ocorrência", MenuMultiUndo: "Desfazer a última adicionada", MenuMultiSkip: "Pular e adicionar a próxima",

	MenuUndo: "Desfazer", MenuRedo: "Refazer", MenuCut: "Recortar", MenuCopy: "Copiar", MenuPaste: "Colar", MenuDelete: "Excluir", MenuSelectAll: "Selecionar tudo",
	MenuInsert: "Inserir", MenuDateShort: "Data e hora (curta)", MenuDateLong: "Data e hora (longa)", MenuDateISO: "Data e hora (ISO 8601)",
	MenuCopyToClipboard: "Copiar para a área de transferência", MenuCopyFullPath: "Caminho completo do arquivo", MenuCopyFileName: "Nome do arquivo",
	MenuCopyDirPath: "Caminho da pasta", MenuIndentSub: "Recuo", MenuIndent: "Aumentar recuo", MenuUnindent: "Diminuir recuo",
	MenuConvertCase: "Converter para", MenuUpperCase: "MAIÚSCULAS", MenuLowerCase: "minúsculas", MenuProperCase: "Iniciais Maiúsculas",
	MenuProperCaseBlend: "Iniciais Maiúsculas (misto)", MenuSentenceCase: "Primeira letra da frase", MenuSentenceCaseBlend: "Primeira letra da frase (misto)",
	MenuInvertCase: "iNVERTER mAIÚSCULAS", MenuRandomCase: "aLeAtÓrIo", MenuLineOperations: "Operações com linhas",
	MenuDuplicateLine: "Duplicar a linha atual", MenuRemoveDupLines: "Remover linhas duplicadas", MenuRemoveConsecutiveDupLines: "Remover duplicadas consecutivas",
	MenuSplitLines: "Dividir linhas", MenuJoinLines: "Juntar linhas", MenuMoveLineUp: "Mover a linha para cima", MenuMoveLineDown: "Mover a linha para baixo",
	MenuDeleteLine: "Excluir a linha atual", MenuCutLine: "Recortar a linha atual", MenuTransposeLine: "Trocar com a linha anterior",
	MenuRemoveEmptyLines: "Remover linhas vazias", MenuRemoveBlankLines: "Remover linhas vazias (também com espaços)",
	MenuInsertLineAbove: "Inserir linha vazia acima", MenuInsertLineBelow: "Inserir linha vazia abaixo",
	MenuReverseLines: "Inverter a ordem das linhas", MenuShuffleLines: "Embaralhar as linhas",
	MenuSortAsc: "Ordenar linhas em ordem crescente", MenuSortDesc: "Ordenar linhas em ordem decrescente",
	MenuSortAscCI: "Ordem crescente ignorando maiúsculas", MenuSortDescCI: "Ordem decrescente ignorando maiúsculas",
	MenuSortIntAsc: "Ordenar como inteiros (crescente)", MenuSortIntDesc: "Ordenar como inteiros (decrescente)",
	MenuSortDecAsc: "Ordenar como decimais (crescente)", MenuSortDecDesc: "Ordenar como decimais (decrescente)",
	MenuSortLenAsc: "Ordenar por comprimento (crescente)", MenuSortLenDesc: "Ordenar por comprimento (decrescente)",
	MenuComment: "Comentários", MenuToggleComment: "Alternar comentário de linha", MenuLineComment: "Comentar linhas",
	MenuLineUncomment: "Descomentar linhas", MenuBlockComment: "Comentário de bloco", MenuBlockUncomment: "Remover comentário de bloco",
	MenuAutoCompletion: "Autocompletar", MenuWordCompletion: "Completar palavra", MenuEOLConversion: "Conversão de fim de linha",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Operações com espaços", MenuTrimTrailing: "Remover espaços no final", MenuTrimLeading: "Remover espaços no início",
	MenuTrimBoth: "Remover espaços no início e no final", MenuEOLToSpace: "Fim de linha para espaço", MenuRemoveBlankEOL: "Remover espaços e fins de linha desnecessários",
	MenuTabToSpace: "Tabulação para espaços", MenuSpaceToTabAll: "Espaços para tabulação (todos)", MenuSpaceToTabLeading: "Espaços para tabulação (no início)",
	MenuColumnEditor: "Editor de colunas...", MenuReadOnly: "Somente leitura",
	DateShortLayout: "15:04 02/01/2006", DateLongLayout: "15:04:05 02/01/2006 (Monday)",

	MenuFind: "Localizar...", MenuFindInFiles: "Localizar em arquivos...", MenuFindNext: "Localizar próximo", MenuFindPrev: "Localizar anterior",
	MenuSelectFindNext: "Selecionar e localizar próximo", MenuSelectFindPrev: "Selecionar e localizar anterior", MenuReplace: "Substituir...",
	MenuIncremental: "Pesquisa incremental", MenuMark: "Marcar...", MenuSearchResults: "Janela de resultados",
	MenuNextResult: "Próximo resultado", MenuPrevResult: "Resultado anterior", MenuGoTo: "Ir para...",
	MenuGotoBrace: "Ir para a chave correspondente", MenuSelectBrace: "Selecionar tudo entre chaves",
	MenuStyleAll: "Destacar todas as ocorrências", MenuStyleOne: "Destacar uma ocorrência", MenuClearStyleSub: "Remover destaque",
	MenuClearAllStyles: "Remover todos os destaques", MenuJumpUp: "Saltar para cima", MenuJumpDown: "Saltar para baixo", MenuBookmark: "Marcadores",
	MenuToggleBookmark: "Alternar marcador", MenuNextBookmark: "Próximo marcador", MenuPrevBookmark: "Marcador anterior",
	MenuClearBookmarks: "Remover todos os marcadores", MenuCutBookmarked: "Recortar linhas marcadas", MenuCopyBookmarked: "Copiar linhas marcadas",
	MenuRemoveBookmarked: "Remover linhas marcadas", MenuRemoveUnbookmarked: "Remover linhas não marcadas", MenuInverseBookmarks: "Inverter marcadores",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Com o estilo %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Remover o estilo %d", n) },

	MenuAlwaysOnTop: "Sempre visível", MenuFullScreen: "Tela cheia", MenuShowSymbol: "Mostrar símbolos",
	MenuShowSpaces: "Espaços e tabulações", MenuShowEOL: "Fins de linha", MenuShowAllChars: "Todos os caracteres",
	MenuShowIndentGuides: "Guias de recuo", MenuShowWrapSymbol: "Símbolo de quebra", MenuZoom: "Zoom", MenuZoomIn: "Aumentar",
	MenuZoomOut: "Diminuir", MenuZoomReset: "Zoom padrão", MenuMoveClone: "Mover/clonar o documento",
	MenuMoveToOtherView: "Mover para a outra visão", MenuCloneToOtherView: "Clonar na outra visão", MenuTab: "Abas",
	MenuNextTab: "Próxima aba", MenuPrevTab: "Aba anterior", MenuMoveTabForward: "Mover a aba para a direita", MenuMoveTabBackward: "Mover a aba para a esquerda",
	MenuWordWrap: "Quebra de linha", MenuLineNumbers: "Números de linha", MenuFocusOtherView: "Ir para a outra visão", MenuFoldAll: "Recolher tudo",
	MenuUnfoldAll: "Expandir tudo", MenuFoldCurrent: "Recolher o nível atual", MenuUnfoldCurrent: "Expandir o nível atual",
	MenuFoldLevel: "Recolher nível", MenuUnfoldLevel: "Expandir nível", MenuSummary: "Resumo...", MenuFolderWorkspace: "Pasta como área de trabalho",
	MenuMonitoring: "Monitoramento (tail -f)", MenuToolbar: "Barra de ferramentas", MenuStatusBar: "Barra de status",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Última aba"
		}
		return fmt.Sprintf("Aba %d", n)
	},

	MenuEncANSI: "Codificar em ANSI", MenuEncUTF8: "Codificar em UTF-8", MenuEncUTF8BOM: "Codificar em UTF-8-BOM", MenuEncUTF16BE: "Codificar em UTF-16 BE BOM",
	MenuEncUTF16LE: "Codificar em UTF-16 LE BOM", MenuCharacterSets: "Conjuntos de caracteres", MenuConvANSI: "Converter para ANSI", MenuConvUTF8: "Converter para UTF-8",
	MenuConvUTF8BOM: "Converter para UTF-8-BOM", MenuConvUTF16BE: "Converter para UTF-16 BE BOM", MenuConvUTF16LE: "Converter para UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Árabe", "Baltic": "Báltico", "Celtic": "Celta", "Cyrillic": "Cirílico",
		"Central European": "Europa central", "Chinese": "Chinês", "Greek": "Grego", "Hebrew": "Hebraico", "Japanese": "Japonês",
		"Korean": "Coreano", "North European": "Europa do norte", "Thai": "Tailandês", "Turkish": "Turco", "Vietnamese": "Vietnamita",
		"Western European": "Europa ocidental"},

	MenuPreferences: "Preferências...", MenuColorScheme: "Esquema de cores", MenuShortcuts: "Atalhos de teclado",
	MenuHashGenerate:  func(name string) string { return "Gerar " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Gerar " + name + " de arquivos..." },
	MenuHashSelection: func(name string) string { return "Copiar " + name + " da seleção" },
	MenuBase64Encode:  "Codificar em Base64", MenuBase64Decode: "Decodificar Base64", MenuURLEncode: "Codificar URL", MenuURLDecode: "Decodificar URL",
	MenuJSONFormat: "Formatar JSON", MenuJSONMinify: "Compactar JSON",
	MenuStartRecording: "Iniciar/parar a gravação", MenuStopRecording: "Parar a gravação", MenuPlayback: "Reproduzir",
	MenuSaveMacro: "Salvar a macro gravada...", MenuRunMacroMulti: "Executar a macro várias vezes...", MenuTrimSave: "Remover espaços no final e salvar",
	MenuRun: "Executar...", MenuOpenInBrowser: "Abrir no navegador", MenuSearchInternet: "Pesquisar na Internet", MenuOpenSelectedFile: "Abrir arquivo (nome selecionado)",
	MenuWindows: "Janelas...", MenuSortTabsByName: "Ordenar abas por nome", MenuSortTabsByPath: "Ordenar abas por caminho",
	MenuHelp: "Ajuda online", MenuHomePage: "Página inicial", MenuAbout: "Sobre o AltNotepad",

	OpenTitle: "Abrir", SaveTitle: "Salvar", SaveAsTitle: "Salvar como", SaveCopyTitle: "Salvar uma cópia como", AllFiles: "Todos os arquivos",
	TextFiles: "Arquivos de texto", Save: "Salvar", DontSave: "Não salvar", PathLabel: "Caminho do arquivo:",
	SaveChangesAsk: func(path string) string { return "Salvar as alterações em «" + path + "»?" },
	CreateFileAsk:  func(path string) string { return "«" + path + "» não existe. Criá-lo?" },
	AlreadyOpen:    func(path string) string { return "«" + path + "» está aberto em outra aba" },
	FileExists:     func(path string) string { return "«" + path + "» já existe" },
	FileNotFound:   func(path string) string { return "Arquivo não encontrado: " + path },
	Unencodable: func(enc string) string {
		return "Alguns caracteres não podem ser salvos em " + enc + ". Salvá-los como «?»?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Carregando %s: %d%%", name, percent) },
	LargeFileMode: "Arquivo grande: destaque de sintaxe, recolhimento e quebra de linha desativados",
	LargeFileMark: "(arquivo grande)",
	ReloadTitle:   "Recarregar", RenameTitle: "Renomear", NewName: "Novo nome:", DeleteFileTitle: "Excluir do disco",
	ReloadAsk: func(path string) string {
		return "Recarregar «" + path + "»? As alterações não salvas serão perdidas."
	},
	DeleteFileAsk:    func(path string) string { return "Excluir «" + path + "» do disco?" },
	FileChangedTitle: "Arquivo alterado", FileDeletedTitle: "Arquivo excluído",
	FileChangedAsk: func(path string) string {
		return "«" + path + "»\n\nEste arquivo foi alterado por outro programa.\nDeseja recarregá-lo?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "«" + path + "»\n\nEste arquivo foi alterado por outro programa.\nRecarregá-lo e perder as alterações do editor?"
	},
	FileDeletedAsk: func(path string) string {
		return "«" + path + "»\n\nEste arquivo não existe mais.\nMantê-lo no editor?"
	},
	NotASession: func(name string) string { return name + " não é um arquivo de sessão" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("tamanho: %d   linhas: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Lin: %d   Col: %d   Pos: %d", line, col, pos) },
	StatusSel:  "Sel:",

	GotoTitle: "Ir para...", GotoLine: "Linha", GotoOffset: "Posição", Go: "Ir",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Você está aqui: %d   Ir para (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Você está aqui: %d   Ir para (0 - %d):", cur, total)
	},
	ColumnTitle: "Editor de colunas e seleção múltipla", ColumnText: "Texto a inserir", ColumnNumbers: "Número a inserir",
	ColumnInitial: "Número inicial:", ColumnIncrease: "Incrementar em:", ColumnRepeat: "Repetir:", ColumnFormat: "Formato:",
	ColumnLeading: "Preencher com:", LeadingNone: "Nada", LeadingZeros: "Zeros", LeadingSpaces: "Espaços",
	FullPath: "Caminho completo", Modified: "Modificado", SummaryChars: "Caracteres (sem fins de linha)", SummaryWords: "Palavras",
	SummaryLines: "Linhas", SummaryNonBlankLines: "Linhas não vazias", SummaryLength: "Tamanho do documento", SummarySelected: "Caracteres selecionados",
	Bytes: "bytes", RunTitle: "Executar", RunLabel: "Programa a executar:", Run: "Executar",
	Name: "Nome", State: "Estado", ModifiedState: "modificado", Activate: "Ativar", CloseWindows: "Fechar",
	Command: "Comando", Shortcut: "Atalho", HashInput: "Texto:", HashEachLine: "Tratar cada linha separadamente",
	CopyToClipboard:   "Copiar",
	CopiedToClipboard: func(s string) string { return "Copiado para a área de transferência: " + s },
	MacroName:         "Nome da macro:", MacroTimes: "Quantas vezes executar (* - até o fim do arquivo):",

	FindTab: "Localizar", ReplaceTab: "Substituir", FindInFilesTab: "Localizar em arquivos", MarkTab: "Marcar", FindWhat: "Localizar:",
	ReplaceWith: "Substituir por:", Filters: "Filtros:", Directory: "Pasta:", CurrentFolder: "Pasta atual",
	MatchCase: "Diferenciar maiúsculas", WholeWord: "Somente palavras inteiras", WrapAround: "Voltar ao início", Backward: "Para trás",
	InSelection: "Na seleção", InSubfolders: "Em todas as subpastas", InHidden: "Em pastas ocultas", BookmarkLine: "Marcar a linha",
	PurgeEach: "Limpar a cada pesquisa", SearchMode: "Modo de pesquisa:", ModeNormal: "Normal", ModeExtended: "Estendido (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Expressão regular", DotAll: ". inclui quebras de linha",
	FindNext: "Localizar próximo", FindPrev: "Localizar anterior", Count: "Contar", FindAllCurrent: "Localizar tudo no documento atual",
	FindAllOpen: "Localizar tudo nos documentos abertos", Replace: "Substituir", ReplaceAll: "Substituir tudo",
	ReplaceAllOpen: "Substituir tudo nos documentos abertos", FindAll: "Localizar tudo", ReplaceInFiles: "Substituir em arquivos",
	MarkAll: "Marcar tudo", ClearMarks: "Remover marcas", CopyMarked: "Copiar o texto marcado", Wrapped: "Fim alcançado, continuado do início",
	ReadOnlyDoc: "O documento é somente leitura",
	NotFound:    func(text string) string { return "Não foi possível encontrar «" + text + "»" },
	CountResult: func(n int) string { return fmt.Sprintf("Total: %d %s", n, ptP(n, "ocorrência", "ocorrências")) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, ptP(n, "ocorrência substituída", "ocorrências substituídas"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, ptP(n, "ocorrência marcada", "ocorrências marcadas"))
	},
	FoundCount: func(n int) string { return fmt.Sprintf("%d %s", n, ptP(n, "resultado", "resultados")) },
	HitsCount:  func(n int) string { return fmt.Sprintf("(%d %s)", n, ptP(n, "resultado", "resultados")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s em %d %s", n, ptP(n, "resultado", "resultados"), files, ptP(files, "arquivo", "arquivos"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s em %d %s", n, ptP(n, "substituição", "substituições"), files, ptP(files, "arquivo", "arquivos"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Substituir todos os «" + find + "» por «" + repl + "» em todos os arquivos de\n" + dir + "?"
	},
	NoSuchFolder:  func(path string) string { return "Pasta inexistente: " + path },
	Searching:     func(path string) string { return "Pesquisando: " + path },
	SearchResults: "Resultados da pesquisa", Stop: "Parar", Stopped: "(parada)", Clear: "Limpar", Line: "Linha",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Pesquisa «%s» (%d %s em %d %s de %d pesquisados)", text, hits, ptP(hits, "resultado", "resultados"), files, ptP(files, "arquivo", "arquivos"), searched)
	},

	Workspace: "Área de trabalho", Refresh: "Atualizar", CollapseAll: "Recolher tudo",

	PageGeneral: "Geral", PageEditing: "Edição", PageDisplay: "Exibição", PageFiles: "Novo documento e arquivos",
	Language: "Idioma:", LanguageSystem: "Como no sistema", Theme: "Tema:", ThemeDark: "Escuro", ThemeLight: "Claro",
	ColorScheme: "Esquema de cores:", SchemeByTheme: "Conforme o tema",
	ShowToolbar: "Mostrar a barra de ferramentas", ShowStatusBar: "Mostrar a barra de status", TabCloseButtons: "Botões de fechar nas abas",
	DoubleClickCloses: "Clique duplo fecha a aba", AlwaysOnTop: "Sempre visível", SingleInstance: "Abrir os arquivos na janela já aberta",
	RememberSession: "Lembrar os arquivos abertos para a próxima sessão", BackupSession: "Manter as alterações não salvas entre sessões",
	BackupEvery: "Backup a cada, segundos:", MaxRecent: "Arquivos recentes na lista:",
	TabSize: "Tamanho da tabulação:", InsertSpaces: "Substituir tabulações por espaços", AutoIndent: "Recuo automático", AutoClose: "Fechar parênteses e aspas",
	AutoCompletion: "Completar palavras ao digitar", SmartHome: "Home vai para o primeiro caractere não vazio", MultiEdit: "Edição múltipla (Ctrl+clique)",
	ScrollPast: "Rolar além da última linha", CopyLineNoSel: "Copiar/recortar a linha sem seleção", CaretWidth: "Largura do cursor:",
	CaretBlink: "Piscar do cursor, ms (0 - sem piscar):", EdgeColumn: "Marca de linhas longas, coluna (0 - nenhuma):", WordChars: "Caracteres de palavra adicionais:",
	FontSize: "Tamanho da fonte:", FontFile: "Arquivo de fonte:", BuiltInFont: "JetBrains Mono (embutida)", Fonts: "Fontes", LineNumbers: "Números de linha",
	BookmarkMargin: "Margem de marcadores", FoldMargin: "Margem de recolhimento", CurrentLine: "Destacar a linha atual", SmartHighlight: "Destaque inteligente",
	SmartMatchCase: "Destaque inteligente: diferenciar maiúsculas", SmartWholeWord: "Destaque inteligente: somente palavras inteiras", BraceMatch: "Destacar chaves correspondentes",
	WrapSymbol: "Mostrar o símbolo de quebra", ChangeHistory: "Histórico de alterações na margem",
	NewDocEOL: "Fim de linha dos novos documentos:", NewDocEncoding: "Codificação dos novos documentos:", ANSICharset: "Conjunto de caracteres ANSI:",
	NewDocLanguage: "Linguagem dos novos documentos:", SystemDefault: "Como no sistema", ByLanguage: "Conforme o idioma da interface",
	LargeFileMB: "Limite de arquivos grandes, MB:", CheckFileChanges: "Detectar arquivos alterados por outros programas", AutoReload: "Recarregá-los sem perguntar se não modificados",

	AboutTitle:       func(name string) string { return "Sobre o " + name },
	AboutDescription: "Um editor de texto e código-fonte", Version: "Versão", Author: "Autor:", License: "Licença:",
	VisitWebsite: "Visitar o site", Close: "Fechar",
}
