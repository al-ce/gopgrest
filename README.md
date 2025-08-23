# gopgrest

A dynamic RESTful HTTP server for a Postgres database. The app will adapt its
URL routing and SQL queries to the provided table schemas. Supports an RSQL
based query language.

## API

| Endpoint                     | Method | Description                 | Request            | Response                                  |
| ---------------------------- | ------ | --------------------------- | ------------------ | ----------------------------------------- |
| `/`                          | GET    | Get structure of all tables | ---                | `application/json` (tables)               |
| `/{tablename}`               | POST   | Insert new row(s)           | `application/json` | `application/json` (new row ids)          |
| `/{tablename}/{id}`          | GET    | Get a row by ID             | ---                | `application/json` (found row)            |
| `/{tablename}?{querystring}` | GET    | Get rows by query params    | ---                | `application/json` (matching rows)        |
| `/{tablename}/{id}`          | PUT    | Update a row by ID          | `application/json` | `row {id} deleted from table {tablename}` |
| `/{tablename}?{querystring}` | PUT    | Update rows by query params | `application/json` | `rows updated in table {tablename}: {n}`  |
| `/{tablename}/{id}`          | DELETE | Delete a row by ID          | ---                | `row {id} deleted from table {tablename}` |
| `/{tablename}?{querystring}` | DELETE | Delete rows by query params | ---                | `rows deleted in table {tablename}: {n}`  |

## REST Query language (based on restSQL)

