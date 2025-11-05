CREATE TABLE IF NOT EXISTS weather(
    id BIGSERIAL PRIMARY KEY,
    temperature REAL NOT NULL,
    timestamp TIME NOT NULL
);