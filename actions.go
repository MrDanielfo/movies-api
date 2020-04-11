package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

/* Propiedades */

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	return session
}

func responseMovie(w http.ResponseWriter, status int, results Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

func responseMovies(w http.ResponseWriter, status int, results []Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

// var movies = Movies{
// 	Movie{"Kill Bill", 2003, "Quentin Tarantino"},
// 	Movie{"Batman Begins", 2005, "Christopher Nolan"},
// 	Movie{"Let the right one in", 2008, "Thomas Alfredson"},
// }

var collection = getSession().DB("movies-go").C("movies")

/* Métodos */

// IndexHandler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hola Mundo desde mi servidor web con Go</h1>")
}

// MoviesHandler
func MoviesHandler(w http.ResponseWriter, r *http.Request) {

	var results []Movie
	err := collection.Find(nil).All(&results)
	if err != nil {
		log.Fatal(err)
	}
	responseMovies(w, 200, results)
	/* fmt.Fprintf(w, "<h1>Hola Mundo desde la página de MovieList</h1>") */
}

// MoviesSingle

func MoviesSingle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	oID := bson.ObjectIdHex(movieID)

	results := Movie{}
	err := collection.FindId(oID).One(&results)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	responseMovie(w, 201, results)
	/* fmt.Fprintf(w, "<h1>Has cargado la película número %s</h1>", movieID) */

}

// func MovieAdd

func MovieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var movieData Movie
	err := decoder.Decode(&movieData)

	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	/* log.Println(movieData) */
	/* Conexión a MongoDB */
	errResponse := collection.Insert(movieData)
	if errResponse != nil {
		w.WriteHeader(500)
		return
	}

	responseMovie(w, 201, movieData)
	// movies = append(movies, movieData)
}

/* Update Movie */

func MovieUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	oID := bson.ObjectIdHex(movieID)
	/* decoder de lo que llega por json */
	decoder := json.NewDecoder(r.Body)
	/* variable movieData */
	var movieData Movie
	errDecoder := decoder.Decode(&movieData)
	if errDecoder != nil {
		panic(errDecoder)
	}

	defer r.Body.Close()
	document := bson.M{"_id": oID}
	change := bson.M{"$set": movieData}
	err := collection.Update(document, change)

	if err != nil {
		w.WriteHeader(404)
		return
	}
	responseMovie(w, 201, movieData)
	/* fmt.Fprintf(w, "<h1>Has cargado la película número %s</h1>", movieID) */
}

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (this *Message) setStatus(data string) {
	this.Status = data
}

func (this *Message) setMessage(data string) {
	this.Message = data
}

/* Delete Movie */

func MovieRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]
	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}
	oID := bson.ObjectIdHex(movieID)
	err := collection.RemoveId(oID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	/* results := Message{
		"success",
		"La película con ID " + movieID + " ha sido borrada",
	} */

	message := new(Message)
	message.setStatus("success")
	message.setMessage("La película con ID " + movieID + " ha sido borrada")

	results := message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)

	/* fmt.Fprintf(w, "<h1>Has cargado la película número %s</h1>", movieID) */
}
