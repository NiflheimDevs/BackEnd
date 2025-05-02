CREATE TABLE IF NOT EXISTS "transaction" (
  "id" serial PRIMARY KEY,
  "from_user_id" int NOT NULL,
  "to_user_id" int NOT NULL,
  "amount" numeric NOT NULL,
  "date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "description" text,
  FOREIGN KEY ("to_user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("from_user_id") REFERENCES "users" ("id")
);