package main

import (
	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app"
	_ "github.com/pedidopago/ms-template/internal/database/statik" // MariaDB migrations
	"github.com/pedidopago/ms-template/internal/meta"
	"github.com/pedidopago/ms-template/internal/xyzservice"
)

func main() {
	app.Run(newService,
		app.WithMariaDBCS("testuser:123456789@tcp(localhost)/ms_xyz?parseTime=true"),
		app.WithAutoMigration(),
		app.WithName(meta.ServiceName())) // default tcp listener addr: 'xyzservice:15055'
}

func newService() app.Service {
	return &xyzservice.Service{}
}
