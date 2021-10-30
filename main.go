package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"weekly-radio-programme/config"
	api2 "weekly-radio-programme/internal/api"
	"weekly-radio-programme/internal/show"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	conn, err := pgx.Connect(ctx, cfg.DbDsn)
	if err != nil {
		panic(err)
	}

	repo := show.NewPostgresRepo(conn)
	srv := show.NewService(repo)
	api := api2.New(srv, cfg.Port)

	if err := api.Run(); err != nil {
		panic(err)
	}
}
