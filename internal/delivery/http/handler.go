go
// HTTP обработчики и роутинг с интегрированным логгированием
package http

import (
	"context"
	"encoding/json"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"flexible_transfer_backend/internal/pkg/logger"
	"net/http"
)

type ExchangeHandler struct {
	uc     *usecase.ExchangeUseCase
	logger logger.Logger
}

func NewExchangeHandler(uc *usecase.ExchangeUseCase, l logger.Logger) *ExchangeHandler {
	return &ExchangeHandler{uc: uc, logger: l}
}

func (h *ExchangeHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var order domain.TradeOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.logger.Error("failed to decode request body", 
			zap.String("error", err.Error()),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateExchangeOrder(ctx, order); err != nil {
		h.logger.Error("order creation failed",
			zap.String("order_id", order.ID),
			zap.String("error", err.Error()),
		)
		respondWithError(w, err)
		return
	}

	h.logger.Info("order created successfully",
		zap.String("order_id", order.ID),
		zap.String("currency_pair", order.CurrencyFrom+"-"+order.CurrencyTo),
	)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}