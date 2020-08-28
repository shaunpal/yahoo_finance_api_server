package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"server/handlers"
)

func main(){
	router := mux.NewRouter()
	router.Host("https://yahoo-finance-api-server.herokuapp.com/")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello"))
	})

	router.HandleFunc("/index/{symbol}", handlers.GetIndex).Methods("GET")
	router.HandleFunc("/equity/{symbol}", handlers.GetEquity).Methods("GET")

	
	http.ListenAndServe(":8080", router)
}