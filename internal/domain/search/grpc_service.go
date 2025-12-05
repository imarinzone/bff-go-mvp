package search

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bff-go-mvp/internal/grpc"
	"bff-go-mvp/internal/model"
	commonpb "bff-go-mvp/proto/common/gen/common"
	discoverpb "bff-go-mvp/proto/discover/gen/discover"
)

// GRPCService implements search.Service using gRPC client
type GRPCService struct {
	client *grpc.DiscoverClient
	logger *zap.Logger
}

// NewGRPCService creates a new gRPC-based search service
func NewGRPCService(client *grpc.DiscoverClient, logger *zap.Logger) *GRPCService {
	return &GRPCService{
		client: client,
		logger: logger,
	}
}

// Search implements the search.Service interface
func (s *GRPCService) Search(ctx context.Context, page, perPage int, req model.SearchRequest) (model.SearchResponse, error) {
	// Convert model request to proto request
	protoReq := s.modelToProtoRequest(req)

	// Call gRPC service
	protoResp, err := s.client.Discover(ctx, protoReq)
	if err != nil {
		s.logger.Error("gRPC search call failed", zap.Error(err))
		return model.SearchResponse{}, err
	}

	// Convert proto response to model response
	return s.protoToModelResponse(protoResp, page, perPage), nil
}

// modelToProtoRequest converts model.SearchRequest to proto DiscoverRequest
func (s *GRPCService) modelToProtoRequest(req model.SearchRequest) *discoverpb.DiscoverRequest {
	// Generate IDs
	transactionID := uuid.New().String()
	messageID := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Build context
	context := &commonpb.Context{
		Action:        ActionDiscover,
		Domain:        DomainDefault,
		BapId:         BapIDDefault,
		BapUri:        BapURIDefault,
		TransactionId: transactionID,
		MessageId:     messageID,
		Timestamp:     timestamp,
		Ttl:           TTLDefault,
		// bpp_id and bpp_uri left empty as they're not known at request time
	}

	// Build message
	message := &discoverpb.Message{}

	// Convert geo_coordinates to Geometry if provided
	if len(req.GeoCoordinates) == 2 {
		message.Geometry = &discoverpb.Geometry{
			Type:        GeometryTypePoint,
			Coordinates: req.GeoCoordinates,
		}
		message.DistanceMeters = req.DistanceMeters
	}

	return &discoverpb.DiscoverRequest{
		Context: context,
		Message: message,
	}
}

// protoToModelResponse converts proto OnDiscoverResponse to model SearchResponse
func (s *GRPCService) protoToModelResponse(protoResp *discoverpb.OnDiscoverResponse, page, perPage int) model.SearchResponse {
	if protoResp == nil || protoResp.Message == nil {
		return model.SearchResponse{
			Total:    0,
			Page:     page,
			PerPage:  perPage,
			Catalogs: []model.Catalog{},
		}
	}

	catalogs := make([]model.Catalog, 0, len(protoResp.Message.Catalogs))
	for _, protoCatalog := range protoResp.Message.Catalogs {
		catalog := s.protoCatalogToModel(protoCatalog)
		catalogs = append(catalogs, catalog)
	}

	return model.SearchResponse{
		Total:    len(catalogs),
		Page:     page,
		PerPage:  perPage,
		Catalogs: catalogs,
	}
}

// protoCatalogToModel converts proto Catalog to model Catalog
func (s *GRPCService) protoCatalogToModel(protoCatalog *discoverpb.Catalog) model.Catalog {
	catalog := model.Catalog{
		ID:         protoCatalog.Id,
		Connectors: []model.Connector{},
		Offers:     []model.Offer{},
	}

	// Convert provider
	if protoCatalog.GetDescriptor_() != nil {
		catalog.Provider = model.Provider{
			ID: protoCatalog.GetId(), // Using catalog ID as provider ID fallback
			Descriptor: model.ProviderDescriptor{
				Name: protoCatalog.GetDescriptor_().GetName(),
			},
		}
	}

	// Convert rating
	if protoCatalog.Validity != nil {
		// Validity is stored but not directly mapped to model.Catalog
		// Catalog doesn't have validity field in model
	}

	// Convert items to connectors
	for _, protoItem := range protoCatalog.Items {
		connector := s.protoItemToConnector(protoItem)
		catalog.Connectors = append(catalog.Connectors, connector)
	}

	// Convert offers
	for _, protoOffer := range protoCatalog.Offers {
		offer := s.protoOfferToModel(protoOffer)
		catalog.Offers = append(catalog.Offers, offer)
	}

	// Set rating from first item if available
	if len(protoCatalog.Items) > 0 && protoCatalog.Items[0].Rating != nil {
		catalog.Rating = &model.Rating{
			Value: protoCatalog.Items[0].Rating.Value,
			Count: int(protoCatalog.Items[0].Rating.Count),
		}
	}

	return catalog
}

