CREATE TABLE messages_log (
    id              VARCHAR(36) PRIMARY KEY,
    chat_id         BIGINT NOT NULL,
    text            TEXT,
    direction       varchar(255) NOT NULL,
    additional_data TEXT
);
alter table messages_log
    add `order` int auto_increment,
    add constraint messages_log_pk
        unique key (`order`);
alter table messages_log
    add time timestamp default current_timestamp() not null;
