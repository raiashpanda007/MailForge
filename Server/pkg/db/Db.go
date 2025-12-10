package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DataBase struct {
	Db    *pgxpool.Pool
	Redis *redis.Client
}

func Db_Init(DbURL string, RedisURL string) (*DataBase, error) {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, DbURL)
	if err != nil {
		slog.Error("UNABLE TO CONNECT TO THE DATABASE ")
		return nil, err
	}

	Redis := redis.NewClient(&redis.Options{
		Addr: RedisURL,
	})

	err = db.Ping(ctx)
	if err != nil {
		slog.Error("DB Ping failed ", slog.Any("ERROR :: ", err))
		return nil, err
	}
	_, err = Redis.Ping(ctx).Result()
	if err != nil {
		slog.Error("REDIS Ping failed ", slog.Any("ERROR :: ", err))
		return nil, err
	}

	slog.Info("CONNECTED TO DB POSTGRES AND REDIS")

	return &DataBase{Db: db, Redis: Redis}, err
}
