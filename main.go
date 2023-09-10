package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	log.Println("test traffic", n)

	quit := make(chan bool)

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-quit:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(10 * time.Second)
	for i := 0; i < n; i++ {
		quit <- true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(true)
}

func main() {
	flagService := NewFlagService()
	router := mux.NewRouter()

	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	router.HandleFunc("/create", flagService.CreateFlag).Methods("POST")
	router.HandleFunc("/update", flagService.UpdateFlag).Methods("POST")
	router.HandleFunc("/get", flagService.GetUserFlagValue).Methods("GET")
	router.HandleFunc("/flags", flagService.ListFlags).Methods("GET")
	router.HandleFunc("/test", TestHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/", fs)

	port := 8080
	fmt.Printf("FF Server started on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), corsOptions(router))
}
