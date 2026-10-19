-- +goose Up
-- Rename Setups to Builds to match the feature naming.
alter table setups rename to builds;
alter table setup_parts rename to build_parts;
alter table build_parts rename column setup_id to build_id;
alter table usage_sessions rename column setup_id to build_id;

-- +goose Down
alter table usage_sessions rename column build_id to setup_id;
alter table build_parts rename column build_id to setup_id;
alter table build_parts rename to setup_parts;
alter table builds rename to setups;
