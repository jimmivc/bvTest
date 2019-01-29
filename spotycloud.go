package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

type Song struct {
	Id 		int `json:"-"`
	Artist 	string `json:"artist"`
	Song 	string `json:"song"`
	Genre 	*Genre `json:"genre"`
	Length 	int `json:"length"`
}

type Genre struct{
	Id		int `json:"-"`
	Name	string `json:"name"`
}

func newSong() *Song{
	return &Song{
		Genre: new(Genre),
	}
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), homepage)

	//mux.HandleFunc(pat.Get("/api/songs"), getAllSongs)
	//mux.HandleFunc(pat.Get("/api/genres"), getAllGenres)
	mux.HandleFunc(pat.Get("/api/songs/artist/:artist"), getSongsByArtist)
	mux.HandleFunc(pat.Get("/api/songs/:name"), getSongsByName)
	mux.HandleFunc(pat.Get("/api/songs/genre/:genre"), getSongsByGenre)
	mux.HandleFunc(pat.Get("/api/genres/summary"), getAllGenresSummary)
	mux.HandleFunc(pat.Get("/api/songs/byLength/:min/:max"), getSongsByLength)



	log.Fatal(http.ListenAndServe("localhost:8005", mux))
}

func homepage(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"http://localhost:8005/api/songs/artist/{artistName}")
	fmt.Fprintln(w,"http://localhost:8005/api/songs/{songName}")
	fmt.Fprintln(w,"http://localhost:8005/api/songs/genre/{genreName}")
	fmt.Fprintln(w,"http://localhost:8005/api/genres/summary/{genreName}")
	fmt.Fprintln(w,"http://localhost:8005/api/songs/byLength/{min}/{max}")
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
	for rows.Next() {
		song := newSong()

		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Genre.Name,&song.Length)
		if err == nil{
			songs = append(songs, *song)
		}
	}

	json.NewEncoder(w).Encode(songs)
	db.Close()
}

func getAllGenres(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	//var artist string
	var genres []Genre

	rows, _ := db.Query("select id, name from Genres")
	for rows.Next() {
		var genre = new(Genre)
		err := rows.Scan(&genre.Id,&genre.Name)
		if err == nil{
			genres = append(genres, *genre)
		}
	}

	json.NewEncoder(w).Encode(genres)
	db.Close()
}

func getSongsByArtist(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	var songs []Song

	artist := pat.Param(r, "artist")
	rows, _ := db.Query("Select  Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
								"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
								"where lower(artist) like lower('%'||?||'%')",artist)


	for rows.Next(){
		song := newSong()

		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Genre.Name,&song.Length)
		if err == nil{
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
		"where lower(song) like lower('%'||?||'%')",songName)


	for rows.Next(){
		song := newSong()

		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Genre.Name,&song.Length)
		if err == nil{
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
		"where lower(Genres.name) like lower('%'||?||'%')",genre)

	for rows.Next(){
		song := newSong()

		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Genre.Name,&song.Length)
		if err == nil{
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

func getAllGenresSummary(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	var genres []struct{
		Id int
		Name string
		CantSongs int
		Length int
	}

	rows, _ := db.Query("select a.id,a.name,COUNT(b.id) as songs,SUM(b.length) as length from Genres as a join Songs as b on a.id = b.genre group by a.id")
	for rows.Next() {
		var summary struct{
			Id int
			Name string
			CantSongs int
			Length int
		}
		err := rows.Scan(&summary.Id,&summary.Name,&summary.CantSongs,&summary.Length)
		if err == nil{
			genres = append(genres,summary)
		}
	}

	json.NewEncoder(w).Encode(genres)
	db.Close()
}

func getSongsByLength(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	var songs []Song

	min := pat.Param(r,"min")
	max := pat.Param(r,"max")

	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
								"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
								"where length>=? and length<=?",min,max)
	for rows.Next() {
		song := newSong()

		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Genre.Name,&song.Length)
		if err == nil{
			songs = append(songs, *song)
		}
	}

	json.NewEncoder(w).Encode(songs)
	db.Close()
}