go
// Дополненные тесты для MatchOrders
func TestOrderRepository_MatchOrders(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error creating mock database: %v", err)
    }
    defer db.Close()

    repo := mysql.NewOrderRepository(db)
    ctx := context.Background()
    
    testOrder := domain.TradeOrder{
        CurrencyFrom: "EUR",
        CurrencyTo:   "USD",
    }

    t.Run("successful match with orders", func(t *testing.T) {
        rows := sqlmock.NewRows([]string{"id", "user_from", "user_to", "amount_from", 
            "currency_from", "currency_to", "status", "created_at", "expires_at"}).
            AddRow("1", "user1", "user2", 50.0, "USD", "EUR", "pending", time.Now(), time.Now().Add(1*time.Hour)).
            AddRow("2", "user3", "user4", 30.0, "USD", "EUR", "pending", time.Now(), time.Now().Add(1*time.Hour))

        mock.ExpectQuery("SELECT \\* FROM trade_orders WHERE currency_from = \\? AND currency_to = \\? AND status = 'pending' LIMIT 50").
            WithArgs("USD", "EUR").
            WillReturnRows(rows)

        orders, err := repo.MatchOrders(ctx, testOrder)
        assert.NoError(t, err)
        assert.Len(t, orders, 2)
    })

    t.Run("database connection error", func(t *testing.T) {
        mock.ExpectQuery("SELECT \\* FROM trade_orders WHERE currency_from = \\? AND currency_to = \\? AND status = 'pending' LIMIT 50").
            WillReturnError(sql.ErrConnDone)

        _, err := repo.MatchOrders(ctx, testOrder)
        assert.ErrorIs(t, err, sql.ErrConnDone)
    })

    t.Run("no matching orders", func(t *testing.T) {
        mock.ExpectQuery("SELECT \\* FROM trade_orders WHERE currency_from = \\? AND currency_to = \\? AND status = 'pending' LIMIT 50").
            WillReturnRows(sqlmock.NewRows(nil))

        orders, err := repo.MatchOrders(ctx, testOrder)
        assert.NoError(t, err)
        assert.Empty(t, orders)
    })
}