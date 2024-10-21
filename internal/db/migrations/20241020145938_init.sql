-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

-- +goose StatementEnd
create extension if not exists "uuid-ossp";

create table organizations (
    id uuid primary key default uuid_generate_v4(),
    name varchar(256) not null,
    created_at timestamptz not null default now()
);

create table roles (
    id uuid primary key default uuid_generate_v4(),
    name varchar(64) not null
);

create table users (
    id uuid primary key default uuid_generate_v4(),
    name varchar(256) not null,
    last_name varchar(256) not null,
    email varchar(256) not null,
    password varchar(256) not null,
    created_at timestamptz not null default now(),
    organization_id uuid NOT NULL,
    role_id uuid not null
);

insert into
    roles (name)
values
    ('ADMIN'),
    ('USER');

-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

-- +goose StatementEnd
drop table users;

drop table organizations;

drop table roles;