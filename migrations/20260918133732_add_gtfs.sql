-- Create "calendar_dates" table
CREATE TABLE "public"."calendar_dates" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "service_id" text NOT NULL,
  "date" date NOT NULL,
  "exception_type" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_calendar_dates_exception_type" CHECK (exception_type = ANY (ARRAY[(1)::bigint, (2)::bigint]))
);
-- Create index "idx_calendar_dates_date" to table: "calendar_dates"
CREATE INDEX "idx_calendar_dates_date" ON "public"."calendar_dates" ("date");
-- Create index "idx_calendar_dates_service_date" to table: "calendar_dates"
CREATE UNIQUE INDEX "idx_calendar_dates_service_date" ON "public"."calendar_dates" ("service_id", "date");
-- Create "calendars" table
CREATE TABLE "public"."calendars" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "service_id" text NOT NULL,
  "monday" bigint NOT NULL,
  "tuesday" bigint NOT NULL,
  "wednesday" bigint NOT NULL,
  "thursday" bigint NOT NULL,
  "friday" bigint NOT NULL,
  "saturday" bigint NOT NULL,
  "sunday" bigint NOT NULL,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_calendars_friday" CHECK (friday = ANY (ARRAY[(0)::bigint, (1)::bigint])),
  CONSTRAINT "chk_calendars_monday" CHECK (monday = ANY (ARRAY[(0)::bigint, (1)::bigint])),
  CONSTRAINT "chk_calendars_saturday" CHECK (saturday = ANY (ARRAY[(0)::bigint, (1)::bigint])),
  CONSTRAINT "chk_calendars_sunday" CHECK (sunday = ANY (ARRAY[(0)::bigint, (1)::bigint])),
  CONSTRAINT "chk_calendars_thursday" CHECK (thursday = ANY (ARRAY[(0)::bigint, (1)::bigint])),
  CONSTRAINT "chk_calendars_tuesday" CHECK (tuesday = ANY (ARRAY[(0)::bigint, (1)::bigint])),
  CONSTRAINT "chk_calendars_wednesday" CHECK (wednesday = ANY (ARRAY[(0)::bigint, (1)::bigint]))
);
-- Create index "idx_calendars_service_id" to table: "calendars"
CREATE UNIQUE INDEX "idx_calendars_service_id" ON "public"."calendars" ("service_id");
-- Create "stops" table
CREATE TABLE "public"."stops" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "stop_id" text NOT NULL,
  "stop_name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_stops_stop_id" to table: "stops"
CREATE UNIQUE INDEX "idx_stops_stop_id" ON "public"."stops" ("stop_id");
-- Create "routes" table
CREATE TABLE "public"."routes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "route_id" text NOT NULL,
  "route_short_name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_routes_route_id" to table: "routes"
CREATE UNIQUE INDEX "idx_routes_route_id" ON "public"."routes" ("route_id");
-- Create "fare_rules" table
CREATE TABLE "public"."fare_rules" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "route_id" text NULL,
  "origin_id" text NULL,
  "destination_id" text NULL,
  "price" numeric NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_fare_rules_destination" FOREIGN KEY ("destination_id") REFERENCES "public"."stops" ("stop_id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_fare_rules_origin" FOREIGN KEY ("origin_id") REFERENCES "public"."stops" ("stop_id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_fare_rules_route" FOREIGN KEY ("route_id") REFERENCES "public"."routes" ("route_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_fare_rules_destination_id" to table: "fare_rules"
CREATE INDEX "idx_fare_rules_destination_id" ON "public"."fare_rules" ("destination_id");
-- Create index "idx_fare_rules_origin_id" to table: "fare_rules"
CREATE INDEX "idx_fare_rules_origin_id" ON "public"."fare_rules" ("origin_id");
-- Create index "idx_fare_rules_route_id" to table: "fare_rules"
CREATE INDEX "idx_fare_rules_route_id" ON "public"."fare_rules" ("route_id");
-- Create "trips" table
CREATE TABLE "public"."trips" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "trip_id" text NOT NULL,
  "route_id" text NOT NULL,
  "service_id" text NOT NULL,
  "direction_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_trips_route" FOREIGN KEY ("route_id") REFERENCES "public"."routes" ("route_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_trips_route_id" to table: "trips"
CREATE INDEX "idx_trips_route_id" ON "public"."trips" ("route_id");
-- Create index "idx_trips_service_id" to table: "trips"
CREATE INDEX "idx_trips_service_id" ON "public"."trips" ("service_id");
-- Create index "idx_trips_trip_id" to table: "trips"
CREATE UNIQUE INDEX "idx_trips_trip_id" ON "public"."trips" ("trip_id");
-- Create "stop_times" table
CREATE TABLE "public"."stop_times" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "trip_id" text NOT NULL,
  "arrival_time" text NULL,
  "departure_time" text NULL,
  "stop_id" text NOT NULL,
  "stop_sequence" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_stop_times_stop" FOREIGN KEY ("stop_id") REFERENCES "public"."stops" ("stop_id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_stop_times_trip" FOREIGN KEY ("trip_id") REFERENCES "public"."trips" ("trip_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_stop_times_stop_id" to table: "stop_times"
CREATE INDEX "idx_stop_times_stop_id" ON "public"."stop_times" ("stop_id");
-- Create index "idx_stop_times_trip_id" to table: "stop_times"
CREATE INDEX "idx_stop_times_trip_id" ON "public"."stop_times" ("trip_id");
-- Create index "idx_stop_times_trip_seq" to table: "stop_times"
CREATE UNIQUE INDEX "idx_stop_times_trip_seq" ON "public"."stop_times" ("trip_id", "stop_sequence");
