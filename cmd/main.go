package main

import (
	"log"
	"net/http"

	"github.com/PeterM45/go-postgres-api/internal/auth"
	"github.com/PeterM45/go-postgres-api/internal/config"
	"github.com/PeterM45/go-postgres-api/internal/database"
	"github.com/PeterM45/go-postgres-api/internal/handler"
	"github.com/PeterM45/go-postgres-api/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jwtAuth := auth.NewJWT(cfg.JWT.SecretKey)
	userHandler := handler.NewUserHandler(db, jwtAuth)

	// Setup router
	router := router.New(userHandler, jwtAuth)

	log.Printf("Server starting on http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router.Setup()); err != nil {
		log.Fatal(err)
	}
}
