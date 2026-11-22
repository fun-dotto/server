-- Modify "notification_target_users" table
ALTER TABLE "public"."notification_target_users" ALTER COLUMN "notification_id" TYPE uuid;
-- Modify "notifications" table
ALTER TABLE "public"."notifications" ALTER COLUMN "id" TYPE uuid;
