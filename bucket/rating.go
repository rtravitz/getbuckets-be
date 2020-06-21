package bucket

import (
	"time"

	"github.com/jmoiron/sqlx"
)

//Rating is information from a user about an invidual bucket
type CleanRating struct {
	ID        int       `db:"id" json:"id"`
	Score     int       `db:"score" json:"score"`
	Locked    bool      `db:"locked" json:"locked"`
	BucketID  int       `db:"bucket_id" json:"bucket_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

//Rating is information from a user about an invidual bucket
type LockRating struct {
	ID        int       `db:"id" json:"id"`
	Locked    bool      `db:"locked" json:"locked"`
	BucketID  int       `db:"bucket_id" json:"bucket_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

//AvgRating are the averages of user responses for a bucket
type AvgRating struct {
	Cleanliness   float64 `json:"cleanliness"`
	LockedPercent float64 `json:"locked_percent"`
	LockRatings   int     `json:"lock_ratings"`
	CleanRatings  int     `json:"clean_ratings"`
}

func SaveCleanlinessRating(db *sqlx.DB, r CleanRating) (CleanRating, error) {
	const q = `
		INSERT INTO clean_ratings (score, bucket_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(q, r.Score, r.BucketID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}

	return r, nil
}

func SaveLockedRating(db *sqlx.DB, r LockRating) (LockRating, error) {
	const q = `
		INSERT INTO lock_ratings (locked, bucket_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(q, r.Locked, r.BucketID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}

	return r, nil
}
