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

create type user_role as enum ('admin', 'user');

create table users (
    id uuid primary key default uuid_generate_v4(),
    name varchar(256),
    last_name varchar(256),
    email varchar(256) not null,
    password text not null,
    created_at timestamptz not null default now(),
    organization_id uuid,
    role user_role not null default 'user'
);

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION create_organization(org_name VARCHAR(256)) RETURNS UUID AS $$ DECLARE new_org_id UUID;

BEGIN
INSERT INTO
    organizations (name)
VALUES
    (org_name) RETURNING id INTO new_org_id;

RETURN new_org_id;

END;

$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

-- +goose StatementEnd
drop table users;

drop table organizations;

drop table roles;

drop type user_role;