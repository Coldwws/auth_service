-- +goose Up
-- Создаём таблицу пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Таблица refresh токенов
CREATE TABLE refresh_tokens (
                                id SERIAL PRIMARY KEY,
                                user_id INT REFERENCES users(id) ON DELETE CASCADE,
                                token TEXT NOT NULL,
                                expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Таблица access к эндпоинтам
CREATE TABLE access_rules (
                              id SERIAL PRIMARY KEY,
                              endpoint_address TEXT NOT NULL,
                              description TEXT,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Индексы
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_access_rules_endpoint ON access_rules(endpoint_address);

-- +goose Down
DROP TABLE IF EXISTS access_rules;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;