package editor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Encoding is how the text of a document is stored in its file
type Encoding struct {
	// ID of the character set, see Charsets: "utf-8", "utf-16le", "windows-1251"...
	ID  string
	BOM bool
}

var (
	EncodingUTF8    = Encoding{ID: "utf-8"}
	EncodingUTF8BOM = Encoding{ID: "utf-8", BOM: true}
	EncodingUTF16LE = Encoding{ID: "utf-16le", BOM: true}
	EncodingUTF16BE = Encoding{ID: "utf-16be", BOM: true}
)

// DefaultANSI is the character set of the files that are not Unicode
// ("ANSI" in the menu), set from the settings
var DefaultANSI = "windows-1252"

// Charset is a character set offered in the Encoding menu
type Charset struct {
	ID    string
	Name  string
	Group string // the submenu of Character sets; "" for the Unicode ones
	enc   encoding.Encoding
}

// Charsets are all the character sets, by groups as in the menu
var Charsets = []Charset{
	{"utf-8", "UTF-8", "", unicode.UTF8},
	{"utf-16le", "UTF-16 LE", "", unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)},
	{"utf-16be", "UTF-16 BE", "", unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)},

	{"iso-8859-6", "ISO 8859-6", "Arabic", charmap.ISO8859_6},
	{"windows-1256", "Windows-1256", "Arabic", charmap.Windows1256},
	{"iso-8859-4", "ISO 8859-4", "Baltic", charmap.ISO8859_4},
	{"iso-8859-13", "ISO 8859-13", "Baltic", charmap.ISO8859_13},
	{"windows-1257", "Windows-1257", "Baltic", charmap.Windows1257},
	{"iso-8859-14", "ISO 8859-14", "Celtic", charmap.ISO8859_14},
	{"iso-8859-5", "ISO 8859-5", "Cyrillic", charmap.ISO8859_5},
	{"koi8-r", "KOI8-R", "Cyrillic", charmap.KOI8R},
	{"koi8-u", "KOI8-U", "Cyrillic", charmap.KOI8U},
	{"x-mac-cyrillic", "Macintosh", "Cyrillic", charmap.MacintoshCyrillic},
	{"ibm855", "OEM 855", "Cyrillic", charmap.CodePage855},
	{"ibm866", "OEM 866", "Cyrillic", charmap.CodePage866},
	{"windows-1251", "Windows-1251", "Cyrillic", charmap.Windows1251},
	{"ibm852", "OEM 852", "Central European", charmap.CodePage852},
	{"iso-8859-2", "ISO 8859-2", "Central European", charmap.ISO8859_2},
	{"windows-1250", "Windows-1250", "Central European", charmap.Windows1250},
	{"gb18030", "GB18030", "Chinese", simplifiedchinese.GB18030},
	{"gbk", "GBK (GB2312)", "Chinese", simplifiedchinese.GBK},
	{"big5", "Big5", "Chinese", traditionalchinese.Big5},
	{"iso-8859-7", "ISO 8859-7", "Greek", charmap.ISO8859_7},
	{"windows-1253", "Windows-1253", "Greek", charmap.Windows1253},
	{"ibm862", "OEM 862", "Hebrew", charmap.CodePage862},
	{"iso-8859-8", "ISO 8859-8", "Hebrew", charmap.ISO8859_8},
	{"windows-1255", "Windows-1255", "Hebrew", charmap.Windows1255},
	{"shift_jis", "Shift-JIS", "Japanese", japanese.ShiftJIS},
	{"euc-jp", "EUC-JP", "Japanese", japanese.EUCJP},
	{"iso-2022-jp", "ISO-2022-JP", "Japanese", japanese.ISO2022JP},
	{"euc-kr", "EUC-KR (Windows-949)", "Korean", korean.EUCKR},
	{"iso-8859-10", "ISO 8859-10", "North European", charmap.ISO8859_10},
	{"windows-874", "Windows-874 (TIS-620)", "Thai", charmap.Windows874},
	{"iso-8859-3", "ISO 8859-3", "Turkish", charmap.ISO8859_3},
	{"iso-8859-9", "ISO 8859-9", "Turkish", charmap.ISO8859_9},
	{"windows-1254", "Windows-1254", "Turkish", charmap.Windows1254},
	{"windows-1258", "Windows-1258", "Vietnamese", charmap.Windows1258},
	{"ibm437", "OEM-US (437)", "Western European", charmap.CodePage437},
	{"ibm850", "OEM 850", "Western European", charmap.CodePage850},
	{"ibm858", "OEM 858", "Western European", charmap.CodePage858},
	{"iso-8859-1", "ISO 8859-1", "Western European", charmap.ISO8859_1},
	{"iso-8859-15", "ISO 8859-15", "Western European", charmap.ISO8859_15},
	{"macintosh", "Macintosh", "Western European", charmap.Macintosh},
	{"windows-1252", "Windows-1252", "Western European", charmap.Windows1252},
}

