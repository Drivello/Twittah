-- auth-init.sql

DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth') THEN
        CREATE DATABASE "auth";
    END IF;
END
$$;

