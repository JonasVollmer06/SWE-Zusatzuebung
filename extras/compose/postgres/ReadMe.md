# PostgreSQL fuer die Go-Zusatzuebung

Dieser Ordner enthaelt die Docker-Compose-Konfiguration, um den bestehenden
PostgreSQL-Server fuer das Fussballer-Datenmodell direkt aus diesem Projekt zu
starten.

Die Dateien orientieren sich am vorherigen Projekt `fussballer`, liegen hier aber
nochmal im aktuellen Projekt, damit die Zusatzuebung eigenstaendig gestartet
werden kann.

## Starten

### Einmaliges Erstsetup

Das einmalige Setup fuer dieses Projekt ist:

```powershell
.\setup.ps1
```

Das Script orientiert sich an der PostgreSQL-ReadMe aus dem vorherigen Projekt
und verwendet bewusst dieselben alten Docker-Volume-Namen:

```text
pg_data
pg_tablespace
pg_init
```

Es fuehrt aus:

- Docker Volumes anlegen, falls sie noch fehlen.
- Projektlokalen Ordner `init` in das Volume `pg_init` kopieren.
- Tablespace-Ordner `/tablespace/fussballer` vorbereiten.
- Bei leerem `pg_data` PostgreSQL einmal ohne TLS starten, damit das Datenverzeichnis
  initialisiert wird.
- TLS-Zertifikate nach `/var/lib/postgresql/18/data` kopieren.
- PostgreSQL normal mit TLS starten.
- Datenbank `fussballer` anlegen, falls sie noch nicht existiert.

Wenn die Datenbank bereits existiert, wird sie nicht neu angelegt.

### Normaler Start

Empfohlen:

```powershell
.\start.ps1
```

Das startet PostgreSQL und prueft kurz den Datenbestand.

Aus diesem Ordner:

```powershell
docker compose up -d
```

Oder aus dem Projektwurzelordner:

```powershell
docker compose -f .\extras\compose\postgres\compose.yml up -d
```

### Datenbank auf CSV-Stand zuruecksetzen

Wenn Requests die Testdaten veraendert haben, kann die Datenbank wieder aus den
CSV-Dateien geladen werden:

```powershell
.\reset-db.ps1
```

Das Script startet PostgreSQL, legt die Datenbank bei Bedarf an und setzt die
Tabellen mit `drop-table.sql`, `create-table.sql` und `copy-csv.sql` auf den
Stand der CSV-Dateien zurueck.

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

Positionen im CSV-Ausgangsbestand:

```powershell
docker exec -e PGPASSWORD=p postgres psql -U postgres -d fussballer -c "select id, nachname, position from fussballer.fussballer order by id;"
```

Dabei enthaelt der CSV-Stand u.a. `id=30` mit `position=TORWART`.

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

Das uebernimmt normalerweise `setup.ps1`.

Der Ordner `init` enthaelt die SQL-, CSV- und TLS-Dateien, die aus dem alten
Projekt uebernommen wurden. `setup.ps1` kopiert diese projektlokalen Dateien in
das Volume `pg_init`.
