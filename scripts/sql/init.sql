create database krisha;



drop table if exists parsers_settings;
drop table if exists apartments;
drop table if exists allowed_chats;
create table if not exists parsers_settings
(
    chat_id               bigint       not null primary key,
    filters               text,
    interval_sec          int          not null default 120,
    aps_limit             int          not null,
    enabled               boolean      not null default false,
    is_granted_explicitly bool                  default true not null,
    chat_name             varchar(255) null,
    start_timestamp       varchar(255)          default current_timestamp,
    curr_aps_count int not null default 0
);

create table if not exists apartments
(
    id        varchar(255) not null primary key,
    data_json text
);

create table if not exists allowed_chats
(
    chat_id bigint not null primary key
);

CREATE TABLE messages_log
(
    id              VARCHAR(36) PRIMARY KEY,
    chat_id         BIGINT                              NOT NULL,
    text            TEXT,
    direction       varchar(255)                        NOT NULL,
    additional_data TEXT,
    time            timestamp default current_timestamp not null
);

create table if not exists known_chats (
    chat_id bigint not null primary key,
    chat_info text
);