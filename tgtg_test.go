package tgtg

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNew_DefaultValues(t *testing.T) {
	client, err := New(nil)
	if err != nil {
		t.Fatalf("New return: %+v, expected no errer.", err)
	}

	if client.BaseURL == nil || client.BaseURL.String() != defaultBaseURL {
		t.Fatalf("New base url: %+v, expected: %+v", client.BaseURL, defaultBaseURL)
	}

	if client.UserAgent == "" || client.UserAgent != defaultUserAgent {
		t.Fatalf("New user agent: %+v, expected: %+v", client.BaseURL, defaultBaseURL)
	}
}

func TestNew_CustomUserAgent(t *testing.T) {
	client, err := New(nil, SetUserAgent("added user agent"))
	if err != nil {
		t.Fatalf("New return: %+v, expected no errer.", err)
	}

	if !strings.Contains(client.UserAgent, "added user agent") {
		t.Fatalf("New user agent: %+v, expected: %+v", client.UserAgent, "user agent")
	}
}

func TestNew_CustomRequestHeaders(t *testing.T) {
	client, err := New(nil, SetRequestHeaders(map[string]string{
		"Accept-Encoding": "chunked",
	}))
	if err != nil {
		t.Fatalf("New return: %+v, expected no errer.", err)
	}

	req, _ := client.NewRequest(http.MethodGet, "/", nil)

	if actual := req.Header.Get("Accept-Encoding"); actual != "chunked" {
		t.Errorf("New header: %+v, expected %+v", actual, "chunked")
	}
}

func TestNewRequest_valid(t *testing.T) {
	url := "no/prefix"
	expectedURL := defaultBaseURL + url

	c := NewClient(nil)

	inBody := &LoginRequest{
		DeviceType: "IOS",
		Email:      "some@email.com",
	}
	expectedBody := `{"device_type":"IOS","email":"some@email.com"}` + "\n"

	req, _ := c.NewRequest(http.MethodPost, url, inBody)
	if req.URL.String() != expectedURL {
		t.Errorf("NewRequest url: %+v, expected: %+v", req.URL, expectedURL)
	}

	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != expectedBody {
		t.Errorf("NewRequest body: %+v, expected: %+v", string(body), expectedBody)
	}
}

func TestDo(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	type test struct {
		Test string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"Test":"test"}`)
	})

	req, _ := client.NewRequest(http.MethodGet, "/", nil)
	body := &test{}
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do: %+v", err)
	}

	expected := &test{"test"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body: %+v, expected: %+v", body, expected)
	}
}

func TestCheckResponse(t *testing.T) {
	testCases := []struct {
		title            string
		response         *http.Response
		expectedResponse *ErrorResponse
	}{
		{
			title: "JSON, only Code - no Message",
			response: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(`
				{
					"errors": [{"code": "INTERNAL_SERVER_ERROR", "message": ""}]
				}`)),
			},
			expectedResponse: &ErrorResponse{
				Errors: []Error{
					{
						Code: "INTERNAL_SERVER_ERROR",
					},
				},
			},
		},
		{
			title: "JSON, Code with Message",
			response: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body: ioutil.NopCloser(strings.NewReader(`
				{
					"errors": [
						{"code": "CODE_1", "message": "error message"},
						{"code": "CODE_2", "message": "error message"}
					]
				}`)),
			},
			expectedResponse: &ErrorResponse{
				Errors: []Error{
					{
						Code:    "CODE_1",
						Message: "error message",
					},
					{
						Code:    "CODE_2",
						Message: "error message",
					},
				},
			},
		},
		{
			title: "String response",
			response: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader("string response")),
			},
			expectedResponse: &ErrorResponse{
				Errors: []Error{
					{
						Code: "string response",
					},
				},
			},
		},
		{
			title: "No response body",
			response: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			expectedResponse: &ErrorResponse{
				Errors: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := CheckResponseForErrors(tc.response)
			if err == nil {
				t.Fatal("Expected error response.")
			}
			tc.expectedResponse.Response = tc.response

			if !reflect.DeepEqual(err, tc.expectedResponse) {
				t.Errorf("Error: %+v, expected: %+v", err, tc.expectedResponse)
			}
		})
	}
}

func setup() (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Fatalf("Request method: %+v, expected: %+v", r.Method, expected)
	}
}
