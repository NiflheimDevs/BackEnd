CREATE TABLE "user_role" (
  "user_id" int,
  "role_id" int,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("role_id") REFERENCES "role" ("id")
);