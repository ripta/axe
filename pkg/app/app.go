package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2/views"
	"github.com/ripta/axe/pkg/iorate"
	"github.com/ripta/axe/pkg/kubelogs"
	"github.com/ripta/axe/pkg/logger"
	"github.com/ripta/axe/pkg/ui"
	"github.com/ripta/axe/pkg/ui/themes"
)

// App is a controller that connects the LogManager (model) and the UI (view).
type App struct {
	App        *views.Application
	UI         *ui.UI
	LogManager *kubelogs.Manager
	l          logger.Interface
}

func New(l logger.Interface, m *kubelogs.Manager) *App {
	style := themes.SolarizedDark()
	app := &views.Application{}

	u := ui.New(app, style)
	app.SetRootWidget(u)

	return &App{
		App:        app,
		UI:         u,
		LogManager: m,
		l:          l,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.App.Start()

	go func() {
		a.App.PostFunc(func() {
			a.UI.SetStatus("SYNCING")
		})
		// a.l.Printf("starting manager")
		if err := a.LogManager.Run(ctx); err != nil {
			a.App.PostFunc(func() {
				a.UI.SetStatus("ERROR")
				a.UI.SetMessage(fmt.Sprintf("Log manager reported: %+v", err))
			})
		} else {
			a.App.PostFunc(func() {
				a.UI.SetStatus("TAILING")
			})
		}
	}()

	rate := iorate.New()
	su := time.Tick(2 * time.Second)

	a.l.Printf("starting UI")
	go func() {
		for {
			select {
			case line := <-a.LogManager.Logs():
				rate.Add(len(line.Bytes))
				// a.l.Printf("%s/%s: %s", line.Namespace, line.Name, string(line.Bytes))
				_ = line
			case <-su:
				r := iorate.HumanizeBytes(rate.Calculate(time.Second))
				a.App.PostFunc(func() {
					a.UI.SetMessage(fmt.Sprintf("Tailing %d containers | Transferring: %s/s", 1, r))
				})
			case <-ctx.Done():
				break
			}
		}
	}()

	return a.App.Wait()
}
