-- +goose Up
-- +goose StatementBegin

-- The slug is a stable, language independent identifier for a type. The
-- frontend uses it to look up a translated name and falls back to the name
-- stored here if it does not know the slug.
ALTER TABLE types ADD COLUMN slug text NOT NULL DEFAULT '';

UPDATE types SET slug = 'misc' WHERE id = 'c29ebd85-812e-4cf6-bfc7-c8368eb83334';
UPDATE types SET slug = 'fruit-vegetables-nuts' WHERE id = 'fe0b085b-2df9-4422-a7cb-7867947719a5';
UPDATE types SET slug = 'canned-goods' WHERE id = '0c9b99fb-b2c8-41e4-8afa-b8cca3ac2ca1';
UPDATE types SET slug = 'sauces-spices-dressings' WHERE id = '7c693d05-4939-44e6-845d-57951720e886';
UPDATE types SET slug = 'drinks-alcohol' WHERE id = '0828b46f-98c9-41ea-9918-164751782861';
UPDATE types SET slug = 'bakery' WHERE id = 'e693272f-4a40-4c0e-9e38-8ebb33004271';
UPDATE types SET slug = 'spreads' WHERE id = 'ab8328c2-29e2-4767-a6fb-27d8e11dc8df';
UPDATE types SET slug = 'coffee-tea' WHERE id = '21b7a2d6-0507-41dc-9a41-4f8a3c86564a';
UPDATE types SET slug = 'cereals-muesli' WHERE id = '97ef6e7e-6c1a-47bc-9d34-35e26a1a0d5c';
UPDATE types SET slug = 'pasta-rice' WHERE id = '1a78a64a-ff86-49db-b64d-45a8b2e76c25';
UPDATE types SET slug = 'cooking-baking' WHERE id = '5d6b6b67-34f3-4a48-bb63-cf65f0f2219d';
UPDATE types SET slug = 'meat-fish' WHERE id = 'd67bd9ce-56f1-4227-885b-0656f74edb22';
UPDATE types SET slug = 'frozen' WHERE id = 'b98f7846-a4cd-4b00-86bf-a6714e982469';
UPDATE types SET slug = 'dairy-chilled' WHERE id = '36298b3b-fcd5-4189-b34f-dae3dea08412';
UPDATE types SET slug = 'sweets-snacks' WHERE id = '13f6bd3e-aeeb-4890-955f-fd91c2450a7e';
UPDATE types SET slug = 'household-baby-pets' WHERE id = 'a14bca10-13b7-4a9c-a663-75a5203c3f09';
UPDATE types SET slug = 'ready-meals-broth-sauce' WHERE id = '3de4b7ac-60be-4d65-8dbf-431f2c6d1270';

-- +goose StatementEnd
