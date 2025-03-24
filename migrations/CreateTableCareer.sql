CREATE TABLE IF NOT EXISTS "career" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "company" varchar NOT NULL,
    "start_date" timestamp NOT NULL,
    "end_date" timestamp,
    "role" varchar NOT NULL,
    "website" varchar,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    CHECK ("end_date" IS NULL OR "start_date" < "end_date")
);