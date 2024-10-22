-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

-- +goose StatementEnd
create extension if not exists "uuid-ossp";

create table organizations (
    id uuid primary key default uuid_generate_v4(),
    name varchar(256),
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

create table tasks (
    id uuid primary key default uuid_generate_v4(),
    name varchar(256),
    description text,
    created_at timestamptz not null default now(),
    user_id uuid,
    organization_id uuid
)

-- Relationships

alter table users add constraint users_organizations_fkey foreign key (organization_id) references organizations(id) on delete cascade;
alter table tasks add constraint tasks_users_fkey foreign key (user_id) references users(id) on delete cascade;
alter table tasks add constraint tasks_organizations_fkey foreign key (organization_id) references organizations(id) on delete cascade;

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION create_organization() RETURNS UUID AS $$ DECLARE new_org_id UUID;

BEGIN
INSERT INTO organizations default values RETURNING id INTO new_org_id;

RETURN new_org_id;

END;

$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION assign_organization_to_admin() RETURNS TRIGGER AS $$
DECLARE
    new_org_id UUID;
BEGIN
    IF NEW.role = 'admin' THEN
        new_org_id := create_organization();
        
        NEW.organization_id := new_org_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER create_organization_for_admin BEFORE
INSERT
    ON users FOR EACH ROW
    WHEN (NEW.role = 'admin') EXECUTE FUNCTION assign_organization_to_admin();
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