CREATE TABLE IF NOT EXISTS "external_transaction" (
  "id" int PRIMARY KEY,
  "user_id" int NOT NULL,
  "type" int NOT NULL,
  "amount" numeric NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);