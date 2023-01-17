-- +goose Up
-- +goose StatementBegin
alter table links add column user_id varchar(50) references users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table links drop column user;
-- +goose StatementEnd
