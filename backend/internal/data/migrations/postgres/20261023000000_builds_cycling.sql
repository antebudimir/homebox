-- +goose Up
-- Migrate existing builds off the old bmx slug to match the sport config
-- naming (BMX -> Cycling). Fresh installs never see the old value.
UPDATE builds SET sport = 'cycling' WHERE sport = 'bmx';

-- +goose Down
UPDATE builds SET sport = 'bmx' WHERE sport = 'cycling';
