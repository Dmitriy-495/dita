package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apihttp "github.com/Dmitriy-495/dita/internal/api/http"
	"github.com/Dmitriy-495/dita/internal/meta"
	"github.com/Dmitriy-495/dita/internal/runtime"
	"github.com/Dmitriy-495/dita/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DITA_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DITA_DATABASE_URL is required")
	}

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	migrationsDir := os.Getenv("DITA_MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	if err := postgres.ApplyMigrations(ctx, db, migrationsDir); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	catalogRepository := postgres.NewCatalogRepository(db)
	catalogService := meta.NewCatalogService(catalogRepository)

	elementRepository := postgres.NewElementRepository(db)
	runtimeService := runtime.NewCatalogService(
		catalogService,
		elementRepository,
	)

	metaHandler := apihttp.NewMetaHandler(catalogService)
	runtimeHandler := apihttp.NewRuntimeHandler(runtimeService)
	handler := apihttp.NewHandler(metaHandler, runtimeHandler)

	addr := os.Getenv("DITA_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("DITA HTTP server listening on %s", addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server: %v", err)
		}
	}()

	<-serverCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown: %v", err)
	}
}
