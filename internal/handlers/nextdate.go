package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/subbbbbaru/first-sample/pkg/log"

	"github.com/subbbbbaru/go_final_project/internal/nextdate"
)

func (h *Handler) NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.URL.Query().Get("now")
	taskParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		log.Error().Println(err)
		http.Error(w, fmt.Sprintf(`invalid "now" parameter: %s`, err.Error()), http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, taskParam, repeatParam)
	if err != nil {
		log.Error().Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Error().Println(err)
	}
}
