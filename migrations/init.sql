
-- CREATE USER niflheim WITH PASSWORD 'niflguard';
-- ALTER ROLE niflheim WITH CREATEDB;

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'bidlancer') THEN
        CREATE DATABASE bidlancer;
    END IF;
END $$;

\c bidlancer;

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
  "wallet" numeric DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "team" (
  "id" serial PRIMARY KEY,
  "title" varchar,
  "description" text
);

CREATE TABLE IF NOT EXISTS "bid" (
  "id" serial PRIMARY KEY,
  "team_id" int NOT NULL,
  "project_id" int NOT NULL,
  "value" numeric NOT NULL,
  "expected_time" timestamp NOT NULL,
  "created_time" timestamp NOT NULL,
  FOREIGN KEY ("team_id") REFERENCES "team" ("id")
);

CREATE TABLE IF NOT EXISTS "project" (
  "id" serial PRIMARY KEY,
  "owner_id" int NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "label" text,
  "selected_bid_id" int UNIQUE,
  "created_time" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_time" timestamp DEFAULT CURRENT_TIMESTAMP,
  "duration" timestamp,
  FOREIGN KEY ("owner_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("selected_bid_id") REFERENCES "bid" ("id")
);

ALTER TABLE bid
ADD FOREIGN KEY ("project_id")
REFERENCES "project" ("id");

CREATE TABLE IF NOT EXISTS "chat" (
  "id" serial PRIMARY KEY,
  "title" varchar,
  "description" text
);

CREATE TABLE IF NOT EXISTS "comment" (
  "id" serial PRIMARY KEY,
  "project_id" int NOT NULL,
  "bid_id" int NOT NULL,
  "content" text,
  "rating" int NOT NULL,
  FOREIGN KEY ("project_id") REFERENCES "project" ("id"),
  FOREIGN KEY ("bid_id") REFERENCES "bid" ("id")
);

CREATE TABLE IF NOT EXISTS "transaction" (
  "id" serial PRIMARY KEY,
  "from_user_id" int NOT NULL,
  "to_user_id" int NOT NULL,
  "amount" numeric NOT NULL,
  "date" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "description" text,
  FOREIGN KEY ("to_user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("from_user_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "message" (
  "id" serial PRIMARY KEY,
  "chat_id" int NOT NULL,
  "sender_id" int,
  "sent_time" timestamp DEFAULT CURRENT_TIMESTAMP,
  "edit_time" timestamp,
  "content" text NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("sender_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "permission" (
  "id" serial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text
);

CREATE TABLE IF NOT EXISTS "role" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "career" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "company" varchar NOT NULL,
    "start_date" timestamp NOT NULL,
    "end_date" timestamp,
    "role" varchar NOT NULL,
    "website" varchar,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    CHECK ("end_date" IS NULL OR "start_date" < "end_date")
);

CREATE TABLE IF NOT EXISTS "role_permission" (
  "role_id" int NOT NULL,
  "permission_id" int NOT NULL,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  FOREIGN KEY ("permission_id") REFERENCES "permission" ("id")
);

CREATE TABLE IF NOT EXISTS "tag" (
  "id" serial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "users_chat" (
  "chat_id" int NOT NULL,
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "project_tag" (
  "project_id" int NOT NULL,
  "tag_id" int NOT NULL,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("project_id") REFERENCES "project" ("id"),
  UNIQUE ("tag_id", "project_id")
);

CREATE TABLE IF NOT EXISTS "users_career_tag" (
  "career_user_id" int NOT NULL,
  "tag_id" int NOT NULL,
  "type" int NOT NULL,
  "level" int,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("career_user_id") REFERENCES "career" ("id"),
  FOREIGN KEY ("career_user_id") REFERENCES "users" ("id"),
  UNIQUE ("career_user_id" , "tag_id" , "type")
);

CREATE TABLE IF NOT EXISTS "users_role" (
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("role_id") REFERENCES "role" ("id")
);

CREATE TABLE IF NOT EXISTS "users_team" (
  "user_id" int NOT NULL,
  "team_id" int NOT NULL,
  "position" varchar,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("team_id") REFERENCES "team" ("id")
);

