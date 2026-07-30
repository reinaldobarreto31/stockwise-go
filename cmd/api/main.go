package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/reinaldobarreto31/stockwise-go/internal/db"
	"github.com/reinaldobarreto31/stockwise-go/internal/handler"
	"github.com/reinaldobarreto31/stockwise-go/internal/middleware"
	"github.com/reinaldobarreto31/stockwise-go/internal/repository"
	"github.com/reinaldobarreto31/stockwise-go/internal/service"
)

func main() {
	// Load .env in development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()
	log.Println("Connected to PostgreSQL")

	// Repositories
	userRepo := repository.NewUserRepository(database)
	productRepo := repository.NewProductRepository(database)
	movementRepo := repository.NewMovementRepository(database)

	// Services
	authSvc := service.NewAuthService(userRepo)
	productSvc := service.NewProductService(productRepo)
	movementSvc := service.NewMovementService(movementRepo, productRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	productHandler := handler.NewProductHandler(productSvc)
	movementHandler := handler.NewMovementHandler(movementSvc)

	// Router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.CORS)

	// Health check (public)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"stockwise-go"}`)
	})

	// API v1
	r.Route("/api", func(r chi.Router) {
		// Public info
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"StockWise API","version":"1.0.0"}`)
		})

		// Auth routes — public
		r.Mount("/auth", authHandler.Router())

		// Protected routes — require valid JWT
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)
			r.Mount("/products", productHandler.Router())
			r.Mount("/movements", movementHandler.Router())
		})
	})

	addr := ":" + port
	log.Printf("StockWise API listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
