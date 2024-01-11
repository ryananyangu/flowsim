CREATE TYPE "screen_types" AS ENUM ('RIS', 'LS', 'RS', 'ES');
CREATE TYPE "direction_types" AS ENUM ('INBOX', 'OUTBOX');
CREATE TABLE "menus" (
    "id" serial PRIMARY KEY,
    "shortcode" varchar(16),
    "country_code" varchar(4),
    "telco" varchar(64),
    "description" text,
    "created_at" timestamp,
    "updated_at" timestamp,
    "createdby" varchar(64),
    "updatedby" varchar
);
CREATE TABLE "screens" (
    "id" serial PRIMARY KEY,
    "menu_id" integer,
    "name" integer,
    "language" varchar(2),
    "screen_type" screen_types,
    "location" int,
    "is_end" bool,
    "back_enabled" bool,
    "exit_enabled" bool,
    "details" jsonb,
    "created_at" timestamp,
    "updated_at" timestamp,
    "createdby" varchar(64),
    "updatedby" varchar(64)
);
CREATE TABLE "messages" (
    "id" serial PRIMARY KEY,
    "screen_id" integer,
    "conversation_id" varchar(120),
    "content" text,
    "source" varchar(16),
    "destination" varchar(16),
    "direction" direction_types,
    "status" varchar(12),
    "status_description" text,
    "created_at" timestamp,
    "updated_at" timestamp,
    "createdby" varchar(64),
    "updatedby" varchar(64),
    "session_data" jsonb
);
ALTER TABLE "screens"
ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");
ALTER TABLE "messages"
ADD FOREIGN KEY ("screen_id") REFERENCES "screens" ("id");