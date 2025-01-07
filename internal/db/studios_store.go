package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xV0lk/htmx-go/models"
)

type Psstudiosstore struct {
	db *pgxpool.Pool
}

func NewPsstudiosstore(db *pgxpool.Pool) *Psstudiosstore {
	return &Psstudiosstore{db}
}

func (s *Psstudiosstore) Close() {
	s.db.Close()
}

type Studiostore interface {
	Createstudio(studio *models.Studio, ctx context.Context) error
	Fetchstudio(id int, ctx context.Context) (*models.Studio, error)
	Fetchstudios(studio *models.Studio, ctx context.Context) ([]*models.Studio, error)
	Updatestudio(studio *models.Studio, ctx context.Context) error
	Deletestudio(id int, ctx context.Context) error
}

func (s *Psstudiosstore) Createstudio(studio *models.Studio, ctx context.Context) error {

	query := `INSERT INTO studios (name, address, email, cut)
	VALUES ($1, $2, $3, $4) RETURNING id`
	err := s.db.QueryRow(ctx, query, studio.Name, studio.Address, studio.Email, studio.Cut).Scan(&studio.ID)
	if err != nil {
		return fmt.Errorf("error al crear estudio: %v", err.Error())
	}
	return nil
}

func (s *Psstudiosstore) Fetchstudio(id int, ctx context.Context) (*models.Studio, error) {

	studio := &models.Studio{}

	query := `SELECT * FROM studios WHERE id = $1`
	err := s.db.QueryRow(ctx, query, id).Scan(&studio.ID, &studio.Name, &studio.Address, &studio.Email, &studio.Cut)
	if err != nil {
		return studio, fmt.Errorf("error al buscar estudio: %v", err.Error())
	}
	return studio, nil
}

func (s *Psstudiosstore) Fetchstudios(studio *models.Studio, ctx context.Context) ([]*models.Studio, error) {

	studios := []*models.Studio{}
	var query string
	var param string

	if studio.Name != "" {
		query = `SELECT * FROM studios WHERE name = $1`
		param = studio.Name
	} else if studio.Email != "" {
		query = `SELECT * FROM studios WHERE email = $1`
		param = studio.Email
	} else if studio.Address != "" {
		query = `SELECT * FROM studios WHERE address = $1`
		param = studio.Address
	} else {
		fmt.Println("No hay parametros de busqueda")
		return studios, nil
	}

	rows, err := s.db.Query(ctx, query, param)
	fmt.Println(studio)
	fmt.Println(rows)
	if err != nil {
		fmt.Errorf("error al buscar estudio:%v", err.Error())
		return studios, nil
	}

	defer rows.Close()

	studios, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[models.Studio])
	if err != nil {
		fmt.Errorf("error %v", err.Error())
		return studios, nil
	}

	return studios, nil
}

func (s *Psstudiosstore) Updatestudio(studio *models.Studio, ctx context.Context) error {

	query := `UPDATE studios SET name = $2, address = $3, email = $4, cut = $5 WHERE id = $1`
	_, err := s.db.Exec(ctx, query, studio.ID, studio.Name, studio.Address, studio.Email, studio.Cut)
	if err != nil {
		return fmt.Errorf("error actualizando estudio: %v", err.Error())
	}
	return nil
}

func (s *Psstudiosstore) Deletestudio(id int, ctx context.Context) error {

	query := `DELETE FROM studios WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando estudio: %v", err.Error())
	}
	return nil
}
