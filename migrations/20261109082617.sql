-- Create "room_reservations" table
CREATE TABLE "public"."room_reservations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "title" text NOT NULL,
  "room_id" uuid NOT NULL,
  "start_time" timestamptz NOT NULL,
  "end_time" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_room_reservations_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE NO ACTION
);
-- Create index "idx_room_reservations_end_time" to table: "room_reservations"
CREATE INDEX "idx_room_reservations_end_time" ON "public"."room_reservations" ("end_time");
-- Create index "idx_room_reservations_room_id" to table: "room_reservations"
CREATE INDEX "idx_room_reservations_room_id" ON "public"."room_reservations" ("room_id");
-- Create index "idx_room_reservations_start_time" to table: "room_reservations"
CREATE INDEX "idx_room_reservations_start_time" ON "public"."room_reservations" ("start_time");
