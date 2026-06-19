-- Aufruf:   psql --dbname=buch --username=postgres --file=/init/buch/sql/copy-csv.sql

SET search_path TO fussballer;

-- https://www.postgresql.org/docs/current/sql-copy.html
COPY fussballer FROM '/init/fussballer/csv/fussballer.csv' (FORMAT csv, DELIMITER ';', HEADER true);
COPY adresse FROM '/init/fussballer/csv/adresse.csv' (FORMAT csv, DELIMITER ';', HEADER true);
COPY auszeichnung FROM '/init/fussballer/csv/auszeichnung.csv' (FORMAT csv, DELIMITER ';', HEADER true);
