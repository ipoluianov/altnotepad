package editor

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

// checkBuffer compares everything the buffer tells with the plain text
func checkBuffer(t *testing.T, b *Buffer, want []byte) {
	t.Helper()
	if b.Len() != len(want) {
		t.Fatalf("Len = %d, want %d", b.Len(), len(want))
	}
	if got := b.Bytes(0, b.Len()); !bytes.Equal(got, want) {
		t.Fatalf("text differs")
	}
	lines := bytes.Split(want, []byte("\n"))
	if b.LineCount() != len(lines) {
		t.Fatalf("LineCount = %d, want %d", b.LineCount(), len(lines))
	}
	off := 0
	for i, l := range lines {
		if s := b.LineStart(i); s != off {
			t.Fatalf("LineStart(%d) = %d, want %d", i, s, off)
		}
		if e := b.LineEnd(i); e != off+len(l) {
			t.Fatalf("LineEnd(%d) = %d, want %d", i, e, off+len(l))
		}
		off += len(l) + 1
	}
	// LineOfOffset on a sample of offsets
	line := 0
	for o := 0; o <= len(want); o += 1 + len(want)/500 {
		for line+1 < len(lines) && b.LineStart(line+1) <= o {
			line++
		}
		for line > 0 && b.LineStart(line) > o {
			line--
		}
		if got := b.LineOfOffset(o); got != line {
			t.Fatalf("LineOfOffset(%d) = %d, want %d", o, got, line)
		}
	}
}

func randomText(r *rand.Rand, n int) []byte {
	words := []string{"a", "bc", "def", "\n", "\n\n", "Привет", "世界", " ", "\t", "x\ny"}
	var sb strings.Builder
	for sb.Len() < n {
		sb.WriteString(words[r.Intn(len(words))])
	}
	return []byte(sb.String())
}

func TestBufferRandomEdits(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	want := randomText(r, 300_000)
	b := NewBufferFromBytes(want)
	checkBuffer(t, b, want)
	for step := 0; step < 400; step++ {
		switch r.Intn(3) {
		case 0, 1:
			off := r.Intn(len(want) + 1)
			n := r.Intn(300)
			if r.Intn(20) == 0 {
				n = r.Intn(200_000)
			}
			ins := randomText(r, n)
			b.Insert(off, ins)
			want = append(want[:off:off], append(ins, want[off:]...)...)
		case 2:
			if len(want) == 0 {
				continue
			}
			off := r.Intn(len(want))
			n := r.Intn(500)
			if r.Intn(20) == 0 {
				n = r.Intn(150_000)
			}
			n = min(n, len(want)-off)
			b.Delete(off, n)
			want = append(want[:off:off], want[off+n:]...)
		}
		if step%40 == 0 {
			checkBuffer(t, b, want)
		}
	}
	checkBuffer(t, b, want)
	// Delete everything
	b.Delete(0, b.Len())
	checkBuffer(t, b, nil)
	b.Insert(0, []byte("a\nb"))
	checkBuffer(t, b, []byte("a\nb"))
}

func TestBufferEmpty(t *testing.T) {
	b := NewBuffer()
	checkBuffer(t, b, nil)
	if b.LineEnd(0) != 0 || b.LineOfOffset(0) != 0 {
		t.Fatal("empty buffer lines")
	}
	b.Insert(0, []byte("\n"))
	checkBuffer(t, b, []byte("\n"))
}

func TestBufferBuilderKeepsRunes(t *testing.T) {
	text := []byte(strings.Repeat("Ж", chunkTarget)) // 2 bytes each: splits land mid-rune
	var bb BufferBuilder
	for i := 0; i < len(text); i += 777 {
		bb.Write(text[i:min(i+777, len(text))])
	}
	b := bb.Buffer()
	checkBuffer(t, b, text)
	for _, c := range b.chunks {
		if len(c.data) > 0 && !bytes.Equal(c.data[:2], []byte("Ж")) {
			t.Fatal("chunk starts inside a rune")
		}
	}
	for off := 0; off < len(text); off += 2 {
		if r, n := b.DecodeRune(off); r != 'Ж' || n != 2 {
			t.Fatalf("DecodeRune(%d) = %q %d", off, r, n)
		}
	}
}

func BenchmarkBufferTyping(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	buf := NewBufferFromBytes(randomText(r, 50<<20))
	off := buf.Len() / 2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Insert(off, []byte("x"))
		off++
		_ = buf.LineOfOffset(off)
		_ = buf.LineStart(buf.LineOfOffset(off))
	}
}

func TestBufferForEachLine(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	text := randomText(r, 400_000)
	b := NewBufferFromBytes(text)
	lines := bytes.Split(text, []byte("\n"))
	for _, first := range []int{0, 1, len(lines) / 2, len(lines) - 1} {
		n := first
		b.ForEachLine(first, func(line int, l []byte) bool {
			if line != n || !bytes.Equal(l, lines[line]) {
				t.Fatalf("line %d differs", line)
			}
			n++
			return true
		})
		if n != len(lines) {
			t.Fatalf("got %d lines, want %d", n, len(lines))
		}
	}
}
