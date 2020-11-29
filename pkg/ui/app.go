package ui

import (
	"github.com/gcla/gowid"
	"github.com/sirupsen/logrus"

	"github.com/ripta/axe/pkg/ui/palettes"
)

func buildFilterWidget()

func buildUI() (*gowid.App, error) {
	return gowid.NewApp(gowid.AppArgs{
		View:         view,
		Palette:      palettes.SolarizedDark(),
		DontActivate: true,
		Log:          logrus.StandardLogger(),
	})
}
