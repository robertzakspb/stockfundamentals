-- +goose Up
ALTER TABLE `user/transaction`
ADD COLUMN timestamp Datetime;

-- +goose Down
ALTER TABLE `user/transaction`
DROP COLUMN timestamp;

