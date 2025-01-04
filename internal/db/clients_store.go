package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xV0lk/htmx-go/models"
)

type ClientStore interface {
	CreateClient(ctx context.Context, client *models.Client) (*models.Client, error)
	FetchClient(id int, ctx context.Context) (*models.Client, error)
	FetchClients(key *models.Client, ctx context.Context) ([]*models.Client, error)
	UpdateClient(clientInfo *models.Client) (*models.Client, error)
	DeleteClient(id int) error
}

type PsclientStore struct {
	db *pgxpool.Pool
}

func NewPsclientStore(db *pgxpool.Pool) *PsclientStore {

	return &PsclientStore{
		db,
	}
}

func (s *PsclientStore) Close() {
	s.db.Close()
}

func (s *PsclientStore) CreateClient(ctx context.Context, client *models.Client) (*models.Client, error) {

	query := `INSERT INTO clients (name, phone, email, notifications, studio_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING *`

	err := s.db.QueryRow(ctx, query, client.Name, client.Phone, client.Email, client.Notifications, client.StudioID).Scan(&client.ID, &client.Name, &client.Phone, &client.Email, &client.Notifications, &client.StudioID)
	if err != nil {
		return client, fmt.Errorf("Inserting Error %w", err)
	}
	return client, nil
}

func (s *PsclientStore) FetchClient(id int, ctx context.Context) (*models.Client, error) {

	client := &models.Client{}

	query := `SELECT * FROM clients WHERE id = $1`

	err := s.db.QueryRow(ctx, query, id).Scan(&client.ID, &client.Name, &client.Phone, &client.Email, &client.Notifications, &client.StudioID)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return client, nil
}

func (s *PsclientStore) FetchClients(key *models.Client, ctx context.Context) ([]*models.Client, error) {

	clients := []*models.Client{}

	var query string
	var param string
	if key.Name != "" {
		query = `SELECT * FROM clients WHERE name = $1`
		param = key.Name
	} else if key.Email != "" {
		query = `SELECT * FROM clients WHERE email = $1`
		param = key.Email
	} else if key.Phone != "" {
		query = `SELECT * FROM clients WHERE phone = $1`
		param = key.Phone
	} else {
		return nil, fmt.Errorf("there is no value to search with")
	}

	// applying the query

	rows, err := s.db.Query(ctx, query, param)
	if err != nil {
		fmt.Errorf("Error in: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	// Collecting the data info clients

	clients, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[models.Client])
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return clients, nil

}

// UpdateTaskTitle updates the title of a task in the SQLStore.
//
// It takes an ID (int) and a title (string) as parameters and returns an Item and an error.

func (s *PsclientStore) UpdateClient(clientInfo *models.Client) (*models.Client, error) {

	client := &models.Client{}

	query := `UPDATE clients SET name = $2, phone = $3, email = $4, notifications = $5, studio_id = $6 WHERE id = $1 RETURNING *`

	fmt.Println(clientInfo.ID)

	row := s.db.QueryRow(context.Background(), query, clientInfo.ID, clientInfo.Name, clientInfo.Phone, clientInfo.Email, clientInfo.Notifications, clientInfo.StudioID)
	err := row.Scan(&client.ID, &client.Name, &client.Phone, &client.Email, &client.Notifications, &client.StudioID)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return client, nil
}

func (s *PsclientStore) DeleteClient(id int) error {

	query := `DELETE FROM clients WHERE Id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return nil

}
