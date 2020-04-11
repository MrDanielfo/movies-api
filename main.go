package main

import (
	"log"
	"net/http"
)

func main() {
	router := newRouter()
	server := http.ListenAndServe(":8080", router)
	log.Fatal(server)
	/* fmt.Println("El servidor está corriendo en el puerto :8080") */
}

/* String de conexión local Mongo

# mongodb://localhost:27017/barger-recipes

*/
