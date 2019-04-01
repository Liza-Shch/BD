package customError

import (
	"net/http"
)

type ErrorType string

const (
	OK                 = "ok"
	NotFound           = "not_found"
	ConflictNickname   = "conflict_nickname"
	ConflictEmail      = "conflict_email"
	ConflictSlug       = "conflict_slug"
	ForumNotFound      = "forum_not_found"
	ThreadNotFound     = "thread_not_found"
	PostNotFound       = "post_not_found"
	DuplicateVote      = "duplicate_vote"
	ConflictPostThread = "conflict_post_thread"
	UnknownError       = "unknown_error"
)

type ErrorBody struct {
	Msg    string
	Status int
}

var errorMap = map[ErrorType]ErrorBody{
	NotFound:           {"Can't find user by nickname: ", http.StatusNotFound},
	ConflictEmail:      {"This email is already registered by user: ", http.StatusConflict},
	ConflictSlug:       {"This slug is already exist", http.StatusConflict},
	ForumNotFound:      {"Can't find forum by slug: ", http.StatusNotFound},
	ThreadNotFound:     {"Can't find thread by slug or id: ", http.StatusNotFound},
	PostNotFound:       {"Can't find post by id: ", http.StatusNotFound},
	UnknownError:       {"Oooooops, some error on server((((", http.StatusInternalServerError},
	ConflictPostThread: {"Parent post was created in another thread", http.StatusConflict},
}

func GetErrorInfo(errorType ErrorType) ErrorBody {
	return errorMap[errorType]
}
