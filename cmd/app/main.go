package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shar1mo/wishlist-api/internal/auth"
	"github.com/shar1mo/wishlist-api/internal/config"
	"github.com/shar1mo/wishlist-api/internal/handler"
	"github.com/shar1mo/wishlist-api/internal/middleware"
	"github.com/shar1mo/wishlist-api/internal/repository/postgres"
	"github.com/shar1mo/wishlist-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, cfg.DBConnString())
	if err != nil {
		log.Fatalf("create db pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	wishlistRepo := postgres.NewWishlistRepository(db)
	itemRepo := postgres.NewItemRepository(db)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL)

	authService := service.NewAuthService(userRepo, jwtManager)
	wishlistService := service.NewWishlistService(wishlistRepo)
	itemService := service.NewItemService(itemRepo)

	authHandler := handler.NewAuthHandler(authService)
	wishlistHandler := handler.NewWishlistHandler(wishlistService, middleware.GetUserIDFromContext)
	itemHandler := handler.NewItemHandler(itemService, middleware.GetUserIDFromContext)

	router := chi.NewRouter()
	router.Get("/health", healthHandler)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))

			r.Route("/wishlists", func(r chi.Router) {
				r.Post("/", wishlistHandler.Create)
				r.Get("/", wishlistHandler.List)
				r.Get("/{wishlistId}", wishlistHandler.GetByID)
				r.Put("/{wishlistId}", wishlistHandler.Update)
				r.Delete("/{wishlistId}", wishlistHandler.Delete)

				r.Post("/{wishlistId}/items", itemHandler.Create)
				r.Get("/{wishlistId}/items", itemHandler.List)
				r.Get("/{wishlistId}/items/{itemId}", itemHandler.GetByID)
				r.Put("/{wishlistId}/items/{itemId}", itemHandler.Update)
				r.Delete("/{wishlistId}/items/{itemId}", itemHandler.Delete)
			})
		})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownCtx.Done()

		log.Println("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("wishlist-api started on :%s", cfg.ServerPort)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}