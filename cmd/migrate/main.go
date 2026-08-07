package main

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"backend-porto/internal/shared/config"

	// pgx's database/sql driver (goose needs a *sql.DB; pgx stays under the hood)
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate <up|down|status|version>")
	}

	config.LoadDotEnv(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required (set it in .env or the environment)")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set dialect: %v", err)
	}

	var runErr error
	switch os.Args[1] {
	case "up":
		runErr = goose.Up(db, "migrations")
	case "down":
		runErr = goose.Down(db, "migrations")
	case "status":
		runErr = goose.Status(db, "migrations")
	case "version":
		runErr = goose.Version(db, "migrations")
	default:
		log.Fatalf("unknown command %q (use up, down, status, version)", os.Args[1])
	}
	if runErr != nil {
		log.Fatalf("migrate %s: %v", os.Args[1], runErr)
	}
}
