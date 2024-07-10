package storage

import (
	"fmt"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func (s *Storage) GetOne() {
	stmt, err := s.Db.Query("select name, country from cities")
	if err != nil {
		panic(err)
	}
	var c models.City
	for stmt.Next() {
		if err = stmt.Scan(&c.Name, &c.Country); err != nil {
			panic(err)
		}
	}
	fmt.Println(c)
}
