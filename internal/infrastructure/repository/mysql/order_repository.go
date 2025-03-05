go
// Реализация репозиториев для работы с MySQL
package mysql

import (
	"context"
	"database/sql"
	"flexible_transfer_backend/internal/core/domain"
	"fmt"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetCurrencyByCode(ctx context.Context, code string) (domain.Currency, error) {
	query := `SELECT code, symbol, min_exchange FROM currencies WHERE code = ?`
	row := r.db.QueryRowContext(ctx, query, code)

	var c domain.Currency
	err := row.Scan(&c.Code, &c.Symbol, &c.MinExchange)
	if err != nil {
		return domain.Currency{}, fmt.Errorf("currency get error: %w", err)
	}
	return c, nil
}

func (r *OrderRepository) SaveExchangeOrder(ctx context.Context, order domain.TradeOrder) error {
	query := `INSERT INTO trade_orders(id, user_from, user_to, amount_from, currency_from, currency_to, status) VALUES(?,?,?,?,?,?,?)`
	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.UserFrom,
		order.UserTo,
		order.AmountFrom,
		order.CurrencyFrom,
		order.CurrencyTo,
		order.Status,
	)
	return err
}

func (r *OrderRepository) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	query := `SELECT * FROM trade_orders WHERE currency_from = ? AND currency_to = ? AND status = 'pending' LIMIT 50`
	rows, err := r.db.QueryContext(ctx, query, order.CurrencyTo, order.CurrencyFrom)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.TradeOrder
	for rows.Next() {
		var o domain.TradeOrder
		err := rows.Scan(&o.ID, &o.UserFrom, &o.UserTo, &o.AmountFrom,
			&o.CurrencyFrom, &o.CurrencyTo, &o.Status, &o.CreatedAt, &o.ExpiresAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}