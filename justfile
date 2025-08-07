# justfile docs  https://just.systems/man/en/
# This file is used for local development convenience
#
# Dependencies:
# - just https://github.com/casey/just?tab=readme-ov-file
# - psql https://archlinux.org/packages/?name=postgresql
# - rainfrog https://github.com/achristmascarl/rainfrog
# - CompileDaemon https://github.com/githubnemo/CompileDaemon
# (OPTIONAL) a go-test pretty printer (will otherwise fallback to `cat`) e.g.
# - tparse https://github.com/mfridman/tparse
# - gotestsum https://github.com/gotestyourself/gotestsum
# - gotestfmt https://github.com/GoTestTools/gotestfmt
# Set this here:

test_parser := "tparse"
# test_parser := "gotestfmt"
# test_parser := "gotestsum"

set quiet := true

# service

PROJECT_NAME := "gopgrest"
API_PORT := "8090"
HOST := "localhost"

# db

DB_CONTAINER := "gopgrest-db"
DB_NAME := "gopgrest"
DB_PORT := "5434"
DB_USER := "postgres"
DB_PASS := "gopgrest"
SCHEMA := "database/schema.sql"

# test db

TEST_DB_CONTAINER := "gopgrest-test-db"
TEST_DB_NAME := "gopgrest_test"
TEST_DB_PORT := "5433"
TEST_DB_USER := "gopgrest_test"
TEST_DB_PASS := "gopgrest_test"
TEST_SCHEMA := "database/test_schema.sql"

# list recipes
default:
    @just --list

###############################################################################
## dev
###############################################################################

# Build the program
[group('dev')]
build:
    go build -o ./{{ PROJECT_NAME }}

# Run the application and watch for changes, recompile/restart on changes
[group('dev')]
watch:
    #!/usr/bin/env sh
    just start
    export API_PORT={{ API_PORT }}
    export DB_PORT={{ DB_PORT }}
    export DB_NAME={{ DB_NAME }}
    export DB_USER={{ DB_USER }}
    export DB_PASS={{ DB_PASS }}
    export HOST={{ HOST }}
    CompileDaemon \
    --build="go build -o {{ PROJECT_NAME }}" \
    --command="./{{ PROJECT_NAME }}"

# Open database with rainfrog
[group('dev')]
rain:
    #!/usr/bin/env sh
    rainfrog \
        --driver="postgres" \
        --username="{{ DB_USER }}" \
        --host="localhost" \
        --port="{{ DB_PORT }}" \
        --database="{{ DB_NAME }}" \
        --password="{{ DB_PASS }}"

###############################################################################
## app
###############################################################################

# Run the app
[group('app')]
run:
    #!/usr/bin/env sh
    just init
    just start

    export API_PORT={{ API_PORT }}
    export DB_PORT={{ DB_PORT }}
    export DB_NAME={{ DB_NAME }}
    export DB_USER={{ DB_USER }}
    export DB_PASS={{ DB_PASS }}
    export HOST={{ HOST }}
    ./{{ PROJECT_NAME }}
    just stop # stop container when the program exits

# Start the container
[group('app')]
start:
    #!/usr/bin/env sh
    # Start container if not running
    if ! docker ps --format json | jq -r .Names | grep -q "^{{ DB_CONTAINER }}$"; then
        echo "Starting {{ PROJECT_NAME }} database..."
        docker start {{ DB_CONTAINER }}
        sleep 1
    else
        echo "{{ PROJECT_NAME }} database already running"
    fi

# Stop the container
[group('app')]
stop:
    #!/usr/bin/env sh
    if docker ps --format json | jq -r .Names | grep -q "^{{ DB_CONTAINER }}$"; then
        echo "Stopping {{ PROJECT_NAME }} database..."
        docker stop {{ DB_CONTAINER }}
        echo "{{ PROJECT_NAME }} database stopped"
    else
        echo "{{ PROJECT_NAME }} database not running"
    fi

###############################################################################
## db
###############################################################################

# Initialize database with schema
[group('db')]
init:
    #!/usr/bin/env sh
    # Init container if not created
    if docker ps --all --format json | jq -r .Names | grep -q "^{{ DB_CONTAINER }}$"; then
        echo "{{ PROJECT_NAME }} database already created"
        exit 0
    fi
    echo "Initializing {{ PROJECT_NAME }} database..."
    docker run -d --name {{ DB_CONTAINER }} \
        -e POSTGRES_DB={{ DB_NAME }} \
        -e POSTGRES_USER={{ DB_USER }} \
        -e POSTGRES_PASSWORD={{ DB_PASS }} \
        -p {{ DB_PORT }}:5432 \
        postgres:15 \
    && \
    sleep 2 && \
    docker cp {{ SCHEMA }} {{ DB_CONTAINER }}:/tmp/schema.sql
    docker exec {{ DB_CONTAINER }} psql -U {{ DB_USER }} -d {{ DB_NAME }} -f /tmp/schema.sql && \
    docker stop {{ DB_CONTAINER }}

