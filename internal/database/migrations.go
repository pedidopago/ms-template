package database

import (
	"embed"
	"io/fs"

	"github.com/pedidopago/ms-grpcd/pkg/grpcd/app"
)

var (
	//go:embed migrations/*
	migrations embed.FS
)

// Migrations embedded files
func Migrations() fs.FS {
	fss, err := fs.Sub(migrations, "migrations")
	if err != nil {
		panic(err)
	}
	return fss
}

func init() {
	app.RegisterMigrations(Migrations())
}
