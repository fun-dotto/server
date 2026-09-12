-- Modify "calendar_dates" table
ALTER TABLE "public"."calendar_dates" DROP CONSTRAINT "fk_calendar_dates_calendar", ALTER COLUMN "date" TYPE date;
-- Modify "calendars" table
ALTER TABLE "public"."calendars" ALTER COLUMN "start_date" TYPE date, ALTER COLUMN "end_date" TYPE date;
-- Modify "fare_rules" table
ALTER TABLE "public"."fare_rules" ALTER COLUMN "route_id" DROP NOT NULL, ALTER COLUMN "origin_id" DROP NOT NULL, ALTER COLUMN "destination_id" DROP NOT NULL;
-- Modify "stop_times" table
ALTER TABLE "public"."stop_times" ALTER COLUMN "arrival_time" DROP NOT NULL, ALTER COLUMN "departure_time" DROP NOT NULL;
-- Modify "trips" table
ALTER TABLE "public"."trips" DROP CONSTRAINT "fk_trips_calendar", ALTER COLUMN "direction_id" DROP NOT NULL;
