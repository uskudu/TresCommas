package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestSave(t *testing.T) {
	mockedURLSaver := &URLSaverMock{
		SaveURLFunc: func(urlToSave string, alias string) error {
			validUrl := regexp.MustCompile(`^https://[a-zA-Z0-9.-]+\.(com|ru|gov|org)$`)

			switch {
			case !validUrl.MatchString(urlToSave):
				return errors.New("invalid request")
			case urlToSave == "":
				return errors.New("invalid request")
			case validUrl.MatchString(urlToSave):
				return nil
			}
			return nil
		},
	}

	handler := Save(slog.Default(), mockedURLSaver)

	tests := []struct {
		name          string
		input         string
		specificAlias any
		want          any
	}{
		{"empty url", "", nil, "invalid request"},
		{"wrong url format", "httpNOTs:google.cum", nil, "failed to save url"},
		{"valid url without specific alias", "https://google.com", "", ""},
		{"valid url with specific alias", "https://nasa.gov", "nasa", "nasa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"url": "%s", "alias": "%v"}`, tt.input, tt.specificAlias)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler(w, req)

			res := w.Result()
			defer res.Body.Close()
			respBody, _ := io.ReadAll(res.Body)
			var respData struct {
				Status string `json:"status"`
				Alias  string `json:"alias"`
			}
			_ = json.Unmarshal(respBody, &respData)

			switch tt.name {
			case "valid url without specific alias":
				if respData.Alias == "" {
					t.Errorf("alias is empty, got %q", respData.Alias)
				}
			case "valid url with specific alias":
				if respData.Alias != tt.want {
					t.Errorf("got alias %q, want %q", respData.Alias, tt.want)
				}
			default:
				if !strings.Contains(string(respBody), tt.want.(string)) {
					t.Errorf("got %s, want %s", string(respBody), tt.want)
				}
			}
		})
	}
}
