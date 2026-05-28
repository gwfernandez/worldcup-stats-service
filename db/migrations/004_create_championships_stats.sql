CREATE TABLE championships_stats (
    year              INTEGER    PRIMARY KEY REFERENCES championships(year),
    total_teams       INTEGER    NOT NULL DEFAULT 0,
    total_matches     INTEGER    NOT NULL DEFAULT 0,
    total_stadiums    INTEGER    NOT NULL DEFAULT 0,
    total_players     INTEGER    NOT NULL DEFAULT 0,
    total_goals       INTEGER    NOT NULL DEFAULT 0,
    champion_code     VARCHAR(3) REFERENCES national_teams(code),
    runner_up_code    VARCHAR(3) REFERENCES national_teams(code),
    third_place_code  VARCHAR(3) REFERENCES national_teams(code),
    fourth_place_code VARCHAR(3) REFERENCES national_teams(code),
    top_scorer_ids    BIGINT[],
    top_scorer_goals  INTEGER    NOT NULL DEFAULT 0
);
