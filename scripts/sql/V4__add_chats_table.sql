create table if not exists known_chats (
    chat_id bigint not null primary key,
    chat_info text
);