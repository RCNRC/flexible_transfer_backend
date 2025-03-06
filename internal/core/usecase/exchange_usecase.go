go
// Бизнес-логика и сценарии использования системы с дополнительными проверками
package usecase

import (
	"context"
	"errors"
	"flexible_transfer_backend/internal/core/domain"
)

func (uc *ExchangeUseCase) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	rate, err := uc.rates.GetCurrentRate(ctx, order.CurrencyFrom, order.CurrencyTo)
	if err != nil {
		return nil, err
	}
	
	if rate <= 0 {
		return nil, errors.New("invalid exchange rate: rate must be positive")
	}

	order.AmountFrom *= rate
	return uc.repo.MatchOrders(ctx, order)
}