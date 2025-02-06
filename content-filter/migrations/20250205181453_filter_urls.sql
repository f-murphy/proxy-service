-- +goose Up
-- +goose StatementBegin
CREATE TABLE filter_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table filter_urls;
-- +goose StatementEnd
