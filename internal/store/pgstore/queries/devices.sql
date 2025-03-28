-- name: CreateDevice :one
INSERT INTO devices (name, brand, state)
VALUES ($1, $2, $3)
RETURNING id, name, brand, state, created_at;

-- name: UpdateDevice :one
UPDATE devices
SET name = $2,
    brand = $3,
    state = $4
WHERE id = $1
RETURNING id, name, brand, state, created_at;

-- name: PatchDevice :one
UPDATE devices
SET name = COALESCE(NULLIF($2, ''), name),
    brand = COALESCE(NULLIF($3, ''), brand),
    state = CASE WHEN $4 = '' THEN state ELSE $4::device_state END
WHERE id = $1
RETURNING id, name, brand, state, created_at;

-- name: GetDeviceById :one
SELECT id, name, brand, state, created_at
FROM devices
WHERE id = $1;

-- name: GetAllDevices :many
SELECT id, name, brand, state, created_at
FROM devices
ORDER BY created_at DESC;
-- name: GetDevicesByBrand :many
SELECT id, name, brand, state, created_at
FROM devices
WHERE LOWER(brand) = LOWER($1)
ORDER BY created_at DESC;

-- name: GetDevicesByState :many
SELECT id, name, brand, state, created_at
FROM devices
WHERE state = $1
ORDER BY created_at DESC;

-- name: GetDevicesByBrandAndState :many
SELECT id, name, brand, state, created_at
FROM devices
WHERE LOWER(brand) = LOWER($1) AND state = $2
ORDER BY created_at DESC;

-- name: DeleteDevice :one
DELETE FROM devices
WHERE id = $1
RETURNING id;