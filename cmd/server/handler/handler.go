package handler

import (
	"encoding/json"
	"errors"
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

		respond(w, b, http.StatusCreated)
	}
}

type RatingReq struct {
	Locked      *bool `json:"locked"`
	Cleanliness *int  `json:"cleanliness"`
}

func SaveRatingHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ratingReq RatingReq
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

		rating := bucket.Rating{BucketID: bucketID}

		var savedRating bucket.Rating
		if ratingReq.Cleanliness != nil && ratingReq.Locked != nil {
			rating.Cleanliness = *ratingReq.Cleanliness
			rating.Locked = *ratingReq.Locked
			savedRating, err = bucket.SaveRating(db, rating)
		} else if ratingReq.Cleanliness != nil {
			rating.Cleanliness = *ratingReq.Cleanliness
			savedRating, err = bucket.SaveCleanlinessRating(db, rating)
		} else if ratingReq.Locked != nil {
			rating.Locked = *ratingReq.Locked
			savedRating, err = bucket.SaveLockedRating(db, rating)
		} else {
			err = errors.New("no valid rating were passed. must provide either cleanliness, locked, or both.")
		}

		if err != nil {
			respondError(w, err)
			return
		}

		respond(w, savedRating, http.StatusCreated)
	}
}
