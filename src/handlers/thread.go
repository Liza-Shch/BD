package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../customError"
	"../dbase"
	"../models"
	"github.com/gorilla/mux"
)

func CreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	forumSlug := vars["slug"]
	thread := models.Thread{}

	err := json.NewDecoder(r.Body).Decode(&thread)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	thread.ForumSlug = forumSlug

	customErr := thread.CreateThread(dbase.DB)
	switch customErr {
	case customError.NotFound:
		sendError(customErr, w, thread.User)
	case customError.ForumNotFound:
		sendError(customErr, w, thread.ForumSlug)
	case customError.ConflictSlug:
		getErr := thread.GetThread(dbase.DB)
		if getErr != customError.OK {
			sendError(getErr, w, thread.Slug)
			return
		}

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(thread)
	default:
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(thread)
	}
}

func ThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	thread := models.Thread{}

	slugOrID := vars["slug_or_id"]
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		thread.Slug = slugOrID
	} else {
		thread.ID = id
	}

	if r.Method == "GET" {
		customErr := thread.GetThread(dbase.DB)
		if customErr != customError.OK {
			if err != nil {
				sendError(customErr, w, thread.Slug)
			} else {
				sendError(customErr, w, strconv.Itoa(thread.ID))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(thread)
		return
	}

	if r.Method == "POST" {
		errBody := json.NewDecoder(r.Body).Decode(&thread)
		defer r.Body.Close()
		if errBody != nil {
			log.Fatal(errBody)
		}

		customErr := thread.UpdateThread(dbase.DB)
		if customErr != customError.OK {
			if err != nil {
				sendError(customErr, w, thread.Slug)
			} else {
				sendError(customErr, w, strconv.Itoa(thread.ID))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(thread)
	}
}

func VoteThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	thread := models.Thread{}

	voice := models.Voice{}
	errBody := json.NewDecoder(r.Body).Decode(&voice)
	defer r.Body.Close()
	if errBody != nil {
		log.Fatal(errBody)
	}

	slugOrID := vars["slug_or_id"]
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		thread.Slug = slugOrID
	} else {
		thread.ID = id
	}

	customErr := thread.VoteThread(dbase.DB, voice)

	switch customErr {
	case customError.NotFound:
		sendError(customErr, w, voice.Nickname)
	case customError.ThreadNotFound:
		if err != nil {
			sendError(customErr, w, thread.Slug)
		} else {
			sendError(customErr, w, strconv.Itoa(thread.ID))
		}
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(thread)
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	thread := models.Thread{}

	slugOrID := vars["slug_or_id"]
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		thread.Slug = slugOrID
	} else {
		thread.ID = id
	}

	customErr := thread.GetThread(dbase.DB)
	if customErr != customError.OK {
		if err != nil {
			sendError(customErr, w, thread.Slug)
		} else {
			sendError(customErr, w, strconv.Itoa(thread.ID))
		}
		return
	}

	// TODO зачем переводить в числа

	descStr := r.URL.Query().Get("desc")
	desc, _ := strconv.ParseBool(descStr)
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	sinceStr := r.URL.Query().Get("since")
	since, _ := strconv.Atoi(sinceStr)
	sort := r.URL.Query().Get("sort")

	posts, customErr := thread.GetPosts(dbase.DB, desc, limit, since, sort)

	if customErr != customError.OK {
		sendError(customErr, w, "")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
