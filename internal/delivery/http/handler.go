// HTTP обработчики с улучшенным логгированием и обработкой ошибок
package http

import (
	"context"
	"encoding/json"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"flexible_transfer_backend/internal/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type ExchangeHandler struct {
	uc     *usecase.ExchangeUseCase
	logger logger.Logger
}

func NewExchangeHandler(uc *usecase.ExchangeUseCase, l logger.Logger) *ExchangeHandler {
	return &ExchangeHandler{uc: uc, logger: l}
}

func respondWithError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrMinimumExchangeAmount):
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, domain.ErrUnsupportedCurrency):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *ExchangeHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var order domain.TradeOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.logger.Error("invalid request format", 
			zap.String("error", err.Error()),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateExchangeOrder(ctx, order); err != nil {
		h.logger.Error("order processing failed",
			zap.String("order_id", order.ID),
			zap.String("currency_pair", order.CurrencyFrom+"-"+order.CurrencyTo),
			zap.Error(err),
		)
		respondWithError(w, err)
		return
	}

	h.logger.Info("new order created",
		zap.String("order_id", order.ID),
		zap.Float64("amount", order.AmountFrom),
		zap.String("currency_pair", order.CurrencyFrom+"-"+order.CurrencyTo),
	)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}