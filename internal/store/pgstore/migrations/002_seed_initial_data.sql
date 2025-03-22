-- Write your migrate up statements here
INSERT INTO devices (name, brand, state) VALUES 
('iPhone 12', 'Apple', 'available'), 
('Galaxy S21', 'Samsung', 'available'),
('Pixel 5', 'Google', 'inactive'),
('iPhone 11', 'Apple', 'inactive'), 
('Galaxy S20', 'Samsung', 'in-use'),
('Pixel 4', 'Google', 'in-use'),
('iPhone 10', 'Apple', 'in-use'), 
('Galaxy S10', 'Samsung', 'available'),
('Pixel 3', 'Google', 'available');



---- create above / drop below ----
DELETE FROM devices WHERE name IN ('iPhone 12', 'Galaxy S21', 'Pixel 5', 'iPhone 11', 'Galaxy S20', 'Pixel 4', 'iPhone 10', 'Galaxy S10', 'Pixel 3');
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