# Remove the database container
[group('db')]
remove:
    docker rm {{ DB_CONTAINER }}

# Execute a psql command in the container database
[group('db')]
exec command flags="":
    docker exec -it {{ DB_CONTAINER }} psql \
        -U {{ DB_USER }} -d {{ DB_NAME }} \
        {{ flags }} \
        --command "{{ command }}"

###############################################################################
## test
###############################################################################

# Run tests
[group('test')]
test path="" parser="":
    #!/usr/bin/env sh
    parser="{{ parser }}"
    if [ -z "{{ parser }}"]; then
        parser="{{ test_parser }}"
    fi
    echo "Test parser: $parser"
    # Clean start for test db
    just tstop
    just tstart && echo "Test db started" || exit 1
    sleep 1;

    export HOST={{ HOST }}
    export TEST_DB_PORT={{ TEST_DB_PORT }}
    export TEST_DB_USER={{ TEST_DB_USER }}
    export TEST_DB_PASS={{ TEST_DB_PASS }}
    export TEST_DB_NAME={{ TEST_DB_NAME }}
    go clean -testcache | exit 1
    if [ -z "{{ path }}" ]; then
        go test -p 1 -v -cover -json ./... | $(which $parser)
    else
        go test -v -cover -json ./{{ path }} | $(which $parser)
    fi
    TEST_RESULT=$?
    just tstop
    exit $TEST_RESULT

# Start test database container
[group('test')]
tstart:
    #!/usr/bin/env sh
    if ! docker ps --format json | jq -r .Names | grep -q "^{{ TEST_DB_CONTAINER }}$"; then
        echo "Starting {{ PROJECT_NAME }} test database..."
        docker run --rm -d --name {{ TEST_DB_CONTAINER }} \
            -e POSTGRES_DB={{ TEST_DB_NAME }} \
            -e POSTGRES_USER={{ TEST_DB_USER }} \
            -e POSTGRES_PASSWORD={{ TEST_DB_PASS }} \
            -p {{ TEST_DB_PORT }}:5432 \
            postgres:15 \
        && \
        sleep 2 && \
        docker cp {{ TEST_SCHEMA }} {{ TEST_DB_CONTAINER }}:/tmp/test_schema.sql
        docker exec {{ TEST_DB_CONTAINER }} psql -U {{ TEST_DB_USER }} -d {{ TEST_DB_NAME }} -f /tmp/test_schema.sql && \
        exit 0 \
        || exit 1
    else
        echo "{{ PROJECT_NAME }} test database already running"
    fi

# Stop test database container
[group('test')]
tstop:
    #!/usr/bin/env sh
    if docker ps --format json | jq -r .Names | grep -q "^{{ TEST_DB_CONTAINER }}$"; then
        echo "Stopping {{ PROJECT_NAME }} test database..."
        docker stop {{ TEST_DB_CONTAINER }}
        echo "{{ PROJECT_NAME }} test database stopped"
    else
        echo "{{ PROJECT_NAME }} test database not running"
    fi

###############################################################################
## api calls
###############################################################################

# pick a single row by id
[group('api')]
pick table id:
    curl -X GET -s http://localhost:{{ API_PORT }}/{{ table }}/{{ id }} \
    | just jqparse

# list sets with optional query params e.g. `list authors 'surname=Plato'`
[group('api')]
list table params='':
    curl -X GET -s http://localhost:{{ API_PORT }}/{{ table }}?{{ params }} \
    | just jqparse

# Insert a row in a table e.g. `insert authors '{"surname": "Plato"}'`
[group('api')]
insert table data:
    curl -X POST -s http://localhost:{{ API_PORT }}/{{ table }} \
    --data '{{ data }}'

# Delete a row in the database by id
[group('api')]
delete table id:
    curl -X DELETE -s http://localhost:{{ API_PORT }}/{{ table }}/{{ id }}

# Update a a row in by id e.g. `update authors 1 '{"surname": "Carson"}'`
[group('api')]
update table id data:
    curl -X PUT -s http://localhost:{{ API_PORT }}/{{ table }}/{{ id }} \
    --data '{{ data }}' | \
    just jqparse

###############################################################################
## helpers
###############################################################################

# parse JSON with jq and handle invalid JSON
[group('helpers')]
jqparse:
    jq -R '. as $line | try (fromjson) catch $line'
