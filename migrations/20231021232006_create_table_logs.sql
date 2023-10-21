-- +goose Up
create table if not exists logs (
     id SERIAL PRIMARY KEY,
     method varchar(100) not null,
     data json,
     created_at timestamptz not null default now()
);

-- +goose Down
DROP TABLE IF EXISTS logs;
