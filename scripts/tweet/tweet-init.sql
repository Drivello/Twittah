-- tweet-init.sql

DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'tweet') THEN
        CREATE DATABASE "tweet";
    END IF;
END
$$;
