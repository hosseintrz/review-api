CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "suggestions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "text" varchar
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "suggestions" ("user_id");

ALTER TABLE "suggestions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
