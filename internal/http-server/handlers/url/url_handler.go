package urlhandler

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"urlShortener/internal/http-server/api/response"
	"urlShortener/internal/lib/urlgen"
	"urlShortener/internal/storage"
	"urlShortener/internal/storage/dbqueries"
	db "urlShortener/internal/storage/postgres"
)

const (
	urlLen = 6
)

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Url   string `json:"url" validate:"required,url"`
		Alias string `json:"alias,omitempty"`
	}

	var req request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		render.JSON(w, r, response.Error("failed to decode request body"))

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

	url, err := db.Queries.InsertUrl(context.Background(), dbqueries.InsertUrlParams{
		Alias: alias,
		Url:   req.Url,
	})
	if errors.Is(err, storage.ErrUrlExists) {
		render.JSON(w, r, response.Error("url alias already exists"))

		return
	} else if err != nil {
		render.JSON(w, r, response.Error("failed to create url"))

		return
	}

	render.JSON(w, r, Response{
		Response: response.Success(),
		Alias:    url.Alias,
	})
}

func Update(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Alias string `json:"alias,omitempty"`
	}

	var req request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		render.JSON(w, r, response.Error("failed to decode request body"))

		return
	}

	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		render.JSON(w, r, response.ValidationError(validateErr))

		return
	}

	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseInt(idStr, 0, 32)
	if err != nil {
		render.JSON(w, r, response.Error("bad id format"))

		return
	}
	_, err = db.Queries.UpdateUrl(context.Background(), dbqueries.UpdateUrlParams{
		ID:    int32(id64),
		Alias: sql.NullString{String: req.Alias, Valid: req.Alias != ""},
	})
	if errors.Is(err, storage.ErrUrlNotFound) {
		render.JSON(w, r, response.Error("url not found"))

		return
	} else if errors.Is(err, storage.ErrUrlExists) {
		render.JSON(w, r, response.Error("url already exists"))

		return
	} else if err != nil {
		render.JSON(w, r, response.Error("failed to get url"))

		return
	}

	render.JSON(w, r, Response{
		Response: response.Success(),
	})
}

func Get(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	url, err := db.Queries.GetUrlByAlias(context.Background(), alias)
	if errors.Is(err, storage.ErrUrlNotFound) {
		render.JSON(w, r, response.Error("url not found"))

		return
	} else if err != nil {
		render.JSON(w, r, response.Error("failed to get url"))
	}

	_, err = db.Queries.UpdateUrl(context.Background(), dbqueries.UpdateUrlParams{
		ID:    url.ID,
		Count: sql.NullInt32{Int32: url.Count + 1, Valid: true},
	})

	http.Redirect(w, r, url.Url, http.StatusFound)
}
