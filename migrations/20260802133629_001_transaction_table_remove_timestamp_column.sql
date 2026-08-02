-- +goose Up
ALTER TABLE `user/transaction`
DROP COLUMN timestamp;

-- +goose Down
ALTER TABLE `user/transaction`
ADD COLUMN timestamp Date;
