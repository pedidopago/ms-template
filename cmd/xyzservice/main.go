package main

import (
	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app"
	_ "github.com/pedidopago/ms-template/internal/database/statik" // MariaDB migrations
	_ "github.com/pedidopago/ms-template/internal/xyzservice"
)

func main() {
	app.LoadAndRunAuto()
}
