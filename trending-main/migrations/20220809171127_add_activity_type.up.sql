CREATE TYPE ACTIVITY_TYPE AS ENUM ('listing', 'sale', 'update_price');

ALTER TABLE activities
    ADD COLUMN activity_type ACTIVITY_TYPE DEFAULT 'sale';

CREATE INDEX ix_activity_type ON activities (activity_type);
