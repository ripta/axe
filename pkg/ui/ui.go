package ui

import (
	"context"

	"github.com/ripta/axe/pkg/kubelogs"
	"github.com/ripta/axe/pkg/logger"
)

type UI struct {
	*kubelogs.Manager
	l logger.Interface
}

func New(l logger.Interface, m *kubelogs.Manager) *UI {
	return &UI{
		Manager: m,
		l:       l,
	}
}

func (u *UI) Run(ctx context.Context) error {
	u.l.Printf("starting manager")
	if err := u.Manager.Run(ctx); err != nil {
		return err
	}

	u.l.Printf("starting UI")
	for {
		select {
		case line := <-u.Manager.Logs():
			// u.l.Printf("%s/%s: %s", line.Namespace, line.Name, string(line.Bytes))
			_ = line
		case <-ctx.Done():
			return nil
		}
	}
}
