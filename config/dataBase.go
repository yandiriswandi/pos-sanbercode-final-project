package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB  *sqlx.DB
	err error
)

func InitDB() {

	err = godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	var err error
	// Ambil value dari environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// dbHost := os.Getenv("PGHOST")
	// dbPort := os.Getenv("PGPORT")
	// dbUser := os.Getenv("PGUSER")
	// dbPassword := os.Getenv("PGPASSWORD")
	// dbName := os.Getenv("PGDATABASE")

	// Bangun connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS category (
		id SERIAL PRIMARY KEY,
		code VARCHAR(100) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS users (
    	id SERIAL PRIMARY KEY,
    	code VARCHAR(100) NOT NULL UNIQUE,
    	name VARCHAR(100) NOT NULL,
    	email VARCHAR(100) NOT NULL UNIQUE,
    	username VARCHAR(50) NOT NULL UNIQUE,
    	password VARCHAR(255) NOT NULL,
    	level INT NOT NULL,
    	address TEXT,
    	phone VARCHAR(20),
    	image VARCHAR(255),
    	status INT DEFAULT 1
	);
	
	CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		code VARCHAR(100) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		image TEXT,
		stock INT DEFAULT 0,
		category_id INT NOT NULL,
		description TEXT,
    
    	CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS cart (
    	id SERIAL PRIMARY KEY,
    	user_id INTEGER NOT NULL,
    	product_id INTEGER NOT NULL,
    	quantity INTEGER NOT NULL,
    	price NUMERIC(12, 2) NOT NULL,
    	subtotal NUMERIC(12, 2) NOT NULL,
    	note TEXT,
   
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES product(id)
);


	`

	DB.MustExec(schema)
	fmt.Println("Database siap digunakan 🚀")
	SeedUsers()
}
