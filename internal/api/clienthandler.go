package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/xV0lk/htmx-go/internal/db"
	"github.com/xV0lk/htmx-go/models"
)

type ClientHandler struct {
	ClientStore db.ClientStore
	Decoder     *schema.Decoder
}

func Newclienthandler(store db.ClientStore, decoder *schema.Decoder) *ClientHandler {

	return &ClientHandler{
		ClientStore: store,
		Decoder:     decoder,
	}
}

func writeJSONResponse(w http.ResponseWriter, status int, client *models.Client) {

	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(client)

}

func validateClient(client *models.Client) error {

	if client.Email == "" {

		return fmt.Errorf("el campo no puede estar vacio")

	} else if !strings.ContainsAny(client.Email, "@") {

		return fmt.Errorf("email no valido")
	}

	if client.Name == "" {

		return fmt.Errorf("el campo no puede estar vacio")

	} else if strings.IndexFunc(client.Name, unicode.IsDigit) != -1 || strings.IndexFunc(client.Name, unicode.IsSymbol) != -1 {

		return fmt.Errorf("el campo no puede tener numeros o caracteres especiales")

	}

	if client.Phone == "" {

		return fmt.Errorf("el campo no puede estar vacio")

	} else if strings.IndexFunc(client.Phone, unicode.IsLetter) != -1 || strings.IndexFunc(client.Phone, unicode.IsSymbol) != -1 {

		return fmt.Errorf("el campo no puede tener letras o caracteres especiales")

	}

	return nil
}

func validateClients(client *models.Client) error {

	if !strings.ContainsAny(client.Email, "@") {

		return fmt.Errorf("email no valido")
	}

	if strings.IndexFunc(client.Name, unicode.IsDigit) != -1 || strings.IndexFunc(client.Name, unicode.IsSymbol) != -1 {

		return fmt.Errorf("el campo no puede tener numeros o caracteres especiales")

	}

	if strings.IndexFunc(client.Phone, unicode.IsLetter) != -1 || strings.IndexFunc(client.Phone, unicode.IsSymbol) != -1 {

		return fmt.Errorf("el campo no puede tener letras o caracteres especiales")

	}

	return nil
}

func (c *ClientHandler) CreateNewClient(w http.ResponseWriter, r *http.Request) {

	newClient := &models.Client{}

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(newClient)

	if err != nil {

		http.Error(w, "Data ingresada no es valida", http.StatusBadRequest)

		return

	}

	err = validateClient(newClient)

	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)

		return

	}

	newClient, err = c.ClientStore.CreateClient(ctx, newClient)

	if err != nil {

		http.Error(w, "Error al crear cliente", http.StatusBadRequest)

	}

	writeJSONResponse(w, http.StatusOK, newClient)

	fmt.Println("Cliente fue creado correctamente")

}

func (c *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {

	cid := chi.URLParam(r, "id")

	id, err := strconv.Atoi(cid)

	if err != nil {

		http.Error(w, "id no valido", http.StatusBadRequest)

	}

	client := &models.Client{}

	ctx := r.Context()

	client, err = c.ClientStore.FetchClient(id, ctx)

	if err != nil {

		http.Error(w, "Error al buscar cliente", http.StatusNotFound)
	}

	writeJSONResponse(w, http.StatusOK, client)

}

func (c *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {

	client := &models.Client{}
	clients := []*models.Client{}

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(client)

	if err != nil {

		http.Error(w, "Data ingresada no es valida", http.StatusBadRequest)

	}

	err = validateClients(client)

	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)

	}

	clients, err = c.ClientStore.FetchClients(client, ctx)

	if err != nil {

		http.Error(w, "Error al buscar cliente", http.StatusNotFound)

	}

	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(clients)
}

func (c *ClientHandler) UpdateClients(w http.ResponseWriter, r *http.Request) {

	clientInfo := &models.Client{}

	err := json.NewDecoder(r.Body).Decode(clientInfo)

	if err != nil {

		http.Error(w, "Data ingresada no es valida", http.StatusBadRequest)

	}

	err = validateClient(clientInfo)

	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)

	}

	clientInfo, err = c.ClientStore.UpdateClient(clientInfo)

	if err != nil {

		http.Error(w, "Error al actualizar cliente", http.StatusBadRequest)
	}

	writeJSONResponse(w, http.StatusOK, clientInfo)
}

func (c *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {

	cid := chi.URLParam(r, "id")

	id, err := strconv.Atoi(cid)

	if err != nil {

		http.Error(w, "id no valido", http.StatusBadRequest)

	}

	err = c.ClientStore.DeleteClient(id)

	if err != nil {

		http.Error(w, "error eliminando cliente", http.StatusBadRequest)

	}

}
