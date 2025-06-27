package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucasfaria/rbac/api-go/internal/auth"
	db "github.com/lucasfaria/rbac/api-go/internal/database/sqlc"
	"github.com/lucasfaria/rbac/api-go/internal/handler"
	authMiddleware "github.com/lucasfaria/rbac/api-go/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(queries *db.Queries, dbConn *sql.DB, jwtSecret, githubClientID, githubClientSecret, githubRedirectURI string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	
	r.Get("/coverage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./coverage.html")
	})

	// Initialize auth service and middleware
	authService := auth.NewAuthService(dbConn, jwtSecret, githubClientID, githubClientSecret, githubRedirectURI)
	authMW := authMiddleware.NewAuthMiddleware(authService)
	authHandler := handler.NewAuthHandler(authService)

	r.Route("/v1", func(r chi.Router) {
		// Public routes
		r.Get("/health", handler.HealthCheck)

		// Authentication routes (public)
		authHandler.RegisterRoutes(r)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireAuth)
			r.Get("/profile", handler.GetProfile(queries))
		})
	})

	return r
}
