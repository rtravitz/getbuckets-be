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

func ShowBucketHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strID := vars["bucket_id"]
		bucketID, err := strconv.Atoi(strID)
		if err != nil {
			respondError(w, err)
			return
		}

		bucket, err := bucket.Show(db, bucketID)
		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, bucket, http.StatusOK)
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

		rated := bucket.RatedBucket{Bucket: b}

		respond(w, rated, http.StatusCreated)
	}
}

func SaveCleanRatingHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ratingReq bucket.CleanRating
		err := json.NewDecoder(r.Body).Decode(&ratingReq)
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

		ratingReq.BucketID = bucketID

		savedRating, err := bucket.SaveCleanlinessRating(db, ratingReq)
		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, savedRating, http.StatusCreated)
	}
}

func SaveLockRatingHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ratingReq bucket.LockRating
		err := json.NewDecoder(r.Body).Decode(&ratingReq)
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

		ratingReq.BucketID = bucketID

		savedRating, err := bucket.SaveLockedRating(db, ratingReq)
		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, savedRating, http.StatusCreated)
	}
}