// protoItemToConnector converts proto Item to model Connector
func (s *GRPCService) protoItemToConnector(protoItem *discoverpb.Item) model.Connector {
	connector := model.Connector{
		ID:       protoItem.Id,
		IsActive: protoItem.IsActive,
	}

	// Convert availability windows
	if len(protoItem.AvailabilityWindow) > 0 {
		connector.AvailabilityWindow = make([]model.AvailabilityWindow, 0, len(protoItem.AvailabilityWindow))
		for _, aw := range protoItem.AvailabilityWindow {
			connector.AvailabilityWindow = append(connector.AvailabilityWindow, model.AvailabilityWindow{
				StartTime: aw.StartTime,
				EndTime:   aw.EndTime,
			})
		}
	}

	// Convert item attributes to connector attributes
	if protoItem.ItemAttributes != nil {
		attrs := protoItem.ItemAttributes
		connector.ConnectorAttributes = model.ConnectorAttributes{
			ConnectorType:        attrs.ConnectorType,
			MaxPowerKW:           attrs.MaxPowerKw,
			MinPowerKW:           attrs.MinPowerKw,
			SocketCount:          int(attrs.SocketCount),
			ReservationSupported: attrs.ReservationSupported,
			OcppID:               attrs.OcppId,
			EvseID:               attrs.EvseId,
			ParkingType:          attrs.ParkingType,
			ConnectorID:          attrs.ConnectorId,
			PowerType:            attrs.PowerType,
			ConnectorFormat:      attrs.ConnectorFormat,
			ChargingSpeed:        attrs.ChargingSpeed,
			Status:               attrs.StationStatus,
			AmenityFeature:       attrs.AmenityFeature,
			RoamingNetwork:       attrs.RoamingNetwork,
		}

		// Service location address can be used for provider address if needed
		// Address is stored in Provider.Address, not ConnectorAttributes
	}

	return connector
}

// protoOfferToModel converts proto Offer to model Offer
func (s *GRPCService) protoOfferToModel(protoOffer *discoverpb.Offer) model.Offer {
	offer := model.Offer{
		ID:    protoOffer.Id,
		Items: protoOffer.Items,
	}

	// Convert descriptor
	if protoOffer.GetDescriptor_() != nil {
		offer.Descriptor = model.OfferDescriptor{
			Name: protoOffer.GetDescriptor_().GetName(),
		}
	}

	// Convert price
	if protoOffer.Price != nil {
		offer.Price = model.Price{
			Currency: protoOffer.Price.Currency,
			Value:    protoOffer.Price.Value,
		}
		if protoOffer.Price.ApplicableQuantity != nil {
			offer.Price.ApplicableQuantity = &model.ApplicableQuantity{
				UnitText:     protoOffer.Price.ApplicableQuantity.UnitText,
				UnitCode:     protoOffer.Price.ApplicableQuantity.UnitCode,
				UnitQuantity: float64(protoOffer.Price.ApplicableQuantity.UnitQuantity),
			}
		}
	}

	// Convert validity
	if protoOffer.Validity != nil {
		offer.Validity = &model.Validity{
			StartDate: protoOffer.Validity.StartDate,
			EndDate:   protoOffer.Validity.EndDate,
		}
	}

	// Convert accepted payment methods
	offer.AcceptedPaymentMethod = protoOffer.AcceptedPaymentMethod

	// Convert offer attributes
	if protoOffer.OfferAttributes != nil {
		offer.OfferAttributes = &model.OfferAttributes{}
		if protoOffer.OfferAttributes.BuyerFinderFee != nil {
			offer.OfferAttributes.BuyerFinderFee = &model.BuyerFinderFee{
				FeeType:  protoOffer.OfferAttributes.BuyerFinderFee.FeeType,
				FeeValue: protoOffer.OfferAttributes.BuyerFinderFee.FeeValue,
			}
		}
		offer.OfferAttributes.IdleFeePolicy = protoOffer.OfferAttributes.IdleFeePolicy
	}

	// Set provider
	offer.Provider = protoOffer.Provider

	return offer
}
