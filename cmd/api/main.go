package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/danielllmuniz/devices-api/internal/api"
	"github.com/danielllmuniz/devices-api/internal/services"
	"github.com/danielllmuniz/devices-api/internal/store/pgstore"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// LOAD ENVIRONMENT VARIABLES
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// LOAD CONTEXT
	ctx := context.Background()

	// DATABASE CONNECTION
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	// START SERVER
	api := api.Api{
		Router:        chi.NewMux(),
		DeviceService: *services.NewDeviceService(pgstore.NewPGDeviceStore(pool)),
	}

	api.BindRoutes()

	if err := http.ListenAndServe(os.Getenv("API_PORT"), api.Router); err != nil {
		panic(err)
	}

}
