CREATE TABLE IF NOT EXISTS "role_permission" (
  "role_id" int NOT NULL,
  "permission_id" int NOT NULL,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  FOREIGN KEY ("permission_id") REFERENCES "permission" ("id") ON DELETE CASCADE
);