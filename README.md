# Programmierworkshop am 19.6.2026

## Namen

TODO: Namen eintragen.

## Link zum Git-Repository

TODO: Repository-Link eintragen, sobald vorhanden.

## KI-Werkzeuge

- Codex / ChatGPT als KI-Agent zur schrittweisen Projektunterstuetzung, Code-Erzeugung,
  Erklaerung fuer Go-Einsteiger und laufenden Pflege dieser ReadMe.

### Agenten

- Codex im lokalen Projektordner `swe_zusatzuebung`.

### Chat-URLs, z.B. https://chatgpt.com

- TODO: Chat-URL eintragen, falls sie fuer die Abgabe benoetigt wird.

## Frameworks und Bibliotheken

### REST-Schnittstelle (Lesen und Neuanlegen)

- Geplant: Go Standardbibliothek `net/http` plus Router `github.com/go-chi/chi/v5`.
- Begründung: `chi` ist schlank, idiomatisch fuer Go und reicht fuer eine kleine REST-API
  mit Routen wie `GET /fussballer`, `GET /fussballer/{id}` und `POST /fussballer`.

### Validierung (nur Neuanlegen)

- Geplant: `github.com/go-playground/validator/v10`.
- Begründung: Damit koennen Pflichtfelder und einfache Regeln am Request-DTO beschrieben
  und beim `POST` geprueft werden.

### OR-Mapping (für PostgreSQL)

- Geplant: PostgreSQL-Anbindung mit `github.com/jackc/pgx/v5`.
- Hinweis: Falls ein echtes OR-Mapping zwingend gefordert wird, wird alternativ oder
  zusaetzlich `gorm.io/gorm` mit PostgreSQL-Treiber bewertet. Fuer den Prototyp ist
  `pgx` einfacher, transparenter und gut testbar.
- Bestehende Datenbasis aus dem Projekt `fussballer`:
  Schema `fussballer`, Haupttabelle `fussballer`, zugehoerige Tabellen `adresse`,
  `auszeichnung` und `fussballer_file`.

### Optional: OIDC mit Keycloak

- Vorerst nicht geplant, da Keycloak laut Aufgabenstellung optional ist und die Zeit
  fuer REST, DB-Zugriff, Validierung und Tests priorisiert wird.

### Einfacher Integrationstest

- Geplant: Go-Tests mit `net/http/httptest` fuer Handler/Router.
- Optional spaeter: Integrationstest gegen eine laufende PostgreSQL-Testdatenbank.

## Prompts/Requests an KI-Agent/en

- Ausgangsprompt: Neues Go-Projekt fuer eine Zusatzuebung. In 4,5 Stunden soll eine
  prototypische REST-Schnittstelle gegen den bestehenden PostgreSQL-DB-Server aus dem
  Projekt `fussballer` entstehen. Gewuenscht sind Router, Service, Repository, Tests,
  Query-Anfragen und eine laufend gepflegte ReadMe. Der Nutzer ist Go-Anfaenger und
  moechte schrittweise begleitet werden.
- KI-Antwort/Entscheidung: Zuerst Projekt einrichten. Go war lokal noch nicht im PATH
  verfuegbar. Empfehlung: offiziellen Windows-Installer von go.dev installieren.
  Danach Go-Modul initialisieren und die Bibliotheken `chi`, `pgx` und `validator`
  einbinden.
- Zusatzrequest: Neben der Abgabe-ReadMe soll eine zweite technische README gepflegt
  werden, die wie eine normale Projekt-README aufgebaut ist. Sie dokumentiert Setup,
  ausgefuehrte Befehle, Projektstruktur, Server-Funktionen, Endpunkte und Tests.
- KI-Antwort/Entscheidung: Anlage der Datei `README-Projekt.md` als technische
  Projektdokumentation, waehrend `ReadMe.md` weiterhin die ILIAS-Abgabevorlage bleibt.
- Nutzer-Update: Go wurde installiert (`go version go1.26.4 windows/amd64`), das
  Modul wurde mit `go mod init swe-zusatzuebung` initialisiert und die Bibliotheken
  `chi`, `pgx` und `validator` wurden per `go get` installiert.
- KI-Antwort/Entscheidung: Naechster sorgfaeltiger Setup-Schritt ist ein minimaler
  HTTP-Server mit `GET /health`, getrennt in `cmd/server`, `internal/config` und
  `internal/server`, inklusive erstem Handler-Test mit `net/http/httptest`.
- Ergebnis: `go test ./...` ist erfolgreich. Der Health-Endpunkt liefert JSON:
  `{"status":"ok"}`.
- Nutzer-Update: Der manuelle Start mit `go run ./cmd/server` und der Aufruf von
  `/health` haben funktioniert.
- KI-Antwort/Entscheidung: Naechster Setup-Schritt ist die PostgreSQL-Anbindung.
  Der bestehende Docker-PostgreSQL-Server aus dem Projekt `fussballer` wurde
  gestartet und als `healthy` geprueft.
