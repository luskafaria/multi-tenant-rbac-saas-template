package main

import (
	"log"
	"net/http"

	_ "github.com/lucasfaria/rbac/api-go/docs"
	"github.com/lucasfaria/rbac/api-go/internal/config"
	"github.com/lucasfaria/rbac/api-go/internal/database"
	"github.com/lucasfaria/rbac/api-go/internal/router"
)

// @title RBAC API
// @version 1.0
// @description This is the API for the RBAC service.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3334
// @BasePath /v1
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	dbConn, queries, err := database.NewDB(cfg.DB_SOURCE)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer dbConn.Close()

	appRouter := router.NewRouter(queries, dbConn, cfg.JWT_SECRET, cfg.GITHUB_OAUTH_CLIENT_ID, cfg.GITHUB_OAUTH_CLIENT_SECRET, cfg.GITHUB_OAUTH_CLIENT_REDIRECT_URI)

	log.Printf("server is running on port %s", cfg.API_PORT)
	err = http.ListenAndServe(":"+cfg.API_PORT, appRouter)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
