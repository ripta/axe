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
		spool, err = ioutil.TempFile("", "axe-*.spool")
		if err != nil {
			return err
		}
		defer func() {
			spool.Close()
			fmt.Printf("Spool: %+v\n", spool.Name())
		}()
	}

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
	su := time.Tick(5 * time.Second)

	a.l.Printf("starting UI")
	go func() {
		for {
			select {
			case line := <-a.LogManager.Logs():
				rate.Add(len(line.Text))
				msg := line.Name + "» " + line.Text
				// a.l.Printf("%s/%s: %s", line.Namespace, line.Name, string(line.Bytes))
				a.UI.PagerAppend(msg)
				if spool != nil {
					spool.WriteString(msg + "\n")
				}
				_ = line
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
