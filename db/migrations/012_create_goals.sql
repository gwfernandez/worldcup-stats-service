DO $$
BEGIN
    CREATE TYPE match_period_type AS ENUM ('regular', 'extra_time');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE match_period_half AS ENUM ('first_half', 'second_half');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE match_period_stoppage AS ENUM ('regular', 'stoppage');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE match_team_condition AS ENUM ('home', 'away');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE SEQUENCE IF NOT EXISTS matches_goals_id_seq;

CREATE TABLE IF NOT EXISTS goals (
    id              BIGINT DEFAULT nextval('matches_goals_id_seq'::regclass) NOT NULL PRIMARY KEY,
    year            INTEGER NOT NULL REFERENCES championships,
    match_id        BIGINT REFERENCES matches,
    player_id       BIGINT REFERENCES players,
    team_code       VARCHAR(3) NOT NULL REFERENCES teams,
    team_condition  match_team_condition,
    shirt_number    INTEGER DEFAULT 0 NOT NULL,
    minute_regular  INTEGER DEFAULT 0 NOT NULL,
    minute_stoppage INTEGER DEFAULT 0 NOT NULL,
    period_type     match_period_type NOT NULL,
    period_half     match_period_half NOT NULL,
    period_stoppage match_period_stoppage NOT NULL,
    penalty         BOOLEAN,
    own_goal        BOOLEAN,
    ref_id          VARCHAR(20),
    CONSTRAINT goals_championship_fkey
        FOREIGN KEY (year, team_code) REFERENCES championships_teams
);

ALTER SEQUENCE matches_goals_id_seq OWNED BY goals.id;

CREATE INDEX IF NOT EXISTS goals_player_id_year_idx
    ON goals (player_id, year);
