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

func CreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	forum := models.Forum{}

	err := json.NewDecoder(r.Body).Decode(&forum)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	customErr := forum.CreateForum(dbase.DB)
	switch customErr {
	case customError.NotFound:
		sendError(customErr, w, forum.User)
		return
	case customError.ConflictSlug:
		customErr = forum.GetForum(dbase.DB)
		if customErr != customError.OK {
			sendError(customErr, w, forum.Slug)
			return
		}

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(forum)
	default:
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(forum)
	}
}

func GetForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug := vars["slug"]
	forum := models.Forum{}
	forum.Slug = slug

	customErr := forum.GetForum(dbase.DB)
	if customErr != customError.OK {
		sendError(customErr, w, forum.Slug)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(forum)
}

func GetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug := vars["slug"]
	descStr := r.URL.Query().Get("desc")
	desc, _ := strconv.ParseBool(descStr)
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	since := r.URL.Query().Get("since")

	forum := models.Forum{}
	forum.Slug = slug
	customErr := forum.GetForum(dbase.DB)
	if customErr != customError.OK {
		sendError(customErr, w, forum.Slug)
		return
	}

	threads, customErr := forum.GetThreads(dbase.DB, desc, limit, since)
	if customErr != customError.OK {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(threads)
}

func GetUsersByForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug := vars["slug"]
	descStr := r.URL.Query().Get("desc")
	desc, _ := strconv.ParseBool(descStr)
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if len(limitStr) != 0 {
		limit, _ = strconv.Atoi(limitStr)
	}
	since := r.URL.Query().Get("since")

	forum := models.Forum{}
	forum.Slug = slug
	customErr := forum.GetForum(dbase.DB)
	if customErr != customError.OK {
		sendError(customErr, w, forum.Slug)
		return
	}

	users, customErr := forum.GetUsers(dbase.DB, desc, limit, since)
	if customErr != customError.OK {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
