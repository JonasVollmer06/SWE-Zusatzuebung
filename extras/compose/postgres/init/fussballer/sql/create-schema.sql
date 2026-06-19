-- Aufruf:   psql --dbname=patient --username=patient --file=/init/patient/sql/create-schema.sql

-- https://www.postgresql.org/docs/devel/app-psql.html
-- https://www.postgresql.org/docs/current/ddl-schemas.html
-- https://www.postgresql.org/docs/current/ddl-schemas.html#DDL-SCHEMAS-CREATE
-- "user-private schema" (Default-Schema: public)
CREATE SCHEMA IF NOT EXISTS AUTHORIZATION fussballer;

ALTER ROLE fussballer SET search_path = 'fussballer';
