package editor

import (
	_ "embed"
	"image"
	"image/color"
	"math"
	"os"
	"sync"

	"github.com/ipoluianov/nui/ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

//go:embed JetBrainsMono-Regular.ttf
var jetBrainsMono []byte

var (
	defaultFontOnce sync.Once
	defaultFont     *sfnt.Font
)

func parsedDefaultFont() *sfnt.Font {
	defaultFontOnce.Do(func() {
		defaultFont, _ = opentype.Parse(jetBrainsMono)
	})
	return defaultFont
}

// Font draws the text of the editor: a monospaced font at a size, with the
// glyphs rendered once and then copied
type Font struct {
	size float64
	sf   *sfnt.Font
	face font.Face
	buf  sfnt.Buffer

	CharWidth  int
	Ascent     int
	Descent    int
	LineHeight int

	glyphs map[glyphKey]*glyph
}

type glyphKey struct {
	r    rune
	bold bool
}

// glyph is the alpha mask of a character, placed relative to the pen at the baseline
type glyph struct {
	mask   []byte
	w, h   int
	dx, dy int
}

// NewFont makes the font at the size in pixels. fontFile is a TrueType or
// OpenType file to use instead of the built-in JetBrains Mono ("" - built-in).
func NewFont(size float64, fontFile string) *Font {
	sf := parsedDefaultFont()
	if fontFile != "" {
		if data, err := os.ReadFile(fontFile); err == nil {
			if f, err := opentype.Parse(data); err == nil {
				sf = f
			} else if coll, err := opentype.ParseCollection(data); err == nil && coll.NumFonts() > 0 {
				if f, err := coll.Font(0); err == nil {
					sf = f
				}
			}
		}
	}
	size = math.Max(4, math.Min(size, 200))
	face, err := opentype.NewFace(sf, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		sf = parsedDefaultFont()
		face, _ = opentype.NewFace(sf, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	}
	f := &Font{size: size, sf: sf, face: face, glyphs: make(map[glyphKey]*glyph)}
	m := face.Metrics()
	f.Ascent = m.Ascent.Ceil()
	f.Descent = m.Descent.Ceil()
	f.LineHeight = f.Ascent + f.Descent + int(math.Round(size*0.25))
	adv, ok := face.GlyphAdvance('0')
	if !ok {
		adv = fixed.I(int(size * 0.6))
	}
	f.CharWidth = max(1, adv.Round())
	return f
}

// Size returns the size of the font in pixels
func (f *Font) Size() float64 { return f.size }

func (f *Font) hasGlyph(r rune) bool {
	i, err := f.sf.GlyphIndex(&f.buf, r)
	return err == nil && i != 0
}

func (f *Font) glyph(r rune, bold bool) *glyph {
	key := glyphKey{r, bold}
	if g, ok := f.glyphs[key]; ok {
		return g
	}
	var g *glyph
	if f.hasGlyph(r) {
		g = f.renderOwn(r)
	} else {
		g = f.renderFallback(r)
	}
	if bold && g != nil {
		g = emboldened(g)
	}
	if len(f.glyphs) > 20000 {
		f.glyphs = make(map[glyphKey]*glyph)
	}
	f.glyphs[key] = g
	return g
}

func (f *Font) renderOwn(r rune) *glyph {
	dr, mask, mp, _, ok := f.face.Glyph(fixed.P(0, 0), r)
	if !ok || dr.Empty() {
		return &glyph{}
	}
	g := &glyph{w: dr.Dx(), h: dr.Dy(), dx: dr.Min.X, dy: dr.Min.Y}
	g.mask = make([]byte, g.w*g.h)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			_, _, _, a := mask.At(mp.X+x, mp.Y+y).RGBA()
			g.mask[y*g.w+x] = uint8(a >> 8)
		}
	}
	return g
}

// renderFallback draws a character the font does not have with the fonts
// of the interface, which fall back to the system ones (e.g. for Chinese)
func (f *Font) renderFallback(r rune) *glyph {
	s := string(r)
	w, h, err := ui.MeasureText(ui.FontFamilyMono, f.size, s)
	if err != nil || w <= 0 || h <= 0 {
		return &glyph{}
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	ui.DrawText(img, s, color.White, ui.FontFamilyMono, f.size, 0, 0, 0, 0, w, h)
	g := &glyph{w: w, h: h, dx: 0}
	// The text is drawn with its top at 0: the baseline is at the ascent
	_, lh, _ := ui.MeasureText(ui.FontFamilyMono, f.size, "Mg")
	desc := lh - f.Ascent
	if desc < 0 {
		desc = 0
	}
	g.dy = -(h - desc)
	g.mask = make([]byte, w*h)
	empty := true
	for i := range g.mask {
		a := img.Pix[i*4+3]
		g.mask[i] = a
		if a != 0 {
			empty = false
		}
	}
	if empty {
		return &glyph{}
	}
	return g
}

// emboldened widens the glyph by a pixel: a bold look from the regular font
func emboldened(g *glyph) *glyph {
	if g.w == 0 {
		return g
	}
	b := &glyph{w: g.w + 1, h: g.h, dx: g.dx, dy: g.dy}
	b.mask = make([]byte, b.w*b.h)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			a := g.mask[y*g.w+x]
			i := y*b.w + x
			if a > b.mask[i] {
				b.mask[i] = a
			}
			if a > b.mask[i+1] {
				b.mask[i+1] = a
			}
		}
	}
	return b
}

