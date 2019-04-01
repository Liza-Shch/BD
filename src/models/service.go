package models

import (
	"database/sql"
	"log"

	"../customError"
)

type Service struct {
	UsersCount   int64 `json:"user"`
	ForumsCount  int64 `json:"forum"`
	ThreadsCount int64 `json:"thread"`
	PostsCount   int64 `json:"post"`
}

func (service *Service) GetInfo(db *sql.DB) customError.ErrorType {
	row, err := db.Query(`SELECT COUNT(*) FROM "user";`)

	if err != nil {
		log.Fatal(err)
		return customError.UnknownError
	}

	row.Next()
	row.Scan(&service.UsersCount)

	row, err = db.Query(`SELECT COUNT(*) FROM "forum";`)

	if err != nil {
		log.Fatal(err)
		return customError.UnknownError
	}

	row.Next()
	row.Scan(&service.ForumsCount)

	row, err = db.Query(`SELECT COUNT(*) FROM "thread";`)

	if err != nil {
		log.Fatal(err)
		return customError.UnknownError
	}

	row.Next()
	row.Scan(&service.ThreadsCount)

	row, err = db.Query(`SELECT COUNT(*) as usersCount FROM "post";`)

	if err != nil {
		log.Fatal(err)
		return customError.UnknownError
	}

	row.Next()
	row.Scan(&service.PostsCount)
	return customError.OK
}

func (service *Service) ClearDB(db *sql.DB) customError.ErrorType {
	_, err := db.Exec(`TRUNCATE "user", "forum", "thread", "post", "vote"`)

	if err != nil {
		log.Fatal(err)
		return customError.UnknownError
	}

	return customError.OK
}
