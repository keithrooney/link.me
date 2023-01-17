-- +goose Up
-- +goose StatementBegin
create index links_title_idx on links (title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index links_title_idx;
-- +goose StatementEnd
