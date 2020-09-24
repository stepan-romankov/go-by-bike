package db

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"log"
)

var dbModels = []interface{}{
	(*models.User)(nil),
	(*models.Bike)(nil),
	(*models.Rental)(nil),
}

func Migrate(postgresUrl string, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, postgresUrl)
	if err != nil {
		log.Panicf("Failed to initialize migrations: %s", err)
	}

	err = m.Up()
	if err != nil {
		log.Panicf("Failed to apply migrations: %s", err)
	}

	return nil
}

func Cleanup(db *pg.DB) error {
	for i := len(dbModels) - 1; i >= 0; i-- {
		model := dbModels[i]
		err := db.Model(model).DropTable(&orm.DropTableOptions{IfExists: true, Cascade: true})
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func Fixtures(db *pg.DB) {
	count, err := db.Model(&models.Bike{}).Count()
	if err != nil {
		panic(err)
	}

	if count == 0 {
		db.Model(&models.Bike{Name: "Henry", Lat: 50.119504, Lon: 8.638137, Rented: false}).Insert()
		db.Model(&models.Bike{Name: "Hans", Lat: 50.119229, Lon: 8.640020, Rented: false}).Insert()
		db.Model(&models.Bike{Name: "Thomas", Lat: 50.120452, Lon: 8.650507, Rented: false}).Insert()
	}

	count, err = db.Model(&models.User{}).Count()
	if count == 0 {
		hash, err := models.HashPassword("test")
		if err != nil {
			panic(err)
		}
		db.Model(&models.User{Login: "user1", Password: hash}).Insert()
		db.Model(&models.User{Login: "user2", Password: hash}).Insert()
		db.Model(&models.User{Login: "user3", Password: hash}).Insert()
	}
}
