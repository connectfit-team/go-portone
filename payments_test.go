package portone_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/connectfit-team/go-portone"
	"github.com/google/go-cmp/cmp"
)

func TestCreatePaymentIntent(t *testing.T) {
	client, mux := mustInitClientWithAuthentication(t)

	mux.HandleFunc("/payments/prepare", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != contentTypeJSON {
			t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Accept") != contentTypeJSON {
			t.Errorf("unexpected accept: %s", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"code": 0,
			"message": "success",
			"response": {
				"merchant_uid": "test_merchant_uid",
				"amount": 1000
			}
		}`))
	})

	resp, err := client.CreatePaymentIntent(context.Background(), portone.CreatePaymentIntentRequest{
		MerchantUID: "test_merchant_uid",
		Amount:      1000,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := portone.CreatePaymentIntentResponse{
		CommonResponse: portone.CommonResponse{
			Code:    0,
			Message: "success",
		},
		Response: struct {
			MerchantUID string `json:"merchant_uid"`
			Amount      int64  `json:"amount"`
		}{
			MerchantUID: "test_merchant_uid",
			Amount:      1000,
		},
	}

	if diff := cmp.Diff(want, resp); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
