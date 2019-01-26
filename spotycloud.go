package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Song struct {
	id int `db:"id" json:"id"`
	artist string `db:"artist" json:"artist"`
	song string `db:"song" json:"song"`
	genre int `db:"genre" json:"genre"`
	length int `db:"length" json:"length"`
}

//func (s Song) Scan(src interface{})error{
//	switch t := src.(type) {
//		case[]byte:
//			return json.Unmarshal(t,&s)
//	default:
//		return errors.New("Invalid type")
//	}
//}

func main() {
	db, _ := sql.Open("sqlite3","./db/jrdd.db")

	//var artist string
	var song Song
	rows, _ := db.Query("Select id,artist,song,genre,length from Songs")

	for rows.Next() {
		rows.Scan(&song.id,&song.artist,&song.song,&song.genre,&song.length)
		fmt.Println(song)
	}
}

//func hello(w http.ResponseWriter, r *http.Request) {
	//name := pat.Param(r, "name")
	//fmt.Fprintf(w, "Hello, %s!", name)

	//db, _ := sql.Open("sqlite3","./db/jrdd.db")

	//var artist string
	//
	//rows, _ := db.Query("Select artist from Songs")
	//
	//for rows.Next() {
	//	rows.Scan(&artist)
	//	fmt.Fprintln(w,artist)
	//}
	//err := db.QueryRow("SELECT id,artist,song,genre,length from Songs where id = 2")
	//	fmt.Fprintln(w,err)

//}
//
//func main() {
//	mux := goji.NewMux()
//	mux.HandleFunc(pat.Get("/hello/:name"), hello)
//
//	http.ListenAndServe("localhost:8000", mux)
//}
