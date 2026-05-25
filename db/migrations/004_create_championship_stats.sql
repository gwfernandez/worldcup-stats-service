CREATE TABLE championship_stats (
    id                BIGINT     PRIMARY KEY REFERENCES championships(id),
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

-- Seed data for championship_stats (1930 - 2022)
INSERT INTO championship_stats (id, total_teams, total_matches, total_stadiums, total_players, total_goals, champion_code, runner_up_code, third_place_code, fourth_place_code, top_scorer_ids, top_scorer_goals) VALUES
(1, 13, 18, 3, 189, 70, 'URU', 'ARG', 'USA', 'YUG', NULL, 8),
(2, 16, 17, 8, 358, 70, 'ITA', 'TCH', 'GER', 'AUT', NULL, 5),
(3, 15, 18, 10, 288, 84, 'ITA', 'HUN', 'BRA', 'SWE', NULL, 7),
(4, 13, 22, 6, 285, 88, 'URU', 'BRA', 'SWE', 'ESP', NULL, 8),
(5, 16, 26, 6, 360, 140, 'FRG', 'HUN', 'AUT', 'URU', NULL, 11),
(6, 16, 35, 12, 320, 126, 'BRA', 'SWE', 'FRA', 'FRG', NULL, 13),
(7, 16, 32, 4, 320, 89, 'BRA', 'TCH', 'CHI', 'YUG', NULL, 4),
(8, 16, 32, 8, 320, 89, 'ENG', 'FRG', 'POR', 'URS', NULL, 9),
(9, 16, 32, 5, 320, 95, 'BRA', 'ITA', 'FRG', 'URU', NULL, 10),
(10, 16, 38, 9, 320, 97, 'FRG', 'NED', 'POL', 'BRA', NULL, 7),
(11, 16, 38, 6, 320, 102, 'ARG', 'NED', 'BRA', 'ITA', NULL, 6),
(12, 24, 52, 17, 528, 146, 'ITA', 'FRG', 'POL', 'FRA', NULL, 6),
(13, 24, 52, 12, 528, 132, 'ARG', 'FRG', 'FRA', 'BEL', NULL, 6),
(14, 24, 52, 12, 528, 115, 'FRG', 'ARG', 'ITA', 'ENG', NULL, 6),
(15, 24, 52, 9, 528, 141, 'BRA', 'ITA', 'SWE', 'BUL', NULL, 6),
(16, 32, 64, 10, 704, 171, 'FRA', 'BRA', 'CRO', 'NED', NULL, 6),
(17, 32, 64, 20, 736, 161, 'BRA', 'GER', 'TUR', 'KOR', NULL, 8),
(18, 32, 64, 12, 736, 147, 'ITA', 'FRA', 'GER', 'POR', NULL, 5),
(19, 32, 64, 10, 736, 145, 'ESP', 'NED', 'GER', 'URU', NULL, 5),
(20, 32, 64, 12, 736, 171, 'GER', 'ARG', 'NED', 'BRA', NULL, 6),
(21, 32, 64, 12, 736, 169, 'FRA', 'CRO', 'BEL', 'ENG', NULL, 6),
(22, 32, 64, 8, 832, 172, 'ARG', 'FRA', 'CRO', 'MAR', NULL, 8);
