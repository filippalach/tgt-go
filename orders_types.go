package tgtg

import "time"

// ActiveOrdersRequest represents a request body to obtain user's active orders.
type ActiveOrdersRequest struct {
	UserID string `json:"user_id"`
}

// InactiveOrdersRequest represents a request body to obtain user's inactive orders.
type InactiveOrdersRequest struct {
	UserID string `json:"user_id"`
	Paging Paging `json:"paging"`
}

// OrdersResponse represents a response body containing Orders details.
type OrdersResponse struct {
	CurrentTime time.Time `json:"current_time"`
	HasMore     bool      `json:"has_more"`
	Orders      []Order   `json:"orders"`
}

// Paging specifies paging options in Order requests.
type Paging struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// Order represents Too Good To Go Order.
type Order struct {
	OrderID                    string         `json:"order_id"`
	State                      string         `json:"state"`
	CancelUntil                time.Time      `json:"cancel_until"`
	RedeemInterval             PickupInterval `json:"redeem_interval"`
	PickupInterval             PickupInterval `json:"pickup_interval"`
	Quantity                   int            `json:"quantity"`
	PriceIncludingTaxes        Price          `json:"price_including_taxes"`
	PriceExcludingTaxes        Price          `json:"price_excluding_taxes"`
	TotalAppliedTaxes          Price          `json:"total_applied_taxes"`
	SalesTaxes                 []SalesTaxes   `json:"sales_taxes"`
	PickupLocation             PickupLocation `json:"pickup_location"`
	IsRated                    bool           `json:"is_rated"`
	TimeOfPurchase             time.Time      `json:"time_of_purchase"`
	StoreID                    string         `json:"store_id"`
	StoreName                  string         `json:"store_name"`
	StoreBranch                string         `json:"store_branch"`
	StoreLogo                  StoreLogo      `json:"store_logo"`
	ItemID                     string         `json:"item_id"`
	ItemName                   string         `json:"item_name"`
	ItemCoverImage             ItemCoverImage `json:"item_cover_image"`
	IsBuffet                   bool           `json:"is_buffet"`
	CanUserSupplyPackaging     bool           `json:"can_user_supply_packaging"`
	PackagingOption            string         `json:"packaging_option"`
	IsStoreWeCare              bool           `json:"is_store_we_care"`
	CanShowBestBeforeExplainer bool           `json:"can_show_best_before_explainer"`
	ShowSalesTaxes             bool           `json:"show_sales_taxes"`
}

// SalesTaxes represents Too Good To Go Sales Taxes.
type SalesTaxes struct {
	TaxDescription string  `json:"tax_description"`
	TaxPercentage  float64 `json:"tax_percentage"`
	TaxAmount      Price   `json:"tax_amount"`
}

// StoreLogo represents Too Good To Go Store Logo.
type StoreLogo struct {
	PictureID              string `json:"picture_id"`
	CurrentURL             string `json:"current_url"`
	IsAutomaticallyCreated bool   `json:"is_automatically_created"`
}

// ItemCoverImage represents Too Good To Go Cover Image.
type ItemCoverImage struct {
	PictureID              string `json:"picture_id"`
	CurrentURL             string `json:"current_url"`
	IsAutomaticallyCreated bool   `json:"is_automatically_created"`
}
