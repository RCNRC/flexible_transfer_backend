go
// Расширенные тесты репозитория заказов
package mysql_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/infrastructure/repository/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Новые тесты для метода SaveExchangeOrder
func TestOrderRepository_SaveExchangeOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	repo := mysql.NewOrderRepository(db)
	ctx := context.Background()
	
	validOrder := domain.TradeOrder{
		ID:           "test123",
		UserFrom:     "user1",
		UserTo:       "user2",
		AmountFrom:   100,
		CurrencyFrom: "USD",
		CurrencyTo:   "EUR",
		Status:       "pending",
	}

	t.Run("successful save", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO trade_orders").
			WithArgs(
				validOrder.ID,
				validOrder.UserFrom,
				validOrder.UserTo,
				validOrder.AmountFrom,
				validOrder.CurrencyFrom,
				validOrder.CurrencyTo,
				validOrder.Status,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveExchangeOrder(ctx, validOrder)
		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO trade_orders").
			WillReturnError(errors.New("connection failed"))

		err := repo.SaveExchangeOrder(ctx, validOrder)
		assert.ErrorContains(t, err, "connection failed")
	})

	t.Run("duplicate entry", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO trade_orders").
			WillReturnError(&mysql.MySQLError{Number: 1062})

		err := repo.SaveExchangeOrder(ctx, validOrder)
		var mysqlErr *mysql.MySQLError
		assert.ErrorAs(t, err, &mysqlErr)
	})
}