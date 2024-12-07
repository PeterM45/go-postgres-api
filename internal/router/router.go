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

	// Public routes (no auth middleware)
	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/api/auth/login", r.userHandler.HandleLogin)
	publicMux.HandleFunc("/api/users", func(w http.ResponseWriter, req *http.Request) {
		// Only allow POST for public access (user creation)
		if req.Method == http.MethodPost {
			r.userHandler.CreateUser(w, req)
			return
		}
		// All other methods need auth
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/users", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/users" {
			http.NotFound(w, req)
			return
		}
		r.handleUsers(w, req)
	})
	protectedMux.HandleFunc("/api/users/", r.handleUserByID)

	// Combine public and protected routes
	mux.Handle("/api/auth/", publicMux)
	mux.Handle("/api/users", publicMux)
	mux.Handle("/api/users/", middleware.AuthMiddleware(r.jwt)(protectedMux))

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

	if id == "" {
		http.NotFound(w, req)
		return
	}

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
