go
// Дополненный тест для обработки ошибок
func TestRespondWithError(t *testing.T) {
    tests := []struct {
        name           string
        err            error
        expectedStatus int
    }{
        {
            name:           "unsupported currency error",
            err:            domain.ErrUnsupportedCurrency,
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:           "invalid exchange rate error",
            err:            domain.ErrInvalidExchangeRate,
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:           "generic error",
            err:            errors.New("unknown error"),
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/", nil)
            w := httptest.NewRecorder()
            
            // Тестируем через создание ошибочного ордера
            handler := http.NewExchangeHandler(new(MockUseCase))
            handler.CreateOrder(w, req) // Провоцируем ошибку с пустым телом
            
            // Принудительный вызов respondWithError
            http.respondWithError(w, tt.err)
            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}

func TestExchangeHandler_CreateOrder_UnprocessableEntity(t *testing.T) {
    mockUC := new(MockUseCase)
    mockUC.On("CreateExchangeOrder", mock.Anything, mock.Anything).
        Return(domain.ErrMinimumExchangeAmount)
    
    handler := http.NewExchangeHandler(mockUC)
    payload := `{"currency_from":"USD","amount_from":5}`
    
    req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(payload))
    w := httptest.NewRecorder()
    
    handler.CreateOrder(w, req)
    
    assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
    assert.Contains(t, w.Body.String(), "amount below minimum exchange")
}