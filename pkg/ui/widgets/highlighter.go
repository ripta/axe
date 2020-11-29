package widgets

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
)

type Highlighter struct {
	views.WidgetWatchers

	lights []int
	curr   int
	kw     []rune

	t views.Text
}

func NewHighlighter() *Highlighter {
	return &Highlighter{}
}

func (h *Highlighter) Append(line string) {
	s := h.t.Text()
	s += line + "\n"
	h.t.SetText(s)

	h.reset()
	h.PostEventWidgetContent(h)
}

func (h *Highlighter) Clear() {
	h.t.SetText("")
	h.curr = 0
	h.kw = nil
	h.lights = nil
}

func (h *Highlighter) Count() int {
	return len(h.lights)
}

func (h *Highlighter) Current() int {
	return h.curr
}

func (h *Highlighter) Draw() {
	h.t.Draw()
}

func (h *Highlighter) HandleEvent(e tcell.Event) bool {
	return h.t.HandleEvent(e)
}

func (h *Highlighter) Highlight(idx int) bool {
	if idx < 0 || idx >= len(h.lights) {
		return false
	}

	h.curr = idx
	h.reset()
	h.PostEventWidgetContent(h)
	return true
}

func (h *Highlighter) Pos(idx int) (int, int, bool) {
	if idx < 0 || idx >= len(h.lights) {
		return 0, 0, false
	}

	var x, y int
	is := h.lights[idx]
	for i, c := range []rune(h.t.Text()) {
		if i == is {
			break
		}
		if c == '\n' {
			x = 0
			y++
			continue
		}
		x += runewidth.RuneWidth(c)
	}

	return x, y, true
}

func (h *Highlighter) Resize() {
	h.t.Resize()
}

func (h *Highlighter) SetKeyword(kw string) {
	h.kw = []rune(kw)
	h.curr = -1
	h.t.SetStyle(h.t.Style())
	if len(kw) == 0 {
		return
	}

	h.reset()
	h.PostEventWidgetContent(h)
}

func (h *Highlighter) Keyword() string {
	return string(h.kw)
}

func (h *Highlighter) SetView(v views.View) {
	h.t.SetView(v)
}

func (h *Highlighter) Size() (int, int) {
	return h.t.Size()
}

func (h *Highlighter) reset() {
	h.lights = nil

	s := h.t.Text()
	kw := string(h.kw)

	for x := 0; len(kw) > 0; {
		i := strings.Index(s, kw)
		if i == -1 {
			break
		}

		is := len([]rune(s[:i])) + x
		h.lights = append(h.lights, is)
		s = s[i+len(kw):]
		x += len(kw) + i
	}

	current := h.t.Style().Background(tcell.ColorYellow)
	reverse := h.t.Style().Reverse(true)
	for i, is := range h.lights {
		for j := range h.kw {
			if i == h.curr {
				h.t.SetStyleAt(is+j, current)
			} else {
				h.t.SetStyleAt(is+j, reverse)
			}
		}
	}
}
