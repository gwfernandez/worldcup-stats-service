ALTER TABLE championships
    ADD COLUMN IF NOT EXISTS confederation_codes VARCHAR(20)[];
