-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(55) NOT NULL UNIQUE,
    mail VARCHAR(55) NOT NULL UNIQUE,
    password VARCHAR(55) NOT NULL,
    banner VARCHAR(255),
    pp_url VARCHAR(255),
    bio TEXT,
    status VARCHAR(25) DEFAULT 'online',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
