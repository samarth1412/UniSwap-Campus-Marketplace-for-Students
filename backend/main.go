package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"uniswap-campus-marketplace/apiresponse"
	"uniswap-campus-marketplace/config"
	"uniswap-campus-marketplace/handlers"
	"uniswap-campus-marketplace/middleware"
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
	listingImageRepo := repository.NewPostgresListingImageRepository(db)
	wishlistRepo := repository.NewPostgresWishlistRepository(db)
	reportRepo := repository.NewPostgresReportRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	listingService := services.NewListingService(listingRepo)
	listingImageService := services.NewListingImageService(listingRepo, listingImageRepo)
	uploadService := services.NewUploadService(listingRepo, listingImageRepo)
	wishlistService := services.NewWishlistService(wishlistRepo, listingRepo)
	reportService := services.NewReportService(reportRepo, listingRepo)

	authHandler := handlers.NewAuthHandler(authService)
	listingHandler := handlers.NewListingHandler(listingService, reportService)
	uploadHandler := handlers.NewUploadHandler(listingImageService, uploadService)
	userHandler := handlers.NewUserHandler(listingService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.healthCheck)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.Handle("/api/auth/me", middleware.Auth(authService)(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/listings", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.Auth(authService)(http.HandlerFunc(listingHandler.Listings)).ServeHTTP(w, r)
			return
		}
		listingHandler.Listings(w, r)
	}))
	mux.Handle("/api/listings/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)
		if strings.HasSuffix(path, "/images") && r.Method == http.MethodPost {
			middleware.Auth(authService)(http.HandlerFunc(uploadHandler.UploadListingImages)).ServeHTTP(w, r)
			return
		}
		if (strings.HasSuffix(path, "/report") && r.Method == http.MethodPost) || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			middleware.Auth(authService)(http.HandlerFunc(listingHandler.ListingByIDRoutes)).ServeHTTP(w, r)
			return
		}
		listingHandler.ListingByIDRoutes(w, r)
	}))
	mux.Handle("/api/uploads/image", middleware.Auth(authService)(http.HandlerFunc(uploadHandler.UploadImage)))
	mux.Handle("/api/users/", http.HandlerFunc(userHandler.UserRoutes))
	mux.Handle("/api/wishlist", middleware.Auth(authService)(http.HandlerFunc(wishlistHandler.Wishlist)))
	mux.Handle("/api/wishlist/", middleware.Auth(authService)(http.HandlerFunc(wishlistHandler.WishlistByID)))
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
		apiresponse.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	apiresponse.WriteSuccess(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
