package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"../customError"
	"../dbase"
	"../models"
	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	user.Nickname = nickname
	if customErr := user.CreateUser(dbase.DB); customErr == customError.ConflictNickname {
		users, customErr := user.GetUsers(dbase.DB)
		if customErr != customError.OK {
			sendError(customErr, w, nickname)
		}

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(users)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func ProfileUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := models.User{}
	user.Nickname = nickname

	if r.Method == "GET" {
		if customErr := user.GetProfile(dbase.DB); customErr != customError.OK {
			sendError(customErr, w, nickname)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&user)
		return
	}

	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&user)
		defer r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		if user.Email == "" && user.Fullname == "" && user.About == "" {
			if customErr := user.GetProfile(dbase.DB); customErr != customError.OK {
				sendError(customErr, w, nickname)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&user)
			return
		}

		if customErr := user.UpdateProfile(dbase.DB); customErr != customError.OK {
			sendError(customErr, w, nickname)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&user)
		return
	}
}

func sendError(err customError.ErrorType, w http.ResponseWriter, identifier string) {
	errorBody := customError.GetErrorInfo(err)
	w.WriteHeader(errorBody.Status)
	msg := errorBody.Msg + identifier
	errJSON := map[string]string{"message": msg}
	json.NewEncoder(w).Encode(errJSON)
}
