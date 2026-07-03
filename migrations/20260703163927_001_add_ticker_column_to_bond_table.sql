-- +goose Up
ALTER TABLE `bonds/bond` ADD COLUMN ticker Text;

-- +goose Down
SELECT 'down SQL query';
ALTER TABLE `bonds/bond` DROP COLUMN ticker;