This project is aiming to implement a URL query parameter parser similar to [restSQL](http://restsql.org/doc/Overview.html).

Query parameters can be added to a GET request after a `?` query separator. Keys and their values are separated by `=`. Multiple subqueries can be joined with `&`.

The following example contains:

- a `select` subquery to select columns to return, with qualifiers and aliases on some columns
- a `left_join` subquery to add two join relations
- a `where` subquery to add conditions to the query

```bash
curl -X GET -s 'http://localhost:8090/books?select=title,genres.name:genre,authors.surname:author&left_join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id&where=born<1900' | jq
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
WHERE authors.born < 1900
```

The following query keys are supported:

| Key          | Description                                    |
| ------------ | ---------------------------------------------- |
| `where`      | add `WHERE` conditions to a `SELECT` query     |
| `select`     | columns to return in a `SELECT` query          |
| `inner join` | add `INNER JOIN` relations to a `SELECT` query |
| `join`       | add `INNER JOIN` relations to a `SELECT` query |
| `left_join`  | add `LEFT JOIN` relations to a `SELECT` query  |
| `right_join` | add `RIGHT JOIN` relations to a `SELECT` query |

Query parameters matching the `where` format for an RSQL query can be added to PUT and DELETE requests to update/delete rows matching the conditions.

```bash
curl -X PUT 'http://localhost:8090/authors?forename==Anne;born<1900' --data '{"forename": "Emily"}'
```

```
rows updated in table authors: 1
```

Unlike a GET request, query parameters in a POST or DELETE request do not require the `where` key or a `=` separator before the column/value conditionals.

### Encoding in the URL query parameters

Percent-encoded characters will be decoded as defined at <a href="https://developer.mozilla.org/en-US/docs/Glossary/Percent-encoding">MDN Web Docs: Percent-encoding</a>

For example, here is a query where space '` `' characters are encoded with `+` and `+` characters are encoded with `%2B`:

```bash
curl -X GET -s 'http://localhost:8090/books?select=title&where=title==Programming:+Principles+and+Practice+Using+C%2B%2B' | jq
```

```json
[
  {
    "title": "Programming: Principles and Practice Using C++"
  }
]
```

## Query parameter specifications

### Where

A `where` key can be added to a GET URL's query parameters to match a SQL `WHERE` clause.

A `where` subquery is in the following format:

```
where={column_name}{operator}{value};...
```

where the right of the `=` is a `;` separated list of conditional expressions equivalent to a `WHERE` clause.

For example, the following SQL query and GET request are equivalent:

```bash
curl -X GET -s 'http://localhost:8090/authors?where=forename==Anne;born>=1900' | jq
```

responds with:

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

These are the currently supported conditional operators:

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

As noted above, a list of `;` separated conditionals can be added to PUT and DELETE requests _without_ the preceding `where=` key/assignment to add a `WHERE` clause to `UPDATE` or `DELETE` queries.

For example, the following SQL query and PUT request are equivalent:

```bash
curl -X PUT 'http://localhost:8090/authors?forename==Anne;born<1900' --data '{"forename": "Emily"}'
```

```sql
UPDATE authors SET forename = 'Emily' WHERE forename = 'Anne' AND born < 1900
```

And the following SQL query and DELETE request are equivalent:

```bash
curl -X DELETE 'http://localhost:8090/books?title=like=Autobiography%'
```

```sql
DELETE FROM books WHERE title LIKE
```

### Select

A `select` key can be added to the URL query to specify columns for the SQL `SELECT` clause. If no columns are specified, the query will be `SELECT *`.

A `select` subquery is in the following format:

```
select=[{column_name}:{alias},... ]
```

where the right of the `=` is a `,` separated list of valid column names and an optional alias, with the column name and alias separated by a `:`.

For example, the following queries and GET request are equivalent:

```bash
curl -X GET -s 'http://localhost:8090/authors?select=surname:last_name,forename' | jq
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
{join_keyword}=[{table}:{left_qualifier}.{left_column}=={right_qualifier}.{right_coulmn};...
```

where right hand side of the subquery `;` separated list of join relations.

The following join keywords are supported:

- `join`
- `inner_join`
- `left_join`
- `right_join`

For example, the following queries and GET request are equivalent:

```bash
curl -X GET -s 'http://localhost:8090/books?select=title,name,surname&left_join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id' | jq
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

Insert a new row or rows into an existing table. Responds with array of new ids.

```http
POST http://{HOST}:{PORT}/{tablename}
Accept: application/json
```

Example (single JSON object):

```bash
curl -X POST -s http://localhost:8090/authors \
      --data '{ "surname": "Woolf", "forename": "Virginia" }'
```
```json
[3]
```

Example (multiple rows from array of JSON objects):

```bash
curl -X POST http://localhost:8090/authors \
      --data '[{ "surname": "Groton", "forename": "Anne" }, { "surname": "Plato" }]'
```

```json
[4, 5]
```

### Get Row (pick)

Get a single row from a table by id as a JSON object.

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

### Get Rows (list)

Get all rows from a table matching a list of optional query parameters as an array of JSON objects.

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

Example (multiple `where` parameters connected by `;`):

```bash
curl -X GET -s 'http://localhost:8090/authors?where=forename==Anne;born>=1900'
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

Update a row by ID or by query parameters, responding with the number of rows updated.

By id:

```http
PUT http://{HOST}:{PORT}/{tablename}/{id}
Accept: application/json
```

```bash
curl -X PUT -s http://localhost:8090/authors/3 --data '{"surname" : "Woolf"}'
```

```
"rows updated in table authors: 1"
```

By query parameters:

```http
PUT http://{HOST}:{PORT}/{tablename}?{querystring}
```

```bash
curl -X PUT 'http://localhost:8090/authors?forename==Anne;born<1900' --data '{"forename": "Emily"}'
# stdout: "rows updated in table authors: 1"

```

### Delete

Delete a row by ID, responding with a message confirming the deleted row.

```http
DELETE http://{HOST}:{PORT}/{tablename}/{id}
```

```bash
curl -X DELETE -s http://localhost:8090/authors/5
# stdout: "row 5 deleted from table authors"
```

By query parameters:

```http
PUT http://{HOST}:{PORT}/{tablename}?{querystring}
```

```bash
curl -X DELETE 'http://localhost:8090/books?title=like=Autobiography%'
# stdout: "rows deleted in table books: 1"

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
