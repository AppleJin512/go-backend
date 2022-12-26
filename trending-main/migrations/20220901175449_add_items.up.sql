CREATE TABLE items
(
    id                BIGSERIAL,
    collection_symbol VARCHAR(500) NOT NULL,
    token_mint        VARCHAR(500) NOT NULL,
    title             TEXT,
    img               TEXT,
    rank              INT,
    attributes        JSONB
);
