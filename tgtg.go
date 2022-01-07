package tgtg

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://apptoogoodtogo.com/api/"
	defaultUserAgent = "TooGoodToGo/21.11.0 (iPhone/iPhone 11 Pro; iOS 14.4.1; Scale/3.00)"

	mediaType = "application/json"
)

// Client represents Too Good To Go client.
type Client struct {
	// HTTP client used to communicate with the Too Good To Go API.
	client *http.Client

	// AccessToken used in subsequent API requests, if set.
	AccessToken string

	// RefreshToken used in refresh API request, if set.
	RefreshToken string

	// UserID is set at the end on successful logging attempt ending with Poll.
	// It is required in multiple Too Good To Go API calls.
	UserID string

	// Base URL for API requests, has to end with "/".
	BaseURL *url.URL

	// User agent for client.
	UserAgent string

	// Services used for communicating with the Too Good To Go API.
	Auth   AuthService
	Items  ItemsService
	Orders OrdersService

	// Optional extra HTTP headers to set on every request to the Too Good To Go API.
	headers map[string]string
}

// NewClient returns a new Too Good To Go API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
	c.Auth = &AuthServiceOp{client: c}
	c.Items = &ItemsServiceOp{client: c}
	c.Orders = &OrdersServiceOp{client: c}

	c.headers = make(map[string]string)

	return c
}

// New returns a new Too Good To Go API client.
//
// Multiple ClientOptions can be passed to configure client.
func New(httpClient *http.Client, options ...ClientOption) (*Client, error) {
	client := NewClient(httpClient)
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// ClientOption are options for New.
type ClientOption func(*Client) error

// SetUserAgent is a ClientOption for setting the user agent.
func SetUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// SetRequestHeaders sets optional HTTP headers on the client that are sent on each HTTP request.
func SetRequestHeaders(headers map[string]string) ClientOption {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

// SetAuthContext sets Too Good To Go API auth context: access token, refresh token and user id.
func (c *Client) SetAuthContext(accessToken, refreshToken, userID string) {
	c.AccessToken = accessToken
	c.RefreshToken = refreshToken
	c.UserID = userID
}

// NewRequest creates Too Good To Go API request. Relative URL has to be provided in url, which will be merged with
// BaseURL of the Client. URL has to start with no "/" prefix, only relative URL is handled. If specified, the value pointing at body would be
// JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var request *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		request, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := &bytes.Buffer{}
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		request, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		request.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		request.Header.Add(k, v)
	}

	request.Header.Set("Accept", mediaType)
	request.Header.Set("User-Agent", c.UserAgent)

	return request, nil
}

// Do sends Too Good To Go API request and returns response. The response is JSON decoded
// and stored in the value pointed to by v, or returned as an error if occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	response, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	err = CheckResponseForErrors(response)
	if err != nil {
		return nil, err
	}

	if v != nil {
		// fmt.Println("123")
		// x, _ := io.ReadAll(response.Body)
		// fmt.Println(string(x))
		err = json.NewDecoder(response.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return c.client.Do(req)
}

// CheckResponseForErrors checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 2xx range. Error response is expected to have either no response
// body, or a JSON response body that maps to ErrorResponse or plain text. Either way Client will log out the message.
func CheckResponseForErrors(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Errors = append(errorResponse.Errors, Error{Code: string(data)})
		}
	}

	return errorResponse
}
