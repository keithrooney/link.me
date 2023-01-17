-- +goose Up
-- +goose StatementBegin
create table users(
    id varchar(50) primary key, 
    username varchar(250) unique, 
    email varchar(80) unique,
    password varchar(500) not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
