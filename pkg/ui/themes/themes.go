package themes

import (
	"github.com/gdamore/tcell/v2"
)

type Alts struct {
	Error   tcell.Style
	Expired tcell.Style
	New     tcell.Style
	Normal  tcell.Style
	OK      tcell.Style
}

type Theme struct {
	Body      tcell.Style
	Title     tcell.Style
	Statusbar Alts
}
