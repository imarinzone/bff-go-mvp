package search

import (
	"context"

	"bff-go-mvp/internal/model"
)

// MockService implements Service and returns static data that matches
// the example in swagger.yaml for POST /v1/search.
type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) Search(ctx context.Context, page, perPage int, req model.SearchRequest) (model.SearchResponse, error) {
	// Static mock response matching the desired structure.
	resp := model.SearchResponse{
		Total:   1,
		Page:    page,
		PerPage: perPage,
		Catalogs: []model.Catalog{
			{
				ID: "catalog-ev-charging-001",
				Provider: model.Provider{
					ID: "ecopower-charging",
					Descriptor: model.ProviderDescriptor{
						Name: "EcoPower Charging Pvt Ltd",
					},
				},
				Address: model.Address{
					Name:           "MG JVLR Jogeshwari Caves Road",
					GeoCoordinates: []float64{12.9716, 77.5946},
				},
				Rating: &model.Rating{
					Value: 4.5,
					Count: 128,
				},
				AvailabilityWindow: []model.AvailabilityWindow{
					{
						StartTime: "06:00:00",
						EndTime:   "22:00:00",
					},
				},
				AvailablePowerType: []string{"DC", "AC"},
				Connectors: []model.Connector{
					{
						ID:       "ev-charger-ccs2-001",
						IsActive: true,
						ConnectorAttributes: model.ConnectorAttributes{
							ConnectorType:        "TYPE 2",
							MaxPowerKW:           60,
							MinPowerKW:           5,
							SocketCount:          2,
							ReservationSupported: true,
							Status:               "Available",
							ChargingSpeed:        "FAST",
							PowerType:            "DC",
							ConnectorFormat:      "CABLE",
						},
					},
				},
				Offers: []model.Offer{
					{
						ID: "offer-ccs2-60kw-kwh",
						Descriptor: model.OfferDescriptor{
							Name: "Per-kWh Tariff - CCS2 60kW",
						},
						Items: []string{
							"ev-charger-ccs2-001",
						},
						Price: model.Price{
							Currency: "INR",
							Value:    18,
							ApplicableQuantity: &model.ApplicableQuantity{
								UnitText:     "Kilowatt Hour",
								UnitCode:     "KWH",
								UnitQuantity: 1,
							},
						},
						Validity: &model.Validity{
							StartDate: "2025-10-01T00:00:00Z",
							EndDate:   "2026-03-31T23:59:59Z",
						},
						AcceptedPaymentMethod: []string{"UPI", "Card", "Wallet"},
						OfferAttributes: &model.OfferAttributes{
							BuyerFinderFee: &model.BuyerFinderFee{
								FeeType:  "PERCENTAGE",
								FeeValue: 2.5,
							},
							IdleFeePolicy: "₹2/min after 10 min post-charge",
						},
						Provider: "ecopower-charging",
					},
				},
			},
		},
	}

	return resp, nil
}
