package model

// This file defines request/response models for the EV Charging HTTP APIs
// described in swagger.yaml. The struct and field names closely follow the
// OpenAPI component schemas while keeping idiomatic Go naming.

// Error represents the standard error response schema.
type Error struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// --- Shared value objects used across multiple APIs ---

type AvailabilityWindow struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type Rating struct {
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

type Address struct {
	Name           string    `json:"name"`
	GeoCoordinates []float64 `json:"geo_coordinates"`
}

type ProviderDescriptor struct {
	Name string `json:"name"`
}

type Provider struct {
	ID         string             `json:"id"`
	Descriptor ProviderDescriptor `json:"descriptor"`
	Address    Address            `json:"address"`
}

type ConnectorAttributes struct {
	ConnectorType        string   `json:"connectorType"`
	MaxPowerKW           float64  `json:"maxPowerKW"`
	MinPowerKW           float64  `json:"minPowerKW"`
	SocketCount          int      `json:"socketCount"`
	ReservationSupported bool     `json:"reservationSupported"`
	OcppID               string   `json:"ocppId,omitempty"`
	EvseID               string   `json:"evseId,omitempty"`
	ParkingType          string   `json:"parkingType,omitempty"`
	ConnectorID          string   `json:"connectorId,omitempty"`
	PowerType            string   `json:"powerType"`
	ConnectorFormat      string   `json:"connectorFormat"`
	ChargingSpeed        string   `json:"chargingSpeed"`
	Status               string   `json:"status"`
	AmenityFeature       []string `json:"amenityFeature,omitempty"`
	RoamingNetwork       string   `json:"roamingNetwork,omitempty"`
}

type Connector struct {
	ID                  string               `json:"id"`
	IsActive            bool                 `json:"isActive"`
	AvailabilityWindow  []AvailabilityWindow `json:"availabilityWindow,omitempty"`
	ConnectorAttributes ConnectorAttributes  `json:"connectorAttributes"`
}

type BuyerFinderFee struct {
	FeeType  string  `json:"feeType"`
	FeeValue float64 `json:"feeValue"`
}

type OfferAttributes struct {
	BuyerFinderFee *BuyerFinderFee `json:"buyerFinderFee,omitempty"`
	IdleFeePolicy  string          `json:"idleFeePolicy,omitempty"`
}

type ApplicableQuantity struct {
	UnitText     string  `json:"unitText"`
	UnitCode     string  `json:"unitCode"`
	UnitQuantity float64 `json:"unitQuantity"`
}

type Price struct {
	Currency           string              `json:"currency"`
	Value              float64             `json:"value"`
	ApplicableQuantity *ApplicableQuantity `json:"applicableQuantity,omitempty"`
}

type Validity struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type OfferDescriptor struct {
	Name string `json:"name"`
}

type Offer struct {
	ID                    string           `json:"id"`
	Descriptor            OfferDescriptor  `json:"descriptor"`
	Items                 []string         `json:"items"`
	Price                 Price            `json:"price"`
	Validity              *Validity        `json:"validity,omitempty"`
	AcceptedPaymentMethod []string         `json:"acceptedPaymentMethod,omitempty"`
	OfferAttributes       *OfferAttributes `json:"offerAttributes,omitempty"`
	Provider              string           `json:"provider"`
}

type Catalog struct {
	ID         string      `json:"id"`
	Provider   Provider    `json:"provider"`
	Rating     *Rating     `json:"rating,omitempty"`
	Connectors []Connector `json:"connectors"`
	Offers     []Offer     `json:"offers"`
}

// --- Search API models ---

type Vehicle struct {
	Make  string `json:"make,omitempty"`
	Model string `json:"model,omitempty"`
	Type  string `json:"type,omitempty"`
}

type SearchFilters struct {
	CPO           string   `json:"cpo,omitempty"`
	ConnectorType string   `json:"connector_type,omitempty"`
	MaxPowerKW    *float64 `json:"max_power_kw,omitempty"`
	Amenities     []string `json:"amenities,omitempty"`
	Vehicle       *Vehicle `json:"vehicle,omitempty"`
}

type SearchSort struct {
	SortBy string `json:"sort_by,omitempty"`
	Order  string `json:"order,omitempty"`
}

type TimeWindow struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// DiscoverRequest represents the full discover request format with context and message
type DiscoverRequest struct {
	Context DiscoverContext `json:"context"`
	Message DiscoverMessage `json:"message"`
}

// DiscoverContext represents the context in a discover request
type DiscoverContext struct {
	TransactionID string `json:"transaction_id"`
	Version       string `json:"version"`
	Action        string `json:"action"`
	Domain        string `json:"domain"`
	BapID         string `json:"bap_id"`
	BapURI        string `json:"bap_uri"`
	MessageID     string `json:"message_id"`
	Timestamp     string `json:"timestamp"`
	TTL           string `json:"ttl"`
}

// DiscoverMessage represents the message payload in a discover request
type DiscoverMessage struct {
	Geometry       *Geometry `json:"geometry,omitempty"`
	DistanceMeters float64   `json:"distance_meters,omitempty"`
}

// Geometry represents a geometry object
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type SearchRequest struct {
	EvseID         string         `json:"evse_id,omitempty"`
	GeoCoordinates []float64      `json:"geo_coordinates,omitempty"`
	DistanceMeters float64        `json:"distance_meters,omitempty"`
	TimeWindow     *TimeWindow    `json:"time_window,omitempty"`
	Filters        *SearchFilters `json:"filters,omitempty"`
	Sort           *SearchSort    `json:"sort,omitempty"`
}

type SearchResponse struct {
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
	Catalogs []Catalog `json:"catalogs"`
}

// --- Estimate API models ---

type Energy struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type EstimateRequest struct {
	EvseID      string      `json:"evse_id"`
	Vehicle     Vehicle     `json:"vehicle"`
	ConnectorID string      `json:"connector_id"`
	TimeWindow  *TimeWindow `json:"time_window,omitempty"`
	Energy      *Energy     `json:"energy,omitempty"`
	OfferID     string      `json:"offer_id,omitempty"`
	Amount      *Amount     `json:"amount,omitempty"`
}

type OrderInfo struct {
	ID     string `json:"id"`
	Mode   string `json:"mode"`
	Status string `json:"status"`
}

type PaymentInfo struct {
	Status string `json:"status"`
}

type ChargingInfo struct {
	Status string `json:"status"`
}

type CancellationFee struct {
	Percentage string `json:"percentage"`
}

type ExternalRef struct {
	MIMEType string `json:"mimetype"`
	URL      string `json:"url"`
}

type CancellationPolicy struct {
	Fee         *CancellationFee `json:"fee,omitempty"`
	ExternalRef *ExternalRef     `json:"externalRef,omitempty"`
}

type PriceComponent struct {
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
}

type EstimateResponse struct {
	Order                      OrderInfo           `json:"order"`
	Amount                     Amount              `json:"amount"`
	DurationInMinutes          string              `json:"durationInMinutes"`
	PercentageOfBatteryCharged string              `json:"percentageOfBatteryCharged"`
	Energy                     *Energy             `json:"energy,omitempty"`
	Validity                   *Validity           `json:"validity,omitempty"`
	PriceComponents            []PriceComponent    `json:"priceComponents,omitempty"`
	Cancellation               *CancellationPolicy `json:"cancellation,omitempty"`
}

// --- Order and charging models ---

type ChargingMetric struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

type ChargingTelemetry struct {
	EventTime string           `json:"eventTime"`
	Metrics   []ChargingMetric `json:"metrics"`
}

type OrderResponse struct {
	Order             OrderInfo          `json:"order"`
	Payment           *PaymentInfo       `json:"payment,omitempty"`
	Charging          *ChargingInfo      `json:"charging,omitempty"`
	ConnectorID       string             `json:"connectorId,omitempty"`
	ConnectorType     string             `json:"connectorType,omitempty"`
	Vehicle           *Vehicle           `json:"vehicle,omitempty"`
	TrackingURL       string             `json:"trackingUrl,omitempty"`
	ChargingTelemetry *ChargingTelemetry `json:"chargingTelemetry,omitempty"`
}

type PriceComponentString struct {
	Type        string `json:"type"`
	Value       string `json:"value"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type PriceComponentFlexible struct {
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	Currency    string      `json:"currency"`
	Description string      `json:"description"`
}

type CancelEstimateResponse struct {
	Order           OrderInfo              `json:"order"`
	Payment         *PaymentInfo           `json:"payment,omitempty"`
	Charging        *ChargingInfo          `json:"charging,omitempty"`
	Validity        *Validity              `json:"validity,omitempty"`
	PriceComponents []PriceComponentString `json:"priceComponents,omitempty"`
}

type CancelResponse struct {
	Order           OrderInfo              `json:"order"`
	Payment         *PaymentInfo           `json:"payment,omitempty"`
	Charging        *ChargingInfo          `json:"charging,omitempty"`
	PriceComponents []PriceComponentString `json:"priceComponents,omitempty"`
}

type StartChargingRequest struct{}

type StartChargingResponse struct {
	Order    OrderInfo     `json:"order"`
	Payment  *PaymentInfo  `json:"payment,omitempty"`
	Charging *ChargingInfo `json:"charging,omitempty"`
}

type StopEstimateResponse struct {
	Order           OrderInfo              `json:"order"`
	Payment         *PaymentInfo           `json:"payment,omitempty"`
	Charging        *ChargingInfo          `json:"charging,omitempty"`
	Validity        *Validity              `json:"validity,omitempty"`
	PriceComponents []PriceComponentString `json:"priceComponents,omitempty"`
}

type StopChargingRequest struct {
	ReasonCode string `json:"reasonCode,omitempty"`
	Message    string `json:"message,omitempty"`
}

type StopChargingResponse struct {
	Order           OrderInfo                `json:"order"`
	Payment         *PaymentInfo             `json:"payment,omitempty"`
	Charging        *ChargingInfo            `json:"charging,omitempty"`
	Validity        *Validity                `json:"validity,omitempty"`
	PriceComponents []PriceComponentFlexible `json:"priceComponents,omitempty"`
}

// --- Feedback and support models ---

type Feedback struct {
	Comments string   `json:"comments,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type RatingRequest struct {
	Value    int       `json:"value"`
	Feedback *Feedback `json:"feedback,omitempty"`
}

type FeedbackForm struct {
	URL          string `json:"url"`
	MIMEType     string `json:"mimeType"`
	SubmissionID string `json:"submissionId"`
}

type RatingResponse struct {
	Order        OrderInfo     `json:"order"`
	FeedbackForm *FeedbackForm `json:"feedbackForm,omitempty"`
}

type SupportResponse struct {
	Order    OrderInfo `json:"order"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	URL      string    `json:"url"`
	Hours    string    `json:"hours"`
	Channels []string  `json:"channels"`
}

// --- Payments ---

type PaymentResponse struct {
	Order                 OrderInfo `json:"order"`
	Amount                Amount    `json:"amount"`
	BeneficiaryID         string    `json:"beneficiaryId"`
	AcceptedPaymentMethod []string  `json:"acceptedPaymentMethod"`
	PaymentURL            string    `json:"paymentUrl"`
	Validity              *Validity `json:"validity,omitempty"`
}
