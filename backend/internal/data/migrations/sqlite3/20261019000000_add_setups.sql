-- +goose Up
-- Setups: named assemblies of parts (e.g. skateboard / bike builds).
-- setup_parts carries added/removed dates so mounting history is preserved.
create table if not exists setups
(
    id          uuid     not null
        primary key,
    created_at  datetime not null,
    updated_at  datetime not null,
    name        text     not null,
    description text,
    is_active   boolean  not null default false,
    sport       text,
    group_setups uuid    not null
        constraint setups_groups_setups
            references groups
            on delete cascade
);

create index if not exists setups_group_setups
    on setups (group_setups);

create table if not exists setup_parts
(
    id           uuid     not null
        primary key,
    created_at   datetime not null,
    updated_at   datetime not null,
    added_date   datetime not null,
    removed_date datetime,
    position     text,
    entity_id    uuid     not null
        constraint setup_parts_entities_setup_parts
            references entities
            on delete cascade,
    setup_id     uuid     not null
        constraint setup_parts_setups_parts
            references setups
            on delete cascade
);

create index if not exists setup_parts_entity_id
    on setup_parts (entity_id);

create index if not exists setup_parts_setup_id
    on setup_parts (setup_id);

create table if not exists usage_sessions
(
    id               uuid     not null
        primary key,
    created_at       datetime not null,
    updated_at       datetime not null,
    date             datetime,
    duration_minutes integer  not null default 0,
    location         text,
    conditions       text,
    notes            text,
    rating           integer,
    setup_id         uuid     not null
        constraint usage_sessions_setups_usage_sessions
            references setups
            on delete cascade
);

create index if not exists usage_sessions_setup_id
    on usage_sessions (setup_id);

-- +goose Down
drop table if exists usage_sessions;
drop table if exists setup_parts;
drop table if exists setups;
