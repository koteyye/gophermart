package gophermart

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
)

// Run временный, для тестового запуска
func Run(args []string) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	ctx := context.Background()

	db, err := newConnection(ctx, config.DSN)
	if err != nil {
		return err
	}

	// Инициализация storage
	storage := storage.NewStorage(db)
	fmt.Println(storage)
	// storage в инициализацию service...


	//для тестового запуска
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {fmt.Print("test")})
	http.ListenAndServe(config.Address, nil)

	return nil
}


func newConnection(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect db: %w", err)
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	sql, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open db with driver: %w", err)
	}

	if err := migrations.Up(ctx, sql); err != nil {
		return nil, fmt.Errorf("can't migration: %w", err)
	}

	return conn, nil
}