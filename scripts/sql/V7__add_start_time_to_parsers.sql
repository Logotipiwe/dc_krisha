alter table krisha.parsers_settings
    add column start_timestamp varchar(255);
update krisha.parsers_settings set parsers_settings.start_timestamp = current_timestamp();