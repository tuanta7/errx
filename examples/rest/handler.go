package main

import (
	"encoding/json"
	"net/http"
	"time"

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
		httpCode, message := errx.Global.GetHTTPResponse(err, language)
		http.Error(w, message, httpCode)
		return
	}

	jsonCounter, _ := json.Marshal(counter)
	_, _ = w.Write(jsonCounter)
}

func (h *Handler) ResetCounter(w http.ResponseWriter, r *http.Request) {
	counterName := r.PathValue("id")
	language := r.URL.Query().Get("language")
	if language == "" {
		language = lang.English.String()
	}

	err := h.uc.SetCounter(counterName, &Counter{
		Value:       1,
		LastUpdated: time.Now().Unix(),
	})
	if err != nil {
		httpCode, message := errx.Global.GetHTTPResponse(err, language)
		http.Error(w, message, httpCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
