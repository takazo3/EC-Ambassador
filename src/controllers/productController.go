package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func CreateProducts(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	go database.ClearCache("Products_Frontend","Products_Backend")

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product
	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	database.DB.Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	go database.ClearCache("Products_Frontend","Products_Backend")

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.Id = uint(id)

	database.DB.Delete(&product)
	
	go database.ClearCache("Products_Frontend","Products_Backend")

	return nil
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()
	result, err := database.Cache.Get(ctx, "Products_Frontend").Result()

	if err != nil {
		fmt.Println(err.Error())
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.Cache.Set(ctx, "Products_Frontend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		err := json.Unmarshal([]byte(result), &products)
		if err != nil {
			panic(err)
		}
	}

	return c.JSON(products)
}

func ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "Product_Backend").Result()

	if err != nil {
		fmt.Println(err.Error())
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		database.Cache.Set(ctx, "Products_Backend", bytes, 30*time.Minute)
	} else {
		json.Unmarshal([]byte(result), &products)

	}

	var searchedProducts []models.Product

	if s := c.Query("s"); s != "" {

		lower := strings.ToLower(s)
		for _, product := range products {
			if (strings.Contains(strings.ToLower(product.Title), lower)) || (strings.Contains(strings.ToLower(product.Description), lower)) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}
	var total = len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 9

	var data []models.Product = searchedProducts
	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage && page > 0 {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}
	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"Last_page": (total + perPage - 1) / perPage,
	})
}
