go
package usecase_test

import (
	"context"
	"errors"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"testing"
	"time"

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
	return m.Called(ctx, order).Error(0)
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

func TestExchangeUseCase_CreateExchangeOrder(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		order       domain.TradeOrder
		repoSetup   func(*MockRepo)
		expectedErr error
	}{
		{
			name: "successful order creation",
			order: domain.TradeOrder{
				ID:           "1",
				CurrencyFrom: "USD",
				AmountFrom:   100,
				ExpiresAt:    now.Add(1 * time.Hour),
			},
			repoSetup: func(mr *MockRepo) {
				mr.On("GetCurrencyByCode", ctx, "USD").Return(domain.Currency{
					Code:        "USD",
					MinExchange: 10,
				}, nil)
				mr.On("SaveExchangeOrder", ctx, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "currency validation error",
			order: domain.TradeOrder{
				CurrencyFrom: "US",
				AmountFrom:   100,
			},
			repoSetup:   func(mr *MockRepo) {},
			expectedErr: domain.ErrInvalidCurrencyCode,
		},
		{
			name: "minimum amount error",
			order: domain.TradeOrder{
				CurrencyFrom: "USD",
				AmountFrom:   5,
				ExpiresAt:    now.Add(1 * time.Hour),
			},
			repoSetup: func(mr *MockRepo) {
				mr.On("GetCurrencyByCode", ctx, "USD").Return(domain.Currency{
					Code:        "USD",
					MinExchange: 10,
				}, nil)
			},
			expectedErr: domain.ErrMinimumExchangeAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockRates := new(MockRateProvider)
			tt.repoSetup(mockRepo)

			uc := usecase.NewExchangeUseCase(mockRepo, mockRates)
			err := uc.CreateExchangeOrder(ctx, tt.order)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}