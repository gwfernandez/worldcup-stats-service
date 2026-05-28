CREATE TABLE championships (
    year              INTEGER      PRIMARY KEY,
    start_date        DATE         NOT NULL,
    end_date          DATE         NOT NULL,
    host_codes        VARCHAR(3)[] NOT NULL,
    champion_code     VARCHAR(3)   REFERENCES teams(code)
);