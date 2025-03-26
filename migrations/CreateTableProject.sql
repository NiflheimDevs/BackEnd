CREATE TABLE IF NOT EXISTS "project" (
  "id" serial PRIMARY KEY,
  "owner_id" int NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "label" int DEFAULT 1,
  "selected_bid_id" int UNIQUE,
  "created_time" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_time" timestamp DEFAULT CURRENT_TIMESTAMP,
  "duration" timestamp,
  FOREIGN KEY ("owner_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("selected_bid_id") REFERENCES "bid" ("id"),
  FOREIGN KEY ("label") REFERENCES "label" ("id"),
);