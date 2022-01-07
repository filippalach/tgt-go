package tgtg

import "time"

// GetItemRequest represents a request body to get item details.
type GetItemRequest struct {
	UserID string  `json:"user_id"`
	Origin *Origin `json:"origin"`
}

// GetItemResponse represents a response body with detailed item info.
type GetItemResponse struct {
	Item  Item  `json:"item"`
	Store Store `json:"store"`

	DisplayName    string         `json:"display_name"`
	Distance       float64        `json:"distance"`
	Favorite       bool           `json:"favorite"`
	InSalesWindow  bool           `json:"in_sales_window"`
	ItemsAvailable int            `json:"items_available"`
	NewItem        bool           `json:"new_item"`
	PickupInterval PickupInterval `json:"pickup_interval"`
	PickupLocation PickupLocation `json:"pickup_location"`
	PurchaseEnd    string         `json:"purchase_end"`
	SharingURL     string         `json:"sharing_url"`
}

// ListItemsRequest represents a request body to list all items.
type ListItemsRequest struct {
	PageSize int `json:"page_size"`
	Page     int `json:"page"`

	UserID string `json:"user_id"`

	Radius int     `json:"radius"`
	Origin *Origin `json:"origin"`

	ItemCategories []string `json:"item_categories"`
	DietCategories []string `json:"diet_categories"`
	PickupEarliest string   `json:"pickup_earliest"`
	PickupLatest   string   `json:"pickup_latest"`
	SearchPhrase   string   `json:"search_phrase"`

	Discover      bool `json:"discover"`
	FavoritesOnly bool `json:"favorites_only"`
	WithStockOnly bool `json:"with_stock_only"`
	HiddenOnly    bool `json:"hidden_only"`
	WeCareOnly    bool `json:"we_care_only"`
}

// ListItemsResponse represents a response body with detailed items list.
type ListItemsResponse struct {
	Items []Items `json:"items"`
}

// FavoriteItemRequest represents a request body to un/set favorite item.
type FavoriteItemRequest struct {
	IsFavorite bool `json:"is_favorite"`
}

// Item represents a Too Good To Go Item details.
type Item struct {
	AverageOverallRating   AverageOverallRating `json:"average_overall_rating"`
	Badges                 []Badges             `json:"badges"`
	Buffet                 bool                 `json:"buffet"`
	CanUserSupplyPackaging bool                 `json:"can_user_supply_packaging"`
	CollectionInfo         string               `json:"collection_info"`
	CoverPicture           Picture              `json:"cover_picture"`
	Description            string               `json:"description"`
	DietCategories         []string             `json:"diet_categories"`
	FavoriteCount          int                  `json:"favorite_count"`
	ItemCategory           string               `json:"item_category"`
	ItemID                 string               `json:"item_id"`
	LogoPicture            Picture              `json:"logo_picture"`
	Name                   string               `json:"name"`
	PackagingOption        string               `json:"packaging_option"`
	PositiveRatingReasons  []string             `json:"positive_rating_reasons"`
	Price                  Price                `json:"price"`
	PriceExcludingTaxes    Price                `json:"price_excluding_taxes"`
	PriceIncludingTaxes    Price                `json:"price_including_taxes"`
	SalesTaxes             []SalesTaxes         `json:"sales_taxes"`
	ShowSalesTaxes         bool                 `json:"show_sales_taxes"`
	TaxAmount              Price                `json:"tax_amount"`
	TaxationPolicy         string               `json:"taxation_policy"`
	ValueExcludingTaxes    Price                `json:"value_excluding_taxes"`
	ValueIncludingTaxes    Price                `json:"value_including_taxes"`
}

