package editor

import (
	"bytes"
	"sort"
	"unicode/utf8"
)

// Buffer holds the text as UTF-8 bytes split into chunks of about
// chunkTarget bytes, so an edit copies one chunk and not the whole text,
// and a file of gigabytes is held without one huge allocation.
// Line breaks are always "\n" inside (see Document for the EOL of the file).
// Offsets are byte offsets from the start of the text.
type Buffer struct {
	chunks []*chunk

	// startOff[i] and startLine[i] are the offset of chunks[i] and the number
	// of line breaks before it; valid for i < validPrefix
	startOff    []int
	startLine   []int
	validPrefix int

	length  int
	nlTotal int

	// Chunks with the line break index built, to drop the oldest ones
	indexed []*chunk
}

type chunk struct {
	data []byte
	nl   int // the number of '\n' in data

	// idx are the positions of the '\n' in data, built on demand
	idx []int32
}

const (
	chunkTarget = 64 << 10
	chunkMax    = 128 << 10
	chunkMin    = 8 << 10

	// Line break indexes kept at once; each is up to 4 bytes per line
	maxIndexedChunks = 512
)

func NewBuffer() *Buffer {
	b := &Buffer{}
	b.chunks = []*chunk{{}}
	return b
}

// NewBufferFromBytes makes a buffer holding a copy of data
func NewBufferFromBytes(data []byte) *Buffer {
	b := &Buffer{}
	b.setChunks(splitIntoChunks(data))
	return b
}

// BufferBuilder builds a buffer from pieces of text appended in order,
// e.g. while reading a file, without holding the whole text twice
type BufferBuilder struct {
	chunks []*chunk
	cur    []byte
}

func (bb *BufferBuilder) Write(p []byte) (int, error) {
	n := len(p)
	for len(p) > 0 {
		if bb.cur == nil {
			bb.cur = make([]byte, 0, chunkTarget)
		}
		space := chunkTarget - len(bb.cur)
		take := min(space, len(p))
		bb.cur = append(bb.cur, p[:take]...)
		p = p[take:]
		if len(bb.cur) >= chunkTarget {
			bb.flush()
		}
	}
	return n, nil
}

func (bb *BufferBuilder) flush() {
	if len(bb.cur) == 0 {
		return
	}
	data := bb.cur
	bb.cur = nil
	// Do not split a UTF-8 sequence between chunks: the tail goes to the next one
	cut := len(data)
	for k := 1; k <= utf8.UTFMax-1 && k <= len(data); k++ {
		if utf8.RuneStart(data[len(data)-k]) {
			if !utf8.FullRune(data[len(data)-k:]) {
				cut = len(data) - k
			}
			break
		}
	}
	if cut == 0 {
		cut = len(data)
	}
	tail := data[cut:]
	data = data[:cut]
	bb.chunks = append(bb.chunks, &chunk{data: data, nl: bytes.Count(data, nlBytes)})
	if len(tail) > 0 {
		bb.cur = make([]byte, 0, chunkTarget)
		bb.cur = append(bb.cur, tail...)
	}
}

// Buffer returns the buffer with everything written
func (bb *BufferBuilder) Buffer() *Buffer {
	bb.flush()
	b := &Buffer{}
	if len(bb.chunks) == 0 {
		bb.chunks = []*chunk{{}}
	}
	b.setChunks(bb.chunks)
	bb.chunks = nil
	return b
}

var nlBytes = []byte{'\n'}

func splitIntoChunks(data []byte) []*chunk {
	if len(data) == 0 {
		return []*chunk{{}}
	}
	var bb BufferBuilder
	bb.Write(data)
	bb.flush()
	return bb.chunks
}

func (b *Buffer) setChunks(chunks []*chunk) {
	b.chunks = chunks
	b.indexed = nil
	b.length = 0
	b.nlTotal = 0
	for _, c := range chunks {
		b.length += len(c.data)
		b.nlTotal += c.nl
	}
	b.validPrefix = 0
}

// Len returns the length of the text in bytes
func (b *Buffer) Len() int {
	return b.length
}

// LineCount returns the number of lines; an empty text has one line
func (b *Buffer) LineCount() int {
	return b.nlTotal + 1
}

func (b *Buffer) invalidateFrom(i int) {
	if i < b.validPrefix {
		b.validPrefix = i
	}
}

