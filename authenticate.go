package portone

import (
	"context"
	"net/http"
	"net/url"
)

type authenticateService struct {
	httpClient *http.Client
	baseURL    *url.URL
}

func newAuthenticateService(baseURL *url.URL) *authenticateService {
	return &authenticateService{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}
}

// GetTokenRequest represents a request for 'POST /users/getToken'.
type GetTokenRequest struct {
	RestAPIKey    string `json:"imp_key"`
	RestAPISecret string `json:"imp_secret"`
}

// GetTokenResponse represents a response for 'POST /users/getToken'.
type GetTokenResponse struct {
	CommonResponse
	Response struct {
		AccessToken string `json:"access_token"`
		Now         int64  `json:"now"`
		ExpiredAt   int64  `json:"expired_at"`
	} `json:"response"`
}

// GetToken returns a new token.
func (as *authenticateService) GetToken(ctx context.Context, req GetTokenRequest) (GetTokenResponse, error) {
	u := as.baseURL.JoinPath("/getToken")
	httpReq, err := newRequest(ctx, http.MethodPost, u.String(), req)
	if err != nil {
		return GetTokenResponse{}, err
	}

	var resp GetTokenResponse
	err = do(as.httpClient, httpReq, &resp)
	if err != nil {
		return GetTokenResponse{}, err
	}

	return resp, nil
}
