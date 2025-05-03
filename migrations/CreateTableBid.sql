CREATE TABLE IF NOT EXISTS "bid" (
  "id" serial PRIMARY KEY,
  "team_id" int NOT NULL,
  "project_id" int NOT NULL,
  "prepayment" numeric NOT NULL,
  "total" numeric NOT NULL,
  "description" text,
  "expected_time" int,
  "created_time" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("team_id") REFERENCES "team" ("id"),
  UNIQUE ("team_id", "project_id")
);
ALTER TABLE bid
ADD FOREIGN KEY ("project_id")
REFERENCES "project" ("id");