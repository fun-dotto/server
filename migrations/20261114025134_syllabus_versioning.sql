-- シラバスのバージョン管理に対応する。
-- 旧 syllabuses.id（教務システムの授業コード）を lesson_id に改名し、id を UUID で採番し直す。

-- 1. 旧 FK を外す
ALTER TABLE "public"."subjects" DROP CONSTRAINT "fk_subjects_syllabus";

-- 2. syllabuses: id → lesson_id に改名し、UUID の id を主キーにする
ALTER TABLE "public"."syllabuses" DROP CONSTRAINT "syllabuses_pkey";
ALTER TABLE "public"."syllabuses" RENAME COLUMN "id" TO "lesson_id";
ALTER TABLE "public"."syllabuses" ADD COLUMN "id" uuid NOT NULL DEFAULT gen_random_uuid(), ADD PRIMARY KEY ("id");

-- 3. year を追加し、紐づく subjects.year から埋める（紐づく科目がない場合は作成年）
ALTER TABLE "public"."syllabuses" ADD COLUMN "year" bigint;
UPDATE "public"."syllabuses" AS sy SET "year" = su."year" FROM "public"."subjects" AS su WHERE su."syllabus_id" = sy."lesson_id";
UPDATE "public"."syllabuses" SET "year" = EXTRACT(YEAR FROM "created_at") WHERE "year" IS NULL;
ALTER TABLE "public"."syllabuses" ALTER COLUMN "year" SET NOT NULL;

-- 4. 追加列（既存行は空文字）・credit の text 化・target_areas の削除
ALTER TABLE "public"."syllabuses"
  ADD COLUMN "display_order" text NOT NULL DEFAULT '',
  ADD COLUMN "list_lesson_name" text NOT NULL DEFAULT '',
  ADD COLUMN "list_en_name" text NOT NULL DEFAULT '',
  ADD COLUMN "list_staff_name" text NOT NULL DEFAULT '',
  ADD COLUMN "list_teaching_form" text NOT NULL DEFAULT '',
  ADD COLUMN "term" text NOT NULL DEFAULT '',
  ADD COLUMN "en_summary" text NOT NULL DEFAULT '',
  ADD COLUMN "en_learning_outcomes" text NOT NULL DEFAULT '',
  ADD COLUMN "en_contents_and_schedule" text NOT NULL DEFAULT '',
  ADD COLUMN "en_evaluation_method" text NOT NULL DEFAULT '',
  ADD COLUMN "en_teaching_language" text NOT NULL DEFAULT '',
  ADD COLUMN "en_prerequisites" text NOT NULL DEFAULT '',
  ADD COLUMN "en_pre_post_learning" text NOT NULL DEFAULT '',
  ADD COLUMN "en_textbooks" text NOT NULL DEFAULT '',
  ADD COLUMN "en_notes" text NOT NULL DEFAULT '',
  ADD COLUMN "en_teaching_and_exam_form" text NOT NULL DEFAULT '',
  ADD COLUMN "diploma_policies" text NOT NULL DEFAULT '',
  ADD COLUMN "diploma_policy_correspondence" text NOT NULL DEFAULT '',
  ADD COLUMN "content_hash" text NOT NULL DEFAULT '',
  ADD COLUMN "fetched_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ALTER COLUMN "credit" TYPE text USING ("credit"::text || '単位'),
  DROP COLUMN "target_areas";
ALTER TABLE "public"."syllabuses"
  ALTER COLUMN "display_order" DROP DEFAULT,
  ALTER COLUMN "list_lesson_name" DROP DEFAULT,
  ALTER COLUMN "list_en_name" DROP DEFAULT,
  ALTER COLUMN "list_staff_name" DROP DEFAULT,
  ALTER COLUMN "list_teaching_form" DROP DEFAULT,
  ALTER COLUMN "term" DROP DEFAULT,
  ALTER COLUMN "en_summary" DROP DEFAULT,
  ALTER COLUMN "en_learning_outcomes" DROP DEFAULT,
  ALTER COLUMN "en_contents_and_schedule" DROP DEFAULT,
  ALTER COLUMN "en_evaluation_method" DROP DEFAULT,
  ALTER COLUMN "en_teaching_language" DROP DEFAULT,
  ALTER COLUMN "en_prerequisites" DROP DEFAULT,
  ALTER COLUMN "en_pre_post_learning" DROP DEFAULT,
  ALTER COLUMN "en_textbooks" DROP DEFAULT,
  ALTER COLUMN "en_notes" DROP DEFAULT,
  ALTER COLUMN "en_teaching_and_exam_form" DROP DEFAULT,
  ALTER COLUMN "diploma_policies" DROP DEFAULT,
  ALTER COLUMN "diploma_policy_correspondence" DROP DEFAULT,
  ALTER COLUMN "content_hash" DROP DEFAULT;
CREATE INDEX "idx_syllabuses_lesson_year_fetched_at" ON "public"."syllabuses" ("lesson_id", "year", "fetched_at");
CREATE UNIQUE INDEX "idx_syllabuses_version" ON "public"."syllabuses" ("lesson_id", "year", "content_hash");

-- 5. subjects.syllabus_id を新しい syllabuses.id（uuid）に付け替える
ALTER TABLE "public"."subjects" ADD COLUMN "syllabus_uuid" uuid;
UPDATE "public"."subjects" AS su SET "syllabus_uuid" = sy."id" FROM "public"."syllabuses" AS sy WHERE sy."lesson_id" = su."syllabus_id";
DROP INDEX "public"."idx_subjects_syllabus_id";
ALTER TABLE "public"."subjects" DROP COLUMN "syllabus_id";
ALTER TABLE "public"."subjects" RENAME COLUMN "syllabus_uuid" TO "syllabus_id";
ALTER TABLE "public"."subjects" ALTER COLUMN "syllabus_id" SET NOT NULL;
CREATE UNIQUE INDEX "idx_subjects_syllabus_id" ON "public"."subjects" ("syllabus_id");
ALTER TABLE "public"."subjects" ADD CONSTRAINT "fk_subjects_syllabus" FOREIGN KEY ("syllabus_id") REFERENCES "public"."syllabuses" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;
