package main

import (
	"github.com/go-faker/faker/v4"
	"ambassador/src/database"
	"ambassador/src/models"
	"math/rand"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		product := models.Product{
			Title: faker.Username(),
			Description:  faker.Username(),
			Image: faker.URL(),
			Price: float64(rand.Intn(90) + 10),
		}

	
		database.DB.Create(&product)
	}
}
