-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS types (
  name text,
  icon text,
  color text,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS lists (
  id text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime,
  PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_lists_deleted_at ON lists(deleted_at);

CREATE TABLE IF NOT EXISTS entries (
  id text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime,
  name text,
  number text,
  bought numeric,
  type_id text,
  list_id text,
  PRIMARY KEY (id),
  CONSTRAINT fk_entries_type FOREIGN KEY (type_id) REFERENCES types(name),
  CONSTRAINT fk_lists_entries FOREIGN KEY (list_id) REFERENCES lists(id)
);
CREATE INDEX IF NOT EXISTS idx_entries_deleted_at ON entries(deleted_at);

CREATE TABLE IF NOT EXISTS users (
  id text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime,
  name text NOT NULL,
  password_hash text,
  is_admin numeric,
  PRIMARY KEY (id),
  CONSTRAINT uni_users_name UNIQUE (name)
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- initial values
INSERT INTO types (name, icon, color) VALUES ('fruit', '🍎', 'crimson') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('vegetable', '🥦', 'green') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('drink', '🍹', 'orange') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('meat', '🍖', 'red') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('snack', '🍿', 'yellow') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('dairy', '🧀', 'gold') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('bread', '🥖', 'saddlebrown') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('condiment', '🧂', 'gray') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('frozen', '❄️', 'lightblue') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('canned', '🥫', 'silver') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('spice', '🌶️', 'darkred') ON CONFLICT(name) DO NOTHING;
INSERT INTO types (name, icon, color) VALUES ('unknown', '🤷‍♀️', 'black') ON CONFLICT(name) DO NOTHING;
-- +goose StatementEnd
