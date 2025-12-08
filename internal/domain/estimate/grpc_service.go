package estimate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bff-go-mvp/internal/grpc"
	"bff-go-mvp/internal/model"
	selectpb "bff-go-mvp/proto/select/gen/select"
)

// GRPCService implements estimate.Service using gRPC client
type GRPCService struct {
	client *grpc.SelectClient
	logger *zap.Logger
}

// NewGRPCService creates a new gRPC-based estimate service
func NewGRPCService(client *grpc.SelectClient, logger *zap.Logger) *GRPCService {
	return &GRPCService{
		client: client,
		logger: logger,
	}
}

// Estimate implements the estimate.Service interface
func (s *GRPCService) Estimate(ctx context.Context, req model.EstimateRequest) (model.EstimateResponse, error) {
	// Convert model request to proto request
	protoReq := s.modelToProtoRequest(req)

	// Call gRPC service
	protoResp, err := s.client.Select(ctx, protoReq)
	if err != nil {
		s.logger.Error("gRPC estimate call failed", zap.Error(err))
		return model.EstimateResponse{}, err
	}

	// Convert proto response to model response
	return s.protoToModelResponse(protoResp), nil
}

// modelToProtoRequest converts model.EstimateRequest to proto SelectRequest
func (s *GRPCService) modelToProtoRequest(req model.EstimateRequest) *selectpb.SelectRequest {
	// Generate IDs
	transactionID := uuid.New().String()
	messageID := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Build context
	protoContext := &selectpb.Context{
		Version:       VersionDefault,
		Action:        ActionSelect,
		Domain:        DomainDefault,
		BapId:         BapIDDefault,
		BapUri:        BapURIDefault,
		TransactionId: transactionID,
		MessageId:     messageID,
		Timestamp:     timestamp,
		Ttl:           TTLDefault,
	}

	// Build message
	message := &selectpb.Message{
		EvseId:      req.EvseID,
		ConnectorId: req.ConnectorID,
		OfferId:     req.OfferID,
	}

	// Convert vehicle if provided
	if req.Vehicle.Make != "" || req.Vehicle.Model != "" || req.Vehicle.Type != "" {
		message.Vehicle = &selectpb.Vehicle{
			Make:  req.Vehicle.Make,
			Model: req.Vehicle.Model,
			Type:  req.Vehicle.Type,
		}
	}

	// Convert time window if provided
	if req.TimeWindow != nil {
		message.TimeWindow = &selectpb.TimeWindow{
			Start: req.TimeWindow.Start,
			End:   req.TimeWindow.End,
		}
	}

	// Convert energy if provided
	if req.Energy != nil {
		message.Energy = &selectpb.Energy{
			Value: req.Energy.Value,
			Unit:  req.Energy.Unit,
		}
	}

	// Convert amount if provided
	if req.Amount != nil {
		message.Amount = &selectpb.Amount{
			Value:    req.Amount.Value,
			Currency: req.Amount.Currency,
		}
	}

	return &selectpb.SelectRequest{
		Context: protoContext,
		Message: message,
	}
}

// protoToModelResponse converts proto OnSelectResponse to model EstimateResponse
func (s *GRPCService) protoToModelResponse(protoResp *selectpb.OnSelectResponse) model.EstimateResponse {
	if protoResp == nil || protoResp.Message == nil || protoResp.Message.Order == nil {
		return model.EstimateResponse{}
	}

	order := protoResp.Message.Order
	resp := model.EstimateResponse{}

	// Convert Order to OrderInfo
	resp.Order = model.OrderInfo{
		ID:     order.Id,
		Mode:   s.getFulfillmentMode(order),
		Status: order.OrderStatus,
	}

	// Convert OrderValue to Amount
	if order.OrderValue != nil {
		resp.Amount = model.Amount{
			Currency: order.OrderValue.Currency,
			Value:    order.OrderValue.Value,
		}

		// Convert price components
		if len(order.OrderValue.Components) > 0 {
			resp.PriceComponents = make([]model.PriceComponent, 0, len(order.OrderValue.Components))
			for _, comp := range order.OrderValue.Components {
				resp.PriceComponents = append(resp.PriceComponents, model.PriceComponent{
					Type:        comp.Type,
					Value:       comp.Value,
					Currency:    comp.Currency,
					Description: comp.Description,
				})
			}
		}
	}

	// Extract validity from accepted offer in order items
	if len(order.OrderItems) > 0 {
		for _, item := range order.OrderItems {
			if item.AcceptedOffer != nil && item.AcceptedOffer.Validity != nil {
				resp.Validity = &model.Validity{
					StartDate: item.AcceptedOffer.Validity.StartDate,
					EndDate:   item.AcceptedOffer.Validity.EndDate,
				}
				break
			}
		}
	}

	// Extract energy from order items (if available in accepted offer)
	// Note: Energy might need to be extracted from order items or fulfillment
	// For now, we'll leave it empty if not directly available

	// Extract cancellation policy from offer attributes
	if len(order.OrderItems) > 0 {
		for _, item := range order.OrderItems {
			if item.AcceptedOffer != nil && item.AcceptedOffer.OfferAttributes != nil {
				// Cancellation policy might be in offer attributes
				// For now, we'll leave it empty as it's not directly in the proto
				break
			}
		}
	}

	// Duration and percentage of battery charged are not directly in the proto
	// These might need to be calculated or extracted from other fields
	// For now, we'll leave them empty

	return resp
}

// getFulfillmentMode extracts the mode from fulfillment
func (s *GRPCService) getFulfillmentMode(order *selectpb.Order) string {
	if order.Fulfillment != nil {
		return order.Fulfillment.Mode
	}
	return ""
}

