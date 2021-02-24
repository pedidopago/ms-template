package xyzservice

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pedidopago/ms-grpcd/pkg/grpcd"
	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app"
	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app/config"
)

type Service struct {
	app.App
	cfg     config.Config
	db      func() *sqlx.DB
	redis   func() *redis.Client
	grpcdcl func() grpcd.Client
}

func (s *Service) Start(c app.ServiceContext) error {
	s.cfg = c.Config()
	s.db = func() *sqlx.DB { return s.MariaDB(s.cfg) }
	s.redis = func() *redis.Client { return s.Redis(s.cfg) }
	s.grpcdcl = func() grpcd.Client { return s.GRPCDClient(s.cfg) }

	if c.Config().AutoMigration() {
		if err := app.MariaDBMigrateUp(c, s.db(), s.cfg.MariaDBCS()); err != nil {
			return err
		}
	}

	// FIXME: lines below
	// xyzpb.RegisterXYZServiceServer(c.Server(), s)
	// go xyzpb.ServerPublish(c, s.cfg)

	// FIXME: remove error below
	return errors.New("FIXME: lines 34 to 40")
	//return nil
}

func (s *Service) Stop(c context.Context) error {
	s.App.Stop()
	return nil
}

func init() {
	app.Register("xyzservice", func() app.Service {
		return &Service{}
	})
}
