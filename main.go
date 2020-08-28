package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"server/handlers"
)

func main(){
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello"))
	})

	router.HandleFunc("/index/{symbol}", handlers.GetIndex).Methods("GET")
	router.HandleFunc("/equity/{symbol}", handlers.GetEquity).Methods("GET")

	server := &http.Server{
		Addr: "http://yahoo-finance-api-server.herokuapp.com:8080",
		Handler: router,
	}
	server.ListenAndServe()
}