CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "created" timestamptz NOT NULL DEFAULT (now()),
  "uid" varchar NOT NULL,
  "type" varchar DEFAULT 'basic',
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "updated" timestamptz NOT NULL DEFAULT (now()),
  "deleted" timestamptz
);

CREATE TABLE "shifts" (
  "id" serial PRIMARY KEY,
  "created" timestamptz NOT NULL DEFAULT (now()),
  "uid" varchar NOT NULL,
  "work_date" timestamptz NOT NULL,
  "shift_length_hours" float NOT NULL,
  "user_id" int NOT NULL,
  "updated" timestamptz DEFAULT (now()),
  "deleted" timestamptz
);

-- Add foriegn keys
ALTER TABLE "shifts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- create function to automatically update "updated" column when UPDATE command is ran
CREATE  FUNCTION update_updated_column_on_users()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- trigger to update "updated" column on users
CREATE TRIGGER update_updated_column_on_users
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_column_on_users();
