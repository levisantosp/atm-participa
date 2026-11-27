-- Modify "reports" table
ALTER TABLE "reports" ADD COLUMN "status" character varying NOT NULL DEFAULT 'open';
