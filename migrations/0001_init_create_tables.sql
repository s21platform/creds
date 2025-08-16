-- +goose Up
CREATE TABLE IF NOT EXISTS tokens
(
    id        SERIAL PRIMARY KEY,
    token UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS credentials
(
    id        SERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS credentials;