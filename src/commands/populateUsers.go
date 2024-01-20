package main

import (
	"github.com/go-faker/faker/v4"
	"ambassador/src/database"
	"ambassador/src/models"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		ambassador := models.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Email:     faker.Email(),
			IsAmbassador: true,
		}

		ambassador.SetPassword("1234")

		database.DB.Create(&ambassador)
	}
}