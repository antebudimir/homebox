-- +goose Up
-- Replace the build is_active boolean with a status enum:
-- assembled (rideable), disassembled (torn down), retired (kept for history).
ALTER TABLE "builds" ADD COLUMN "status" character varying NOT NULL DEFAULT 'assembled';
UPDATE "builds" SET "status" = CASE WHEN "is_active" THEN 'assembled' ELSE 'disassembled' END;
ALTER TABLE "builds" DROP COLUMN "is_active";

-- +goose Down
ALTER TABLE "builds" ADD COLUMN "is_active" boolean NOT NULL DEFAULT false;
UPDATE "builds" SET "is_active" = CASE WHEN "status" = 'assembled' THEN true ELSE false END;
ALTER TABLE "builds" DROP COLUMN "status";