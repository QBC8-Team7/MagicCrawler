ALTER TABLE ad
    ADD CONSTRAINT ad_publisher_ad_key_key UNIQUE (publisher_ad_key);