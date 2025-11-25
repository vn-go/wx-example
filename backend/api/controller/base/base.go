package base

import "core"

type Base struct {
	Svc *core.Service
}

func (b *Base) New() error {
	b.Svc = core.Services
	return nil
}
