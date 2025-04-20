CREATE TABLE IF NOT EXISTS "team" (
  "id" serial PRIMARY KEY,
  "type" int DEFAULT 0,
  "title" varchar,
  "description" text,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);