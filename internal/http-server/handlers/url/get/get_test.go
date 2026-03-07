package get

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sptringTresRestAPI/internal/storage"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/context"
)

func TestGet(t *testing.T) {
	mockedURLGetter := &URLGetterMock{
		GetURLFunc: func(alias string) (string, error) {
			switch alias {
			case "google":
				return "https://google.com", nil
			case "nasa":
				return "https://nasa.gov", nil
			case "":
				return "alias is required", nil
			default:
				return "", storage.ErrURLNotFound
			}
		},
	}

	handler := Get(slog.Default(), mockedURLGetter)

	tests := []struct {
		name  string
		input string
		want  any
	}{
		{"empty alias", "", "alias is required"},
		{"not existing alias", "notExistingAlias", storage.ErrURLNotFound},
		{"impossible alias", "!(*&@$@(!@%(*)!%#(*)#!%*(@#)*(", storage.ErrURLNotFound},
		{"existing alias 1", "google", "https://google.com"},
		{"existing alias 2", "nasa", "https://nasa.gov"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tt.input)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()
		})
	}
}
