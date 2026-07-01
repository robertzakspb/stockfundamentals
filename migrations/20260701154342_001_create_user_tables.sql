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

-- +goose Down
DROP TABLE `user/account`;
DROP TABLE `user/account_market_value_history`;