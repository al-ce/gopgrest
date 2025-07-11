-- init_db.sql

DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'ftrack') THEN
      CREATE USER ftrack WITH PASSWORD 'ftrack';
   END IF;
END
$do$;

CREATE DATABASE ftrack WITH OWNER ftrack ENCODING 'UTF8';

\connect ftrack

ALTER DATABASE ftrack SET timezone TO 'UTC';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

GRANT ALL PRIVILEGES ON DATABASE ftrack TO ftrack;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ftrack;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ftrack;
GRANT ALL PRIVILEGES ON SCHEMA public TO ftrack;

CREATE TABLE IF NOT EXISTS sets (
   name varchar(50) unique not null,
   date timestamp not null default (now()),
   weight decimal not null default 0,
   unit varchar(4) not null,
   reps smallint not null,
   set smallint not null,
   notes text,
   split_day varchar(20),
   program varchar(50),
   tags text
)
