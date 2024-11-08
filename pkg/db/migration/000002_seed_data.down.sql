-- remove favorite ads
delete from favorite_ads
where user_id in ('319280055', '350321889', '1280709698');

-- remove pictures
delete from ad_picture
where ad_id in (select id from ad where publisher_ad_key in ('wzqmlftl', '445904922'));

-- remove prices
delete from price
where ad_id in (select id from ad where publisher_ad_key in ('wzqmlftl', '445904922'));

-- remove ads
delete from ad
where publisher_ad_key in ('wzqmlftl', '445904922');

-- remove publishers
delete from publisher
where name in ('divar', 'sheypoor');

-- remove users
delete from "user"
where tg_id in ('319280055', '350321889', '1280709698');