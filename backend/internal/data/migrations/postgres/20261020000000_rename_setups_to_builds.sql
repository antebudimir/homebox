-- +goose Up
-- Rename Setups to Builds to match the feature naming.
ALTER TABLE "setups" RENAME TO "builds";
ALTER TABLE "setup_parts" RENAME TO "build_parts";
ALTER TABLE "setup_parts" RENAME COLUMN "setup_id" TO "build_id";
ALTER TABLE "usage_sessions" RENAME COLUMN "setup_id" TO "build_id";
ALTER TABLE "builds" RENAME COLUMN "group_setups" TO "group_builds";

ALTER TABLE "builds" RENAME CONSTRAINT "setups_groups_setups" TO "builds_groups_builds";
ALTER TABLE "build_parts" RENAME CONSTRAINT "setup_parts_entities_setup_parts" TO "build_parts_entities_build_parts";
ALTER TABLE "build_parts" RENAME CONSTRAINT "setup_parts_setups_parts" TO "build_parts_builds_parts";
ALTER TABLE "usage_sessions" RENAME CONSTRAINT "usage_sessions_setups_usage_sessions" TO "usage_sessions_builds_usage_sessions";

ALTER INDEX "setups_group_setups" RENAME TO "builds_group_builds";
ALTER INDEX "setup_parts_entity_id" RENAME TO "build_parts_entity_id";
ALTER INDEX "setup_parts_setup_id" RENAME TO "build_parts_build_id";
ALTER INDEX "usage_sessions_setup_id" RENAME TO "usage_sessions_build_id";

-- +goose Down
ALTER INDEX "usage_sessions_build_id" RENAME TO "usage_sessions_setup_id";
ALTER INDEX "build_parts_build_id" RENAME TO "setup_parts_setup_id";
ALTER INDEX "build_parts_entity_id" RENAME TO "setup_parts_entity_id";
ALTER INDEX "builds_group_builds" RENAME TO "setups_group_setups";

ALTER TABLE "usage_sessions" RENAME CONSTRAINT "usage_sessions_builds_usage_sessions" TO "usage_sessions_setups_usage_sessions";
ALTER TABLE "build_parts" RENAME CONSTRAINT "build_parts_builds_parts" TO "setup_parts_setups_parts";
ALTER TABLE "build_parts" RENAME CONSTRAINT "build_parts_entities_build_parts" TO "setup_parts_entities_setup_parts";
ALTER TABLE "builds" RENAME CONSTRAINT "builds_groups_builds" TO "setups_groups_setups";

ALTER TABLE "usage_sessions" RENAME COLUMN "build_id" TO "setup_id";
ALTER TABLE "build_parts" RENAME COLUMN "build_id" TO "setup_id";
ALTER TABLE "builds" RENAME COLUMN "group_builds" TO "group_setups";
ALTER TABLE "build_parts" RENAME TO "setup_parts";
ALTER TABLE "builds" RENAME TO "setups";
