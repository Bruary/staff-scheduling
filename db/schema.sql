CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "created" timestamptz NOT NULL DEFAULT (now()),
  "uid" varchar NOT NULL,
  "type" varchar NOT NULL DEFAULT 'basic',
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

ALTER TABLE "shifts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
