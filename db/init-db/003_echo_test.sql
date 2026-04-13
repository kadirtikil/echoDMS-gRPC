CREATE DATABASE echo_test;

\c echo_test

CREATE TYPE "role" AS ENUM (
  'admin',
  'editor',
  'viewer'
);

CREATE TYPE "document_format" AS ENUM (
  'typst',
  'latex',
  'markdown'
);

CREATE TYPE "compiler_name" AS ENUM (
  'typst',
  'pdflatex',
  'xelatex',
  'lualatex',
  'pandoc'
);

CREATE TYPE "output_format" AS ENUM (
  'pdf',
  'png',
  'svg',
  'html'
);

CREATE TYPE "collaborator_status" AS ENUM (
  'invited',
  'accepted',
  'pending',
  'declined',
  'blocked'
);

CREATE TABLE "documents" (
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

CREATE TABLE "pages" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "content" text,
  "sort_order" int NOT NULL DEFAULT 0,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "document_collaborators" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "user_id" varchar(255) NOT NULL,
  "role" role DEFAULT 'viewer',
  "status" collaborator_status,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "document_versions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "version_number" int NOT NULL,
  "snapshot" jsonb NOT NULL,
  "created_by" varchar(255) NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "compiled_outputs" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "document_id" uuid NOT NULL,
  "version_id" uuid,
  "output_format" output_format NOT NULL,
  "file_path" varchar(500) NOT NULL,
  "file_size_bytes" bigint,
  "compiled_at" timestamptz DEFAULT (now())
);

CREATE UNIQUE INDEX idx_doc_collab_unique ON "document_collaborators" ("document_id", "user_id");

CREATE UNIQUE INDEX idx_doc_version_unique ON "document_versions" ("document_id", "version_number");

ALTER TABLE "pages" ADD CONSTRAINT fk_pages_document FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "document_collaborators" ADD CONSTRAINT fk_collab_document FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "document_versions" ADD CONSTRAINT fk_version_document FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "compiled_outputs" ADD CONSTRAINT fk_output_document FOREIGN KEY ("document_id") REFERENCES "documents" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "compiled_outputs" ADD CONSTRAINT fk_output_version FOREIGN KEY ("version_id") REFERENCES "document_versions" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- Documents
INSERT INTO "documents" ("id", "title", "description", "format", "compiler", "compiler_version", "output_format", "owner_id")
VALUES
  ('a1b2c3d4-1111-1111-1111-000000000001', 'Thesis Draft', 'My master thesis on distributed systems', 'latex', 'pdflatex', '3.14', 'pdf', 'a1b2c3d4-1111-1111-1111-000000000002'),
  ('a1b2c3d4-1111-1111-1111-000000000002', 'Meeting Notes', 'Weekly team sync notes', 'markdown', 'pandoc', '2.19', 'html', 'a1b2c3d4-1111-1111-1111-000000000002'),
  ('a1b2c3d4-1111-1111-1111-000000000003', 'Invoice Template', 'Standard invoice layout', 'typst', 'typst', '0.11', 'pdf', 'a1b2c3d4-1111-1111-1111-000000000002');

-- Pages
INSERT INTO "pages" ("id", "document_id", "content", "sort_order")
VALUES
  ('b1b2c3d4-2222-2222-2222-000000000001', 'a1b2c3d4-1111-1111-1111-000000000001', '\\section{Introduction}\nThis thesis explores...', 0),
  ('b1b2c3d4-2222-2222-2222-000000000002', 'a1b2c3d4-1111-1111-1111-000000000001', '\\section{Background}\nDistributed systems are...', 1),
  ('b1b2c3d4-2222-2222-2222-000000000003', 'a1b2c3d4-1111-1111-1111-000000000001', '\\section{Methodology}\nWe propose a novel approach...', 2),
  ('b1b2c3d4-2222-2222-2222-000000000004', 'a1b2c3d4-1111-1111-1111-000000000002', '# Agenda\n- Sprint review\n- Blockers\n- Next steps', 0),
  ('b1b2c3d4-2222-2222-2222-000000000005', 'a1b2c3d4-1111-1111-1111-000000000003', '#let invoice_number = 1001\n#let client = "Acme Corp"', 0);

-- Collaborators
INSERT INTO "document_collaborators" ("id", "document_id", "user_id", "role", "status")
VALUES
  ('c1b2c3d4-3333-3333-3333-000000000001', 'a1b2c3d4-1111-1111-1111-000000000001', 'user-002', 'editor', 'accepted'),
  ('c1b2c3d4-3333-3333-3333-000000000002', 'a1b2c3d4-1111-1111-1111-000000000001', 'user-003', 'viewer', 'pending'),
  ('c1b2c3d4-3333-3333-3333-000000000003', 'a1b2c3d4-1111-1111-1111-000000000002', 'user-001', 'admin', 'accepted');

-- Document versions
INSERT INTO "document_versions" ("id", "document_id", "version_number", "snapshot", "created_by")
VALUES
  ('d1b2c3d4-4444-4444-4444-000000000001', 'a1b2c3d4-1111-1111-1111-000000000001', 1, '{"pages": 2, "title": "Thesis Draft v1"}', 'user-001'),
  ('d1b2c3d4-4444-4444-4444-000000000002', 'a1b2c3d4-1111-1111-1111-000000000001', 2, '{"pages": 3, "title": "Thesis Draft v2"}', 'user-001'),
  ('d1b2c3d4-4444-4444-4444-000000000003', 'a1b2c3d4-1111-1111-1111-000000000002', 1, '{"pages": 1, "title": "Meeting Notes v1"}', 'user-002');

-- Compiled outputs
INSERT INTO "compiled_outputs" ("id", "document_id", "version_id", "output_format", "file_path", "file_size_bytes")
VALUES
  ('e1b2c3d4-5555-5555-5555-000000000001', 'a1b2c3d4-1111-1111-1111-000000000001', 'd1b2c3d4-4444-4444-4444-000000000002', 'pdf', '/outputs/thesis-v2.pdf', 245000),
  ('e1b2c3d4-5555-5555-5555-000000000002', 'a1b2c3d4-1111-1111-1111-000000000002', 'd1b2c3d4-4444-4444-4444-000000000003', 'html', '/outputs/meeting-notes.html', 12000),
  ('e1b2c3d4-5555-5555-5555-000000000003', 'a1b2c3d4-1111-1111-1111-000000000003', NULL, 'pdf', '/outputs/invoice-1001.pdf', 98000);