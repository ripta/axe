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
	sb.SetStatus("START-UP", themes.AltTypeError)

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
		return u.handleAppEventKeys(te) || u.handleScrollEventKeys(te)
	}
	return false
}

func (u *UI) handleAppEventKeys(ek *tcell.EventKey) bool {
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

func (u *UI) handleScrollEventKeys(ek *tcell.EventKey) bool {
	switch ek.Key() {
	case tcell.KeyCtrlD:
		u.autoscroll = false
		u.pager.ScrollPageDown(1)
		return true
	case tcell.KeyCtrlU:
		u.autoscroll = false
		u.pager.ScrollPageUp(1)
		return true
	case tcell.KeyDown:
		u.autoscroll = false
		u.pager.ScrollDown(1)
		return true
	case tcell.KeyUp:
		u.autoscroll = false
		u.pager.ScrollUp(1)
		return true
	case tcell.KeyRune:
		switch ek.Rune() {
		case 'f':
			u.autoscroll = !u.autoscroll
			return true
		case 'j':
			u.autoscroll = false
			u.pager.ScrollPageDown(1)
			return true
		case 'k':
			u.autoscroll = false
			u.pager.ScrollPageUp(1)
			return true
		}
	}
	return false
}

func (u *UI) PagerAppend(line string) {
	u.pager.Append(line)
	if u.autoscroll {
		u.statusbar.SetStatus("FOLLOW", themes.AltTypeNew)
		u.pager.ScrollToEnd()
	} else {
		u.statusbar.SetStatus("NOFOLLOW", themes.AltTypeNormal)
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

func (u *UI) SetStatus(s string, a themes.AltType) {
	u.statusbar.SetStatus(s, a)
}
