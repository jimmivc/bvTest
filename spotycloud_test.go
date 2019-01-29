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
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.

	req, err := http.NewRequest("GET", "/api/songs/artist/The Black Eyed Peas",nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/artist/:artist"), getSongsByArtist)
	mux.ServeHTTP(rr,req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"artist":"The Black Eyed Peas","song":"I Gotta Feeling","genre":{"name":"Rap"},"length":219}`

	if !strings.Contains(rr.Body.String(),expected) {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			rr.Body.String(), expected)
	}
}


func TestFindSongByName(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/songs/Macarena",nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/:name"), getSongsByName)
	mux.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"artist":"Los Del Rio","song":"Macarena","genre":{"name":"Pop"},"length":159}`

	if !strings.Contains(rr.Body.String(),expected) {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			rr.Body.String(), expected)
	}
}

func TestFindSongByGenre(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/songs/genre/latin pop rock",nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/genre/:genre"), getSongsByGenre)
	mux.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"artist":"Los Waldners","song":"Horacio","genre":{"name":"Latin Pop Rock"},"length":165}`

	if !strings.Contains(rr.Body.String(),expected) {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			rr.Body.String(), expected)
	}
}

func TestListGenresSummary(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/genres/summary",nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/genres/summary"), getAllGenresSummary)
	mux.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}


func TestFindSongByLength(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/songs/byLength/180/190",nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/songs/byLength/:min/:max"), getSongsByLength)
	mux.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"artist":"LMFAO","song":"Party Rock Anthem","genre":{"name":"Rap"},"length":189}`

	if !strings.Contains(rr.Body.String(),expected) {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			rr.Body.String(), expected)
	}
}