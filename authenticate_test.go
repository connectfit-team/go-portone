package portone_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/connectfit-team/go-portone"
	"github.com/google/go-cmp/cmp"
)

const (
	testRestAPIKey    = "test_rest_api_key"
	testRestAPISecret = "test_rest_api_secret"
)

func TestGetToken(t *testing.T) {
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

	resp, err := client.GetToken(context.Background(), portone.GetTokenRequest{
		RestAPIKey:    testRestAPIKey,
		RestAPISecret: testRestAPISecret,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := portone.GetTokenResponse{
		CommonResponse: portone.CommonResponse{
			Code:    0,
			Message: "success",
		},
		Response: struct {
			AccessToken string `json:"access_token"`
			Now         int64  `json:"now"`
			ExpiredAt   int64  `json:"expired_at"`
		}{
			AccessToken: "test_access_token",
			Now:         1600_000_000,
			ExpiredAt:   1700_000_000,
		},
	}

	if diff := cmp.Diff(want, resp); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
