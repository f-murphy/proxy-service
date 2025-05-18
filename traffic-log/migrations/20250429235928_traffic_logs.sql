-- +goose Up
-- +goose StatementBegin
CREATE TABLE traffic_logs (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    client_ip VARCHAR(50),
    method VARCHAR(10),
    url TEXT,
    response_code INT,
    bytes_sent BIGINT,
    bytes_received BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE traffic_logs;
-- +goose StatementEnd
