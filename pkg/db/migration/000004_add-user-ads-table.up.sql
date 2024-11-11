-- Table for storing user and ad relationships
create table if not exists user_ads
(
    user_id varchar(31) references "user" (tg_id) on delete cascade,
    ad_id   bigint references ad (id) on delete cascade
);