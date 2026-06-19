$composeDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$composeFile = Join-Path $composeDir "compose.yml"
$composeNoTlsFile = Join-Path $composeDir "compose.notls.yml"
$initDir = Join-Path $composeDir "init"

function Invoke-Checked {
    param(
        [Parameter(Mandatory = $true)]
        [string[]] $Command
    )

    Write-Host "> $($Command -join ' ')"
    & $Command[0] $Command[1..($Command.Length - 1)]
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed: $($Command -join ' ')"
    }
}

function Wait-PostgresHealthy {
    Write-Host "Waiting for PostgreSQL container to become healthy..."

    for ($i = 1; $i -le 30; $i++) {
        $status = docker inspect --format "{{.State.Health.Status}}" postgres 2>$null
        if ($LASTEXITCODE -eq 0 -and $status -eq "healthy") {
            Write-Host "PostgreSQL is healthy."
            return
        }

        Start-Sleep -Seconds 2
    }

    throw "PostgreSQL did not become healthy in time."
}

function Test-PgDataInitialized {
    docker run `
        -v pg_data:/var/lib/postgresql/18/data `
        --rm `
        --entrypoint= `
        dhi.io/postgres:18.3-debian13 `
        /bin/bash -c "test -s /var/lib/postgresql/18/data/PG_VERSION" *> $null

    return $LASTEXITCODE -eq 0
}

function Test-FussballerDatabaseExists {
    $result = docker compose -f $composeFile exec -T -e PGPASSWORD=p db `
        psql --dbname=postgres --username=postgres --tuples-only --no-align `
        --command "select 1 from pg_database where datname = 'fussballer';"

    if ($LASTEXITCODE -ne 0) {
        throw "Could not check whether database fussballer exists."
    }

    return ($result -join "").Trim() -eq "1"
}

function Initialize-Volumes {
    Invoke-Checked @("docker", "volume", "create", "pg_data")
    Invoke-Checked @("docker", "volume", "create", "pg_tablespace")
    Invoke-Checked @("docker", "volume", "create", "pg_init")

    $initMount = "${initDir}:/tmp/init:ro"

    Invoke-Checked @(
        "docker", "run",
        "-v", "pg_init:/init",
        "-v", "pg_tablespace:/tablespace",
        "-v", $initMount,
        "--rm",
        "-u", "0",
        "--entrypoint=",
        "dhi.io/postgres:18.3-debian13",
        "/bin/bash",
        "-c",
        "cp -r /tmp/init/* /init && mkdir -p /tablespace/fussballer && chown -R postgres:postgres /init /tablespace && chmod 400 /init/*/sql/* /init/*/csv/* /init/tls/*"
    )
}

function Initialize-TlsIfNeeded {
    if (Test-PgDataInitialized) {
        Write-Host "pg_data is already initialized. Skipping first-start TLS bootstrap."
        return
    }

    Write-Host "pg_data is empty. Starting PostgreSQL once without TLS to initialize data directory."
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "-f", $composeNoTlsFile, "up", "-d")
    Wait-PostgresHealthy

    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "db", "bash", "-c", "cp /init/tls/* /var/lib/postgresql/18/data")
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "-f", $composeNoTlsFile, "down")
}

function Start-Postgres {
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "up", "-d")
    Wait-PostgresHealthy
}

function Ensure-DatabaseExists {
    if (Test-FussballerDatabaseExists) {
        Write-Host "Database fussballer already exists."
        return
    }

    Write-Host "Creating database fussballer from project-local SQL files."

    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "-e", "PGPASSWORD=p", "db", "psql", "--dbname=postgres", "--username=postgres", "--file=/init/fussballer/sql/create-db.sql")
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "-e", "PGPASSWORD=p", "db", "psql", "--dbname=fussballer", "--username=fussballer", "--file=/init/fussballer/sql/create-schema.sql")
}

function Reset-DatabaseFromCsv {
    Write-Host "Resetting database fussballer from project-local CSV files."

    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "-e", "PGPASSWORD=p", "db", "psql", "--dbname=fussballer", "--username=fussballer", "--file=/init/fussballer/sql/drop-table.sql")
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "-e", "PGPASSWORD=p", "db", "psql", "--dbname=fussballer", "--username=fussballer", "--file=/init/fussballer/sql/create-table.sql")
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "-e", "PGPASSWORD=p", "db", "psql", "--dbname=fussballer", "--username=postgres", "--file=/init/fussballer/sql/copy-csv.sql")
}

function Show-RowCount {
    Invoke-Checked @("docker", "compose", "-f", $composeFile, "exec", "-T", "-e", "PGPASSWORD=p", "db", "psql", "--dbname=fussballer", "--username=postgres", "--command=select count(*) from fussballer.fussballer;")
}
