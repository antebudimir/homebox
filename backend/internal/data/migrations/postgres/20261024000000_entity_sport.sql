-- +goose Up
-- Add a sport column to entities so the build part picker can group parts by
-- sport. Backfill the seeded build part templates by stripping their name
-- prefix into the sport column ("Skateboard Deck" -> "Deck" + sport
-- "skateboarding"). Generic entities keep their name and stay sportless.
ALTER TABLE entities ADD COLUMN sport character varying(50) NULL;

UPDATE entities SET sport = 'skateboarding', name = substr(name, 12) WHERE name LIKE 'Skateboard %';
UPDATE entities SET sport = 'cycling', name = substr(name, 6) WHERE name LIKE 'Bike %';

-- +goose Down
UPDATE entities SET name = 'Bike ' || name WHERE sport = 'cycling';
UPDATE entities SET name = 'Skateboard ' || name WHERE sport = 'skateboarding';
UPDATE entities SET sport = NULL;
ALTER TABLE entities DROP COLUMN sport;
