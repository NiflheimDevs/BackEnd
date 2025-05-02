CREATE TABLE IF NOT EXISTS "permission" (
  "id" int PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "description" text
);