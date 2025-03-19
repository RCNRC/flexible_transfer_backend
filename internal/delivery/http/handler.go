package http

import (
	"context"
	"encoding/json"
	"errors"
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

func (h *ExchangeHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.TradeOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateExchangeOrder(r.Context(), order); err != nil {
		h.logger.Error("order creation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "order_created"})
}

func (h *ExchangeHandler) MatchOrders(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	// Реализация парсинга параметров и вызова usecase
}
