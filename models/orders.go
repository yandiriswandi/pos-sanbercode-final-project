package models

type Order struct {
	ID              int     `db:"id" json:"id"`
	UserID          int     `db:"user_id" json:"user_id"` // relasi ke tabel users
	Code            string  `db:"code" json:"code"`       // kode unik order (opsional tapi bagus punya)
	Price           float64 `db:"price" json:"price"`     // harga satuan saat dimasukkan
	Subtotal        float64 `db:"sub_total" json:"sub_total"`
	Quantity        float64 `db:"quantity" json:"quantity"` // total belanja
	ProductID       int     `db:"product_id" json:"product_id"`
	Status          string  `db:"status" json:"status"` // status pesanan (pending, paid, shipped, etc.)
	Note            string  `db:"note" json:"note"`
	TransactionDate string  `db:"transaction_date" json:"transaction_date"` // catatan tambahan (opsional)
}
