-- +goose Up
-- Replace the build is_active boolean with a status enum:
-- assembled (rideable), disassembled (torn down), retired (kept for history).
alter table builds add column status text not null default 'assembled';
update builds set status = case when is_active then 'assembled' else 'disassembled' end;
alter table builds drop column is_active;

-- +goose Down
alter table builds add column is_active boolean not null default false;
update builds set is_active = case when status = 'assembled' then true else false end;
alter table builds drop column status;