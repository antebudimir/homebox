-- +goose Up
-- Add a nullable default sport to entity templates so entities created via
-- the create-from-template flow inherit the template's sport. Nullable and
-- untouched by upstream code paths, so generic templates stay generic.
ALTER TABLE entity_templates ADD COLUMN default_sport TEXT;

-- +goose Down
ALTER TABLE entity_templates DROP COLUMN default_sport;
