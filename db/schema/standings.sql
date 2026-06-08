CREATE TABLE standings (
    team_code VARCHAR(3) PRIMARY KEY REFERENCES teams(code),
    matches_played INTEGER DEFAULT 0 NOT NULL,
    wins INTEGER DEFAULT 0 NOT NULL,
    draws INTEGER DEFAULT 0 NOT NULL,
    losses INTEGER DEFAULT 0 NOT NULL,
    goals_for INTEGER DEFAULT 0 NOT NULL,
    goals_against INTEGER DEFAULT 0 NOT NULL,
    goal_difference INTEGER DEFAULT 0 NOT NULL,
    points INTEGER DEFAULT 0 NOT NULL,
    unified_points INTEGER DEFAULT 0 NOT NULL,
    position INTEGER NOT NULL,
    unified_position INTEGER NOT NULL
);
