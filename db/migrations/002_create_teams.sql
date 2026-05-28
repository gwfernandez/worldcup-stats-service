CREATE TABLE teams (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    dissolution_date DATE,
    confederation_code VARCHAR(20) NOT NULL REFERENCES confederations(code),
    federation_name VARCHAR(150) NOT NULL,
    federation_code VARCHAR(10) NOT NULL
);
