package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"uniswap-campus-marketplace/config"
	"uniswap-campus-marketplace/handlers"
	"uniswap-campus-marketplace/repository"
	"uniswap-campus-marketplace/services"

	_ "github.com/lib/pq"
)

type app struct {
	cfg *config.Config
	db  *sql.DB
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := config.OpenDB(ctx, cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	a := &app{cfg: cfg, db: db}
	userRepo := repository.NewPostgresUserRepository(db)
	listingRepo := repository.NewPostgresListingRepository(db)
	reportRepo := repository.NewPostgresReportRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	listingService := services.NewListingService(listingRepo)
	reportService := services.NewReportService(reportRepo, listingRepo)

	authHandler := handlers.NewAuthHandler(authService)
	listingHandler := handlers.NewListingHandler(listingService, reportService)
	uploadHandler := handlers.NewUploadHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.healthCheck)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/listings", listingHandler.Listings)
	mux.HandleFunc("/api/listings/", listingHandler.ListingByIDRoutes)
	mux.HandleFunc("/api/uploads/image", uploadHandler.UploadImage)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("server running on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func (a *app) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiError("method not allowed"))
		return
	}

	writeJSON(w, http.StatusOK, apiSuccess(map[string]string{
		"status": "ok",
	}))
}

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func apiSuccess(data interface{}) apiResponse {
	return apiResponse{Success: true, Data: data}
}

func apiError(message string) apiResponse {
	return apiResponse{Success: false, Error: message}
}

func writeJSON(w http.ResponseWriter, status int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
