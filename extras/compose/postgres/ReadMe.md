# PostgreSQL fuer die Go-Zusatzuebung

Dieser Ordner enthaelt die Docker-Compose-Konfiguration, um den bestehenden
PostgreSQL-Server fuer das Fussballer-Datenmodell direkt aus diesem Projekt zu
starten.

Die Dateien orientieren sich am vorherigen Projekt `fussballer`, liegen hier aber
nochmal im aktuellen Projekt, damit die Zusatzuebung eigenstaendig gestartet
werden kann.

## Starten

Aus diesem Ordner:

```powershell
docker compose up -d
```

Oder aus dem Projektwurzelordner:

```powershell
docker compose -f .\extras\compose\postgres\compose.yml up -d
```

Status pruefen:

```powershell
docker ps
```

Erwartung:

```text
postgres ... 0.0.0.0:5432->5432/tcp ... healthy
```

## Stoppen

Aus diesem Ordner:

```powershell
docker compose down
```

Oder aus dem Projektwurzelordner:

```powershell
docker compose -f .\extras\compose\postgres\compose.yml down
```

## Verbindung

Die Go-Anwendung nutzt standardmaessig:

```text
postgres://fussballer:p@localhost:5432/fussballer?sslmode=disable
```

Bestandteile:

- Benutzer: `fussballer`
- Passwort: `p`
- Host: `localhost`
- Port: `5432`
- Datenbank: `fussballer`

## Daten pruefen

Wenn der Container laeuft:

```powershell
docker exec -e PGPASSWORD=p postgres psql -U postgres -d fussballer -c "select count(*) from fussballer.fussballer;"
```

Erwartung im aktuellen Testbestand:

```text
count
-----
7
```

## Named Volumes

Die Compose-Datei verwendet dieselben Docker Named Volumes wie das vorherige
Projekt:

```text
pg_data
pg_tablespace
pg_init
```

Falls sie auf einem Rechner noch nicht existieren, koennen sie so angelegt werden:

```powershell
docker volume create pg_data
docker volume create pg_tablespace
docker volume create pg_init
```

Der Ordner `init` enthaelt die SQL-, CSV- und TLS-Dateien, die aus dem alten
Projekt uebernommen wurden. Bei einem bereits eingerichteten Rechner sind diese
Daten normalerweise schon im Volume `pg_init` vorhanden.
