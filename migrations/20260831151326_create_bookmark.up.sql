-- create "bookmarks" table
CREATE TABLE "public"."bookmarks" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "code" text NOT NULL,
  "description" text NULL,
  "url" text NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_bookmarks_deleted_at" to table: "bookmarks"
CREATE INDEX "idx_bookmarks_deleted_at" ON "public"."bookmarks" ("deleted_at");
