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
# test_parser := ```
#     which tparse || which gotestsum || which gotestfmt || which cat
# ```

set quiet := true

# service

PROJECT_NAME := "ftrack"
API_PORT := "8090"

# db

DB_PASS := "ftrack"
DB_NAME := "ftrack"
DB_PORT := "5432"
DB_USER := "postgres"
INIT_DB := "database/init_db.sql"
SCHEMA := "database/schema.sql"

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
    sudo -u postgres psql -d {{ DB_NAME}} -f {{ SCHEMA }}

# Drop `sets` table from database
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
## api calls
###############################################################################


[group('api')]
list:
    curl -X GET -s http://localhost:{{ API_PORT }}/ | jq
