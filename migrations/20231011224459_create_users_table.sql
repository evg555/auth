-- +goose Up
create table if not exists users (
    id SERIAL PRIMARY KEY,
    name varchar(100) not null,
    email varchar(100) not null,
    password varchar(255) not null,
    role int2 not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);

-- +goose Down
drop table if exists users;
