// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package db

import (
	"context"
	"time"

	"github.com/lib/pq"
)

const createShift = `-- name: CreateShift :one
INSERT INTO shifts (
        uid, work_date, shift_length_hours, user_id
    ) VALUES (
        $1, $2, $3, $4
    ) RETURNING id, created, uid, work_date, shift_length_hours, user_id, updated, deleted
`

type CreateShiftParams struct {
	Uid              string
	WorkDate         time.Time
	ShiftLengthHours float64
	UserID           int32
}

func (q *Queries) CreateShift(ctx context.Context, arg CreateShiftParams) (Shift, error) {
	row := q.db.QueryRowContext(ctx, createShift,
		arg.Uid,
		arg.WorkDate,
		arg.ShiftLengthHours,
		arg.UserID,
	)
	var i Shift
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.WorkDate,
		&i.ShiftLengthHours,
		&i.UserID,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
        uid, first_name, last_name, email, password
    ) VALUES (
        $1, $2, $3, $4, $5
    ) RETURNING id, created, uid, type, first_name, last_name, email, password, updated, deleted
`

type CreateUserParams struct {
	Uid       string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Uid,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.Type,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const deleteShift = `-- name: DeleteShift :one
UPDATE shifts SET deleted = now() WHERE uid = $1 and deleted IS NULL RETURNING id, created, uid, work_date, shift_length_hours, user_id, updated, deleted
`

func (q *Queries) DeleteShift(ctx context.Context, uid string) (Shift, error) {
	row := q.db.QueryRowContext(ctx, deleteShift, uid)
	var i Shift
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.WorkDate,
		&i.ShiftLengthHours,
		&i.UserID,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
UPDATE users SET deleted = now() WHERE email = $1 RETURNING id, created, uid, type, first_name, last_name, email, password, updated, deleted
`

func (q *Queries) DeleteUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, deleteUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.Type,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const getAllUsersWithShifts = `-- name: GetAllUsersWithShifts :many
SELECT users.id, users.created, users.uid, users.email, users."type",users.first_name, users.last_name, sum(shifts.shift_length_hours) as total_hours from users 
Inner Join shifts on shifts.user_id = users.id
WHERE shifts.work_date >= $1 AND shifts.work_date <= $2 AND users.deleted IS NULL AND shifts.deleted IS NULL
group by users.id
Order by total_hours DESC
`

type GetAllUsersWithShiftsParams struct {
	WorkDate   time.Time
	WorkDate_2 time.Time
}

type GetAllUsersWithShiftsRow struct {
	ID         int32
	Created    time.Time
	Uid        string
	Email      string
	Type       string
	FirstName  string
	LastName   string
	TotalHours int64
}