// Items represents a Too Good To Go Items details.
type Items struct {
	Item  Item  `json:"item"`
	Store Store `json:"store"`

	DisplayName    string         `json:"display_name"`
	PickupInterval PickupInterval `json:"pickup_interval"`
	PickupLocation Location       `json:"pickup_location"`
	PurchaseEnd    time.Time      `json:"purchase_end"`
	ItemsAvailable int            `json:"items_available"`
	SoldOutAt      time.Time      `json:"sold_out_at"`
	Distance       float64        `json:"distance"`
	Favorite       bool           `json:"favorite"`
	InSalesWindow  bool           `json:"in_sales_window"`
	NewItem        bool           `json:"new_item"`
}

// AverageOverallRating represents a Too Good To Go AverageOverallRating details.
type AverageOverallRating struct {
	AverageOverallRating float64 `json:"average_overall_rating"`
	MonthCount           int     `json:"month_count"`
	RatingCount          int     `json:"rating_count"`
}

// Badges represents a Too Good To Go Badges details.
type Badges struct {
	BadgeType   string `json:"badge_type"`
	MonthCount  int    `json:"month_count"`
	Percentage  int    `json:"percentage"`
	RatingGroup string `json:"rating_group"`
	UserCount   int    `json:"user_count"`
}

// Picture represents a Too Good To Go Picture details.
type Picture struct {
	CurrentURL string `json:"current_url"`
	PictureID  string `json:"picture_id"`
}

// Price represents a Too Good To Go Price details.
type Price struct {
	Code       string `json:"code"`
	Decimals   int    `json:"decimals"`
	MinorUnits int    `json:"minor_units"`
}

// PickupInterval represents a Too Good To Go Pickup Interval details.
type PickupInterval struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
}

// PickupLocation represents a Too Good To Go Pickup Location details.
type PickupLocation struct {
	Address  Address  `json:"address"`
	Location Location `json:"location"`
}

// Address represents a Too Good To Go Address details.
type Address struct {
	AddressLine string  `json:"address_line"`
	City        string  `json:"city"`
	Country     Country `json:"country"`
	PostalCode  string  `json:"postal_code"`
}

// Country represents a Too Good To Go Country details.
type Country struct {
	IsoCode string `json:"iso_code"`
	Name    string `json:"name"`
}

// Location represents a Too Good To Go Location details which has same structure as Origin.
type Location Origin

// Store represents a Too Good To Go Store details.
type Store struct {
	Branch        string        `json:"branch"`
	CoverPicture  Picture       `json:"cover_picture"`
	Description   string        `json:"description"`
	Distance      float64       `json:"distance"`
	FavoriteCount int           `json:"favorite_count"`
	Hidden        bool          `json:"hidden"`
	Items         []StoreItems  `json:"items"`
	LogoPicture   Picture       `json:"logo_picture"`
	Milestones    []Milestones  `json:"milestones"`
	StoreID       string        `json:"store_id"`
	StoreLocation StoreLocation `json:"store_location"`
	StoreName     string        `json:"store_name"`
	StoreTimeZone string        `json:"store_time_zone"`
	TaxIdentifier string        `json:"tax_identifier"`
	WeCare        bool          `json:"we_care"`
	Website       string        `json:"website"`
}

// StoreItems represents a Too Good To Go Store Items details.
type StoreItems struct {
	DisplayName    string         `json:"display_name"`
	Distance       float64        `json:"distance"`
	Favorite       bool           `json:"favorite"`
	InSalesWindow  bool           `json:"in_sales_window"`
	Item           Item           `json:"item"`
	ItemsAvailable int            `json:"items_available"`
	NewItem        bool           `json:"new_item"`
	PickupInterval PickupInterval `json:"pickup_interval"`
	PickupLocation PickupLocation `json:"pickup_location"`
	PurchaseEnd    string         `json:"purchase_end"`
}

// Milestones represents a Too Good To Go Milestones details.
type Milestones struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// StoreLocation represents a Too Good To Go Store Location details which has same structure as PickupLocation.
type StoreLocation PickupLocation

// Origin represents a Too Good To Go Origin details.
type Origin struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
