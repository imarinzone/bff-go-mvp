package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/search"
	"bff-go-mvp/internal/httpx"
	"bff-go-mvp/internal/model"
)

// SearchHandler handles POST /v1/search requests.
type SearchHandler struct {
	service search.Service
	logger  *zap.Logger
}

func NewSearchHandler(service search.Service, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		service: service,
		logger:  logger,
	}
}

// SearchChargingConnectors handles the search API.
// @Summary Search for EV charging connectors
// @Description Search charging locations/connectors by EVSE ID or geo-coordinates with optional filters.
// @Tags Search
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 20, max 100)"
// @Param request body model.SearchRequest true "Search request payload"
// @Success 200 {object} model.SearchResponse
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /v1/search [post]
func (h *SearchHandler) SearchChargingConnectors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Pagination params with defaults
	page := 1
	perPage := 20

	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 {
			page = p
		}
	}

	if v := r.URL.Query().Get("per_page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 && p <= 100 {
			perPage = p
		}
	}

	// Try to decode as DiscoverRequest first (full format with context and message)
	var discoverReq model.DiscoverRequest
	var req model.SearchRequest

	// Read body once and try both formats
	bodyBytes, err := readBody(r)
	if err != nil {
		h.logger.Warn("failed to read request body", zap.Error(err))
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	// Try DiscoverRequest format first
	var discoverContext *model.DiscoverContext
	if err := json.Unmarshal(bodyBytes, &discoverReq); err == nil && discoverReq.Context.Action != "" {
		// Convert DiscoverRequest to SearchRequest
		req = convertDiscoverToSearchRequest(discoverReq)
		// Extract context for passing to service
		discoverContext = &discoverReq.Context
	} else {
		// Fall back to SearchRequest format
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			h.logger.Warn("failed to decode search request", zap.Error(err))
			httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
			return
		}
	}

	// Validate oneOf: either evse_id or (geo_coordinates + distance_meters)
	hasEvse := req.EvseID != ""
	hasGeo := len(req.GeoCoordinates) == 2 && req.DistanceMeters > 0

	if (hasEvse && hasGeo) || (!hasEvse && !hasGeo) {
		httpx.WriteError(
			w,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"Either evse_id or geo_coordinates with distance_meters is required (but not both).",
		)
		return
	}

	// Add discover context to request context if available
	ctx := r.Context()
	if discoverContext != nil {
		ctx = context.WithValue(ctx, "discover_context", discoverContext)
	}

	resp, err := h.service.Search(ctx, page, perPage, req)
	if err != nil {
		h.logger.Error("search service failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	// Set mock transaction id header for now
	w.Header().Set("X-Transaction-Id", "mock-transaction-id")
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// readBody reads the request body and returns it as bytes
// It also restores the body so it can be read again if needed
func readBody(r *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return bodyBytes, nil
}

// convertDiscoverToSearchRequest converts a DiscoverRequest to SearchRequest
func convertDiscoverToSearchRequest(discoverReq model.DiscoverRequest) model.SearchRequest {
	req := model.SearchRequest{
		DistanceMeters: discoverReq.Message.DistanceMeters,
	}

	// Extract coordinates from geometry if present
	if discoverReq.Message.Geometry != nil && len(discoverReq.Message.Geometry.Coordinates) >= 2 {
		req.GeoCoordinates = discoverReq.Message.Geometry.Coordinates[:2]
	}

	return req
}
