CREATE USER authvar_dev_user WITH PASSWORD 'authvar_dev_paSSword';

CREATE DATABASE authvar_db OWNER authvar_dev_user;

CREATE TABLE IF NOT EXISTS public.auth_phone
(
    id              BIGSERIAL PRIMARY KEY,
    phone           VARCHAR(20) NOT NULL,
    code            INTEGER NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    valid_till      TIMESTAMP WITH TIME ZONE NOT NULL NOT NULL,
    confirmed_at    TIMESTAMP WITH TIME ZONE
);

ALTER TABLE public.auth_phone OWNER TO authvar_dev_user;

CREATE TABLE IF NOT EXISTS public.user
(
    id              BIGSERIAL PRIMARY KEY,
    phone           VARCHAR(20) NOT NULL UNIQUE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    delete_at       TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE public.user OWNER TO authvar_dev_user;