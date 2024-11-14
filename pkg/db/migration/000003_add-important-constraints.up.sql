-- Make user_id and ad_id NOT NULL and add unique constraint in favorite_ads table
alter table favorite_ads
    alter column user_id set not null,
    alter column ad_id set not null,
    add constraint unique_user_ad_favorite unique (user_id, ad_id);

-- Make user_id and ad_id NOT NULL and add unique constraint in user_ads table
alter table user_ads
    alter column user_id set not null,
    alter column ad_id set not null,
    add constraint unique_user_ad_user_ads unique (user_id, ad_id);