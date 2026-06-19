SET default_tablespace = fussballerspace;

CREATE SCHEMA IF NOT EXISTS AUTHORIZATION fussballer;

ALTER ROLE fussballer SET search_path = 'fussballer';
set search_path to 'fussballer';

CREATE TYPE position_enum AS ENUM ('TORWART', 'VERTEIDIGER', 'MITTELFELDSPIELER', 'STUERMER');

CREATE TABLE IF NOT EXISTS fussballer (
    id            integer GENERATED ALWAYS AS IDENTITY(START WITH 1000) PRIMARY KEY,
    version       integer NOT NULL DEFAULT 0,
    nachname      text NOT NULL,
    nationalitaet text NOT NULL,
    position      position_enum,
    geburtsdatum  date NOT NULL,
    username      text NOT NULL UNIQUE,
    erzeugt       timestamp NOT NULL DEFAULT NOW(),
    aktualisiert  timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS adresse (
    id             integer GENERATED ALWAYS AS IDENTITY(START WITH 1000) PRIMARY KEY,
    plz            text NOT NULL,
    ort            text NOT NULL,
    bundesland     text,
    fussballer_id  integer NOT NULL UNIQUE REFERENCES fussballer ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS auszeichnung (
    id             integer GENERATED ALWAYS AS IDENTITY(START WITH 1000) PRIMARY KEY,
    bezeichnung    text NOT NULL,
    saison         text NOT NULL,
    fussballer_id  integer NOT NULL REFERENCES fussballer ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS auszeichnung_fussballer_id_idx ON auszeichnung(fussballer_id);

CREATE TABLE IF NOT EXISTS fussballer_file (
    id              integer GENERATED ALWAYS AS IDENTITY(START WITH 1000) PRIMARY KEY,
    data            bytea NOT NULL,
    filename        text NOT NULL,
    mimetype        text,
    fussballer_id         integer NOT NULL UNIQUE REFERENCES fussballer ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS fussballer_file_fussballer_id_idx ON fussballer_file(fussballer_id);
