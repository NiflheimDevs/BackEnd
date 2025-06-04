CREATE TABLE IF NOT EXISTS "notification"{
  "id" SERIAL PRIMARY KEY,
  "user_id" int,
  "content" TEXT,
  "send_time" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "read_time" TIMESTAMP WITH TIME ZONE,
  "isread" BOOLEAN DEFAULT FALSE,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
}