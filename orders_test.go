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

func TestOrdersService_Active(t *testing.T) {
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
			"startup_data": {
				"user": {
					"user_id": "1"
				}
			}
		}
		`)
	})

	_, _, err := client.Auth.Poll(context.Background(), pollRequest)
	if err != nil {
		t.Errorf("Auth.Poll returned error: %+v", err)
	}

	mux.HandleFunc("/order/v6/active", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprintf(w, `
		{
			"has_more": true,
			"orders": [
				{
					"order_id": "order_id",
					"state": "REDEEMED",
					"quantity": 1,
					"store_id": "111111"
				}
			]
		}
		`)
	})

	actual, _, err := client.Orders.Active(context.Background(), nil)
	if err != nil {
		t.Errorf("Orders.Active returned error: %+v", err)
	}

	expected := &OrdersResponse{
		HasMore: true,
		Orders: []Order{
			{
				OrderID:  "order_id",
				State:    "REDEEMED",
				Quantity: 1,
				StoreID:  "111111",
			},
		},
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Orders.Active returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Orders.Active returned: %+v, expected: %+v", actual, expected)
	}
}

func TestOrdersService_ActiveArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Orders.Active(context.Background(), nil)
	if !strings.Contains(err.Error(), "activeOrdersRequest") {
		t.Errorf("Error: %+v did not contain activeOrdersRequest", err)
	}
}

func TestOrdersService_Inactive(t *testing.T) {
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
			"startup_data": {
				"user": {
					"user_id": "1"
				}
			}
		}
		`)
	})

	_, _, err := client.Auth.Poll(context.Background(), pollRequest)
	if err != nil {
		t.Errorf("Auth.Poll returned error: %+v", err)
	}

	inactiveRequest := &InactiveOrdersRequest{
		Paging: Paging{
			Page: 0,
			Size: 200,
		},
	}

	mux.HandleFunc("/order/v6/inactive", func(w http.ResponseWriter, r *http.Request) {
		req := &InactiveOrdersRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, inactiveRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, inactiveRequest)
		}

		fmt.Fprintf(w, `
		{
			"has_more": true,
			"orders": [
				{
					"order_id": "order_id",
					"state": "REDEEMED",
					"quantity": 1,
					"store_id": "111111"
				}
			]
		}
		`)
	})

	actual, _, err := client.Orders.Inactive(context.Background(), inactiveRequest)
	if err != nil {
		t.Errorf("Orders.Inactive returned error: %+v", err)
	}

	expected := &OrdersResponse{
		HasMore: true,
		Orders: []Order{
			{
				OrderID:  "order_id",
				State:    "REDEEMED",
				Quantity: 1,
				StoreID:  "111111",
			},
		},
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Orders.Inactive returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Orders.Inactive returned: %+v, expected: %+v", actual, expected)
	}
}

func TestOrdersService_InactiveArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Orders.Inactive(context.Background(), nil)
	if !strings.Contains(err.Error(), "inactiveOrdersRequest") {
		t.Errorf("Error: %+v did not contain inactiveOrdersRequest", err)
	}

	_, _, err = client.Orders.Inactive(context.Background(), &InactiveOrdersRequest{})
	if !strings.Contains(err.Error(), "inactiveOrdersRequest.UserID") {
		t.Errorf("Error: %+v did not contain inactiveOrdersRequest.UserID", err)
	}
}
