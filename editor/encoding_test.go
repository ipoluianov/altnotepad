package editor

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

func roundTrip(t *testing.T, raw []byte, wantEnc string, wantEOL EOL, wantText string) {
	t.Helper()
	lf, err := LoadBytes(raw, Encoding{})
	if err != nil {
		t.Fatal(err)
	}
	if lf.Encoding.Name() != wantEnc {
		t.Errorf("encoding %s, want %s", lf.Encoding.Name(), wantEnc)
	}
	if lf.EOL != wantEOL {
		t.Errorf("EOL %v, want %v", lf.EOL, wantEOL)
	}
	if got := lf.Buffer.String(); got != wantText {
		t.Errorf("text %q, want %q", got, wantText)
	}
	var out bytes.Buffer
	if err := WriteEncoded(&out, lf.Buffer, lf.Encoding, lf.EOL, false); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out.Bytes(), raw) {
		t.Errorf("saved %q, want %q", out.Bytes(), raw)
	}
}

func TestEncodings(t *testing.T) {
	text := "Привет, мир!\nВторая строка с текстом на русском языке\n"
	crlf := strings.ReplaceAll(text, "\n", "\r\n")

	roundTrip(t, []byte(text), "UTF-8", EOLLF, text)
	roundTrip(t, append([]byte{0xEF, 0xBB, 0xBF}, crlf...), "UTF-8-BOM", EOLCRLF, text)

	cp1251, _ := charmap.Windows1251.NewEncoder().Bytes([]byte(crlf))
	roundTrip(t, cp1251, "Windows-1251", EOLCRLF, text)

	koi, _ := charmap.KOI8R.NewEncoder().Bytes([]byte(text))
	roundTrip(t, koi, "KOI8-R", EOLLF, text)

	u16, _ := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Bytes([]byte(crlf))
	roundTrip(t, u16, "UTF-16 LE BOM", EOLCRLF, text)

	u16be, _ := unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewEncoder().Bytes([]byte(text))
	roundTrip(t, u16be, "UTF-16 BE BOM", EOLLF, text)

	// Old Mac line breaks
	roundTrip(t, []byte("a\rb\rc"), "UTF-8", EOLCR, "a\nb\nc")
}

func TestMixedEOLAndBinary(t *testing.T) {
	lf, err := LoadBytes([]byte("a\r\nb\nc\r\nd\re"), Encoding{})
	if err != nil {
		t.Fatal(err)
	}
	if lf.EOL != EOLCRLF || !lf.MixedEOL || lf.Buffer.String() != "a\nb\nc\nd\ne" {
		t.Errorf("mixed: %v %v %q", lf.EOL, lf.MixedEOL, lf.Buffer.String())
	}
	bin := []byte{0, 1, 2, 'a', 0xFF, 0xFE, 'b', 0}
	lf, err = LoadBytes(bin, Encoding{})
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	WriteEncoded(&out, lf.Buffer, lf.Encoding, lf.EOL, false)
	if !bytes.Equal(out.Bytes(), bin) {
		t.Errorf("binary round trip: %v (%s)", out.Bytes(), lf.Encoding.Name())
	}
}

func TestUnencodable(t *testing.T) {
	b := NewBufferFromBytes([]byte("Привет 中文"))
	var out bytes.Buffer
	err := WriteEncoded(&out, b, Encoding{ID: "windows-1251"}, EOLLF, false)
	if !errors.Is(err, ErrUnencodable) {
		t.Fatalf("err %v", err)
	}
	out.Reset()
	if err := WriteEncoded(&out, b, Encoding{ID: "windows-1251"}, EOLLF, true); err != nil {
		t.Fatal(err)
	}
	dec, _ := charmap.Windows1251.NewDecoder().Bytes(out.Bytes())
	if string(dec) != "Привет ??" {
		t.Errorf("replaced: %q", dec)
	}
}

func TestReinterpret(t *testing.T) {
	// A CP1251 file read as ANSI 1252 by mistake, then read again as 1251
	raw, _ := charmap.Windows1251.NewEncoder().Bytes([]byte("Тест"))
	lf, _ := LoadBytes(raw, Encoding{ID: "windows-1252"})
	back, err := Reinterpret(lf.Buffer, Encoding{ID: "windows-1252"}, Encoding{ID: "windows-1251"})
	if err != nil {
		t.Fatal(err)
	}
	if back.Buffer.String() != "Тест" {
		t.Errorf("got %q", back.Buffer.String())
	}
}
