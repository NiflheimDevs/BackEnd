CREATE TABLE IF NOT EXISTS "users" (
  "id" serial PRIMARY KEY,
  "firstname" varchar,
  "lastname" varchar,
  "username" varchar NOT NULL UNIQUE,
  "password" bytea NOT NULL,
  "email" varchar,
  "is_verified" bool DEFAULT FALSE,
  "bio" text,
  "phone" varchar NOT NULL,
  "wallet" numeric DEFAULT 0,
  "created_time" timestamp DEFAULT CURRENT_TIMESTAMP
);