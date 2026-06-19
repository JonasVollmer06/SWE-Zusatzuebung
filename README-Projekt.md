# Fussballer REST API in Go

Dieses Projekt entsteht im Rahmen der Zusatzuebung "Programmierworkshop" am
19.06.2026. Ziel ist eine prototypische REST-Schnittstelle in Go gegen den
bestehenden PostgreSQL-DB-Server aus dem vorherigen Projekt `fussballer`.

Die Datei wird waehrend der Entwicklung laufend gepflegt. Sie beschreibt, welche
Befehle ausgefuehrt wurden, wie das Projekt aufgebaut ist und was der Server
kann.

## Aktueller Stand

- Projektordner ist angelegt: `swe_zusatzuebung`
- Abgabe-ReadMe des Dozenten wird separat gepflegt: `ReadMe.md`
- Technische Projekt-README wird hier gepflegt: `README-Projekt.md`
- Docker ist installiert.
- Go ist installiert: `go version go1.26.4 windows/amd64`.
- Go-Modul ist initialisiert: `module swe-zusatzuebung`.
- Bibliotheken `chi`, `pgx` und `validator` wurden installiert.
- Minimaler HTTP-Server mit `GET /health` ist implementiert.
- Erster Handler-Test fuer `GET /health` ist implementiert und erfolgreich.
- PostgreSQL aus dem alten Projekt wurde per Docker Compose gestartet.
- Die Datenbank `fussballer` ist erreichbar und enthaelt aktuell 7 Datensaetze in
  `fussballer.fussballer`.
- Der Go-Server baut beim Start eine PostgreSQL-Verbindung auf und prueft sie mit
  einem Ping.
- `.gitignore` wurde fuer Go, lokale Umgebungsdateien, Logs und Editor-Dateien
  angelegt.
- Erste Models fuer `Fussballer`, `Adresse`, `Auszeichnung`, `Position` und
  `CreateFussballerRequest` wurden angelegt.
- PostgreSQL kann jetzt direkt aus diesem Projekt gestartet werden:
  `extras/compose/postgres/compose.yml`.
- Beim Serverstart wird ein Banner fuer die Fussballer REST API ausgegeben.
- Bestehendes Datenmodell wurde aus dem Projekt `fussballer` analysiert.

## Voraussetzungen

### Go installieren

Go wird als Plattform fuer die Zusatzuebung verwendet.

Download:

```text
https://go.dev/dl/
```

Empfohlener Installer fuer Windows 10+ 64-bit:

```text
go1.26.4.windows-amd64.msi
```

Nach der Installation PowerShell neu oeffnen und pruefen:

```powershell
go version
```

Erwartetes Ergebnis:

```text
go version go1.26.4 windows/amd64
```

### Docker

Docker wird fuer den bestehenden PostgreSQL-Server aus dem alten Projekt genutzt.

Geprueft mit:

```powershell
docker --version
```

Ergebnis:

```text
Docker version 29.5.2, build 79eb04c
```

## Verwendete Bibliotheken

Geplant sind diese Go-Bibliotheken:

```text
github.com/go-chi/chi/v5
github.com/jackc/pgx/v5
github.com/go-playground/validator/v10
```

### Warum diese Bibliotheken?

- `chi`: schlanker HTTP-Router fuer REST-Endpunkte.
- `pgx`: PostgreSQL-Treiber fuer direkten, gut nachvollziehbaren DB-Zugriff.
- `validator`: Validierung von JSON-Requests beim Neuanlegen.

Keycloak/OIDC wird vorerst nicht eingebaut, weil es laut Aufgabenstellung
optional ist.

## Bisher ausgefuehrte Befehle

Im Projekt `swe_zusatzuebung`:

```powershell
Get-ChildItem -Force
Get-Content -Raw -LiteralPath .\ReadMe.md
go version
docker --version
go mod init swe-zusatzuebung
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/go-playground/validator/v10
go get github.com/jackc/pgx/v5/pgxpool@v5.10.0
go test ./...
```

Hinweis: In der Codex-Sitzung war der frisch gesetzte Go-PATH noch nicht sichtbar.
Deshalb wurden `gofmt` und `go test` dort einmal mit absolutem Pfad ausgefuehrt:

