-- +goose Up
-- +goose StatementBegin
CREATE TABLE filter_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL
);

CREATE TABLE blacklist (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL CHECK (type IN ('url', 'keyword')),
    value VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO blacklist (type, value) VALUES 
('url', 'example.com/badpage'),
('keyword', 'hack'),
('keyword', 'malware');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table filter_urls;
drop table blacklist;
-- +goose StatementEnd