func (b *Buffer) ensurePrefix() {
	n := len(b.chunks)
	if b.validPrefix >= n && len(b.startOff) == n {
		return
	}
	if cap(b.startOff) < n {
		off := make([]int, n, n+n/4+4)
		line := make([]int, n, n+n/4+4)
		copy(off, b.startOff[:min(b.validPrefix, len(b.startOff))])
		copy(line, b.startLine[:min(b.validPrefix, len(b.startLine))])
		b.startOff, b.startLine = off, line
	} else {
		b.startOff = b.startOff[:n]
		b.startLine = b.startLine[:n]
	}
	i := min(b.validPrefix, len(b.startOff))
	off, line := 0, 0
	if i > 0 {
		off = b.startOff[i-1] + len(b.chunks[i-1].data)
		line = b.startLine[i-1] + b.chunks[i-1].nl
	}
	for ; i < n; i++ {
		b.startOff[i] = off
		b.startLine[i] = line
		off += len(b.chunks[i].data)
		line += b.chunks[i].nl
	}
	b.validPrefix = n
}

// chunkAt returns the index of the chunk holding offset off; the end of
// the text is in the last chunk
func (b *Buffer) chunkAt(off int) int {
	b.ensurePrefix()
	i := sort.Search(len(b.chunks), func(i int) bool { return b.startOff[i] > off }) - 1
	if i < 0 {
		i = 0
	}
	// Skip empty chunks: the offset belongs to the chunk that has it
	for i < len(b.chunks)-1 && off-b.startOff[i] >= len(b.chunks[i].data) {
		i++
	}
	return i
}

func (b *Buffer) lineIndex(c *chunk) []int32 {
	if c.idx != nil || c.nl == 0 {
		return c.idx
	}
	idx := make([]int32, 0, c.nl)
	data := c.data
	pos := 0
	for {
		i := bytes.IndexByte(data[pos:], '\n')
		if i < 0 {
			break
		}
		idx = append(idx, int32(pos+i))
		pos += i + 1
	}
	c.idx = idx
	b.indexed = append(b.indexed, c)
	if len(b.indexed) > maxIndexedChunks {
		drop := len(b.indexed) - maxIndexedChunks/2
		for _, old := range b.indexed[:drop] {
			old.idx = nil
		}
		b.indexed = append(b.indexed[:0], b.indexed[drop:]...)
	}
	return idx
}

// LineStart returns the offset of the first byte of the line
func (b *Buffer) LineStart(line int) int {
	if line <= 0 {
		return 0
	}
	if line > b.nlTotal {
		return b.length
	}
	b.ensurePrefix()
	// The chunk with the line break number `line` (1-based)
	i := sort.Search(len(b.chunks), func(i int) bool { return b.startLine[i]+b.chunks[i].nl >= line })
	c := b.chunks[i]
	k := line - b.startLine[i] - 1
	idx := b.lineIndex(c)
	return b.startOff[i] + int(idx[k]) + 1
}

// LineEnd returns the offset of the line break of the line, or the end of
// the text for the last line
func (b *Buffer) LineEnd(line int) int {
	if line >= b.nlTotal {
		return b.length
	}
	if line < 0 {
		line = 0
	}
	return b.LineStart(line+1) - 1
}

// LineOfOffset returns the line holding the offset
func (b *Buffer) LineOfOffset(off int) int {
	if off <= 0 {
		return 0
	}
	if off >= b.length {
		return b.nlTotal
	}
	i := b.chunkAt(off)
	c := b.chunks[i]
	rel := off - b.startOff[i]
	if c.nl == 0 {
		return b.startLine[i]
	}
	idx := b.lineIndex(c)
	// The number of line breaks before rel
	n := sort.Search(len(idx), func(k int) bool { return int(idx[k]) >= rel })
	return b.startLine[i] + n
}

// ByteAt returns the byte at the offset, 0 outside the text
func (b *Buffer) ByteAt(off int) byte {
	if off < 0 || off >= b.length {
		return 0
	}
	i := b.chunkAt(off)
	return b.chunks[i].data[off-b.startOff[i]]
}

