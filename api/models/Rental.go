package models

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"log"
)

type Rental struct {
	Id        uint32 `pg:",pk"`
	BikeId    uint32 `pg:",notnull"`
	Bike      *Bike  `pg:"rel:has-one"`
	UserID    uint32 `pg:",notnull"`
	User      *User  `pg:"rel:has-one"`
	Completed bool   `pg:"default:false,notnull"`
}

type RentalError struct {
	Err error
}

func (r *RentalError) Error() string {
	return r.Err.Error()
}

func Rent(db *pg.DB, userId uint32, bikeId uint32) (*Rental, error) {
	var rentalId uint32
	_, err := db.
		WithParam("userID", userId).
		WithParam("bikeID", bikeId).
		QueryOne(pg.Scan(&rentalId), "INSERT INTO rentals(bike_id, user_id, completed) SELECT ?bikeID, ?userID, false WHERE "+
			"NOT EXISTS (SELECT 1 FROM rentals WHERE (bike_id = ?bikeID OR user_id = ?userID) AND completed = false) "+
			"RETURNING id")

	if err != nil {
		if err == pg.ErrNoRows {
			log.Printf("Failed to rent bike %d for user %d", userId, bikeId)
			return nil, &RentalError{Err: errors.New("can't rent bike")}
		}
		log.Printf("Unexpected error: %s", err)
		return nil, err
	}

	log.Printf("Bike %d rented for user %d ", bikeId, userId)
	return &Rental{Id: rentalId, UserID: userId, BikeId: bikeId}, nil
}

func Return(db *pg.DB, userId uint32) error {
	result, err := db.
		WithParam("userID", userId).
		Exec("UPDATE rentals SET completed = true WHERE user_id = ?userID AND completed = false")

	if err != nil {
		log.Printf("Unexpected error %s", err)
		return err
	}

	if result.RowsAffected() == 0 {
		log.Printf("Failed to return bike for user %d", userId)
		return &RentalError{Err: errors.New("can't return bike")}
	}
	log.Printf("Bike for user %d returned", userId)
	return nil
}
