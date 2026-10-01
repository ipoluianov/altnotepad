package forms

import "fmt"

var zh = Strings{
	Error:      "错误",
	NewDocName: "新建",
	NormalText: "普通文本",

	MenuFile: "文件", MenuEdit: "编辑", MenuSearch: "搜索", MenuView: "视图", MenuEncoding: "编码", MenuLanguage: "语言",
	MenuSettings: "设置", MenuTools: "工具", MenuMacro: "宏", MenuRunMenu: "运行", MenuWindow: "窗口",

	MenuNew: "新建", MenuOpen: "打开...", MenuOpenContainingFolderSub: "打开所在文件夹",
	MenuOpenContainingFolder: "文件管理器", MenuOpenDefaultViewer: "用默认程序打开", MenuOpenFolderWorkspace: "将文件夹作为工作区打开...",
	MenuReload: "从磁盘重新加载", MenuSave: "保存", MenuSaveAs: "另存为...", MenuSaveCopy: "保存副本为...", MenuSaveAll: "全部保存",
	MenuRename: "重命名...", MenuClose: "关闭", MenuCloseAll: "全部关闭", MenuCloseMore: "关闭多个文档",
	MenuCloseAllButActive: "关闭当前以外的所有文档", MenuCloseLeft: "关闭左侧所有文档", MenuCloseRight: "关闭右侧所有文档",
	MenuCloseUnchanged: "关闭所有未修改的文档", MenuDeleteFile: "从磁盘删除...", MenuLoadSession: "加载会话...",
	MenuSaveSession: "保存会话...", MenuRecentFiles: "最近的文件", MenuClearRecent: "清空最近文件列表",
	MenuRestoreClosed: "恢复最近关闭的文件", MenuExit: "退出",
	MenuPrint: "打印...", MenuExportHTML: "导出为 HTML...", MenuFunctionList: "函数列表", Filter: "筛选", MenuDocumentMap: "文档地图",
	MenuChangeHistory: "更改历史", MenuNextChange: "转到下一处更改", MenuPrevChange: "转到上一处更改", MenuClearChanges: "清除更改历史",
	SavedTo:         func(path string) string { return "已保存到 " + path },
	MenuMultiSelect: "多重选择", MenuMultiAll: "选择所有匹配项", MenuMultiAllCase: "选择所有匹配项（区分大小写）",
	MenuMultiNext: "添加下一个匹配项", MenuMultiUndo: "撤销最近添加的选择", MenuMultiSkip: "跳过当前并添加下一个",

	MenuUndo: "撤销", MenuRedo: "重做", MenuCut: "剪切", MenuCopy: "复制", MenuPaste: "粘贴", MenuDelete: "删除", MenuSelectAll: "全选",
	MenuInsert: "插入", MenuDateShort: "日期和时间（短）", MenuDateLong: "日期和时间（长）", MenuDateISO: "日期和时间（ISO 8601）",
	MenuCopyToClipboard: "复制到剪贴板", MenuCopyFullPath: "文件完整路径", MenuCopyFileName: "文件名",
	MenuCopyDirPath: "文件夹路径", MenuIndentSub: "缩进", MenuIndent: "增加缩进", MenuUnindent: "减少缩进",
	MenuConvertCase: "转换大小写", MenuUpperCase: "大写", MenuLowerCase: "小写", MenuProperCase: "单词首字母大写",
	MenuProperCaseBlend: "单词首字母大写（混合）", MenuSentenceCase: "句首字母大写", MenuSentenceCaseBlend: "句首字母大写（混合）",
	MenuInvertCase: "大小写反转", MenuRandomCase: "随机大小写", MenuLineOperations: "行操作",
	MenuDuplicateLine: "复制当前行", MenuRemoveDupLines: "删除重复行", MenuRemoveConsecutiveDupLines: "删除连续的重复行",
	MenuSplitLines: "拆分行", MenuJoinLines: "合并行", MenuMoveLineUp: "上移当前行", MenuMoveLineDown: "下移当前行",
	MenuDeleteLine: "删除当前行", MenuCutLine: "剪切当前行", MenuTransposeLine: "与上一行交换",
	MenuRemoveEmptyLines: "删除空行", MenuRemoveBlankLines: "删除空行（包括仅含空白的行）",
	MenuInsertLineAbove: "在上方插入空行", MenuInsertLineBelow: "在下方插入空行",
	MenuReverseLines: "反转行顺序", MenuShuffleLines: "随机排列行",
	MenuSortAsc: "按升序排列行", MenuSortDesc: "按降序排列行",
	MenuSortAscCI: "升序排列（忽略大小写）", MenuSortDescCI: "降序排列（忽略大小写）",
	MenuSortIntAsc: "按整数升序排列", MenuSortIntDesc: "按整数降序排列",
	MenuSortDecAsc: "按小数升序排列", MenuSortDecDesc: "按小数降序排列",
	MenuSortLenAsc: "按长度升序排列", MenuSortLenDesc: "按长度降序排列",
	MenuComment: "注释", MenuToggleComment: "切换行注释", MenuLineComment: "添加行注释",
	MenuLineUncomment: "取消行注释", MenuBlockComment: "块注释", MenuBlockUncomment: "取消块注释",
	MenuAutoCompletion: "自动完成", MenuWordCompletion: "单词补全", MenuEOLConversion: "换行符转换",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "空白字符操作", MenuTrimTrailing: "去除行尾空白", MenuTrimLeading: "去除行首空白",
	MenuTrimBoth: "去除行首和行尾空白", MenuEOLToSpace: "换行符转为空格", MenuRemoveBlankEOL: "去除多余的空白和换行符",
	MenuTabToSpace: "制表符转为空格", MenuSpaceToTabAll: "空格转为制表符（全部）", MenuSpaceToTabLeading: "空格转为制表符（行首）",
	MenuColumnEditor: "列编辑器...", MenuReadOnly: "只读",
	DateShortLayout: "15:04 2006/01/02", DateLongLayout: "2006年1月2日 15:04:05 (Monday)",

	MenuFind: "查找...", MenuFindInFiles: "在文件中查找...", MenuFindNext: "查找下一个", MenuFindPrev: "查找上一个",
	MenuSelectFindNext: "选择并查找下一个", MenuSelectFindPrev: "选择并查找上一个", MenuReplace: "替换...",
	MenuIncremental: "增量查找", MenuMark: "标记...", MenuSearchResults: "搜索结果窗口",
	MenuNextResult: "下一个结果", MenuPrevResult: "上一个结果", MenuGoTo: "转到...",
	MenuGotoBrace: "转到匹配的括号", MenuSelectBrace: "选择括号之间的内容",
	MenuStyleAll: "高亮所有匹配项", MenuStyleOne: "高亮一个匹配项", MenuClearStyleSub: "清除高亮",
	MenuClearAllStyles: "清除所有高亮", MenuJumpUp: "向上跳转", MenuJumpDown: "向下跳转", MenuBookmark: "书签",
	MenuToggleBookmark: "切换书签", MenuNextBookmark: "下一个书签", MenuPrevBookmark: "上一个书签",
	MenuClearBookmarks: "清除所有书签", MenuCutBookmarked: "剪切书签行", MenuCopyBookmarked: "复制书签行",
	MenuRemoveBookmarked: "删除书签行", MenuRemoveUnbookmarked: "删除非书签行", MenuInverseBookmarks: "反转书签",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("使用样式 %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("清除样式 %d", n) },

	MenuAlwaysOnTop: "总在最前", MenuFullScreen: "全屏", MenuShowSymbol: "显示符号",
	MenuShowSpaces: "空格和制表符", MenuShowEOL: "换行符", MenuShowAllChars: "所有字符",
	MenuShowIndentGuides: "缩进参考线", MenuShowWrapSymbol: "自动换行符号", MenuZoom: "缩放", MenuZoomIn: "放大",
	MenuZoomOut: "缩小", MenuZoomReset: "恢复默认缩放", MenuMoveClone: "移动/克隆文档",
	MenuMoveToOtherView: "移动到另一视图", MenuCloneToOtherView: "克隆到另一视图", MenuTab: "标签页",
	MenuNextTab: "下一个标签页", MenuPrevTab: "上一个标签页", MenuMoveTabForward: "标签页右移", MenuMoveTabBackward: "标签页左移",
	MenuWordWrap: "自动换行", MenuLineNumbers: "行号", MenuFocusOtherView: "切换到另一视图", MenuFoldAll: "全部折叠",
	MenuUnfoldAll: "全部展开", MenuFoldCurrent: "折叠当前层级", MenuUnfoldCurrent: "展开当前层级",
	MenuFoldLevel: "折叠层级", MenuUnfoldLevel: "展开层级", MenuSummary: "摘要...", MenuFolderWorkspace: "文件夹工作区",
	MenuMonitoring: "监视 (tail -f)", MenuToolbar: "工具栏", MenuStatusBar: "状态栏",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "最后一个标签页"
		}
		return fmt.Sprintf("标签页 %d", n)
	},

	MenuEncANSI: "以 ANSI 编码", MenuEncUTF8: "以 UTF-8 编码", MenuEncUTF8BOM: "以 UTF-8-BOM 编码", MenuEncUTF16BE: "以 UTF-16 BE BOM 编码",
	MenuEncUTF16LE: "以 UTF-16 LE BOM 编码", MenuCharacterSets: "字符集", MenuConvANSI: "转为 ANSI", MenuConvUTF8: "转为 UTF-8",
	MenuConvUTF8BOM: "转为 UTF-8-BOM", MenuConvUTF16BE: "转为 UTF-16 BE BOM", MenuConvUTF16LE: "转为 UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "阿拉伯语", "Baltic": "波罗的语", "Celtic": "凯尔特语", "Cyrillic": "西里尔文",
		"Central European": "中欧", "Chinese": "中文", "Greek": "希腊语", "Hebrew": "希伯来语", "Japanese": "日语",
		"Korean": "韩语", "North European": "北欧", "Thai": "泰语", "Turkish": "土耳其语", "Vietnamese": "越南语",
		"Western European": "西欧"},

	MenuPreferences: "首选项...", MenuColorScheme: "配色方案", MenuShortcuts: "快捷键",
	MenuHashGenerate:  func(name string) string { return "生成 " + name + "..." },
	MenuHashFiles:     func(name string) string { return "生成文件的 " + name + "..." },
	MenuHashSelection: func(name string) string { return "复制选区的 " + name },
	MenuBase64Encode:  "Base64 编码", MenuBase64Decode: "Base64 解码", MenuURLEncode: "URL 编码", MenuURLDecode: "URL 解码",
	MenuJSONFormat: "格式化 JSON", MenuJSONMinify: "压缩 JSON",
	MenuStartRecording: "开始/停止录制", MenuStopRecording: "停止录制", MenuPlayback: "回放",
	MenuSaveMacro: "保存录制的宏...", MenuRunMacroMulti: "多次运行宏...", MenuTrimSave: "去除行尾空白并保存",
	MenuRun: "运行...", MenuOpenInBrowser: "在浏览器中打开", MenuSearchInternet: "在互联网上搜索", MenuOpenSelectedFile: "打开文件（所选文件名）",
	MenuWindows: "窗口...", MenuSortTabsByName: "按名称排列标签页", MenuSortTabsByPath: "按路径排列标签页",
	MenuHelp: "在线帮助", MenuHomePage: "主页", MenuAbout: "关于 AltNotepad",

	OpenTitle: "打开", SaveTitle: "保存", SaveAsTitle: "另存为", SaveCopyTitle: "保存副本为", AllFiles: "所有文件",
	TextFiles: "文本文件", Save: "保存", DontSave: "不保存", PathLabel: "文件路径：",
	SaveChangesAsk: func(path string) string { return "是否保存对“" + path + "”的更改？" },
	CreateFileAsk:  func(path string) string { return "“" + path + "”不存在。是否创建？" },
	AlreadyOpen:    func(path string) string { return "“" + path + "”已在另一个标签页中打开" },
	FileExists:     func(path string) string { return "“" + path + "”已存在" },
	FileNotFound:   func(path string) string { return "找不到文件：" + path },
	Unencodable: func(enc string) string {
		return "部分字符无法以 " + enc + " 保存。是否将其保存为“?”？"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("正在加载 %s：%d%%", name, percent) },
	LargeFileMode: "大文件：已禁用语法高亮、折叠和自动换行",
	LargeFileMark: "（大文件）",
	ReloadTitle:   "重新加载", RenameTitle: "重命名", NewName: "新名称：", DeleteFileTitle: "从磁盘删除",
	ReloadAsk: func(path string) string {
		return "是否重新加载“" + path + "”？未保存的更改将丢失。"
	},
	DeleteFileAsk:    func(path string) string { return "是否从磁盘删除“" + path + "”？" },
	FileChangedTitle: "文件已更改", FileDeletedTitle: "文件已删除",
	FileChangedAsk: func(path string) string {
		return "“" + path + "”\n\n此文件已被其他程序修改。\n是否重新加载？"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "“" + path + "”\n\n此文件已被其他程序修改。\n是否重新加载并丢失编辑器中的更改？"
	},
	FileDeletedAsk: func(path string) string {
		return "“" + path + "”\n\n此文件已不存在。\n是否保留在编辑器中？"
	},
	NotASession: func(name string) string { return name + " 不是会话文件" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("长度：%d   行数：%d", length, lines) },
	StatusPos: func(line, col, pos int) string {
		return fmt.Sprintf("行：%d   列：%d   位置：%d", line, col, pos)
	},
	StatusSel: "选中：",

	GotoTitle: "转到...", GotoLine: "行", GotoOffset: "偏移", Go: "转到",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("当前：%d   转到 (1 - %d)：", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("当前：%d   转到 (0 - %d)：", cur, total)
	},
	ColumnTitle: "列编辑器与多重选择", ColumnText: "插入文本", ColumnNumbers: "插入数字",
	ColumnInitial: "起始数字：", ColumnIncrease: "增量：", ColumnRepeat: "重复：", ColumnFormat: "格式：",
	ColumnLeading: "前导填充：", LeadingNone: "无", LeadingZeros: "零", LeadingSpaces: "空格",
	FullPath: "完整路径", Modified: "修改时间", SummaryChars: "字符数（不含换行符）", SummaryWords: "单词数",
	SummaryLines: "行数", SummaryNonBlankLines: "非空行数", SummaryLength: "文档长度", SummarySelected: "选中字符数",
	Bytes: "字节", RunTitle: "运行", RunLabel: "要运行的程序：", Run: "运行",
	Name: "名称", State: "状态", ModifiedState: "已修改", Activate: "激活", CloseWindows: "关闭",
	Command: "命令", Shortcut: "快捷键", HashInput: "文本：", HashEachLine: "逐行分别处理",
	CopyToClipboard:   "复制",
	CopiedToClipboard: func(s string) string { return "已复制到剪贴板：" + s },
	MacroName:         "宏名称：", MacroTimes: "运行次数（* - 直到文件末尾）：",

	FindTab: "查找", ReplaceTab: "替换", FindInFilesTab: "在文件中查找", MarkTab: "标记", FindWhat: "查找：",
	ReplaceWith: "替换为：", Filters: "筛选：", Directory: "文件夹：", CurrentFolder: "当前文件夹",
	MatchCase: "区分大小写", WholeWord: "全词匹配", WrapAround: "循环查找", Backward: "向后查找",
	InSelection: "在选区中", InSubfolders: "包含所有子文件夹", InHidden: "包含隐藏文件夹", BookmarkLine: "标记书签",
	PurgeEach: "每次查找前清除", SearchMode: "查找模式：", ModeNormal: "普通", ModeExtended: "扩展 (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "正则表达式", DotAll: ". 匹配换行符",
	FindNext: "查找下一个", FindPrev: "查找上一个", Count: "计数", FindAllCurrent: "在当前文档中查找全部",
	FindAllOpen: "在所有打开的文档中查找全部", Replace: "替换", ReplaceAll: "全部替换",
	ReplaceAllOpen: "在所有打开的文档中全部替换", FindAll: "查找全部", ReplaceInFiles: "在文件中替换",
	MarkAll: "全部标记", ClearMarks: "清除标记", CopyMarked: "复制标记的文本", Wrapped: "已到达末尾，从开头继续",
	ReadOnlyDoc: "文档为只读",
	NotFound:    func(text string) string { return "找不到“" + text + "”" },
	CountResult: func(n int) string { return fmt.Sprintf("共 %d 个匹配项", n) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("已替换 %d 处", n)
	},
	MarkedCount: func(n int) string { return fmt.Sprintf("已标记 %d 个匹配项", n) },
	FoundCount:  func(n int) string { return fmt.Sprintf("%d 个结果", n) },
	HitsCount:   func(n int) string { return fmt.Sprintf("（%d 个结果）", n) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("在 %d 个文件中找到 %d 个结果", files, n)
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("在 %d 个文件中替换了 %d 处", files, n)
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "是否将以下文件夹中所有文件里的“" + find + "”替换为“" + repl + "”？\n" + dir
	},
	NoSuchFolder:  func(path string) string { return "文件夹不存在：" + path },
	Searching:     func(path string) string { return "正在搜索：" + path },
	SearchResults: "搜索结果", Stop: "停止", Stopped: "（已停止）", Clear: "清除", Line: "行",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("搜索“%s”（在已搜索的 %d 个文件中的 %d 个文件里找到 %d 个结果）", text, searched, files, hits)
	},

	Workspace: "工作区", Refresh: "刷新", CollapseAll: "全部折叠",

	PageGeneral: "常规", PageEditing: "编辑", PageDisplay: "显示", PageFiles: "新建文档和文件",
	Language: "语言：", LanguageSystem: "跟随系统", Theme: "主题：", ThemeDark: "深色", ThemeLight: "浅色",
	ColorScheme: "配色方案：", SchemeByTheme: "跟随主题",
	ShowToolbar: "显示工具栏", ShowStatusBar: "显示状态栏", TabCloseButtons: "标签页上显示关闭按钮",
	DoubleClickCloses: "双击关闭标签页", AlwaysOnTop: "总在最前", SingleInstance: "在已打开的窗口中打开文件",
	RememberSession: "记住打开的文件以供下次使用", BackupSession: "在会话之间保留未保存的更改",
	BackupEvery: "备份间隔（秒）：", MaxRecent: "最近文件数量：",
	TabSize: "制表符宽度：", InsertSpaces: "用空格替换制表符", AutoIndent: "自动缩进", AutoClose: "自动补全括号和引号",
	AutoCompletion: "输入时补全单词", SmartHome: "Home 键跳到首个非空白字符", MultiEdit: "多处编辑（Ctrl+单击）",
	ScrollPast: "允许滚动超过最后一行", CopyLineNoSel: "无选区时复制/剪切整行", CaretWidth: "光标宽度：",
	CaretBlink: "光标闪烁间隔，毫秒（0 - 不闪烁）：", EdgeColumn: "长行标记列（0 - 无）：", WordChars: "额外的单词字符：",
	FontSize: "字号：", FontFile: "字体文件：", BuiltInFont: "JetBrains Mono（内置）", Fonts: "字体", LineNumbers: "行号",
	BookmarkMargin: "书签边栏", FoldMargin: "折叠边栏", CurrentLine: "高亮当前行", SmartHighlight: "智能高亮",
	SmartMatchCase: "智能高亮：区分大小写", SmartWholeWord: "智能高亮：全词匹配", BraceMatch: "高亮匹配的括号",
	WrapSymbol: "显示自动换行符号", ChangeHistory: "在边栏显示更改历史",
	NewDocEOL: "新文档的换行符：", NewDocEncoding: "新文档的编码：", ANSICharset: "ANSI 字符集：",
	NewDocLanguage: "新文档的语言：", SystemDefault: "跟随系统", ByLanguage: "根据界面语言",
	LargeFileMB: "大文件阈值（MB）：", CheckFileChanges: "检测被其他程序修改的文件", AutoReload: "未修改时自动重新加载",

	AboutTitle:       func(name string) string { return "关于 " + name },
	AboutDescription: "文本和源代码编辑器", Version: "版本", Author: "作者：", License: "许可证：",
	VisitWebsite: "访问网站", Close: "关闭",
}
