CREATE TABLE "external_transaction" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "type" int,
  "amount" numeric,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);