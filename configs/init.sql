create database service;

drop table if exists users;
create table users
(
    id            serial primary key,
    username      varchar(32) unique       not null,
    phone         char(11) unique          not null,

    password_hash char(60)                 not null,
    avatar        char(55)                 null,
    state         smallint                 not null default 0,

    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default ('now'::text)::timestamp(0) with time zone,
    deleted_at    timestamp with time zone null
);

create index users_deleted_at_key on users (deleted_at);

drop table if exists super_users;
create table super_users
(
    id         serial primary key,
    username   varchar(32) unique       not null,
    created_at timestamp with time zone not null default now(),
    foreign key (username) references users (username) on delete cascade on update cascade
);
