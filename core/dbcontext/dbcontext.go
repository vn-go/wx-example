package dbcontext

import (
	"core/config"

	"github.com/vn-go/dx"
)

func OpneDb() (*dx.DB, error) {
	
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		return nil, err
	}

	return dx.Open(cfg.Database.Driver, cfg.Database.DSN)
}
