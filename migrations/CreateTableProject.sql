CREATE TABLE "project" (
  "id" serial PRIMARY KEY,
  "owner_id" int,
  "title" text,
  "description" text,
  "selected_bid_id" int UNIQUE,
  "created_time" timestamp,
  "updated_time" timestamp,
  "duration" timestamp,
  FOREIGN KEY ("owner_id") REFERENCES "user" ("id")
);