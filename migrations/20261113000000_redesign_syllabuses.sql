-- 科目テーブル・シラバステーブル再設計（docs/design/20261109_subject_syllabus_redesign.md）
-- syllabuses の主キーを教務システムの ID から UUID に変更し、旧 id を lesson_id に改名する。
-- subjects.syllabus_id は新しい UUID の syllabuses.id を参照する。

-- 1. syllabuses に UUID の id を追加して採番し、旧 id を lesson_id に改名する
ALTER TABLE "syllabuses" DROP CONSTRAINT "syllabuses_pkey" CASCADE;
ALTER TABLE "syllabuses" RENAME COLUMN "id" TO "lesson_id";
ALTER TABLE "syllabuses" ADD COLUMN "id" uuid NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE "syllabuses" ADD PRIMARY KEY ("id");

-- 2. year を追加し、紐づく subjects.year から埋める（紐づく科目がない行は作成日時の年度）
ALTER TABLE "syllabuses" ADD COLUMN "year" bigint NULL;
UPDATE "syllabuses" AS sy SET "year" = s."year" FROM "subjects" AS s WHERE s."syllabus_id" = sy."lesson_id";
UPDATE "syllabuses" SET "year" = EXTRACT(YEAR FROM COALESCE("created_at", CURRENT_TIMESTAMP) - INTERVAL '3 months')::bigint WHERE "year" IS NULL;
ALTER TABLE "syllabuses" ALTER COLUMN "year" SET NOT NULL;

-- 3. subjects に新しい syllabus_id を追加し、旧 syllabus_id（= lesson_id）で JOIN して設定する
ALTER TABLE "subjects" RENAME COLUMN "syllabus_id" TO "legacy_syllabus_id";
ALTER TABLE "subjects" ADD COLUMN "syllabus_id" uuid NULL;
UPDATE "subjects" AS s SET "syllabus_id" = sy."id" FROM "syllabuses" AS sy WHERE sy."lesson_id" = s."legacy_syllabus_id";
ALTER TABLE "subjects" ALTER COLUMN "syllabus_id" SET NOT NULL;

-- 4. FK・ユニークインデックスを張り替え、旧列を削除する
DROP INDEX IF EXISTS "idx_subjects_syllabus_id";
ALTER TABLE "subjects" DROP COLUMN "legacy_syllabus_id";
CREATE UNIQUE INDEX "idx_subjects_syllabus_id" ON "subjects" ("syllabus_id");
ALTER TABLE "subjects" ADD CONSTRAINT "fk_subjects_syllabus" FOREIGN KEY ("syllabus_id") REFERENCES "syllabuses" ("id") ON UPDATE CASCADE ON DELETE NO ACTION;

-- 5. credit を text に変換する（2 → 2単位）
ALTER TABLE "syllabuses" ALTER COLUMN "credit" TYPE text USING ("credit"::text || '単位');

-- 対象領域は「対象コース・領域」として target_courses にまとめる
UPDATE "syllabuses" SET "target_courses" = CASE WHEN "target_courses" = '' THEN "target_areas" ELSE "target_courses" || E'\n' || "target_areas" END WHERE "target_areas" <> '';
ALTER TABLE "syllabuses" DROP COLUMN "target_areas";

-- 取得元の列を追加する
ALTER TABLE "syllabuses"
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
  ADD COLUMN "content_hash" text NULL,
  ADD COLUMN "fetched_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- 既存行は取得元のハッシュを持たないため、移行前の行であることを示す値にする
-- （次回のシラバス取得ジョブで新しいバージョンとして追加される）
UPDATE "syllabuses" SET "content_hash" = 'legacy';
ALTER TABLE "syllabuses" ALTER COLUMN "content_hash" SET NOT NULL;

CREATE UNIQUE INDEX "idx_syllabuses_lesson_id_year_content_hash" ON "syllabuses" ("lesson_id", "year", "content_hash");
CREATE INDEX "idx_syllabuses_lesson_id_year_fetched_at" ON "syllabuses" ("lesson_id", "year", "fetched_at");
