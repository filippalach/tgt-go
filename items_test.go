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

func TestItemsService_List(t *testing.T) {
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

	listRequest := &ListItemsRequest{
		Radius:   40,
		PageSize: 400,
		Page:     1,
		Origin: &Origin{
			Latitude:  45.624,
			Longitude: 9.282,
		},
	}

	mux.HandleFunc("/item/v7/", func(w http.ResponseWriter, r *http.Request) {
		req := &ListItemsRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, listRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, listRequest)
		}

		fmt.Fprintf(w, `
		{
			"items": [
				{
					"item": {
						"name": "name-1"
					},
					"store": {
						"description": "description-1"
					},
					"display_name": "name-1",
					"items_available": 1
				},
				{
					"item": {
						"name": "name-2"
					},
					"store": {
						"description": "description-2"
					},
					"display_name": "name-2",
					"items_available": 2
				}
			]
		}
		`)
	})

	actual, _, err := client.Items.List(context.Background(), listRequest)
	if err != nil {
		t.Errorf("Items.List returned error: %+v", err)
	}

	expected := &ListItemsResponse{
		Items: []Items{
			{
				Item: Item{
					Name: "name-1",
				},
				Store: Store{
					Description: "description-1",
				},
				DisplayName:    "name-1",
				ItemsAvailable: 1,
			},
			{
				Item: Item{
					Name: "name-2",
				},
				Store: Store{
					Description: "description-2",
				},
				DisplayName:    "name-2",
				ItemsAvailable: 2,
			},
		},
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Items.List returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Items.List returned: %+v, expected: %+v", actual, expected)
	}
}

func TestItemsService_ListNoUserID(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	listRequest := &ListItemsRequest{
		Radius:   40,
		PageSize: 400,
		Page:     1,
		Origin: &Origin{
			Latitude:  45.624,
			Longitude: 9.282,
		},
	}

	_, _, err := client.Items.List(context.Background(), listRequest)
	if err == nil {
		t.Fatal("Items.List returned no error.")
	}
}

func TestItemsService_ListArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Items.List(context.Background(), nil)
	if !strings.Contains(err.Error(), "listItemsRequest") {
		t.Errorf("Error: %+v did not contain listItemsRequest", err)
	}

	_, _, err = client.Items.List(context.Background(), &ListItemsRequest{PageSize: 1})
	if !strings.Contains(err.Error(), "listItemsRequest.UserID") {
		t.Errorf("Error: %+v did not contain listItemsRequest.UserID", err)
	}
}

func TestItemsService_Get(t *testing.T) {
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

	getRequest := &GetItemRequest{
		Origin: &Origin{
			Latitude:  45.624,
			Longitude: 9.282,
		},
	}

	mux.HandleFunc("/item/v7/1", func(w http.ResponseWriter, r *http.Request) {
		req := &GetItemRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, getRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, getRequest)
		}

		fmt.Fprintf(w, `
		{
			"item": {
				"item_id": "1",
				"name": "name"
			},
			"store": {
				"description": "description"
			},
			"display_name": "name",
			"favorite": true
		}
		`)
	})

	actual, _, err := client.Items.Get(context.Background(), getRequest, "1")
	if err != nil {
		t.Errorf("Items.Get returned error: %+v", err)
	}

	expected := &GetItemResponse{
		Item: Item{
			ItemID: "1",
			Name:   "name",
		},
		Store: Store{
			Description: "description",
		},
		DisplayName: "name",
		Favorite:    true,
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("Items.Get returned: %+v, expected: %+v", actual, expected)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Items.Get returned: %+v, expected: %+v", actual, expected)
	}
}

func TestItemsService_GetArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, _, err := client.Items.Get(context.Background(), nil, "1")
	if !strings.Contains(err.Error(), "getItemRequest") {
		t.Errorf("Error: %+v did not contain getItemRequest", err)
	}

	_, _, err = client.Items.Get(context.Background(), &GetItemRequest{}, "")
	if !strings.Contains(err.Error(), "itemID") {
		t.Errorf("Error: %+v did not contain itemID", err)
	}

	_, _, err = client.Items.Get(context.Background(), &GetItemRequest{Origin: &Origin{Latitude: 1}}, "1")
	if !strings.Contains(err.Error(), "getItemRequest.UserID") {
		t.Errorf("Error: %+v did not contain getItemRequest.UserID", err)
	}
}

func TestItemsService_Favorite(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	favoriteRequest := &FavoriteItemRequest{
		IsFavorite: true,
	}

	mux.HandleFunc("/item/v7/1/setFavorite", func(w http.ResponseWriter, r *http.Request) {
		req := &FavoriteItemRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Fatalf("Decode json: %+v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(req, favoriteRequest) {
			t.Errorf("Request body: %+v, expected: %+v", req, favoriteRequest)
		}
	})

	_, err := client.Items.Favorite(context.Background(), favoriteRequest, "1")
	if err != nil {
		t.Errorf("Items.Favorite returned error: %+v", err)
	}
}

func TestItemsService_FavoriteArgumentError(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, err := client.Items.Favorite(context.Background(), nil, "1")
	if !strings.Contains(err.Error(), "favoriteItemRequest") {
		t.Errorf("Error: %+v did not contain favoriteItemRequest", err)
	}

	_, err = client.Items.Favorite(context.Background(), &FavoriteItemRequest{}, "")
	if !strings.Contains(err.Error(), "itemID") {
		t.Errorf("Error: %+v did not contain itemID", err)
	}
}
