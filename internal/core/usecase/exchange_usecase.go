go
// Бизнес-логика и сценарии использования системы
package usecase

import (
	"context"
	"flexible_transfer_backend/internal/core/domain"
)

// ExchangeRepository определяет контракт для работы с данными
type ExchangeRepository interface {
	GetCurrencyByCode(ctx context.Context, code string) (domain.Currency, error)
	SaveExchangeOrder(ctx context.Context, order domain.TradeOrder) error
	MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error)
}

// RateProvider определяет контракт для получения курсов
type RateProvider interface {
	GetCurrentRate(ctx context.Context, from, to string) (float64, error)
}

// ExchangeUseCase реализует бизнес-логику обмена
type ExchangeUseCase struct {
	repo   ExchangeRepository
	rates  RateProvider
}

func NewExchangeUseCase(r ExchangeRepository, rp RateProvider) *ExchangeUseCase {
	return &ExchangeUseCase{
		repo:  r,
		rates: rp,
	}
}

// CreateExchangeOrder обрабатывает запрос на создание ордера
func (uc *ExchangeUseCase) CreateExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	if err := order.Validate(); err != nil {
		return err
	}
	
	fromCurrency, err := uc.repo.GetCurrencyByCode(ctx, order.CurrencyFrom)
	if err != nil {
		return err
	}
	
	if order.AmountFrom < fromCurrency.MinExchange {
		return domain.ErrMinimumExchangeAmount
	}
	
	return uc.repo.SaveExchangeOrder(ctx, order)
}

// MatchOrders ищет подходящие совпадения для ордера
func (uc *ExchangeUseCase) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	rate, err := uc.rates.GetCurrentRate(ctx, order.CurrencyFrom, order.CurrencyTo)
	if err != nil {
		return nil, err
	}
	
	order.AmountFrom *= rate
	return uc.repo.MatchOrders(ctx, order)
}