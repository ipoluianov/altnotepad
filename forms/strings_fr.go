package forms

import (
	"fmt"

	"github.com/ipoluianov/nui/ui/i18n"
)

func frP(n int, one, many string) string { return i18n.Plural("fr", n, one, many) }

var fr = Strings{
	Error:      "Erreur",
	NewDocName: "nouveau",
	NormalText: "Texte normal",

	MenuFile: "Fichier", MenuEdit: "Édition", MenuSearch: "Recherche", MenuView: "Affichage", MenuEncoding: "Encodage", MenuLanguage: "Langage",
	MenuSettings: "Paramètres", MenuTools: "Outils", MenuMacro: "Macro", MenuRunMenu: "Exécution", MenuWindow: "Fenêtres",

	MenuNew: "Nouveau", MenuOpen: "Ouvrir...", MenuOpenContainingFolderSub: "Ouvrir le dossier du document",
	MenuOpenContainingFolder: "Gestionnaire de fichiers", MenuOpenDefaultViewer: "Ouvrir avec le programme par défaut", MenuOpenFolderWorkspace: "Ouvrir un dossier comme espace de travail...",
	MenuReload: "Recharger depuis le disque", MenuSave: "Enregistrer", MenuSaveAs: "Enregistrer sous...", MenuSaveCopy: "Enregistrer une copie sous...", MenuSaveAll: "Enregistrer tout",
	MenuRename: "Renommer...", MenuClose: "Fermer", MenuCloseAll: "Fermer tout", MenuCloseMore: "Fermer plusieurs documents",
	MenuCloseAllButActive: "Fermer tout sauf le document actif", MenuCloseLeft: "Fermer tout à gauche", MenuCloseRight: "Fermer tout à droite",
	MenuCloseUnchanged: "Fermer les documents non modifiés", MenuDeleteFile: "Supprimer du disque...", MenuLoadSession: "Charger une session...",
	MenuSaveSession: "Enregistrer la session...", MenuRecentFiles: "Fichiers récents", MenuClearRecent: "Vider la liste des fichiers récents",
	MenuRestoreClosed: "Rouvrir le dernier fichier fermé", MenuExit: "Quitter",
	MenuPrint: "Imprimer...", MenuExportHTML: "Exporter en HTML...", MenuFunctionList: "Liste des fonctions", Filter: "Filtre", MenuDocumentMap: "Carte du document",
	MenuChangeHistory: "Historique des modifications", MenuNextChange: "Modification suivante", MenuPrevChange: "Modification précédente", MenuClearChanges: "Effacer l'historique des modifications",
	SavedTo:         func(path string) string { return "Enregistré dans " + path },
	MenuMultiSelect: "Sélection multiple", MenuMultiAll: "Sélectionner toutes les occurrences", MenuMultiAllCase: "Sélectionner toutes les occurrences (respecter la casse)",
	MenuMultiNext: "Ajouter l'occurrence suivante", MenuMultiUndo: "Annuler la dernière sélection ajoutée", MenuMultiSkip: "Ignorer et ajouter la suivante",

	MenuUndo: "Annuler", MenuRedo: "Rétablir", MenuCut: "Couper", MenuCopy: "Copier", MenuPaste: "Coller", MenuDelete: "Supprimer", MenuSelectAll: "Tout sélectionner",
	MenuInsert: "Insérer", MenuDateShort: "Date et heure (courte)", MenuDateLong: "Date et heure (longue)", MenuDateISO: "Date et heure (ISO 8601)",
	MenuCopyToClipboard: "Copier dans le presse-papiers", MenuCopyFullPath: "Chemin complet du fichier", MenuCopyFileName: "Nom du fichier",
	MenuCopyDirPath: "Chemin du dossier", MenuIndentSub: "Indentation", MenuIndent: "Augmenter l'indentation", MenuUnindent: "Diminuer l'indentation",
	MenuConvertCase: "Convertir la casse", MenuUpperCase: "MAJUSCULES", MenuLowerCase: "minuscules", MenuProperCase: "Première Lettre En Majuscule",
	MenuProperCaseBlend: "Première Lettre En Majuscule (mixte)", MenuSentenceCase: "Majuscule en début de phrase", MenuSentenceCaseBlend: "Majuscule en début de phrase (mixte)",
	MenuInvertCase: "iNVERSER lA cASSE", MenuRandomCase: "cAsSe AlÉaToIrE", MenuLineOperations: "Opérations sur les lignes",
	MenuDuplicateLine: "Dupliquer la ligne", MenuRemoveDupLines: "Supprimer les lignes en double", MenuRemoveConsecutiveDupLines: "Supprimer les doublons consécutifs",
	MenuSplitLines: "Couper les lignes", MenuJoinLines: "Joindre les lignes", MenuMoveLineUp: "Monter la ligne", MenuMoveLineDown: "Descendre la ligne",
	MenuDeleteLine: "Supprimer la ligne", MenuCutLine: "Couper la ligne", MenuTransposeLine: "Échanger avec la ligne précédente",
	MenuRemoveEmptyLines: "Supprimer les lignes vides", MenuRemoveBlankLines: "Supprimer les lignes vides (y compris d'espaces)",
	MenuInsertLineAbove: "Insérer une ligne vide au-dessus", MenuInsertLineBelow: "Insérer une ligne vide en dessous",
	MenuReverseLines: "Inverser l'ordre des lignes", MenuShuffleLines: "Mélanger les lignes",
	MenuSortAsc: "Trier les lignes par ordre croissant", MenuSortDesc: "Trier les lignes par ordre décroissant",
	MenuSortAscCI: "Tri croissant sans tenir compte de la casse", MenuSortDescCI: "Tri décroissant sans tenir compte de la casse",
	MenuSortIntAsc: "Trier comme entiers (croissant)", MenuSortIntDesc: "Trier comme entiers (décroissant)",
	MenuSortDecAsc: "Trier comme décimaux (croissant)", MenuSortDecDesc: "Trier comme décimaux (décroissant)",
	MenuSortLenAsc: "Trier par longueur (croissant)", MenuSortLenDesc: "Trier par longueur (décroissant)",
	MenuComment: "Commentaires", MenuToggleComment: "Commenter/décommenter les lignes", MenuLineComment: "Commenter les lignes",
	MenuLineUncomment: "Décommenter les lignes", MenuBlockComment: "Commentaire de bloc", MenuBlockUncomment: "Retirer le commentaire de bloc",
	MenuAutoCompletion: "Auto-complétion", MenuWordCompletion: "Compléter le mot", MenuEOLConversion: "Conversion des fins de ligne",
	MenuEOLWindows: "Windows (CR LF)", MenuEOLUnix: "Unix (LF)", MenuEOLMac: "Macintosh (CR)",
	MenuBlankOperations: "Opérations sur les espaces", MenuTrimTrailing: "Supprimer les espaces en fin de ligne", MenuTrimLeading: "Supprimer les espaces en début de ligne",
	MenuTrimBoth: "Supprimer les espaces au début et à la fin", MenuEOLToSpace: "Fins de ligne en espaces", MenuRemoveBlankEOL: "Supprimer les espaces et fins de ligne inutiles",
	MenuTabToSpace: "Tabulations en espaces", MenuSpaceToTabAll: "Espaces en tabulations (tout)", MenuSpaceToTabLeading: "Espaces en tabulations (début de ligne)",
	MenuColumnEditor: "Éditeur de colonnes...", MenuReadOnly: "Lecture seule",
	DateShortLayout: "15:04 02/01/2006", DateLongLayout: "15:04:05 02/01/2006 (Monday)",

	MenuFind: "Rechercher...", MenuFindInFiles: "Rechercher dans les fichiers...", MenuFindNext: "Suivant", MenuFindPrev: "Précédent",
	MenuSelectFindNext: "Sélectionner et rechercher le suivant", MenuSelectFindPrev: "Sélectionner et rechercher le précédent", MenuReplace: "Remplacer...",
	MenuIncremental: "Recherche incrémentale", MenuMark: "Marquer...", MenuSearchResults: "Fenêtre des résultats",
	MenuNextResult: "Résultat suivant", MenuPrevResult: "Résultat précédent", MenuGoTo: "Aller à...",
	MenuGotoBrace: "Aller à l'accolade correspondante", MenuSelectBrace: "Sélectionner entre les accolades",
	MenuStyleAll: "Styler toutes les occurrences", MenuStyleOne: "Styler une occurrence", MenuClearStyleSub: "Effacer le style",
	MenuClearAllStyles: "Effacer tous les styles", MenuJumpUp: "Aller vers le haut", MenuJumpDown: "Aller vers le bas", MenuBookmark: "Signets",
	MenuToggleBookmark: "Basculer le signet", MenuNextBookmark: "Signet suivant", MenuPrevBookmark: "Signet précédent",
	MenuClearBookmarks: "Effacer tous les signets", MenuCutBookmarked: "Couper les lignes marquées", MenuCopyBookmarked: "Copier les lignes marquées",
	MenuRemoveBookmarked: "Supprimer les lignes marquées", MenuRemoveUnbookmarked: "Supprimer les lignes non marquées", MenuInverseBookmarks: "Inverser les signets",
	MenuUsingStyle: func(n int) string { return fmt.Sprintf("Avec le style %d", n) },
	MenuClearStyle: func(n int) string { return fmt.Sprintf("Effacer le style %d", n) },

	MenuAlwaysOnTop: "Toujours au premier plan", MenuFullScreen: "Plein écran", MenuShowSymbol: "Afficher les symboles",
	MenuShowSpaces: "Espaces et tabulations", MenuShowEOL: "Fins de ligne", MenuShowAllChars: "Tous les caractères",
	MenuShowIndentGuides: "Guides d'indentation", MenuShowWrapSymbol: "Symbole de retour à la ligne", MenuZoom: "Zoom", MenuZoomIn: "Zoom avant",
	MenuZoomOut: "Zoom arrière", MenuZoomReset: "Zoom par défaut", MenuMoveClone: "Déplacer/cloner le document",
	MenuMoveToOtherView: "Déplacer vers l'autre vue", MenuCloneToOtherView: "Cloner dans l'autre vue", MenuTab: "Onglets",
	MenuNextTab: "Onglet suivant", MenuPrevTab: "Onglet précédent", MenuMoveTabForward: "Déplacer l'onglet vers la droite", MenuMoveTabBackward: "Déplacer l'onglet vers la gauche",
	MenuWordWrap: "Retour à la ligne", MenuLineNumbers: "Numéros de ligne", MenuFocusOtherView: "Passer à l'autre vue", MenuFoldAll: "Tout replier",
	MenuUnfoldAll: "Tout déplier", MenuFoldCurrent: "Replier le niveau actuel", MenuUnfoldCurrent: "Déplier le niveau actuel",
	MenuFoldLevel: "Replier le niveau", MenuUnfoldLevel: "Déplier le niveau", MenuSummary: "Résumé...", MenuFolderWorkspace: "Dossier comme espace de travail",
	MenuMonitoring: "Surveillance (tail -f)", MenuToolbar: "Barre d'outils", MenuStatusBar: "Barre d'état",
	MenuTabN: func(n int) string {
		if n == 9 {
			return "Dernier onglet"
		}
		return fmt.Sprintf("Onglet %d", n)
	},

	MenuEncANSI: "Encoder en ANSI", MenuEncUTF8: "Encoder en UTF-8", MenuEncUTF8BOM: "Encoder en UTF-8-BOM", MenuEncUTF16BE: "Encoder en UTF-16 BE BOM",
	MenuEncUTF16LE: "Encoder en UTF-16 LE BOM", MenuCharacterSets: "Jeux de caractères", MenuConvANSI: "Convertir en ANSI", MenuConvUTF8: "Convertir en UTF-8",
	MenuConvUTF8BOM: "Convertir en UTF-8-BOM", MenuConvUTF16BE: "Convertir en UTF-16 BE BOM", MenuConvUTF16LE: "Convertir en UTF-16 LE BOM",
	CharsetGroups: map[string]string{"Arabic": "Arabe", "Baltic": "Balte", "Celtic": "Celte", "Cyrillic": "Cyrillique",
		"Central European": "Europe centrale", "Chinese": "Chinois", "Greek": "Grec", "Hebrew": "Hébreu", "Japanese": "Japonais",
		"Korean": "Coréen", "North European": "Europe du Nord", "Thai": "Thaï", "Turkish": "Turc", "Vietnamese": "Vietnamien",
		"Western European": "Europe occidentale"},

	MenuPreferences: "Préférences...", MenuColorScheme: "Palette de couleurs", MenuShortcuts: "Raccourcis clavier",
	MenuHashGenerate:  func(name string) string { return "Générer " + name + "..." },
	MenuHashFiles:     func(name string) string { return "Générer " + name + " de fichiers..." },
	MenuHashSelection: func(name string) string { return "Copier " + name + " de la sélection" },
	MenuBase64Encode:  "Encoder en Base64", MenuBase64Decode: "Décoder le Base64", MenuURLEncode: "Encoder l'URL", MenuURLDecode: "Décoder l'URL",
	MenuJSONFormat: "Formater le JSON", MenuJSONMinify: "Compacter le JSON",
	MenuStartRecording: "Démarrer/arrêter l'enregistrement", MenuStopRecording: "Arrêter l'enregistrement", MenuPlayback: "Exécuter",
	MenuSaveMacro: "Enregistrer la macro...", MenuRunMacroMulti: "Exécuter la macro plusieurs fois...", MenuTrimSave: "Supprimer les espaces en fin de ligne et enregistrer",
	MenuRun: "Exécuter...", MenuOpenInBrowser: "Ouvrir dans le navigateur", MenuSearchInternet: "Rechercher sur Internet", MenuOpenSelectedFile: "Ouvrir le fichier (nom sélectionné)",
	MenuWindows: "Fenêtres...", MenuSortTabsByName: "Trier les onglets par nom", MenuSortTabsByPath: "Trier les onglets par chemin",
	MenuHelp: "Aide en ligne", MenuHomePage: "Page d'accueil", MenuAbout: "À propos d'AltNotepad",

	OpenTitle: "Ouvrir", SaveTitle: "Enregistrer", SaveAsTitle: "Enregistrer sous", SaveCopyTitle: "Enregistrer une copie sous", AllFiles: "Tous les fichiers",
	TextFiles: "Fichiers texte", Save: "Enregistrer", DontSave: "Ne pas enregistrer", PathLabel: "Chemin du fichier :",
	SaveChangesAsk: func(path string) string { return "Enregistrer les modifications de « " + path + " » ?" },
	CreateFileAsk:  func(path string) string { return "« " + path + " » n'existe pas. Le créer ?" },
	AlreadyOpen:    func(path string) string { return "« " + path + " » est ouvert dans un autre onglet" },
	FileExists:     func(path string) string { return "« " + path + " » existe déjà" },
	FileNotFound:   func(path string) string { return "Fichier introuvable : " + path },
	Unencodable: func(enc string) string {
		return "Certains caractères ne peuvent pas être enregistrés en " + enc + ". Les enregistrer comme « ? » ?"
	},
	Loading:       func(name string, percent int) string { return fmt.Sprintf("Chargement de %s : %d %%", name, percent) },
	LargeFileMode: "Gros fichier : coloration syntaxique, repli et retour à la ligne désactivés",
	LargeFileMark: "(gros fichier)",
	ReloadTitle:   "Recharger", RenameTitle: "Renommer", NewName: "Nouveau nom :", DeleteFileTitle: "Supprimer du disque",
	ReloadAsk: func(path string) string {
		return "Recharger « " + path + " » ? Les modifications non enregistrées seront perdues."
	},
	DeleteFileAsk:    func(path string) string { return "Supprimer « " + path + " » du disque ?" },
	FileChangedTitle: "Fichier modifié", FileDeletedTitle: "Fichier supprimé",
	FileChangedAsk: func(path string) string {
		return "« " + path + " »\n\nCe fichier a été modifié par un autre programme.\nVoulez-vous le recharger ?"
	},
	FileChangedModifiedAsk: func(path string) string {
		return "« " + path + " »\n\nCe fichier a été modifié par un autre programme.\nLe recharger et perdre les modifications faites dans l'éditeur ?"
	},
	FileDeletedAsk: func(path string) string {
		return "« " + path + " »\n\nCe fichier n'existe plus.\nLe garder dans l'éditeur ?"
	},
	NotASession: func(name string) string { return name + " n'est pas un fichier de session" },

	StatusSize: func(length, lines int) string { return fmt.Sprintf("longueur : %d   lignes : %d", length, lines) },
	StatusPos:  func(line, col, pos int) string { return fmt.Sprintf("Li : %d   Col : %d   Pos : %d", line, col, pos) },
	StatusSel:  "Sél :",

	GotoTitle: "Aller à...", GotoLine: "Ligne", GotoOffset: "Position", Go: "Aller",
	GotoLineInfo: func(cur, total int) string {
		return fmt.Sprintf("Vous êtes ici : %d   Aller à (1 - %d) :", cur, total)
	},
	GotoOffsetInfo: func(cur, total int) string {
		return fmt.Sprintf("Vous êtes ici : %d   Aller à (0 - %d) :", cur, total)
	},
	ColumnTitle: "Éditeur de colonnes et de sélection multiple", ColumnText: "Texte à insérer", ColumnNumbers: "Nombre à insérer",
	ColumnInitial: "Nombre initial :", ColumnIncrease: "Augmenter de :", ColumnRepeat: "Répéter :", ColumnFormat: "Format :",
	ColumnLeading: "Compléter par :", LeadingNone: "Rien", LeadingZeros: "Zéros", LeadingSpaces: "Espaces",
	FullPath: "Chemin complet", Modified: "Modifié", SummaryChars: "Caractères (sans fins de ligne)", SummaryWords: "Mots",
	SummaryLines: "Lignes", SummaryNonBlankLines: "Lignes non vides", SummaryLength: "Longueur du document", SummarySelected: "Caractères sélectionnés",
	Bytes: "octets", RunTitle: "Exécuter", RunLabel: "Programme à exécuter :", Run: "Exécuter",
	Name: "Nom", State: "État", ModifiedState: "modifié", Activate: "Activer", CloseWindows: "Fermer",
	Command: "Commande", Shortcut: "Raccourci", HashInput: "Texte :", HashEachLine: "Traiter chaque ligne séparément",
	CopyToClipboard:   "Copier",
	CopiedToClipboard: func(s string) string { return "Copié dans le presse-papiers : " + s },
	MacroName:         "Nom de la macro :", MacroTimes: "Nombre d'exécutions (* - jusqu'à la fin du fichier) :",

	FindTab: "Rechercher", ReplaceTab: "Remplacer", FindInFilesTab: "Rechercher dans les fichiers", MarkTab: "Marquer", FindWhat: "Rechercher :",
	ReplaceWith: "Remplacer par :", Filters: "Filtres :", Directory: "Dossier :", CurrentFolder: "Dossier actuel",
	MatchCase: "Respecter la casse", WholeWord: "Mot entier uniquement", WrapAround: "Boucler", Backward: "Vers le haut",
	InSelection: "Dans la sélection", InSubfolders: "Dans tous les sous-dossiers", InHidden: "Dans les dossiers cachés", BookmarkLine: "Marquer la ligne",
	PurgeEach: "Purger à chaque recherche", SearchMode: "Mode :", ModeNormal: "Normal", ModeExtended: "Étendu (\\n, \\r, \\t, \\0, \\x...)",
	ModeRegex: "Expression régulière", DotAll: ". trouve les sauts de ligne",
	FindNext: "Suivant", FindPrev: "Précédent", Count: "Compter", FindAllCurrent: "Tout trouver dans le document",
	FindAllOpen: "Tout trouver dans les documents ouverts", Replace: "Remplacer", ReplaceAll: "Tout remplacer",
	ReplaceAllOpen: "Tout remplacer dans les documents ouverts", FindAll: "Tout trouver", ReplaceInFiles: "Remplacer dans les fichiers",
	MarkAll: "Tout marquer", ClearMarks: "Effacer les marques", CopyMarked: "Copier le texte marqué", Wrapped: "Fin atteinte, recherche reprise au début",
	ReadOnlyDoc: "Le document est en lecture seule",
	NotFound:    func(text string) string { return "Texte introuvable : « " + text + " »" },
	CountResult: func(n int) string { return fmt.Sprintf("Nombre : %d %s", n, frP(n, "occurrence", "occurrences")) },
	ReplacedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, frP(n, "occurrence remplacée", "occurrences remplacées"))
	},
	MarkedCount: func(n int) string {
		return fmt.Sprintf("%d %s", n, frP(n, "occurrence marquée", "occurrences marquées"))
	},
	FoundCount: func(n int) string { return fmt.Sprintf("%d %s", n, frP(n, "résultat trouvé", "résultats trouvés")) },
	HitsCount:  func(n int) string { return fmt.Sprintf("(%d %s)", n, frP(n, "résultat", "résultats")) },
	FoundInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s dans %d %s", n, frP(n, "résultat", "résultats"), files, frP(files, "fichier", "fichiers"))
	},
	ReplacedInFiles: func(n, files int) string {
		return fmt.Sprintf("%d %s dans %d %s", n, frP(n, "remplacement", "remplacements"), files, frP(files, "fichier", "fichiers"))
	},
	ReplaceInFilesAsk: func(find, repl, dir string) string {
		return "Remplacer tous les « " + find + " » par « " + repl + " » dans tous les fichiers de\n" + dir + " ?"
	},
	NoSuchFolder:  func(path string) string { return "Dossier introuvable : " + path },
	Searching:     func(path string) string { return "Recherche : " + path },
	SearchResults: "Résultats de la recherche", Stop: "Arrêter", Stopped: "(arrêtée)", Clear: "Effacer", Line: "Ligne",
	ResultsHeader: func(text string, hits, files, searched int) string {
		return fmt.Sprintf("Recherche « %s » (%d %s dans %d %s sur %d parcourus)", text, hits, frP(hits, "résultat", "résultats"), files, frP(files, "fichier", "fichiers"), searched)
	},

	Workspace: "Espace de travail", Refresh: "Actualiser", CollapseAll: "Tout replier",

	PageGeneral: "Général", PageEditing: "Édition", PageDisplay: "Affichage", PageFiles: "Nouveau document et fichiers",
	Language: "Langue :", LanguageSystem: "Comme le système", Theme: "Thème :", ThemeDark: "Sombre", ThemeLight: "Clair",
	ColorScheme: "Palette de couleurs :", SchemeByTheme: "Selon le thème",
	ShowToolbar: "Afficher la barre d'outils", ShowStatusBar: "Afficher la barre d'état", TabCloseButtons: "Boutons de fermeture sur les onglets",
	DoubleClickCloses: "Le double-clic ferme l'onglet", AlwaysOnTop: "Toujours au premier plan", SingleInstance: "Ouvrir les fichiers dans la fenêtre existante",
	RememberSession: "Mémoriser les fichiers ouverts pour la prochaine session", BackupSession: "Conserver les modifications non enregistrées entre les sessions",
	BackupEvery: "Sauvegarde toutes les, secondes :", MaxRecent: "Fichiers récents dans la liste :",
	TabSize: "Taille des tabulations :", InsertSpaces: "Remplacer les tabulations par des espaces", AutoIndent: "Indentation automatique", AutoClose: "Fermer les parenthèses et guillemets",
	AutoCompletion: "Complétion des mots pendant la saisie", SmartHome: "Début va au premier caractère non blanc", MultiEdit: "Édition multiple (Ctrl+Clic)",
	ScrollPast: "Défiler au-delà de la dernière ligne", CopyLineNoSel: "Copier/couper la ligne sans sélection", CaretWidth: "Largeur du curseur :",
	CaretBlink: "Clignotement du curseur, ms (0 - aucun) :", EdgeColumn: "Repère des lignes longues, colonne (0 - aucun) :", WordChars: "Caractères de mot supplémentaires :",
	FontSize: "Taille de la police :", FontFile: "Fichier de police :", BuiltInFont: "JetBrains Mono (intégrée)", Fonts: "Polices", LineNumbers: "Numéros de ligne",
	BookmarkMargin: "Marge des signets", FoldMargin: "Marge de repli", CurrentLine: "Surligner la ligne courante", SmartHighlight: "Surlignage intelligent",
	SmartMatchCase: "Surlignage intelligent : respecter la casse", SmartWholeWord: "Surlignage intelligent : mot entier", BraceMatch: "Surligner les accolades correspondantes",
	WrapSymbol: "Afficher le symbole de retour à la ligne", ChangeHistory: "Historique des modifications dans la marge",
	NewDocEOL: "Fins de ligne des nouveaux documents :", NewDocEncoding: "Encodage des nouveaux documents :", ANSICharset: "Jeu de caractères ANSI :",
	NewDocLanguage: "Langage des nouveaux documents :", SystemDefault: "Comme le système", ByLanguage: "Selon la langue de l'interface",
	LargeFileMB: "Limite des gros fichiers, Mo :", CheckFileChanges: "Détecter les fichiers modifiés par d'autres programmes", AutoReload: "Les recharger sans demander s'ils ne sont pas modifiés",

	AboutTitle:       func(name string) string { return "À propos de " + name },
	AboutDescription: "Un éditeur de texte et de code source", Version: "Version", Author: "Auteur :", License: "Licence :",
	VisitWebsite: "Visiter le site", Close: "Fermer",
}
