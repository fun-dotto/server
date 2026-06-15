-- Modify "timetable_items" table
ALTER TABLE "public"."timetable_items" ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP;
-- Modify "announcements" table
ALTER TABLE "public"."announcements" ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP;
-- Modify "syllabuses" table
ALTER TABLE "public"."syllabuses" ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP;
-- Modify "subjects" table
ALTER TABLE "public"."subjects" DROP CONSTRAINT "fk_subjects_syllabus", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD CONSTRAINT "fk_subjects_syllabus" FOREIGN KEY ("syllabus_id") REFERENCES "public"."syllabuses" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Modify "cancelled_classes" table
ALTER TABLE "public"."cancelled_classes" DROP CONSTRAINT "fk_cancelled_classes_subject", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD CONSTRAINT "fk_cancelled_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Drop index "idx_course_registrations_subject_id" from table: "course_registrations"
DROP INDEX "public"."idx_course_registrations_subject_id";
-- Drop index "idx_course_registrations_user_id" from table: "course_registrations"
DROP INDEX "public"."idx_course_registrations_user_id";
-- Modify "course_registrations" table
ALTER TABLE "public"."course_registrations" DROP CONSTRAINT "course_registrations_pkey", DROP CONSTRAINT "fk_course_registrations_subject", ALTER COLUMN "id" DROP NOT NULL, ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD PRIMARY KEY ("user_id", "subject_id"), ADD CONSTRAINT "fk_course_registrations_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION, ADD CONSTRAINT "fk_course_registrations_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Modify "faculties" table
ALTER TABLE "public"."faculties" ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP;
-- Modify "rooms" table
ALTER TABLE "public"."rooms" ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP;
-- Drop index "idx_faculty_rooms_faculty_year" from table: "faculty_rooms"
DROP INDEX "public"."idx_faculty_rooms_faculty_year";
-- Drop index "idx_faculty_rooms_room_year" from table: "faculty_rooms"
DROP INDEX "public"."idx_faculty_rooms_room_year";
-- Modify "faculty_rooms" table
ALTER TABLE "public"."faculty_rooms" DROP CONSTRAINT "faculty_rooms_pkey", DROP CONSTRAINT "fk_faculty_rooms_faculty", DROP CONSTRAINT "fk_faculty_rooms_room", ALTER COLUMN "id" DROP NOT NULL, ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD PRIMARY KEY ("faculty_id", "room_id", "year"), ADD CONSTRAINT "fk_faculty_rooms_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE CASCADE ON DELETE NO ACTION, ADD CONSTRAINT "fk_faculty_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Drop index "idx_fcm_tokens_updated_at" from table: "fcm_tokens"
DROP INDEX "public"."idx_fcm_tokens_updated_at";
-- Modify "fcm_tokens" table
ALTER TABLE "public"."fcm_tokens" ALTER COLUMN "created_at" DROP NOT NULL, ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" DROP NOT NULL, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD CONSTRAINT "fk_fcm_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "makeup_classes" table
ALTER TABLE "public"."makeup_classes" DROP CONSTRAINT "fk_makeup_classes_subject", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD CONSTRAINT "fk_makeup_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Modify "notifications" table
ALTER TABLE "public"."notifications" ALTER COLUMN "ap_ns_badge" TYPE bigint, ALTER COLUMN "android_ttl_seconds" TYPE bigint, ADD COLUMN "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP, ADD COLUMN "updated_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP;
-- Modify "notification_target_users" table
ALTER TABLE "public"."notification_target_users" DROP CONSTRAINT "fk_notification_target_users_notification", DROP CONSTRAINT "fk_notification_target_users_user", ADD CONSTRAINT "fk_notification_target_users_notification" FOREIGN KEY ("notification_id") REFERENCES "public"."notifications" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_notification_target_users_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "room_changes" table
ALTER TABLE "public"."room_changes" DROP CONSTRAINT "fk_room_changes_new_room", DROP CONSTRAINT "fk_room_changes_original_room", DROP CONSTRAINT "fk_room_changes_subject", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "created_at" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "updated_at" SET DEFAULT CURRENT_TIMESTAMP, ADD CONSTRAINT "fk_room_changes_new_room" FOREIGN KEY ("new_room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE NO ACTION, ADD CONSTRAINT "fk_room_changes_original_room" FOREIGN KEY ("original_room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE NO ACTION, ADD CONSTRAINT "fk_room_changes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Modify "subject_eligible_attributes" table
ALTER TABLE "public"."subject_eligible_attributes" DROP CONSTRAINT "fk_subjects_eligible_attributes", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ADD CONSTRAINT "fk_subjects_eligible_attributes" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Modify "subject_faculties" table
ALTER TABLE "public"."subject_faculties" DROP CONSTRAINT "fk_subjects_faculties", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ADD CONSTRAINT "fk_subjects_faculties" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
-- Modify "subject_requirements" table
ALTER TABLE "public"."subject_requirements" DROP CONSTRAINT "fk_subjects_requirements", ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ADD CONSTRAINT "fk_subjects_requirements" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
