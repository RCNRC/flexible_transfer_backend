go
package http_test

import (
	"bytes"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"flexible_transfer_backend/internal/delivery/http"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) CreateExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	return m.Called(ctx, order).Error(0)
}

func (m *MockUseCase) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	args := m.Called(ctx, order)
	return args.Get(0).([]domain.TradeOrder), args.Error(1)
}

func TestExchangeHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		mockSetup      func(*MockUseCase)
		expectedStatus int
	}{
		{
			name:    "successful request",
			payload: `{"id":"1","currency_from":"USD","amount_from":100,"currency_to":"EUR"}`,
			mockSetup: func(m *MockUseCase) {
				m.On("CreateExchangeOrder", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			payload:        `{invalid}`,
			mockSetup:      func(m *MockUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "validation error",
			payload: `{"currency_from":"US","amount_from":100}`,
			mockSetup: func(m *MockUseCase) {
				m.On("CreateExchangeOrder", mock.Anything, mock.Anything).Return(domain.ErrInvalidCurrencyCode)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockUseCase)
			tt.mockSetup(mockUC)

			handler := http.NewExchangeHandler(mockUC)
			req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(tt.payload))
			w := httptest.NewRecorder()

			handler.CreateOrder(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUC.AssertExpectations(t)
		})
	}
}