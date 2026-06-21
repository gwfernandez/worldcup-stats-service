DO $$
BEGIN
    CREATE TYPE player_position_code AS ENUM (
        'GK', 'DF', 'MF', 'FW', 'CB', 'RB', 'LB', 'SW', 'RWB', 'LWB',
        'DM', 'CM', 'AM', 'RM', 'LM', 'CF', 'SS', 'RW', 'LW', 'RF', 'LF'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END
$$;

ALTER TABLE players
ADD COLUMN IF NOT EXISTS position player_position_code;

ALTER TABLE players
ADD COLUMN IF NOT EXISTS list_championships INTEGER[];

UPDATE players
SET list_championships = '{}'
WHERE list_championships IS NULL;

ALTER TABLE players
ALTER COLUMN list_championships SET NOT NULL;
