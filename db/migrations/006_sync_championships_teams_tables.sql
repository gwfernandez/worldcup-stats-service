CREATE TABLE IF NOT EXISTS championships_teams_stats (
    year            INTEGER           NOT NULL,
    team_code       VARCHAR(3)        NOT NULL,
    group_code      VARCHAR(2)
        CONSTRAINT championship_team_stats_group_code_check
            CHECK ((group_code)::TEXT ~ '^[A-L1-9]$'::TEXT),
    stage_reached   stage_reached_type,
    matches_played  INTEGER DEFAULT 0 NOT NULL,
    wins            INTEGER DEFAULT 0 NOT NULL,
    draws           INTEGER DEFAULT 0 NOT NULL,
    losses          INTEGER DEFAULT 0 NOT NULL,
    goals_for       INTEGER DEFAULT 0 NOT NULL,
    goals_against   INTEGER DEFAULT 0 NOT NULL,
    goal_difference INTEGER GENERATED ALWAYS AS ((goals_for - goals_against)) STORED,
    points          INTEGER DEFAULT 0 NOT NULL,
    unified_points  INTEGER DEFAULT 0 NOT NULL,
    position        INTEGER,
    CONSTRAINT pk_championship_team_stats
        PRIMARY KEY (year, team_code),
    CONSTRAINT fk_championship_team_stats_teams
        FOREIGN KEY (year, team_code) REFERENCES championships_teams
);

CREATE TABLE IF NOT EXISTS managers (
    id          BIGSERIAL PRIMARY KEY,
    first_name  VARCHAR(100),
    last_name   VARCHAR(100),
    nationality TEXT,
    wikipedia   TEXT,
    ref_id      VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS championships_managers (
    year       INTEGER    NOT NULL,
    team_code  VARCHAR(3) NOT NULL,
    manager_id BIGINT     NOT NULL REFERENCES managers(id),
    CONSTRAINT pk_championships_managers
        PRIMARY KEY (year, team_code, manager_id),
    CONSTRAINT fk_championships_managers_teams
        FOREIGN KEY (year, team_code) REFERENCES championships_teams
);
