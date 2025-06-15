package web

import (
	"encoding/json"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/gorilla/mux"
)

type apiHypha struct {
	Name    string `json:"name"`
	Text    string `json:"text,omitempty"`
	IsMedia bool   `json:"is_media"`
}

type uploadHyphaReq struct {
	Text string `json:"text"`
}

func initAPI(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/hypha", handlerAPIListHyphae).Methods(http.MethodGet)
	api.HandleFunc("/hypha/{name}", handlerAPIHypha).Methods(http.MethodGet, http.MethodPost)
}

func handlerAPIListHyphae(w http.ResponseWriter, rq *http.Request) {
	var names []string
	for h := range hyphae.YieldExistingHyphae() {
		names = append(names, h.CanonicalName())
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(names)
}

func handlerAPIHypha(w http.ResponseWriter, rq *http.Request) {
	name := util.CanonicalName(mux.Vars(rq)["name"])
	h := hyphae.ByName(name)

	if rq.Method == http.MethodPost {
		var data uploadHyphaReq
		if err := json.NewDecoder(rq.Body).Decode(&data); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		u := user.FromRequest(rq)
		if err := shroom.UploadText(h, []byte(data.Text), "", u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case *hyphae.TextualHypha:
		text, err := hyphae.FetchMycomarkupFile(h)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondHypha(w, h.CanonicalName(), text, false)
	case *hyphae.MediaHypha:
		text := ""
		if h.HasTextFile() {
			var err error
			text, err = hyphae.FetchMycomarkupFile(h)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		respondHypha(w, h.CanonicalName(), text, true)
	}
}

func respondHypha(w http.ResponseWriter, name, text string, isMedia bool) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(apiHypha{
		Name:    name,
		Text:    text,
		IsMedia: isMedia,
	})
}
