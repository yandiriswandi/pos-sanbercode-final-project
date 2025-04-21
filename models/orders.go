package models

import "time"

type Order struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"` // relasi ke tabel users
	CoID            int       `json:"cart_id"`
	Code            string    `db:"code" json:"code"`   // kode unik order (opsional tapi bagus punya)
	Price           float64   `db:"price" json:"price"` // harga satuan saat dimasukkan
	Subtotal        float64   `db:"subtotal" json:"subtotal"`
	Quantity        float64   `db:"quantity" json:"quantity"` // total belanja
	ProductID       int       `db:"product_id" json:"product_id"`
	Status          string    `db:"status" json:"status"` // status pesanan (pending, paid, shipped, etc.)
	UserName        string    `json:"user_name"`
	Address         string    `db:"address" json:"address"`
	ProductName     string    `json:"product_name"` // quantity * price
	Note            string    `db:"note" json:"note"`
	TransactionDate time.Time `db:"transaction_date" json:"transaction_date"` // catatan tambahan (opsional)
}
