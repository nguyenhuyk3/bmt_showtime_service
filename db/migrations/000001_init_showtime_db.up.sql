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
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "auditoriums" (
  "id" serial PRIMARY KEY,
  "cinema_id" int NOT NULL,
  "name" text NOT NULL,
  "seat_capacity" int NOT NULL DEFAULT 0,
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
  "is_deleted" boolean NOT NULL DEFAULT false,
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

CREATE TABLE "film_infos" (
  "id" serial PRIMARY KEY,
  "film_id" int UNIQUE NOT NULL,
  "duration" interval NOT NULL
);

CREATE INDEX ON "cinemas" ("id", "name");

CREATE INDEX ON "auditoriums" ("id", "cinema_id");

CREATE UNIQUE INDEX ON "seats" ("auditorium_id", "seat_number");

CREATE INDEX ON "showtimes" ("film_id");

CREATE INDEX ON "showtimes" ("auditorium_id");

CREATE UNIQUE INDEX ON "showtime_seats" ("showtime_id", "seat_id");

ALTER TABLE "auditoriums" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id");

ALTER TABLE "seats" ADD FOREIGN KEY ("auditorium_id") REFERENCES "auditoriums" ("id");

ALTER TABLE "showtimes" ADD FOREIGN KEY ("auditorium_id") REFERENCES "auditoriums" ("id");

ALTER TABLE "showtime_seats" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id");

ALTER TABLE "showtime_seats" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("id");


-- Addition commands
ALTER TABLE showtimes
ADD CONSTRAINT show_date_not_in_past
CHECK (show_date::date >= CURRENT_DATE);

ALTER TABLE showtimes
ADD CONSTRAINT valid_showtime_duration
CHECK (start_time <= end_time);

INSERT INTO
    "cinemas" ("name", "city", "location")
VALUES
    (
        'CGV Landmark',
        'HO_CHI_MINH',
        'Vincom Landmark 81'
    );

INSERT INTO film_infos (film_id, duration)
VALUES
  (1, INTERVAL '1 hour 31 minutes'),
  (2, INTERVAL '1 hour 45 minutes'),
  (3, INTERVAL '2 hours 10 minutes');

INSERT INTO
    "auditoriums" ("cinema_id", "name", "seat_capacity")
SELECT
    c.id,
    'Room ' || i,
    70
FROM
    generate_series (1, 5) AS i,
    cinemas c
WHERE
    c.name = 'CGV Landmark';

DO $$
DECLARE
    aud auditoriumS%ROWTYPE;
    row_labels TEXT[] := ARRAY['A','B','C','D','E','F','G','H'];
    seat_idx INT;
    row_label TEXT;
    seat_number TEXT;
    seat_type seat_types;
BEGIN
    FOR aud IN SELECT * FROM auditoriums WHERE cinema_id = (SELECT id FROM cinemas WHERE name = 'CGV Landmark') LOOP
        seat_idx := 0;
        FOR i IN 1..array_length(row_labels, 1) LOOP
            row_label := row_labels[i];
            FOR j IN 0..9 LOOP
                seat_number := row_label || j;
                
                IF seat_idx < 30 THEN
                    seat_type := 'standard';
                ELSIF seat_idx < 50 THEN
                    seat_type := 'vip';
                ELSE
                    seat_type := 'coupled';
                END IF;

                INSERT INTO seats (auditorium_id, seat_number, seat_type, price)
                VALUES (aud.id, seat_number, seat_type, 
                    CASE seat_type
                        WHEN 'standard' THEN 50000
                        WHEN 'vip' THEN 80000
                        WHEN 'coupled' THEN 100000
                    END
                );

                seat_idx := seat_idx + 1;
            END LOOP;
        END LOOP;
    END LOOP;
END $$;