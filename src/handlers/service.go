package handlers

import (
	"encoding/json"
	"net/http"

	"../customError"
	"../dbase"
	"../models"
)

func ServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	service := models.Service{}
	customErr := service.GetInfo(dbase.DB)

	if customErr != customError.OK {
		sendError(customErr, w, "")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}

func ServiceClear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	service := models.Service{}
	customErr := service.ClearDB(dbase.DB)

	if customErr != customError.OK {
		sendError(customErr, w, "")
		return
	}
}
