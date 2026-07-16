-- +goose Up
ALTER TABLE `user/account` 
ADD COLUMN is_open Bool;

-- +goose Down
ALTER TABLE `user/account` 
DROP COLUMN is_open;
