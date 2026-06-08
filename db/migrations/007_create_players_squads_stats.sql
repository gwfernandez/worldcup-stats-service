CREATE TABLE IF NOT EXISTS players (
    id              BIGSERIAL PRIMARY KEY,
    first_name      VARCHAR(100),
    last_name       VARCHAR(100),
    nationality     TEXT,
    birth_date      DATE,
    death_date      DATE,
    wikipedia       TEXT,
    ref_id          VARCHAR(20),
    CONSTRAINT players_name_not_empty
        CHECK (
            NULLIF(TRIM(COALESCE(first_name, '')), '') IS NOT NULL
            OR NULLIF(TRIM(COALESCE(last_name, '')), '') IS NOT NULL
        )
);

CREATE TABLE IF NOT EXISTS squads (
    year        INTEGER    NOT NULL,
    team_code   VARCHAR(3) NOT NULL,
    player_id   BIGINT     NOT NULL REFERENCES players(id),
    shirt_number INTEGER,
    position    VARCHAR(50),
    CONSTRAINT pk_squads PRIMARY KEY (year, team_code, player_id),
    CONSTRAINT fk_squads_teams
        FOREIGN KEY (year, team_code) REFERENCES championships_teams
);

CREATE TABLE IF NOT EXISTS squads_stats (
    year         INTEGER    NOT NULL,
    team_code    VARCHAR(3) NOT NULL,
    player_id    BIGINT     NOT NULL,
    appearances  INTEGER    NOT NULL DEFAULT 0,
    goals        INTEGER    NOT NULL DEFAULT 0,
    CONSTRAINT pk_squads_stats PRIMARY KEY (year, team_code, player_id),
    CONSTRAINT fk_squads_stats_squads
        FOREIGN KEY (year, team_code, player_id) REFERENCES squads,
    CONSTRAINT squads_stats_appearances_non_negative CHECK (appearances >= 0),
    CONSTRAINT squads_stats_goals_non_negative CHECK (goals >= 0)
);
