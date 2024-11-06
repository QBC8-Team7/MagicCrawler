-- Clear previous data
truncate ad, ad_picture, favorite_ads, price, publisher, "user" cascade;

-- Insert users
insert into "user" (tg_id, role, watchlist_period)
values ('319280055', 'super_admin', 0),
       ('350321889', 'admin', 0),
       ('1280709698', 'simple', 0);

-- Insert publishers
insert into publisher (name, url)
values ('Divar', 'https://divar.ir'),
       ('Sheypoor', 'https://sheypoor.com');

-- Retrieve publisher IDs
with divar_publisher as (select id from publisher where name = 'Divar'),
     sheypoor_publisher as (select id from publisher where name = 'Sheypoor')

-- Insert ads
insert
into ad (publisher_ad_key, publisher_id, updated_at, published_at, category, author, url, title, description, city,
         neighborhood, house_type, meterage, rooms_count, year, floor, total_floors, has_warehouse, has_elevator, lat,
         lng)
values ('wZqMlftL', (select id from divar_publisher), now(), '2024-11-05 18:00:00', 'mortgage', 'آژانس نگین اندرزگو',
        'https://divar.ir/v/wZqMlftL', 'رهن و اجاره آپارتمان ۱۲۰ متر برج باغ حمکت', 'به نام خالق دلها', 'تهران', 'حکمت',
        'apartment', 120, 2, 1396, 4, 11, true, true, 35.805331, 51.440023),
       ('445904922', (select id from sheypoor_publisher), now(), '2024-11-05 15:00:00', 'buy', 'کاربر شیپور',
        'https://www.sheypoor.com/v/445904922', 'ویلا مدرن استخردار', 'ویلا دوبلکس مدرن در بهترین لوکیشن جنگلی',
        'محمودآباد', 'تشبندان', 'villa', 350, 3, 1402, 0, 0, true, false, 0, 0);

-- Insert prices
insert into price (ad_id, fetched_at, has_price, normal_price, mortgage)
values ((select id from ad where publisher_ad_key = 'wZqMlftL'), now(), true, 70000000, 600000000),
       ((select id from ad where publisher_ad_key = 'wZqMlftL'), now() + interval '1 hour', true, 55000000, 650000000);

insert into price (ad_id, fetched_at, has_price, total_price, price_per_meter)
values ((select id from ad where publisher_ad_key = '445904922'), now(), true, 4000000000, 11429000);


-- Insert pictures
insert into ad_picture (ad_id, url)
values ((select id from ad where publisher_ad_key = 'wZqMlftL'),
        'https://s100.divarcdn.com/static/photo/neda/post/hdTM1MQLQ5KBObt0tTQn2g/fa1c7495-5ae8-4747-b631-f1af6c785c33.jpg'),
       ((select id from ad where publisher_ad_key = '445904922'),
        'https://cdn.sheypoor.com/imgs/2024/11/04/445904922/1500x936_Sw/445904922_bc4e875cc2de97f7b03e8b97230d14e7.webp');

-- Insert favorite ads for two users, adding the first ad as their favorite
insert into favorite_ads (user_id, ad_id)
values
    ('319280055', (select id from ad where publisher_ad_key = 'wZqMlftL')),
    ('350321889', (select id from ad where publisher_ad_key = '445904922'));

-- Insert favorite ads for the remaining user, adding both ads as favorites
insert into favorite_ads (user_id, ad_id)
values
    ('1280709698', (select id from ad where publisher_ad_key = 'wZqMlftL')),
    ('1280709698', (select id from ad where publisher_ad_key = '445904922'));