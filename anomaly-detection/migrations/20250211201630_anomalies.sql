-- +goose Up
-- +goose StatementBegin
CREATE TABLE anomalies (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    details TEXT,
    source_ip INET NOT NULL,
    dest_ip INET,
    port INTEGER,
    protocol VARCHAR(50),
    data_volume BIGINT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS anomalies;
-- +goose StatementEnd
