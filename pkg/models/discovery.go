package models

// DiscoveryRequest represents the JSON request structure
type DiscoveryRequest struct {
	Context Context `json:"context"`
	Message Message `json:"message"`
}

// Context represents the context in the request
type Context struct {
	Version       string   `json:"version"`
	Action        string   `json:"action"`
	Domain        string   `json:"domain"`
	Location      Location `json:"location"`
	BapID         string   `json:"bap_id"`
	BapURI        string   `json:"bap_uri"`
	BppID         string   `json:"bpp_id"`
	BppURI        string   `json:"bpp_uri"`
	TransactionID string   `json:"transaction_id"`
	MessageID     string   `json:"message_id"`
	Timestamp     string   `json:"timestamp"`
	TTL           string   `json:"ttl"`
	SchemaContext []string `json:"schema_context"`
}

// Location represents location information
type Location struct {
	Country Country `json:"country"`
	City    City    `json:"city"`
}

// Country represents country information
type Country struct {
	Code string `json:"code"`
}

// City represents city information
type City struct {
	Code string `json:"code"`
}

// Message represents the message payload
type Message struct {
	Catalogs []Catalog `json:"catalogs"`
}

// Catalog represents a catalog of items
type Catalog struct {
	Context    string     `json:"@context"`
	Type       string     `json:"@type"`
	Descriptor Descriptor `json:"beckn:descriptor"`
	Provider   Provider   `json:"beckn:provider"`
	Items      []Item     `json:"beckn:items"`
}

// Descriptor represents descriptor information
type Descriptor struct {
	Name      string `json:"name"`
	ShortDesc string `json:"short_desc,omitempty"`
}

// Provider represents provider information
type Provider struct {
	Context    string     `json:"@context"`
	Type       string     `json:"@type"`
	ID         string     `json:"beckn:id"`
	Descriptor Descriptor `json:"beckn:descriptor"`
	Address    Address    `json:"beckn:address"`
	Contact    []Contact  `json:"beckn:contact"`
}

// Address represents address information
type Address struct {
	AreaCode string  `json:"area_code,omitempty"`
	City     City    `json:"city,omitempty"`
	State    State   `json:"state,omitempty"`
	Country  Country `json:"country,omitempty"`
	Full     string  `json:"full,omitempty"`
}

// State represents state information
type State struct {
	Code string `json:"code"`
}

// Contact represents contact information
type Contact struct {
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
}

// Item represents an item in the catalog
type Item struct {
	Context        string                 `json:"@context"`
	Type           string                 `json:"@type"`
	ID             string                 `json:"beckn:id"`
	Descriptor     Descriptor             `json:"beckn:descriptor"`
	Category       Category               `json:"beckn:category"`
	ItemAttributes map[string]interface{} `json:"beckn:itemAttributes"`
	AvailableAt    []LocationItem         `json:"beckn:availableAt"`
	Offers         []Offer                `json:"beckn:offers"`
	Rating         *Rating                `json:"beckn:rating,omitempty"`
}

// Category represents category information
type Category struct {
	Type       string     `json:"@type"`
	ID         string     `json:"id"`
	Descriptor Descriptor `json:"descriptor"`
}

// LocationItem represents location for an item
type LocationItem struct {
	Type    string  `json:"@type"`
	Geo     Geo     `json:"geo,omitempty"`
	Address Address `json:"address,omitempty"`
}

// Geo represents geographic coordinates
type Geo struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Offer represents an offer
type Offer struct {
	Context               string                 `json:"@context"`
	Type                  string                 `json:"@type"`
	ID                    string                 `json:"beckn:id"`
	Descriptor            Descriptor             `json:"beckn:descriptor"`
	Price                 Price                  `json:"beckn:price"`
	Validity              Validity               `json:"beckn:validity"`
	AcceptedPaymentMethod []string               `json:"beckn:acceptedPaymentMethod"`
	OfferAttributes       map[string]interface{} `json:"beckn:offerAttributes,omitempty"`
}

// Price represents pricing information
type Price struct {
	Currency           string             `json:"currency"`
	Value              float64            `json:"value"`
	ApplicableQuantity ApplicableQuantity `json:"applicableQuantity"`
}

// ApplicableQuantity represents quantity information
type ApplicableQuantity struct {
	UnitText     string `json:"unitText"`
	UnitCode     string `json:"unitCode"`
	UnitQuantity int    `json:"unitQuantity"`
}

// Validity represents validity period
type Validity struct {
	Type      string `json:"@type"`
	StartDate string `json:"schema:startDate"`
	EndDate   string `json:"schema:endDate"`
}

// Rating represents rating information
type Rating struct {
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

// DiscoveryResponse represents the JSON response structure
type DiscoveryResponse struct {
	Context Context `json:"context"`
	Message Message `json:"message"`
}

