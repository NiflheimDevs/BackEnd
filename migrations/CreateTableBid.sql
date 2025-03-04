CREATE TABLE IF NOT EXISTS "bid" (
  "id" serial PRIMARY KEY,
  "team_id" int NOT NULL,
  "project_id" int NOT NULL,
  "value" numeric NOT NULL,
  "expected_time" timestamp NOT NULL,
  "created_time" timestamp NOT NULL,
  FOREIGN KEY ("team_id") REFERENCES "team" ("id")
);

ALTER TABLE bid
ADD FOREIGN KEY ("project_id")
REFERENCES "project" ("id");