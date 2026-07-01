-- +goose Up
CREATE TABLE `marketdata/time_series`(
    figi Text,
    close_price Double,
    date Date,
    PRIMARY KEY(figi, date)
);

-- +goose Down
DROP TABLE `marketdata/time_series`;

