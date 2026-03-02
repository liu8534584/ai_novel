/*
 Navicat Premium Dump SQL

 Source Server         : pgsql
 Source Server Type    : PostgreSQL
 Source Server Version : 160011 (160011)
 Source Host           : localhost:5432
 Source Catalog        : ai_novel
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160011 (160011)
 File Encoding         : 65001

 Date: 02/03/2026 11:43:24
*/


-- ----------------------------
-- Type structure for halfvec
-- ----------------------------
DROP TYPE IF EXISTS "public"."halfvec";
CREATE TYPE "public"."halfvec" (
  INPUT = "public"."halfvec_in",
  OUTPUT = "public"."halfvec_out",
  RECEIVE = "public"."halfvec_recv",
  SEND = "public"."halfvec_send",
  TYPMOD_IN = "public"."halfvec_typmod_in",
  INTERNALLENGTH = VARIABLE,
  STORAGE = external,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."halfvec" OWNER TO "postgres";

-- ----------------------------
-- Type structure for sparsevec
-- ----------------------------
DROP TYPE IF EXISTS "public"."sparsevec";
CREATE TYPE "public"."sparsevec" (
  INPUT = "public"."sparsevec_in",
  OUTPUT = "public"."sparsevec_out",
  RECEIVE = "public"."sparsevec_recv",
  SEND = "public"."sparsevec_send",
  TYPMOD_IN = "public"."sparsevec_typmod_in",
  INTERNALLENGTH = VARIABLE,
  STORAGE = external,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."sparsevec" OWNER TO "postgres";

-- ----------------------------
-- Type structure for vector
-- ----------------------------
DROP TYPE IF EXISTS "public"."vector";
CREATE TYPE "public"."vector" (
  INPUT = "public"."vector_in",
  OUTPUT = "public"."vector_out",
  RECEIVE = "public"."vector_recv",
  SEND = "public"."vector_send",
  TYPMOD_IN = "public"."vector_typmod_in",
  INTERNALLENGTH = VARIABLE,
  STORAGE = external,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."vector" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for books_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."books_id_seq";
CREATE SEQUENCE "public"."books_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."books_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for chapter_health_scores_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."chapter_health_scores_id_seq";
CREATE SEQUENCE "public"."chapter_health_scores_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."chapter_health_scores_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for chapter_versions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."chapter_versions_id_seq";
CREATE SEQUENCE "public"."chapter_versions_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."chapter_versions_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for chapters_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."chapters_id_seq";
CREATE SEQUENCE "public"."chapters_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."chapters_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for character_anchors_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."character_anchors_id_seq";
CREATE SEQUENCE "public"."character_anchors_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."character_anchors_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for character_state_records_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."character_state_records_id_seq";
CREATE SEQUENCE "public"."character_state_records_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."character_state_records_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for characters_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."characters_id_seq";
CREATE SEQUENCE "public"."characters_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."characters_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for foreshadowings_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."foreshadowings_id_seq";
CREATE SEQUENCE "public"."foreshadowings_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."foreshadowings_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for ooc_scores_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ooc_scores_id_seq";
CREATE SEQUENCE "public"."ooc_scores_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."ooc_scores_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for outline_versions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."outline_versions_id_seq";
CREATE SEQUENCE "public"."outline_versions_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."outline_versions_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for prompt_templates_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."prompt_templates_id_seq";
CREATE SEQUENCE "public"."prompt_templates_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."prompt_templates_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for story_contradictions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."story_contradictions_id_seq";
CREATE SEQUENCE "public"."story_contradictions_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."story_contradictions_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for story_events_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."story_events_id_seq";
CREATE SEQUENCE "public"."story_events_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."story_events_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for vector_records_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."vector_records_id_seq";
CREATE SEQUENCE "public"."vector_records_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."vector_records_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for vs_characters_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."vs_characters_id_seq";
CREATE SEQUENCE "public"."vs_characters_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."vs_characters_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for vs_history_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."vs_history_id_seq";
CREATE SEQUENCE "public"."vs_history_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."vs_history_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for vs_outlines_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."vs_outlines_id_seq";
CREATE SEQUENCE "public"."vs_outlines_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."vs_outlines_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for vs_world_rules_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."vs_world_rules_id_seq";
CREATE SEQUENCE "public"."vs_world_rules_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."vs_world_rules_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Table structure for books
-- ----------------------------
DROP TABLE IF EXISTS "public"."books";
CREATE TABLE "public"."books" (
  "id" int8 NOT NULL DEFAULT nextval('books_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "title" text COLLATE "pg_catalog"."default" NOT NULL,
  "author" text COLLATE "pg_catalog"."default",
  "genre" text COLLATE "pg_catalog"."default",
  "tags" text COLLATE "pg_catalog"."default",
  "language" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "total_chapters" int4 DEFAULT 1,
  "status" text COLLATE "pg_catalog"."default" DEFAULT 'draft'::text,
  "world_setting" text COLLATE "pg_catalog"."default",
  "current_state" text COLLATE "pg_catalog"."default",
  "llm_config" text COLLATE "pg_catalog"."default",
  "prompt_bindings" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."books" OWNER TO "postgres";

-- ----------------------------
-- Table structure for chapter_health_scores
-- ----------------------------
DROP TABLE IF EXISTS "public"."chapter_health_scores";
CREATE TABLE "public"."chapter_health_scores" (
  "id" int8 NOT NULL DEFAULT nextval('chapter_health_scores_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "chapter_id" int8,
  "ooc_score" numeric,
  "event_consistency" numeric,
  "foreshadowing" numeric,
  "total_health" numeric,
  "audit_report" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."chapter_health_scores" OWNER TO "postgres";

-- ----------------------------
-- Table structure for chapter_versions
-- ----------------------------
DROP TABLE IF EXISTS "public"."chapter_versions";
CREATE TABLE "public"."chapter_versions" (
  "id" int8 NOT NULL DEFAULT nextval('chapter_versions_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "chapter_id" int8,
  "version" int8,
  "title" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "word_count" int8,
  "summary" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."chapter_versions" OWNER TO "postgres";

-- ----------------------------
-- Table structure for chapters
-- ----------------------------
DROP TABLE IF EXISTS "public"."chapters";
CREATE TABLE "public"."chapters" (
  "id" int8 NOT NULL DEFAULT nextval('chapters_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "title" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "order" int8,
  "objective" text COLLATE "pg_catalog"."default",
  "summary" text COLLATE "pg_catalog"."default",
  "user_intent" text COLLATE "pg_catalog"."default",
  "outline" text COLLATE "pg_catalog"."default",
  "is_outline_confirmed" bool DEFAULT false,
  "current_version" int4 DEFAULT 1
)
;
ALTER TABLE "public"."chapters" OWNER TO "postgres";

-- ----------------------------
-- Table structure for character_anchors
-- ----------------------------
DROP TABLE IF EXISTS "public"."character_anchors";
CREATE TABLE "public"."character_anchors" (
  "id" int8 NOT NULL DEFAULT nextval('character_anchors_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "character_id" int8,
  "personality_labels" text COLLATE "pg_catalog"."default",
  "core_motivation" text COLLATE "pg_catalog"."default",
  "behavior_bottom_line" text COLLATE "pg_catalog"."default",
  "decision_tendency" text COLLATE "pg_catalog"."default",
  "emotional_triggers" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."character_anchors" OWNER TO "postgres";

-- ----------------------------
-- Table structure for character_state_records
-- ----------------------------
DROP TABLE IF EXISTS "public"."character_state_records";
CREATE TABLE "public"."character_state_records" (
  "id" int8 NOT NULL DEFAULT nextval('character_state_records_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "character_id" int8,
  "chapter_id" int8,
  "state" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."character_state_records" OWNER TO "postgres";

-- ----------------------------
-- Table structure for characters
-- ----------------------------
DROP TABLE IF EXISTS "public"."characters";
CREATE TABLE "public"."characters" (
  "id" int8 NOT NULL DEFAULT nextval('characters_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "name" text COLLATE "pg_catalog"."default",
  "role" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "dynamic_state" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."characters" OWNER TO "postgres";

-- ----------------------------
-- Table structure for foreshadowings
-- ----------------------------
DROP TABLE IF EXISTS "public"."foreshadowings";
CREATE TABLE "public"."foreshadowings" (
  "id" int8 NOT NULL DEFAULT nextval('foreshadowings_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "chapter_id" int8,
  "chapter_index" int8,
  "event_type" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "involved_characters" text COLLATE "pg_catalog"."default",
  "direct_consequence" text COLLATE "pg_catalog"."default",
  "unresolved_impact" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "importance" int8,
  "last_referenced_chapter" int8,
  "resolved_chapter_index" int8,
  "resolve_reason" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."foreshadowings" OWNER TO "postgres";

-- ----------------------------
-- Table structure for ooc_scores
-- ----------------------------
DROP TABLE IF EXISTS "public"."ooc_scores";
CREATE TABLE "public"."ooc_scores" (
  "id" int8 NOT NULL DEFAULT nextval('ooc_scores_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "character_id" int8,
  "chapter_id" int8,
  "personality_consistency" numeric,
  "motivation_consistency" numeric,
  "emotional_reasonability" numeric,
  "cost_missing" numeric,
  "total_score" numeric,
  "conclusion" text COLLATE "pg_catalog"."default",
  "explanation" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."ooc_scores" OWNER TO "postgres";

-- ----------------------------
-- Table structure for outline_versions
-- ----------------------------
DROP TABLE IF EXISTS "public"."outline_versions";
CREATE TABLE "public"."outline_versions" (
  "id" int8 NOT NULL DEFAULT nextval('outline_versions_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "version" int8,
  "world_view" text COLLATE "pg_catalog"."default",
  "outline" text COLLATE "pg_catalog"."default",
  "characters" text COLLATE "pg_catalog"."default",
  "titles" text COLLATE "pg_catalog"."default",
  "is_selected" bool DEFAULT false,
  "is_locked" bool DEFAULT false
)
;
ALTER TABLE "public"."outline_versions" OWNER TO "postgres";

-- ----------------------------
-- Table structure for prompt_templates
-- ----------------------------
DROP TABLE IF EXISTS "public"."prompt_templates";
CREATE TABLE "public"."prompt_templates" (
  "id" int8 NOT NULL DEFAULT nextval('prompt_templates_id_seq'::regclass),
  "key" text COLLATE "pg_catalog"."default" NOT NULL,
  "title" text COLLATE "pg_catalog"."default" NOT NULL,
  "category" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "updated_at" timestamptz(6),
  "source" text COLLATE "pg_catalog"."default",
  "enabled" bool DEFAULT true
)
;
ALTER TABLE "public"."prompt_templates" OWNER TO "postgres";

-- ----------------------------
-- Table structure for story_contradictions
-- ----------------------------
DROP TABLE IF EXISTS "public"."story_contradictions";
CREATE TABLE "public"."story_contradictions" (
  "id" int8 NOT NULL DEFAULT nextval('story_contradictions_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "chapter_id" int8,
  "type" text COLLATE "pg_catalog"."default",
  "severity" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "reference" text COLLATE "pg_catalog"."default",
  "suggestion" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."story_contradictions" OWNER TO "postgres";

-- ----------------------------
-- Table structure for story_events
-- ----------------------------
DROP TABLE IF EXISTS "public"."story_events";
CREATE TABLE "public"."story_events" (
  "id" int8 NOT NULL DEFAULT nextval('story_events_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "chapter_id" int8,
  "chapter_index" int8,
  "event_type" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "involved_characters" text COLLATE "pg_catalog"."default",
  "direct_consequence" text COLLATE "pg_catalog"."default",
  "unresolved_impact" text COLLATE "pg_catalog"."default",
  "importance" int8
)
;
ALTER TABLE "public"."story_events" OWNER TO "postgres";

-- ----------------------------
-- Table structure for vector_records
-- ----------------------------
DROP TABLE IF EXISTS "public"."vector_records";
CREATE TABLE "public"."vector_records" (
  "id" int8 NOT NULL DEFAULT nextval('vector_records_id_seq'::regclass),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "book_id" int8,
  "chapter_id" int8,
  "category" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "embedding" text COLLATE "pg_catalog"."default",
  "metadata" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."vector_records" OWNER TO "postgres";

-- ----------------------------
-- Table structure for vs_characters
-- ----------------------------
DROP TABLE IF EXISTS "public"."vs_characters";
CREATE TABLE "public"."vs_characters" (
  "id" int8 NOT NULL DEFAULT nextval('vs_characters_id_seq'::regclass),
  "content" text COLLATE "pg_catalog"."default",
  "metadata" jsonb,
  "vector" "public"."vector"
)
;
ALTER TABLE "public"."vs_characters" OWNER TO "postgres";

-- ----------------------------
-- Table structure for vs_history
-- ----------------------------
DROP TABLE IF EXISTS "public"."vs_history";
CREATE TABLE "public"."vs_history" (
  "id" int8 NOT NULL DEFAULT nextval('vs_history_id_seq'::regclass),
  "content" text COLLATE "pg_catalog"."default",
  "metadata" jsonb,
  "vector" "public"."vector"
)
;
ALTER TABLE "public"."vs_history" OWNER TO "postgres";

-- ----------------------------
-- Table structure for vs_outlines
-- ----------------------------
DROP TABLE IF EXISTS "public"."vs_outlines";
CREATE TABLE "public"."vs_outlines" (
  "id" int8 NOT NULL DEFAULT nextval('vs_outlines_id_seq'::regclass),
  "content" text COLLATE "pg_catalog"."default",
  "metadata" jsonb,
  "vector" "public"."vector"
)
;
ALTER TABLE "public"."vs_outlines" OWNER TO "postgres";

-- ----------------------------
-- Table structure for vs_world_rules
-- ----------------------------
DROP TABLE IF EXISTS "public"."vs_world_rules";
CREATE TABLE "public"."vs_world_rules" (
  "id" int8 NOT NULL DEFAULT nextval('vs_world_rules_id_seq'::regclass),
  "content" text COLLATE "pg_catalog"."default",
  "metadata" jsonb,
  "vector" "public"."vector"
)
;
ALTER TABLE "public"."vs_world_rules" OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_numeric, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_numeric, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_halfvec"(_numeric, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_int4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_int4, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_halfvec"(_int4, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_float4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_float4, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_halfvec"(_float4, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_float8, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_float8, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_halfvec"(_float8, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_int4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_int4, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_sparsevec"(_int4, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_numeric, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_numeric, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_sparsevec"(_numeric, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_float8, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_float8, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_sparsevec"(_float8, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_float4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_float4, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_sparsevec"(_float4, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_float4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_float4, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_vector"(_float4, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_numeric, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_numeric, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_vector"(_numeric, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_float8, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_float8, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_vector"(_float8, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_int4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_int4, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."array_to_vector"(_int4, int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for binary_quantize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."binary_quantize"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."binary_quantize"("public"."halfvec")
  RETURNS "pg_catalog"."bit" AS '$libdir/vector', 'halfvec_binary_quantize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."binary_quantize"("public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for binary_quantize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."binary_quantize"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."binary_quantize"("public"."vector")
  RETURNS "pg_catalog"."bit" AS '$libdir/vector', 'binary_quantize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."binary_quantize"("public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for cosine_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."cosine_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."cosine_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_cosine_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."cosine_distance"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for cosine_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."cosine_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."cosine_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_cosine_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."cosine_distance"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for cosine_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."cosine_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."cosine_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'cosine_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."cosine_distance"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec"("public"."halfvec", int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec"("public"."halfvec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_accum
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_accum"(_float8, "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_accum"(_float8, "public"."halfvec")
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'halfvec_accum'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_accum"(_float8, "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_add
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_add"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_add"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_add'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_add"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_avg
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_avg"(_float8);
CREATE OR REPLACE FUNCTION "public"."halfvec_avg"(_float8)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_avg'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_avg"(_float8) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_cmp"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_cmp"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'halfvec_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_cmp"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_combine
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_combine"(_float8, _float8);
CREATE OR REPLACE FUNCTION "public"."halfvec_combine"(_float8, _float8)
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'vector_combine'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_combine"(_float8, _float8) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_concat
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_concat"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_concat"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_concat'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_concat"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_eq
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_eq"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_eq"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_eq'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_eq"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_ge"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_ge"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_ge"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_gt"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_gt"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_gt"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_in"(cstring, oid, int4);
CREATE OR REPLACE FUNCTION "public"."halfvec_in"(cstring, oid, int4)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_in"(cstring, oid, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_l2_squared_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_l2_squared_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_l2_squared_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l2_squared_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_l2_squared_distance"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_le"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_le"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_le"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_lt"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_lt"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_lt"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_mul
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_mul"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_mul"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_mul'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_mul"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_ne
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_ne"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_ne"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_ne'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_ne"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_negative_inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_negative_inner_product"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_negative_inner_product"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_negative_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_negative_inner_product"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_out"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_out"("public"."halfvec")
  RETURNS "pg_catalog"."cstring" AS '$libdir/vector', 'halfvec_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_out"("public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_recv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_recv"(internal, oid, int4);
CREATE OR REPLACE FUNCTION "public"."halfvec_recv"(internal, oid, int4)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_recv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_recv"(internal, oid, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_send
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_send"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_send"("public"."halfvec")
  RETURNS "pg_catalog"."bytea" AS '$libdir/vector', 'halfvec_send'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_send"("public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_spherical_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_spherical_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_spherical_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_spherical_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_spherical_distance"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_sub
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_sub"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_sub"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_sub'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_sub"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_to_float4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_to_float4"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec_to_float4"("public"."halfvec", int4, bool)
  RETURNS "pg_catalog"."_float4" AS '$libdir/vector', 'halfvec_to_float4'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_to_float4"("public"."halfvec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_to_sparsevec"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec_to_sparsevec"("public"."halfvec", int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'halfvec_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_to_sparsevec"("public"."halfvec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_to_vector"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec_to_vector"("public"."halfvec", int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'halfvec_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_to_vector"("public"."halfvec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for halfvec_typmod_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_typmod_in"(_cstring);
CREATE OR REPLACE FUNCTION "public"."halfvec_typmod_in"(_cstring)
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'halfvec_typmod_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."halfvec_typmod_in"(_cstring) OWNER TO "postgres";

-- ----------------------------
-- Function structure for hamming_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hamming_distance"(bit, bit);
CREATE OR REPLACE FUNCTION "public"."hamming_distance"(bit, bit)
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'hamming_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."hamming_distance"(bit, bit) OWNER TO "postgres";

-- ----------------------------
-- Function structure for hnsw_bit_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnsw_bit_support"(internal);
CREATE OR REPLACE FUNCTION "public"."hnsw_bit_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'hnsw_bit_support'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."hnsw_bit_support"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for hnsw_halfvec_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnsw_halfvec_support"(internal);
CREATE OR REPLACE FUNCTION "public"."hnsw_halfvec_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'hnsw_halfvec_support'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."hnsw_halfvec_support"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for hnsw_sparsevec_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnsw_sparsevec_support"(internal);
CREATE OR REPLACE FUNCTION "public"."hnsw_sparsevec_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'hnsw_sparsevec_support'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."hnsw_sparsevec_support"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for hnswhandler
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnswhandler"(internal);
CREATE OR REPLACE FUNCTION "public"."hnswhandler"(internal)
  RETURNS "pg_catalog"."index_am_handler" AS '$libdir/vector', 'hnswhandler'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."hnswhandler"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."inner_product"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."inner_product"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."inner_product"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."inner_product"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."inner_product"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."inner_product"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."inner_product"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."inner_product"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."inner_product"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for ivfflat_bit_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."ivfflat_bit_support"(internal);
CREATE OR REPLACE FUNCTION "public"."ivfflat_bit_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'ivfflat_bit_support'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."ivfflat_bit_support"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for ivfflat_halfvec_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."ivfflat_halfvec_support"(internal);
CREATE OR REPLACE FUNCTION "public"."ivfflat_halfvec_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'ivfflat_halfvec_support'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."ivfflat_halfvec_support"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for ivfflathandler
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."ivfflathandler"(internal);
CREATE OR REPLACE FUNCTION "public"."ivfflathandler"(internal)
  RETURNS "pg_catalog"."index_am_handler" AS '$libdir/vector', 'ivfflathandler'
  LANGUAGE c VOLATILE
  COST 1;
ALTER FUNCTION "public"."ivfflathandler"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for jaccard_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."jaccard_distance"(bit, bit);
CREATE OR REPLACE FUNCTION "public"."jaccard_distance"(bit, bit)
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'jaccard_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."jaccard_distance"(bit, bit) OWNER TO "postgres";

-- ----------------------------
-- Function structure for l1_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l1_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l1_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l1_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l1_distance"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l1_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l1_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l1_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l1_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l1_distance"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l1_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l1_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."l1_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'l1_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l1_distance"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l2_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l2_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_distance"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l2_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l2_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_distance"("public"."halfvec", "public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."l2_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'l2_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_distance"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_norm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_norm"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l2_norm"("public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l2_norm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_norm"("public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_norm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_norm"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l2_norm"("public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l2_norm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_norm"("public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_normalize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_normalize"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."l2_normalize"("public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'l2_normalize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_normalize"("public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_normalize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_normalize"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l2_normalize"("public"."sparsevec")
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec_l2_normalize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_normalize"("public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for l2_normalize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_normalize"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l2_normalize"("public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_l2_normalize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."l2_normalize"("public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec"("public"."sparsevec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."sparsevec"("public"."sparsevec", int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec"("public"."sparsevec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_cmp"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_cmp"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'sparsevec_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_cmp"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_eq
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_eq"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_eq"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_eq'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_eq"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_ge"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_ge"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_ge"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_gt"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_gt"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_gt"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_in"(cstring, oid, int4);
CREATE OR REPLACE FUNCTION "public"."sparsevec_in"(cstring, oid, int4)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_in"(cstring, oid, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_l2_squared_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_l2_squared_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_l2_squared_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l2_squared_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_l2_squared_distance"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_le"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_le"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_le"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_lt"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_lt"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_lt"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_ne
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_ne"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_ne"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_ne'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_ne"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_negative_inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_negative_inner_product"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_negative_inner_product"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_negative_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_negative_inner_product"("public"."sparsevec", "public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_out"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_out"("public"."sparsevec")
  RETURNS "pg_catalog"."cstring" AS '$libdir/vector', 'sparsevec_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_out"("public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_recv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_recv"(internal, oid, int4);
CREATE OR REPLACE FUNCTION "public"."sparsevec_recv"(internal, oid, int4)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec_recv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_recv"(internal, oid, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_send
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_send"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_send"("public"."sparsevec")
  RETURNS "pg_catalog"."bytea" AS '$libdir/vector', 'sparsevec_send'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_send"("public"."sparsevec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_to_halfvec"("public"."sparsevec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."sparsevec_to_halfvec"("public"."sparsevec", int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'sparsevec_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_to_halfvec"("public"."sparsevec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_to_vector"("public"."sparsevec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."sparsevec_to_vector"("public"."sparsevec", int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'sparsevec_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_to_vector"("public"."sparsevec", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for sparsevec_typmod_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_typmod_in"(_cstring);
CREATE OR REPLACE FUNCTION "public"."sparsevec_typmod_in"(_cstring)
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'sparsevec_typmod_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."sparsevec_typmod_in"(_cstring) OWNER TO "postgres";

-- ----------------------------
-- Function structure for subvector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."subvector"("public"."halfvec", int4, int4);
CREATE OR REPLACE FUNCTION "public"."subvector"("public"."halfvec", int4, int4)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_subvector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."subvector"("public"."halfvec", int4, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for subvector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."subvector"("public"."vector", int4, int4);
CREATE OR REPLACE FUNCTION "public"."subvector"("public"."vector", int4, int4)
  RETURNS "public"."vector" AS '$libdir/vector', 'subvector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."subvector"("public"."vector", int4, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector"("public"."vector", int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector"("public"."vector", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_accum
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_accum"(_float8, "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_accum"(_float8, "public"."vector")
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'vector_accum'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_accum"(_float8, "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_add
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_add"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_add"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_add'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_add"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_avg
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_avg"(_float8);
CREATE OR REPLACE FUNCTION "public"."vector_avg"(_float8)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_avg'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_avg"(_float8) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_cmp"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_cmp"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'vector_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_cmp"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_combine
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_combine"(_float8, _float8);
CREATE OR REPLACE FUNCTION "public"."vector_combine"(_float8, _float8)
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'vector_combine'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_combine"(_float8, _float8) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_concat
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_concat"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_concat"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_concat'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_concat"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_dims
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_dims"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_dims"("public"."vector")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'vector_dims'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_dims"("public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_dims
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_dims"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."vector_dims"("public"."halfvec")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'halfvec_vector_dims'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_dims"("public"."halfvec") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_eq
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_eq"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_eq"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_eq'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_eq"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_ge"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_ge"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_ge"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_gt"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_gt"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_gt"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_in"(cstring, oid, int4);
CREATE OR REPLACE FUNCTION "public"."vector_in"(cstring, oid, int4)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_in"(cstring, oid, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_l2_squared_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_l2_squared_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_l2_squared_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_l2_squared_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_l2_squared_distance"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_le"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_le"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_le"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_lt"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_lt"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_lt"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_mul
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_mul"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_mul"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_mul'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_mul"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_ne
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_ne"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_ne"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_ne'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_ne"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_negative_inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_negative_inner_product"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_negative_inner_product"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_negative_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_negative_inner_product"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_norm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_norm"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_norm"("public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_norm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_norm"("public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_out"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_out"("public"."vector")
  RETURNS "pg_catalog"."cstring" AS '$libdir/vector', 'vector_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_out"("public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_recv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_recv"(internal, oid, int4);
CREATE OR REPLACE FUNCTION "public"."vector_recv"(internal, oid, int4)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_recv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_recv"(internal, oid, int4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_send
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_send"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_send"("public"."vector")
  RETURNS "pg_catalog"."bytea" AS '$libdir/vector', 'vector_send'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_send"("public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_spherical_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_spherical_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_spherical_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_spherical_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_spherical_distance"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_sub
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_sub"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_sub"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_sub'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_sub"("public"."vector", "public"."vector") OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_to_float4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_to_float4"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector_to_float4"("public"."vector", int4, bool)
  RETURNS "pg_catalog"."_float4" AS '$libdir/vector', 'vector_to_float4'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_to_float4"("public"."vector", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_to_halfvec"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector_to_halfvec"("public"."vector", int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'vector_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_to_halfvec"("public"."vector", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_to_sparsevec"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector_to_sparsevec"("public"."vector", int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'vector_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_to_sparsevec"("public"."vector", int4, bool) OWNER TO "postgres";

-- ----------------------------
-- Function structure for vector_typmod_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_typmod_in"(_cstring);
CREATE OR REPLACE FUNCTION "public"."vector_typmod_in"(_cstring)
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'vector_typmod_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."vector_typmod_in"(_cstring) OWNER TO "postgres";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."books_id_seq"
OWNED BY "public"."books"."id";
SELECT setval('"public"."books_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."chapter_health_scores_id_seq"
OWNED BY "public"."chapter_health_scores"."id";
SELECT setval('"public"."chapter_health_scores_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."chapter_versions_id_seq"
OWNED BY "public"."chapter_versions"."id";
SELECT setval('"public"."chapter_versions_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."chapters_id_seq"
OWNED BY "public"."chapters"."id";
SELECT setval('"public"."chapters_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."character_anchors_id_seq"
OWNED BY "public"."character_anchors"."id";
SELECT setval('"public"."character_anchors_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."character_state_records_id_seq"
OWNED BY "public"."character_state_records"."id";
SELECT setval('"public"."character_state_records_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."characters_id_seq"
OWNED BY "public"."characters"."id";
SELECT setval('"public"."characters_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."foreshadowings_id_seq"
OWNED BY "public"."foreshadowings"."id";
SELECT setval('"public"."foreshadowings_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ooc_scores_id_seq"
OWNED BY "public"."ooc_scores"."id";
SELECT setval('"public"."ooc_scores_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."outline_versions_id_seq"
OWNED BY "public"."outline_versions"."id";
SELECT setval('"public"."outline_versions_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."prompt_templates_id_seq"
OWNED BY "public"."prompt_templates"."id";
SELECT setval('"public"."prompt_templates_id_seq"', 22, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."story_contradictions_id_seq"
OWNED BY "public"."story_contradictions"."id";
SELECT setval('"public"."story_contradictions_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."story_events_id_seq"
OWNED BY "public"."story_events"."id";
SELECT setval('"public"."story_events_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."vector_records_id_seq"
OWNED BY "public"."vector_records"."id";
SELECT setval('"public"."vector_records_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."vs_characters_id_seq"
OWNED BY "public"."vs_characters"."id";
SELECT setval('"public"."vs_characters_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."vs_history_id_seq"
OWNED BY "public"."vs_history"."id";
SELECT setval('"public"."vs_history_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."vs_outlines_id_seq"
OWNED BY "public"."vs_outlines"."id";
SELECT setval('"public"."vs_outlines_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."vs_world_rules_id_seq"
OWNED BY "public"."vs_world_rules"."id";
SELECT setval('"public"."vs_world_rules_id_seq"', 1, false);

-- ----------------------------
-- Indexes structure for table books
-- ----------------------------
CREATE INDEX "idx_books_deleted_at" ON "public"."books" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_books_genre" ON "public"."books" USING btree (
  "genre" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_books_status" ON "public"."books" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_books_title" ON "public"."books" USING btree (
  "title" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table books
-- ----------------------------
ALTER TABLE "public"."books" ADD CONSTRAINT "books_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table chapter_health_scores
-- ----------------------------
CREATE INDEX "idx_chapter_health_scores_book_id" ON "public"."chapter_health_scores" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_chapter_health_scores_chapter_id" ON "public"."chapter_health_scores" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_chapter_health_scores_deleted_at" ON "public"."chapter_health_scores" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chapter_health_scores
-- ----------------------------
ALTER TABLE "public"."chapter_health_scores" ADD CONSTRAINT "chapter_health_scores_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table chapter_versions
-- ----------------------------
CREATE INDEX "idx_chapter_versions_chapter_id" ON "public"."chapter_versions" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_chapter_versions_deleted_at" ON "public"."chapter_versions" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chapter_versions
-- ----------------------------
ALTER TABLE "public"."chapter_versions" ADD CONSTRAINT "chapter_versions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table chapters
-- ----------------------------
CREATE INDEX "idx_chapters_book_id" ON "public"."chapters" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_chapters_deleted_at" ON "public"."chapters" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chapters
-- ----------------------------
ALTER TABLE "public"."chapters" ADD CONSTRAINT "chapters_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table character_anchors
-- ----------------------------
CREATE INDEX "idx_character_anchors_character_id" ON "public"."character_anchors" USING btree (
  "character_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_character_anchors_deleted_at" ON "public"."character_anchors" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table character_anchors
-- ----------------------------
ALTER TABLE "public"."character_anchors" ADD CONSTRAINT "character_anchors_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table character_state_records
-- ----------------------------
CREATE INDEX "idx_character_state_records_chapter_id" ON "public"."character_state_records" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_character_state_records_character_id" ON "public"."character_state_records" USING btree (
  "character_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_character_state_records_deleted_at" ON "public"."character_state_records" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table character_state_records
-- ----------------------------
ALTER TABLE "public"."character_state_records" ADD CONSTRAINT "character_state_records_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table characters
-- ----------------------------
CREATE INDEX "idx_characters_book_id" ON "public"."characters" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_characters_deleted_at" ON "public"."characters" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table characters
-- ----------------------------
ALTER TABLE "public"."characters" ADD CONSTRAINT "characters_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table foreshadowings
-- ----------------------------
CREATE INDEX "idx_foreshadowings_book_id" ON "public"."foreshadowings" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_foreshadowings_chapter_id" ON "public"."foreshadowings" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_foreshadowings_deleted_at" ON "public"."foreshadowings" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_foreshadowings_status" ON "public"."foreshadowings" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table foreshadowings
-- ----------------------------
ALTER TABLE "public"."foreshadowings" ADD CONSTRAINT "foreshadowings_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ooc_scores
-- ----------------------------
CREATE INDEX "idx_ooc_scores_chapter_id" ON "public"."ooc_scores" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ooc_scores_character_id" ON "public"."ooc_scores" USING btree (
  "character_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ooc_scores_deleted_at" ON "public"."ooc_scores" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ooc_scores
-- ----------------------------
ALTER TABLE "public"."ooc_scores" ADD CONSTRAINT "ooc_scores_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table outline_versions
-- ----------------------------
CREATE INDEX "idx_outline_versions_book_id" ON "public"."outline_versions" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_outline_versions_deleted_at" ON "public"."outline_versions" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table outline_versions
-- ----------------------------
ALTER TABLE "public"."outline_versions" ADD CONSTRAINT "outline_versions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table prompt_templates
-- ----------------------------
ALTER TABLE "public"."prompt_templates" ADD CONSTRAINT "prompt_templates_key_key" UNIQUE ("key");

-- ----------------------------
-- Primary Key structure for table prompt_templates
-- ----------------------------
ALTER TABLE "public"."prompt_templates" ADD CONSTRAINT "prompt_templates_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table story_contradictions
-- ----------------------------
CREATE INDEX "idx_story_contradictions_book_id" ON "public"."story_contradictions" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_story_contradictions_chapter_id" ON "public"."story_contradictions" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_story_contradictions_deleted_at" ON "public"."story_contradictions" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table story_contradictions
-- ----------------------------
ALTER TABLE "public"."story_contradictions" ADD CONSTRAINT "story_contradictions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table story_events
-- ----------------------------
CREATE INDEX "idx_story_events_book_id" ON "public"."story_events" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_story_events_chapter_id" ON "public"."story_events" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_story_events_deleted_at" ON "public"."story_events" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table story_events
-- ----------------------------
ALTER TABLE "public"."story_events" ADD CONSTRAINT "story_events_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table vector_records
-- ----------------------------
CREATE INDEX "idx_vector_records_book_id" ON "public"."vector_records" USING btree (
  "book_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_vector_records_category" ON "public"."vector_records" USING btree (
  "category" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_vector_records_chapter_id" ON "public"."vector_records" USING btree (
  "chapter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_vector_records_deleted_at" ON "public"."vector_records" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table vector_records
-- ----------------------------
ALTER TABLE "public"."vector_records" ADD CONSTRAINT "vector_records_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table vs_characters
-- ----------------------------
ALTER TABLE "public"."vs_characters" ADD CONSTRAINT "vs_characters_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table vs_history
-- ----------------------------
ALTER TABLE "public"."vs_history" ADD CONSTRAINT "vs_history_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table vs_outlines
-- ----------------------------
ALTER TABLE "public"."vs_outlines" ADD CONSTRAINT "vs_outlines_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table vs_world_rules
-- ----------------------------
ALTER TABLE "public"."vs_world_rules" ADD CONSTRAINT "vs_world_rules_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table chapter_versions
-- ----------------------------
ALTER TABLE "public"."chapter_versions" ADD CONSTRAINT "fk_chapters_versions" FOREIGN KEY ("chapter_id") REFERENCES "public"."chapters" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table chapters
-- ----------------------------
ALTER TABLE "public"."chapters" ADD CONSTRAINT "fk_books_chapters" FOREIGN KEY ("book_id") REFERENCES "public"."books" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table characters
-- ----------------------------
ALTER TABLE "public"."characters" ADD CONSTRAINT "fk_books_characters" FOREIGN KEY ("book_id") REFERENCES "public"."books" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
