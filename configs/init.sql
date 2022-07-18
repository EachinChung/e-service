create database if not exists `e-service`;

drop table if exists `e-service`.users;
create table `e-service`.users
(
    id            int unsigned primary key auto_increment,
    user_id       bigint unsigned                            not null comment '用户ID',

    phone         char(11)                                   not null comment '手机号',
    email         varchar(64)                                null comment '邮箱',

    username      varchar(32) binary                         not null comment '用户名',
    password_hash char(60)                                   not null comment '密码',

    avatar        char(55)                                   null comment '头像',

    state         tinyint unsigned default 0                 not null comment '状态',

    created_at    datetime         default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at    datetime         default CURRENT_TIMESTAMP not null comment '更新时间',
    deleted_at    datetime         default null comment '删除时间',

    unique index (user_id),
    unique index (phone),
    unique index (email),

    index (deleted_at),
    index (created_at)
) comment '用户表';