- Ergebnis: Die Datenbank `fussballer` ist erreichbar; die Tabelle
  `fussballer.fussballer` enthaelt 7 Datensaetze. Der Go-Server nutzt nun
  `pgxpool`, liest `DATABASE_URL` aus der Umgebung oder verwendet lokal
  `postgres://fussballer:p@localhost:5432/fussballer?sslmode=disable` und prueft
  die Verbindung beim Start mit einem Ping. `go test ./...` ist erfolgreich.
- Zusatzrequest: Git soll als fester Teil des Projekts genutzt werden:
  passende Branches, kleine sinnvolle Commits und Push optional manuell durch den
  Nutzer.
- KI-Antwort/Entscheidung: Fuer den naechsten fachlichen Schritt wurde der
  Feature-Branch `codex-fussballer-read-api` angelegt. Die Git-Arbeitsweise wurde
  in `README-Projekt.md` dokumentiert.
- Zusatzinput: Vorlesungsmitschrieb `erklaerung.txt` bestaetigt REST-Schnittstelle
  mit DB-Zugriff, Nutzung des bestehenden DB-Servers, Git-Repository, laufenden
  Docker-Betrieb und Dokumentation der KI-Requests in der Abgabe-ReadMe.
- Nutzerrequest: Sinnvolle `.gitignore` hinzufuegen und langsam mit dem Serveraufbau
  weitermachen; Frage, ob zuerst die Models sinnvoll sind.
- KI-Antwort/Entscheidung: Ja, zuerst schlanke Models, weil sie die gemeinsame
  Datenstruktur fuer Repository, Service und Router bilden. Angelegt wurden
  `.gitignore`, `internal/fussballer/model.go` und ein erster Test
  `internal/fussballer/model_test.go` fuer gueltige Positionswerte.
- Nutzerrequest: PostgreSQL soll nicht nur ueber das alte Projekt gestartet werden,
  sondern auch direkt aus diesem Zusatzuebungsprojekt. Orientierung an der alten
  PostgreSQL-ReadMe.
- KI-Antwort/Entscheidung: Die PostgreSQL-Compose-Umgebung wurde nach
  `extras/compose/postgres` uebernommen, inklusive `compose.yml`, `password.txt`,
  `init`-Dateien und eigener kurzer Start-ReadMe.
- Ergebnis: Der vorher aus dem alten Projekt gestartete Container wurde gestoppt und
  PostgreSQL erfolgreich aus `swe_zusatzuebung\extras\compose\postgres` gestartet.
  Der Container ist `healthy`; die Abfrage auf `fussballer.fussballer` liefert
  weiterhin 7 Datensaetze.
- Nutzerrequest: Beim Serverstart soll ein Banner ausgegeben werden.
- KI-Antwort/Entscheidung: Das Banner wurde als eigene Funktion
  `server.PrintBanner()` in `internal/server/banner.go` umgesetzt und in
  `cmd/server/main.go` beim Start aufgerufen.
- Ergebnis: `go test ./...` ist erfolgreich; ein kurzer Smoke-Test von `/health`
  war ebenfalls erfolgreich.
- Nutzerrequest: PostgreSQL-Setup soll per Script fuer dieses Projekt neu aufgebaut
  werden, aber mit den alten Volume-Namen `pg_data`, `pg_tablespace` und `pg_init`.
- KI-Antwort/Entscheidung: `extras/compose/postgres/setup.ps1` wurde angelegt.
  Das Script erstellt/prueft die alten Docker-Volumes, kopiert die Init-Dateien aus
  diesem Projekt in `pg_init`, bereitet den Tablespace vor, startet bei leerem
  `pg_data` einmal ohne TLS und danach normal mit TLS. Die SQL-Initialisierung wird
  nur ausgefuehrt, wenn die Datenbank `fussballer` noch nicht existiert.
- Ergebnis: `setup.ps1` wurde erfolgreich ausgefuehrt. Es erkannte die vorhandene
  Datenbank, uebersprang die SQL-Neuanlage und pruefte `select count(*) from
  fussballer.fussballer;` mit Ergebnis 7. `go test ./...` ist erfolgreich.
- Nutzerrequest: Linter und Formatierungstool einbauen, pruefen und ausfuehren.
- KI-Antwort/Entscheidung: Fuer den ersten Schritt werden Go-Bordmittel verwendet:
  `gofmt` als Formatierer und `go vet` als Linter/statische Pruefung. Dafuer wurden
  die Skripte `scripts/format.ps1`, `scripts/lint.ps1`, `scripts/check.ps1` und
  `scripts/go-tools.ps1` angelegt.
- Ergebnis: PowerShell-Syntaxcheck, Formatierung, Linting und Gesamtcheck mit Tests
  wurden erfolgreich ausgefuehrt.
- Nutzerrequest: GitHub Actions einbauen, sodass bei Pushes Tests, Linter und
  Formatierer ausgefuehrt werden.
