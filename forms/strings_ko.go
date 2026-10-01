package forms

import "fmt"

var ko = Strings{
	Error:      "오류",
	NewDocName: "새 문서",
	NormalText: "일반 텍스트",

	MenuFile: "파일", MenuEdit: "편집", MenuSearch: "검색", MenuView: "보기", MenuEncoding: "인코딩", MenuLanguage: "언어",
	MenuSettings: "설정", MenuTools: "도구", MenuMacro: "매크로", MenuRunMenu: "실행", MenuWindow: "창",

	MenuNew: "새로 만들기", MenuOpen: "열기...", MenuOpenContainingFolderSub: "파일이 있는 폴더 열기",
	MenuOpenContainingFolder: "파일 관리자", MenuOpenDefaultViewer: "기본 프로그램으로 열기", MenuOpenFolderWorkspace: "폴더를 작업 공간으로 열기...",
	MenuReload: "디스크에서 다시 불러오기", MenuSave: "저장", MenuSaveAs: "다른 이름으로 저장...", MenuSaveCopy: "사본 저장...", MenuSaveAll: "모두 저장",
	MenuRename: "이름 바꾸기...", MenuClose: "닫기", MenuCloseAll: "모두 닫기", MenuCloseMore: "여러 문서 닫기",
	MenuCloseAllButActive: "현재 문서 외 모두 닫기", MenuCloseLeft: "왼쪽 모두 닫기", MenuCloseRight: "오른쪽 모두 닫기",
	MenuCloseUnchanged: "변경되지 않은 문서 모두 닫기", MenuDeleteFile: "디스크에서 삭제...", MenuLoadSession: "세션 불러오기...",
	MenuSaveSession: "세션 저장...", MenuRecentFiles: "최근 파일", MenuClearRecent: "최근 파일 목록 지우기",
	MenuRestoreClosed: "마지막으로 닫은 파일 복원", MenuExit: "끝내기",
	MenuPrint: "인쇄...", MenuExportHTML: "HTML로 내보내기...", MenuFunctionList: "함수 목록", Filter: "필터", MenuDocumentMap: "문서 지도",
	MenuChangeHistory: "변경 기록", MenuNextChange: "다음 변경으로 이동", MenuPrevChange: "이전 변경으로 이동", MenuClearChanges: "변경 기록 지우기",
	SavedTo:         func(path string) string { return path + "에 저장했습니다" },
	MenuMultiSelect: "다중 선택", MenuMultiAll: "일치하는 항목 모두 선택", MenuMultiAllCase: "일치하는 항목 모두 선택(대/소문자 구분)",
	MenuMultiNext: "다음 일치 항목 추가", MenuMultiUndo: "마지막 추가 취소", MenuMultiSkip: "현재 항목 건너뛰고 다음 추가",

	MenuUndo: "실행 취소", MenuRedo: "다시 실행", MenuCut: "잘라내기", MenuCopy: "복사", MenuPaste: "붙여넣기", MenuDelete: "삭제", MenuSelectAll: "모두 선택",
	MenuInsert: "삽입", MenuDateShort: "날짜와 시간(짧게)", MenuDateLong: "날짜와 시간(길게)", MenuDateISO: "날짜와 시간(ISO 8601)",
	MenuCopyToClipboard: "클립보드에 복사", MenuCopyFullPath: "파일 전체 경로", MenuCopyFileName: "파일 이름",
	MenuCopyDirPath: "폴더 경로", MenuIndentSub: "들여쓰기", MenuIndent: "들여쓰기 늘리기", MenuUnindent: "들여쓰기 줄이기",
	MenuConvertCase: "대/소문자 변환", MenuUpperCase: "대문자", MenuLowerCase: "소문자", MenuProperCase: "단어 첫 글자 대문자",
	MenuProperCaseBlend: "단어 첫 글자 대문자(혼합)", MenuSentenceCase: "문장 첫 글자 대문자", MenuSentenceCaseBlend: "문장 첫 글자 대문자(혼합)",
	MenuInvertCase: "대/소문자 반전", MenuRandomCase: "무작위", MenuLineOperations: "줄 작업",
	MenuDuplicateLine: "현재 줄 복제", MenuRemoveDupLines: "중복 줄 제거", MenuRemoveConsecutiveDupLines: "연속된 중복 줄 제거",
	MenuSplitLines: "줄 나누기", MenuJoinLines: "줄 합치기", MenuMoveLineUp: "줄을 위로 이동", MenuMoveLineDown: "줄을 아래로 이동",
	MenuDeleteLine: "현재 줄 삭제", MenuCutLine: "현재 줄 잘라내기", MenuTransposeLine: "이전 줄과 바꾸기",
	MenuRemoveEmptyLines: "빈 줄 제거", MenuRemoveBlankLines: "빈 줄 제거(공백만 있는 줄 포함)",
	MenuInsertLineAbove: "위에 빈 줄 삽입", MenuInsertLineBelow: "아래에 빈 줄 삽입",
	MenuReverseLines: "줄 순서 뒤집기", MenuShuffleLines: "줄 무작위로 섞기",
	MenuSortAsc: "줄 오름차순 정렬", MenuSortDesc: "줄 내림차순 정렬",
	MenuSortAscCI: "오름차순 정렬(대/소문자 무시)", MenuSortDescCI: "내림차순 정렬(대/소문자 무시)",
	MenuSortIntAsc: "정수로 오름차순 정렬", MenuSortIntDesc: "정수로 내림차순 정렬",
	MenuSortDecAsc: "소수로 오름차순 정렬", MenuSortDecDesc: "소수로 내림차순 정렬",
	MenuSortLenAsc: "길이 오름차순 정렬", MenuSortLenDesc: "길이 내림차순 정렬",
	MenuComment: "주석", MenuToggleComment: "줄 주석 전환", MenuLineComment: "줄 주석 달기",
	MenuLineUncomment: "줄 주석 해제", MenuBlockComment: "블록 주석", MenuBlockUncomment: "블록 주석 해제",
	MenuAutoCompletion: "자동 완성", MenuWordCompletion: "단어 완성", MenuEOLConversion: "줄 끝 변환",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "공백 작업", MenuTrimTrailing: "줄 끝 공백 제거", MenuTrimLeading: "줄 앞 공백 제거",
	MenuTrimBoth: "줄 앞뒤 공백 제거", MenuEOLToSpace: "줄 끝을 공백으로", MenuRemoveBlankEOL: "불필요한 공백과 줄 끝 제거",
	MenuTabToSpace: "탭을 공백으로", MenuSpaceToTabAll: "공백을 탭으로(모두)", MenuSpaceToTabLeading: "공백을 탭으로(줄 앞)",
	MenuColumnEditor: "열 편집기...", MenuReadOnly: "읽기 전용",
	DateShortLayout: "15:04 2006-01-02", DateLongLayout: "2006년 1월 2일 15:04:05 (Monday)",

	MenuFind: "찾기...", MenuFindInFiles: "파일에서 찾기...", MenuFindNext: "다음 찾기", MenuFindPrev: "이전 찾기",
	MenuSelectFindNext: "선택하고 다음 찾기", MenuSelectFindPrev: "선택하고 이전 찾기", MenuReplace: "바꾸기...",
	MenuIncremental: "증분 검색", MenuMark: "표시...", MenuSearchResults: "검색 결과 창",
	MenuNextResult: "다음 결과", MenuPrevResult: "이전 결과", MenuGoTo: "이동...",
	MenuGotoBrace: "짝이 맞는 괄호로 이동", MenuSelectBrace: "괄호 사이 모두 선택",
	MenuStyleAll: "일치하는 항목 모두 강조", MenuStyleOne: "일치하는 항목 하나 강조", MenuClearStyleSub: "강조 해제",
	MenuClearAllStyles: "모든 강조 해제", MenuJumpUp: "위로 이동", MenuJumpDown: "아래로 이동", MenuBookmark: "책갈피",
	MenuToggleBookmark: "책갈피 전환", MenuNextBookmark: "다음 책갈피", MenuPrevBookmark: "이전 책갈피",
	MenuClearBookmarks: "모든 책갈피 지우기", MenuCutBookmarked: "책갈피 줄 잘라내기", MenuCopyBookmarked: "책갈피 줄 복사",
	MenuRemoveBookmarked: "책갈피 줄 제거", MenuRemoveUnbookmarked: "책갈피 없는 줄 제거", MenuInverseBookmarks: "책갈피 반전",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("스타일 %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("스타일 %d 해제", n) },

	MenuAlwaysOnTop: "항상 위", MenuFullScreen: "전체 화면", MenuShowSymbol: "기호 표시",
	MenuShowSpaces: "공백과 탭", MenuShowEOL: "줄 끝", MenuShowAllChars: "모든 문자",
	MenuShowIndentGuides: "들여쓰기 안내선", MenuShowWrapSymbol: "줄 바꿈 기호", MenuZoom: "확대/축소", MenuZoomIn: "확대",
	MenuZoomOut: "축소", MenuZoomReset: "기본 크기", MenuMoveClone: "문서 이동/복제",
	MenuMoveToOtherView: "다른 보기로 이동", MenuCloneToOtherView: "다른 보기로 복제", MenuTab: "탭",
	MenuNextTab: "다음 탭", MenuPrevTab: "이전 탭", MenuMoveTabForward: "탭을 오른쪽으로 이동", MenuMoveTabBackward: "탭을 왼쪽으로 이동",
	MenuWordWrap: "자동 줄 바꿈", MenuLineNumbers: "줄 번호", MenuFocusOtherView: "다른 보기로 전환", MenuFoldAll: "모두 접기",
	MenuUnfoldAll: "모두 펼치기", MenuFoldCurrent: "현재 수준 접기", MenuUnfoldCurrent: "현재 수준 펼치기",
	MenuFoldLevel: "수준 접기", MenuUnfoldLevel: "수준 펼치기", MenuSummary: "요약...", MenuFolderWorkspace: "폴더 작업 공간",
	MenuMonitoring: "모니터링 (tail -f)", MenuToolbar: "도구 모음", MenuStatusBar: "상태 표시줄",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "마지막 탭"
		}
		return fmt.Sprintf("탭 %d", n)
	},

	MenuEncANSI: "ANSI로 인코딩", MenuEncUTF8: "UTF-8로 인코딩", MenuEncUTF8BOM: "UTF-8-BOM으로 인코딩", MenuEncUTF16BE: "UTF-16 BE BOM으로 인코딩",
	MenuEncUTF16LE: "UTF-16 LE BOM으로 인코딩", MenuCharacterSets: "문자 집합", MenuConvANSI: "ANSI로 변환", MenuConvUTF8: "UTF-8로 변환",
	MenuConvUTF8BOM: "UTF-8-BOM으로 변환", MenuConvUTF16BE: "UTF-16 BE BOM으로 변환", MenuConvUTF16LE: "UTF-16 LE BOM으로 변환",
	CharsetGroups: map[string]string{"Arabic": "아랍어", "Baltic": "발트어", "Celtic": "켈트어", "Cyrillic": "키릴 문자",
		"Central European": "중앙 유럽", "Chinese": "중국어", "Greek": "그리스어", "Hebrew": "히브리어", "Japanese": "일본어",
		"Korean": "한국어", "North European": "북유럽", "Thai": "태국어", "Turkish": "터키어", "Vietnamese": "베트남어",
		"Western European": "서유럽"},

	MenuPreferences: "환경 설정...", MenuColorScheme: "색 구성표", MenuShortcuts: "바로 가기 키",
	MenuHashGenerate:  func(name string) string { return name + " 생성..." },
	MenuHashFiles:     func(name string) string { return "파일의 " + name + " 생성..." },
	MenuHashSelection: func(name string) string { return "선택 영역의 " + name + " 복사" },
	MenuBase64Encode:  "Base64 인코딩", MenuBase64Decode: "Base64 디코딩", MenuURLEncode: "URL 인코딩", MenuURLDecode: "URL 디코딩",
	MenuJSONFormat: "JSON 서식 지정", MenuJSONMinify: "JSON 압축",
	MenuStartRecording: "기록 시작/중지", MenuStopRecording: "기록 중지", MenuPlayback: "재생",
	MenuSaveMacro: "기록한 매크로 저장...", MenuRunMacroMulti: "매크로 여러 번 실행...", MenuTrimSave: "줄 끝 공백 제거 후 저장",
	MenuRun: "실행...", MenuOpenInBrowser: "브라우저에서 열기", MenuSearchInternet: "인터넷에서 검색", MenuOpenSelectedFile: "파일 열기(선택한 이름)",
	MenuWindows: "창...", MenuSortTabsByName: "이름순으로 탭 정렬", MenuSortTabsByPath: "경로순으로 탭 정렬",
	MenuHelp: "온라인 도움말", MenuHomePage: "홈페이지", MenuAbout: "AltNotepad 정보",

	OpenTitle: "열기", SaveTitle: "저장", SaveAsTitle: "다른 이름으로 저장", SaveCopyTitle: "사본 저장", AllFiles: "모든 파일",
	TextFiles: "텍스트 파일", Save: "저장", DontSave: "저장 안 함", PathLabel: "파일 경로:",
	SaveChangesAsk: func(path string) string { return "\"" + path + "\"의 변경 내용을 저장하시겠습니까?" },
	CreateFileAsk:  func(path string) string { return "\"" + path + "\"이(가) 없습니다. 만드시겠습니까?" },
	AlreadyOpen:    func(path string) string { return "\"" + path + "\"이(가) 다른 탭에서 열려 있습니다" },
	FileExists:     func(path string) string { return "\"" + path + "\"이(가) 이미 있습니다" },
	FileNotFound:   func(path string) string { return "파일을 찾을 수 없습니다: " + path },
	Unencodable: func(enc string) string {
		return "일부 문자는 " + enc + "(으)로 저장할 수 없습니다. \"?\"로 저장하시겠습니까?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("%s 불러오는 중: %d%%", name, percent) },
	LargeFileMode: "큰 파일: 구문 강조, 접기, 줄 바꿈이 꺼져 있습니다",
	LargeFileMark: "(큰 파일)",
	ReloadTitle:   "다시 불러오기", RenameTitle: "이름 바꾸기", NewName: "새 이름:", DeleteFileTitle: "디스크에서 삭제",
	ReloadAsk: func(path string) string {
		return "\"" + path + "\"을(를) 다시 불러오시겠습니까? 저장하지 않은 변경 내용은 사라집니다."
	},
	DeleteFileAsk:    func(path string) string { return "\"" + path + "\"을(를) 디스크에서 삭제하시겠습니까?" },
	FileChangedTitle: "파일이 변경됨", FileDeletedTitle: "파일이 삭제됨",
	FileChangedAsk: func(path string) string {
		return "\"" + path + "\"\n\n다른 프로그램이 이 파일을 변경했습니다.\n다시 불러오시겠습니까?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "\"" + path + "\"\n\n다른 프로그램이 이 파일을 변경했습니다.\n다시 불러오고 편집기의 변경 내용을 버리시겠습니까?"
	},
	FileDeletedAsk: func(path string) string {
		return "\"" + path + "\"\n\n이 파일이 더 이상 존재하지 않습니다.\n편집기에 유지하시겠습니까?"
	},
	NotASession: func(name string) string { return name + "은(는) 세션 파일이 아닙니다" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("길이: %d   줄: %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("줄: %d   열: %d   위치: %d", line, col, pos) },
	StatusSel:  "선택:",

	GotoTitle: "이동...", GotoLine: "줄", GotoOffset: "오프셋", Go: "이동",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("현재 위치: %d   이동할 위치 (1 - %d):", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("현재 위치: %d   이동할 위치 (0 - %d):", cur, total)
	},
	ColumnTitle: "열 편집기와 다중 선택", ColumnText: "삽입할 텍스트", ColumnNumbers: "삽입할 숫자",
	ColumnInitial: "시작 숫자:", ColumnIncrease: "증가값:", ColumnRepeat: "반복:", ColumnFormat: "형식:",
	ColumnLeading: "앞 채우기:", LeadingNone: "없음", LeadingZeros: "0", LeadingSpaces: "공백",
	FullPath: "전체 경로", Modified: "수정한 날짜", SummaryChars: "문자 수(줄 끝 제외)", SummaryWords: "단어 수",
	SummaryLines: "줄 수", SummaryNonBlankLines: "비어 있지 않은 줄 수", SummaryLength: "문서 길이", SummarySelected: "선택한 문자 수",
	Bytes: "바이트", RunTitle: "실행", RunLabel: "실행할 프로그램:", Run: "실행",
	Name: "이름", State: "상태", ModifiedState: "수정됨", Activate: "활성화", CloseWindows: "닫기",
	Command: "명령", Shortcut: "바로 가기 키", HashInput: "텍스트:", HashEachLine: "각 줄을 따로 처리",
	CopyToClipboard:   "복사",
	CopiedToClipboard: func(s string) string { return "클립보드에 복사했습니다: " + s },
	MacroName:         "매크로 이름:", MacroTimes: "실행 횟수(* - 파일 끝까지):",

	FindTab: "찾기", ReplaceTab: "바꾸기", FindInFilesTab: "파일에서 찾기", MarkTab: "표시", FindWhat: "찾을 내용:",
	ReplaceWith: "바꿀 내용:", Filters: "필터:", Directory: "폴더:", CurrentFolder: "현재 폴더",
	MatchCase: "대/소문자 구분", WholeWord: "단어 단위로", WrapAround: "처음부터 다시", Backward: "뒤로",
	InSelection: "선택 영역에서", InSubfolders: "모든 하위 폴더에서", InHidden: "숨김 폴더에서", BookmarkLine: "줄에 책갈피",
	PurgeEach: "검색할 때마다 지우기", SearchMode: "검색 모드:", ModeNormal: "일반", ModeExtended: "확장 (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "정규식", DotAll: ".이 줄 바꿈과 일치",
	FindNext: "다음 찾기", FindPrev: "이전 찾기", Count: "개수", FindAllCurrent: "현재 문서에서 모두 찾기",
	FindAllOpen: "열린 모든 문서에서 찾기", Replace: "바꾸기", ReplaceAll: "모두 바꾸기",
	ReplaceAllOpen: "열린 모든 문서에서 바꾸기", FindAll: "모두 찾기", ReplaceInFiles: "파일에서 바꾸기",
	MarkAll: "모두 표시", ClearMarks: "표시 지우기", CopyMarked: "표시한 텍스트 복사", Wrapped: "끝에 도달하여 처음부터 계속했습니다",
	ReadOnlyDoc: "문서가 읽기 전용입니다",
	NotFound:    func(text string) string { return "\"" + text + "\"을(를) 찾을 수 없습니다" },
	CountResult: func(n int) string { return fmt.Sprintf("합계: %d개 일치", n) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d개를 바꿨습니다", n)
	},
	MarkedCount: func(n int) string { return fmt.Sprintf("%d개 일치 항목을 표시했습니다", n) },
	FoundCount:  func(n int) string { return fmt.Sprintf("%d개 결과", n) },
	HitsCount:   func(n int) string { return fmt.Sprintf("(%d개 결과)", n) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("파일 %d개에서 %d개 결과", files, n)
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("파일 %d개에서 %d개 바꿈", files, n)
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "다음 폴더의 모든 파일에서 \"" + find + "\"을(를) \"" + repl + "\"(으)로 바꾸시겠습니까?\n" + dir
	},
	NoSuchFolder:  func(path string) string { return "폴더가 없습니다: " + path },
	Searching:     func(path string) string { return "검색 중: " + path },
	SearchResults: "검색 결과", Stop: "중지", Stopped: "(중지됨)", Clear: "지우기", Line: "줄",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("\"%s\" 검색 (검색한 파일 %d개 중 %d개에서 %d개 결과)", text, searched, files, hits)
	},

	Workspace: "작업 공간", Refresh: "새로 고침", CollapseAll: "모두 접기",

	PageGeneral: "일반", PageEditing: "편집", PageDisplay: "표시", PageFiles: "새 문서와 파일",
	Language: "언어:", LanguageSystem: "시스템 설정", Theme: "테마:", ThemeDark: "어둡게", ThemeLight: "밝게",
	ColorScheme: "색 구성표:", SchemeByTheme: "테마에 따라",
	ShowToolbar: "도구 모음 표시", ShowStatusBar: "상태 표시줄 표시", TabCloseButtons: "탭에 닫기 단추 표시",
	DoubleClickCloses: "두 번 클릭하여 탭 닫기", AlwaysOnTop: "항상 위", SingleInstance: "이미 열린 창에서 파일 열기",
	RememberSession: "다음 세션을 위해 열린 파일 기억", BackupSession: "저장하지 않은 변경 내용을 세션 간에 유지",
	BackupEvery: "백업 간격(초):", MaxRecent: "최근 파일 개수:",
	TabSize: "탭 크기:", InsertSpaces: "탭을 공백으로 바꾸기", AutoIndent: "자동 들여쓰기", AutoClose: "괄호와 따옴표 자동 닫기",
	AutoCompletion: "입력하면서 단어 완성", SmartHome: "Home 키로 첫 번째 공백 아닌 문자로 이동", MultiEdit: "다중 편집(Ctrl+클릭)",
	ScrollPast: "마지막 줄 너머로 스크롤", CopyLineNoSel: "선택 없이 줄 복사/잘라내기", CaretWidth: "커서 너비:",
	CaretBlink: "커서 깜박임, ms(0 - 깜박이지 않음):", EdgeColumn: "긴 줄 표시 열(0 - 없음):", WordChars: "추가 단어 문자:",
	FontSize: "글꼴 크기:", FontFile: "글꼴 파일:", BuiltInFont: "JetBrains Mono(내장)", Fonts: "글꼴", LineNumbers: "줄 번호",
	BookmarkMargin: "책갈피 여백", FoldMargin: "접기 여백", CurrentLine: "현재 줄 강조", SmartHighlight: "스마트 강조",
	SmartMatchCase: "스마트 강조: 대/소문자 구분", SmartWholeWord: "스마트 강조: 단어 단위로", BraceMatch: "짝이 맞는 괄호 강조",
	WrapSymbol: "줄 바꿈 기호 표시", ChangeHistory: "여백에 변경 기록 표시",
	NewDocEOL: "새 문서의 줄 끝:", NewDocEncoding: "새 문서의 인코딩:", ANSICharset: "ANSI 문자 집합:",
	NewDocLanguage: "새 문서의 언어:", SystemDefault: "시스템 설정", ByLanguage: "인터페이스 언어에 따라",
	LargeFileMB: "큰 파일 기준(MB):", CheckFileChanges: "다른 프로그램이 변경한 파일 감지", AutoReload: "수정하지 않았으면 묻지 않고 다시 불러오기",

	AboutTitle:       func(name string) string { return name + " 정보" },
	AboutDescription: "텍스트 및 소스 코드 편집기", Version: "버전", Author: "만든 이:", License: "라이선스:",
	VisitWebsite: "웹사이트 방문", Close: "닫기",
}
