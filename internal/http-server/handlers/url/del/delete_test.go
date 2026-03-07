package del

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestDelete(t *testing.T) {
	mockedURLGetter := &AliasDeleterMock{
		DeleteAliasFunc: func(alias string) error {
			switch alias {
			case "":
				return errors.New("alias is empty")
			case "notExistingAlias":
				return errors.New("failed to delete alias")
			case "google":
				return nil
			}
			return nil
		},
	}

	handler := Delete(slog.Default(), mockedURLGetter)

	tests := []struct {
		name  string
		input string
		want  any
	}{
		{"empty alias", "", errors.New("alias is empty")},
		{"not existing alias", "notExistingAlias", errors.New("failed to delete alias")},
		{"existing alias", "google", ""},
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
