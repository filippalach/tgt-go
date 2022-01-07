package tgtg

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestArgumentError(t *testing.T) {
	expected := "foo argument is invalid (reason: bar)"
	err := NewArgumentError("foo", "bar")
	if actual := err.Error(); actual != expected {
		t.Errorf("ArgumentError: %+v; expected: %+v", actual, expected)
	}
}

func TestErrorResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	loginRequest := &LoginRequest{}

	mux.HandleFunc("/auth/v3/authByEmail", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `
		{
			"errors": [{"code": "INTERNAL_SERVER_ERROR"}]
		}
		`)
	})

	_, _, err := client.Auth.Login(context.Background(), loginRequest)
	if err != err.(*ErrorResponse) {
		t.Errorf("Auth.Login return: %+v, expected: %+v", err, err.(*ErrorResponse))
	}

	if !strings.Contains(err.Error(), "INTERNAL_SERVER_ERROR") {
		t.Errorf("Error: %+v did not contain INTERNAL_SERVER_ERROR", err)
	}
}
