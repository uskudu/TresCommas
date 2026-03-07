package save

import (
	"log/slog"
	"net/http"
	resp "sptringTresRestAPI/internal/lib/api/response"
	"sptringTresRestAPI/internal/lib/random"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 7

//go:generate moq -out save_mock.go . URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func Save(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.Save"
		log := slog.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error("failed to save url", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved", slog.String("url", req.URL))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
		return
	}
}
