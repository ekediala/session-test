package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func (server *Server) registerExternalApiRoutes() {
	server.router.Route("/api/v1", func(r chi.Router) {
		config := huma.DefaultConfig("My API", "1.0.0")
		config.Servers = []*huma.Server{
			{URL: "https://example.com/api"},
		}
		api := humachi.New(r, config)
		server.apiApi = api

		// Register operations...
		huma.Register(api, huma.Operation{
			Path:        "/set-cookie",
			Method:      http.MethodGet,
			OperationID: "set-cookie",
		}, server.SetCookie)

		huma.Register(api, huma.Operation{
			Path:        "/get-cookie/{name}",
			Method:      http.MethodGet,
			OperationID: "get-cookie",
			Middlewares: huma.Middlewares{
				server.Auth,
			},
		}, server.GetCookie)
	})
}

func (server *Server) registerDashboardApiRoutes() {
	server.router.Route("/ui", func(r chi.Router) {
		config := huma.DefaultConfig("My API", "1.0.0")
		config.Servers = []*huma.Server{
			{URL: "https://example.com/api"},
		}
		api := humachi.New(r, config)
		server.dashboardAPI = api

		// Register operations...
		huma.Get(api, "/demo", func(ctx context.Context, input *struct{}) (*struct{}, error) {
			// TODO: Implement me!
			return nil, nil
		})
	})
}
