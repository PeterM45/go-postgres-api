package router

import (
	"net/http"
	"strings"

	"github.com/PeterM45/go-postgres-api/internal/auth"
	"github.com/PeterM45/go-postgres-api/internal/handler"
	"github.com/PeterM45/go-postgres-api/internal/middleware"
)

type Router struct {
	userHandler *handler.UserHandler
	jwt         *auth.JWT
}

func New(userHandler *handler.UserHandler, jwt *auth.JWT) *Router {
	return &Router{
		userHandler: userHandler,
		jwt:         jwt,
	}
}

func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/auth/login", r.userHandler.HandleLogin)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("/api/users", r.handleUsers)
	protected.HandleFunc("/api/users/", r.handleUserByID) // Will handle /api/users/[id]/*

	// Apply auth middleware to protected routes
	mux.Handle("/api/users", middleware.AuthMiddleware(r.jwt)(protected))
	mux.Handle("/api/users/", middleware.AuthMiddleware(r.jwt)(protected))

	return mux
}

func (r *Router) handleUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.userHandler.GetUsers(w, req)
	case http.MethodPost:
		r.userHandler.CreateUser(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleUserByID(w http.ResponseWriter, req *http.Request) {
	// Extract ID from path: /api/users/[id]
	path := strings.TrimPrefix(req.URL.Path, "/api/users/")
	id := strings.Split(path, "/")[0]

	switch req.Method {
	case http.MethodGet:
		r.userHandler.GetUser(w, req, id)
	case http.MethodPut:
		r.userHandler.UpdateUser(w, req, id)
	case http.MethodDelete:
		r.userHandler.DeleteUser(w, req, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
