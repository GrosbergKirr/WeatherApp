package storage

import (
	//"database/sql"
	//"database/sql"
	"fmt"
	"log/slog"

	//_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	Db *sqlx.DB
}

func InitStorage(log *slog.Logger, user, pass, addr, name, mode string) *Storage {
	dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, pass, addr, name, mode)
	db, err := sqlx.Connect("postgres", dbPath)
	if err != nil {
		log.Error("Failed to initialize storage", slog.Any("err", err))
		//panic(err)

	}
	sqlDB := db.DB
	err = goose.Up(sqlDB, "migrations")
	if err != nil {
		log.Error("Migration error: %e", slog.Any("err", err))
		return nil
	}
	log.Info("storage initialized")
	return &Storage{Db: db}
}
