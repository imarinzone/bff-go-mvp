package router

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"bff-go-mvp/internal/config"
	"bff-go-mvp/internal/domain/estimate"
	"bff-go-mvp/internal/domain/feedback"
	"bff-go-mvp/internal/domain/orders"
	"bff-go-mvp/internal/domain/payment"
	"bff-go-mvp/internal/domain/search"
	"bff-go-mvp/internal/domain/support"
	"bff-go-mvp/internal/grpc"
	"bff-go-mvp/internal/handler"
)

// New constructs the main HTTP router, wiring all handlers and middleware.
func New(cfg *config.Config, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(loggingMiddleware(logger))
	r.Use(recoveryMiddleware(logger))

	// Services
	searchService := chooseSearchService(cfg, logger)
	estimateService := chooseEstimateService(cfg, logger)
	paymentService := choosePaymentService(cfg, logger)
	ordersService := chooseOrdersService(cfg, logger)
	lifecycleService := chooseOrdersLifecycleService(cfg, logger)
	feedbackService := chooseFeedbackService(cfg, logger)
	supportService := chooseSupportService(cfg, logger)

	// Handlers
	searchHandler := handler.NewSearchHandler(searchService, logger)
	estimateHandler := handler.NewEstimateHandler(estimateService, logger)
	paymentHandler := handler.NewPaymentHandler(paymentService, logger)
	ordersHandler := handler.NewOrdersHandler(ordersService, logger)
	ordersLifecycleHandler := handler.NewOrdersLifecycleHandler(lifecycleService, logger)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService, logger)
	supportHandler := handler.NewSupportHandler(supportService, logger)

	// Routes from swagger.yaml
	r.HandleFunc("/v1/search", searchHandler.SearchChargingConnectors).Methods(http.MethodPost)
	r.HandleFunc("/v1/estimate", estimateHandler.GetEstimates).Methods(http.MethodPost)
	r.HandleFunc("/v1/orders/{order_id}/payment", paymentHandler.InitiatePayment).Methods(http.MethodPost)
	r.HandleFunc("/v1/orders/{order_id}", ordersHandler.GetOrder).Methods(http.MethodGet)
	r.HandleFunc("/v1/orders/{order_id}/cancel", ordersLifecycleHandler.EstimateCancel).Methods(http.MethodGet)
	r.HandleFunc("/v1/orders/{order_id}/cancel", ordersLifecycleHandler.Cancel).Methods(http.MethodPost)
	r.HandleFunc("/v1/orders/{order_id}/stop", ordersLifecycleHandler.EstimateStop).Methods(http.MethodGet)
	r.HandleFunc("/v1/orders/{order_id}/stop", ordersLifecycleHandler.StopCharging).Methods(http.MethodPut)
	r.HandleFunc("/v1/orders/{order_id}/start", ordersLifecycleHandler.StartCharging).Methods(http.MethodPut)
	r.HandleFunc("/v1/orders/{order_id}/rating", feedbackHandler.SetOrderRating).Methods(http.MethodPost)
	r.HandleFunc("/v1/orders/{order_id}/support", supportHandler.GetOrderSupport).Methods(http.MethodGet)

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Swagger UI (served by swaggo/http-swagger).
	// The docs are registered via the blank import of `internal/docs` in cmd/api/main.go.
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

func chooseSearchService(cfg *config.Config, logger *zap.Logger) search.Service {
	// Create gRPC client for discover service
	grpcClient, err := grpc.NewDiscoverClient(cfg.GRPC.DiscoverServiceAddress)
	if err != nil {
		logger.Error("Failed to create gRPC discover client, falling back to mock", zap.Error(err))
		return search.NewMockService()
	}

	// Create gRPC service
	return search.NewGRPCService(grpcClient, logger)
}

func chooseEstimateService(cfg *config.Config, logger *zap.Logger) estimate.Service {
	_ = logger
	_ = cfg
	return estimate.NewMockService()
}

func choosePaymentService(cfg *config.Config, logger *zap.Logger) payment.Service {
	_ = logger
	_ = cfg
	return payment.NewMockService()
}

func chooseOrdersService(cfg *config.Config, logger *zap.Logger) orders.Service {
	_ = logger
	_ = cfg
	return orders.NewMockService()
}

func chooseOrdersLifecycleService(cfg *config.Config, logger *zap.Logger) orders.LifecycleService {
	_ = logger
	_ = cfg
	return orders.NewMockLifecycleService()
}

func chooseFeedbackService(cfg *config.Config, logger *zap.Logger) feedback.Service {
	_ = logger
	_ = cfg
	return feedback.NewMockService()
}

func chooseSupportService(cfg *config.Config, logger *zap.Logger) support.Service {
	_ = logger
	_ = cfg
	return support.NewMockService()
}

// loggingMiddleware logs HTTP requests.
func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

// recoveryMiddleware recovers from panics.
func recoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
