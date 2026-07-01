-- +goose Up
CREATE TABLE `user/account`(
    id Uuid,
    opening_date Date,
    type Text,
    broker Text,
    holder Text,
    primary_currency Text,
    cash_balance Double,
    PRIMARY KEY(id)
);

CREATE TABLE `user/account_market_value_history`(
    account_id Uuid,
    date Date,
    currency Text,
    eod_value Double,
    PRIMARY KEY(account_id, date, currency)
);

CREATE TABLE `user/transaction`(
    id Uuid,
    account_id Uuid,
    figi Text,
    type Text,
    timestamp Date,
    side Text,
    quantity Double,
    price_per_unit Double,
    currency Text,
    description Text,
    PRIMARY KEY(id)
);

CREATE TABLE `user/transaction_lot`(
    id Uuid,
    transaction_id Uuid,
    stock_lot_id Uuid,
    bond_lot_id Uuid,
    date Date,
    quantity Double,
    PRIMARY KEY(id)
);

CREATE TABLE `user/bond_position_lot`(
    id Uuid,
    figi Text,
    isin Text,
    opening_date Datetime,
    modification_date Datetime,
    account_id Uuid,
    quantity Double,
    price_per_unit_percentage Double,
    PRIMARY KEY(figi, account_id)
);

CREATE TABLE `user/stock_position_lot`(
    id Uuid,
    figi Text,
    account_id Uuid,
    created_at Datetime,
    updated_at Datetime,
    quantity Double,
    price_per_unit Double,
    currency Text,
    is_closed Bool,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE `user/account`;
DROP TABLE `user/account_market_value_history`;
DROP TABLE `user/transaction`;
DROP TABLE TABLE `user/transaction_lot`;
DROP TABLE TABLE `user/bond_position_lot`;
DROP TABLE TABLE `user/stock_position_lot`;