func (q *Queries) GetAllUsersWithShifts(ctx context.Context, arg GetAllUsersWithShiftsParams) ([]GetAllUsersWithShiftsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsersWithShifts, arg.WorkDate, arg.WorkDate_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersWithShiftsRow
	for rows.Next() {
		var i GetAllUsersWithShiftsRow
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Uid,
			&i.Email,
			&i.Type,
			&i.FirstName,
			&i.LastName,
			&i.TotalHours,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getShiftByUid = `-- name: GetShiftByUid :one
SELECT id, created, uid, work_date, shift_length_hours, user_id, updated, deleted FROM shifts WHERE uid = $1 AND deleted IS NULL
`

func (q *Queries) GetShiftByUid(ctx context.Context, uid string) (Shift, error) {
	row := q.db.QueryRowContext(ctx, getShiftByUid, uid)
	var i Shift
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.WorkDate,
		&i.ShiftLengthHours,
		&i.UserID,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
Select id, created, uid, type, first_name, last_name, email, password, updated, deleted from users where email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.Type,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const getUserByUid = `-- name: GetUserByUid :one
Select id, created, uid, type, first_name, last_name, email, password, updated, deleted from users where uid = $1
`

func (q *Queries) GetUserByUid(ctx context.Context, uid string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUid, uid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.Type,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const getUserShifts = `-- name: GetUserShifts :many
SELECT id, created, uid, work_date, shift_length_hours, user_id, updated, deleted FROM shifts WHERE user_id = (SELECT id FROM users WHERE email = $1) AND deleted IS NULL
`

func (q *Queries) GetUserShifts(ctx context.Context, email string) ([]Shift, error) {
	rows, err := q.db.QueryContext(ctx, getUserShifts, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shift
	for rows.Next() {
		var i Shift
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Uid,
			&i.WorkDate,
			&i.ShiftLengthHours,
			&i.UserID,
			&i.Updated,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersShiftsByDateRange = `-- name: GetUsersShiftsByDateRange :many
SELECT id, created, uid, work_date, shift_length_hours, user_id, updated, deleted FROM shifts WHERE user_id IN (SELECT id FROM users WHERE email = ANY($1::varchar[])) AND work_date >= $2 AND work_date <= $3 AND deleted IS NULL ORDER BY work_date DESC
`

type GetUsersShiftsByDateRangeParams struct {
	Column1    []string
	WorkDate   time.Time
	WorkDate_2 time.Time
}

func (q *Queries) GetUsersShiftsByDateRange(ctx context.Context, arg GetUsersShiftsByDateRangeParams) ([]Shift, error) {
	rows, err := q.db.QueryContext(ctx, getUsersShiftsByDateRange, pq.Array(arg.Column1), arg.WorkDate, arg.WorkDate_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shift
	for rows.Next() {
		var i Shift
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Uid,
			&i.WorkDate,
			&i.ShiftLengthHours,
			&i.UserID,
			&i.Updated,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateShiftLength = `-- name: UpdateShiftLength :one
UPDATE shifts SET shift_length_hours = $1 WHERE uid = $2 AND deleted IS NULL RETURNING id, created, uid, work_date, shift_length_hours, user_id, updated, deleted
`

type UpdateShiftLengthParams struct {
	ShiftLengthHours float64
	Uid              string
}

func (q *Queries) UpdateShiftLength(ctx context.Context, arg UpdateShiftLengthParams) (Shift, error) {
	row := q.db.QueryRowContext(ctx, updateShiftLength, arg.ShiftLengthHours, arg.Uid)
	var i Shift
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.WorkDate,
		&i.ShiftLengthHours,
		&i.UserID,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const updateShiftUserId = `-- name: UpdateShiftUserId :one
UPDATE shifts SET user_id = (SELECT id FROM users WHERE email = $1) WHERE shifts.uid = $2 AND deleted IS NULL RETURNING id, created, uid, work_date, shift_length_hours, user_id, updated, deleted
`

type UpdateShiftUserIdParams struct {
	Email string
	Uid   string
}

func (q *Queries) UpdateShiftUserId(ctx context.Context, arg UpdateShiftUserIdParams) (Shift, error) {
	row := q.db.QueryRowContext(ctx, updateShiftUserId, arg.Email, arg.Uid)
	var i Shift
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.WorkDate,
		&i.ShiftLengthHours,
		&i.UserID,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const updateShiftWorkDate = `-- name: UpdateShiftWorkDate :one
UPDATE shifts SET work_date = $1 WHERE uid = $2 AND deleted IS NULL RETURNING id, created, uid, work_date, shift_length_hours, user_id, updated, deleted
`

type UpdateShiftWorkDateParams struct {
	WorkDate time.Time
	Uid      string
}

func (q *Queries) UpdateShiftWorkDate(ctx context.Context, arg UpdateShiftWorkDateParams) (Shift, error) {
	row := q.db.QueryRowContext(ctx, updateShiftWorkDate, arg.WorkDate, arg.Uid)
	var i Shift
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.WorkDate,
		&i.ShiftLengthHours,
		&i.UserID,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}

const updateUserPermissionLevel = `-- name: UpdateUserPermissionLevel :one
UPDATE users SET type=$1 WHERE email=$2 RETURNING id, created, uid, type, first_name, last_name, email, password, updated, deleted
`

type UpdateUserPermissionLevelParams struct {
	Type  string
	Email string
}

func (q *Queries) UpdateUserPermissionLevel(ctx context.Context, arg UpdateUserPermissionLevelParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPermissionLevel, arg.Type, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Uid,
		&i.Type,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Updated,
		&i.Deleted,
	)
	return i, err
}
