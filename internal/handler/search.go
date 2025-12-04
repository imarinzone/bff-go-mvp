package handler

import (
	"encoding/json"
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

	var req model.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode search request", zap.Error(err))
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
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

	resp, err := h.service.Search(r.Context(), page, perPage, req)
	if err != nil {
		h.logger.Error("search service failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	// Set mock transaction id header for now
	w.Header().Set("X-Transaction-Id", "mock-transaction-id")
	httpx.WriteJSON(w, http.StatusOK, resp)
}
