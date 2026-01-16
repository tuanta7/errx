package main

import (
	"encoding/json"
	"net/http"

	"github.com/tuanta7/errx"
	lang "golang.org/x/text/language"
)

type Handler struct {
	uc *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetCounter(w http.ResponseWriter, r *http.Request) {
	counterName := r.PathValue("id")
	language := r.URL.Query().Get("language")
	if language == "" {
		language = lang.English.String()
	}

	counter, err := h.uc.GetCounter(counterName)
	if err != nil {
		httpCode, message := errx.HTTPResponse(err, language)
		http.Error(w, message, httpCode)
		return
	}

	jsonCounter, _ := json.Marshal(counter)
	_, _ = w.Write(jsonCounter)
}
