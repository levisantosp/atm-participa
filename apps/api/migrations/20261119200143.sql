-- Drop index "issue_upvotes" from table: "issues"
DROP INDEX "issue_upvotes";
-- Create index "issue_upvotes_id" to table: "issues"
CREATE INDEX "issue_upvotes_id" ON "issues" ("upvotes" DESC, "id" DESC);
