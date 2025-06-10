package create

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"urlShortener/internal/http-server/api/response"
	"urlShortener/internal/http-server/handlers"
	"urlShortener/internal/lib/logger/sl"
	"urlShortener/internal/lib/urlgen"
	"urlShortener/internal/storage"
	"urlShortener/internal/storage/dbqueries"
)

const (
	urlLen = 6
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

func Create(h handlers.Handlers, w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.url.create"

	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		msg := "failed to decode request body"

		h.Logger.Error(msg, sl.Err(err))

		render.JSON(w, r, response.Error(msg))

		return
	}

	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		render.JSON(w, r, response.ValidationError(validateErr))

		return
	}

	alias := req.Alias
	if alias == "" {
		alias = urlgen.New(urlLen)
	}

	url, err := h.Queries.InsertUrl(context.Background(), dbqueries.InsertUrlParams{
		Alias: alias,
		Url:   req.Url,
	})
	if errors.Is(err, storage.ErrUrlExists) {
		msg := "url alias already exists"
		h.Logger.Error(msg, sl.Err(err))
		render.JSON(w, r, response.Error(msg))
	}
	if err != nil {
		msg := "failed to create url"
		h.Logger.Error(msg, sl.Err(err))
		render.JSON(w, r, response.Error(msg))
	}

	render.JSON(w, r, Response{
		Response: response.Success(),
		Alias:    url.Alias,
	})
}
