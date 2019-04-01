package models

import (
	"database/sql"
	"log"

	"../customError"
)

type Post struct {
	ID        int    `json:"id"`
	Parent    int    `json:"parent"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	IsEdited  bool   `json:"isEdited"`
	ForumSlug string `json:"forum"`
	ThreadID  int    `json:"thread"`
	Created   string `json:"created"`
}

func (post *Post) CreatePost(db *sql.DB, thread *Thread) customError.ErrorType {
	post.ForumSlug = thread.ForumSlug
	post.ThreadID = thread.ID

	// TODO: сделать функцию проверки унивесрально для всех сущностей на наличие валидного автора
	// (может отнаследовать сущности от общего класса расширение в типе структуры, тогда метод будет досутпен)

	err := db.QueryRow(`SELECT "nickname" FROM "user" WHERE "nickname" = $1;`, post.Author).Scan(&thread.User)
	if err != nil {
		return customError.NotFound
	}

	if post.Parent != 0 {
		parentPost := Post{}
		parentPost.ID = post.Parent
		if customErr := parentPost.GetPost(db); customErr != customError.OK {
			return customError.ConflictPostThread
		}

		if parentPost.ThreadID != post.ThreadID {
			return customError.ConflictPostThread
		}
	}

	// TODO: убрать условие, дичь вообще

	if post.Created == "" {
		err = db.QueryRow(`INSERT INTO "post"(author, message, "forumSlug", tid, parent, path) VALUES($1, $2, $3, $4, $5, (SELECT path FROM post WHERE pid = $5) || (select currval('post_pid_seq'))) RETURNING pid, parent;`,
			post.Author, post.Message, post.ForumSlug, post.ThreadID, post.Parent).Scan(&post.ID, &post.Parent)
	} else {
		err = db.QueryRow(`INSERT INTO "post"(author, message, "forumSlug", tid, created, parent, path) VALUES($1, $2, $3, $4, $5, $6, (SELECT path FROM post WHERE pid = $5) || (select currval('post_id_seq'))) RETURNING pid, parent;`,
			post.Author, post.Message, post.ForumSlug, post.ThreadID, post.Created, post.Parent).Scan(&post.ID, &post.Parent)
	}

	if err != nil {
		log.Fatal(err)
		return customError.UnknownError
	}

	db.Exec(`UPDATE "forum" SET "posts" = "posts" + 1 WHERE "slug" = $1;`, post.ForumSlug)

	return customError.OK
}

func (post *Post) GetPost(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT author, message, "forumSlug", tid, created, "isEdited" FROM "post" WHERE "pid" = $1;`,
		post.ID).Scan(&post.Author, &post.Message, &post.ForumSlug, &post.ThreadID, &post.Created, &post.IsEdited)

	if err != nil {
		return customError.PostNotFound
	}

	return customError.OK
}

func (post *Post) UpdatePost(db *sql.DB) customError.ErrorType {
	result, err := db.Exec(`
			UPDATE "post" SET
			"message" = (CASE WHEN $1 = '' or $1 = "message" THEN "message" ELSE $1 END),
		    "isEdited" = (CASE WHEN $1 = '' or $1 = "message" THEN "isEdited" ELSE 'true' END)
		    WHERE "pid" = $2;`,
		post.Message, post.ID)

	if err != nil {
		log.Fatal(err)
	}

	if rowCount, _ := result.RowsAffected(); rowCount == 0 {
		return customError.PostNotFound
	}

	customErr := post.GetPost(db)
	return customErr
}
