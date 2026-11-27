-- Modify "reports" table
ALTER TABLE "reports" DROP CONSTRAINT "report_target_matches_entity", ADD CONSTRAINT "report_target_matches_entity" CHECK ((((entity)::text = 'user'::text) AND (target_issue_id IS NULL)) OR (((entity)::text = 'issue'::text) AND (target_user_id IS NULL))), DROP COLUMN "entity_id";
