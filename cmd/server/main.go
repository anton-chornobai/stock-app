package main

import (
	"fmt"
	"net/http"
	"todo-app/internal/handlers"
)

func main() {
	fmt.Println("Started the server")
	router := http.NewServeMux()
	router.HandleFunc("/todos", handler.GetStocks)
	router.HandleFunc("POST /todos", handler.PostStock)
	router.HandleFunc("GET /stocks", handler.GetSingleStock)

	s := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