// Bytes returns a copy of the text from start to end
func (b *Buffer) Bytes(start, end int) []byte {
	start = max(0, start)
	end = min(end, b.length)
	if end <= start {
		return []byte{}
	}
	out := make([]byte, 0, end-start)
	b.forEachPiece(start, end, func(p []byte) bool {
		out = append(out, p...)
		return true
	})
	return out
}

// View returns the text from start to end; when it is within one chunk it
// is not copied, so it must not be changed or kept after the next edit
func (b *Buffer) View(start, end int) []byte {
	start = max(0, start)
	end = min(end, b.length)
	if end <= start {
		return nil
	}
	i := b.chunkAt(start)
	rel := start - b.startOff[i]
	c := b.chunks[i]
	if rel+(end-start) <= len(c.data) {
		return c.data[rel : rel+(end-start)]
	}
	return b.Bytes(start, end)
}

// String returns the whole text
func (b *Buffer) String() string {
	return string(b.Bytes(0, b.length))
}

// forEachPiece calls fn for the pieces of the text from start to end in order,
// without copying; fn returns false to stop
func (b *Buffer) forEachPiece(start, end int, fn func(p []byte) bool) {
	start = max(0, start)
	end = min(end, b.length)
	if end <= start {
		return
	}
	i := b.chunkAt(start)
	off := start
	for off < end && i < len(b.chunks) {
		c := b.chunks[i]
		rel := off - b.startOff[i]
		n := min(len(c.data)-rel, end-off)
		if n > 0 {
			if !fn(c.data[rel : rel+n]) {
				return
			}
			off += n
		}
		i++
	}
}

// ForEachPiece calls fn for the pieces of the text from start to end in order;
// the pieces must not be changed or kept. fn returns false to stop.
func (b *Buffer) ForEachPiece(start, end int, fn func(p []byte) bool) {
	b.forEachPiece(start, end, fn)
}

// DecodeRune returns the character at the offset and its length in bytes
func (b *Buffer) DecodeRune(off int) (rune, int) {
	if off < 0 || off >= b.length {
		return utf8.RuneError, 0
	}
	i := b.chunkAt(off)
	c := b.chunks[i]
	rel := off - b.startOff[i]
	if rel+utf8.UTFMax <= len(c.data) || i == len(b.chunks)-1 {
		return utf8.DecodeRune(c.data[rel:])
	}
	return utf8.DecodeRune(b.Bytes(off, min(off+utf8.UTFMax, b.length)))
}

// DecodeLastRune returns the character that ends at the offset and its length
func (b *Buffer) DecodeLastRune(off int) (rune, int) {
	if off <= 0 || off > b.length {
		return utf8.RuneError, 0
	}
	start := max(0, off-utf8.UTFMax)
	return utf8.DecodeLastRune(b.View(start, off))
}

// Insert inserts text at the offset
func (b *Buffer) Insert(off int, text []byte) {
	if len(text) == 0 {
		return
	}
	off = max(0, min(off, b.length))
	i := b.chunkAt(off)
	// At the start of a chunk the text may as well go to the end of the previous one
	b.ensurePrefix()
	if i > 0 && off == b.startOff[i] && len(b.chunks[i-1].data)+len(text) <= chunkMax &&
		len(b.chunks[i].data)+len(text) > chunkMax {
		i--
	}
	c := b.chunks[i]
	rel := off - b.startOff[i]
	nl := bytes.Count(text, nlBytes)

	if len(c.data)+len(text) <= chunkMax {
		// Chunk data is never changed in place: views of it stay valid
		data := make([]byte, 0, max(len(c.data)+len(text), min(chunkMax, len(c.data)+len(text)+256)))
		data = append(data, c.data[:rel]...)
		data = append(data, text...)
		data = append(data, c.data[rel:]...)
		b.replaceChunk(i, []*chunk{{data: data, nl: c.nl + nl}})
	} else {
		all := make([]byte, 0, len(c.data)+len(text))
		all = append(all, c.data[:rel]...)
		all = append(all, text...)
		all = append(all, c.data[rel:]...)
		b.replaceChunk(i, splitIntoChunks(all))
	}
	b.length += len(text)
	b.nlTotal += nl
}

