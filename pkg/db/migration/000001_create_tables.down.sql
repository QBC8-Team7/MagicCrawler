-- Drop tables if they already exist to avoid conflicts
drop table if exists favorite_ads, ad_picture, price, ad, publisher, "user", user_ads cascade;
drop type if exists user_role, ad_category, house_type cascade;
drop index if exists idx_ad_publisher, idx_price_ad, idx_favorite_ads_user cascade;
