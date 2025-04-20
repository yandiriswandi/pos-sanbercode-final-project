package models

type Cart struct {
	ID          int     `db:"id" json:"id"`
	UserID      int     `db:"user_id" json:"user_id"`       // user pemilik item
	ProductID   int     `db:"product_id" json:"product_id"` // produk yang dimasukkan
	Quantity    int     `db:"quantity" json:"quantity"`     // jumlah produk
	Price       float64 `db:"price" json:"price"`           // harga satuan saat dimasukkan
	Subtotal    float64 `db:"subtotal" json:"subtotal"`     // quantity * price
	Note        string  `db:"note" json:"note"`
	UserName    string  `json:"user_name"`
	ProductName string  `json:"product_name"` // quantity * price
}
