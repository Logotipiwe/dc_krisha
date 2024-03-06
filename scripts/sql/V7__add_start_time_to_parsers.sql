alter table parsers_settings
    add column start_timestamp varchar(255);
update parsers_settings set parsers_settings.start_timestamp = current_timestamp();