package themes

import (
	"github.com/gdamore/tcell/v2"
)

type AltType string

var (
	AltTypeError   = AltType("Error")
	AltTypeExpired = AltType("Expired")
	AltTypeNew     = AltType("New")
	AltTypeNormal  = AltType("Normal")
	AltTypeOK      = AltType("OK")
)

type Alts struct {
	Error   tcell.Style
	Expired tcell.Style
	New     tcell.Style
	Normal  tcell.Style
	OK      tcell.Style
}

func (a Alts) Select(t AltType) tcell.Style {
	switch t {
	case AltTypeError:
		return a.Error
	case AltTypeExpired:
		return a.Expired
	case AltTypeNew:
		return a.New
	case AltTypeOK:
		return a.OK
	default:
		return a.Normal
	}
}

type Theme struct {
	Body      tcell.Style
	Title     tcell.Style
	Statusbar Alts
}
