package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

type Song struct {
	Id int `json:"id"`
	Artist string `json:"artist"`
	Song string `json:"song"`
	Genre int `json:"genre"`
	Length int `json:"length"`
}

//func (s Song) Scan(src interface{})error{
//	switch t := src.(type) {
//		case[]byte:
//			return json.Unmarshal(t,&s)
//	default:
//		return errors.New("Invalid type")
//	}
//}

func getSongs(w http.ResponseWriter,r *http.Request) {
	db, _ := sql.Open("sqlite3","./db/jrdd.db")

	//var artist string
	var song Song
	rows, _ := db.Query("Select id,artist,song,genre,length from Songs")

	for rows.Next() {
		rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre,&song.Length)
		jsonOut, _ := json.Marshal(song)

		fmt.Fprintln(w,string(jsonOut))
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

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/songs"), getSongs)

	http.ListenAndServe("localhost:8000", mux)
}
