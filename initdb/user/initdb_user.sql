-- Ensure 'user' database exists
DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user') THEN
      PERFORM dblink_exec('dbname=postgres user=user password=user', 'CREATE DATABASE "user"');
   END IF;
END
$do$;

-- Connect to 'user' database and create 'users' and 'follows' tables if not exist
\c "user"

-- Minimalistic user table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY,
    username VARCHAR(64) NOT NULL
);

-- Follow relationships table
CREATE TABLE IF NOT EXISTS follows (
    id SERIAL PRIMARY KEY,
    follower_id VARCHAR(64) NOT NULL REFERENCES users(id),
    followee_id VARCHAR(64) NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
