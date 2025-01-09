package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/moroz/oauth-tutorial/handlers"
)

const ListenOn = ":3000"

func main() {
	router := handlers.Router()
	log.Printf("Listening on %s", ListenOn)
	log.Fatal(http.ListenAndServe(ListenOn, router))
}
