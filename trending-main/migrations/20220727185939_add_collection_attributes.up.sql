CREATE TABLE IF NOT EXISTS attributes
(
    id           BIGSERIAL,
    symbol       VARCHAR(500) NOT NULL,
    trait_type   VARCHAR(500) NOT NULL,
    value        TEXT,
    count        INT       DEFAULT 0,
    date_updated TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX ix_attributes_symbol ON attributes (symbol, date_updated);
CREATE UNIQUE INDEX uix_attributes ON attributes (symbol, trait_type, value);

CREATE OR REPLACE FUNCTION trigger_auto_set_date_fields()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.date_updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_set_date_fields
    BEFORE UPDATE OR INSERT
    ON attributes
    FOR EACH ROW
EXECUTE PROCEDURE trigger_auto_set_date_fields();