- KI-Antwort/Entscheidung: Ein einfacher Workflow `.github/workflows/ci.yml` wurde
  angelegt. Er nutzt `actions/checkout`, `actions/setup-go` mit `go.mod`,
  `go mod download`, einen `gofmt`-Formatcheck, `go vet ./...` und `go test ./...`.
  Lokal wurde zusaetzlich `scripts/format-check.ps1` ergaenzt und in
  `scripts/check.ps1` eingebunden.
- Ergebnis: PowerShell-Syntaxcheck und `scripts/check.ps1` wurden lokal erfolgreich
  ausgefuehrt.
- Nutzerrequest: Aktuellen Stand holen und ein Repository fuer den Lesezugriff
  erstellen, orientiert am Hono-Projekt, aber passend fuer Go.
- KI-Antwort/Entscheidung: Der Stand war bereits aktuell. Das Hono-Projekt wurde
  fuer `findById`, Query-Parameter-Suche und Count als fachliche Vorlage gelesen.
  In Go wurde `internal/fussballer/repository.go` mit `pgxpool`, SQL-Queries,
  Fehlern fuer `not found` und ungueltige Suchparameter sowie Suchkriterien fuer
  `nachname`, `nationalitaet` und `position` umgesetzt.
- Ergebnis: `go test ./...` ist erfolgreich; Unit-Tests fuer den dynamischen
  WHERE-Klausel-Aufbau wurden ergaenzt.
- Nutzerrequest: Nach gleichem Vorgehen soll nun der Router fuer den Lesezugriff
  entstehen, wieder orientiert am Hono-Projekt.
- KI-Antwort/Entscheidung: Da im aktuellen Git-Stand noch kein Read-Service
  vorhanden war, wurde zuerst `internal/fussballer/service.go` ergaenzt und
  danach `internal/fussballer/router.go` umgesetzt. Der Router bietet
  `GET /fussballer/{id}` mit ETag und `If-None-Match`, `GET /fussballer` mit
  Query-Parametern `nachname`, `nationalitaet`, `position`, Pagination ueber
  `page` und `size` sowie `count-only`.
- Ergebnis: `cmd/server/main.go` verdrahtet nun Repository, Read-Service und
  Fussballer-Router. Router- und Service-Tests wurden ergaenzt; `go test ./...`
- Nutzerrequest: Orientierung am alten Hono-Projekt und Umsetzung eines
  Read-Service fuer den Lesezugriff, inklusive Tests falls sinnvoll.
- KI-Antwort/Entscheidung: Der Read-Service wurde in Go als Schicht ueber dem
  Repository angelegt. Er validiert IDs, Suchparameter und Pagination, ruft
  Repository-Methoden fuer `findById`, Suche und Count auf und liefert einen
  Slice mit `content` und `totalElements`.
- Ergebnis: Service-Tests mit Fake-Repository wurden ergaenzt; `go test ./...`
  ist erfolgreich.
- Nutzerrequest: Aktuellen Stand von `main` holen und danach schrittweise die
  Write-Komponenten bauen. Zuerst soll nur das Repository fuer Schreibzugriff
  ergaenzt werden. Integrationstests und Bruno-Requests sollen auf die TODO-Liste.
- KI-Antwort/Entscheidung: `main` wurde per Fast-Forward auf `origin/main`
  aktualisiert und ein neuer Branch `codex-write-repository` erstellt.
  `Repository.Create(...)` wurde in `internal/fussballer/repository.go`
  implementiert. Die Methode arbeitet mit einer Transaktion, fuegt zuerst einen
  Datensatz in `fussballer.fussballer` ein und legt optional eine Adresse in
  `fussballer.adresse` an.
- Ergebnis: `scripts/check.ps1` ist erfolgreich. TODOs fuer Write-Service,
  Write-Router, Integrationstests und Bruno-Collection wurden in
  `README-Projekt.md` ergaenzt.
- Nutzerrequest: Write-Service mit Validierung und Unit-Tests bauen, orientiert am
  Read-Service. Fuer diesen Schritt soll ein neuer Branch ohne KI-Hinweis im Namen
  verwendet werden.
- KI-Antwort/Entscheidung: Neuer Branch `write-service-validation` wurde erstellt.
  `internal/fussballer/write_service.go` implementiert `WriteService.Create(...)`
  ueber einem `WriteRepository`. Fuer Pflichtfelder wird `validator` genutzt,
  zusaetzlich wird die `Position` fachlich geprueft und String-Felder werden per
  `strings.TrimSpace` normalisiert.
- Ergebnis: Unit-Tests in `internal/fussballer/write_service_test.go` pruefen
  erfolgreichen Create-Aufruf, fehlende Pflichtfelder, ungueltige Position,
  ungueltige Adresse und Weitergabe von Repository-Fehlern. `scripts/check.ps1`
  ist erfolgreich.
