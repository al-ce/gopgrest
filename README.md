# gopgrest

A dynamic RESTful HTTP server for a Postgres database. The app will adapt its
URL routing and SQL queries to the provided table schemas. Supports an RSQL
based query language.

## API

The following endpoints are valid for each table in the database with an `id`:

| Endpoint                     | Method | Description                   | Request            | Response                                   |
| ---------------------------- | ------ | ----------------------------- | ------------------ | ------------------------------------------ |
| `/`                          | GET    | Get structure of all tables   | ---                | `application/json` (tables)                |
| `/{tablename}`               | POST   | Insert a new row              | `application/json` | `rows created in table {tablename}: [ids]` |
| `/{tablename}/{id}`          | GET    | Get a row by ID               | ---                | `application/json` (found row)             |
| `/{tablename}?{querystring}` | GET    | List rows matching RSQL query | ---                | `application/json` (matching rows)         |
| `/{tablename}/{id}`          | PUT    | Update a row by ID            | `application/json` | `application/json` (updated row)           |
| `/{tablename}/{id}`          | DELETE | Delete a row by ID            | ---                | `row {id} deleted from table {tablename}`  |

## REST Query language (based on RSQL)

The project has a goal to be compatible with the [rsql-jpa-specification](https://github.com/perplexhub/rsql-jpa-specification) with additional features.

Query parameters can be added to a GET request after a `?` query separator. Keys and their values are separated by `=`. Multiple subqueries can be joined with `&`.

The following example contains:

- a `fields` subquery to select columns to return, with qualifiers and aliases on some fields
- a `left_join` subquery to add two join relations
- a `filter` subquery to add conditional filters to the query

```bash
curl -X GET -s 'http://localhost:8090/books?fields=title,genres.name:genre,authors.surname:author&left_join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id&filter=born<1900' | jq
```

```json
[
  {
    "author": "Brontë",
    "genre": "Epistolary",
    "title": "The Tenant of Wildfell Hall"
  },
  {
    "author": "Woolf",
    "genre": "Modernism",
    "title": "To The Lighthouse"
  },
  {
    "author": "Woolf",
    "genre": null,
    "title": "Mrs. Dalloway"
  }
]
```

The JSON output from the request is marshalled from the rows returned by this SQL query:

```sql
SELECT title,name AS genre, surname AS author FROM books
LEFT JOIN authors ON books.author_id = authors.id
LEFT JOIN genres ON books.genre_id = genres.id
WHERE authors.born < ($1)
```

The following query keys are supported:

| Key              | Description                                |
| ---------------- | ------------------------------------------ |
| `filter`         | add `WHERE` conditions to a `SELECT` query |
| `fields`         | columns to return in a `SELECT` query      |
| `join` (various) | add `JOIN` relations to a `SELECT` query   |

### Filter

A `filter` key can be added to the URL query to match a SQL `WHERE` clause.

A filter subquery is in the following format:

```
filter=[{column_name}{operator}{value};...]
```

where the right hand side is a `;` separated list of conditional expressions equivalent to a `WHERE` clause.

For example, the following SQL query and GET request are equivalent:

```bash
curl -X GET -s 'http://localhost:8090/authors?filter=forename==Anne;born>=1900' | jq
```

```json
{
  "1": {
    "born": 1950,
    "died": null,
    "forename": "Anne",
    "id": 1,
    "surname": "Carson"
  }
}
```

```sql
SELECT * FROM author
WHERE forename = 'Anne' AND born >= 1900
```

```
 id | surname | forename | born | died
----+---------+----------+------+------
  1 | Carson  | Anne     | 1950 |
(1 row)

```

These are the currently supported filter operators:

| Operator      | SQL equivalent |
| ------------- | -------------- |
| `==`          | `=`            |
| `!=`          | `!=`           |
| `=in=`        | `IN`           |
| `=out=`       | `NOT IN`       |
| `=like=`      | `LIKE`         |
| `=!like=`     | `NOT LIKE`     |
| `=notlike=`   | `NOT LIKE`     |
| `=nk=`        | `NOT LIKE`     |
| `=isnull=`    | `IS NULL`      |
| `=na=`        | `IS NULL`      |
| `=isnotnull=` | `IS NOT NULL`  |
| `=notnull=`   | `IS NOT NULL`  |
| `=nn=`        | `IS NOT NULL`  |
| `=!null=`     | `IS NOT NULL`  |
| `=le=`        | `<=`           |
| `<=`          | `<=`           |
| `=ge=`        | `>=`           |
| `>=`          | `>=`           |
| `=lt=`        | `<`            |
| `<`           | `<`            |
| `=gt=`        | `>`            |
| `>`           | `>`            |

### Fields

A `fields` key can be added to the URL query to specify columns for the SQL `SELECT` clause. If no fields are specified, the query will be `SELECT *`.

A fields subquery is in the following format:

```
fields=[{column_name}:{alias},... ]
```

where the right hand side of the subquery is a comma separated list of valid column names and an optional alias, with the column name and alias separated by a `:`.

For example, the following queries and GET request are equivalent:

```bash
curl -X GET -s 'http://localhost:8090/authors?fields=surname:last_name,forename' | jq
```

```json
[
  {
    "forename": "Anne",
    "last_name": "Carson"
  },
  {
    "forename": "Anne",
    "last_name": "Brontë"
  },
  {
    "forename": "Virginia",
    "last_name": "Woolf"
  }
]
```

```sql
SELECT surname, forename FROM authors
```

```
 surname | forename
---------+----------
 Woolf   | Virginia
 Brontë  | Anne
 Carson  | Anne
(3 rows)

```

### Joins

Joins can be added to the URL query to add a Join statement to the `SELECT` query.

A join subquery is in the following format:

```
{join_keyword}=[{table}:{left_qualifier}.{left_column}=={right_qualifier}.{right_coulmn};...]
```

where right hand side of the subquery `;` separated list of join relations.

The following join keywords are supported:

- `join`
- `inner_join`
- `left_join`
- `right_join`

For example, the following queries and GET request are equivalent:

```bash
curl -X GET -s 'http://localhost:8090/books?fields=title,name,surname&left_join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id' | jq
```

```json
[
  {
    "name": "Romance",
    "surname": "Carson",
    "title": "Autobiography of Red"
  },
  {
    "name": "Epistolary",
    "surname": "Brontë",
    "title": "The Tenant of Wildfell Hall"
  },
  {
    "name": "Modernism",
    "surname": "Woolf",
    "title": "To The Lighthouse"
  },
  {
    "name": null,
    "surname": "Woolf",
    "title": "Mrs. Dalloway"
  }
]
```

```sql
SELECT name, surname, title FROM books
LEFT JOIN authors on books.author_id=authors.id
LEFT JOIN genres ON books.genre_id=genres.id"
```

```
    name    | surname |            title
------------+---------+-----------------------------
 Romance    | Carson  | Autobiography of Red
 Epistolary | Brontë  | The Tenant of Wildfell Hall
 Modernism  | Woolf   | To The Lighthouse
            | Woolf   | Mrs. Dalloway
(4 rows)
```

## Example usage

### Get table structures

Get a JSON object that describes all the tables in the database with their column names and column types:

```http
GET http://{HOST}:{PORT}/
```

Example:

```bash
curl -X GET 'http://localhost:8090/' | jq
```

```json
{
  "authors": [
    {
      "col_name": "id",
      "col_type": "int32"
    },
    {
      "col_name": "surname",
      "col_type": "string"
    },
    {
      "col_name": "forename",
      "col_type": "string"
    },
    {
      "col_name": "born",
      "col_type": "int16"
    },
    {
      "col_name": "died",
      "col_type": "int16"
    }
  ],
  "books": [
    {
      "col_name": "id",
      "col_type": "int32"
    },
    {
      "col_name": "title",
      "col_type": "string"
    },
    {
      "col_name": "author_id",
      "col_type": "int32"
    }
  ]
}
```

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
[
  {
    "born": 1882,
    "died": 1941,
    "forename": "Virginia",
    "id": 3,
    "surname": "Woolf"
  },
  {
    "born": 1820,
    "died": 1849,
    "forename": "Anne",
    "id": 2,
    "surname": "Brontë"
  },
  {
    "born": 1950,
    "died": null,
    "forename": "Anne",
    "id": 1,
    "surname": "Carson"
  }
]
```

Example (multiple filter parameters connected by `;`):

```bash
curl -X GET -s 'http://localhost:8090/authors?filter=forename==Anne;born>=1900'
```

```json
[
  {
    "born": 1950,
    "died": null,
    "forename": "Anne",
    "id": 1,
    "surname": "Carson"
  }
]
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

## Quick setup/usage

Use recipes in the `justfile` with [casey/just](https://github.com/casey/just) as a task runner.

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

## Limitations:

- Insert, update, or delete requests _cannot_ be made on tables without an `id`
  column
- Read requests ("pick" or "list") _can_ be made on tables without an `id`
  column
- The app will route a request with a RESTful HTTP method + path combination
  for any valid table found in the following example query:

```sql
SELECT tablename FROM Pg_catalog.pg_tables
WHERE schemaname='public'"
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

## Why this project?

This project allows me to get a backend server going as soon as I have my
tables defined for a database. If I decide I need to make changes to the table
structures, then I don't need to make any changes to the backend. This allows
me to perform simple CRUD operations right away and put off writing a more
robust backend until I know exactly what I need.

This makes `gopgrest` good for simple data retrieval on a home server, like
tracking exercise data, managing a personal library, language learning, etc.;
or as a placeholder backend for local development on a frontend application.

I would not use this for a project that publicly exposes sensitive personal
data.

## Security measures

These are the steps I took to be mindful of SQL injection in the service layer:

1. Statements use parameter placeholders for all value
2. Queries can only be made to tables from a whitelist, i.e.
   `Pg_catalog.pg_tables`
3. Columns in JSON requests are checked against the associated table. Any
   invalid column name will throw an error before the query is executed
