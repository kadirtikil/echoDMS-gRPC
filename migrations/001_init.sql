CREATE TYPE IF NOT EXISTS "role" AS ENUM (
  'admin',
  'editor',
  'viewer'
);

CREATE TYPE IF NOT EXISTS "document_format" AS ENUM (
  'typst',
  'latex',
  'markdown'
);

CREATE TYPE IF NOT EXISTS "compiler_name" AS ENUM (
  'typst',
  'pdflatex',
  'xelatex',
  'lualatex',
  'pandoc'
);

CREATE TYPE IF NOT EXISTS "output_format" AS ENUM (
  'pdf',
  'png',
  'svg',
  'html'
);

CREATE TYPE IF NOT EXISTS "collaborator_status" AS ENUM (
  'invited',
  'accepted',
  'pending',
  'declined',
  'blocked'
);

CREATE TABLE IF NOT EXISTS "documents" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar(255) NOT NULL,
  "description" text,
  "format" document_format NOT NULL,
  "compiler" compiler_name NOT NULL,
  "compiler_version" varchar(20),
  "output_format" output_format DEFAULT 'pdf',
  "owner_id" varchar(255) NOT NULL,
  "is_archived" boolean DEFAULT false,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "pages" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "content" text,
  "sort_order" int NOT NULL DEFAULT 0,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "document_collaborators" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "user_id" varchar(255) NOT NULL,
  "role" role DEFAULT 'viewer',
  "status" collaborator_status,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "document_versions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "version_number" int NOT NULL,
  "snapshot" jsonb NOT NULL,
  "created_by" varchar(255) NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "compiled_outputs" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "version_id" uuid,
  "output_format" output_format NOT NULL,
  "file_path" varchar(500) NOT NULL,
  "file_size_bytes" bigint,
  "compiled_at" timestamptz DEFAULT (now())
);

CREATE UNIQUE INDEX IF NOT EXISTS ON "document_collaborators" ("document_id", "user_id");

CREATE UNIQUE INDEX IF NOT EXISTS ON "document_versions" ("document_id", "version_number");

ALTER TABLE "pages" ADD FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "document_collaborators" ADD FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "document_versions" ADD FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "compiled_outputs" ADD FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "compiled_outputs" ADD FOREIGN KEY ("version_id") REFERENCES "document_versions" ("id") DEFERRABLE INITIALLY IMMEDIATE;
