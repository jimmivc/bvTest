package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

func main() {
	database, _ := sql.Open("sqlite3","./db/jrdd.db")

}
