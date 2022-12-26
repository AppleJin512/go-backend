ALTER TABLE collections
    DROP COLUMN watchlist_count;
ALTER TABLE collections
    ADD COLUMN meta_symbol VARCHAR(500) DEFAULT NULL;
ALTER TABLE collections
    ADD COLUMN update_authority VARCHAR(500) DEFAULT NULL;