// GlyphWidth returns the width of the drawn character in pixels, for the
// characters that do not fit a cell, e.g. of a fallback font
func (f *Font) GlyphWidth(r rune) int {
	g := f.glyph(r, false)
	return g.dx + g.w
}

// DrawRune draws the character with the pen at x on the baseline y, in the
// clip rectangle of dst
func (f *Font) DrawRune(dst *image.RGBA, clip image.Rectangle, x, y int, r rune, col color.RGBA, bold bool) {
	g := f.glyph(r, bold)
	if g == nil || g.w == 0 {
		return
	}
	blitMask(dst, clip, x+g.dx, y+g.dy, g.w, g.h, g.mask, col)
}

// blitMask blends col into dst through the alpha mask at (x, y)
func blitMask(dst *image.RGBA, clip image.Rectangle, x, y, w, h int, mask []byte, col color.RGBA) {
	clip = clip.Intersect(dst.Rect)
	x0, y0 := max(x, clip.Min.X), max(y, clip.Min.Y)
	x1, y1 := min(x+w, clip.Max.X), min(y+h, clip.Max.Y)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	cr, cg, cb := uint32(col.R), uint32(col.G), uint32(col.B)
	ca := uint32(col.A)
	for yy := y0; yy < y1; yy++ {
		row := dst.Pix[(yy-dst.Rect.Min.Y)*dst.Stride+(x0-dst.Rect.Min.X)*4:]
		m := mask[(yy-y)*w+(x0-x):]
		for xx := 0; xx < x1-x0; xx++ {
			a := uint32(m[xx])
			if a == 0 {
				continue
			}
			if ca != 255 {
				a = a * ca / 255
			}
			p := row[xx*4 : xx*4+4 : xx*4+4]
			if a == 255 {
				p[0], p[1], p[2], p[3] = uint8(cr), uint8(cg), uint8(cb), 255
				continue
			}
			ia := 255 - a
			p[0] = uint8((cr*a + uint32(p[0])*ia) / 255)
			p[1] = uint8((cg*a + uint32(p[1])*ia) / 255)
			p[2] = uint8((cb*a + uint32(p[2])*ia) / 255)
			p[3] = 255
		}
	}
}

// fillRect fills the rectangle of dst within clip with col, blending when it is translucent
func fillRect(dst *image.RGBA, clip image.Rectangle, x, y, w, h int, col color.RGBA) {
	r := image.Rect(x, y, x+w, y+h).Intersect(clip).Intersect(dst.Rect)
	if r.Empty() {
		return
	}
	if col.A == 255 {
		row := dst.Pix[(r.Min.Y-dst.Rect.Min.Y)*dst.Stride+(r.Min.X-dst.Rect.Min.X)*4:]
		n := r.Dx() * 4
		line := row[:n]
		for i := 0; i < n; i += 4 {
			line[i], line[i+1], line[i+2], line[i+3] = col.R, col.G, col.B, 255
		}
		for yy := r.Min.Y + 1; yy < r.Max.Y; yy++ {
			copy(dst.Pix[(yy-dst.Rect.Min.Y)*dst.Stride+(r.Min.X-dst.Rect.Min.X)*4:], line)
		}
		return
	}
	a := uint32(col.A)
	ia := 255 - a
	cr, cg, cb := uint32(col.R)*a, uint32(col.G)*a, uint32(col.B)*a
	for yy := r.Min.Y; yy < r.Max.Y; yy++ {
		row := dst.Pix[(yy-dst.Rect.Min.Y)*dst.Stride+(r.Min.X-dst.Rect.Min.X)*4:]
		for xx := 0; xx < r.Dx(); xx++ {
			p := row[xx*4 : xx*4+4 : xx*4+4]
			p[0] = uint8((cr + uint32(p[0])*ia) / 255)
			p[1] = uint8((cg + uint32(p[1])*ia) / 255)
			p[2] = uint8((cb + uint32(p[2])*ia) / 255)
		}
	}
}
