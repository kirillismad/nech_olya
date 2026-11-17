BEGIN;

ALTER TABLE "notes"
    DROP COLUMN "archived_at";

COMMIT;
