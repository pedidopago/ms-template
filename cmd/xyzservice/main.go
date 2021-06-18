package main

import (
	_ "embed"

	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app"
	_ "github.com/pedidopago/ms-template/internal/database" // MariaDB migrations
	_ "github.com/pedidopago/ms-template/internal/xyzservice"
)

func main() {
	app.LoadAndRunAuto()
}
