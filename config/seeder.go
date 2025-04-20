package config

import (
	"fmt"
	"log"

	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
	"github.com/yandiriswandi/pos-sanbercode-final-project/utils"
)

func SeedUsers() {
	users := []models.User{
		{
			Name:     "Admin",
			Code:     "USR/GH/01",
			Email:    "admin@example.com",
			Username: "admin",
			Password: utils.HashPassword("admin1234"),
			Level:    1,
			Address:  "Jakarta",
			Phone:    "081234567890",
			Image:    "default.jpg",
			Status:   1,
		},
		{
			Name:     "user",
			Code:     "USR/GH/02",
			Email:    "user@example.com",
			Username: "user",
			Password: utils.HashPassword("user1234"),
			Level:    2,
			Address:  "Jakarta",
			Phone:    "081234567890",
			Image:    "default.jpg",
			Status:   2,
		},
	}

	for _, user := range users {
		// Cek apakah username sudah ada di database
		var count int
		err := DB.Get(&count, `SELECT COUNT(*) FROM users WHERE username = $1`, user.Username)
		if err != nil {
			fmt.Printf("Error checking username %s: %v", user.Username, err)
			continue
		}

		if count > 0 {
			fmt.Printf("Username %s already exists, skipping seeding.\n", user.Username)
			continue // Skip jika username sudah ada
		}

		// Insert data jika username belum ada
		_, err = DB.Exec(`
            INSERT INTO users (code, name, email, username, password, level, address, phone, image, status) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, user.Code, user.Name, user.Email, user.Username, user.Password, user.Level, user.Address, user.Phone, user.Image, user.Status)

		if err != nil {
			log.Printf("Error seeding user %s: %v", user.Username, err)
		} else {
			fmt.Println("User seeded:", user.Username)
		}
	}
}
