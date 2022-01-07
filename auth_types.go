package tgtg

// LoginRequest represents a request body to initiate authentication process.
type LoginRequest struct {
	DeviceType string `json:"device_type"`
	Email      string `json:"email"`
}

// LoginResponse represents a response body in authentication process.
type LoginResponse struct {
	PollingID string `json:"polling_id"`
	State     string `json:"state"`
}

// PollRequest represents a request to check current authentication status.
type PollRequest struct {
	DeviceType string `json:"device_type"`
	Email      string `json:"email"`
	PollingID  string `json:"request_polling_id"`
}

// PollResponse represents a response body in polling process.
type PollResponse struct {
	AccessToken    string      `json:"access_token"`
	RefreshToken   string      `json:"refresh_token"`
	AccessTokenTTL int         `json:"access_token_ttl_seconds"`
	StartupData    StartupData `json:"startup_data"`
}

// RefreshTokensRequest represents a request body to refresh Too Good To Go tokens.
type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokensResponse represents Too Good To Go Tokens returned in authentication process.
type RefreshTokensResponse struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	AccessTokenTTL int    `json:"access_token_ttl_seconds"`
}

// SignupRequest represents a request body to create a new account.
type SignupRequest struct {
	CountryID             string `json:"country_id"`
	DeviceType            string `json:"device_type"`
	Email                 string `json:"email"`
	Name                  string `json:"name"`
	NewsletterOptIn       bool   `json:"newsletter_opt_in"`
	PushNotificationOptIn bool   `json:"push_notification_opt_in"`
}

// SignupResponse represents a response body in signup process.
type SignupResponse struct {
	Login Login `json:"login_response"`
}

// Login is a part of SignupResponse and represents a response body in sucessful signup process.
type Login struct {
	AccessToken    string      `json:"access_token"`
	RefreshToken   string      `json:"refresh_token"`
	AccessTokenTTL int         `json:"access_token_ttl_seconds"`
	StartupData    StartupData `json:"startup_data"`
}

// StartupData contains info about freshly authnticated user.
type StartupData struct {
	User User `json:"user"`
}

// User contains info about user's ID.
type User struct {
	UserID string `json:"user_id"`
}
