package get

import (
	"errors"
	"log/slog"
	"net/http"
	resp "sptringTresRestAPI/internal/lib/api/response"
	"sptringTresRestAPI/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	URL string `json:"url"`
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func Get(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.Get"
		log = slog.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request", err)

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		url, err := urlGetter.GetURL(req.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("no saved url for the alias", req.Alias)
				render.JSON(w, r, resp.Error("url not found"))
				return
			}
			log.Error("failed to get url", err)

			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}

		log.Info("url was found", req.Alias, url)

		render.JSON(w, r, Response{
			Response: resp.OK(),
			URL:      url,
		})
		return
	}
}
