package widgets

import (
	"fmt"

	"github.com/gdamore/tcell/v2/views"
	"github.com/ripta/axe/pkg/ui/themes"
)

type Statusbar struct {
	views.BoxLayout

	app    *views.Application
	styles themes.Alts

	status  *views.Text
	message *views.Text
	scroll  *views.Text
}

func NewStatusbar(app *views.Application, style themes.Theme) *Statusbar {
	status := views.NewText()
	status.SetStyle(style.Statusbar.New)

	message := views.NewText()
	message.SetStyle(style.Statusbar.New)

	scroll := views.NewText()
	scroll.SetStyle(style.Statusbar.New)

	bar := &Statusbar{
		app:    app,
		styles: style.Statusbar,

		status:  status,
		message: message,
		scroll:  scroll,
	}

	bar.AddWidget(status, 0)
	bar.AddWidget(message, 1)
	bar.AddWidget(scroll, 0)
	return bar
}

func (bar *Statusbar) SetMessage(s string) {
	bar.message.SetText(" " + s + " ")
}

func (bar *Statusbar) SetScrollPercentage(pct int) {
	bar.scroll.SetText(fmt.Sprintf(" %d%% ", pct))
}

func (bar *Statusbar) SetStatus(s string) {
	if len(s) > 8 {
		s = s[:8]
	}
	bar.status.SetText(fmt.Sprintf(" %-8s ", s))
	bar.status.SetStyle(bar.styles.Normal)
}
