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


func newSong() *Song{
	return &Song{
		0,
		"",
		"",
		&Genre{},
		0,
	}
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/songs"), getAllSongs)
	mux.HandleFunc(pat.Get("/songs/artist/:artist"), getSongsByArtist)
	mux.HandleFunc(pat.Get("/songs/:name"), getSongsByName)
	mux.HandleFunc(pat.Get("/songs/genre/:genre"), getSongsByGenre)
	mux.HandleFunc(pat.Get("/genre/:name"), getAllSongs)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

const (
	dbdriver = "sqlite3"
	dbsource = "./db/jrdd.db"
)

func initDb() *sql.DB  {
	db,err := sql.Open(dbdriver,dbsource)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func getAllSongs(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	//var artist string
	var songs []Song

	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
								"from Songs INNER JOIN Genres On Songs.genre = Genres.id")
	song := newSong()
	for rows.Next() {
		var genre Genre
		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&genre.Id,&genre.Name,&song.Length)
		if err == nil{
			song.Genre = &genre
			songs = append(songs, *song)
		}
	}

	song.Genre.Name = "DAMMM"

	json.NewEncoder(w).Encode(songs)
	db.Close()
}
func getSongsByArtist(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	var songs []Song

	artist := pat.Param(r, "artist")
	rows, _ := db.Query("Select  Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
								"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
								"where lower(artist) like lower(?)",artist)

	song := newSong()

	for rows.Next(){
		var genre Genre
		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&genre.Id,&genre.Name,&song.Length)
		if err == nil{
			song.Genre = &genre
			songs = append(songs, *song)
		}
	}
	if len(songs)>0 {
		json.NewEncoder(w).Encode(songs)
	}else {
		w.WriteHeader(http.StatusNotFound)
	}
	db.Close()
}


func getSongsByName(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	var songs []Song

	songName := pat.Param(r, "name")
	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
		"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
		"where lower(song) like lower(?)",songName)

	song := newSong()

	for rows.Next(){
		var genre Genre
		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&genre.Id,&genre.Name,&song.Length)
		if err == nil{
			song.Genre = &genre
			songs = append(songs, *song)
		}
	}
	if len(songs)>0 {
		json.NewEncoder(w).Encode(songs)
	}else {
		w.WriteHeader(http.StatusNotFound)
	}
	db.Close()
}

func getSongsByGenre(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	db := initDb()
	var songs []Song
	genre := pat.Param(r, "genre")
	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
		"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
		"where lower(Genres.name) like lower(?)",genre)

	song := newSong()
	for rows.Next(){
		var genre Genre
		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&genre.Id,&genre.Name,&song.Length)
		if err == nil{
			song.Genre = &genre
			songs = append(songs, *song)
		}
	}
	if len(songs)>0 {
		json.NewEncoder(w).Encode(songs)
	}else {
		w.WriteHeader(http.StatusNotFound)
	}
	db.Close()
}