CREATE TABLE IF NOT EXISTS "career" (
    "user_id" int NOT NULL,
    "company" varchar NOT NULL,
    "start_date" datetime NOT NULL,
    "end_date" datetime,
    "role" varchar NOT NULL.
    "website" varchar,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    CHECK ("end_date" IS NULL OR "start_date" < "end_date")
);