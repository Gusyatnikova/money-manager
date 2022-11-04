-------------------------------------------------
-- Create table users
-------------------------------------------------

create table if not exists public.user
(
    id       varchar(26) not null primary key,
    user_id  text        not null unique,
    amount   numeric,
    created  timestamp   without time zone not null default now(),
    updated  timestamp   without time zone default null
);

-------------------------------------------------
-- Create table reserve
-------------------------------------------------

create table if not exists public.reserve
(
    user_id    varchar(26) not null references public.user(id),
    service_id text        not null,
    order_id   text        not null,
    amount     numeric     not null,
    created    timestamp   without time zone not null default now(),
    primary key (user_id, service_id, order_id)
);



