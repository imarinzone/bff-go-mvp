package search

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bff-go-mvp/internal/grpc"
	"bff-go-mvp/internal/model"
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
	// Extract context from the request context if available
	var discoverContext *model.DiscoverContext
	if ctxVal := ctx.Value("discover_context"); ctxVal != nil {
		if dc, ok := ctxVal.(*model.DiscoverContext); ok {
			discoverContext = dc
		}
	}

	// Convert model request to proto request
	protoReq := s.modelToProtoRequest(req, discoverContext)

	// Log the gRPC request
	s.logger.Info("Sending gRPC discover request to downstream",
		zap.String("service", "discover"),
		zap.String("request", protoReq.String()),
	)

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
func (s *GRPCService) modelToProtoRequest(req model.SearchRequest, discoverContext *model.DiscoverContext) *discoverpb.DiscoverRequest {
	// Use provided context values or generate defaults
	var version, action, domain, bapID, bapURI, transactionID, messageID, timestamp, ttl string

	if discoverContext != nil {
		version = discoverContext.Version
		action = discoverContext.Action
		domain = discoverContext.Domain
		bapID = discoverContext.BapID
		bapURI = discoverContext.BapURI
		transactionID = discoverContext.TransactionID
		messageID = discoverContext.MessageID
		timestamp = discoverContext.Timestamp
		ttl = discoverContext.TTL
	}

	// Use defaults if not provided
	if version == "" {
		version = VersionDefault
	}
	if action == "" {
		action = ActionDiscover
	}
	if domain == "" {
		domain = DomainDefault
	}
	if bapID == "" {
		bapID = BapIDDefault
	}
	if bapURI == "" {
		bapURI = BapURIDefault
	}
	if transactionID == "" {
		transactionID = uuid.New().String()
	}
	if messageID == "" {
		messageID = uuid.New().String()
	}
	if timestamp == "" {
		timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if ttl == "" {
		ttl = TTLDefault
	}

	// Build context
	context := &discoverpb.Context{
		Version:       version,
		Action:        action,
		Domain:        domain,
		BapId:         bapID,
		BapUri:        bapURI,
		TransactionId: transactionID,
		MessageId:     messageID,
		Timestamp:     timestamp,
		Ttl:           ttl,
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

	// Set evse_id if provided
	if req.EvseID != "" {
		message.EvseId = req.EvseID
	}

	// Convert time_window if provided
	if req.TimeWindow != nil {
		message.TimeWindow = &discoverpb.TimeWindow{
			Start: req.TimeWindow.Start,
			End:   req.TimeWindow.End,
		}
	}

	// Convert filters if provided
	if req.Filters != nil {
		filters := &discoverpb.Filters{
			Cpo:           req.Filters.CPO,
			ConnectorType: req.Filters.ConnectorType,
			Amenities:     req.Filters.Amenities,
		}

		// Handle MaxPowerKW
		if req.Filters.MaxPowerKW != nil {
			filters.MaxPowerKw = *req.Filters.MaxPowerKW
		}

		// Convert vehicle if provided
		if req.Filters.Vehicle != nil {
			filters.Vehicle = &discoverpb.Vehicle{
				Make:  req.Filters.Vehicle.Make,
				Model: req.Filters.Vehicle.Model,
				Type:  req.Filters.Vehicle.Type,
			}
		}

		message.Filters = filters
	}

	// Convert sort if provided
	if req.Sort != nil {
		message.Sort = &discoverpb.Sort{
			SortBy: req.Sort.SortBy,
			Order:  req.Sort.Order,
		}
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
	if protoCatalog.GetDescriptorData() != nil {
		catalog.Provider = model.Provider{
			ID: protoCatalog.GetId(), // Using catalog ID as provider ID fallback
			Descriptor: model.ProviderDescriptor{
				Name: protoCatalog.GetDescriptorData().GetName(),
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
	if protoOffer.GetDescriptorData() != nil {
		offer.Descriptor = model.OfferDescriptor{
			Name: protoOffer.GetDescriptorData().GetName(),
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
