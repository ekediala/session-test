package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/securecookie"
)

type Server struct {
	router       *chi.Mux
	sessionKeys  []securecookie.Codec
	dashboardAPI huma.API
	apiApi       huma.API
}

func (server *Server) registerRoutes() {
	server.registerExternalApiRoutes()
	server.registerDashboardApiRoutes()
}

func (server *Server) start() error {
	keys := []string{
		os.Getenv("SESSION_SECRET_ONE"),
		os.Getenv("SESSION_SECRET_TWO"),
	}

	for _, secret := range keys {
		if secret == "" {
			return errors.New("Missing environment variable")
		}
	}

	codec1 := securecookie.New([]byte(keys[0]), securecookie.GenerateRandomKey(32))
	codec2 := securecookie.New([]byte(keys[1]), securecookie.GenerateRandomKey(32))

	server.sessionKeys = []securecookie.Codec{codec1, codec2}

	server.router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	server.registerRoutes()
	return http.ListenAndServe("localhost:8888", server.router)
}

func New(router *chi.Mux) *Server {
	return &Server{
		router: router,
	}
}
