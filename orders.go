package tgtg

import (
	"context"
	"fmt"
	"net/http"
)

const ordersBasePath = "order/v6"

// OrdersService is an interface for interfacing with the Orders endpoints of the Too Good To Go API.
type OrdersService interface {
	Active(context.Context, *ActiveOrdersRequest) (*OrdersResponse, *http.Response, error)
	Inactive(context.Context, *InactiveOrdersRequest) (*OrdersResponse, *http.Response, error)
}

// OrdersServiceOp handles communication with the Orders related methods of the Too Good To Go API.
type OrdersServiceOp struct {
	client *Client
}

var _ OrdersService = &OrdersServiceOp{}

// Active handles getting active orders for given user.
func (s *OrdersServiceOp) Active(ctx context.Context, activeOrdersRequest *ActiveOrdersRequest) (*OrdersResponse, *http.Response, error) {
	if activeOrdersRequest == nil {
		if s.client.UserID == "" {
			return nil, nil, NewArgumentError("activeOrdersRequest.UserID", "must not be nil - client has no UserID")
		}
		// if UserID was not passed, but we have already authenticated with API, use UserID stored in client.
		activeOrdersRequest = &ActiveOrdersRequest{UserID: s.client.UserID}
	}

	url := fmt.Sprintf("%s/active", ordersBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, activeOrdersRequest)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.AccessToken))

	activeOrdersResponse := &OrdersResponse{}
	response, err := s.client.Do(ctx, req, activeOrdersResponse)
	if err != nil {
		return nil, nil, err
	}

	return activeOrdersResponse, response, nil
}

// Inactive handles getting inactive/past orders for given user.
func (s *OrdersServiceOp) Inactive(ctx context.Context, inactiveOrdersRequest *InactiveOrdersRequest) (*OrdersResponse, *http.Response, error) {
	if inactiveOrdersRequest == nil {
		return nil, nil, NewArgumentError("inactiveOrdersRequest", "must not be nil")
	}

	if inactiveOrdersRequest.UserID == "" {
		if s.client.UserID == "" {
			return nil, nil, NewArgumentError("inactiveOrdersRequest.UserID", "must not be nil - client has no user id set - please log in using Auth service first or provide UserID in inactiveOrdersRequest")
		}
		// if UserID not passed, but we have already authenticated with API, use user id stored in client.
		inactiveOrdersRequest.UserID = s.client.UserID
	}

	url := fmt.Sprintf("%s/inactive", ordersBasePath)
	req, err := s.client.NewRequest(http.MethodPost, url, inactiveOrdersRequest)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.AccessToken))

	inactiveOrdersResponse := &OrdersResponse{}
	response, err := s.client.Do(ctx, req, inactiveOrdersResponse)
	if err != nil {
		return nil, nil, err
	}

	return inactiveOrdersResponse, response, nil
}
