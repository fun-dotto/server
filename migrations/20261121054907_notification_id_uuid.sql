-- 外部キーは text と uuid の列の間には張れないため、型を変える間だけ外す
ALTER TABLE "public"."notification_target_users" DROP CONSTRAINT "fk_notification_target_users_notification";
-- Modify "notifications" table
ALTER TABLE "public"."notifications" ALTER COLUMN "id" TYPE uuid USING "id"::uuid;
-- Modify "notification_target_users" table
ALTER TABLE "public"."notification_target_users" ALTER COLUMN "notification_id" TYPE uuid USING "notification_id"::uuid, ADD CONSTRAINT "fk_notification_target_users_notification" FOREIGN KEY ("notification_id") REFERENCES "public"."notifications" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
