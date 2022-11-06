-------------------------------------------------
-- Create table users
-------------------------------------------------

create table if not exists public.user
(
    id       text        primary key,
    amount   numeric,
    created  timestamp   without time zone not null default now(),
    updated  timestamp   without time zone default null
    constraint amount_nonnegative check (amount >= 0)
    );

-------------------------------------------------
-- Create table reserve
-------------------------------------------------

create table if not exists public.reserve
(
    user_id    text        not null references public.user(id),
    service_id text        not null,
    order_id   text        not null,
    amount     numeric     not null,
    created    timestamp   without time zone not null default now(),
    primary key (user_id, service_id, order_id),
    constraint amount_nonnegative check (amount >= 0)
);

-------------------------------------------------
-- Create table report
-------------------------------------------------

create table if not exists public.report
(
    id         varchar(26) primary key,
    service_id text        not null,
    amount     numeric     not null,
    created    timestamp   without time zone not null default now(),
    constraint amount_nonnegative check (amount >= 0)
);