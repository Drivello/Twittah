-- user-init.sql

DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user') THEN
        CREATE DATABASE "user";
    END IF;
END
$$;

\connect user

-- 02-create-user-table.sql
CREATE TABLE IF NOT EXISTS "user" (
    id BIGINT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 03-create-follow-table.sql
CREATE TABLE IF NOT EXISTS "follow" (
    id BIGINT PRIMARY KEY,
    follower_id BIGINT NOT NULL,
    followee_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES "user"(id),
    CONSTRAINT fk_followee FOREIGN KEY (followee_id) REFERENCES "user"(id)
);
