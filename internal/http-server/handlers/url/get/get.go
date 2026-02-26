package get

import (
	"errors"
	"log/slog"
	"net/http"
	resp "sptringTresRestAPI/internal/lib/api/response"
	"sptringTresRestAPI/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

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
		log := slog.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("url not found"))
				return
			}
			log.Error(
				"failed to get url",
				slog.String("error", err.Error()),
				slog.String("alias", alias),
			)
			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}

		log.Info("url found", slog.String(alias, url))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			URL:      url,
		})
		return
	}
}
