-- migrate:up
pragma foreign_keys = on;

create table if not exists users(
  id integer primary key autoincrement,
  email text not null unique,
  created_at datetime default current_timestamp,
  updated_at datetime default current_timestamp
);

create table if not exists sessions(
  user_id integer not null references users(id) on delete cascade,
  token text not null,
  created_at datetime default current_timestamp,
  primary key(user_id)
);

create table if not exists verification_codes(
  code text default (lower(hex(randomblob(16)))),
  user_id integer not null references users(id) on delete cascade,
  created_at datetime default current_timestamp,
  primary key(user_id)
);

create table if not exists monitors(
  id integer primary key autoincrement,
  name text not null,
  url text not null,
  protocol text default 'https',
  interval integer default 300, -- seconds
  user_id integer not null references users(id) on delete cascade,
  created_at datetime default current_timestamp,
  updated_at datetime default current_timestamp
);

create table if not exists monitor_data(
  id integer primary key autoincrement,
  status_code integer not null,
  response_time float not null,
  monitor_id integer not null references monitors(id) on delete cascade,
  created_at datetime default current_timestamp
);

-- migrate:down
drop table if exists monitor_data;
drop table if exists monitors;
drop table if exists verification_codes;
drop table if exists sessions;
drop table if exists users;
