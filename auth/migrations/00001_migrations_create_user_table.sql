-- +goose Up
CREATE table auth_user (
    id serial primary key,
    name varchar(50) not null,
    role int not null default 1,
    email varchar(50) not null unique,
    password varchar(60) not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
DROP table auth_user;
