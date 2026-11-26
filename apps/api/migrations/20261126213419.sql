-- Modify "issues" table
ALTER TABLE "issues" ADD COLUMN "cep" character varying NOT NULL, ADD COLUMN "address" character varying NOT NULL, ADD COLUMN "neighborhood" character varying NOT NULL;
