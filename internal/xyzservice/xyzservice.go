package xyzservice

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app"
	"github.com/urfave/cli/v2"
)

type Service struct {
	db *sqlx.DB
}

func (s *Service) Start(c app.ServiceContext) error {
	var err error

	s.db, err = app.MariaDB(c.MariaDBCS())
	if err != nil {
		return err
	}

	if c.AutoMigrate() {
		if err := app.MariaDBMigrateUp(c, s.db, c.MariaDBCS()); err != nil {
			return err
		}
	}

	// TODO: register rpc service
	return nil
}

func (s *Service) Stop(c context.Context) error {
	if s.db != nil {
		_ = s.db.Close()
	}
	return nil
}

func (s *Service) BuildExtraSettings(c *cli.Context) interface{} {
	// here you can initialize a custom struct to use in Start()
	return nil
}
