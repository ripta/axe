package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/ripta/axe/pkg/iorate"
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

	rate := iorate.New()
	su := time.Tick(time.Second)

	u.l.Printf("starting UI")
	for {
		select {
		case line := <-u.Manager.Logs():
			rate.Add(len(line.Bytes))
			// u.l.Printf("%s/%s: %s", line.Namespace, line.Name, string(line.Bytes))
			_ = line
		case <-su:
			fmt.Printf("\rRate: %s/s\033[0K", iorate.HumanizeBytes(rate.Calculate(time.Second)))
		case <-ctx.Done():
			return nil
		}
	}
}
