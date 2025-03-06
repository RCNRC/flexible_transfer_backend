// Реализация репозитория MySQL с улучшенной логикой матчинга ордеров
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"flexible_transfer_backend/internal/core/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) MatchOrders(ctx context.Context, order domain.TradeOrder) ([]domain.TradeOrder, error) {
	query := `
		SELECT * FROM trade_orders 
		WHERE currency_from = ? 
		AND currency_to = ? 
		AND status = 'pending'
		AND amount_from >= ?
		AND expires_at > ?
		ORDER BY created_at ASC
		LIMIT 50`
	
	rows, err := r.db.QueryContext(ctx, query, 
		order.CurrencyTo, 
		order.CurrencyFrom,
		order.AmountFrom*0.95, // +/-5% допустимое отклонение
		time.Now().UTC(),
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.TradeOrder{}, nil
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var orders []domain.TradeOrder
	for rows.Next() {
		var o domain.TradeOrder
		err := rows.Scan(&o.ID, &o.UserFrom, &o.UserTo, &o.AmountFrom,
			&o.CurrencyFrom, &o.CurrencyTo, &o.Status, &o.CreatedAt, &o.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("row scanning error: %w", err)
		}
		orders = append(orders, o)
	}
	
	return orders, nil
}