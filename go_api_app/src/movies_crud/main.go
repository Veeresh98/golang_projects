/*
author @veereshm

*/

package main

import (
	"encoding/json" //to encode my code into json when i send it to postman
	"fmt"           // to get the logs and errors
	"math/rand"
	"net/http" // to create a hhtp
	"strconv"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

// note i'm not going to use aa database while i'll be using a structs and methods

//creating a struct called movie which will have movie id, title director etc

type Movie struct {
	ID       string    `json:"id"`
	isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// each movie will have separate director

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

//creating a getmovies fucntion and passing w and r as responsewriter and request
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	log.Infof("get the w")
	json.NewEncoder(w).Encode(movies)
	//a function defined in the encoding/json package which gets a JSON encoding of any type and encodes/writes
	//it any writable stream that implements a io. Writer interface.
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	parameters := mux.Vars(r)  //
	for index, item := range movies {
		if item.ID == parameters["id"] {
			movies = append(movies[index:], movies[index+1:]...)
			log.Infof("see the movie list ", movies)
			break
		}
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	parameters := mux.Vars(r)
	for _, item := range movies {
		if item.ID == parameters["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body) // we are going tosend the movie to the body of the postman api
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	log.Infof("see the movie ID ", movie.ID)
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set the content type
	w.Header().Set("content-type", "application/json")
	//params
	parameters := mux.Vars(r)
	// itirate through the 	movies
	for index, item := range movies {
		for item.ID == parameters["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			log.Infof("see the movie list ", movies)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = parameters["id"]
			log.Infof("see the movie ID ", movie.ID)
			movies = append(movies, movie)
			log.Infof("see the movie list ", movies)
			json.NewEncoder(w).Encode(&movie)
			test := json.NewEncoder(w).Encode(&movie)
			log.Infof("see the output of the json", test)
			return
		}
	}

}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", isbn: "23214", Title: "One piece movie 1", Director: &Director{Firstname: "Enricha", Lastname: "Oda"}})
	movies = append(movies, Movie{ID: "2", isbn: "23235", Title: "One piece movie 2", Director: &Director{Firstname: "Enricha", Lastname: "Oda"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POSt")
	r.HandleFunc("/movies", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movies", updateMovie).Methods("PUT")

	fmt.Println("Let's start, Starting of my crud APi on my localhost")
	log.Fatal(http.ListenAndServe(":8000", r))

}