// CharsetByID returns the character set; ok is false for an unknown ID
func CharsetByID(id string) (Charset, bool) {
	id = strings.ToLower(id)
	for _, c := range Charsets {
		if c.ID == id {
			return c, true
		}
	}
	return Charset{}, false
}

// IsUnicode reports whether the encoding is UTF-8 or UTF-16
func (e Encoding) IsUnicode() bool {
	return e.ID == "utf-8" || e.ID == "utf-16le" || e.ID == "utf-16be"
}

// Name returns the name of the encoding as shown in the status bar
func (e Encoding) Name() string {
	switch e.ID {
	case "utf-8":
		if e.BOM {
			return "UTF-8-BOM"
		}
		return "UTF-8"
	case "utf-16le":
		return "UTF-16 LE BOM"
	case "utf-16be":
		return "UTF-16 BE BOM"
	}
	if c, ok := CharsetByID(e.ID); ok {
		return c.Name
	}
	return e.ID
}

func (e Encoding) bom() []byte {
	if !e.BOM {
		return nil
	}
	switch e.ID {
	case "utf-8":
		return []byte{0xEF, 0xBB, 0xBF}
	case "utf-16le":
		return []byte{0xFF, 0xFE}
	case "utf-16be":
		return []byte{0xFE, 0xFF}
	}
	return nil
}

func (e Encoding) charset() encoding.Encoding {
	if c, ok := CharsetByID(e.ID); ok {
		return c.enc
	}
	return unicode.UTF8
}

// LoadedFile is a file read by LoadFile
type LoadedFile struct {
	Buffer   *Buffer
	Encoding Encoding
	EOL      EOL
	// MixedEOL: the file had different line breaks, all are now EOL
	MixedEOL bool
}

// LoadFile reads a file. With a zero enc the encoding is detected.
// progress, if not nil, is called with the bytes read so far.
func LoadFile(path string, enc Encoding, progress func(done, total int64)) (*LoadedFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}
	return loadFrom(f, st.Size(), enc, progress)
}

// LoadBytes reads a file from memory, see LoadFile
func LoadBytes(data []byte, enc Encoding) (*LoadedFile, error) {
	return loadFrom(bytes.NewReader(data), int64(len(data)), enc, nil)
}

func loadFrom(r io.ReadSeeker, size int64, enc Encoding, progress func(done, total int64)) (*LoadedFile, error) {
	head := make([]byte, 64<<10)
	n, err := io.ReadFull(r, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}
	head = head[:n]

	skip := 0
	detect := enc.ID == ""
	if detect {
		enc, skip = detectEncoding(head, size)
	} else {
		skip = len(enc.bom())
		if skip > 0 && !bytes.HasPrefix(head, enc.bom()) {
			skip = 0
		}
	}
	if _, err := r.Seek(int64(skip), io.SeekStart); err != nil {
		return nil, err
	}

	// A file detected as UTF-8 may turn out invalid further on: then it is read again as ANSI
	checkUTF8 := detect && enc.ID == "utf-8"
	lf, err := readDecoded(r, size, enc, checkUTF8, progress)
	if errors.Is(err, errNotUTF8) {
		enc = Encoding{ID: guessANSI(head)}
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		lf, err = readDecoded(r, size, enc, false, progress)
	}
	if err != nil {
		return nil, err
	}
	lf.Encoding = enc
	return lf, nil
}

