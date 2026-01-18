package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tuanta7/errx"
	"golang.org/x/text/language"
)

func main() {
	errx.Global.RegisterMessage(ErrCounterNotFound, language.English.String(), "Counter not found")
	errx.Global.RegisterMessage(ErrCounterNotFound, language.Vietnamese.String(), "Không tìm thấy bộ đếm")
	errx.Global.RegisterHTTPErrorCode(ErrCounterNotFound, http.StatusNotFound)

	cache := NewCache()
	repo := NewRepository(cache)
	uc := NewUseCase(repo)
	handler := NewHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /counters/{id}", handler.GetCounter)
	mux.HandleFunc("PUT /counters/{id}", handler.ResetCounter)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server started at port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
