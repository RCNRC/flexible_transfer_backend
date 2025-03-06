go
// Дополненный пакет тестов для сценариев использования
package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetCurrencyByCode(ctx context.Context, code string) (domain.Currency, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(domain.Currency), args.Error(1)
}

func (m *MockRepo) SaveExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockRepo) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	args := m.Called(ctx, order)
	return args.Get(0).([]domain.TradeOrder), args.Error(1)
}

type MockRateProvider struct {
	mock.Mock
}

func (m *MockRateProvider) GetCurrentRate(ctx context.Context, from, to string) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}

// Новые тесты для метода CreateExchangeOrder
func TestExchangeUseCase_CreateExchangeOrder(t *testing.T) {
	ctx := context.Background()
	validOrder := domain.TradeOrder{
		ID:           "test123",
		UserFrom:     "user1",
		CurrencyFrom: "USD",
		CurrencyTo:   "EUR",
		AmountFrom:   100,
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	tests := []struct {
		name        string
		order       domain.TradeOrder
		mockRepo    func(*MockRepo)
		expectedErr error
	}{
		{
			name:  "successful order creation",
			order: validOrder,
			mockRepo: func(mr *MockRepo) {
				mr.On("GetCurrencyByCode", ctx, "USD").Return(
					domain.Currency{Code: "USD", MinExchange: 10}, nil)
				mr.On("SaveExchangeOrder", ctx, validOrder).Return(nil)
			},
		},
		{
			name:  "invalid order validation",
			order: domain.TradeOrder{AmountFrom: -50},
			mockRepo: func(mr *MockRepo) {
				// No calls expected
			},
			expectedErr: errors.New("invalid amount"),
		},
		{
			name:  "currency not found",
			order: validOrder,
			mockRepo: func(mr *MockRepo) {
				mr.On("GetCurrencyByCode", ctx, "USD").Return(
					domain.Currency{}, domain.ErrUnsupportedCurrency)
			},
			expectedErr: domain.ErrUnsupportedCurrency,
		},
		{
			name:  "save order error",
			order: validOrder,
			mockRepo: func(mr *MockRepo) {
				mr.On("GetCurrencyByCode", ctx, "USD").Return(
					domain.Currency{Code: "USD", MinExchange: 10}, nil)
				mr.On("SaveExchangeOrder", ctx, validOrder).Return(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tt.mockRepo(mockRepo)

			uc := usecase.NewExchangeUseCase(mockRepo, nil)
			err := uc.CreateExchangeOrder(ctx, tt.order)

			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}