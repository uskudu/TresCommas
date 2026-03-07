package redirect

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sptringTresRestAPI/internal/storage"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRedirect(t *testing.T) {
	mockedURLGetter := &URLGetterMock{
		GetURLFunc: func(alias string) (string, error) {
			switch alias {
			case "google":
				return "https://google.com", nil
			case "nasa":
				return "https://nasa.gov", nil
			case "":
				return "", errors.New("alias is required")
			default:
				return "", storage.ErrURLNotFound
			}
		},
	}

	handler := Redirect(slog.Default(), mockedURLGetter)

	tests := []struct {
		name     string
		input    string
		wantBody string
		wantURL  string
	}{
		{"empty alias", "", "alias is required", ""},
		{"not existing alias", "notExistingAlias", "url not found", ""},
		{"existing alias 1", "google", "", "https://google.com"},
		{"existing alias 2", "nasa", "", "https://nasa.gov"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("alias", tt.input)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			if tt.wantBody != "" && !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("got body %q, want %q", string(body), tt.wantBody)
			}
			if tt.wantURL != "" && res.StatusCode == http.StatusFound {
				loc, _ := res.Location()
				if loc.String() != tt.wantURL {
					t.Errorf("got redirect to %q, want %q", loc.String(), tt.wantURL)
				}
			}
		})
	}
}
