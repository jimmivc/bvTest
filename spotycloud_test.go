package main

import (
	"goji.io"
	"goji.io/pat"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFindSongByArtist(t *testing.T)  {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/artist/:artist"), getSongsByArtist)

	statusAndResponseGetTest(mux,"/api/songs/artist/The Black Eyed Peas",`{"artist":"The Black Eyed Peas","song":"I Gotta Feeling","genre":{"name":"Rap"},"length":219}`,t)
}

func TestFindSongByName(t *testing.T)  {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/:name"), getSongsByName)

	statusAndResponseGetTest(mux,"/api/songs/Macarena",`{"artist":"Los Del Rio","song":"Macarena","genre":{"name":"Pop"},"length":159}`,t)
}

func TestFindSongByGenre(t *testing.T)  {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/genre/:genre"), getSongsByGenre)

	statusAndResponseGetTest(mux,"/api/songs/genre/latin pop rock",`{"artist":"Los Waldners","song":"Horacio","genre":{"name":"Latin Pop Rock"},"length":165}`,t)
}

func TestListGenresSummary(t *testing.T)  {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/genre/summary"), getAllGenresSummary)

	statusAndResponseGetTest(mux,"/api/genre/summary","",t)
}


func TestFindSongByLength(t *testing.T)  {

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/byLength/:min/:max"), getSongsByLength)
	statusAndResponseGetTest(mux,"/api/songs/byLength/180/190",`{"artist":"LMFAO","song":"Party Rock Anthem","genre":{"name":"Rap"},"length":189}`,t)

}

func statusAndResponseGetTest(mux *goji.Mux,path string, expected string,t *testing.T){
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", path,nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mux.ServeHTTP(rr,req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(),expected) {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			rr.Body.String(), expected)
	}
}