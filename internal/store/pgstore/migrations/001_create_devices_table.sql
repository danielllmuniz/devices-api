-- Write your migrate up statements here
CREATE TYPE device_state AS ENUM ('available', 'in-use', 'inactive');
CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    state device_state NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
---- create above / drop below ----
DROP TABLE IF EXISTS devices;
DROP TYPE IF EXISTS device_state;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
