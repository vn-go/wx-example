package base

import (
	"core"
)

type Base struct {
	Svc *core.Service
	// db of current tenant

}

func (b *Base) New() error {
	b.Svc = core.Services
	return nil
}
