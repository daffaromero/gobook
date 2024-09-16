CREATE TABLE books (
  "id" uuid PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "author" VARCHAR(255) NOT NULL,
  "category_id" uuid,
  "description" TEXT DEFAULT '-',
  "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMPTZ DEFAULT NULL
);