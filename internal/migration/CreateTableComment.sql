CREATE TABLE "comment" (
  "id" serial PRIMARY KEY,
  "from_user_id" int,
  "to_user_id" int,
  "content" text,
  "rating" int,
  FOREIGN KEY ("from_user_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("to_user_id") REFERENCES "user" ("id")
);