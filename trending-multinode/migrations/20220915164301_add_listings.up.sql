CREATE TABLE listings
(
    id                    BIGSERIAL,
    signature             VARCHAR(500),
    symbol                VARCHAR(500),
    block_time            TIMESTAMP,
    price                 NUMERIC       DEFAULT 0,
    name                  VARCHAR(255)  DEFAULT ''::CHARACTER VARYING,
    mint_address          VARCHAR(64)   DEFAULT ''::CHARACTER VARYING,
    uri                   VARCHAR(500)  DEFAULT ''::CHARACTER VARYING,
    seller                VARCHAR(64)   DEFAULT NULL::CHARACTER VARYING,
    auction_house_address VARCHAR(64)   DEFAULT NULL::CHARACTER VARYING,
    seller_referral       VARCHAR(64)   DEFAULT NULL::CHARACTER VARYING,
    rank                  INT,
    attributes            JSONB         DEFAULT '[]',
    PRIMARY KEY (id)
);
