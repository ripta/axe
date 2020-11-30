package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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

	debug bool
	l     logger.Interface
}

func New(l logger.Interface, m *kubelogs.Manager, debug bool) *App {
	style := themes.SolarizedDark()
	app := &views.Application{}

	u := ui.New(app, style)
	app.SetRootWidget(u)

	return &App{
		App:        app,
		UI:         u,
		LogManager: m,

		debug: debug,
		l:     l,
	}
}

func (a *App) Run(ctx context.Context) error {
	var spool *os.File
	var err error

	if a.debug {
		spool, err = ioutil.TempFile("", "axe-*.log")
		if err != nil {
			return err
		}
		defer func() {
			spool.Close()
			fmt.Printf("Log file: %+v\n", spool.Name())
		}()
	}

	a.App.Start()

	go func() {
		a.App.PostFunc(func() {
			a.UI.SetStatus("SYNCING", themes.AltTypeNew)
		})
		// a.l.Printf("starting manager")
		if err := a.LogManager.Run(ctx); err != nil {
			a.App.PostFunc(func() {
				a.UI.SetStatus("ERROR", themes.AltTypeError)
				a.UI.SetMessage(fmt.Sprintf("Log manager reported: %+v", err))
			})
		} else {
			a.App.PostFunc(func() {
				a.UI.SetStatus("TAILING", themes.AltTypeOK)
			})
		}
	}()

	rate := iorate.New()
	su := time.Tick(5 * time.Second)

	a.l.Printf("starting UI")
	go func() {
		for {
			select {
			case line := <-a.LogManager.Logs():
				msg := line.Name + "] " + line.Text

				rate.Add(len(msg))
				a.App.PostFunc(func() {
					a.UI.PagerAppend(msg)
					if spool != nil {
						spool.WriteString(msg + "\n")
					}
				})
			case <-su:
				activeCnt, allCnt := a.LogManager.ContainerCount()
				r := iorate.HumanizeBytes(rate.Calculate(time.Second))
				a.App.PostFunc(func() {
					l := iorate.HumanizeBytes(float64(a.UI.PagerLen()))
					a.UI.SetMessage(fmt.Sprintf("%d/%d containers | %s (%s/s)", activeCnt, allCnt, l, r))
				})
			case <-ctx.Done():
				break
			}
		}
	}()

	return a.App.Wait()
}
