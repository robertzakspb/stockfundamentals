-- +goose Up
CREATE TABLE `forex/fx_rate`(
    currency_1 Text,
    currency_2 Text,
    date Date,
    rate Double,
    PRIMARY KEY(currency_1, currency_2, date)
);

-- +goose Down
DROP TABLE `forex/fx_rate`;
