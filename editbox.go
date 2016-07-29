package main

import (
	"github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
	"unicode/utf8"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func rune_advance_len(r rune, pos int) int {
	if r == '\t' {
		return tabstop_length - pos%tabstop_length
	}
	return runewidth.RuneWidth(r)
}

func voffset_coffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset += 1
		voffset += rune_advance_len(r, voffset)
	}
	return
}

func byte_slice_grow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byte_slice_remove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byte_slice_insert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byte_slice_grow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferred_horizontal_threshold = 5
const tabstop_length = 8

type EditBox struct {
	lines          [][]byte
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
}

// Draws the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.line_voffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstop_length
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.line_voffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if eb.line_voffset != 0 {
		termbox.SetCell(x, y, '←', coldef, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferred_horizontal_threshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.line_voffset != 0 {
		threshold = width - ht
	}
	if eb.cursor_voffset-eb.line_voffset >= threshold {
		eb.line_voffset = eb.cursor_voffset + (ht - width + 1)
	}

	if eb.line_voffset != 0 && eb.cursor_voffset-eb.line_voffset < ht {
		eb.line_voffset = eb.cursor_voffset - ht
		if eb.line_voffset < 0 {
			eb.line_voffset = 0
		}
	}
}

func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursor_boffset = boffset
	eb.cursor_voffset, eb.cursor_coffset = voffset_coffset(eb.text, boffset)
}

func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursor_boffset:])
}

func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursor_boffset])
}

func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursor_boffset - size)
}

func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursor_boffset + size)
}

func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursor_boffset]
}

func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byte_slice_insert(eb.text, eb.cursor_boffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursor_voffset - eb.line_voffset
}

func (eb *EditBox) NewLine() {
	eb.lines = append(eb.lines, eb.text)
	eb.text = make([]byte, 0)
	eb.MoveCursorTo(0)
}

var edit_box EditBox

type Frame struct {
	x  int
	y  int
	w  int
	h  int
	fg termbox.Attribute
	bg termbox.Attribute
}

func (f Frame) Draw() {
	// unicode box drawing chars around the edit box
	for i := f.x + 1; i < f.w; i++ {
		termbox.SetCell(i, f.y, '─', f.fg, f.bg)
		termbox.SetCell(i, f.y+f.h, '─', f.fg, f.bg)
	}
	for i := f.y + 1; i < f.h-1; i++ {
		termbox.SetCell(f.x, i, '|', f.fg, f.bg)
		termbox.SetCell(f.x+f.w-1, i, '|', f.fg, f.bg)
	}
	termbox.SetCell(f.x, f.y, '┌', f.fg, f.bg)
	termbox.SetCell(f.x+f.w-1, f.y, '┐', f.fg, f.bg)

	for n, line := range edit_box.lines {
		t := line
		lx := 0
		tabstop := 0
		for {
			if len(t) == 0 {
				break
			}

			rx := lx
			if lx == tabstop {
				tabstop += tabstop_length
			}

			r, size := utf8.DecodeRune(t)
			if r == '\t' {
				for ; lx < tabstop; lx++ {
					rx = lx
					if rx >= 0 {
						termbox.SetCell(f.x+rx+1, f.y+n+1, ' ', f.fg, f.bg)
					}
				}
			} else {
				if rx >= 0 {
					termbox.SetCell(f.x+rx+1, f.y+n+1, r, f.fg, f.bg)
				}
				lx += runewidth.RuneWidth(r)
			}

			t = t[size:]
		}
	}
}

type drawScreenFunc func(w, h, x, y int, coldef termbox.Attribute)

func drawScreen2(w, h, x, y int, coldef termbox.Attribute) {
	edit_box_width := w
	midy := y + (h / 2)
	midx := x

	// unicode box drawing chars around the edit box
	termbox.SetCell(midx, midy, '│', coldef, coldef)
	termbox.SetCell(midx+edit_box_width-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx, midy-(h/2), '┌', coldef, coldef)
	termbox.SetCell(midx, midy+(h/2), '└', coldef, coldef)
	termbox.SetCell(midx+edit_box_width-1, midy-(h/2), '┐', coldef, coldef)
	termbox.SetCell(midx+edit_box_width-1, midy+(h/2), '┘', coldef, coldef)
	fill(midx+1, midy-(h/2), edit_box_width-2, 1, termbox.Cell{Ch: '─'})
	fill(midx+1, midy+(h/2), edit_box_width-2, 1, termbox.Cell{Ch: '─'})

	edit_box.Draw(midx+1, midy, edit_box_width-2, 1)
	termbox.SetCursor(midx+1+edit_box.CursorX(), midy)

	tbprint(midx+1, midy-2, termbox.ColorYellow, coldef, "Press ESC to quit")
}

func drawScreen1(w, h, x, y int, coldef termbox.Attribute) {
	frame := Frame{
		x:  x,
		y:  y,
		w:  w,
		h:  h,
		fg: coldef,
		bg: coldef,
	}
	frame.Draw()
}

func redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	w, h := termbox.Size()

	drawScreen1(w-2, h-1, 1, 1, coldef)
	drawScreen2(w-2, 3, 1, h-3, coldef)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	redraw_all()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				edit_box.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				edit_box.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				edit_box.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				edit_box.DeleteRuneForward()
			case termbox.KeyTab:
				edit_box.InsertRune('\t')
			case termbox.KeySpace:
				edit_box.InsertRune(' ')
			case termbox.KeyCtrlK:
				edit_box.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				edit_box.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				edit_box.MoveCursorToEndOfTheLine()
			case termbox.KeyEnter:
				edit_box.NewLine()
			default:
				if ev.Ch != 0 {
					edit_box.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all()
	}
}
