ALTER TABLE "axes_resource_resource_rel" DROP CONSTRAINT "axes_resource_resource_rel_child_id_fkey" CASCADE;
ALTER TABLE "axes_resource_resource_rel" DROP CONSTRAINT "axes_resource_resource_rel_parent_id_fkey" CASCADE;
ALTER TABLE "axes_resource_user_rel" DROP CONSTRAINT "axes_resource_user_rel_resource_id_fkey" CASCADE;
ALTER TABLE "axes_resource_user_rel" DROP CONSTRAINT "axes_resource_user_rel_user_id_fkey" CASCADE;

DROP INDEX "axes_resource_type_idx" CASCADE;
DROP INDEX "axes_resource_resource_rel_child_id_idx" CASCADE;
DROP INDEX "axes_resource_resource_rel_child_id_type_idx" CASCADE;
DROP INDEX "axes_resource_resource_rel_parent_id_idx" CASCADE;
DROP INDEX "axes_resource_resource_rel_parent_id_type_idx" CASCADE;
DROP INDEX "axes_resource_user_rel_resource_id_idx" CASCADE;
DROP INDEX "axes_resource_user_rel_resource_id_user_id_idx" CASCADE;
DROP INDEX "axes_resource_user_rel_user_id_idx" CASCADE;
DROP INDEX "axes_user_display_name_idx" CASCADE;
DROP INDEX "axes_user_email_idx" CASCADE;
DROP INDEX "axes_user_mobile_idx" CASCADE;
DROP INDEX "axes_user_origin_id_idx" CASCADE;
DROP INDEX "axes_user_username_idx" CASCADE;

ALTER TABLE "axes_resource" DROP CONSTRAINT "axes_resource_pkey" CASCADE;
ALTER TABLE "axes_user" DROP CONSTRAINT "axes_user_pkey" CASCADE;

DROP TABLE "axes_resource" CASCADE;
DROP TABLE "axes_resource_resource_rel" CASCADE;
DROP TABLE "axes_resource_user_rel" CASCADE;
DROP TABLE "axes_user" CASCADE;

CREATE TABLE "axes_resource" (
    "id" serial8 NOT NULL,
    "type" int8 NOT NULL,
    CONSTRAINT "axes_resource_pkey" PRIMARY KEY ("id")
)WITHOUT OIDS;
CREATE INDEX "axes_resource_type_idx" ON "axes_resource" USING btree ("type" "pg_catalog"."int8_ops" ASC NULLS LAST);
COMMENT ON COLUMN "axes_resource"."type" IS '类型';

CREATE TABLE "axes_resource_resource_rel" (
    "id" serial8 NOT NULL,
    "parent_id" int8,
    "child_id" int8 NOT NULL,
    "type" varchar(1024) COLLATE "default"
)
WITHOUT OIDS;
CREATE INDEX "axes_resource_resource_rel_child_id_idx" ON "axes_resource_resource_rel" USING btree ("child_id" "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "axes_resource_resource_rel_child_id_type_idx" ON "axes_resource_resource_rel" USING btree ("child_id" "pg_catalog"."int8_ops" ASC NULLS LAST, "type" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "axes_resource_resource_rel_parent_id_idx" ON "axes_resource_resource_rel" USING btree ("parent_id" "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "axes_resource_resource_rel_parent_id_type_idx" ON "axes_resource_resource_rel" USING btree ("parent_id" "pg_catalog"."int8_ops" ASC NULLS LAST, "type" "pg_catalog"."text_ops" ASC NULLS LAST);

CREATE TABLE "axes_resource_user_rel" (
    "resource_id" int8 NOT NULL,
    "user_id" int8 NOT NULL
)
WITHOUT OIDS;
CREATE INDEX "axes_resource_user_rel_resource_id_idx" ON "axes_resource_user_rel" USING btree ("resource_id" "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "axes_resource_user_rel_resource_id_user_id_idx" ON "axes_resource_user_rel" USING btree ("resource_id" "pg_catalog"."int8_ops" ASC NULLS LAST, "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST);
CREATE INDEX "axes_resource_user_rel_user_id_idx" ON "axes_resource_user_rel" USING btree ("user_id" "pg_catalog"."int8_ops" ASC NULLS LAST);

CREATE TABLE "axes_user" (
    "id" serial8 NOT NULL,
    "username" varchar(100) COLLATE "default" NOT NULL,
    "display_name" varchar(255) COLLATE "default",
    "password" bytea,
    "mobile" varchar(255) COLLATE "default",
    "email" varchar(255) COLLATE "default",
    "user_source_id" int8 NOT NULL,
    "origin_id" text COLLATE "default",
    "created_at" timestamptz DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    CONSTRAINT "axes_user_pkey" PRIMARY KEY ("id")
)
WITHOUT OIDS;
CREATE INDEX "axes_user_display_name_idx" ON "axes_user" USING btree ("display_name" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "axes_user_email_idx" ON "axes_user" USING btree ("email" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "axes_user_mobile_idx" ON "axes_user" USING btree ("mobile" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "axes_user_origin_id_idx" ON "axes_user" USING btree ("origin_id" "pg_catalog"."text_ops" ASC NULLS LAST);
CREATE INDEX "axes_user_username_idx" ON "axes_user" USING btree ("username" "pg_catalog"."text_ops" ASC NULLS LAST);
COMMENT ON COLUMN "axes_user"."user_source_id" IS '用户源ID';
COMMENT ON COLUMN "axes_user"."origin_id" IS '用户在源内的ID';


ALTER TABLE "axes_resource_resource_rel" ADD CONSTRAINT "axes_resource_resource_rel_child_id_fkey" FOREIGN KEY ("child_id") REFERENCES "axes_resource" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "axes_resource_resource_rel" ADD CONSTRAINT "axes_resource_resource_rel_parent_id_fkey" FOREIGN KEY ("parent_id") REFERENCES "axes_resource" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "axes_resource_user_rel" ADD CONSTRAINT "axes_resource_user_rel_resource_id_fkey" FOREIGN KEY ("resource_id") REFERENCES "axes_resource" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "axes_resource_user_rel" ADD CONSTRAINT "axes_resource_user_rel_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "axes_user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

