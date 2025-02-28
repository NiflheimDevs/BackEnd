CREATE TABLE "internal_transaction" (
  "id" serial PRIMARY KEY,
  "from_user_id" int,
  "to_user_id" int,
  "amount" numeric,
  FOREIGN KEY ("to_user_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("from_user_id") REFERENCES "user" ("id")
);