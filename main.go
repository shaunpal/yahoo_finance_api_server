package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"server/handlers"
	"os"
)

func main(){
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello"))
	})

	router.HandleFunc("/index/{symbol}", handlers.GetIndex).Methods("GET")
	router.HandleFunc("/equity/{symbol}", handlers.GetEquity).Methods("GET")


	port := ":" + os.Getenv("PORT")
	http.ListenAndServe(port, router)
}