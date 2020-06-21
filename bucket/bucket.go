package bucket

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

//Bucket represents a bucket on the map
type Bucket struct {
	ID        int       `db:"id" json:"id"`
	Lat       float64   `db:"lat" json:"lat"`
	Lng       float64   `db:"lng" json:"lng"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type RatedBucket struct {
	Bucket
	AverageRating AvgRating `json:"average_rating"`
}

func processRatedBucketRow(row *sql.Row, b *RatedBucket) error {
	var cleanliness sql.NullFloat64
	err := row.Scan(&b.ID, &b.Lat, &b.Lng, &b.CreatedAt, &b.UpdatedAt,
		&cleanliness, &b.AverageRating.LockedPercent, &b.AverageRating.LockRatings, &b.AverageRating.CleanRatings)
	if err != nil {
		return err
	}

	if cleanliness.Valid {
		b.AverageRating.Cleanliness = cleanliness.Float64
	}

	return nil
}

func Show(db *sqlx.DB, bucketID int) (RatedBucket, error) {
	var b RatedBucket
	const q = `
    SELECT 
      buckets.id, buckets.lat, buckets.lng, buckets.created_at, buckets.updated_at, 
      AVG(clean_ratings.score) AS cleanliness,
      (((COUNT(*) FILTER (WHERE "locked")) / CAST(COUNT(*) AS DECIMAL)) * 100) AS locked_percent,
      COUNT(distinct lock_ratings.id) as lock_ratings,
      COUNT(distinct clean_ratings.id) as clean_ratings
    FROM buckets
    LEFT JOIN clean_ratings ON buckets.id = clean_ratings.bucket_id
    LEFT JOIN lock_ratings ON buckets.id = lock_ratings.bucket_id
    GROUP BY buckets.id
    HAVING buckets.id = $1;
	`

	row := db.QueryRow(q, bucketID)
	err := processRatedBucketRow(row, &b)
	if err != nil {
		return b, err
	}

	return b, nil
}

//List retrieves a list of existing buckets from the database
func List(db *sqlx.DB) ([]RatedBucket, error) {
	var buckets []RatedBucket
	const q = `
    SELECT 
      buckets.id, buckets.lat, buckets.lng, buckets.created_at, buckets.updated_at, 
      AVG(clean_ratings.score) AS cleanliness,
      (((COUNT(*) FILTER (WHERE "locked")) / CAST(COUNT(*) AS DECIMAL)) * 100) AS locked_percent,
      COUNT(distinct lock_ratings.id) as lock_ratings,
      COUNT(distinct clean_ratings.id) as clean_ratings
    FROM buckets
    LEFT JOIN clean_ratings ON buckets.id = clean_ratings.bucket_id
    LEFT JOIN lock_ratings ON buckets.id = lock_ratings.bucket_id
    GROUP BY buckets.id;
	`

	rows, err := db.Query(q)
	if err != nil {
		return buckets, err
	}

	for rows.Next() {
		var b RatedBucket
		var cleanliness sql.NullFloat64
		err = rows.Scan(&b.ID, &b.Lat, &b.Lng, &b.CreatedAt, &b.UpdatedAt,
			&cleanliness, &b.AverageRating.LockedPercent, &b.AverageRating.LockRatings, &b.AverageRating.CleanRatings)
		if err != nil {
			return buckets, err
		}

		if cleanliness.Valid {
			b.AverageRating.Cleanliness = cleanliness.Float64
		}

		buckets = append(buckets, b)
	}

	return buckets, nil
}

//Save saves a new bucket
func (b *Bucket) Save(db *sqlx.DB) error {
	const q = `INSERT INTO buckets (lat, lng) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := db.QueryRow(q, b.Lat, b.Lng).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
