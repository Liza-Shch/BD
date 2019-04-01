package models

import (
	"database/sql"
	"log"
	"strconv"

	"../customError"
)

type Forum struct {
	Title        string `json:"title"`
	User         string `json:"user"`
	Slug         string `json:"slug"`
	PostsCount   int    `json:"posts"`
	ThreadsCount int    `json:"threads"`
}

func (forum *Forum) CreateForum(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT "nickname" FROM "user" WHERE "nickname" = $1;`,
		forum.User).Scan(&forum.User)

	if err != nil {
		return customError.NotFound
	}

	_, err = db.Exec(`INSERT INTO "forum"(title, author, slug) VALUES($1, $2, $3);`,
		forum.Title, forum.User, forum.Slug)

	if err != nil {
		return customError.ConflictSlug
	}

	return customError.OK
}

func (forum *Forum) GetForum(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT "slug", "title", "author", "posts", "threads" FROM "forum" WHERE "slug" = $1;`,
		forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.PostsCount, &forum.ThreadsCount)

	if err != nil {
		return customError.ForumNotFound
	}

	return customError.OK
}

func (forum *Forum) GetThreads(db *sql.DB, desc bool, limit int, since string) ([]Thread, customError.ErrorType) {
	requestSQL := `
	SELECT "tid", "title", "author", "message", "forumSlug", "slug", "created"
	FROM "thread"
	WHERE "forumSlug" = $1 `

	descSQL := " ASC "
	orderSign := " >= "

	if desc {
		descSQL = " DESC "
		orderSign = " <= "
	}

	if since != "" {
		sinceSQL := ` and "created" ` + orderSign + "'" + since + "'"
		requestSQL += sinceSQL
	}

	orderSQL := ` ORDER BY "created" `

	orderSQL += descSQL

	requestSQL += orderSQL

	limitSQL := " LIMIT " + strconv.Itoa(limit) + ";"

	requestSQL += limitSQL

	rows, _ := db.Query(requestSQL, forum.Slug)

	defer rows.Close()

	threads := []Thread{}
	buf := Thread{}

	for rows.Next() {
		err := rows.Scan(&buf.ID, &buf.Title, &buf.User, &buf.Message, &buf.ForumSlug, &buf.Slug, &buf.Created)
		if err != nil {
			log.Fatal(err)
			return nil, customError.ThreadNotFound
		}
		threads = append(threads, buf)
	}

	return threads, customError.OK
}

func (forum *Forum) GetUsers(db *sql.DB, desc bool, limit int, since string) ([]User, customError.ErrorType) {
	requestSQL := `
			select distinct "nickname", "fullname", "email", "about" 
			from (
				select "author"
				from "thread"
				where "forumSlug" = $1
			union
				select "author"
				from "post"
				where "forumSlug" = $1
			) usersForum
			join "user" u on usersForum.author = u.nickname
		`

	descSQL := " ASC "
	orderSign := " > "

	if desc {
		descSQL = " DESC "
		orderSign = " < "
	}

	if since != "" {
		sinceSQL := ` and "nickname" ` + orderSign + "'" + since + "'"
		requestSQL += sinceSQL
	}

	orderSQL := ` ORDER BY "nickname" `

	orderSQL += descSQL

	requestSQL += orderSQL

	limitSQL := " LIMIT " + strconv.Itoa(limit)

	requestSQL += limitSQL
	requestSQL += " ;"

	rows, _ := db.Query(requestSQL, forum.Slug)

	defer rows.Close()

	users := []User{}
	buf := User{}

	for rows.Next() {
		err := rows.Scan(&buf.Nickname, &buf.Fullname, &buf.Email, &buf.About)
		if err != nil {
			log.Fatal(err)
			return nil, customError.NotFound
		}
		users = append(users, buf)
	}

	return users, customError.OK
}
