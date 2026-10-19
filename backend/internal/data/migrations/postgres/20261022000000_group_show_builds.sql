-- +goose Up
-- Per-group visibility toggle for the Builds feature (sidebar entry).
ALTER TABLE "groups" ADD COLUMN "show_builds" boolean NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE "groups" DROP COLUMN "show_builds";