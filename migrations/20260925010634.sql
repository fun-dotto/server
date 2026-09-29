-- Modify "timetable_item_rooms" table
ALTER TABLE "public"."timetable_item_rooms" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