```powershell
& 'C:\Program Files\Go\bin\gofmt.exe' -w .\cmd\server\main.go .\internal\config\config.go .\internal\server\router.go .\internal\server\router_test.go
& 'C:\Program Files\Go\bin\go.exe' test ./...
```

Testergebnis:

```text
?   	swe-zusatzuebung/cmd/server	[no test files]
?   	swe-zusatzuebung/internal/config	[no test files]
ok  	swe-zusatzuebung/internal/server
```

PostgreSQL-Container aus diesem Projekt starten:

```powershell
docker compose up -d
```

Ausgefuehrt im Ordner:

```text
C:\Users\jv10s\SWE\Projekte\swe_zusatzuebung\extras\compose\postgres
```

Alternativ aus dem Projektwurzelordner:

```powershell
docker compose -f .\extras\compose\postgres\compose.yml up -d
```

Docker-Status:

```text
postgres   dhi.io/postgres:18.3-debian13   0.0.0.0:5432->5432/tcp   Up ... (healthy)
```

DB-Pruefung im Container:

```powershell
docker exec -e PGPASSWORD=p postgres psql -U postgres -d fussballer -c "select count(*) from fussballer.fussballer;"
```

Ergebnis:

```text
count
-----
7
```

Im alten Projekt `fussballer` wurden Dateien gelesen, um das Datenmodell zu
verstehen:

```powershell
Get-Content -Raw -LiteralPath C:\Users\jv10s\SWE\Projekte\fussballer\prisma\schema.prisma
Get-Content -Raw -LiteralPath C:\Users\jv10s\SWE\Projekte\fussballer\extras\compose\postgres\compose.yml
Get-Content -Raw -LiteralPath C:\Users\jv10s\SWE\Projekte\fussballer\extras\compose\postgres\init\fussballer\sql\create-table.sql
Get-Content -Raw -LiteralPath C:\Users\jv10s\SWE\Projekte\fussballer\.env
```

Die PostgreSQL-Compose-Umgebung wurde aus dem alten Projekt in dieses Projekt
uebernommen:

```text
extras/compose/postgres/compose.yml
extras/compose/postgres/password.txt
extras/compose/postgres/init/
extras/compose/postgres/ReadMe.md
```

Zusaetzlicher Vorlesungsmitschrieb wurde gelesen:

```powershell
Get-Content -Raw -LiteralPath C:\Users\jv10s\Desktop\SWE\erklaerung.txt
```

Kernaussagen daraus:

- REST-Schnittstelle mit Datenbankzugriff ist gefordert.
- Der DB-Server aus den vorherigen Abgaben soll verwendet werden.
- Git-Repository soll eingerichtet und genutzt werden.
- Docker soll vernuenftig laufen.
- Die Abgabe-ReadMe soll KI-Werkzeug und Requests dokumentieren.

## Setup-Befehle

Ausgefuehrt:

```powershell
go mod init swe-zusatzuebung
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/go-playground/validator/v10
go get github.com/jackc/pgx/v5/pgxpool@v5.10.0
```

Server starten:

```powershell
go run ./cmd/server
```

Tests ausfuehren:

```powershell
go test ./...
```

Datenbank starten:

```powershell
docker compose -f .\extras\compose\postgres\compose.yml up -d
```

Datenbank stoppen:

```powershell
docker compose -f .\extras\compose\postgres\compose.yml down
```

## Git-Arbeitsweise

Wir arbeiten nicht dauerhaft direkt auf `main`, sondern nutzen Feature-Branches
fuer zusammenhaengende Schritte.

Aktueller Feature-Branch:

```text
codex-fussballer-read-api
```

Grundregel:

- Kleine, sinnvolle Commits nach stabilen Zwischenschritten.
- Vor einem Commit sollten Tests laufen: `go test ./...`.
- Push nach GitHub kann manuell erfolgen.

Nuetzliche Git-Befehle:

```powershell
git status
git add .
git commit -m "Kurze Beschreibung"
git push -u origin codex-fussballer-read-api
```

## Geplante Projektstruktur

```text
swe_zusatzuebung/
  cmd/
    server/
      main.go
  extras/
    compose/
      postgres/
        compose.yml
        password.txt
        ReadMe.md
        init/
  internal/
    config/
      config.go
    database/
      database.go
    server/
      banner.go
      router.go
      router_test.go
    fussballer/
      model.go
      model_test.go
      repository.go
      service.go
      router.go
      validation.go
  ReadMe.md
  README-Projekt.md
  go.mod
  go.sum
```

