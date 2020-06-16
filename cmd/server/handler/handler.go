package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func SaveRatingHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rating bucket.Rating
		err := json.NewDecoder(r.Body).Decode(&rating)
		if err != nil {
			respondError(w, err)
			return
		}

		vars := mux.Vars(r)
		strID := vars["bucket_id"]
		bucketID, err := strconv.Atoi(strID)
		if err != nil {
			respondError(w, err)
			return
		}

		b := bucket.Bucket{ID: bucketID}
		err = b.SaveRating(db, rating)
		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, rating, http.StatusCreated)
	}
}
