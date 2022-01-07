package tgtg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAuthService_Login(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	loginRequest := &LoginRequest{
		DeviceType: "IOS",
		Email:      "some@email.com",
	}

	mux.HandleFunc("/auth/v3/authByEmail", func(w http.ResponseWriter, r *http.Request) {
		req := &LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, loginRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, loginRequest)
		}

		fmt.Fprintf(w, `
		{
			"polling_id": "polling_id", 
			"state": "state"
		}
		`)
	})

	actual, _, err := client.Auth.Login(context.Background(), loginRequest)
	if err != nil {
		t.Errorf("Auth.Login returned error: %+v", err)
	}

	expected := &LoginResponse{
		PollingID: "polling_id",
		State:     "state",
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Auth.Login returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Auth.Login returned: %+v, expected: %+v", actual, expected)
	}
}

func TestAuthService_LoginArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Auth.Login(context.Background(), nil)
	if !strings.Contains(err.Error(), "loginRequest") {
		t.Errorf("Error: %+v did not contain loginRequest", err)
	}
}

func TestAuthService_Poll(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	pollRequest := &PollRequest{
		DeviceType: "IOS",
		Email:      "some@email.com",
		PollingID:  "polling_id",
	}

	mux.HandleFunc("/auth/v3/authByRequestPollingId", func(w http.ResponseWriter, r *http.Request) {
		req := &PollRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, pollRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, pollRequest)
		}

		fmt.Fprintf(w, `
		{
			"access_token": "access_token", 
			"refresh_token": "refresh_token",
			"access_token_ttl_seconds": 172800,
			"startup_data": {
				"user": {
					"user_id": "1"
				}
			}
		}
		`)
	})

	actual, _, err := client.Auth.Poll(context.Background(), pollRequest)
	if err != nil {
		t.Errorf("Auth.Poll returned error: %+v", err)
	}

	expected := &PollResponse{
		AccessToken:    "access_token",
		RefreshToken:   "refresh_token",
		AccessTokenTTL: 172800,
		StartupData: StartupData{
			User: User{
				UserID: "1",
			},
		},
	}

	if client.AccessToken != "access_token" {
		t.Errorf("Auth.Poll returned: %+v, expected: access_token", client.AccessToken)
	}

	if client.RefreshToken != "refresh_token" {
		t.Errorf("Auth.Poll returned: %+v, expected: refresh_token", client.RefreshToken)
	}

	if client.UserID != "1" {
		t.Errorf("Auth.Poll returned: %+v, expected: 1", client.UserID)
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Auth.Poll returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Auth.Poll returned: %+v, expected: %+v", actual, expected)
	}
}

func TestAuthService_PollArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Auth.Poll(context.Background(), nil)
	if !strings.Contains(err.Error(), "pollRequest") {
		t.Errorf("Error: %+v did not contain pollRequest", err)
	}
}

func TestAuthService_Refresh(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	refreshRequest := &RefreshTokensRequest{
		RefreshToken: "refresh_token",
	}

	mux.HandleFunc("/auth/v3/token/refresh", func(w http.ResponseWriter, r *http.Request) {
		req := &RefreshTokensRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, refreshRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, refreshRequest)
		}

		fmt.Fprintf(w, `
		{
			"access_token": "access_token", 
			"refresh_token": "refresh_token",
			"access_token_ttl_seconds": 172800
		}
		`)
	})

	actual, _, err := client.Auth.Refresh(context.Background(), refreshRequest)
	if err != nil {
		t.Errorf("Auth.Refresh returned error: %+v", err)
	}

	expected := &RefreshTokensResponse{
		AccessToken:    "access_token",
		RefreshToken:   "refresh_token",
		AccessTokenTTL: 172800,
	}

	if client.AccessToken != "access_token" {
		t.Errorf("Auth.Refresh returned: %+v, expected: access_token", client.AccessToken)
	}

	if client.RefreshToken != "refresh_token" {
		t.Errorf("Auth.Refresh returned: %+v, expected: refresh_token", client.RefreshToken)
	}

	if client.UserID != "" {
		t.Errorf("Auth.Refresh returned: %+v, expected: empty string", client.UserID)
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Auth.Refresh returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Auth.Refresh returned: %+v, expected: %+v", actual, expected)
	}
}

func TestAuthService_RefreshArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Auth.Refresh(context.Background(), nil)
	if !strings.Contains(err.Error(), "refreshRequest") {
		t.Errorf("Error: %+v did not contain refreshRequest", err)
	}
}

func TestAuthService_Signup(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	signupRequest := &SignupRequest{
		CountryID:             "country_id",
		DeviceType:            "IOS",
		Email:                 "some@email.com",
		Name:                  "name",
		NewsletterOptIn:       true,
		PushNotificationOptIn: true,
	}

	mux.HandleFunc("/auth/v3/signUpByEmail", func(w http.ResponseWriter, r *http.Request) {
		req := &SignupRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, signupRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, signupRequest)
		}

		fmt.Fprintf(w, `
		{
			"login_response": {
				"access_token": "access_token", 
				"refresh_token": "refresh_token",
				"access_token_ttl_seconds": 172800,
				"startup_data": {
					"user": {
						"user_id": "id"
					}
				}
			}
		}
		`)
	})

	actual, _, err := client.Auth.Signup(context.Background(), signupRequest)
	if err != nil {
		t.Errorf("Auth.Signup returned error: %+v", err)
	}

	expected := &SignupResponse{
		Login: Login{
			AccessToken:    "access_token",
			RefreshToken:   "refresh_token",
			AccessTokenTTL: 172800,
			StartupData: StartupData{
				User: User{
					UserID: "id",
				},
			},
		},
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Auth.Signup returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Auth.Signup returned: %+v, expected: %+v", actual, expected)
	}
}

func TestAuthService_SignupArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Auth.Signup(context.Background(), nil)
	if !strings.Contains(err.Error(), "signupRequest") {
		t.Errorf("Error: %+v did not contain signupRequest", err)
	}
}
