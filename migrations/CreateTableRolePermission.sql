CREATE TABLE "role_permission" (
  "role_id" int,
  "permission_id" int,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  FOREIGN KEY ("permission_id") REFERENCES "permission" ("id")
);