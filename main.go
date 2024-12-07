package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	genapi "github.com/h-okay/golang-api-example/gen/api"
	gendb "github.com/h-okay/golang-api-example/gen/db"
	"github.com/h-okay/golang-api-example/internal/handlers"
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5432/userdb"
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Run migrations
	err = runMigrations(dbPool)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	swagger, err := genapi.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(swagger)
	})
	router.Handle("/swagger/", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    "/swagger/",
		SpecURL: "/swagger/doc.json",
	}, nil))

	queries := gendb.New(dbPool)
	userHandlers := handlers.NewUserHandlers(queries)

	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(
		swagger,
		&nethttpmiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
					return nil
				},
			},
		},
	)

	apiServer := genapi.HandlerWithOptions(
		genapi.NewStrictHandler(userHandlers, nil),
		genapi.ChiServerOptions{
			BaseURL:    "/api/v1",
			BaseRouter: router,
			Middlewares: []genapi.MiddlewareFunc{
				validator,
			},
		},
	)

	addr := ":8000"
	httpServer := http.Server{
		Addr:    addr,
		Handler: apiServer,
	}

	log.Println("Server listening on", addr)
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func runMigrations(pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)

	// Create a driver for golang-migrate
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///root/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	// Run the migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully.")
	return nil
}
