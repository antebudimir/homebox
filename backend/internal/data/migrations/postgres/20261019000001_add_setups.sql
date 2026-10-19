-- +goose Up
-- Setups: named assemblies of parts (e.g. skateboard / bike builds).
-- setup_parts carries added/removed dates so mounting history is preserved.
CREATE TABLE IF NOT EXISTS "setups" ("id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, "is_active" boolean NOT NULL DEFAULT false, "sport" character varying NULL, "group_setups" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "setups_groups_setups" FOREIGN KEY ("group_setups") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
CREATE INDEX IF NOT EXISTS "setups_group_setups" ON "setups" ("group_setups");
CREATE TABLE IF NOT EXISTS "setup_parts" ("id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "added_date" timestamptz NOT NULL, "removed_date" timestamptz NULL, "position" character varying NULL, "entity_id" uuid NOT NULL, "setup_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "setup_parts_entities_setup_parts" FOREIGN KEY ("entity_id") REFERENCES "entities" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "setup_parts_setups_parts" FOREIGN KEY ("setup_id") REFERENCES "setups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
CREATE INDEX IF NOT EXISTS "setup_parts_entity_id" ON "setup_parts" ("entity_id");
CREATE INDEX IF NOT EXISTS "setup_parts_setup_id" ON "setup_parts" ("setup_id");
CREATE TABLE IF NOT EXISTS "usage_sessions" ("id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "date" timestamptz NULL, "duration_minutes" bigint NOT NULL DEFAULT 0, "location" character varying NULL, "conditions" character varying NULL, "notes" character varying NULL, "rating" bigint NULL, "setup_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "usage_sessions_setups_usage_sessions" FOREIGN KEY ("setup_id") REFERENCES "setups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
CREATE INDEX IF NOT EXISTS "usage_sessions_setup_id" ON "usage_sessions" ("setup_id");

-- +goose Down
DROP TABLE IF EXISTS "usage_sessions";
DROP TABLE IF EXISTS "setup_parts";
DROP TABLE IF EXISTS "setups";
