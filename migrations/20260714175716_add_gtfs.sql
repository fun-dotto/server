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
  "start_date" text NOT NULL,
  "end_date" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_calendars_service_id" to table: "calendars"
CREATE UNIQUE INDEX "idx_calendars_service_id" ON "public"."calendars" ("service_id");
-- Create "calendar_dates" table
CREATE TABLE "public"."calendar_dates" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "service_id" text NOT NULL,
  "date" text NOT NULL,
  "exception_type" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_calendar_dates_calendar" FOREIGN KEY ("service_id") REFERENCES "public"."calendars" ("service_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_calendar_dates_date" to table: "calendar_dates"
CREATE INDEX "idx_calendar_dates_date" ON "public"."calendar_dates" ("date");
-- Create index "idx_calendar_dates_service_id" to table: "calendar_dates"
CREATE INDEX "idx_calendar_dates_service_id" ON "public"."calendar_dates" ("service_id");
-- Modify "cancelled_classes" table
ALTER TABLE "public"."cancelled_classes" DROP CONSTRAINT "fk_cancelled_classes_subject", ADD CONSTRAINT "fk_cancelled_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "course_registrations" table
ALTER TABLE "public"."course_registrations" DROP CONSTRAINT "fk_course_registrations_subject", DROP CONSTRAINT "fk_course_registrations_user", ADD CONSTRAINT "fk_course_registrations_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_course_registrations_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "faculty_rooms" table
ALTER TABLE "public"."faculty_rooms" DROP CONSTRAINT "fk_faculty_rooms_faculty", DROP CONSTRAINT "fk_faculty_rooms_room", ADD CONSTRAINT "fk_faculty_rooms_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_faculty_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
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
  "route_id" text NOT NULL,
  "origin_id" text NOT NULL,
  "destination_id" text NOT NULL,
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
-- Modify "makeup_classes" table
ALTER TABLE "public"."makeup_classes" DROP CONSTRAINT "fk_makeup_classes_subject", ADD CONSTRAINT "fk_makeup_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "room_changes" table
ALTER TABLE "public"."room_changes" DROP CONSTRAINT "fk_room_changes_subject", ADD CONSTRAINT "fk_room_changes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Create "trips" table
CREATE TABLE "public"."trips" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "trip_id" text NOT NULL,
  "route_id" text NOT NULL,
  "service_id" text NOT NULL,
  "direction_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_trips_calendar" FOREIGN KEY ("service_id") REFERENCES "public"."calendars" ("service_id") ON UPDATE CASCADE ON DELETE CASCADE,
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
  "arrival_time" text NOT NULL,
  "departure_time" text NOT NULL,
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
-- Modify "subject_eligible_attributes" table
ALTER TABLE "public"."subject_eligible_attributes" DROP CONSTRAINT "fk_subjects_eligible_attributes", ADD CONSTRAINT "fk_subjects_eligible_attributes" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "subject_faculties" table
ALTER TABLE "public"."subject_faculties" DROP CONSTRAINT "fk_subject_faculties_faculty", DROP CONSTRAINT "fk_subjects_faculties", ADD CONSTRAINT "fk_subject_faculties_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_subjects_faculties" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "subject_requirements" table
ALTER TABLE "public"."subject_requirements" DROP CONSTRAINT "fk_subjects_requirements", ADD CONSTRAINT "fk_subjects_requirements" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "timetable_items" table
ALTER TABLE "public"."timetable_items" DROP CONSTRAINT "fk_timetable_items_subject", ADD CONSTRAINT "fk_timetable_items_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "timetable_item_rooms" table
ALTER TABLE "public"."timetable_item_rooms" DROP CONSTRAINT "fk_timetable_item_rooms_room", DROP CONSTRAINT "fk_timetable_items_rooms", ADD CONSTRAINT "fk_timetable_item_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_timetable_items_rooms" FOREIGN KEY ("timetable_item_id") REFERENCES "public"."timetable_items" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
