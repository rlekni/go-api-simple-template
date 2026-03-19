-- name: CreateTeaBlend :one
INSERT INTO tea_blends (
  name,
  description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetTeaBlend :one
SELECT * FROM tea_blends
WHERE id = $1 LIMIT 1;

-- name: ListTeaBlends :many
SELECT * FROM tea_blends
ORDER BY name;

-- name: UpdateTeaBlend :one
UPDATE tea_blends
  set name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTeaBlend :exec
DELETE FROM tea_blends
WHERE id = $1;

-- name: CreateLocation :one
INSERT INTO locations (
  tea_blend_id,
  name,
  quantity
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateLocationQuantity :one
UPDATE locations
  set quantity = $2,
  updated_at = NOW()
WHERE tea_blend_id = $1
RETURNING *;
