CREATE TABLE national_teams (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(3) NOT NULL UNIQUE,
    dissolution_date DATE,
    confederation_id BIGINT NOT NULL REFERENCES confederations(id),
    federation_name VARCHAR(150) NOT NULL,
    federation_code VARCHAR(10) NOT NULL
);
