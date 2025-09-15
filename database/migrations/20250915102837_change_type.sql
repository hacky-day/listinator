-- +goose Up
-- +goose StatementBegin

ALTER TABLE types RENAME TO types_old;

CREATE TABLE types (
  id text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime,
  name text NOT NULL,
  immutable numeric DEFAULT 0,
  color text NOT NULL,
  priority integer,
  PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_types_deleted_at ON types(deleted_at);

-- initial values
INSERT INTO types (id, created_at, updated_at, deleted_at, name, immutable, color, priority)
VALUES
    ('c29ebd85-812e-4cf6-bfc7-c8368eb83334', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '📦 Miscellaneous', 1, 'gray', 0),
    ('fe0b085b-2df9-4422-a7cb-7867947719a5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍎 Fruit & Vegetables & Nuts', 0, 'crimson', 10),
    ('0c9b99fb-b2c8-41e4-8afa-b8cca3ac2ca1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🥫 Canned Goods', 0, 'silver', 20),
    ('7c693d05-4939-44e6-845d-57951720e886', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🌶️ Sauces & Spices & Dressings', 0, 'darkred', 30),
    ('0828b46f-98c9-41ea-9918-164751782861', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍹 Drinks & Alcohol', 0, 'orange', 40),
    ('e693272f-4a40-4c0e-9e38-8ebb33004271', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🥖 Bakery', 0, 'saddlebrown', 50),
    ('ab8328c2-29e2-4767-a6fb-27d8e11dc8df', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🥜 Spreads', 0, 'peru', 60),
    ('21b7a2d6-0507-41dc-9a41-4f8a3c86564a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '☕ Coffee & Tea', 0, 'brown', 70),
    ('97ef6e7e-6c1a-47bc-9d34-35e26a1a0d5c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🥣 Cereals & Muesli', 0, 'goldenrod', 80),
    ('1a78a64a-ff86-49db-b64d-45a8b2e76c25', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍝 Pasta & Rice', 0, 'khaki', 90),
    ('5d6b6b67-34f3-4a48-bb63-cf65f0f2219d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍳 Cooking & Baking', 0, 'coral', 100),
    ('d67bd9ce-56f1-4227-885b-0656f74edb22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍖 Meat & Fish', 0, 'red', 110),
    ('b98f7846-a4cd-4b00-86bf-a6714e982469', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '❄️ Frozen', 0, 'lightblue', 120),
    ('36298b3b-fcd5-4189-b34f-dae3dea08412', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🧀 Dairy & Chilled', 0, 'gold', 130),
    ('13f6bd3e-aeeb-4890-955f-fd91c2450a7e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍬 Sweets & Snacks', 0, 'yellow', 140),
    ('a14bca10-13b7-4a9c-a663-75a5203c3f09', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🏠 Household & Baby & Pets', 0, 'teal', 150),
    ('3de4b7ac-60be-4d65-8dbf-431f2c6d1270', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '🍲 Ready Meals & Broth & Sauce', 0, 'olive', 170);

ALTER TABLE entries RENAME TO entries_old;
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
  CONSTRAINT fk_entries_type FOREIGN KEY (type_id) REFERENCES types(id),
  CONSTRAINT fk_lists_entries FOREIGN KEY (list_id) REFERENCES lists(id)
);
CREATE INDEX IF NOT EXISTS idx_entries_deleted_at ON entries(deleted_at);

-- Fruit + Vegetables → Fruit & Vegetables & Nuts
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, 'fe0b085b-2df9-4422-a7cb-7867947719a5', e.list_id
FROM entries_old e WHERE type_id IN ('fruit', 'vegetable');

-- Drinks → Drinks & Alcohol
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, '0828b46f-98c9-41ea-9918-164751782861', e.list_id
FROM entries_old e WHERE type_id = 'drink';

-- Meat → Meat & Fish
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, 'd67bd9ce-56f1-4227-885b-0656f74edb22', e.list_id
FROM entries_old e WHERE type_id = 'meat';

-- Snacks → Sweets & Snacks
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, '13f6bd3e-aeeb-4890-955f-fd91c2450a7e', e.list_id
FROM entries_old e WHERE type_id = 'snack';

-- Dairy → Dairy & Chilled
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, '36298b3b-fcd5-4189-b34f-dae3dea08412', e.list_id
FROM entries_old e WHERE type_id = 'dairy';

-- Bread → Bakery
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, 'e693272f-4a40-4c0e-9e38-8ebb33004271', e.list_id
FROM entries_old e WHERE type_id = 'bread';

-- Condiment + Spice → Sauces & Spices & Dressings
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, '7c693d05-4939-44e6-845d-57951720e886', e.list_id
FROM entries_old e WHERE type_id IN ('condiment', 'spice');

-- Frozen → Frozen
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, 'b98f7846-a4cd-4b00-86bf-a6714e982469', e.list_id
FROM entries_old e WHERE type_id = 'frozen';

-- Canned → Canned Goods
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, '0c9b99fb-b2c8-41e4-8afa-b8cca3ac2ca1', e.list_id
FROM entries_old e WHERE type_id = 'canned';

-- Unknown → Miscellaneous
INSERT INTO entries (id, created_at, updated_at, deleted_at, name, number, bought, type_id, list_id)
SELECT id, created_at, updated_at, deleted_at, name, number, bought, 'c29ebd85-812e-4cf6-bfc7-c8368eb83334', e.list_id
FROM entries_old e WHERE type_id = 'unknown';

DROP TABLE entries_old;
DROP TABLE types_old;
-- +goose StatementEnd
