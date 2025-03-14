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
  "photo" varchar,
  "resume" varchar,
  "wallet" numeric DEFAULT 0
);