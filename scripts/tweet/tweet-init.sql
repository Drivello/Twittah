-- tweet-init.sql

DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'tweet') THEN
        CREATE DATABASE "tweet";
    END IF;
END
$$;

\connect tweet

-- 02-create-tweets-table.sql
CREATE TABLE IF NOT EXISTS tweets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
