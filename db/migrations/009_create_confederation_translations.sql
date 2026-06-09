CREATE TABLE confederation_translations (
    confederation_code VARCHAR(20)  NOT NULL REFERENCES confederations(code),
    language           VARCHAR(10)  NOT NULL,
    name               VARCHAR(100) NOT NULL,
    PRIMARY KEY (confederation_code, language)
);
