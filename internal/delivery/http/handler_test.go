go
// Расширенные тесты HTTP обработчиков
package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"flexible_transfer_backend/internal/delivery/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) CreateExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockUseCase) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	args := m.Called(ctx, order)
	return args.Get(0).([]domain.TradeOrder), args.Error(1)
}

// Новый тест для успешного создания ордера
func TestExchangeHandler_CreateOrder_Success(t *testing.T) {
	mockUC := new(MockUseCase)
	handler := http.NewExchangeHandler(mockUC)
	
	payload := `{
		"id": "test123",
		"user_from": "user1",
		"currency_from": "USD",
		"currency_to": "EUR",
		"amount_from": 100,
		"expires_at": "2024-01-01T00:00:00Z"
	}`

	mockUC.On("CreateExchangeOrder", mock.Anything, mock.Anything).Return(nil)
	
	req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.CreateOrder(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, payload, w.Body.String())
	mockUC.AssertExpectations(t)
}