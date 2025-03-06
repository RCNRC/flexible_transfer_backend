CREATE DATABASE IF NOT EXISTS flex_exchange;

CREATE TABLE IF NOT EXISTS currencies (
    code VARCHAR(3) PRIMARY KEY,
    symbol VARCHAR(3) NOT NULL,
    min_exchange DECIMAL(15,2) DEFAULT 1.00
);

CREATE TABLE IF NOT EXISTS trade_orders (
    id VARCHAR(36) PRIMARY KEY,
    user_from VARCHAR(255) NOT NULL,
    user_to VARCHAR(255) NOT NULL,
    amount_from DECIMAL(15,2) NOT NULL,
    currency_from VARCHAR(3) NOT NULL,
    currency_to VARCHAR(3) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP
);