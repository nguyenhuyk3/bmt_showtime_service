CREATE TYPE "seat_statuses" AS ENUM (
  'available',
  'reserved',
  'booked'
);

CREATE TYPE "cities" AS ENUM (
  'HO_CHI_MINH',
  'HA_NOI',
  'DONG_NAI'
);

CREATE TYPE "seat_types" AS ENUM (
  'standard',
  'coupled',
  'vip'
);

CREATE TABLE "cinemas" (
  "id" serial PRIMARY KEY NOT NULL,
  "name" text UNIQUE NOT NULL,
  "city" cities NOT NULL,
  "location" text NOT NULL,
  "is_released" bool NOT NULL DEFAULT true,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "auditoriums" (
  "id" serial PRIMARY KEY,
  "cinema_id" int NOT NULL,
  "name" text NOT NULL,
  "seat_capacity" int NOT NULL DEFAULT 0,
  "is_released" bool NOT NULL DEFAULT true,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "seats" (
  "id" serial PRIMARY KEY,
  "auditorium_id" int NOT NULL,
  "seat_number" varchar(5) NOT NULL,
  "seat_type" seat_types NOT NULL,
  "price" int NOT NULL DEFAULT 0,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "showtimes" (
  "id" serial PRIMARY KEY,
  "film_id" int NOT NULL,
  "auditorium_id" int NOT NULL,
  "show_date" date NOT NULL,
  "start_time" timestamp NOT NULL,
  "end_time" timestamp NOT NULL,
  "is_released" boolean NOT NULL DEFAULT false,
  "changed_by" varchar(32) NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "showtime_seats" (
  "id" serial PRIMARY KEY,
  "showtime_id" int NOT NULL,
  "seat_id" int NOT NULL,
  "status" seat_statuses NOT NULL,
  "booked_by" varchar(64) NOT NULL DEFAULT '',
  "created_at" timestamp DEFAULT (now()),
  "booked_at" timestamp
);

CREATE TABLE "outboxes" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
  "aggregated_type" varchar(64) NOT NULL,
  "aggregated_id" int NOT NULL,
  "event_type" varchar(64) NOT NULL,
  "payload" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "cinemas" ("id", "name");

CREATE INDEX ON "auditoriums" ("id", "cinema_id");

CREATE UNIQUE INDEX ON "seats" ("auditorium_id", "seat_number");

CREATE INDEX ON "showtimes" ("film_id");

CREATE INDEX ON "showtimes" ("auditorium_id");

CREATE UNIQUE INDEX ON "showtime_seats" ("showtime_id", "seat_id");

CREATE INDEX ON "outboxes" ("aggregated_type", "aggregated_id");

ALTER TABLE "auditoriums" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id");

ALTER TABLE "seats" ADD FOREIGN KEY ("auditorium_id") REFERENCES "auditoriums" ("id");

ALTER TABLE "showtimes" ADD FOREIGN KEY ("auditorium_id") REFERENCES "auditoriums" ("id");

ALTER TABLE "showtime_seats" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id");

ALTER TABLE "showtime_seats" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("id");
