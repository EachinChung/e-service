create database service;

drop table if exists super_users;
drop table if exists users;
create table users
(
    id            serial primary key,
    eid           varchar(32) unique       not null,
    phone         char(11) unique          not null,

    password_hash char(60)                 not null,
    nickname      varchar(32)              not null,
    avatar        char(55)                 null,
    state         smallint                 not null default 0,

    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default ('now'::text)::timestamp(0) with time zone,
    deleted_at    timestamp with time zone null
);

create index users_deleted_at_key on users (deleted_at);

INSERT INTO users (eid, phone, password_hash, nickname)
VALUES ('Eachin', '13711164450', '$2a$10$U6IWJS3.fy1wUa/I2GnHeOQXGU.VZMVirjEO.xb/meuUraBCpzo2i', 'Eachin');


drop table if exists super_users;
create table super_users
(
    id         serial primary key,
    eid        varchar(32) unique       not null,
    created_at timestamp with time zone not null default now(),
    foreign key (eid) references users (eid) on delete cascade on update cascade
);

INSERT INTO public.super_users (eid)
VALUES ('Eachin');
