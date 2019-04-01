package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../customError"
	"../dbase"
	"../models"
	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	thread := models.Thread{}
	//TODO: вынести это везде в отдельную функцию
	slugOrID := vars["slug_or_id"]
	id, err := strconv.Atoi(slugOrID)

	threadIdentyIsID := true
	if err != nil {
		thread.Slug = slugOrID
		threadIdentyIsID = false
	} else {
		thread.ID = id
	}

	posts := []models.Post{}
	err = json.NewDecoder(r.Body).Decode(&posts)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	customErr := thread.GetThread(dbase.DB)
	if customErr != customError.OK {
		if threadIdentyIsID {
			sendError(customErr, w, strconv.Itoa(thread.ID))
		} else {
			sendError(customErr, w, thread.Slug)
		}
		return
	}

	for i, _ := range posts {
		customErr := posts[i].CreatePost(dbase.DB, &thread)
		switch customErr {
		case customError.NotFound:
			sendError(customErr, w, posts[i].Author)
			return
		case customError.ThreadNotFound:
			if threadIdentyIsID {
				sendError(customErr, w, strconv.Itoa(thread.ID))
			} else {
				sendError(customErr, w, thread.Slug)
			}
			return
		case customError.ConflictPostThread:
			sendError(customErr, w, "")
			return
		case customError.UnknownError:
			sendError(customErr, w, "")
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(posts)
}

func PostDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	post := models.Post{}
	post.ID = id
	if r.Method == "GET" {
		relatedStr := r.URL.Query().Get("related")
		related := strings.Split(relatedStr, ",")

		customErr := post.GetPost(dbase.DB)
		if customErr != customError.OK {
			sendError(customErr, w, idStr)
			return
		}

		answer := map[string]interface{}{}
		for _, relation := range related {
			switch relation {
			case "user":
				user := models.User{}
				user.Nickname = post.Author
				customErr := user.GetProfile(dbase.DB)
				if customErr != customError.OK {
					sendError(customErr, w, user.Nickname)
					return
				}

				answer["author"] = user
			case "thread":
				thread := models.Thread{}
				thread.ID = post.ThreadID
				customErr := thread.GetThread(dbase.DB)
				if customErr != customError.OK {
					sendError(customErr, w, strconv.Itoa(thread.ID))
					return
				}

				answer["thread"] = thread

			case "forum":
				forum := models.Forum{}
				forum.Slug = post.ForumSlug
				customErr := forum.GetForum(dbase.DB)
				if customErr != customError.OK {
					sendError(customErr, w, forum.Slug)
					return
				}

				answer["forum"] = forum
			}
		}

		answer["post"] = post

		answerJSON, _ := json.Marshal(answer)
		w.WriteHeader(http.StatusOK)
		w.Write(answerJSON)
		return
	}

	if r.Method == "POST" {
		errBody := json.NewDecoder(r.Body).Decode(&post)
		defer r.Body.Close()
		if errBody != nil {
			log.Fatal(errBody)
		}

		customErr := post.UpdatePost(dbase.DB)
		if customErr != customError.OK {
			sendError(customErr, w, strconv.Itoa(post.ID))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}
