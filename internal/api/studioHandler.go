package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/xV0lk/htmx-go/internal/db"
	"github.com/xV0lk/htmx-go/models"
)

type Studiohandler struct {
	studiostore db.Studiostore
	decoder     *schema.Decoder
}

func Newstudiohandler(studiostore db.Studiostore, decoder *schema.Decoder) Studiohandler {

	return Studiohandler{
		studiostore: studiostore,
		decoder:     decoder,
	}
}

func writeshortJSONResponse(w http.ResponseWriter, id int, s string) {

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(id)
	w.Write([]byte(s))
}

func writeJSONresponse(w http.ResponseWriter, id int, s *models.Studio) {

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(id)
	json.NewEncoder(w).Encode(s)
}

func (s *Studiohandler) Createnewstudio(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	studio := &models.Studio{}

	err := json.NewDecoder(r.Body).Decode(studio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.Validatestudio(studio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.studiostore.Createstudio(studio, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeshortJSONResponse(w, http.StatusOK, "Estudio creado satisfactoriamente")

}

func (s *Studiohandler) Fetchstudio(w http.ResponseWriter, r *http.Request) {

	studio := &models.Studio{}
	ctx := r.Context()

	sid := chi.URLParam(r, "id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	studio, err = s.studiostore.Fetchstudio(id, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONresponse(w, http.StatusOK, studio)
}

func (s *Studiohandler) Fetchstudios(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	studio := &models.Studio{}
	studios := []*models.Studio{}

	if r.URL.Query().Get("email") != "" {
		studio.Email = r.URL.Query().Get("email")
		if !strings.ContainsAny(studio.Email, "@") {
			http.Error(w, "ingrese un email valido", http.StatusBadRequest)
			return
		}
	} else if r.URL.Query().Get("address") != "" {
		studio.Address = r.URL.Query().Get("address")
	} else {
		http.Error(w, "no existen datos para realizar busqueda", http.StatusBadRequest)
		return
	}

	studios, err := s.studiostore.Fetchstudios(studio, ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if len(studios) == 0 {
		http.Error(w, "estudio no encontrado", http.StatusBadRequest)
		return

	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(studios)

}

func (s *Studiohandler) Updatestudio(w http.ResponseWriter, r *http.Request) {

	studio := &models.Studio{}
	ctx := r.Context()
	sid := chi.URLParam(r, "id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(studio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	studio.ID = id

	err = s.studiostore.Updatestudio(studio, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONresponse(w, http.StatusOK, studio)
}

func (s *Studiohandler) Deletestudio(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	sid := chi.URLParam(r, "id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.studiostore.Deletestudio(id, ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeshortJSONResponse(w, http.StatusOK, "Estudio eliminado correctamente")

}
