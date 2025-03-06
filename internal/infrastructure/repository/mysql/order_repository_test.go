go
package mysql_test

import (
	"context"
	"database/sql"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/infrastructure/repository/mysql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_SaveExchangeOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := mysql.NewOrderRepository(db)
	ctx := context.Background()
	now := time.Now()

	order := domain.TradeOrder{
		ID:           "1",
		UserFrom:     "user1",
		UserTo:       "user2",
		AmountFrom:   100,
		CurrencyFrom: "USD",
		CurrencyTo:   "EUR",
		Status:       "pending",
		CreatedAt:    now,
		ExpiresAt:    now.Add(1 * time.Hour),
	}

	t.Run("successful save", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO trade_orders").
			WithArgs(order.ID, order.UserFrom, order.UserTo, order.AmountFrom,
				order.CurrencyFrom, order.CurrencyTo, order.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveExchangeOrder(ctx, order)
		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO trade_orders").
			WillReturnError(sql.ErrConnDone)

		err := repo.SaveExchangeOrder(ctx, order)
		assert.Error(t, err)
	})
}

func TestOrderRepository_GetCurrencyByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := mysql.NewOrderRepository(db)
	ctx := context.Background()

	t.Run("currency found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"code", "symbol", "min_exchange"}).
			AddRow("USD", "$", 10.0)

		mock.ExpectQuery("SELECT code, symbol, min_exchange FROM currencies WHERE code = ?").
			WithArgs("USD").
			WillReturnRows(rows)

		currency, err := repo.GetCurrencyByCode(ctx, "USD")
		assert.NoError(t, err)
		assert.Equal(t, "USD", currency.Code)
	})

	t.Run("currency not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT code, symbol, min_exchange FROM currencies WHERE code = ?").
			WithArgs("XXX").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetCurrencyByCode(ctx, "XXX")
		assert.Error(t, err)
	})
}