package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type GetCookieParams struct {
	Name string `path:"name"`
}

func (server *Server) SetCookie(ctx context.Context, input *struct{}) (*Response, error) {
	resp, err := server.GetSuccessResponseData(Payload{
		Message: "Hello!",
	})

	if err != nil {
		message := http.StatusText(http.StatusInternalServerError)
		return nil, huma.Error500InternalServerError(message, errors.New(message))
	}

	return &resp, nil
}

func (server *Server) GetCookie(ctx context.Context, input *GetCookieParams) (*Response, error) {
	resp, err := server.GetSuccessResponseData(Payload{
		Message: fmt.Sprintf("Hello, %s!", input.Name),
	})

	if err != nil {
		message := http.StatusText(http.StatusInternalServerError)
		return nil, huma.Error500InternalServerError(message, errors.New(message))
	}

	return &resp, nil
}
