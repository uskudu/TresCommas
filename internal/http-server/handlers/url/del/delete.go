package del

import (
	"log/slog"
	"net/http"
	resp "sptringTresRestAPI/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Result string `json:"result"`
}

type AliasDeleter interface {
	DeleteAlias(alias string) error
}

func Delete(log *slog.Logger, aliasDeleter AliasDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.Delete"
		log := slog.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		err := aliasDeleter.DeleteAlias(alias)
		if err != nil {
			log.Error(
				"failed to delete alias",
				slog.String("error", err.Error()),
				slog.String("alias", alias),
			)
			render.JSON(w, r, resp.Error("failed to delete alias"))
			return
		}

		log.Info("alias deleted", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Result:   "successfully deleted",
		})
		return
	}
}
