package portone

import (
	"context"
	"net/http"
	"net/url"
)

type paymentsService struct {
	httpClient *http.Client
	baseURL    *url.URL
}

func newPaymentsService(baseURL *url.URL, httpClient *http.Client) *paymentsService {
	return &paymentsService{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// CreatePaymentIntentRequest represents a request for 'POST /payments/prepare'.
type CreatePaymentIntentRequest struct {
	MerchantUID string `json:"merchant_uid"`
	Amount      int64  `json:"amount"`
}

// CreatePaymentIntentResponse represents a response of 'POST /payments/prepare'.
type CreatePaymentIntentResponse struct {
	CommonResponse
	Response struct {
		MerchantUID string `json:"merchant_uid"`
		Amount      int64  `json:"amount"`
	} `json:"response"`
}

// CreatePaymentIntent creates a new payment intent.
func (ps *paymentsService) CreatePaymentIntent(ctx context.Context, req CreatePaymentIntentRequest) (CreatePaymentIntentResponse, error) {
	u := ps.baseURL.JoinPath("/prepare")
	httpReq, err := newRequest(ctx, http.MethodPost, u.String(), req)
	if err != nil {
		return CreatePaymentIntentResponse{}, err
	}

	var resp CreatePaymentIntentResponse
	err = do(ps.httpClient, httpReq, &resp)
	if err != nil {
		return CreatePaymentIntentResponse{}, err
	}

	return resp, nil
}

type GetPaymentResponse struct {
	CommonResponse
	Response struct {
		ImpUID            string          `json:"imp_uid"`
		MerchantUID       string          `json:"merchant_uid"`
		PayMethod         string          `json:"pay_method"`
		Channel           string          `json:"channel"`
		PgProvider        string          `json:"pg_provider"`
		EmbPgProvider     string          `json:"emb_pg_provider"`
		PgTid             string          `json:"pg_tid"`
		PgID              string          `json:"pg_id"`
		Escrow            bool            `json:"escrow"`
		ApplyNum          string          `json:"apply_num"`
		BankCode          string          `json:"bank_code"`
		BankName          string          `json:"bank_name"`
		CardCode          string          `json:"card_code"`
		CardName          string          `json:"card_name"`
		CardQuota         int             `json:"card_quota"`
		CardNumber        string          `json:"card_number"`
		CardType          int             `json:"card_type"`
		VbankCode         string          `json:"vbank_code"`
		VbankName         string          `json:"vbank_name"`
		VbankNum          string          `json:"vbank_num"`
		VbankHolder       string          `json:"vbank_holder"`
		VbankDate         int             `json:"vbank_date"`
		VbankIssuedAt     int             `json:"vbank_issued_at"`
		Name              string          `json:"name"`
		Amount            int             `json:"amount"`
		CancelAmount      int             `json:"cancel_amount"`
		Currency          string          `json:"currency"`
		BuyerName         string          `json:"buyer_name"`
		BuyerEmail        string          `json:"buyer_email"`
		BuyerTel          string          `json:"buyer_tel"`
		BuyerAddr         string          `json:"buyer_addr"`
		BuyerPostcode     string          `json:"buyer_postcode"`
		CustomData        string          `json:"custom_data"`
		UserAgent         string          `json:"user_agent"`
		Status            string          `json:"status"`
		StartedAt         int             `json:"started_at"`
		PaidAt            int             `json:"paid_at"`
		FailedAt          int             `json:"failed_at"`
		CancelledAt       int             `json:"cancelled_at"`
		FailReason        string          `json:"fail_reason"`
		CancelReason      string          `json:"cancel_reason"`
		ReceiptURL        string          `json:"receipt_url"`
		CancelHistory     []CancelHistory `json:"cancel_history"`
		CancelReceiptUrls []string        `json:"cancel_receipt_urls"`
		CashReceiptIssued bool            `json:"cash_receipt_issued"`
		CustomerUID       string          `json:"customer_uid"`
		CustomerUIDUsage  string          `json:"customer_uid_usage"`
	} `json:"response"`
}

type CancelHistory struct {
	PgTid       string `json:"pg_tid"`
	Amount      int    `json:"amount"`
	CancelledAt int    `json:"cancelled_at"`
	Reason      string `json:"reason"`
	ReceiptURL  string `json:"receipt_url"`
}

func (ps *paymentsService) GetPayment(ctx context.Context, paymentID string) (GetPaymentResponse, error) {
	u := ps.baseURL.JoinPath(paymentID)
	httpReq, err := newRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return GetPaymentResponse{}, err
	}

	var resp GetPaymentResponse
	err = do(ps.httpClient, httpReq, &resp)
	if err != nil {
		return GetPaymentResponse{}, err
	}

	return resp, nil
}
