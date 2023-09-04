package portone

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL          = "https://api.iamport.kr"
	authenticateServicePath = "/users"
	paymentsServicePath     = "/payments"
)

var (
	defaultClientTimeout = 10 * time.Second
)

var defaultClientConfig = clientConfig{
	baseURL: defaultBaseURL,
	timeout: defaultClientTimeout,
}

type clientConfig struct {
	baseURL string
	timeout time.Duration
}

type ClientOption func(*clientConfig)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) {
		c.timeout = timeout
	}
}

// Client is a client for Portone API.
type Client struct {
	clientConfig

	*authenticateService
	*paymentsService
}

// NewClient returns a new PortOne API client.
func NewClient(restAPIKey, restAPISecret string, opts ...ClientOption) (*Client, error) {
	cfg := defaultClientConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	u, err := url.Parse(cfg.baseURL)
	if err != nil {
		return nil, err
	}

	authenticateServiceBaseURL := u.JoinPath(authenticateServicePath)
	authenticateService := newAuthenticateService(authenticateServiceBaseURL)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: defaultClientTimeout,
		Transport: &roundTripperWithToken{
			authenticateService: authenticateService,
			restAPIKey:          restAPIKey,
			restAPISecret:       restAPISecret,
		},
	}

	paymentsServiceBaseURL := u.JoinPath(paymentsServicePath)
	paymentsService := newPaymentsService(paymentsServiceBaseURL, httpClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientConfig:        cfg,
		authenticateService: authenticateService,
		paymentsService:     paymentsService,
	}, nil
}

func newRequest(ctx context.Context, method, urlStr string, body any) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, urlStr, buf)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	return httpReq, nil
}

func do(httpClient *http.Client, req *http.Request, respBody any) error {
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(respBody)
	if err != nil {
		return err
	}

	return nil
}

type roundTripperWithToken struct {
	authenticateService *authenticateService
	restAPIKey          string
	restAPISecret       string
	accessToken         string
	expireAt            int64
}

func (rt *roundTripperWithToken) RoundTrip(req *http.Request) (*http.Response, error) {
	if !rt.isAuthenticated() || rt.isAccessTokenExpired() {
		resp, err := rt.authenticateService.GetToken(context.Background(), GetTokenRequest{
			RestAPIKey:    rt.restAPIKey,
			RestAPISecret: rt.restAPISecret,
		})
		if err != nil {
			return nil, err
		}

		rt.setToken(resp.Response.AccessToken, resp.Response.ExpiredAt)
	}

	req.Header.Set("Authorization", rt.accessToken)
	return http.DefaultTransport.RoundTrip(req)
}

func (rt *roundTripperWithToken) isAuthenticated() bool {
	return rt.accessToken != ""
}

func (rt *roundTripperWithToken) isAccessTokenExpired() bool {
	return rt.isAuthenticated() && rt.expireAt <= time.Now().Unix()
}

func (rt *roundTripperWithToken) setToken(token string, expireAt int64) {
	rt.accessToken = token
	rt.expireAt = expireAt
}

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
