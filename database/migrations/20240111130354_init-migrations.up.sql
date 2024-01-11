CREATE TYPE "screen_types" AS ENUM ('RIS', 'LS', 'RS', 'ES');
CREATE TYPE "direction_types" AS ENUM ('INBOX', 'OUTBOX');
CREATE TABLE "menus" (
    "id" serial PRIMARY KEY,
    "shortcode" varchar(16) NOT NULL,
    "country_code" varchar(3) NOT NULL,
    "telco" varchar(64) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "createdby" varchar(64) NOT NULL,
    "updatedby" varchar(64) NOT NULL
);
CREATE TABLE "screens" (
    "id" serial PRIMARY KEY,
    "menu_id" integer NOT NULL,
    "name" varchar(16) NOT NULL,
    "language" varchar(2) NOT NULL,
    "screen_type" screen_types NOT NULL,
    "location" int NOT NULL,
    "is_end" bool NOT NULL,
    "back_enabled" bool NOT NULL,
    "exit_enabled" bool NOT NULL,
    "details" jsonb NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "createdby" varchar(64) NOT NULL,
    "updatedby" varchar(64) NOT NULL
);
CREATE TABLE "messages" (
    "id" serial PRIMARY KEY,
    "screen_id" integer NOT NULL,
    "conversation_id" varchar(120) NOT NULL,
    "content" text,
    "source" varchar(16) NOT NULL,
    "destination" varchar(16) NOT NULL,
    "direction" direction_types NOT NULL,
    "status" varchar(12) NOT NULL,
    "status_description" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "createdby" varchar(64) NOT NULL,
    "updatedby" varchar(64) NOT NULL,
    "session_data" jsonb
);
ALTER TABLE "screens"
ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");
ALTER TABLE "messages"
ADD FOREIGN KEY ("screen_id") REFERENCES "screens" ("id");