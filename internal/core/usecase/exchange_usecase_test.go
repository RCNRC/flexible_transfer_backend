package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"flexible_transfer_backend/internal/config"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCurrencyByCode(ctx context.Context, code string) (domain.Currency, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(domain.Currency), args.Error(1)
}

func (m *MockRepository) SaveExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockRepository) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	args := m.Called(ctx, order)
	return args.Get(0).([]domain.TradeOrder), args.Error(1)
}

// Добавляем моки для новых зависимостей

Коммит: c72f502e5719ed238721fa12ab42dc50b465c142
