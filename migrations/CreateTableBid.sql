CREATE TABLE IF NOT EXISTS "bid" (
  "id" serial PRIMARY KEY,
  "team_id" int NOT NULL,
  "project_id" int NOT NULL,
  "value" numeric NOT NULL,
  "expected_time" TIMESTAMP WITH TIME ZONE NOT NULL,
  "created_time" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("team_id") REFERENCES "team" ("id"),
  UNIQUE ("team_id", "project_id")
);
ALTER TABLE bid
ADD FOREIGN KEY ("project_id")
REFERENCES "project" ("id");