// Delete removes n bytes starting at the offset
func (b *Buffer) Delete(off, n int) {
	off = max(0, off)
	end := min(off+n, b.length)
	if end <= off {
		return
	}
	first := b.chunkAt(off)
	last := b.chunkAt(end - 1)
	relStart := off - b.startOff[first]
	relEnd := end - b.startOff[last]

	head := b.chunks[first].data[:relStart]
	tail := b.chunks[last].data[relEnd:]
	removedNL := 0
	for i := first; i <= last; i++ {
		removedNL += b.chunks[i].nl
	}
	data := make([]byte, 0, len(head)+len(tail))
	data = append(data, head...)
	data = append(data, tail...)
	keptNL := bytes.Count(data, nlBytes)
	b.length -= end - off
	b.nlTotal -= removedNL - keptNL

	repl := []*chunk{{data: data, nl: keptNL}}
	// Small chunks are joined with a neighbor
	if len(data) < chunkMin && len(b.chunks) > last-first+1 {
		if last+1 < len(b.chunks) && len(b.chunks[last+1].data)+len(data) <= chunkTarget {
			next := b.chunks[last+1]
			joined := make([]byte, 0, len(data)+len(next.data))
			joined = append(joined, data...)
			joined = append(joined, next.data...)
			repl = []*chunk{{data: joined, nl: keptNL + next.nl}}
			last++
		} else if first > 0 && len(b.chunks[first-1].data)+len(data) <= chunkTarget {
			prev := b.chunks[first-1]
			joined := make([]byte, 0, len(data)+len(prev.data))
			joined = append(joined, prev.data...)
			joined = append(joined, data...)
			repl = []*chunk{{data: joined, nl: keptNL + prev.nl}}
			first--
		}
	}
	if len(repl[0].data) == 0 && len(b.chunks)-(last-first+1) > 0 {
		repl = nil
	}
	b.replaceChunks(first, last, repl)
}

func (b *Buffer) replaceChunk(i int, repl []*chunk) {
	b.replaceChunks(i, i, repl)
}

// replaceChunks replaces chunks[first..last] with repl
func (b *Buffer) replaceChunks(first, last int, repl []*chunk) {
	for i := first; i <= last; i++ {
		b.chunks[i].idx = nil
	}
	n := last - first + 1
	if len(repl) == n {
		copy(b.chunks[first:], repl)
	} else {
		chunks := make([]*chunk, 0, len(b.chunks)-n+len(repl))
		chunks = append(chunks, b.chunks[:first]...)
		chunks = append(chunks, repl...)
		chunks = append(chunks, b.chunks[last+1:]...)
		b.chunks = chunks
	}
	if len(b.chunks) == 0 {
		b.chunks = []*chunk{{}}
	}
	b.invalidateFrom(first)
	// Forget the dropped indexes
	kept := b.indexed[:0]
	for _, c := range b.indexed {
		if c.idx != nil {
			kept = append(kept, c)
		}
	}
	b.indexed = kept
}

// Clone returns a copy of the buffer sharing the (never changed in place) chunk data
func (b *Buffer) Clone() *Buffer {
	nb := &Buffer{}
	chunks := make([]*chunk, len(b.chunks))
	for i, c := range b.chunks {
		chunks[i] = &chunk{data: c.data, nl: c.nl}
	}
	nb.setChunks(chunks)
	return nb
}

// forEachLine calls fn for the lines from the line first on, in order, with
// the text of each line without its line break; fn returns false to stop.
// The text must not be changed or kept.
func (b *Buffer) forEachLine(first int, fn func(line int, text []byte) bool) {
	if first > b.nlTotal {
		return
	}
	first = max(first, 0)
	off := b.LineStart(first)
	i := b.chunkAt(off)
	rel := off - b.startOff[i]
	line := first
	var partial []byte
	for ; i < len(b.chunks); i++ {
		data := b.chunks[i].data[rel:]
		rel = 0
		for {
			j := bytes.IndexByte(data, '\n')
			if j < 0 {
				partial = append(partial, data...)
				break
			}
			var text []byte
			if len(partial) > 0 {
				partial = append(partial, data[:j]...)
				text = partial
			} else {
				text = data[:j]
			}
			if !fn(line, text) {
				return
			}
			partial = partial[:0]
			line++
			data = data[j+1:]
		}
	}
	fn(line, partial)
}

// ForEachLine calls fn for the lines from first on, see forEachLine
func (b *Buffer) ForEachLine(first int, fn func(line int, text []byte) bool) {
	b.forEachLine(first, fn)
}
