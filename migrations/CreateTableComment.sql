CREATE TABLE IF NOT EXISTS "comment" (
  "id" serial PRIMARY KEY,
  "project_id" int NOT NULL,
  "bid_id" int NOT NULL,
  "content" text,
  "rating" int NOT NULL,
  FOREIGN KEY ("project_id") REFERENCES "project" ("id"),
  FOREIGN KEY ("bid_id") REFERENCES "bid" ("id")
);