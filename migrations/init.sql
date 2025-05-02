
-- CREATE USER niflheim WITH PASSWORD 'niflguard';
-- ALTER ROLE niflheim WITH CREATEDB;


DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'bidlancer') THEN
        CREATE DATABASE bidlancer;
    END IF;
END $$;

\c bidlancer;

CREATE TABLE IF NOT EXISTS "permission" (
  "id" int PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "description" text
);

CREATE TABLE IF NOT EXISTS "role" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "label" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text,
  "price" numeric NOT NULL
);

CREATE TABLE IF NOT EXISTS "chat" (
  "id" serial PRIMARY KEY,
  "title" varchar,
  "description" text
);

CREATE TABLE IF NOT EXISTS "tag" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "team" (
  "id" serial PRIMARY KEY,
  "type" int DEFAULT 0,
  "title" varchar,
  "description" text,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "users" (
  "id" serial PRIMARY KEY,
  "firstname" varchar,
  "lastname" varchar,
  "username" varchar NOT NULL UNIQUE,
  "password" bytea NOT NULL,
  "email" varchar,
  "is_verified" bool DEFAULT FALSE,
  "bio" text,
  "phone" varchar NOT NULL,
  "wallet" numeric DEFAULT 0,
  "created_time" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "bid" (
  "id" serial PRIMARY KEY,
  "team_id" int NOT NULL,
  "project_id" int NOT NULL,
  "value" numeric NOT NULL,
  "expected_time" TIMESTAMP WITH TIME ZONE NOT NULL,
  "created_time" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("team_id") REFERENCES "team" ("id"),
  UNIQUE ("team_id", "project_id")
);

CREATE TABLE IF NOT EXISTS "project" (
  "id" serial PRIMARY KEY,
  "owner_id" int NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "label" int NOT NULL,
  "selected_bid_id" int UNIQUE,
  "status" int DEFAULT 0,
  "created_time" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_time" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "duration" TIMESTAMP WITH TIME ZONE,
  FOREIGN KEY ("owner_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("selected_bid_id") REFERENCES "bid" ("id"),
  FOREIGN KEY ("label") REFERENCES "label" ("id")
);

ALTER TABLE bid
ADD FOREIGN KEY ("project_id")
REFERENCES "project" ("id");

CREATE TABLE IF NOT EXISTS "comment" (
  "id" serial PRIMARY KEY,
  "project_id" int NOT NULL,
  "bid_id" int NOT NULL,
  "content" text,
  "rating" int NOT NULL,
  FOREIGN KEY ("project_id") REFERENCES "project" ("id"),
  FOREIGN KEY ("bid_id") REFERENCES "bid" ("id"),
  UNIQUE ("bid_id", "project_id")
);

CREATE TABLE IF NOT EXISTS "transaction" (
  "id" serial PRIMARY KEY,
  "from_user_id" int NOT NULL,
  "to_user_id" int NOT NULL,
  "amount" numeric NOT NULL,
  "date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "description" text,
  FOREIGN KEY ("to_user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("from_user_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "message" (
  "id" serial PRIMARY KEY,
  "chat_id" int NOT NULL,
  "sender_id" int,
  "sent_time" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "edit_time" TIMESTAMP WITH TIME ZONE,
  "content" text NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("sender_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "career" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "company" varchar NOT NULL,
    "start_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "end_date" TIMESTAMP WITH TIME ZONE,
    "role" varchar NOT NULL,
    "website" varchar,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    CHECK ("end_date" IS NULL OR "start_date" < "end_date")
);

CREATE TABLE IF NOT EXISTS "role_permission" (
  "role_id" int NOT NULL,
  "permission_id" int NOT NULL,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  FOREIGN KEY ("permission_id") REFERENCES "permission" ("id") ON DELETE CASCADE,
  UNIQUE ("permission_id", "role_id")
);

CREATE TABLE IF NOT EXISTS "users_chat" (
  "chat_id" int NOT NULL,
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  UNIQUE ("chat_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "project_tag" (
  "project_id" int NOT NULL,
  "tag_id" int NOT NULL,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
  UNIQUE ("tag_id", "project_id")
);

CREATE TABLE IF NOT EXISTS "users_career_tag" (
  "career_user_id" int NOT NULL,
  "tag_id" int NOT NULL,
  "type" int NOT NULL,
  "level" int,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  UNIQUE ("career_user_id" , "tag_id" , "type")
);

CREATE TABLE IF NOT EXISTS "users_team" (
  "user_id" int NOT NULL,
  "team_id" int NOT NULL,
  "position" varchar,
  "role_id" int NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("team_id") REFERENCES "team" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  UNIQUE ("team_id", "user_id")
);

CREATE FUNCTION delete_cascade_for_career() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_career_tag AS uct WHERE uct.career_user_id = OLD.id AND uct.type = 1;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER career_delete_trigger
BEFORE DELETE ON career
FOR EACH ROW EXECUTE FUNCTION delete_cascade_for_career();

CREATE FUNCTION delete_cascade_for_user() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_career_tag AS uct WHERE uct.career_user_id = OLD.id AND uct.type = 0;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_delete_trigger
BEFORE DELETE ON users
FOR EACH ROW EXECUTE FUNCTION delete_cascade_for_user();

CREATE OR REPLACE FUNCTION enforce_fk_constraint_on_users_career_tag() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.type = 1 THEN
        IF NOT EXISTS (SELECT 1 FROM career WHERE id = NEW.career_user_id) THEN
            RAISE EXCEPTION 'Invalid reference: % does not exist in career', NEW.career_user_id;
        END IF;
    ELSIF NEW.type = 0 THEN
        IF NOT EXISTS (SELECT 1 FROM users WHERE id = NEW.career_user_id) THEN
            RAISE EXCEPTION 'Invalid reference: % does not exist in users', NEW.career_user_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_fk_users_career_tag
BEFORE INSERT OR UPDATE ON users_career_tag
FOR EACH ROW EXECUTE FUNCTION enforce_fk_constraint_on_users_career_tag();


DO $$
BEGIN
  IF current_setting('app.skip_seeds', true) IS DISTINCT FROM 'on' THEN
    RAISE NOTICE 'skip_seeds is: %', current_setting('app.skip_seeds', true);
    INSERT INTO users (username,phone,password) VALUES 
    ('God','666','\x48656c6c6f20776f726c64'),
    ('deleted user','666','\x48656c6c6f20776f726c64');

    INSERT INTO tag (id,name) VALUES
    (0 , 'Data Structures'),
    (1 , 'Algorithms'),
    (2 , 'Operating Systems'),
    (3 , 'Computer Networks'),
    (4 , 'Database Management'),
    (5 , 'Cybersecurity'),
    (6 , 'Software Engineering'),
    (7 , 'Artificial Intelligence'),
    (8 , 'Machine Learning'),
    (9 , 'Deep Learning'),
    (10 , 'Cloud Computing'),
    (11 , 'Internet of Things (IoT)'),
    (12 , 'Embedded Systems'),
    (13 , 'Blockchain'),
    (14 , 'Quantum Computing'),
    (15 , 'DevOps'),
    (16 , 'System Design'),
    (17 , 'Parallel Computing'),
    (18 , 'High-Performance Computing'),
    (19 , 'Frontend Development'),
    (20 , 'Backend Development'),
    (21 , 'Full-Stack Development'),
    (57 , 'UI/UX Design'),
    (22 , 'Responsive Design'),
    (23 , 'Web Performance Optimization'),
    (24 , 'SEO'),
    (25 , 'Progressive Web Apps'),
    (26 , 'Serverless Architecture'),
    (27 , 'APIs & Microservices'),
    (28 , 'Authentication & Authorization'),
    (29 , 'Web Security'),
    (30 , 'WordPress'),
    (31 , 'Python'),
    (32 , 'JavaScript'),
    (33 , 'TypeScript'),
    (34 , 'Java'),
    (35 , 'C'),
    (36 , 'C++'),
    (37 , 'C#'),
    (38 , 'Go'),
    (39 , 'Rust'),
    (40 , 'PHP'),
    (41 , 'Swift'),
    (42 , 'Kotlin'),
    (43 , 'Dart'),
    (44 , 'Ruby'),
    (45 , 'Perl'),
    (46 , 'R'),
    (47 , 'MATLAB'),
    (48 , 'Scala'),
    (49 , 'Shell Scripting'),
    (50 , 'SQL'),
    (51 , 'NoSQL'),
    (52 , 'GraphQL'),
    (53 , 'HTML'),
    (54 , 'CSS'),
    (55 , 'Verilog'),
    (56 , 'VHDL');

    INSERT INTO label("id","name","description","price") VALUES
    (0 , 'Free','',0),
    (1 , 'Bold','important',75000),
    (2 , 'Urgent','',200000);


  END IF;
END $$;
-- ALTER SEQUENCE permission_id_seq RESTART WITH 0;
-- INSERT INTO "permission" ("name" , "description") VALUES
-- ("ADD_MEMBER", "adds member"),
-- ("REMOVE_MEMEBER", "removes a member"),
-- ("EDIT_INFO", "edit title, bio and ..."),
-- ("BIDDER", "the one who bids"),
-- ("EDIT_NICKNAME", "for teams, it works for positions. for groups and etc for nickname"),
-- ("EDIT_ROLE", "able to change the roles");
