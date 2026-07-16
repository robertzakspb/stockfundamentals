-- +goose Up
UPDATE `user/account` SET is_open = TRUE;

-- +goose Down

