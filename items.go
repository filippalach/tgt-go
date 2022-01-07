package tgtg

import (
	"context"
	"fmt"
	"net/http"
)

const itemsBasePath = "item/v7"

// ItemsService is an interface for interfacing with the Items endpoints of the Too Good To Go API.
type ItemsService interface {
	List(context.Context, *ListItemsRequest) (*ListItemsResponse, *http.Response, error)
	Get(context.Context, *GetItemRequest, string) (*GetItemResponse, *http.Response, error)
	Favorite(context.Context, *FavoriteItemRequest, string) (*http.Response, error)
}

// ItemsServiceOp handles communication with the Items related methods of the Too Good To Go API.
type ItemsServiceOp struct {
	client *Client
}

var _ ItemsService = &ItemsServiceOp{}

// List handles listing all items.
func (s *ItemsServiceOp) List(ctx context.Context, listItemsRequest *ListItemsRequest) (*ListItemsResponse, *http.Response, error) {
	if listItemsRequest == nil {
		return nil, nil, NewArgumentError("listItemsRequest", "must not be nil")
	}

	if listItemsRequest.UserID == "" {
		if s.client.UserID == "" {
			return nil, nil, NewArgumentError("listItemsRequest.UserID", "must not be nil - client has no user id set - please log in using Auth service first or provide UserID in listItemsRequest")
		}
		// if UserID not passed, but we have already authenticated with API, use user id stored in client.
		listItemsRequest.UserID = s.client.UserID
	}

	url := fmt.Sprintf("%s/", itemsBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, listItemsRequest)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.AccessToken))

	listItemsResponse := &ListItemsResponse{}
	response, err := s.client.Do(ctx, req, listItemsResponse)
	if err != nil {
		return nil, nil, err
	}

	return listItemsResponse, response, nil
}

// Get handles getting specific information about particular item entry.
func (s *ItemsServiceOp) Get(ctx context.Context, getItemRequest *GetItemRequest, itemID string) (*GetItemResponse, *http.Response, error) {
	if itemID == "" {
		return nil, nil, NewArgumentError("itemID", "must not be nil")
	}

	if getItemRequest == nil {
		return nil, nil, NewArgumentError("getItemRequest", "must not be nil")
	}

	if getItemRequest.UserID == "" {
		if s.client.UserID == "" {
			return nil, nil, NewArgumentError("getItemRequest.UserID", "must not be nil - client has no user id set - please log in using Auth service first or provide UserID in getItemRequest")
		}
		// if UserID not passed, but we have already authenticated with API, use user id stored in client.
		getItemRequest.UserID = s.client.UserID
	}

	url := fmt.Sprintf("%s/%s", itemsBasePath, itemID)
	req, err := s.client.NewRequest(http.MethodPost, url, getItemRequest)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.AccessToken))

	getItemResp := &GetItemResponse{}
	response, err := s.client.Do(ctx, req, getItemResp)
	if err != nil {
		return nil, nil, err
	}

	return getItemResp, response, nil
}

// Favorite handles setting particular item entry as favorite.
func (s *ItemsServiceOp) Favorite(ctx context.Context, favoriteItemRequest *FavoriteItemRequest, itemID string) (*http.Response, error) {
	if itemID == "" {
		return nil, NewArgumentError("itemID", "must not be nil")
	}

	if favoriteItemRequest == nil {
		return nil, NewArgumentError("favoriteItemRequest", "must not be nil")
	}

	url := fmt.Sprintf("%s/%s/setFavorite", itemsBasePath, itemID)
	req, err := s.client.NewRequest(http.MethodPost, url, favoriteItemRequest)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.AccessToken))

	response, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
