-- +goose Up
-- +goose StatementBegin
create table links(
    id varchar(50) primary key, 
    title varchar(250) not null, 
    href varchar(500) not null, 
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table links;
-- +goose StatementEnd
