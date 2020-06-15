package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rtravitz/getbuckets-be/bucket"
)

func BucketsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buckets, err := bucket.List(db)
		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, buckets, http.StatusOK)
	}
}

func SaveBucketHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b bucket.Bucket
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			respondError(w, err)
			return
		}

		err = b.Save(db)
		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, b, http.StatusCreated)
	}
}
