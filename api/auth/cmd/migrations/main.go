package main

import (
	"TeamProject/api/auth/internal/config"
	"flag"

	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationPath string

	flag.StringVar(&migrationPath, "migrationPath", "", "path to the migrations directory")

	cfg := config.MustConfig()

	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DB_name, cfg.DB.SSL_mode,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		panic(err)
	}
}
