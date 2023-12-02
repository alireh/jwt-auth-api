package main

import (
	"fmt"
	"jwt-auth-api/app"
	"jwt-auth-api/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
    router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/user/getall", controllers.GetUsers).Methods("GET")

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, // Allowing only get, just an example
	})

    //http.ListenAndServe(":8000", corsOpts.Handler(router))

	fmt.Println(port)

	err :=  http.ListenAndServe(":8000", c.Handler(router))
	// err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
