package models

import (
	"errors"
	"strings"
	"unicode"
)

type Studio struct {
	ID      int    `db:"id,omitempty"`
	Name    string `db:"name"`
	Address string `db:"address"`
	Email   string `db:"email"`
	Cut     int    `db:"cut"`
}

type NewStudio struct {
	Name    string
	Address string
	Email   string
	Cut     int
}

func Validatestudio(s *Studio) error {

	if s.Name == "" {
		return errors.New("el campo no puede estar vacio")
	} else if strings.IndexFunc(s.Name, unicode.IsNumber) != -1 || strings.IndexFunc(s.Name, unicode.IsSymbol) != -1 {
		return errors.New("el nombre no puede tener numeros o caracteres especiales")
	} else if s.Address == "" {
		return errors.New("el campo no puede estar vacio")
	} else if s.Email == "" {
		return errors.New("el campo no puede estar vacio")
	} else if !strings.ContainsAny(s.Email, "@") {
		return errors.New("ingrese un email valido")
	} else {

		return nil
	}

}
