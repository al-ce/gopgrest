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

test_parser := ```
    which gotestfmt || which tparse || which gotestsum || which cat
```

set quiet := true

# service

PROJECT_NAME := "ftrack"
API_PORT := "8090"
HOST := "localhost"

# db

DB_NAME := "ftrack"
DB_PORT := "5432"
DB_USER := "postgres"
DB_PASS := "ftrack"
INIT_DB := "database/init_db.sql"
SCHEMA := "database/schema.sql"

# test db

TEST_DB_NAME := "ftrack_test"
TEST_DB_PORT := "5433"
TEST_DB_USER := "ftrack_test"
TEST_DB_PASS := "ftrack_test"
TEST_DB_CONTAINER := "ftrack-test-db"
TEST_SCHEMA := "database/test_schema.sql"

default:
    @just --list

###############################################################################
## dev
###############################################################################

# Run the application and watch for changes, recompile/restart on changes
[group('dev')]
watch:
    #!/usr/bin/env sh
    export API_PORT={{ API_PORT }}
    export DB_NAME={{ DB_NAME }}
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
## db
###############################################################################

# Initialize database with schema
[group('db')]
init:
    #!/usr/bin/env sh
    sudo -u postgres psql -f {{ INIT_DB }}
    sudo -u postgres psql -d {{ DB_NAME }} -f {{ SCHEMA }}

# Drop db
[group('db')]
drop:
    #!/usr/bin/env sh
    sudo -u postgres psql -c "DROP DATABASE IF EXISTS {{ DB_NAME }};"

# Execute a psql command in the database
[group('db')]
exec command flags="":
    sudo -u postgres psql -U {{ DB_USER }} -d {{ DB_NAME }} \
        {{ flags }} --command "{{ command }}"

###############################################################################
## test
###############################################################################

# Run tests
[group('test')]
test path="":
    #!/usr/bin/env sh
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
        go test -p 1 -v -cover -json ./... | {{ test_parser }}
    else
        go test -v -cover -json ./{{ path }} | {{ test_parser }}
    fi
    TEST_RESULT=$?
    just tstop
    exit $TEST_RESULT

# Start test database container
[group('test')]
tstart:
    #!/usr/bin/env sh
    if ! docker ps --format json | jq -r .Names | grep -q "^{{ TEST_DB_CONTAINER }}$"; then
        echo "Starting test database..."
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
        echo "Test database already running"
    fi

# Stop test database container
[group('test')]
tstop:
    #!/usr/bin/env sh
    if docker ps --format json | jq -r .Names | grep -q "^{{ TEST_DB_CONTAINER }}$"; then
        echo "Stopping test database..."
        docker stop {{ TEST_DB_CONTAINER }}
        echo "Test database stopped"
    else
        echo "Test database not running"
    fi

###############################################################################
## api calls
###############################################################################

# pick gets a row by id
[group('api')]
pick table id:
    curl -X GET -s http://localhost:{{ API_PORT }}/{{ table }}/{{ id }} \
    | just jqparse

# list sets filtered by optional query params
[group('api')]
list table params='':
    curl -X GET -s http://localhost:{{ API_PORT }}/{{ table }}?{{ params }} \
    | just jqparse

# insert a row in the specified table
[group('api')]
insert table data:
    curl -X POST -s http://localhost:{{ API_PORT }}/{{ table }} \
    --data '{{ data }}'

# delete a set in the database by id
[group('api')]
delete table id:
    curl -X DELETE -s http://localhost:{{ API_PORT }}/{{ table }}/{{ id }}

# delete a set in the database by id
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

