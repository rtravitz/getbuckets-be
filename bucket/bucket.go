package bucket

import (
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

//List retrieves a list of existing buckets from the database
func List(db *sqlx.DB) ([]Bucket, error) {
	buckets := []Bucket{}
	const q = `SELECT * FROM buckets`

	if err := db.Select(&buckets, q); err != nil {
		return nil, err
	}

	return buckets, nil
}
