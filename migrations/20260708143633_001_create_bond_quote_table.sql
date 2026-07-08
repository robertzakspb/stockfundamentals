-- +goose Up
CREATE TABLE `marketdata/bond_quote`(
    ticker Text,
    timestamp Datetime,
    price_as_percentage Double,
    PRIMARY KEY(ticker, timestamp)
);

-- +goose Down
DROP TABLE `marketdata/bond_quote`;
