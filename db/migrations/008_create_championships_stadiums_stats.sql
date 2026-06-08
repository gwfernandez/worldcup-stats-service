ALTER TABLE stadiums
    ADD COLUMN IF NOT EXISTS city_name TEXT,
    ADD COLUMN IF NOT EXISTS country VARCHAR(3) REFERENCES teams(code),
    ADD COLUMN IF NOT EXISTS capacity INTEGER,
    ADD COLUMN IF NOT EXISTS wikipedia_link TEXT,
    ADD COLUMN IF NOT EXISTS ref_id VARCHAR(20);

CREATE TABLE IF NOT EXISTS championships_stadiums_stats (
    year           INTEGER DEFAULT 0 NOT NULL REFERENCES championships(year),
    stadium_id     BIGINT            NOT NULL REFERENCES stadiums(id),
    matches_played INTEGER DEFAULT 0 NOT NULL,
    CONSTRAINT pk_championships_stadiums_stats
        PRIMARY KEY (year, stadium_id)
);
