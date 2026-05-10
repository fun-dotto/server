-- Modify "fcm_tokens" table
ALTER TABLE "public"."fcm_tokens" DROP CONSTRAINT "fk_fcm_tokens_user", ALTER COLUMN "created_at" SET NOT NULL, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" SET NOT NULL, ALTER COLUMN "updated_at" DROP DEFAULT;
-- Create index "idx_fcm_tokens_updated_at" to table: "fcm_tokens"
CREATE INDEX "idx_fcm_tokens_updated_at" ON "public"."fcm_tokens" ("updated_at");
-- Modify "announcements" table
ALTER TABLE "public"."announcements" ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT;
-- Modify "syllabuses" table
ALTER TABLE "public"."syllabuses" ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT;
-- Modify "subjects" table
ALTER TABLE "public"."subjects" DROP CONSTRAINT "fk_subjects_syllabus", ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT, ADD CONSTRAINT "fk_subjects_syllabus" FOREIGN KEY ("syllabus_id") REFERENCES "public"."syllabuses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "cancelled_classes" table
ALTER TABLE "public"."cancelled_classes" DROP CONSTRAINT "fk_cancelled_classes_subject", ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT, ADD CONSTRAINT "fk_cancelled_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "course_registrations" table
ALTER TABLE "public"."course_registrations" DROP CONSTRAINT "course_registrations_pkey", DROP CONSTRAINT "fk_course_registrations_subject", DROP CONSTRAINT "fk_course_registrations_user", ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT, ADD COLUMN "id" uuid NOT NULL, ADD PRIMARY KEY ("id"), ADD CONSTRAINT "fk_course_registrations_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "idx_course_registrations_subject_id" to table: "course_registrations"
CREATE INDEX "idx_course_registrations_subject_id" ON "public"."course_registrations" ("subject_id");
-- Create index "idx_course_registrations_user_id" to table: "course_registrations"
CREATE INDEX "idx_course_registrations_user_id" ON "public"."course_registrations" ("user_id");
-- Modify "faculties" table
ALTER TABLE "public"."faculties" ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT;
-- Modify "rooms" table
ALTER TABLE "public"."rooms" ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT;
-- Modify "faculty_rooms" table
ALTER TABLE "public"."faculty_rooms" DROP CONSTRAINT "faculty_rooms_pkey", DROP CONSTRAINT "fk_faculty_rooms_faculty", DROP CONSTRAINT "fk_faculty_rooms_room", ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT, ADD COLUMN "id" uuid NOT NULL, ADD PRIMARY KEY ("id"), ADD CONSTRAINT "fk_faculty_rooms_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_faculty_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "idx_faculty_rooms_faculty_year" to table: "faculty_rooms"
CREATE UNIQUE INDEX "idx_faculty_rooms_faculty_year" ON "public"."faculty_rooms" ("faculty_id", "year");
-- Create index "idx_faculty_rooms_room_year" to table: "faculty_rooms"
CREATE UNIQUE INDEX "idx_faculty_rooms_room_year" ON "public"."faculty_rooms" ("room_id", "year");
-- Modify "makeup_classes" table
ALTER TABLE "public"."makeup_classes" DROP CONSTRAINT "fk_makeup_classes_subject", ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT, ADD CONSTRAINT "fk_makeup_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "notifications" table
ALTER TABLE "public"."notifications" ALTER COLUMN "id" TYPE text, ALTER COLUMN "id" DROP DEFAULT, DROP COLUMN "created_at", DROP COLUMN "updated_at", ALTER COLUMN "ap_ns_badge" TYPE integer, ALTER COLUMN "android_ttl_seconds" TYPE integer;
-- Modify "notification_target_users" table
ALTER TABLE "public"."notification_target_users" DROP CONSTRAINT "fk_notification_target_users_notification", DROP CONSTRAINT "fk_notification_target_users_user", ALTER COLUMN "notification_id" TYPE text, ADD CONSTRAINT "fk_notification_target_users_notification" FOREIGN KEY ("notification_id") REFERENCES "public"."notifications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_notification_target_users_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "room_changes" table
ALTER TABLE "public"."room_changes" DROP CONSTRAINT "fk_room_changes_new_room", DROP CONSTRAINT "fk_room_changes_original_room", DROP CONSTRAINT "fk_room_changes_subject", ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT, ADD CONSTRAINT "fk_room_changes_new_room" FOREIGN KEY ("new_room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_room_changes_original_room" FOREIGN KEY ("original_room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_room_changes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "subject_eligible_attributes" table
ALTER TABLE "public"."subject_eligible_attributes" DROP CONSTRAINT "fk_subjects_eligible_attributes", ALTER COLUMN "id" DROP DEFAULT, ADD CONSTRAINT "fk_subjects_eligible_attributes" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "subject_faculties" table
ALTER TABLE "public"."subject_faculties" DROP CONSTRAINT "fk_subjects_faculties", ALTER COLUMN "id" DROP DEFAULT, ADD CONSTRAINT "fk_subjects_faculties" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "subject_requirements" table
ALTER TABLE "public"."subject_requirements" DROP CONSTRAINT "fk_subjects_requirements", ALTER COLUMN "id" DROP DEFAULT, ADD CONSTRAINT "fk_subjects_requirements" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "timetable_items" table
ALTER TABLE "public"."timetable_items" ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "created_at" DROP DEFAULT, ALTER COLUMN "updated_at" DROP DEFAULT;
-- Modify "timetable_item_rooms" table
ALTER TABLE "public"."timetable_item_rooms" DROP CONSTRAINT "timetable_item_rooms_pkey", DROP CONSTRAINT "fk_timetable_item_rooms_timetable_item", ALTER COLUMN "timetable_item_id" DROP DEFAULT, ALTER COLUMN "room_id" DROP DEFAULT, ADD COLUMN "id" uuid NOT NULL, ADD PRIMARY KEY ("id"), ADD CONSTRAINT "fk_timetable_items_rooms" FOREIGN KEY ("timetable_item_id") REFERENCES "public"."timetable_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "idx_timetable_item_rooms_room_id" to table: "timetable_item_rooms"
CREATE INDEX "idx_timetable_item_rooms_room_id" ON "public"."timetable_item_rooms" ("room_id");
-- Create index "idx_timetable_item_rooms_timetable_item_id" to table: "timetable_item_rooms"
CREATE INDEX "idx_timetable_item_rooms_timetable_item_id" ON "public"."timetable_item_rooms" ("timetable_item_id");
-- Drop "commons" table
DROP TABLE "public"."commons";
