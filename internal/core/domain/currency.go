go
// Доменные сущности системы и бизнес-правила
package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidCurrencyCode   = errors.New("invalid currency code")
	ErrInvalidExchangeRate   = errors.New("invalid exchange rate")
	ErrUnsupportedCurrency   = errors.New("unsupported currency")
	ErrMinimumExchangeAmount = errors.New("amount below minimum exchange")
)

// Currency представляет валюту по стандарту ISO 4217
type Currency struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	MinExchange float64 `json:"min_exchange"`
}

// Validate проверяет корректность валюты
func (c Currency) Validate() error {
	if len(c.Code) != 3 {
		return ErrInvalidCurrencyCode
	}
	if c.MinExchange < 0 {
		return ErrInvalidCurrencyCode
	}
	return nil
}

// ExchangeRate содержит актуальные курсы обмена
type ExchangeRate struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Rate      float64   `json:"rate"`
	Timestamp time.Time `json:"timestamp"`
}

// TradeOrder представляет запрос на обмен валюты
type TradeOrder struct {
	ID            string    `json:"id"`
	UserFrom      string    `json:"user_from"`
	UserTo        string    `json:"user_to"`
	AmountFrom    float64   `json:"amount_from"`
	CurrencyFrom  string    `json:"currency_from"`
	CurrencyTo    string    `json:"currency_to"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

func (o TradeOrder) Validate() error {
	if time.Now().After(o.ExpiresAt) {
		return errors.New("order expired")
	}
	if o.AmountFrom <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}