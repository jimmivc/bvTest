package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

type Song struct {
	Id 		int `json:"id"`
	Artist 	string `json:"artist"`
	Song 	string `json:"song"`
	Genre 	*Genre `json:"genre"`
	Length 	int `json:"length"`
}

type Genre struct{
	Id		int `json:"id"`
	Name	string `json:"name"`
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/songs"), getAllSongs)
	mux.HandleFunc(pat.Get("/songs/artist/:artist"), getAllSongs)
	mux.HandleFunc(pat.Get("/songs/name/:name"), getAllSongs)
	mux.HandleFunc(pat.Get("/songs/genre/:genre"), getAllSongs)
	mux.HandleFunc(pat.Get("/genre/:name"), getAllSongs)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}


func getAllSongs(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	db, _ := sql.Open("sqlite3","./db/jrdd.db")
	//var artist string
	song := newSong()
	var songs []Song

	rows, _ := db.Query("Select id,artist,song,genre,length from Songs")

	for rows.Next() {

		rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Length)
		songs= append(songs, *song)
	}
	json.NewEncoder(w).Encode(songs)

}
func newSong() *Song{
	return &Song{
		0,
		"",
		"",
		&Genre{},
		0,
	}
}
func getSongByArtist(w http.ResponseWriter,r *http.Request){
	//artist := pat.Param(r,"artist")

}


func getSongsByName(w http.ResponseWriter,r *http.Request)  {
	//w.Header().Set("Content-Type","application/json")
	//db, _ := sql.Open("sqlite3","./db/jrdd.db")
	//
	////var artist string
	//var song Song
	//rows, _ := db.Query("Select id,artist,song,genre,length from Songs")
	//
	//for rows.Next() {
	//
	//	rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre,&song.Length)
	//	jsonOut, _ := json.Marshal(song)
	//
	//	fmt.Fprintln(w,string(jsonOut))
	//}
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
