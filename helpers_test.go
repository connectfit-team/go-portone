package portone_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/connectfit-team/go-portone"
)

const (
	contentTypeJSON = "application/json"
)

func mustInitClient(t *testing.T) (*portone.Client, *http.ServeMux) {
	t.Helper()

	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client, err := portone.NewClient(testRestAPIKey, testRestAPISecret, portone.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}

	return client, mux
}

func mustInitClientWithAuthentication(t *testing.T) (*portone.Client, *http.ServeMux) {
	t.Helper()

	client, mux := mustInitClient(t)

	mux.HandleFunc("/users/getToken", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != contentTypeJSON {
			t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Accept") != contentTypeJSON {
			t.Errorf("unexpected accept: %s", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"code": 0,
			"message": "success",
			"response": {
				"access_token": "test_access_token",
				"now": 1600000000,
				"expired_at": 1700000000
			}
		}`))
	})

	return client, mux
}
