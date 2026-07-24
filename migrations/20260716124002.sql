-- Modify "cancelled_classes" table
ALTER TABLE "public"."cancelled_classes" DROP CONSTRAINT "fk_cancelled_classes_subject", ADD CONSTRAINT "fk_cancelled_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "course_registrations" table
ALTER TABLE "public"."course_registrations" DROP CONSTRAINT "fk_course_registrations_subject", DROP CONSTRAINT "fk_course_registrations_user", ADD CONSTRAINT "fk_course_registrations_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_course_registrations_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "faculty_rooms" table
ALTER TABLE "public"."faculty_rooms" DROP CONSTRAINT "fk_faculty_rooms_faculty", DROP CONSTRAINT "fk_faculty_rooms_room", ADD CONSTRAINT "fk_faculty_rooms_faculty" FOREIGN KEY ("faculty_id") REFERENCES "public"."faculties" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_faculty_rooms_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "makeup_classes" table
ALTER TABLE "public"."makeup_classes" DROP CONSTRAINT "fk_makeup_classes_subject", ADD CONSTRAINT "fk_makeup_classes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "room_changes" table
ALTER TABLE "public"."room_changes" DROP CONSTRAINT "fk_room_changes_subject", ADD CONSTRAINT "fk_room_changes_subject" FOREIGN KEY ("subject_id") REFERENCES "public"."subjects" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
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
