go
package usecase_test

import (
    // ... существующий импорт ...
    "flexible_transfer_backend/internal/core/usecase" // Добавлен явный импорт
)

// Добавленные тесты для метода MatchOrders
func TestExchangeUseCase_MatchOrders(t *testing.T) {
    ctx := context.Background()
    now := time.Now()

    tests := []struct {
        name          string
        order         domain.TradeOrder
        mockRepoSetup func(*MockRepo)
        mockRateSetup func(*MockRateProvider)
        expectedLen   int
        expectedErr   error
    }{
        {
            name: "successful match with rate conversion",
            order: domain.TradeOrder{
                CurrencyFrom: "USD",
                CurrencyTo:   "EUR",
                AmountFrom:   100,
                ExpiresAt:    now.Add(1 * time.Hour),
            },
            mockRateSetup: func(mr *MockRateProvider) {
                mr.On("GetCurrentRate", ctx, "USD", "EUR").Return(0.92, nil)
            },
            mockRepoSetup: func(mr *MockRepo) {
                mr.On("MatchOrders", ctx, mock.Anything).Return([]domain.TradeOrder{
                    {ID: "2", CurrencyFrom: "EUR", CurrencyTo: "USD", AmountFrom: 92},
                }, nil)
            },
            expectedLen: 1,
            expectedErr: nil,
        },
        {
            name: "rate provider error",
            order: domain.TradeOrder{
                CurrencyFrom: "USD",
                CurrencyTo:   "EUR",
            },
            mockRateSetup: func(mr *MockRateProvider) {
                mr.On("GetCurrentRate", ctx, "USD", "EUR").Return(0.0, domain.ErrInvalidExchangeRate)
            },
            mockRepoSetup: func(mr *MockRepo) {},
            expectedErr:   domain.ErrInvalidExchangeRate,
        },
        {
            name: "repository error",
            order: domain.TradeOrder{
                CurrencyFrom: "USD",
                CurrencyTo:   "EUR",
                AmountFrom:   100,
            },
            mockRateSetup: func(mr *MockRateProvider) {
                mr.On("GetCurrentRate", ctx, "USD", "EUR").Return(0.92, nil)
            },
            mockRepoSetup: func(mr *MockRepo) {
                mr.On("MatchOrders", ctx, mock.Anything).Return(nil, sql.ErrConnDone)
            },
            expectedErr: sql.ErrConnDone,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockRepo)
            mockRates := new(MockRateProvider)
            
            tt.mockRepoSetup(mockRepo)
            tt.mockRateSetup(mockRates)

            uc := usecase.NewExchangeUseCase(mockRepo, mockRates)
            result, err := uc.MatchOrders(ctx, tt.order)

            if tt.expectedErr != nil {
                assert.ErrorIs(t, err, tt.expectedErr)
            } else {
                assert.NoError(t, err)
                assert.Len(t, result, tt.expectedLen)
                // Проверка конвертации суммы
                if tt.name == "successful match with rate conversion" {
                    assert.Equal(t, 92.0, result[0].AmountFrom)
                }
            }
            mockRepo.AssertExpectations(t)
            mockRates.AssertExpectations(t)
        })
    }
}