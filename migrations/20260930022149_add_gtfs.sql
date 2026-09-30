-- Create "fare_prices" table
CREATE TABLE "public"."fare_prices" (
  "fare_id" text NOT NULL,
  "rider_category" text NOT NULL DEFAULT 'adult',
  "price" numeric NOT NULL,
  PRIMARY KEY ("fare_id", "rider_category")
);
-- Create "trips" table
CREATE TABLE "public"."trips" (
  "trip_id" text NOT NULL,
  "route_id" text NOT NULL,
  "service_id" text NOT NULL,
  "direction_id" bigint NULL,
  PRIMARY KEY ("trip_id"),
  CONSTRAINT "chk_trips_direction_id" CHECK (direction_id = ANY (ARRAY[(0)::bigint, (1)::bigint]))
);
-- Create index "idx_trips_route_id" to table: "trips"
CREATE INDEX "idx_trips_route_id" ON "public"."trips" ("route_id");
-- Create index "idx_trips_service_id" to table: "trips"
CREATE INDEX "idx_trips_service_id" ON "public"."trips" ("service_id");
-- Create "calendars" table
CREATE TABLE "public"."calendars" (
  "service_id" text NOT NULL,
  "monday" boolean NOT NULL,
  "tuesday" boolean NOT NULL,
  "wednesday" boolean NOT NULL,
  "thursday" boolean NOT NULL,
  "friday" boolean NOT NULL,
  "saturday" boolean NOT NULL,
  "sunday" boolean NOT NULL,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  PRIMARY KEY ("service_id"),
  CONSTRAINT "fk_trips_calendar" FOREIGN KEY ("service_id") REFERENCES "public"."trips" ("trip_id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "calendar_dates" table
CREATE TABLE "public"."calendar_dates" (
  "service_id" text NOT NULL,
  "date" date NOT NULL,
  "exception_type" bigint NOT NULL,
  PRIMARY KEY ("service_id", "date"),
  CONSTRAINT "fk_calendar_dates_calendar" FOREIGN KEY ("service_id") REFERENCES "public"."calendars" ("service_id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_calendar_dates_exception_type" CHECK (exception_type = ANY (ARRAY[(1)::bigint, (2)::bigint]))
);
-- Create "fare_rules" table
CREATE TABLE "public"."fare_rules" (
  "fare_id" text NOT NULL,
  "route_id" text NULL,
  "origin_id" text NULL,
  "destination_id" text NULL,
  PRIMARY KEY ("fare_id")
);
-- Create index "idx_fare_rules_destination_id" to table: "fare_rules"
CREATE INDEX "idx_fare_rules_destination_id" ON "public"."fare_rules" ("destination_id");
-- Create index "idx_fare_rules_origin_id" to table: "fare_rules"
CREATE INDEX "idx_fare_rules_origin_id" ON "public"."fare_rules" ("origin_id");
-- Create index "idx_fare_rules_route_id" to table: "fare_rules"
CREATE INDEX "idx_fare_rules_route_id" ON "public"."fare_rules" ("route_id");
-- Create "routes" table
CREATE TABLE "public"."routes" (
  "route_id" text NOT NULL,
  "route_short_name" text NOT NULL,
  PRIMARY KEY ("route_id"),
  CONSTRAINT "fk_fare_rules_route" FOREIGN KEY ("route_id") REFERENCES "public"."fare_rules" ("fare_id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_trips_route" FOREIGN KEY ("route_id") REFERENCES "public"."trips" ("trip_id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "stops" table
CREATE TABLE "public"."stops" (
  "stop_id" text NOT NULL,
  "stop_name" text NOT NULL,
  "stop_lat" numeric(9,6) NOT NULL,
  "stop_lon" numeric(9,6) NOT NULL,
  "zone_id" text NULL,
  PRIMARY KEY ("stop_id")
);
-- Create index "idx_stops_zone_id" to table: "stops"
CREATE INDEX "idx_stops_zone_id" ON "public"."stops" ("zone_id");
-- Create "stop_times" table
CREATE TABLE "public"."stop_times" (
  "trip_id" text NOT NULL,
  "arrival_time" text NOT NULL,
  "departure_time" text NOT NULL,
  "stop_id" text NOT NULL,
  "stop_sequence" bigint NOT NULL,
  PRIMARY KEY ("trip_id", "stop_sequence"),
  CONSTRAINT "fk_stop_times_stop" FOREIGN KEY ("stop_id") REFERENCES "public"."stops" ("stop_id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_stop_times_trip" FOREIGN KEY ("trip_id") REFERENCES "public"."trips" ("trip_id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_stop_times_stop_id" to table: "stop_times"
CREATE INDEX "idx_stop_times_stop_id" ON "public"."stop_times" ("stop_id");
