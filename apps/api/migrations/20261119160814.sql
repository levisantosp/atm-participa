-- Create index "issue_status_id" to table: "issues"
CREATE INDEX "issue_status_id" ON "issues" ("status", "id" DESC);
-- Create index "issue_upvotes" to table: "issues"
CREATE INDEX "issue_upvotes" ON "issues" ("upvotes" DESC);
