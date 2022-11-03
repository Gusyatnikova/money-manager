-------------------------------------------------
-- Create table user
-------------------------------------------------

create table if not exists public.user
(
    id       bigserial    not null primary key,
    user_id  varchar(27) not null unique,
    amount   numeric,
    created  timestamp    without time zone not null default now(),
    updated  timestamp    without time zone null
);
