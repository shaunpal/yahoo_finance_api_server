package yahoo_finance_api_server

import (
	"net/http"
	"os"

	"./handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	router.HandleFunc("/index/{symbol}", handlers.GetIndex).Methods("GET")
	router.HandleFunc("/equity/{symbol}", handlers.GetEquity).Methods("GET")

	port := ":" + os.Getenv("PORT")
	http.ListenAndServe(port, router)
}
