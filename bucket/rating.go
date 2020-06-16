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
	BucketID      int     `json:"bucket_id"`
	Cleanliness   float64 `json:"cleanliness"`
	LockedPercent float64 `json:"locked_percent"`
}

//SaveRating persists a rating for a bucket
func SaveRating(db *sqlx.DB, r Rating) (Rating, error) {
	const q = `
		INSERT INTO ratings (cleanliness, locked, bucket_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, update_at
	`

	err := db.QueryRow(q, r.Cleanliness, r.Locked, r.BucketID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}

	return r, nil
}

//GetAverageRating returns the averages of ratings for a bucket
func GetAverageRating(db *sqlx.DB, bucketID int) (AvgRating, error) {
	const q = `
		SELECT * FROM ratings
		WHERE bucket_id = $1
	`
	var ratings []Rating
	err := db.Select(&ratings, q, bucketID)
	if err != nil {
		return AvgRating{}, err
	}

	return calcAvgRatings(ratings), nil
}

func calcAvgRatings(ratings []Rating) AvgRating {
	var sum int
	var locked int
	totalRatings := float64(len(ratings))

	for _, r := range ratings {
		sum += r.Cleanliness
		if r.Locked == true {
			locked++
		}
	}

	return AvgRating{
		LockedPercent: (float64(locked) / totalRatings) * 100,
		Cleanliness:   float64(sum) / totalRatings,
	}
}
