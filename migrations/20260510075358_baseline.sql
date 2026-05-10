-- Create "announcements" table
CREATE TABLE "public"."announcements" (
  "id" uuid NOT NULL,
  "title" text NOT NULL,
  "url" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "available_from" timestamptz NOT NULL,
  "available_until" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_announcements_available_from" to table: "announcements"
CREATE INDEX "idx_announcements_available_from" ON "public"."announcements" ("available_from");
-- Create index "idx_announcements_available_until" to table: "announcements"
CREATE INDEX "idx_announcements_available_until" ON "public"."announcements" ("available_until");
-- Add new schema named "atlas_schema_revisions"
CREATE SCHEMA "atlas_schema_revisions";
-- Create "fcm_tokens" table
CREATE TABLE "public"."fcm_tokens" (
  "token" text NOT NULL,
  "user_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("token")
);
-- Create index "idx_fcm_tokens_updated_at" to table: "fcm_tokens"
CREATE INDEX "idx_fcm_tokens_updated_at" ON "public"."fcm_tokens" ("updated_at");
-- Create index "idx_fcm_tokens_user_id" to table: "fcm_tokens"
CREATE INDEX "idx_fcm_tokens_user_id" ON "public"."fcm_tokens" ("user_id");
-- Create "syllabuses" table
CREATE TABLE "public"."syllabuses" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "en_name" text NOT NULL,
  "grades" text NOT NULL,
  "credit" bigint NOT NULL,
  "faculty_names" text NOT NULL,
  "practical_home_faculty_category" text NOT NULL,
  "multiple_person_teaching_form" text NOT NULL,
  "teaching_form" text NOT NULL,
  "summary" text NOT NULL,
  "learning_outcomes" text NOT NULL,
  "assignments" text NOT NULL,
  "evaluation_method" text NOT NULL,
  "textbooks" text NOT NULL,
  "reference_books" text NOT NULL,
  "prerequisites" text NOT NULL,
  "pre_learning" text NOT NULL,
  "post_learning" text NOT NULL,
  "notes" text NOT NULL,
  "keywords" text NOT NULL,
  "target_courses" text NOT NULL,
  "target_areas" text NOT NULL,
  "classifications" text NOT NULL,
  "teaching_language" text NOT NULL,
  "contents_and_schedule" text NOT NULL,
  "teaching_and_exam_form" text NOT NULL,
  "dsop_subject" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "subjects" table
CREATE TABLE "public"."subjects" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "year" bigint NOT NULL,
  "semester" text NOT NULL,
  "credit" bigint NOT NULL,
  "classification" text NOT NULL,
  "cultural_subject_category" text NOT NULL,
  "syllabus_id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subjects_syllabus" FOREIGN KEY ("syllabus_id") REFERENCES "public"."syllabuses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_subjects_syllabus_id" to table: "subjects"
CREATE UNIQUE INDEX "idx_subjects_syllabus_id" ON "public"."subjects" ("syllabus_id");
-- Create "cancelled_classes" table
CREATE TABLE "public"."cancelled_classes" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "date" date NOT NULL,
  "period" text NOT NULL,
  "comment" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_cancelled_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_cancelled_classes_date" to table: "cancelled_classes"
CREATE INDEX "idx_cancelled_classes_date" ON "public"."cancelled_classes" ("date");
-- Create index "idx_cancelled_classes_subject_id" to table: "cancelled_classes"
CREATE INDEX "idx_cancelled_classes_subject_id" ON "public"."cancelled_classes" ("subject_id");
-- Create "course_registrations" table
CREATE TABLE "public"."course_registrations" (
  "id" uuid NOT NULL,
  "user_id" text NOT NULL,
  "subject_id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_course_registrations_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_course_registrations_subject_id" to table: "course_registrations"
CREATE INDEX "idx_course_registrations_subject_id" ON "public"."course_registrations" ("subject_id");
-- Create index "idx_course_registrations_user_id" to table: "course_registrations"
CREATE INDEX "idx_course_registrations_user_id" ON "public"."course_registrations" ("user_id");
-- Create "faculties" table
CREATE TABLE "public"."faculties" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "rooms" table
CREATE TABLE "public"."rooms" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "floor" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "faculty_rooms" table
CREATE TABLE "public"."faculty_rooms" (
  "id" uuid NOT NULL,
  "faculty_id" uuid NOT NULL,
  "room_id" uuid NOT NULL,
  "year" bigint NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_faculty_rooms_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_faculty_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_faculty_rooms_faculty_year" to table: "faculty_rooms"
CREATE UNIQUE INDEX "idx_faculty_rooms_faculty_year" ON "public"."faculty_rooms" ("faculty_id", "year");
-- Create index "idx_faculty_rooms_room_year" to table: "faculty_rooms"
CREATE UNIQUE INDEX "idx_faculty_rooms_room_year" ON "public"."faculty_rooms" ("room_id", "year");
-- Create "makeup_classes" table
CREATE TABLE "public"."makeup_classes" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "date" date NOT NULL,
  "period" text NOT NULL,
  "comment" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_makeup_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_makeup_classes_date" to table: "makeup_classes"
CREATE INDEX "idx_makeup_classes_date" ON "public"."makeup_classes" ("date");
-- Create index "idx_makeup_classes_subject_id" to table: "makeup_classes"
CREATE INDEX "idx_makeup_classes_subject_id" ON "public"."makeup_classes" ("subject_id");
-- Create "notifications" table
CREATE TABLE "public"."notifications" (
  "id" text NOT NULL,
  "title" text NOT NULL,
  "body" text NOT NULL,
  "image_url" text NULL,
  "analytics_label" text NULL,
  "ap_ns_badge" integer NULL,
  "ap_ns_sound" text NULL,
  "ap_ns_content_available" boolean NULL,
  "android_channel_id" text NULL,
  "android_priority" text NULL,
  "android_ttl_seconds" integer NULL,
  "webpush_link" text NULL,
  "url" text NULL,
  "notify_after" timestamptz NOT NULL,
  "notify_before" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_notifications_notify_after" to table: "notifications"
CREATE INDEX "idx_notifications_notify_after" ON "public"."notifications" ("notify_after");
-- Create index "idx_notifications_notify_before" to table: "notifications"
CREATE INDEX "idx_notifications_notify_before" ON "public"."notifications" ("notify_before");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" text NOT NULL,
  "email" text NOT NULL,
  "grade" text NULL,
  "course" text NULL,
  "class" text NULL,
  PRIMARY KEY ("id")
);
-- Create "notification_target_users" table
CREATE TABLE "public"."notification_target_users" (
  "notification_id" text NOT NULL,
  "user_id" text NOT NULL,
  "notified_at" timestamptz NULL,
  PRIMARY KEY ("notification_id", "user_id"),
  CONSTRAINT "fk_notification_target_users_notification" FOREIGN KEY ("notification_id") REFERENCES "public"."notifications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_notification_target_users_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_notification_target_users_notified_at" to table: "notification_target_users"
CREATE INDEX "idx_notification_target_users_notified_at" ON "public"."notification_target_users" ("notified_at");
-- Create "room_changes" table
CREATE TABLE "public"."room_changes" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "date" date NOT NULL,
  "period" text NOT NULL,
  "original_room_id" uuid NOT NULL,
  "new_room_id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_room_changes_new_room" FOREIGN KEY ("new_room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_room_changes_original_room" FOREIGN KEY ("original_room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_room_changes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_room_changes_date" to table: "room_changes"
CREATE INDEX "idx_room_changes_date" ON "public"."room_changes" ("date");
-- Create index "idx_room_changes_subject_id" to table: "room_changes"
CREATE INDEX "idx_room_changes_subject_id" ON "public"."room_changes" ("subject_id");
-- Create "subject_eligible_attributes" table
CREATE TABLE "public"."subject_eligible_attributes" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "grade" text NOT NULL,
  "class" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subjects_eligible_attributes" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_subject_eligible_attributes_subject_id" to table: "subject_eligible_attributes"
CREATE INDEX "idx_subject_eligible_attributes_subject_id" ON "public"."subject_eligible_attributes" ("subject_id");
-- Create "subject_faculties" table
CREATE TABLE "public"."subject_faculties" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "faculty_id" uuid NOT NULL,
  "is_primary" boolean NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subject_faculties_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_subjects_faculties" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_subject_faculties_subject_id" to table: "subject_faculties"
CREATE INDEX "idx_subject_faculties_subject_id" ON "public"."subject_faculties" ("subject_id");
-- Create "subject_requirements" table
CREATE TABLE "public"."subject_requirements" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "course" text NOT NULL,
  "requirement_type" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subjects_requirements" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_subject_requirements_subject_id" to table: "subject_requirements"
CREATE INDEX "idx_subject_requirements_subject_id" ON "public"."subject_requirements" ("subject_id");
-- Create "timetable_items" table
CREATE TABLE "public"."timetable_items" (
  "id" uuid NOT NULL,
  "subject_id" uuid NOT NULL,
  "day_of_week" text NULL,
  "period" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_timetable_items_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_timetable_items_subject_id" to table: "timetable_items"
CREATE INDEX "idx_timetable_items_subject_id" ON "public"."timetable_items" ("subject_id");
-- Create "timetable_item_rooms" table
CREATE TABLE "public"."timetable_item_rooms" (
  "id" uuid NOT NULL,
  "timetable_item_id" uuid NOT NULL,
  "room_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_timetable_item_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_timetable_items_rooms" FOREIGN KEY ("timetable_item_id") REFERENCES "public"."timetable_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_timetable_item_rooms_room_id" to table: "timetable_item_rooms"
CREATE INDEX "idx_timetable_item_rooms_room_id" ON "public"."timetable_item_rooms" ("room_id");
-- Create index "idx_timetable_item_rooms_timetable_item_id" to table: "timetable_item_rooms"
CREATE INDEX "idx_timetable_item_rooms_timetable_item_id" ON "public"."timetable_item_rooms" ("timetable_item_id");
