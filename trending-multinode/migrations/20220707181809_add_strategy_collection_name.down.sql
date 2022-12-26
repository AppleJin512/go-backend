ALTER TABLE collections
    DROP COLUMN meta_symbol;
ALTER TABLE collections
    DROP COLUMN update_authority;
ALTER TABLE collections
    ADD COLUMN watchlist_count INTEGER DEFAULT 0;

