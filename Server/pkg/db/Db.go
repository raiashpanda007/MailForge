package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type DataBase struct {
	Db *pgxpool.Pool
}

func Db_Init(URL string) (*DataBase, error) {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, URL)
	if err != nil {
		slog.Error("UNABLE TO CONNECT TO THE DATABASE ")
		return nil, err
	}

	err = db.Ping(ctx)

	if err != nil {
		slog.Error("DB Ping failed ", slog.Any("ERROR :: ", err))
		return nil, err
	}

	slog.Info("Connect to db Postgresql")

	return &DataBase{Db: db}, err
}
