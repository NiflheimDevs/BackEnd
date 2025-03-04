CREATE TABLE IF NOT EXISTS "project" (
  "id" serial PRIMARY KEY,
  "owner_id" int NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "selected_bid_id" int UNIQUE,
  "created_time" timestamp NOT NULL,
  "updated_time" timestamp,
  "duration" timestamp,
  FOREIGN KEY ("owner_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("selected_bid_id") REFERENCES "bid" ("id")
);
