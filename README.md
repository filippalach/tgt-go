# Too Good To Go API client (tgt-go)

[![Coverage](https://codecov.io/gh/filippalach/tgt-go/branch/master/graphs/badge.svg?branch=main)](https://codecov.io/gh/filippalach/tgt-go)
[![Build Status](https://github.com/filippalach/tgt-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/filippalach/tgt-go/actions)
[![GoDoc](https://godoc.org/github.com/filippalach/tgt-go?status.svg)](https://godoc.org/github.com/filippalach/tgt-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nwillc/snipgo)](https://goreportcard.com/report/github.com/filippalach/tgt-go/)

tgt-go is an unofficial Go client library for accessing the Too Good To Go API.

No official Too Good To Go API documentation is available.
<br></br>

## Install
```sh
go get github.com/filippalach/tgt-go@vX.Y.Z
```

where X.Y.Z is the [version](https://github.com/filippalach/tgt-go/releases) you need.

or
```sh
go get github.com/filippalach/tgt-go
```
for non Go modules usage or latest version.
<br></br>

## Usage

Create a new Too Good To Go client, then use the exposed services (Auth, Items, Orders) to
access different parts of the API.
<br></br>

### Capabilities

This package provides following functionalities through exposed services:

<ul>
  <li>Auth Service</li>
    <ul>
      <li>Login - initiate auth process - /auth/vX/authByEmail</li>
      <li>Poll - finish auth process - /auth/vX/authByRequestPollingId</li>
      <li>Refresh - refresh tokens - /auth/vX/token/refresh</li>
      <li>Signup - create account - /auth/vX/signUpByEmail</li>
    </ul>
  </li>
  <li>Items Service</li>
    <ul>
      <li>List - fetch items - /items/vX/</li>
      <li>Get - fetch specific item - /items/vX/{item_id} </li>
      <li>Favorite - un/set specific item as favorite - /items/vX/{item_id}/setFavorite </li>
    </ul>
  </li>
  <li>Orders Service</li>
    <ul>
      <li>Active - fetch active orders - /order/vX/active</li>
      <li>Inactive - fetch past/inactive orders - /order/vX/inactive</li>
    </ul>
  </li>
</ul>

Apart from that, exported methods such as: SetAuthContext, NewRequest, Do, CheckResponseForErrors can be used to form request from scratch, if service capabilites would happen to be insufficient in any case.
<br></br>

### Authentication

Too Good To Go application uses custom authentication process. It consists of 3 steps:

* Initiating Login process with email (Login method of Auth service)
* Clicking link in email received (requires manual intervention)
* Obtaining tokens (Poll method of Auth service)

After successful authentication process with client's Auth service, Too Good To Go auth context will be set
and used in every subsequent request, consisting of
* access_token (needed for subsequent requests using Items, Orders service)
* refresh_token (needed for refreshing tokens)
* user_id (needed for subsequent requests using Items, Orders service)
<br></br>

## Usage example

```go
package main

import (
	"context"
	"log"
	"os"
	"time"

	tgtg "github.com/filippalach/tgt-go"
)

func main() {
	// Create client with custom user agent and headers. Pass no ClentOpts to get default client.
	client, err := tgtg.New(
		nil,
		tgtg.SetUserAgent("<your custom user agent>"),
		tgtg.SetRequestHeaders(map[string]string{"<header>": "<value>"}))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Initiate login process with login request.
	loginReq := &tgtg.LoginRequest{
		DeviceType: "<IOS|ANDROID>",
		Email:      "<your_email>",
	}
	loginResp, _, err := client.Auth.Login(context.TODO(), loginReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Poll to obtain Too Good To Go tokens.
	pollReq := &tgtg.PollRequest{
		DeviceType: "<IOS|ANDROID>",
		Email:      "<your_email>",

		// Polling ID is returned by Login request.
		PollingID: loginResp.PollingID,
	}

	// Give client time for you to finish authentication process via received email.
	log.Println("Check your mailbox, finish login process and wait 2 minutes...")
	time.Sleep(2 * time.Minute)
	pollResp, _, err := client.Auth.Poll(context.TODO(), pollReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Refresh tokens - nil body - client's refresh token will be used by default
	// since they already got set earlier with sucessful Poll method.
	refreshResp, _, err := client.Auth.Refresh(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Refresh tokens - with body - specific refresh token will be used.
	refreshReq := &tgtg.RefreshTokensRequest{
		RefreshToken: refreshResp.RefreshToken,
	}
	_, _, err = client.Auth.Refresh(context.TODO(), refreshReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// List items.
	listItemsReq := &tgtg.ListItemsRequest{
		// UserID has to be specified if no SetAuthContext was called
		// either manually, or through the Auth service usage.
		UserID:   pollResp.StartupData.User.UserID,
		Radius:   10,
		PageSize: 100,
		Page:     1,
		Origin: &tgtg.Origin{
			Latitude:  50.000,
			Longitude: 20.000,
		},
	}
	listItemsResp, _, err := client.Items.List(context.TODO(), listItemsReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Get speicifc's item ID.
	itemID := listItemsResp.Items[0].Item.ItemID

	// Get specific item
	itemReq := &tgtg.GetItemRequest{
		// UserID has to be specified if no SetAuthContext was called either
		// manually, or through the Auth service usage.
		UserID: pollResp.StartupData.User.UserID,
	}
	_, _, err = client.Items.Get(context.TODO(), itemReq, itemID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// ...process received item.

	// Un/Set specific item as favorite.
	favItemReq := &tgtg.FavoriteItemRequest{
		IsFavorite: true,
	}
	_, err = client.Items.Favorite(context.TODO(), favItemReq, itemID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// ...process received response.

	// Check active orders.
	activeOrdersReq := &tgtg.ActiveOrdersRequest{
		// UserID has to be specified if no SetAuthContext was called either
		// manually, or through the Auth service usage.
		UserID: pollResp.StartupData.User.UserID,
	}
	_, _, err = client.Orders.Active(context.TODO(), activeOrdersReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// ...process received response.

	// Check all inactive/past orders.
	inactiveOrdersReq := &tgtg.InactiveOrdersRequest{
		// UserID has to be specified if no SetAuthContext was called either
		// manually, or through the Auth service usage.
		UserID: pollResp.StartupData.User.UserID,
		Paging: tgtg.Paging{
			// 0 for all - this API endpoint works really weirdly.
			Page: 0,
			Size: 200,
		},
	}
	_, _, err = client.Orders.Inactive(context.TODO(), inactiveOrdersReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// ...process received response.

	// Signup request.
	signupReq := &tgtg.SignupRequest{
		CountryID:             "<ISO 3166-1 APLHA-2 country code - e.x. GB, IT>",
		DeviceType:            "<ANDROID|IOS>",
		Email:                 "<your_email>",
		Name:                  "<your_name>",
		NewsletterOptIn:       false,
		PushNotificationOptIn: true,
	}
	_, _, err = client.Auth.Signup(context.TODO(), signupReq)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// ...process received response.
}
```

## Versioning

Each version of the client is tagged and the version is updated accordingly.

To see the list of past versions, run `git tag`.
<br></br>

## Documentation

For details on all the functionality in this library, see the [GoDoc](https://godoc.org/github.com/filippalach/tgt-go) documentation.
<br></br>

## Contributions

Too Good To Go App is subject of constant changes - some might be tricky to track - that's why we love pull requests!