### Bedeutung der Ordner

- `cmd/server`: Einstiegspunkt der Anwendung. Hier startet der HTTP-Server.
- `extras/compose/postgres`: Docker-Compose-Setup fuer PostgreSQL mit Init-Dateien.
- `internal/config`: Konfiguration, z.B. Port und Datenbank-URL.
- `internal/database`: Aufbau und Pruefung der PostgreSQL-Verbindung.
- `internal/server`: Allgemeiner HTTP-Router, Health Check und Startbanner.
- `internal/fussballer`: Fachlogik fuer Fussballer.
- `repository.go`: Datenbankzugriff.
- `service.go`: Geschaeftslogik zwischen Router und Repository.
- `router.go`: REST-Routen und HTTP-Handler.
- `model.go`: Datenstrukturen fuer Fussballer, Adresse und Auszeichnungen.
- `model_test.go`: Erste Tests fuer fachliche Konstanten, aktuell Positionswerte.
- `validation.go`: Regeln fuer neue Fussballer.

## Datenbankgrundlage

Das bestehende Projekt `fussballer` verwendet PostgreSQL mit dem Schema
`fussballer`.

Wichtige Tabellen:

- `fussballer`
- `adresse`
- `auszeichnung`
- `fussballer_file`

Wichtige Felder der Tabelle `fussballer`:

- `id`
- `version`
- `nachname`
- `nationalitaet`
- `position`
- `geburtsdatum`
- `username`
- `erzeugt`
- `aktualisiert`

Aktuell angelegte Go-Models:

- `Fussballer`
- `Adresse`
- `Auszeichnung`
- `Position`
- `CreateFussballerRequest`
- `CreateAdresseRequest`

Erlaubte Werte fuer `position`:

- `TORWART`
- `VERTEIDIGER`
- `MITTELFELDSPIELER`
- `STUERMER`

## Geplante REST-Endpunkte

### Health Check

```http
GET /health
```

Zweck: Pruefen, ob der Server laeuft.

Aktuelle Antwort:

```json
{"status":"ok"}
```

### Fussballer nach ID lesen

```http
GET /fussballer/{id}
```

Zweck: Einen einzelnen Fussballer anhand der ID lesen.

### Fussballer suchen

```http
GET /fussballer
```

Moegliche Query-Parameter:

```text
nachname
nationalitaet
position
```

Beispiel:

```http
GET /fussballer?position=TORWART
```

### Fussballer neu anlegen

```http
POST /fussballer
Content-Type: application/json
```

Zweck: Einen neuen Fussballer anlegen.

Validierung:

- `nachname` ist Pflichtfeld.
- `nationalitaet` ist Pflichtfeld.
- `geburtsdatum` ist Pflichtfeld.
- `username` ist Pflichtfeld.
- `position` muss einer der erlaubten Enum-Werte sein.

## Was der Server koennen soll

Aktuell kann der Server:

- per HTTP starten,
- beim Start ein Banner ausgeben,
- beim Start eine PostgreSQL-Verbindung aufbauen und pruefen,
- einen Health-Endpunkt anbieten.

Am Ende soll der Server zusaetzlich:

- Fussballer aus PostgreSQL lesen,
- Fussballer ueber Query-Parameter suchen,
- neue Fussballer per JSON anlegen,
- Eingaben beim Neuanlegen validieren,
- Fehler als sinnvolle HTTP-Statuscodes zurueckgeben,
- einfache Tests fuer Router/Handler enthalten.

## Teststrategie

Geplant:

- Unit-/Handler-Tests mit Go und `net/http/httptest`.
- Repository-Tests optional gegen laufende PostgreSQL-Datenbank.

Spaetere Testbefehle:

```powershell
go test ./...
```

## Offene Punkte

- DB-Verbindung ist grundlegend konfiguriert.
- PostgreSQL-Compose-Setup ist im aktuellen Projekt vorhanden und getestet.
- Repository fuer `GET /fussballer/{id}` implementieren.
- Service fuer `GET /fussballer/{id}` implementieren.
- Router fuer `GET /fussballer/{id}` implementieren.
- Tests ergaenzen.
