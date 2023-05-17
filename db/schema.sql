CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE users(
  id integer primary key autoincrement,
  email text not null unique,
  created_at datetime default current_timestamp,
  updated_at datetime default current_timestamp
);
CREATE TABLE sessions(
  user_id integer not null references users(id) on delete cascade,
  token text not null,
  created_at datetime default current_timestamp,
  primary key(user_id)
);
CREATE TABLE verification_codes(
  code text default (lower(hex(randomblob(16)))),
  user_id integer not null references users(id) on delete cascade,
  created_at datetime default current_timestamp,
  primary key(user_id)
);
CREATE TABLE monitors(
  id integer primary key autoincrement,
  name text not null,
  url text not null,
  protocol text default 'https',
  interval integer default 300, -- seconds
  user_id integer not null references users(id) on delete cascade,
  created_at datetime default current_timestamp,
  updated_at datetime default current_timestamp
);
CREATE TABLE monitor_data(
  id integer primary key autoincrement,
  status_code integer not null,
  response_time float not null,
  monitor_id integer not null references monitors(id) on delete cascade,
  created_at datetime default current_timestamp
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20230401182455');
