package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

var UserAgent = "gbizinfo-cli/dev"

const defaultTimeout = 30 * time.Second

type Client struct {
	HTTP    *http.Client
	BaseURL string
	Token   string
	Verbose bool
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: defaultTimeout},
		BaseURL: baseURL,
		Token:   token,
	}
}

func (c *Client) Get(path string, result any) error {
	url := c.BaseURL + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("X-hojinInfo-api-token", c.Token)
	}

	if c.Verbose {
		debugLogRequest(req)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの送信に失敗: %w", err)
	}

	if c.Verbose {
		debugLogResponse(resp)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	return resp, nil
}

func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	message := string(body)
	var errResp struct {
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
		message = errResp.Message
	}

	switch resp.StatusCode {
	case 401, 403:
		return &cerrors.AuthError{Message: message}
	case 404:
		return &cerrors.NotFoundError{Resource: "リソース", ID: resp.Request.URL.Path}
	case 429:
		return &cerrors.RateLimitError{Message: message}
	default:
		return &cerrors.APIError{StatusCode: resp.StatusCode, Message: message}
	}
}

func debugLogRequest(req *http.Request) {
	fmt.Fprintf(os.Stderr, "> %s %s\n", req.Method, req.URL.String())
	for key, vals := range req.Header {
		for _, v := range vals {
			if strings.EqualFold(key, "X-hojinInfo-api-token") {
				if len(v) > 4 {
					v = v[:4] + "******"
				}
			}
			fmt.Fprintf(os.Stderr, "> %s: %s\n", key, v)
		}
	}
	fmt.Fprintln(os.Stderr)
}

func debugLogResponse(resp *http.Response) {
	fmt.Fprintf(os.Stderr, "< %s\n", resp.Status)
	for key, vals := range resp.Header {
		for _, v := range vals {
			fmt.Fprintf(os.Stderr, "< %s: %s\n", key, v)
		}
	}
	fmt.Fprintln(os.Stderr)
}
