use `krisha`;

drop table if exists parsers_settings;
drop table if exists apartments;
drop table if exists allowed_chats;
create table if not exists parsers_settings (
    chat_id bigint not null primary key,
    filters text,
    interval_sec int not null default 120,
    enabled bool not null default false
);

create table if not exists apartments (
    id varchar(255) not null primary key,
    data_json text
);

create table if not exists allowed_chats (
    chat_id bigint not null primary key
);
