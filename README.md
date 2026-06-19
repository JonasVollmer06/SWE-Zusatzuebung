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
