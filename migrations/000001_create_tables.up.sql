-- Set timezone to Tehran (UTC+3:30)
alter database "magic-crawler" set timezone to 'Asia/Tehran';

-- Define enums
create type user_role as enum ('super_admin', 'admin', 'simple');
create type ad_category as enum ('rent', 'buy', 'mortgage', 'other');
create type house_type as enum ('apartment', 'villa', 'other');

-- Table for storing publishers
create table if not exists publisher
(
    id   serial primary key,
    name varchar(31) not null,
    url  varchar(63) not null
);

-- Table for storing ads
create table if not exists ad
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
    year             int check (year >= 0),
    floor            int,
    total_floors     int,
    has_warehouse    boolean,
    has_elevator     boolean,
    lat              decimal(9, 6) check (lat between -90 and 90),
    lng              decimal(9, 6) check (lng between -180 and 180)
);

-- Table for storing user information
create table if not exists "user"
(
    tg_id            varchar(31) primary key,
    role             user_role,
    watchlist_period int
);

-- Table for storing prices
create table if not exists price
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
create table if not exists ad_picture
(
    id    bigserial primary key,
    ad_id bigint references ad (id) on delete cascade,
    url   varchar(255)
);

-- Table for storing user's favorite ads
create table if not exists favorite_ads
(
    id      bigserial primary key,
    user_id varchar(31) references "user" (tg_id) on delete cascade,
    ad_id   bigint references ad (id) on delete cascade
);

-- Indexes for optimized searches
create index if not exists idx_ad_publisher on ad (publisher_ad_key);
create index if not exists idx_price_ad on price (ad_id);
create index if not exists idx_favorite_ads_user on favorite_ads (user_id);