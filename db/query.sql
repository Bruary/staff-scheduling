-- name: CreateUser :one
INSERT INTO users (
        uid, first_name, last_name, email, password
    ) VALUES (
        $1, $2, $3, $4, $5
    ) RETURNING *;

-- name: GetUserByEmail :one
Select * from users where email = $1;

-- name: GetUserByUid :one
Select * from users where uid = $1;

-- name: UpdateUserPermissionLevel :one
UPDATE users SET type=$1 WHERE email=$2 RETURNING *;

-- name: DeleteUser :one
UPDATE users SET deleted = now() WHERE email = $1 RETURNING *;

-- name: CreateShift :one
INSERT INTO shifts (
        uid, work_date, shift_length_hours, user_id
    ) VALUES (
        $1, $2, $3, $4
    ) RETURNING *;

-- name: GetUserShifts :many
SELECT * FROM shifts WHERE user_id = (SELECT id FROM users WHERE email = $1) AND deleted IS NULL;

-- name: DeleteShift :one
UPDATE shifts SET deleted = now() WHERE uid = $1 and deleted IS NULL RETURNING *;

-- name: UpdateShiftUserId :one
UPDATE shifts SET user_id = (SELECT id FROM users WHERE email = $1) WHERE shifts.uid = $2 AND deleted IS NULL RETURNING *;

-- name: UpdateShiftWorkDate :one
UPDATE shifts SET work_date = $1 WHERE uid = $2 AND deleted IS NULL RETURNING *;

-- name: UpdateShiftLength :one
UPDATE shifts SET shift_length_hours = $1 WHERE uid = $2 AND deleted IS NULL RETURNING *;

-- name: GetShiftByUid :one
SELECT * FROM shifts WHERE uid = $1 AND deleted IS NULL;