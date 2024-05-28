package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gorilla/securecookie"
)

var (
	ErrInvalidCookie   error = errors.New("Session expired. Please login and try again.")
	ErrForbiddenAccess error = errors.New("You do not have permission to acess this resource.")
)

type SessionData struct {
	Value string
}

type Payload struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Response struct {
	Body      Payload     `json:"body"`
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (server *Server) ForbiddenResponse(ctx huma.Context) {
	huma.WriteErr(server.apiApi,
		ctx,
		http.StatusForbidden,
		http.StatusText(http.StatusForbidden),
		ErrForbiddenAccess,
	)
}

func (server *Server) UnauthorisedResponse(ctx huma.Context) {
	huma.WriteErr(server.apiApi,
		ctx,
		http.StatusUnauthorized,
		http.StatusText(http.StatusUnauthorized),
		ErrInvalidCookie,
	)
}

func (server *Server) GetSuccessResponseData(payload Payload) (Response, error) {
	data := SessionData{
		Value: "random_data",
	}
	encoded, err := securecookie.EncodeMulti("session", data, server.sessionKeys...)

	if err != nil {
		return Response{}, err
	}

	cookie := http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
		// Secure:   true,
		HttpOnly: true,
		Domain:   "localhost",
		MaxAge:   300,
		// SameSite: http.SameSiteNoneMode,
	}

	return Response{Body: payload, SetCookie: cookie}, nil
}

func (server *Server) GetUnauthorisedResponse() (Response, error) {
	userId := "random_user"
	encoded, err := securecookie.EncodeMulti("session", userId, server.sessionKeys...)
	if err != nil {
		return Response{}, nil
	}
	SetCookie := http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Domain:   "localhost",
		MaxAge:   -300,
	}

	payload := Payload{
		Message: http.StatusText(http.StatusUnauthorized),
	}

	return Response{payload, SetCookie}, nil
}

func (server *Server) Auth(ctx huma.Context, next func(huma.Context)) {
	cookie, err := huma.ReadCookie(ctx, "session")
	fmt.Println("data", cookie, err)
	if err != nil {
		server.UnauthorisedResponse(ctx)
		return
	}

	if err := cookie.Valid(); err != nil {
		server.UnauthorisedResponse(ctx)
		return
	}

	var value SessionData
	err = securecookie.DecodeMulti("session", cookie.Value, &value, server.sessionKeys...)
	fmt.Println("error", err)
	if err != nil {
		server.ForbiddenResponse(ctx)
		return
	}

	fmt.Println("data", value.Value)
	next(ctx)
}
