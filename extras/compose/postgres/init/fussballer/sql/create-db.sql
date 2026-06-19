-- Aufruf:   psql --dbname=postgres --username=postgres --file=/init/patient/sql/create-db.sql

-- https://www.postgresql.org/docs/current/sql-createuser.html
-- https://www.postgresql.org/docs/current/sql-createrole.html
CREATE USER fussballer PASSWORD 'p';

-- https://www.postgresql.org/docs/current/sql-createdatabase.html
CREATE DATABASE fussballer;

-- https://www.postgresql.org/docs/current/role-attributes.html
-- https://www.postgresql.org/docs/current/ddl-priv.html
-- https://www.postgresql.org/docs/current/sql-grant.html
GRANT ALL ON DATABASE fussballer TO fussballer;

-- https://www.postgresql.org/docs/current/sql-createtablespace.html
CREATE TABLESPACE fussballerspace OWNER fussballer LOCATION '/tablespace/fussballer';
