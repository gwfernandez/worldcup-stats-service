CREATE TABLE championships (
    id                BIGSERIAL    PRIMARY KEY,
    year              INTEGER      NOT NULL UNIQUE,
    start_date        DATE         NOT NULL,
    end_date          DATE         NOT NULL,
    host_nation_codes VARCHAR(3)[] NOT NULL,
    champion_code     VARCHAR(3)   REFERENCES national_teams(code)
);

-- Seed data for championships (1930 - 2022)
INSERT INTO championships (year, start_date, end_date, host_nation_codes, champion_code) VALUES
(1930, '1930-07-13', '1930-07-30', ARRAY['URU'], 'URU'),
(1934, '1934-05-27', '1934-06-10', ARRAY['ITA'], 'ITA'),
(1938, '1938-06-04', '1938-06-19', ARRAY['FRA'], 'ITA'),
(1950, '1950-06-24', '1950-07-16', ARRAY['BRA'], 'URU'),
(1954, '1954-06-16', '1954-07-04', ARRAY['SUI'], 'FRG'),
(1958, '1958-06-08', '1958-06-29', ARRAY['SWE'], 'BRA'),
(1962, '1962-05-30', '1962-06-17', ARRAY['CHI'], 'BRA'),
(1966, '1966-07-11', '1966-07-30', ARRAY['ENG'], 'ENG'),
(1970, '1970-05-31', '1970-06-21', ARRAY['MEX'], 'BRA'),
(1974, '1974-06-13', '1974-07-07', ARRAY['FRG'], 'FRG'),
(1978, '1978-06-01', '1978-06-25', ARRAY['ARG'], 'ARG'),
(1982, '1982-06-13', '1982-07-11', ARRAY['ESP'], 'ITA'),
(1986, '1986-05-31', '1986-06-29', ARRAY['MEX'], 'ARG'),
(1990, '1990-06-08', '1990-07-08', ARRAY['ITA'], 'FRG'),
(1994, '1994-06-17', '1994-07-17', ARRAY['USA'], 'BRA'),
(1998, '1998-06-10', '1998-07-12', ARRAY['FRA'], 'FRA'),
(2002, '2002-05-31', '2002-06-30', ARRAY['KOR', 'JPN'], 'BRA'),
(2006, '2006-06-09', '2006-07-09', ARRAY['GER'], 'ITA'),
(2010, '2010-06-11', '2010-07-11', ARRAY['RSA'], 'ESP'),
(2014, '2014-06-12', '2014-07-13', ARRAY['BRA'], 'GER'),
(2018, '2018-06-14', '2018-07-15', ARRAY['RUS'], 'FRA'),
(2022, '2022-11-20', '2022-12-18', ARRAY['QAT'], 'ARG');
