package models

import (
	"database/sql"

	"../customError"
)

type Voice struct {
	ID       int
	Nickname string `json:"nickname"`
	Thread   int
	Voice    int `json:"voice"`
}

func (voice *Voice) vote(db *sql.DB) customError.ErrorType {
	err := db.QueryRow(`SELECT "nickname" FROM "user" WHERE "nickname" = $1;`,
		voice.Nickname).Scan(&voice.Nickname)

	if err != nil {
		return customError.NotFound
	}

	currentVoice := 0
	err = db.QueryRow(`SELECT "voice" FROM "vote" WHERE author = $1 AND tid = $2;`,
		voice.Nickname, voice.Thread).Scan(&currentVoice)

	if currentVoice == voice.Voice {
		return customError.DuplicateVote
	}

	result, _ := db.Exec(`UPDATE "vote" SET voice = $1 WHERE author = $2 AND tid = $3;`,
		voice.Voice, voice.Nickname, voice.Thread)

	if rowCount, _ := result.RowsAffected(); rowCount != 0 {
		voice.Voice *= 2
		return customError.OK
	}

	_, err = db.Exec(`INSERT INTO "vote"(author, tid, voice) VALUES($1, $2, $3);`,
		voice.Nickname, voice.Thread, voice.Voice)

	return customError.OK
}
