-- Add up migration script here
create table categories (
  id integer not null primary key autoincrement,
  name text not null unique, -- name of the category
  url text not null unique, -- url of the category
  is_completed integer not null default 0 -- is completed
);

create table business_in_category (
  id integer not null primary key autoincrement,
  category_id integer not null references categories(id), -- the category of the business
  name text not null, -- name of the business
  url text not null unique, -- url of the business
  page integer not null default 1 -- page index where the business is located
);

create table business_info (
  id integer not null primary key autoincrement,
  trade_name text unique, -- business trade name
  business_name text unique, -- business name
  address text unique, -- business address
  phone_number text unique, -- business phone number
  email text unique, -- business email
  website text unique, -- business
  social_media text unique, -- business social media
  canonical_link text not null references business_in_category(url), -- canonical link
  rating real, -- business rating
  description text -- business description
);