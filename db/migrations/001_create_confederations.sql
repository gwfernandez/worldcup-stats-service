CREATE TABLE confederations (
    id      BIGSERIAL    PRIMARY KEY,
    code    VARCHAR(20)  NOT NULL UNIQUE,
    name    VARCHAR(100) NOT NULL
);
