CREATE TABLE IF NOT EXISTS "team" (
  "id" serial PRIMARY KEY,
  "type" int DEFAULT 0,
  "title" varchar,
  "description" text
);