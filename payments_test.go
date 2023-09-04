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

func TestGetPayment(t *testing.T) {
	client, mux := mustInitClientWithAuthentication(t)

	mux.HandleFunc("/payments/test_imp_uid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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
				"imp_uid": "test_imp_uid",
				"merchant_uid": "test_merchant_uid",
				"pay_method": "card",
				"channel": "pc",
				"pg_provider": "nice",
				"pg_tid": "test_pg_tid",
				"pg_id": "test_pg_id",
				"pg_provider_uid": "test_pg_provider_uid",
				"escrow": false,
				"apply_num": "test_apply_num",
				"bank_code": "test_bank_code",
				"bank_name": "test_bank_name",
				"card_code": "test_card_code",
				"card_name": "test_card_name",
				"card_quota": 0,
				"vbank_code": "test_vbank_code",
				"vbank_name": "test_vbank_name",
				"vbank_num": "test_vbank_num",
				"vbank_holder": "test_vbank_holder",
				"vbank_date": 0,
				"name": "test_name",
				"amount": 1000,
				"cancel_amount": 0,
				"buyer_name": "test_buyer_name",
				"buyer_email": "test_buyer_email",
				"buyer_tel": "test_buyer_tel",
				"buyer_addr": "test_buyer_addr",
				"buyer_postcode": "test_buyer_postcode",
				"custom_data": "test_custom_data",
				"user_agent": "test_user_agent",
				"status": "ready",
				"started_at": 0,
				"paid_at": 0,
				"failed_at": 0,
				"cancelled_at": 0,
				"fail_reason": "test_fail_reason",
				"cancel_reason": "test_cancel_reason",
				"receipt_url": "test_receipt_url",
				"cancel_history": [
					{
						"pg_tid": "test_pg_tid",
						"amount": 1000,
						"cancelled_at": 0,
						"reason": "test_reason",
						"receipt_url": "test_receipt_url"
					}
				],
				"cancel_receipt_urls": [
					"test_cancel_receipt_urls"
				],
				"cash_receipt_issued": false,
				"customer_uid": "test_customer_uid",
				"customer_uid_usage": "test_customer_uid_usage"
			}
		}`))
	})

	resp, err := client.GetPayment(context.Background(), "test_imp_uid")
	if err != nil {
		t.Fatal(err)
	}

	want := portone.GetPaymentResponse{
		CommonResponse: portone.CommonResponse{
			Code:    0,
			Message: "success",
		},
		Response: struct {
			ImpUID            string                  `json:"imp_uid"`
			MerchantUID       string                  `json:"merchant_uid"`
			PayMethod         string                  `json:"pay_method"`
			Channel           string                  `json:"channel"`
			PgProvider        string                  `json:"pg_provider"`
			EmbPgProvider     string                  `json:"emb_pg_provider"`
			PgTid             string                  `json:"pg_tid"`
			PgID              string                  `json:"pg_id"`
			Escrow            bool                    `json:"escrow"`
			ApplyNum          string                  `json:"apply_num"`
			BankCode          string                  `json:"bank_code"`
			BankName          string                  `json:"bank_name"`
			CardCode          string                  `json:"card_code"`
			CardName          string                  `json:"card_name"`
			CardQuota         int                     `json:"card_quota"`
			CardNumber        string                  `json:"card_number"`
			CardType          string                  `json:"card_type"`
			VbankCode         string                  `json:"vbank_code"`
			VbankName         string                  `json:"vbank_name"`
			VbankNum          string                  `json:"vbank_num"`
			VbankHolder       string                  `json:"vbank_holder"`
			VbankDate         int                     `json:"vbank_date"`
			VbankIssuedAt     int                     `json:"vbank_issued_at"`
			Name              string                  `json:"name"`
			Amount            int                     `json:"amount"`
			CancelAmount      int                     `json:"cancel_amount"`
			Currency          string                  `json:"currency"`
			BuyerName         string                  `json:"buyer_name"`
			BuyerEmail        string                  `json:"buyer_email"`
			BuyerTel          string                  `json:"buyer_tel"`
			BuyerAddr         string                  `json:"buyer_addr"`
			BuyerPostcode     string                  `json:"buyer_postcode"`
			CustomData        string                  `json:"custom_data"`
			UserAgent         string                  `json:"user_agent"`
			Status            string                  `json:"status"`
			StartedAt         int                     `json:"started_at"`
			PaidAt            int                     `json:"paid_at"`
			FailedAt          int                     `json:"failed_at"`
			CancelledAt       int                     `json:"cancelled_at"`
			FailReason        string                  `json:"fail_reason"`
			CancelReason      string                  `json:"cancel_reason"`
			ReceiptURL        string                  `json:"receipt_url"`
			CancelHistory     []portone.CancelHistory `json:"cancel_history"`
			CancelReceiptUrls []string                `json:"cancel_receipt_urls"`
			CashReceiptIssued bool                    `json:"cash_receipt_issued"`
			CustomerUID       string                  `json:"customer_uid"`
			CustomerUIDUsage  string                  `json:"customer_uid_usage"`
		}{
			ImpUID:            "test_imp_uid",
			MerchantUID:       "test_merchant_uid",
			PayMethod:         "card",
			Channel:           "pc",
			PgProvider:        "nice",
			PgTid:             "test_pg_tid",
			PgID:              "test_pg_id",
			Escrow:            false,
			ApplyNum:          "test_apply_num",
			BankCode:          "test_bank_code",
			BankName:          "test_bank_name",
			CardCode:          "test_card_code",
			CardName:          "test_card_name",
			CardQuota:         0,
			VbankCode:         "test_vbank_code",
			VbankName:         "test_vbank_name",
			VbankNum:          "test_vbank_num",
			VbankHolder:       "test_vbank_holder",
			VbankDate:         0,
			Name:              "test_name",
			Amount:            1000,
			CancelAmount:      0,
			BuyerName:         "test_buyer_name",
			BuyerEmail:        "test_buyer_email",
			BuyerTel:          "test_buyer_tel",
			BuyerAddr:         "test_buyer_addr",
			BuyerPostcode:     "test_buyer_postcode",
			CustomData:        "test_custom_data",
			UserAgent:         "test_user_agent",
			Status:            "ready",
			StartedAt:         0,
			PaidAt:            0,
			FailedAt:          0,
			CancelledAt:       0,
			FailReason:        "test_fail_reason",
			CancelReason:      "test_cancel_reason",
			ReceiptURL:        "test_receipt_url",
			CancelHistory:     []portone.CancelHistory{{PgTid: "test_pg_tid", Amount: 1000, CancelledAt: 0, Reason: "test_reason", ReceiptURL: "test_receipt_url"}},
			CancelReceiptUrls: []string{"test_cancel_receipt_urls"},
			CashReceiptIssued: false,
			CustomerUID:       "test_customer_uid",
			CustomerUIDUsage:  "test_customer_uid_usage",
		},
	}

	if diff := cmp.Diff(want, resp); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
