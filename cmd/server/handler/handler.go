package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rtravitz/getbuckets-be/bucket"
)

func processCoords(param string) (bucket.BoundingBox, error) {
	coords := strings.Split(param, ",")
	var processed []float64
	for _, coord := range coords {
		res, err := strconv.ParseFloat(coord, 64)
		if err != nil {
			return bucket.BoundingBox{}, err
		}
		processed = append(processed, res)
	}

	return bucket.BoundingBox{
		SWLng: processed[0],
		SWLat: processed[1],
		NELng: processed[2],
		NELat: processed[3],
	}, nil
}

func BucketsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bboxParam, ok := r.URL.Query()["bbox"]

		if !ok || len(bboxParam[0]) < 1 {
			log.Println("url param 'bbox' is missing")
			return
		}

		bbox, err := processCoords(bboxParam[0])
		if err != nil {
			respondError(w, err)
			return
		}

		buckets, err := bucket.ListInBox(db, bbox)
		if err != nil {
			respondError(w, err)
			return
		}
		if len(buckets) == 0 {
			respond(w, make([]bucket.RatedBucket, 0), http.StatusOK)
		} else {
			respond(w, buckets, http.StatusOK)
		}
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
