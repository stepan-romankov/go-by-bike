package models

import (
	"github.com/go-pg/pg/v10"
)

type Bike struct {
	ID     uint32  `json:"id"`
	Name   string  `json:"name" pg:",notnull"`
	Lat    float32 `json:"latitude" pg:",notnull"`
	Lon    float32 `json:"longitude" pg:",notnull"`
	Rented bool    `pg:"-"`
}

func GetBikes(db *pg.DB) ([]Bike, error) {
	var bikes []Bike

	err := db.Model(&bikes).
		ColumnExpr("bike.*").
		ColumnExpr("COALESCE(r.bike_id, 0) > 0 AS rented").
		Join("LEFT JOIN rentals AS r").
		JoinOn("r.bike_id = bike.id").
		JoinOn("completed = ?", false).
		Order("bike.id ASC").
		Select()

	return bikes, err
}
