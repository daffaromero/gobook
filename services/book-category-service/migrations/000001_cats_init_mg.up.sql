CREATE TABLE book_categories (
  "id" uuid PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT DEFAULT '-',
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);