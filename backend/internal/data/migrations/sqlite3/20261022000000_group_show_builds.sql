-- +goose Up
-- Per-group visibility toggle for the Builds feature (sidebar entry).
alter table groups add column show_builds boolean not null default true;

-- +goose Down
alter table groups drop column show_builds;