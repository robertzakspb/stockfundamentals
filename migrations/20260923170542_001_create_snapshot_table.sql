-- +goose Up
CREATE TABLE `user/stock_lot_snapshot`(
    figi Text,
    isin Text,
    date Date,
    quantity Double,
    account_id Uuid,
    PRIMARY KEY(account_id, figi, date)
);

-- +goose Down
DROP TABLE `user/stock_lot_snapshot`
