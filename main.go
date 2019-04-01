package main

import (
	"log"
	"net/http"

	"./src/dbase"
	"./src/router"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	db, err := dbase.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = dbase.DumpDB(db)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	router.Routing(r)

	http.ListenAndServe(":5000", r)
}
