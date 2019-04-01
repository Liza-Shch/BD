package models

import (
	"database/sql"

	"../customError"
)

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

func (user *User) CreateUser(db *sql.DB) customError.ErrorType {
	_, err := db.Exec(`INSERT INTO "user"(nickname, fullname, about, email) VALUES($1, $2, $3, $4);`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		return customError.ConflictNickname
	}

	return customError.OK
}

func (user *User) GetUsers(db *sql.DB) ([]User, customError.ErrorType) {
	rows, err := db.Query(`SELECT "nickname", "fullname", "about", "email" FROM "user" WHERE "nickname" = $1 or "email" = $2;`,
		user.Nickname, user.Email)

	if err != nil {
		return nil, customError.NotFound
	}

	defer rows.Close()

	users := []User{}
	buf := User{}
	for rows.Next() {
		err := rows.Scan(&buf.Nickname, &buf.Fullname, &buf.About, &buf.Email)
		if err != nil {
			return nil, customError.NotFound
		}
		users = append(users, buf)
	}

	return users, customError.OK
}

func (user *User) GetProfile(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT "nickname", "fullname", "about", "email" FROM "user" WHERE "nickname" = $1;`,
		user.Nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return customError.NotFound
	}

	return customError.OK
}

func (user *User) UpdateProfile(db *sql.DB) customError.ErrorType {
	result, err := db.Exec(`UPDATE "user" SET `+
		`"fullname" = (CASE WHEN $1 = '' THEN "fullname" ELSE $1 END),`+
		`"about" = (CASE WHEN $2 = '' THEN "about" ELSE $2 END),`+
		`"email" = (CASE WHEN $3 = '' THEN "email" ELSE $3 END) `+
		`WHERE "nickname" = $4;`,
		&user.Fullname, &user.About, &user.Email, &user.Nickname)

	if err != nil {
		return customError.ConflictEmail
	}

	if rowCount, _ := result.RowsAffected(); rowCount == 0 {
		return customError.NotFound
	}

	customErr := user.GetProfile(db)
	return customErr
}

func (user *User) GetEmail(db *sql.DB) (string, customError.ErrorType) {
	nickname := ""
	err := db.QueryRow(`SELECT "nickname" FROM "user" WHERE "email" = $1;`,
		user.Email).Scan(&nickname)

	if err != nil {
		return "", customError.NotFound
	}

	return nickname, customError.OK
}
