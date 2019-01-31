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
	"strconv"
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
	fmt.Fprintln(w,"http://localhost:8005/api/genres/summary")
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

func getSongsByArtist(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	artist := pat.Param(r, "artist")
	rows, _ := db.Query("Select  Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
								"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
								"where lower(artist) like lower('%'||?||'%')",artist)


	songs := mapSongsByQueryResult(rows)

	if len(songs)>0 {
		json.NewEncoder(w).Encode(songs)
	}else {
		fmt.Fprint(w,"Not Found")
	}
	db.Close()
}


func getSongsByName(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	songName := pat.Param(r, "name")
	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
		"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
		"where lower(song) like lower('%'||?||'%')",songName)


	searchSongsListResponse(mapSongsByQueryResult(rows),w)
	db.Close()
}

func getSongsByGenre(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	db := initDb()
	genre := pat.Param(r, "genre")
	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
		"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
		"where lower(Genres.name) like lower('%'||?||'%')",genre)

	searchSongsListResponse(mapSongsByQueryResult(rows),w)

	db.Close()
}

func getAllGenresSummary(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	db := initDb()

	var genres []struct{
		Name string
		CantSongs int
		Length int
	}

	rows, _ := db.Query("select a.name,COUNT(b.id) as songs,SUM(b.length) as length from Genres as a join Songs as b on a.id = b.genre group by a.name")
	for rows.Next() {
		var summary struct{
			Name string
			CantSongs int
			Length int
		}
		err := rows.Scan(&summary.Name,&summary.CantSongs,&summary.Length)
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

	min := pat.Param(r,"min")
	max := pat.Param(r,"max")

	minimus,errmin := strconv.Atoi(min)
	maximus,errmax := strconv.Atoi(max);

	if errmin != nil || errmax !=nil{
		fmt.Fprint(w,"Input format error")
		return
	}

	if maximus-minimus<0{
		fmt.Fprint(w,"Invalid Range (min/max)")
		return
	}


	rows, _ := db.Query("Select Songs.id,Songs.artist,Songs.song,Genres.id,Genres.name,Songs.length " +
								"from Songs INNER JOIN Genres On Songs.genre = Genres.id " +
								"where length>=? and length<=?",min,max)

	searchSongsListResponse(mapSongsByQueryResult(rows),w)

	db.Close()
}

func mapSongsByQueryResult(rows *sql.Rows) []Song  {
	var songs []Song
	for rows.Next() {
		song := newSong()
		err := rows.Scan(&song.Id,&song.Artist,&song.Song,&song.Genre.Id,&song.Genre.Name,&song.Length)
		if err == nil{
			songs = append(songs, *song)
		}
	}

	return songs
}

func searchSongsListResponse(songs []Song, w http.ResponseWriter){
	if len(songs)>0 {
		json.NewEncoder(w).Encode(songs)
	}else {
		w.WriteHeader(http.StatusNotFound)
	}
}