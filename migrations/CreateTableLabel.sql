CREATE TABLE IF NOT EXISTS "label" (
    "id" serial PRIMARY KEY,
    "name" varchar NOT NULL,
    "description" text,
    "price" numeric NOT NULL
);