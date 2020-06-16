package bucket

import (
	"time"

	"github.com/jmoiron/sqlx"
)

//Rating is information from a user about an invidual bucket
type Rating struct {
	ID          int       `db:"id" json:"id"`
	Cleanliness int       `db:"cleanliness" json:"cleanliness"`
	Locked      bool      `db:"locked" json:"locked"`
	BucketID    int       `db:"bucket_id" json:"bucket_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

//AvgRating are the averages of user responses for a bucket
type AvgRating struct {
	Cleanliness   float64 `json:"cleanliness"`
	LockedPercent float64 `json:"locked_percent"`
	LockRatings   int     `json:"lock_ratings"`
	CleanRatings  int     `json:"clean_ratings"`
}

func SaveCleanlinessRating(db *sqlx.DB, r Rating) (Rating, error) {
	const q = `
		INSERT INTO ratings (cleanliness, bucket_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(q, r.Cleanliness, r.BucketID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}

	return r, nil
}

func SaveLockedRating(db *sqlx.DB, r Rating) (Rating, error) {
	const q = `
		INSERT INTO ratings (locked, bucket_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(q, r.Locked, r.BucketID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}

	return r, nil
}

//SaveRating persists a rating for a bucket
func SaveRating(db *sqlx.DB, r Rating) (Rating, error) {
	const q = `
		INSERT INTO ratings (cleanliness, locked, bucket_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRow(q, r.Cleanliness, r.Locked, r.BucketID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}

	return r, nil
}
