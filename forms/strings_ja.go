package forms

import "fmt"

var ja = Strings{
	Error:      "エラー",
	NewDocName: "新規",
	NormalText: "標準テキスト",

	MenuFile: "ファイル", MenuEdit: "編集", MenuSearch: "検索", MenuView: "表示", MenuEncoding: "エンコード", MenuLanguage: "言語",
	MenuSettings: "設定", MenuTools: "ツール", MenuMacro: "マクロ", MenuRunMenu: "実行", MenuWindow: "ウィンドウ",

	MenuNew: "新規作成", MenuOpen: "開く...", MenuOpenContainingFolderSub: "ファイルのあるフォルダーを開く",
	MenuOpenContainingFolder: "ファイルマネージャー", MenuOpenDefaultViewer: "既定のアプリで開く", MenuOpenFolderWorkspace: "フォルダーをワークスペースとして開く...",
	MenuReload: "ディスクから再読み込み", MenuSave: "保存", MenuSaveAs: "名前を付けて保存...", MenuSaveCopy: "コピーを保存...", MenuSaveAll: "すべて保存",
	MenuRename: "名前の変更...", MenuClose: "閉じる", MenuCloseAll: "すべて閉じる", MenuCloseMore: "複数の文書を閉じる",
	MenuCloseAllButActive: "現在の文書以外をすべて閉じる", MenuCloseLeft: "左側をすべて閉じる", MenuCloseRight: "右側をすべて閉じる",
	MenuCloseUnchanged: "未変更の文書をすべて閉じる", MenuDeleteFile: "ディスクから削除...", MenuLoadSession: "セッションを読み込む...",
	MenuSaveSession: "セッションを保存...", MenuRecentFiles: "最近使ったファイル", MenuClearRecent: "最近使ったファイルの一覧を消去",
	MenuRestoreClosed: "最後に閉じたファイルを復元", MenuExit: "終了",
	MenuPrint: "印刷...", MenuExportHTML: "HTML にエクスポート...", MenuFunctionList: "関数一覧", Filter: "フィルター", MenuDocumentMap: "ドキュメントマップ",
	MenuChangeHistory: "変更履歴", MenuNextChange: "次の変更へ移動", MenuPrevChange: "前の変更へ移動", MenuClearChanges: "変更履歴を消去",
	SavedTo:         func(path string) string { return path + " に保存しました" },
	MenuMultiSelect: "複数選択", MenuMultiAll: "すべての一致を選択", MenuMultiAllCase: "すべての一致を選択（大文字と小文字を区別）",
	MenuMultiNext: "次の一致を追加", MenuMultiUndo: "最後の追加を取り消す", MenuMultiSkip: "現在をスキップして次を追加",

	MenuUndo: "元に戻す", MenuRedo: "やり直し", MenuCut: "切り取り", MenuCopy: "コピー", MenuPaste: "貼り付け", MenuDelete: "削除", MenuSelectAll: "すべて選択",
	MenuInsert: "挿入", MenuDateShort: "日付と時刻（短い形式）", MenuDateLong: "日付と時刻（長い形式）", MenuDateISO: "日付と時刻（ISO 8601）",
	MenuCopyToClipboard: "クリップボードにコピー", MenuCopyFullPath: "ファイルのフルパス", MenuCopyFileName: "ファイル名",
	MenuCopyDirPath: "フォルダーのパス", MenuIndentSub: "インデント", MenuIndent: "インデントを増やす", MenuUnindent: "インデントを減らす",
	MenuConvertCase: "大文字/小文字の変換", MenuUpperCase: "大文字", MenuLowerCase: "小文字", MenuProperCase: "単語の先頭を大文字",
	MenuProperCaseBlend: "単語の先頭を大文字（混在）", MenuSentenceCase: "文の先頭を大文字", MenuSentenceCaseBlend: "文の先頭を大文字（混在）",
	MenuInvertCase: "大文字/小文字を反転", MenuRandomCase: "ランダム", MenuLineOperations: "行の操作",
	MenuDuplicateLine: "現在の行を複製", MenuRemoveDupLines: "重複行を削除", MenuRemoveConsecutiveDupLines: "連続する重複行を削除",
	MenuSplitLines: "行を分割", MenuJoinLines: "行を結合", MenuMoveLineUp: "行を上へ移動", MenuMoveLineDown: "行を下へ移動",
	MenuDeleteLine: "現在の行を削除", MenuCutLine: "現在の行を切り取り", MenuTransposeLine: "前の行と入れ替え",
	MenuRemoveEmptyLines: "空行を削除", MenuRemoveBlankLines: "空行を削除（空白のみの行を含む）",
	MenuInsertLineAbove: "上に空行を挿入", MenuInsertLineBelow: "下に空行を挿入",
	MenuReverseLines: "行の順序を反転", MenuShuffleLines: "行をランダムに並べ替え",
	MenuSortAsc: "行を昇順で並べ替え", MenuSortDesc: "行を降順で並べ替え",
	MenuSortAscCI: "昇順で並べ替え（大文字と小文字を無視）", MenuSortDescCI: "降順で並べ替え（大文字と小文字を無視）",
	MenuSortIntAsc: "整数として昇順で並べ替え", MenuSortIntDesc: "整数として降順で並べ替え",
	MenuSortDecAsc: "小数として昇順で並べ替え", MenuSortDecDesc: "小数として降順で並べ替え",
	MenuSortLenAsc: "長さの昇順で並べ替え", MenuSortLenDesc: "長さの降順で並べ替え",
	MenuComment: "コメント", MenuToggleComment: "行コメントの切り替え", MenuLineComment: "行をコメントにする",
	MenuLineUncomment: "行のコメントを解除", MenuBlockComment: "ブロックコメント", MenuBlockUncomment: "ブロックコメントを解除",
	MenuAutoCompletion: "自動補完", MenuWordCompletion: "単語を補完", MenuEOLConversion: "改行コードの変換",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "空白の操作", MenuTrimTrailing: "行末の空白を削除", MenuTrimLeading: "行頭の空白を削除",
	MenuTrimBoth: "行頭と行末の空白を削除", MenuEOLToSpace: "改行を空白に変換", MenuRemoveBlankEOL: "不要な空白と改行を削除",
	MenuTabToSpace: "タブを空白に変換", MenuSpaceToTabAll: "空白をタブに変換（すべて）", MenuSpaceToTabLeading: "空白をタブに変換（行頭）",
	MenuColumnEditor: "列エディター...", MenuReadOnly: "読み取り専用",
	DateShortLayout: "15:04 2006/01/02", DateLongLayout: "2006年1月2日 15:04:05 (Monday)",

	MenuFind: "検索...", MenuFindInFiles: "ファイル内を検索...", MenuFindNext: "次を検索", MenuFindPrev: "前を検索",
	MenuSelectFindNext: "選択して次を検索", MenuSelectFindPrev: "選択して前を検索", MenuReplace: "置換...",
	MenuIncremental: "インクリメンタル検索", MenuMark: "マーク...", MenuSearchResults: "検索結果ウィンドウ",
	MenuNextResult: "次の結果", MenuPrevResult: "前の結果", MenuGoTo: "移動...",
	MenuGotoBrace: "対応する括弧へ移動", MenuSelectBrace: "括弧の間をすべて選択",
	MenuStyleAll: "すべての一致を強調表示", MenuStyleOne: "一つの一致を強調表示", MenuClearStyleSub: "強調表示を解除",
	MenuClearAllStyles: "すべての強調表示を解除", MenuJumpUp: "上へジャンプ", MenuJumpDown: "下へジャンプ", MenuBookmark: "ブックマーク",
	MenuToggleBookmark: "ブックマークの切り替え", MenuNextBookmark: "次のブックマーク", MenuPrevBookmark: "前のブックマーク",
	MenuClearBookmarks: "すべてのブックマークを解除", MenuCutBookmarked: "ブックマーク行を切り取り", MenuCopyBookmarked: "ブックマーク行をコピー",
	MenuRemoveBookmarked: "ブックマーク行を削除", MenuRemoveUnbookmarked: "ブックマークのない行を削除", MenuInverseBookmarks: "ブックマークを反転",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("スタイル %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("スタイル %d を解除", n) },

	MenuAlwaysOnTop: "常に手前に表示", MenuFullScreen: "全画面表示", MenuShowSymbol: "記号を表示",
	MenuShowSpaces: "空白とタブ", MenuShowEOL: "改行", MenuShowAllChars: "すべての文字",
	MenuShowIndentGuides: "インデントガイド", MenuShowWrapSymbol: "折り返し記号", MenuZoom: "ズーム", MenuZoomIn: "拡大",
	MenuZoomOut: "縮小", MenuZoomReset: "既定のズーム", MenuMoveClone: "文書の移動/複製",
	MenuMoveToOtherView: "もう一方のビューへ移動", MenuCloneToOtherView: "もう一方のビューへ複製", MenuTab: "タブ",
	MenuNextTab: "次のタブ", MenuPrevTab: "前のタブ", MenuMoveTabForward: "タブを右へ移動", MenuMoveTabBackward: "タブを左へ移動",
	MenuWordWrap: "右端で折り返す", MenuLineNumbers: "行番号", MenuFocusOtherView: "もう一方のビューへ切り替え", MenuFoldAll: "すべて折りたたむ",
	MenuUnfoldAll: "すべて展開", MenuFoldCurrent: "現在のレベルを折りたたむ", MenuUnfoldCurrent: "現在のレベルを展開",
	MenuFoldLevel: "レベルを折りたたむ", MenuUnfoldLevel: "レベルを展開", MenuSummary: "概要...", MenuFolderWorkspace: "フォルダーワークスペース",
	MenuMonitoring: "監視 (tail -f)", MenuToolbar: "ツールバー", MenuStatusBar: "ステータスバー",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "最後のタブ"
		}
		return fmt.Sprintf("タブ %d", n)
	},

	MenuEncANSI: "ANSI でエンコード", MenuEncUTF8: "UTF-8 でエンコード", MenuEncUTF8BOM: "UTF-8-BOM でエンコード", MenuEncUTF16BE: "UTF-16 BE BOM でエンコード",
	MenuEncUTF16LE: "UTF-16 LE BOM でエンコード", MenuCharacterSets: "文字セット", MenuConvANSI: "ANSI に変換", MenuConvUTF8: "UTF-8 に変換",
	MenuConvUTF8BOM: "UTF-8-BOM に変換", MenuConvUTF16BE: "UTF-16 BE BOM に変換", MenuConvUTF16LE: "UTF-16 LE BOM に変換",
	CharsetGroups: map[string]string{"Arabic": "アラビア語", "Baltic": "バルト語", "Celtic": "ケルト語", "Cyrillic": "キリル文字",
		"Central European": "中央ヨーロッパ", "Chinese": "中国語", "Greek": "ギリシャ語", "Hebrew": "ヘブライ語", "Japanese": "日本語",
		"Korean": "韓国語", "North European": "北ヨーロッパ", "Thai": "タイ語", "Turkish": "トルコ語", "Vietnamese": "ベトナム語",
		"Western European": "西ヨーロッパ"},

	MenuPreferences: "環境設定...", MenuColorScheme: "配色", MenuShortcuts: "キーボードショートカット",
	MenuHashGenerate:  func(name string) string { return name + " を生成..." },
	MenuHashFiles:     func(name string) string { return "ファイルの " + name + " を生成..." },
	MenuHashSelection: func(name string) string { return "選択範囲の " + name + " をコピー" },
	MenuBase64Encode:  "Base64 エンコード", MenuBase64Decode: "Base64 デコード", MenuURLEncode: "URL エンコード", MenuURLDecode: "URL デコード",
	MenuJSONFormat: "JSON を整形", MenuJSONMinify: "JSON を圧縮",
	MenuStartRecording: "記録の開始/停止", MenuStopRecording: "記録を停止", MenuPlayback: "再生",
	MenuSaveMacro: "記録したマクロを保存...", MenuRunMacroMulti: "マクロを複数回実行...", MenuTrimSave: "行末の空白を削除して保存",
	MenuRun: "実行...", MenuOpenInBrowser: "ブラウザーで開く", MenuSearchInternet: "インターネットで検索", MenuOpenSelectedFile: "ファイルを開く（選択した名前）",
	MenuWindows: "ウィンドウ...", MenuSortTabsByName: "タブを名前順に並べ替え", MenuSortTabsByPath: "タブをパス順に並べ替え",
	MenuHelp: "オンラインヘルプ", MenuHomePage: "ホームページ", MenuAbout: "AltNotepad について",

	OpenTitle: "開く", SaveTitle: "保存", SaveAsTitle: "名前を付けて保存", SaveCopyTitle: "コピーを保存", AllFiles: "すべてのファイル",
	TextFiles: "テキストファイル", Save: "保存", DontSave: "保存しない", PathLabel: "ファイルのパス：",
	SaveChangesAsk: func(path string) string { return "「" + path + "」への変更を保存しますか？" },
	CreateFileAsk:  func(path string) string { return "「" + path + "」は存在しません。作成しますか？" },
	AlreadyOpen:    func(path string) string { return "「" + path + "」は別のタブで開かれています" },
	FileExists:     func(path string) string { return "「" + path + "」は既に存在します" },
	FileNotFound:   func(path string) string { return "ファイルが見つかりません：" + path },
	Unencodable: func(enc string) string {
		return "一部の文字は " + enc + " で保存できません。「?」として保存しますか？"
	},
	Loading: func(name string, percent int) string {
		return fmt.Sprintf("%s を読み込み中：%d%%", name, percent)
	},
	LargeFileMode: "大きなファイル：構文の強調表示、折りたたみ、折り返しは無効です",
	LargeFileMark: "（大きなファイル）",
	ReloadTitle:   "再読み込み", RenameTitle: "名前の変更", NewName: "新しい名前：", DeleteFileTitle: "ディスクから削除",
	ReloadAsk: func(path string) string {
		return "「" + path + "」を再読み込みしますか？保存されていない変更は失われます。"
	},
	DeleteFileAsk:    func(path string) string { return "「" + path + "」をディスクから削除しますか？" },
	FileChangedTitle: "ファイルが変更されました", FileDeletedTitle: "ファイルが削除されました",
	FileChangedAsk: func(path string) string {
		return "「" + path + "」\n\nこのファイルは他のプログラムによって変更されました。\n再読み込みしますか？"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "「" + path + "」\n\nこのファイルは他のプログラムによって変更されました。\n再読み込みしてエディターでの変更を破棄しますか？"
	},
	FileDeletedAsk: func(path string) string {
		return "「" + path + "」\n\nこのファイルはもう存在しません。\nエディターに残しますか？"
	},
	NotASession: func(name string) string { return name + " はセッションファイルではありません" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("長さ：%d   行数：%d", length, lines) },
	StatusPos: func(line, col, pos int) string {
		return fmt.Sprintf("行：%d   列：%d   位置：%d", line, col, pos)
	},
	StatusSel: "選択：",

	GotoTitle: "移動...", GotoLine: "行", GotoOffset: "オフセット", Go: "移動",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("現在：%d   移動先 (1 - %d)：", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("現在：%d   移動先 (0 - %d)：", cur, total)
	},
	ColumnTitle: "列エディターと複数選択", ColumnText: "挿入するテキスト", ColumnNumbers: "挿入する数値",
	ColumnInitial: "開始番号：", ColumnIncrease: "増分：", ColumnRepeat: "繰り返し：", ColumnFormat: "形式：",
	ColumnLeading: "先頭の埋め文字：", LeadingNone: "なし", LeadingZeros: "ゼロ", LeadingSpaces: "空白",
	FullPath: "フルパス", Modified: "更新日時", SummaryChars: "文字数（改行を除く）", SummaryWords: "単語数",
	SummaryLines: "行数", SummaryNonBlankLines: "空でない行数", SummaryLength: "文書の長さ", SummarySelected: "選択された文字数",
	Bytes: "バイト", RunTitle: "実行", RunLabel: "実行するプログラム：", Run: "実行",
	Name: "名前", State: "状態", ModifiedState: "変更あり", Activate: "アクティブにする", CloseWindows: "閉じる",
	Command: "コマンド", Shortcut: "ショートカット", HashInput: "テキスト：", HashEachLine: "各行を個別に処理",
	CopyToClipboard:   "コピー",
	CopiedToClipboard: func(s string) string { return "クリップボードにコピーしました：" + s },
	MacroName:         "マクロ名：", MacroTimes: "実行回数（* - ファイルの末尾まで）：",

	FindTab: "検索", ReplaceTab: "置換", FindInFilesTab: "ファイル内を検索", MarkTab: "マーク", FindWhat: "検索する文字列：",
	ReplaceWith: "置換後の文字列：", Filters: "フィルター：", Directory: "フォルダー：", CurrentFolder: "現在のフォルダー",
	MatchCase: "大文字と小文字を区別", WholeWord: "単語単位", WrapAround: "末尾で先頭に戻る", Backward: "後方へ",
	InSelection: "選択範囲内", InSubfolders: "すべてのサブフォルダー", InHidden: "隠しフォルダー", BookmarkLine: "行をブックマーク",
	PurgeEach: "検索ごとにクリア", SearchMode: "検索モード：", ModeNormal: "標準", ModeExtended: "拡張 (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "正規表現", DotAll: ". が改行にも一致",
	FindNext: "次を検索", FindPrev: "前を検索", Count: "カウント", FindAllCurrent: "現在の文書ですべて検索",
	FindAllOpen: "開いているすべての文書で検索", Replace: "置換", ReplaceAll: "すべて置換",
	ReplaceAllOpen: "開いているすべての文書で置換", FindAll: "すべて検索", ReplaceInFiles: "ファイル内で置換",
	MarkAll: "すべてマーク", ClearMarks: "マークを解除", CopyMarked: "マークしたテキストをコピー", Wrapped: "末尾に達したため先頭から続行しました",
	ReadOnlyDoc: "文書は読み取り専用です",
	NotFound:    func(text string) string { return "「" + text + "」が見つかりません" },
	CountResult: func(n int) string { return fmt.Sprintf("合計：%d 件の一致", n) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d 件を置換しました", n)
	},
	MarkedCount: func(n int) string { return fmt.Sprintf("%d 件の一致をマークしました", n) },
	FoundCount:  func(n int) string { return fmt.Sprintf("%d 件", n) },
	HitsCount:   func(n int) string { return fmt.Sprintf("（%d 件）", n) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d 個のファイルで %d 件", files, n)
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d 個のファイルで %d 件を置換", files, n)
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "次のフォルダー内のすべてのファイルで「" + find + "」を「" + repl + "」に置換しますか？\n" + dir
	},
	NoSuchFolder:  func(path string) string { return "フォルダーが存在しません：" + path },
	Searching:     func(path string) string { return "検索中：" + path },
	SearchResults: "検索結果", Stop: "停止", Stopped: "（停止）", Clear: "クリア", Line: "行",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("「%s」の検索（検索した %d 個のうち %d 個のファイルで %d 件）", text, searched, files, hits)
	},

	Workspace: "ワークスペース", Refresh: "更新", CollapseAll: "すべて折りたたむ",

	PageGeneral: "全般", PageEditing: "編集", PageDisplay: "表示", PageFiles: "新規文書とファイル",
	Language: "言語：", LanguageSystem: "システムと同じ", Theme: "テーマ：", ThemeDark: "ダーク", ThemeLight: "ライト",
	ColorScheme: "配色：", SchemeByTheme: "テーマに従う",
	ShowToolbar: "ツールバーを表示", ShowStatusBar: "ステータスバーを表示", TabCloseButtons: "タブに閉じるボタンを表示",
	DoubleClickCloses: "ダブルクリックでタブを閉じる", AlwaysOnTop: "常に手前に表示", SingleInstance: "開いているウィンドウでファイルを開く",
	RememberSession: "開いているファイルを次回のために記憶", BackupSession: "保存していない変更をセッション間で保持",
	BackupEvery: "バックアップ間隔（秒）：", MaxRecent: "最近使ったファイルの数：",
	TabSize: "タブのサイズ：", InsertSpaces: "タブを空白に置き換える", AutoIndent: "自動インデント", AutoClose: "括弧と引用符を自動で閉じる",
	AutoCompletion: "入力中に単語を補完", SmartHome: "Home キーで最初の非空白文字へ移動", MultiEdit: "複数箇所の編集（Ctrl+クリック）",
	ScrollPast: "最終行より先へスクロール", CopyLineNoSel: "選択なしで行をコピー/切り取り", CaretWidth: "カーソルの幅：",
	CaretBlink: "カーソルの点滅、ミリ秒（0 - 点滅なし）：", EdgeColumn: "長い行の目印の列（0 - なし）：", WordChars: "単語に含める追加の文字：",
	FontSize: "フォントサイズ：", FontFile: "フォントファイル：", BuiltInFont: "JetBrains Mono（内蔵）", Fonts: "フォント", LineNumbers: "行番号",
	BookmarkMargin: "ブックマーク余白", FoldMargin: "折りたたみ余白", CurrentLine: "現在の行を強調表示", SmartHighlight: "スマートハイライト",
	SmartMatchCase: "スマートハイライト：大文字と小文字を区別", SmartWholeWord: "スマートハイライト：単語単位", BraceMatch: "対応する括弧を強調表示",
	WrapSymbol: "折り返し記号を表示", ChangeHistory: "余白に変更履歴を表示",
	NewDocEOL: "新規文書の改行コード：", NewDocEncoding: "新規文書のエンコード：", ANSICharset: "ANSI 文字セット：",
	NewDocLanguage: "新規文書の言語：", SystemDefault: "システムと同じ", ByLanguage: "インターフェースの言語に従う",
	LargeFileMB: "大きなファイルのしきい値（MB）：", CheckFileChanges: "他のプログラムによる変更を検出", AutoReload: "未変更なら確認せずに再読み込み",

	AboutTitle:       func(name string) string { return name + " について" },
	AboutDescription: "テキストとソースコードのエディター", Version: "バージョン", Author: "作者：", License: "ライセンス：",
	VisitWebsite: "ウェブサイトを開く", Close: "閉じる",
}
