-- Rename the account's user foreign key without losing existing associations.
ALTER TABLE "accounts" DROP CONSTRAINT "accounts_users_accounts";
ALTER TABLE "accounts" RENAME COLUMN "user_accounts" TO "user_id";
ALTER TABLE "accounts"
  ADD CONSTRAINT "accounts_users_accounts"
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
  ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER INDEX "account_user_accounts" RENAME TO "account_user_id";
