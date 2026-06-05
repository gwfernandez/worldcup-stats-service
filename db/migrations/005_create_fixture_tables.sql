DO $$
BEGIN
    CREATE TYPE stage_reached_type AS ENUM (
        'group_stage',
        'second_group_stage',
        'round_of_16',
        'quarter_finals',
        'semi_finals',
        'third_place',
        'final'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE stage_type AS ENUM (
        'group',
        'knockout'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS stadiums (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL
);

CREATE TABLE IF NOT EXISTS championships_teams (
    year INTEGER NOT NULL REFERENCES championships(year),
    team_code VARCHAR(3) NOT NULL REFERENCES teams(code),
    PRIMARY KEY (year, team_code)
);

CREATE TABLE IF NOT EXISTS championships_groups_stats (
    year            INTEGER            NOT NULL,
    stage           stage_reached_type NOT NULL,
    group_code      VARCHAR(2)         NOT NULL
        CONSTRAINT championships_groups_stats_group_code_check
            CHECK ((group_code)::TEXT ~ '^[A-L1-9]$'::TEXT),
    team_code       VARCHAR(3)         NOT NULL,
    matches_played  INTEGER DEFAULT 0  NOT NULL,
    wins            INTEGER DEFAULT 0  NOT NULL,
    draws           INTEGER DEFAULT 0  NOT NULL,
    losses          INTEGER DEFAULT 0  NOT NULL,
    goals_for       INTEGER DEFAULT 0  NOT NULL,
    goals_against   INTEGER DEFAULT 0  NOT NULL,
    goal_difference INTEGER GENERATED ALWAYS AS ((goals_for - goals_against)) STORED,
    points          INTEGER DEFAULT 0  NOT NULL,
    unified_points  INTEGER DEFAULT 0  NOT NULL,
    position        INTEGER,
    CONSTRAINT pk_championship_group_stats
        PRIMARY KEY (year, stage, group_code, team_code),
    CONSTRAINT fk_championship_group_stats_teams
        FOREIGN KEY (year, team_code) REFERENCES championships_teams
);

CREATE TABLE IF NOT EXISTS matches (
    id                        BIGSERIAL PRIMARY KEY,
    year                      INTEGER NOT NULL,
    stage                     stage_reached_type,
    group_code                VARCHAR(2)
        CONSTRAINT matches_group_code_check
            CHECK ((group_code)::TEXT ~ '^[A-L1-9]$'::TEXT),
    replayed                  BOOLEAN DEFAULT FALSE NOT NULL,
    replay_of                 BIGINT REFERENCES matches,
    match_date                DATE,
    match_time                TIME,
    stadium_id                BIGINT REFERENCES stadiums,
    home_team_code            VARCHAR(3) NOT NULL,
    away_team_code            VARCHAR(3) NOT NULL,
    home_team_score           INTEGER,
    away_team_score           INTEGER,
    extra_time                BOOLEAN DEFAULT FALSE NOT NULL,
    penalty_shootout          BOOLEAN DEFAULT FALSE NOT NULL,
    home_team_score_penalties INTEGER,
    away_team_score_penalties INTEGER,
    home_team_win             BOOLEAN,
    away_team_win             BOOLEAN,
    draw                      BOOLEAN,
    ref_id                    VARCHAR(20),
    stage_type                stage_type NOT NULL,
    CONSTRAINT matches_away_team_fk
        FOREIGN KEY (year, away_team_code) REFERENCES championships_teams,
    CONSTRAINT matches_home_team_fk
        FOREIGN KEY (year, home_team_code) REFERENCES championships_teams
);
