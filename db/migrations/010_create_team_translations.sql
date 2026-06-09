CREATE TABLE team_translations (
    team_code VARCHAR(3)   NOT NULL REFERENCES teams(code),
    language  VARCHAR(10)  NOT NULL,
    name      VARCHAR(100) NOT NULL,
    PRIMARY KEY (team_code, language)
);
