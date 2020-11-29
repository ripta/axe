package palettes

import "github.com/gcla/gowid"

// Colors from https://ethanschoonover.com/solarized/
var (
	solarizedBlack   = gowid.MakeColor("#002b36")
	solarizedBlue    = gowid.MakeColor("#268bd2")
	solarizedCyan    = gowid.MakeColor("#2aa198")
	solarizedGray1   = gowid.MakeColor("#586e75")
	solarizedGray2   = gowid.MakeColor("#657b83")
	solarizedGray3   = gowid.MakeColor("#839496")
	solarizedGray4   = gowid.MakeColor("#93a1a1")
	solarizedGreen   = gowid.MakeColor("#859900")
	solarizedMagenta = gowid.MakeColor("#d33682")
	solarizedOrange  = gowid.MakeColor("#cb4b16")
	solarizedPurple  = gowid.MakeColor("#6c71c4")
	solarizedRed     = gowid.MakeColor("#dc322f")
	solarizedWhite   = gowid.MakeColor("#fdf6e3")
	solarizedYellow  = gowid.MakeColor("#B58900")
	solarizedZeta    = gowid.MakeColor("#79e11a")
)

func SolarizedDark() gowid.Palette {
	return gowid.Palette{
		"default": gowid.MakePaletteEntry(solarizedWhite, solarizedBlack),

		"button":              gowid.MakePaletteEntry(solarizedBlack, solarizedGray3),
		"button-focus":        gowid.MakePaletteEntry(solarizedWhite, solarizedMagenta),
		"button-selected":     gowid.MakePaletteEntry(solarizedWhite, solarizedGray2),
		"dialog":              gowid.MakePaletteEntry(solarizedBlack, solarizedYellow),
		"dialog-button":       gowid.MakePaletteEntry(solarizedYellow, solarizedBlack),
		"filter-intermediate": gowid.MakePaletteEntry(solarizedBlack, solarizedOrange),
		"filter-invalid":      gowid.MakePaletteEntry(solarizedBlack, solarizedRed),
		"filter-menu":         gowid.MakePaletteEntry(solarizedWhite, solarizedBlack),
		"filter-valid":        gowid.MakePaletteEntry(solarizedBlack, solarizedGreen),
		"spinner":             gowid.MakePaletteEntry(solarizedYellow, solarizedBlack),
		"title":               gowid.MakePaletteEntry(solarizedRed, solarizedZeta),
	}
}
