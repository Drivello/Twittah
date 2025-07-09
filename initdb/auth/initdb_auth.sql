-- Ensure 'auth' database exists
DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth') THEN
      PERFORM dblink_exec('dbname=postgres user=auth password=auth', 'CREATE DATABASE auth');
   END IF;
END
$do$;

-- Connect to 'auth' database and create 'user' table if it does not exist
\c auth

CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
