CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "username" varchar NOT NULL,
                         "password" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "movies" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar NOT NULL,
                          "year" integer NOT NULL,
                          "director" varchar,
                          "rating" numeric(2,1)
);

CREATE TABLE "reviews" (
                               "id" bigserial PRIMARY KEY,
                               "user_id" bigserial NOT NULL,
                               "movie_id" bigserial NOT NULL,
                               "rating" numeric(2,1),
                               "text" varchar
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "reviews" ("user_id", "movie_id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id");
