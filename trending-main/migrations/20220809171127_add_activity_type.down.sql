DROP INDEX ix_activity_type;

ALTER TABLE activities
    DROP COLUMN activity_type;

DROP TYPE ACTIVITY_TYPE;
