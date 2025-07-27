-- init_db.sql
do $do$
begin
    if not exists (
        select
        from
            pg_catalog.pg_roles
        where
            rolname = 'ftrack') then
    create user ftrack with password 'ftrack';
end if;
end
$do$;

create database ftrack with owner ftrack encoding 'utf8';

\connect ftrack
alter database ftrack set timezone to 'utc';

create extension if not exists "uuid-ossp";

grant all privileges on database ftrack to ftrack;

grant all privileges on all tables in schema public to ftrack;

grant all privileges on all sequences in schema public to ftrack;

grant all privileges on schema public to ftrack;
