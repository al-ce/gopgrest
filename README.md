# gopgrest

A dynamic RESTful HTTP server for a Postgres database. The app will adapt its
URL routing and SQL queries to the provided table schemas.

## Limitations:

- Insert, update, or delete requests _cannot_ be made on tables without an `id`
  column
- Read requests ("pick" or "list") _can_ be made on tables without an `id`
  column
- The app will route a request with a RESTful HTTP method + path combination
  for any valid table found in the following example query:

```sql
SELECT tablename FROM Pg_catalog.pg_tables WHERE schemaname='public'"
```

```
 tablename
-----------
 authors
 books
(2 rows)
```

- Requests with JSON content (insert/update) or query params (list) must use
  valid column names and corresponding column types

## API

The following endpoints are valid for each table in the database with an `id`:

| Endpoint                     | Method | Description              | Request            | Response                                   |
| ---------------------------- | ------ | ------------------------ | ------------------ | ------------------------------------------ |
| `/{tablename}`               | POST   | Insert a new row         | `application/json` | `rows created in table {tablename}: [ids]` |
| `/{tablename}/{id}`          | GET    | Get a row by ID          | ---                | `application/json` (found row)             |
| `/{tablename}?{querystring}` | GET    | List rows matching query | ---                | `application/json` (matching rows)         |
| `/{tablename}/{id}`          | PUT    | Update a row by ID       | `application/json` | `application/json` (updated row)           |
| `/{tablename}/{id}`          | DELETE | Delete a row by ID       | ---                | `row {id} deleted from table {tablename}`  |

### Insert

Insert a new row or rows into an existing table.

```http
POST http://{HOST}:{PORT}/{tablename}
Accept: application/json
```

Example (single JSON object):

```bash
curl -X POST -s http://localhost:8090/authors \
      --data '{ "surname": "Woolf", "forename": "Virginia" }'
# rows created in table authors: [1]⏎
```

Example (multiple rows from array of JSON objects):

```bash
curl -X POST http://localhost:8090/authors \
      --data '[{ "surname": "Groton", "forename": "Anne" }, { "surname": "Plato" }]'
# rows created in table authors: [2 3]⏎
```

### Pick

Get a single row from a table by id in a JSON response

```http
GET http://{HOST}:{PORT}/{tablename}/{id}
```

Example:

```bash
curl -X GET -s http://localhost:8090/authors/1 | jq
```

```json
{
  "forename": "",
  "id": 1,
  "surname": "Plato"
}
```

### List

Get all rows from a table matching a list of optional query parameters.

```http
GET http://{HOST}:{PORT}/{tablename}?{querystring}
```

Example (no parameters):

```bash
curl -X GET -s http://localhost:8090/authors
```

```json
{
  "1": {
    "forename": "Virginia",
    "id": 1,
    "surname": "Woolf"
  },
  "2": {
    "forename": "Anne",
    "id": 2,
    "surname": "Carson"
  },
  "3": {
    "forename": "Anne",
    "id": 3,
    "surname": "Brontë"
  }
}
```

Example (multiple parameters connected by `&`):

```bash
curl -X GET -s 'http://localhost:8090/authors?forename=Anne&surname=Carson'
```

```json
{
  "3": {
    "forename": "Anne",
    "id": 3,
    "surname": "Carson"
  }
}
```

### Update

```http
PUT http://{HOST}:{PORT}/{tablename}/{id}
Accept: application/json
```

```bash
curl -X PUT -s http://localhost:8090/authors/5 --data '{"surname" : "Woolf"}'
```

```json
{
  "forename": "Virginia",
  "id": 5,
  "surname": "Woolf"
}
```

### Delete

```http
DELETE http://{HOST}:{PORT}/{tablename}/{id}
```

```bash
curl -X DELETE -s http://localhost:8090/authors/5
# example stdout: row 5 deleted from table authors
```

## Setup

Build and run the project with the following environment variables:

```bash
go build -o gopgrest
export HOST={{ HOST }}          # The host for your server, e.g. localhost
export API_PORT={{ API_PORT }}  # The port to run the server on, e.g. 8090
export DB_NAME={{ DB_NAME }}    # The name of your Postgres database
export DB_PASS={{ DB_PASS }}    # The password to your Postgres database

./gopgrest                      # Run the build output
```

### Quick setup/usage

With [casey/just](https://github.com/casey/just) as a task runner:

- Define a schema (e.g. as above) in `./database/schema.sql`
- `just run` (initialize a docker container database and run the program)
- `just insert authors '{"surname": "Woolf", "forename": "Virginia" }'`
- `just list authors 'surname=Carson`
- `just exec "select * from authors"` (query the container database directly)
- etc.

```just
Available recipes:
    default               # list recipes

    [api]
    delete table id       # Delete a row in the database by id
    insert table data     # Insert a row in a table e.g. `insert authors '{"surname": "Plato"}'`
    list table params=''  # list sets with optional query params e.g. `list authors 'surname=Plato'`
    pick table id         # pick a single row by id
    update table id data  # Update a a row in by id e.g. `update authors 1 '{"surname": "Carson"}'`

    [app]
    run                   # Run the app
    start                 # Start the container
    stop                  # Stop the container

    [db]
    exec command flags="" # Execute a psql command in the container database
    init                  # Initialize database with schema
    remove                # Remove the database container

    [dev]
    build                 # Build the program
    rain                  # Open database with rainfrog
    watch                 # Run the application and watch for changes, recompile/restart on changes

    [helpers]
    jqparse               # parse JSON with jq and handle invalid JSON

    [test]
    test path=""          # Run tests
    tstart                # Start test database container
    tstop                 # Stop test database container
```
