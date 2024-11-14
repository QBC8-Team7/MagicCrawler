-- Drop unique constraint and remove NOT NULL for user_id and ad_id in favorite_ads table
alter table favorite_ads
    drop constraint if exists unique_user_ad_favorite,
    alter column user_id drop not null,
    alter column ad_id drop not null;

-- Drop unique constraint and remove NOT NULL for user_id and ad_id in user_ads table
alter table user_ads
    drop constraint if exists unique_user_ad_user_ads,
    alter column user_id drop not null,
    alter column ad_id drop not null;