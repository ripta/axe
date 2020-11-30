package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/ripta/axe/pkg/ui/themes"
	"github.com/ripta/axe/pkg/ui/widgets"
)

type UI struct {
	views.BoxLayout

	app       *views.Application
	statusbar *widgets.Statusbar

	autoscroll bool
	pager      *widgets.Pager
}

func New(app *views.Application, style themes.Theme) *UI {
	pg := widgets.NewPager(app)

	sb := widgets.NewStatusbar(app, style)
	sb.SetStatus("START-UP")

	u := &UI{
		app:       app,
		statusbar: sb,

		autoscroll: true,
		pager:      pg,
	}

	u.SetOrientation(views.Vertical)
	u.AddWidget(pg, 1)
	u.AddWidget(sb, 0)
	return u
}

func (u *UI) HandleEvent(e tcell.Event) bool {
	switch te := e.(type) {
	case *tcell.EventKey:
		return u.handleEventKey(te)
	}
	return false
}

func (u *UI) handleEventKey(ek *tcell.EventKey) bool {
	switch ek.Key() {
	case tcell.KeyCtrlC:
		u.app.Quit()
		return true
	case tcell.KeyRune:
		switch ek.Rune() {
		case 'q':
			u.app.Quit()
			return true
		}
	}
	return false
}

func (u *UI) PagerAppend(line string) {
	u.pager.Append(line)
	if u.autoscroll {
		u.pager.ScrollToEnd()
	}

	pct := u.pager.GetScrollPercentage()
	u.statusbar.SetScrollPercentage(int(pct * 100))
}

func (u *UI) PagerLen() int {
	return u.pager.Len()
}

func (u *UI) SetMessage(s string) {
	u.statusbar.SetMessage(s)
}

func (u *UI) SetStatus(s string) {
	u.statusbar.SetStatus(s)
}
