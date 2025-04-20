package models

type Product struct {
	ID           int    `db:"id" json:"id"`
	Code         string `db:"code" json:"code"`
	Name         string `db:"name" json:"name" binding:"required"`
	Image        string `db:"image" json:"image"`
	Stock        int    `db:"stock" json:"stock"`
	CategoryID   int    `db:"category_id" json:"category_id" binding:"required"`
	Description  string `db:"description" json:"description"`
	CategoryName string `db:"category_name" json:"category_name"`
}

type Category struct {
	ID          int    `db:"id" json:"id"`
	Code        string `db:"code" json:"code"`
	Name        string `db:"name" json:"name" binding:"required"`
	Description string `db:"description" json:"description"`
}
