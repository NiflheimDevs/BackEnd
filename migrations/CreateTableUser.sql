CREATE TABLE IF NOT EXISTS "user" (
  "id" serial PRIMARY KEY,
  "firstname" varchar,
  "lastname" varchar,
  "username" varchar NOT NULL,
  "pasword" varchar NOT NULL,
  "email" varchar,
  "is_verified" bool DEFAULT FALSE,
  "bio" text,
  "phone" numeric NOT NULL,
  "photo" varchar,
  "resume" varchar,
  "wallet" numeric DEFAULT 0
);