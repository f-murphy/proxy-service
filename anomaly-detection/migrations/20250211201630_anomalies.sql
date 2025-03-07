-- +goose Up
-- +goose StatementBegin
CREATE TABLE request_logs (
    ip_address VARCHAR(45) PRIMARY KEY,
    request_count INTEGER NOT NULL,
    last_request TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS request_logs;
-- +goose StatementEnd
