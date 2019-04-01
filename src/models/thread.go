package models

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"../customError"
)

type Thread struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	User      string `json:"author"`
	ForumSlug string `json:"forum"`
	Message   string `json:"message"`
	Slug      string `json:"slug"`
	Created   string `json:"created"`
	Votes     int    `json:"votes"`
}

func (thread *Thread) CreateThread(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT "slug" FROM "forum" WHERE "slug" = $1;`,
		thread.ForumSlug).Scan(&thread.ForumSlug)

	if err != nil {
		return customError.ForumNotFound
	}

	err = db.QueryRow(`SELECT "nickname" FROM "user" WHERE "nickname" = $1;`,
		thread.User).Scan(&thread.User)

	if err != nil {
		return customError.NotFound
	}

	if thread.Slug == "" {
		err = db.QueryRow(`INSERT INTO "thread"(title, author, message, "forumSlug") VALUES($1, $2, $3, $4) RETURNING tid;`,
			thread.Title, thread.User, thread.Message, thread.ForumSlug).Scan(&thread.ID)
	} else {
		if thread.Created == "" {
			err = db.QueryRow(`INSERT INTO "thread"(title, author, message, "forumSlug", slug) VALUES($1, $2, $3, $4, $5) RETURNING tid, "slug";`,
				thread.Title, thread.User, thread.Message, thread.ForumSlug, thread.Slug).Scan(&thread.ID, &thread.Slug)
		} else {
			err = db.QueryRow(`INSERT INTO "thread"(title, author, message, "forumSlug", slug, created) VALUES($1, $2, $3, $4, $5, $6) RETURNING tid, "slug";`,
				thread.Title, thread.User, thread.Message, thread.ForumSlug, thread.Slug, thread.Created).Scan(&thread.ID, &thread.Slug)
		}
	}

	if err != nil {
		return customError.ConflictSlug
	}

	db.Exec(`UPDATE "forum" SET "threads" = "threads" + 1 WHERE "slug" = $1;`, thread.ForumSlug)

	return customError.OK
}

func (thread *Thread) GetThread(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT tid, title, author, message, "forumSlug", slug, created, votes FROM "thread" WHERE "slug" = $1 AND "slug" <> '' OR "tid" = $2;`,
		thread.Slug, thread.ID).Scan(&thread.ID, &thread.Title, &thread.User, &thread.Message, &thread.ForumSlug,
		&thread.Slug, &thread.Created, &thread.Votes)

	if err != nil {
		return customError.ThreadNotFound
	}

	return customError.OK
}

func (thread *Thread) UpdateThread(db *sql.DB) customError.ErrorType {
	result, _ := db.Exec(`UPDATE "thread" SET `+
		`"message" = (CASE WHEN $1 = '' THEN "message" ELSE $1 END),`+
		`"title" = (CASE WHEN $2 = '' THEN "title" ELSE $2 END) `+
		`WHERE "slug" = $3 AND "slug" <> '' OR "tid" = $4;`,
		&thread.Message, &thread.Title, &thread.Slug, &thread.ID)

	if rowCount, _ := result.RowsAffected(); rowCount == 0 {
		return customError.ThreadNotFound
	}

	customErr := thread.GetThread(db)
	return customErr
}

func (thread *Thread) VoteThread(db *sql.DB, voice Voice) customError.ErrorType {
	err := db.QueryRow(`SELECT tid, author FROM "thread" WHERE "slug" = $1 AND "slug" <> '' OR "tid" = $2;`,
		thread.Slug, thread.ID).Scan(&voice.Thread, &thread.User)

	if err != nil {
		return customError.ThreadNotFound
	}

	if strings.EqualFold(thread.User, voice.Nickname) {
		return customError.OK
	}

	customErr := voice.vote(db)

	if customErr != customError.OK {
		if customErr == customError.DuplicateVote {
			customErr = thread.GetThread(db)
			return customErr
		}
		return customErr
	}

	result, _ := db.Exec(`UPDATE "thread" SET votes = votes + $1 WHERE "slug" = $2 AND "slug" <> '' OR "tid" = $3;`,
		&voice.Voice, &thread.Slug, &thread.ID)

	if rowCount, _ := result.RowsAffected(); rowCount == 0 {
		return customError.ThreadNotFound
	}

	customErr = thread.GetThread(db)

	return customErr
}

func (thread *Thread) GetPosts(db *sql.DB, desc bool, limit int, since int, sort string) ([]Post, customError.ErrorType) {
	requestSQL := `
		SELECT pid, parent, author, message, "forumSlug", tid, created
		FROM "post"
		WHERE "tid" = $1
	`

	descSQL := " ASC "
	orderSign := " > "

	if desc {
		descSQL = " DESC "
		orderSign = " < "
	}

	var rows *sql.Rows
	//var err error
	switch sort {
	case "", "flat":
		if since != 0 {
			sinceSQL := ` and pid ` + orderSign + "'" + strconv.Itoa(since) + "'"
			requestSQL += sinceSQL
		}

		orderSQL := ` ORDER BY created `
		orderSQL += descSQL

		orderSQL += `, pid `
		orderSQL += descSQL

		requestSQL += orderSQL

		limitSQL := " LIMIT " + strconv.Itoa(limit)

		requestSQL += limitSQL
		requestSQL += " ;"
		//rows, _ = db.Query(requestSQL, thread.ID)
	case "tree":
		if since != 0 {
			sinceSQL := ` and path ` + orderSign + `(SELECT path FROM post WHERE pid = ` + strconv.Itoa(since) + ") "
			requestSQL += sinceSQL
		}

		orderSQL := ` ORDER BY path `
		orderSQL += descSQL

		requestSQL += orderSQL

		limitSQL := " LIMIT " + strconv.Itoa(limit)
		requestSQL += limitSQL
		requestSQL += " ;"
		//rows, _ = db.Query(requestSQL, thread.ID)
	case "parent_tree":
		parentPathSQL := `
			and path[1] IN (
				SELECT parent_post.pid FROM "post" AS parent_post
				WHERE parent_post.tid = $1 AND parent_post.parent = 0
		`
		requestSQL += parentPathSQL

		if since != 0 {
			sinceSQL := ` AND parent_post.path[1] ` + orderSign + `(SELECT path[1] FROM post WHERE pid = ` + strconv.Itoa(since) + ") "
			requestSQL += sinceSQL
		}

		orderByParentSQL := ` ORDER BY parent_post.pid `
		orderByParentSQL += descSQL
		requestSQL += orderByParentSQL

		limitSQL := " LIMIT " + strconv.Itoa(limit)
		requestSQL += limitSQL
		requestSQL += ")"

		orderSQL := ` ORDER BY path[1] `
		orderSQL += descSQL
		orderSQL += `, path;`
		requestSQL += orderSQL
	}

	var err error
	rows, err = db.Query(requestSQL, thread.ID)
	if err != nil {
		log.Fatal(err)
		return nil, customError.ThreadNotFound
	}

	defer rows.Close()

	posts := []Post{}
	buf := Post{}

	for rows.Next() {
		err := rows.Scan(&buf.ID, &buf.Parent, &buf.Author, &buf.Message, &buf.ForumSlug, &buf.ThreadID, &buf.Created)
		if err != nil {
			log.Fatal(err)
			return nil, customError.NotFound
		}
		posts = append(posts, buf)
	}

	return posts, customError.OK
}
