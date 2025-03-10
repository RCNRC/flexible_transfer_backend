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
		SELECT id, user_from, user_to, amount_from, 
		       currency_from, currency_to, status, 
		       created_at, expires_at 
		FROM trade_orders 
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
		order.AmountFrom*0.95,
		time.Now().UTC(),
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.TradeOrder{}, nil
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Printf("error closing rows: %v\n", closeErr)
		}
	}()

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
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	
	return orders, nil
}