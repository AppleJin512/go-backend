ALTER TABLE activities
    ADD COLUMN seller VARCHAR(64) DEFAULT NULL;
ALTER TABLE activities
    ADD COLUMN auction_house_address VARCHAR(64) DEFAULT NULL;
ALTER TABLE activities
    ADD COLUMN seller_referral VARCHAR(64) DEFAULT NULL;