var errNotUTF8 = errors.New("not UTF-8")

func readDecoded(r io.Reader, size int64, enc Encoding, checkUTF8 bool, progress func(done, total int64)) (*LoadedFile, error) {
	var src io.Reader = r
	if enc.ID != "utf-8" {
		src = transform.NewReader(r, enc.charset().NewDecoder())
	}
	norm := &eolNormalizer{}
	var validator utf8Validator
	buf := make([]byte, 1<<20)
	var done int64
	for {
		n, err := src.Read(buf)
		if n > 0 {
			p := buf[:n]
			if checkUTF8 && !validator.valid(p) {
				return nil, errNotUTF8
			}
			norm.write(p)
			done += int64(n)
			if progress != nil {
				progress(min(done, size), size)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	if checkUTF8 && !validator.finish() {
		return nil, errNotUTF8
	}
	norm.finish()
	lf := &LoadedFile{Buffer: norm.out.Buffer()}
	lf.EOL, lf.MixedEOL = norm.eol()
	return lf, nil
}

// utf8Validator checks UTF-8 text given in pieces that may split a character
type utf8Validator struct {
	tail []byte
}

func (v *utf8Validator) valid(p []byte) bool {
	if len(v.tail) > 0 {
		need := min(len(p), utf8.UTFMax-len(v.tail))
		joined := append(append([]byte(nil), v.tail...), p[:need]...)
		r, size := utf8.DecodeRune(joined)
		if r == utf8.RuneError && size <= 1 {
			if !utf8.FullRune(joined) && need == len(p) {
				v.tail = joined
				return true
			}
			return false
		}
		p = p[size-len(v.tail):]
		v.tail = nil
	}
	// Keep an unfinished character at the end for the next piece
	cut := len(p)
	for k := 1; k < utf8.UTFMax && k <= len(p); k++ {
		if utf8.RuneStart(p[len(p)-k]) {
			if !utf8.FullRune(p[len(p)-k:]) {
				cut = len(p) - k
			}
			break
		}
	}
	if !utf8.Valid(p[:cut]) {
		return false
	}
	if cut < len(p) {
		v.tail = append([]byte(nil), p[cut:]...)
	}
	return true
}

func (v *utf8Validator) finish() bool {
	return len(v.tail) == 0
}

// eolNormalizer turns the line breaks into "\n", counting each kind
type eolNormalizer struct {
	out       BufferBuilder
	pendingCR bool
	crlf, lf  int
	cr        int
	tmp       []byte
}

func (n *eolNormalizer) write(p []byte) {
	if n.pendingCR {
		n.pendingCR = false
		if len(p) > 0 && p[0] == '\n' {
			n.crlf++
			n.out.Write(nlBytes)
			p = p[1:]
		} else {
			n.cr++
			n.out.Write(nlBytes)
		}
	}
	if bytes.IndexByte(p, '\r') < 0 {
		n.lf += bytes.Count(p, nlBytes)
		n.out.Write(p)
		return
	}
	tmp := n.tmp[:0]
	for i := 0; i < len(p); i++ {
		c := p[i]
		switch c {
		case '\r':
			if i+1 < len(p) {
				if p[i+1] == '\n' {
					n.crlf++
					i++
				} else {
					n.cr++
				}
				tmp = append(tmp, '\n')
			} else {
				n.pendingCR = true
			}
		case '\n':
			n.lf++
			tmp = append(tmp, '\n')
		default:
			tmp = append(tmp, c)
		}
	}
	n.out.Write(tmp)
	n.tmp = tmp
}

func (n *eolNormalizer) finish() {
	if n.pendingCR {
		n.pendingCR = false
		n.cr++
		n.out.Write(nlBytes)
	}
}

// eol returns the most used line break; the system one if there are none
func (n *eolNormalizer) eol() (EOL, bool) {
	kinds := 0
	for _, c := range []int{n.crlf, n.lf, n.cr} {
		if c > 0 {
			kinds++
		}
	}
	mixed := kinds > 1
	switch {
	case n.crlf == 0 && n.lf == 0 && n.cr == 0:
		return DefaultEOL(), false
	case n.crlf >= n.lf && n.crlf >= n.cr:
		return EOLCRLF, mixed
	case n.lf >= n.cr:
		return EOLLF, mixed
	}
	return EOLCR, mixed
}

// detectEncoding looks at the start of a file for a BOM or UTF-16 text;
// returns the encoding and the length of the BOM to skip
func detectEncoding(head []byte, size int64) (Encoding, int) {
	switch {
	case bytes.HasPrefix(head, []byte{0xEF, 0xBB, 0xBF}):
		return EncodingUTF8BOM, 3
	case bytes.HasPrefix(head, []byte{0xFF, 0xFE}):
		return EncodingUTF16LE, 2
	case bytes.HasPrefix(head, []byte{0xFE, 0xFF}):
		return EncodingUTF16BE, 2
	}
	if le, be := utf16Guess(head); le {
		return Encoding{ID: "utf-16le"}, 0
	} else if be {
		return Encoding{ID: "utf-16be"}, 0
	}
	// UTF-8 when valid (checked while reading), ANSI otherwise
	return EncodingUTF8, 0
}

// utf16Guess: UTF-16 text without a BOM has zero bytes on one side of
// most of the ASCII characters
func utf16Guess(head []byte) (le, be bool) {
	n := len(head) &^ 1
	if n < 16 {
		return false, false
	}
	n = min(n, 4096)
	evenZero, oddZero := 0, 0
	for i := 0; i < n; i += 2 {
		if head[i] == 0 {
			evenZero++
		}
		if head[i+1] == 0 {
			oddZero++
		}
	}
	pairs := n / 2
	if oddZero > pairs*6/10 && evenZero < pairs/10 {
		return true, false
	}
	if evenZero > pairs*6/10 && oddZero < pairs/10 {
		return false, true
	}
	return false, false
}

// guessANSI picks the 8-bit character set of a text that is not UTF-8:
// a Russian text is recognized in the usual Cyrillic encodings, anything
// else is DefaultANSI
func guessANSI(head []byte) string {
	high := 0
	for _, c := range head {
		if c >= 0x80 {
			high++
		}
	}
	if high < 8 {
		return DefaultANSI
	}
	best, bestScore := DefaultANSI, 0
	for _, id := range []string{"windows-1251", "koi8-r", "ibm866"} {
		cs, _ := CharsetByID(id)
		decoded, err := cs.enc.NewDecoder().Bytes(head)
		if err != nil {
			continue
		}
		score := 0
		for _, r := range string(decoded) {
			if strings.ContainsRune(frequentCyrillic, r) {
				score++
			}
		}
		if score > bestScore {
			best, bestScore = id, score
		}
	}
	// Most of the non-ASCII characters must be the frequent Russian letters
	if bestScore*10 < high*6 {
		return DefaultANSI
	}
	if cs, ok := CharsetByID(DefaultANSI); ok && cs.Group == "Cyrillic" && best == "windows-1251" {
		return DefaultANSI
	}
	return best
}

const frequentCyrillic = "оеаинтсрвлкмдпуяыьгзбчйхжшюцщэфъё"

// ErrUnencodable: the text has characters the encoding cannot store
var ErrUnencodable = errors.New("the text has characters that cannot be stored in this encoding")

// SaveFile writes the buffer to the file with the encoding and line breaks.
// It writes a temporary file next to it first, then replaces the file, so
// a failed save does not lose the old file. With replaceUnsupported the
// characters the encoding cannot store are saved as '?', otherwise such a
// text gives ErrUnencodable.
func SaveFile(path string, b *Buffer, enc Encoding, eol EOL, replaceUnsupported bool) error {
	if real, err := filepath.EvalSymlinks(path); err == nil {
		path = real
	}
	mode := os.FileMode(0644)
	st, statErr := os.Stat(path)
	if statErr == nil {
		mode = st.Mode().Perm()
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		// No rights to create files in the folder: write the file in place
		if statErr == nil {
			return saveInPlace(path, b, enc, eol, replaceUnsupported)
		}
		return err
	}
	tmpName := tmp.Name()
	ok := false
	defer func() {
		if !ok {
			tmp.Close()
			os.Remove(tmpName)
		}
	}()
	if err := WriteEncoded(tmp, b, enc, eol, replaceUnsupported); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, mode); err != nil && statErr == nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		ok = true
		if statErr == nil {
			return saveInPlace(path, b, enc, eol, replaceUnsupported)
		}
		return err
	}
	ok = true
	return nil
}

func saveInPlace(path string, b *Buffer, enc Encoding, eol EOL, replaceUnsupported bool) error {
	// Check the encoding first, so the file is not truncated for nothing
	if !replaceUnsupported {
		if err := WriteEncoded(io.Discard, b, enc, eol, false); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	if err := WriteEncoded(f, b, enc, eol, replaceUnsupported); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// WriteEncoded writes the text with the encoding and the line breaks
func WriteEncoded(w io.Writer, b *Buffer, enc Encoding, eol EOL, replaceUnsupported bool) error {
	bw := bufio.NewWriterSize(w, 1<<20)
	if bom := enc.bom(); bom != nil {
		if _, err := bw.Write(bom); err != nil {
			return err
		}
	}
	var out io.Writer = bw
	var tw *transform.Writer
	var sanitize func(p []byte) []byte
	if enc.ID != "utf-8" {
		tw = transform.NewWriter(bw, enc.charset().NewEncoder())
		out = tw
		if replaceUnsupported {
			sanitize = unsupportedReplacer(enc.charset())
		}
	}
	eolBytes := eol.Bytes()
	var werr error
	b.ForEachPiece(0, b.Len(), func(p []byte) bool {
		if sanitize != nil {
			p = sanitize(p)
		}
		if eol == EOLLF {
			_, werr = out.Write(p)
			return werr == nil
		}
		for len(p) > 0 {
			i := bytes.IndexByte(p, '\n')
			if i < 0 {
				_, werr = out.Write(p)
				return werr == nil
			}
			if _, werr = out.Write(p[:i]); werr != nil {
				return false
			}
			if _, werr = out.Write(eolBytes); werr != nil {
				return false
			}
			p = p[i+1:]
		}
		return true
	})
	if werr == nil && tw != nil {
		werr = tw.Close()
	}
	if werr != nil {
		if strings.Contains(werr.Error(), "unsupported") || strings.Contains(werr.Error(), "rune not supported") {
			return ErrUnencodable
		}
		return werr
	}
	return bw.Flush()
}

// Reinterpret returns the text as if its bytes in the encoding from were
// read in the encoding to (Encoding > "Encode in" of the menu)
func Reinterpret(b *Buffer, from, to Encoding) (*LoadedFile, error) {
	var raw bytes.Buffer
	if err := WriteEncoded(&raw, b, from, EOLLF, true); err != nil {
		return nil, err
	}
	data := raw.Bytes()
	if bom := from.bom(); bom != nil {
		data = data[len(bom):]
	}
	if to.ID == "utf-8" && !utf8.Valid(data) {
		// Invalid bytes are kept as the replacement character
		data = bytes.ToValidUTF8(data, []byte("�"))
	}
	return LoadBytes(data, Encoding{ID: to.ID, BOM: false})
}

// unsupportedReplacer returns a function that replaces the characters the
// encoding cannot store with '?', as Notepad++ does. A piece may end in the
// middle of a character: that character is checked as it comes.
func unsupportedReplacer(enc encoding.Encoding) func(p []byte) []byte {
	ok := map[rune]bool{}
	encoder := enc.NewEncoder()
	var carry []byte
	return func(p []byte) []byte {
		if len(carry) > 0 {
			p = append(carry, p...)
			carry = nil
		}
		out := make([]byte, 0, len(p))
		for i := 0; i < len(p); {
			c := p[i]
			if c < utf8.RuneSelf {
				out = append(out, c)
				i++
				continue
			}
			if !utf8.FullRune(p[i:]) {
				carry = append([]byte(nil), p[i:]...)
				break
			}
			r, size := utf8.DecodeRune(p[i:])
			good, known := ok[r]
			if !known {
				_, err := encoder.Bytes(p[i : i+size])
				good = err == nil
				ok[r] = good
			}
			if good {
				out = append(out, p[i:i+size]...)
			} else {
				out = append(out, '?')
			}
			i += size
		}
		return out
	}
}
