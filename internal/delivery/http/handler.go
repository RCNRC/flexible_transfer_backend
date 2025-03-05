go
// HTTP обработчики и роутинг
package http

import (
	"context"
	"encoding/json"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"net/http"
)

type ExchangeHandler struct {
	uc *usecase.ExchangeUseCase
}

func NewExchangeHandler(uc *usecase.ExchangeUseCase) *ExchangeHandler {
	return &ExchangeHandler{uc: uc}
}

func (h *ExchangeHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var order domain.TradeOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateExchangeOrder(ctx, order); err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func respondWithError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInvalidCurrencyCode:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case domain.ErrMinimumExchangeAmount:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}