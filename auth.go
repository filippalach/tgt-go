package tgtg

import (
	"context"
	"fmt"
	"net/http"
)

const authBasePath = "auth/v3"

// AuthService is an interface for interfacing with the Auth related endpoints of the Too Good To Go API.
type AuthService interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, *http.Response, error)
	Poll(context.Context, *PollRequest) (*PollResponse, *http.Response, error)
	Refresh(context.Context, *RefreshTokensRequest) (*RefreshTokensResponse, *http.Response, error)
	Signup(context.Context, *SignupRequest) (*SignupResponse, *http.Response, error)
}

// AuthServiceOp handles communication with the Auth related methods of the Too Good To Go API.
type AuthServiceOp struct {
	client *Client
}

var _ AuthService = &AuthServiceOp{}

// Login handles initiating authentication process. It does not return tokens.
// After Login user have to manually finish logging process using received mail.
//
// Calling this endpoint frequently will cause 429 TOO_MANY_REQUESTS HTTP error,
// which will decay after ~15mins.
func (s *AuthServiceOp) Login(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, *http.Response, error) {
	if loginRequest == nil {
		return nil, nil, NewArgumentError("loginRequest", "must not be nil")
	}

	url := fmt.Sprintf("%s/authByEmail", authBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, loginRequest)
	if err != nil {
		return nil, nil, err
	}

	loginResponse := &LoginResponse{}
	response, err := s.client.Do(ctx, req, loginResponse)
	if err != nil {
		return nil, nil, err
	}

	return loginResponse, response, nil
}

// Poll handles checking current authentication status.
//
// If the login was finished using email, it return tokens,
// otherwise it errors out with EOF, empty body HTTP 202 ACCEPTED response.
func (s *AuthServiceOp) Poll(ctx context.Context, pollRequest *PollRequest) (*PollResponse, *http.Response, error) {
	if pollRequest == nil {
		return nil, nil, NewArgumentError("pollRequest", "must not be nil")
	}

	url := fmt.Sprintf("%s/authByRequestPollingId", authBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, pollRequest)
	if err != nil {
		return nil, nil, err
	}

	pollResponse := &PollResponse{}
	response, err := s.client.Do(ctx, req, pollResponse)
	if err != nil {
		return nil, nil, err
	}
	s.client.SetAuthContext(pollResponse.AccessToken, pollResponse.RefreshToken, pollResponse.StartupData.User.UserID)

	return pollResponse, response, nil
}

// Refresh handles refreshing tokens.
//
// If refreshRequest does not specify refresh token to use, client will attempt to use one already stored.
// For that to work authenticate client first using Login and Poll, or use SetAuthContext manually.
func (s *AuthServiceOp) Refresh(ctx context.Context, refreshRequest *RefreshTokensRequest) (*RefreshTokensResponse, *http.Response, error) {
	if refreshRequest == nil {
		if s.client.RefreshToken == "" {
			return nil, nil, NewArgumentError("refreshRequest", "must not be nil - client has no refresh token")
		}
		// if refresh token was not passed, but we have already authenticated with API, use refresh token stored in client.
		refreshRequest = &RefreshTokensRequest{RefreshToken: s.client.RefreshToken}
	}

	url := fmt.Sprintf("%s/token/refresh", authBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, refreshRequest)
	if err != nil {
		return nil, nil, err
	}

	refreshResponse := &RefreshTokensResponse{}
	response, err := s.client.Do(ctx, req, refreshResponse)
	if err != nil {
		return nil, nil, err
	}
	s.client.SetAuthContext(refreshResponse.AccessToken, refreshResponse.RefreshToken, s.client.UserID)

	return refreshResponse, response, nil
}

// Signup handles signing up a new account.
func (s *AuthServiceOp) Signup(ctx context.Context, signupRequest *SignupRequest) (*SignupResponse, *http.Response, error) {
	if signupRequest == nil {
		return nil, nil, NewArgumentError("signupRequest", "must not be nil")
	}

	url := fmt.Sprintf("%s/signUpByEmail", authBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, signupRequest)
	if err != nil {
		return nil, nil, err
	}

	signupResponse := &SignupResponse{}
	response, err := s.client.Do(ctx, req, signupResponse)
	if err != nil {
		return nil, nil, err
	}
	s.client.SetAuthContext(signupResponse.Login.AccessToken, signupResponse.Login.RefreshToken, signupResponse.Login.StartupData.User.UserID)

	return signupResponse, response, nil
}
