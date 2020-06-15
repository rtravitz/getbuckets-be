package handler

import (
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
