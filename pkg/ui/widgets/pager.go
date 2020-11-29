package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Pager struct {
	views.WidgetWatchers

	app *views.Application
	h   *Highlighter
	v   views.View
	vp  views.ViewPort
}

func NewPager(app *views.Application) *Pager {
	p := &Pager{
		app: app,
		h:   NewHighlighter(),
	}
	p.h.SetView(&p.vp)
	return p
}

func (p *Pager) Append(line string) {
	p.h.Append(line)
	w, h := p.h.Size()
	p.vp.SetContentSize(w, h, true)
	p.vp.ValidateView()
}

func (p *Pager) Clear() {
	p.h.Clear()
	w, h := p.h.Size()
	p.vp.SetContentSize(w, h, true)
	p.vp.ValidateView()
}

func (p *Pager) Draw() {
	if p.v == nil {
		return
	}
	p.h.Draw()
}

func (p *Pager) GetScrollPercentage() float64 {
	_, vy, _, _ := p.vp.GetVisible()
	_, cy := p.vp.GetContentSize()
	_, vh := p.vp.Size()
	return float64(vy) / float64(cy-vh)
}

func (p *Pager) HandleEvent(e tcell.Event) bool {
	return false
}

func (p *Pager) HighlightNext() bool {
	return p.Highlight(p.h.Current() + 1)
}

func (p *Pager) HighlightPrev() bool {
	return p.Highlight(p.h.Current() - 1)
}

func (p *Pager) Highlight(idx int) bool {
	c := p.h.Count()
	if c == 0 {
		return false
	}

	if idx >= c {
		idx = 0
	} else if idx < 0 {
		idx = c - 1
	}

	if ok := p.h.Highlight(idx); !ok {
		return false
	}

	x, y, ok := p.h.Pos(idx)
	if !ok {
		return false
	}

	p.vp.Center(x, y)
	p.PostEventWidgetContent(p)
	return true
}

func (p *Pager) Resize() {
	w, h := p.v.Size()
	p.vp.Resize(0, 0, w, h)
	p.vp.ValidateView()
}

func (p *Pager) ScrollDown(rows int) {
	p.vp.ScrollDown(rows)
}

func (p *Pager) ScrollToBeginning() {
	_, h := p.v.Size()
	p.vp.ScrollUp(h)
}

func (p *Pager) ScrollToEnd() {
	_, h := p.v.Size()
	p.vp.ScrollDown(h)
}

func (p *Pager) ScrollUp(rows int) {
	p.vp.ScrollUp(rows)
}

func (p *Pager) SetKeyword(kw string) {
	p.h.SetKeyword(kw)
	p.PostEventWidgetContent(p)
}

func (p *Pager) Keyword() string {
	return p.h.Keyword()
}

func (p *Pager) SetView(v views.View) {
	p.v = v
	p.vp.SetView(v)
	if v == nil {
		return
	}
	p.Resize()
}

func (p *Pager) Size() (int, int) {
	w, h := p.v.Size()
	if w > 2 {
		w = 2
	}
	if h > 2 {
		h = 2
	}
	return w, h
}
