// Реализация бизнес-логики сценариев обмена валют
package usecase

import (
	"context"
	"errors"
	"flexible_transfer_backend/internal/core/domain"
	"fmt"
	"time"
)

type ExchangeRepository interface {
	GetCurrencyByCode(context.Context, string) (domain.Currency, error)
	SaveExchangeOrder(context.Context, domain.TradeOrder) error
	MatchOrders(context.Context, domain.TradeOrder) ([]domain.TradeOrder, error)
}

type RateProvider interface {
	GetCurrentRate(context.Context, string, string) (float64, error)
}

type ExchangeUseCase struct {
	repo  ExchangeRepository
	rates RateProvider
}

func NewExchangeUseCase(r ExchangeRepository, rp RateProvider) *ExchangeUseCase {
	return &ExchangeUseCase{repo: r, rates: rp}
}

// CreateExchangeOrder создает новый ордер с валидацией минимального лимита
func (uc *ExchangeUseCase) CreateExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	currency, err := uc.repo.GetCurrencyByCode(ctx, order.CurrencyFrom)
	if err != nil {
		return fmt.Errorf("currency validation error: %w", err)
	}

	if order.AmountFrom < currency.MinExchange {
		return domain.ErrMinimumExchangeAmount
	}

	order.CreatedAt = time.Now().UTC()
	if order.ExpiresAt.IsZero() {
		order.ExpiresAt = order.CreatedAt.Add(24 * time.Hour)
	}

	return uc.repo.SaveExchangeOrder(ctx, order)
}

// MatchOrders реализует логику поиска совпадающих ордеров
func (uc *ExchangeUseCase) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	rate, err := uc.rates.GetCurrentRate(ctx, order.CurrencyFrom, order.CurrencyTo)
	if err != nil {
		return nil, fmt.Errorf("rate provider error: %w", err)
	}
	
	if rate <= 0 {
		return nil, errors.New("invalid exchange rate: rate must be positive")
	}

	order.AmountFrom *= rate
	return uc.repo.MatchOrders(ctx, order)
}