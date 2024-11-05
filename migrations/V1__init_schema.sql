-- Some queries to check the connection
-- create table "test"(
--     id bigserial primary key,
--     created_at timestamp default now(),
--     name varchar(255)
-- );
--
-- select * from "test";
--
-- insert into test (name) values ('ali'), ('mahdi');
--
-- drop table test;

-- Drop tables if they already exist to avoid conflicts
drop table if exists favorite_ads, ad_picture, price, ad, publisher, "user" cascade;
drop type if exists user_role, ad_category, house_type cascade;
drop index if exists idx_ad_publisher, idx_price_ad, idx_favorite_ads_user cascade;

-- Define enums
create type user_role as enum ('super_admin', 'admin', 'simple');
create type ad_category as enum ('rent', 'buy', 'mortgage');
create type house_type as enum ('apartment', 'villa');

-- Table for storing publishers
create table publisher
(
    id   serial primary key,
    name varchar(31) not null,
    url  varchar(63) not null
);

-- Table for storing ads
create table ad
(
    id               bigserial primary key,
    publisher_ad_key varchar(255) unique not null,
    publisher_id     int                 references publisher (id) on delete set null,
    created_at       timestamp default now(),
    updated_at       timestamp,
    published_at     timestamp,
    category         ad_category,
    author           varchar(63),
    url              varchar(255),
    title            varchar(255),
    description      text,
    city             varchar(63),
    neighborhood     varchar(63),
    house_type       house_type,
    meterage         int check (meterage >= 0),
    rooms_count      int check (rooms_count >= 0),
    age              int check (age >= 0),
    floor            int,
    has_warehouse    boolean,
    has_elevator     boolean,
    lat              decimal(9, 6) check (lat between -90 and 90),
    lng              decimal(9, 6) check (lng between -180 and 180)
);

-- Table for storing user information
create table "user"
(
    tg_id            varchar(31) primary key,
    role             user_role,
    watchlist_period int
);

-- Table for storing prices
create table price
(
    id              serial primary key,
    ad_id           bigint references ad (id) on delete cascade,
    fetched_at      timestamp default now(),
    has_price       boolean,
    total_price     bigint check (total_price >= 0),
    price_per_meter bigint check (price_per_meter >= 0),
    mortgage        bigint check (mortgage >= 0),
    normal_price    bigint check (normal_price >= 0),
    weekend_price   bigint check (weekend_price >= 0)
);

-- Table for storing ad pictures
create table ad_picture
(
    id    bigserial primary key,
    ad_id bigint references ad (id) on delete cascade,
    url   varchar(255)
);

-- Table for storing user's favorite ads
create table favorite_ads
(
    id      bigserial primary key,
    user_id varchar(31) references "user" (tg_id) on delete cascade,
    ad_id   bigint references ad (id) on delete cascade
);

-- Indexes for optimized searches
create index idx_ad_publisher on ad (publisher_ad_key);
create index idx_price_ad on price (ad_id);
create index idx_favorite_ads_user on favorite_ads (user_id);