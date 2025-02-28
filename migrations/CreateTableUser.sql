CREATE TABLE "user" (
  "id" serial PRIMARY KEY,
  "firstname" varchar,
  "lastname" varchar,
  "username" varchar,
  "pasword" varchar,
  "email" varchar,
  "is_verified" bool,
  "phone" numeric,
  "photo" varchar,
  "resume" varchar,
  "wallet" numeric
);