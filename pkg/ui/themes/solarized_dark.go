package themes

import "github.com/gdamore/tcell/v2"

var (
	solarizedBlack   = tcell.GetColor("#002b36")
	solarizedBlue    = tcell.GetColor("#268bd2")
	solarizedCyan    = tcell.GetColor("#2aa198")
	solarizedGray1   = tcell.GetColor("#586e75")
	solarizedGray2   = tcell.GetColor("#657b83")
	solarizedGray3   = tcell.GetColor("#839496")
	solarizedGray4   = tcell.GetColor("#93a1a1")
	solarizedGreen   = tcell.GetColor("#859900")
	solarizedMagenta = tcell.GetColor("#d33682")
	solarizedOrange  = tcell.GetColor("#cb4b16")
	solarizedPurple  = tcell.GetColor("#6c71c4")
	solarizedRed     = tcell.GetColor("#dc322f")
	solarizedWhite   = tcell.GetColor("#fdf6e3")
	solarizedYellow  = tcell.GetColor("#B58900")
	solarizedZeta    = tcell.GetColor("#79e11a")
)

var (
	solarizedForeground = solarizedGray3
	solarizedBackground = solarizedBlack
	solarizedSelected   = tcell.GetColor("#073642")
)

func SolarizedDark() Theme {
	base := tcell.StyleDefault.Background(solarizedBackground).Foreground(solarizedForeground)
	return Theme{
		Body: base,
		Statusbar: Alts{
			Error:   base.Foreground(solarizedRed),
			Expired: base.Foreground(solarizedRed).Reverse(true),
			New:     base.Foreground(solarizedCyan),
			Normal:  base,
			OK:      base.Foreground(solarizedGreen),
		},
		Title: base.Background(solarizedSelected),
	}
}
