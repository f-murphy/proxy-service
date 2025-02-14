-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklist (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL CHECK (type IN ('url', 'keyword')),
    value VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE blocked_ips (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(50) NOT NULL UNIQUE,
    blocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO blacklist (type, value) VALUES 
('url', 'example.com/badpage'),
('keyword', 'hack'),
('keyword', 'malware');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table blocked_ips;
drop table blacklist;
-- +goose StatementEnd
