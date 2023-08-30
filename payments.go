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
