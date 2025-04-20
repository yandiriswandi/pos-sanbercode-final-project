package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Code     string `db:"code" json:"code"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Level    int    `db:"level" json:"level"`
	Address  string `db:"address" json:"address"`
	Phone    string `db:"phone" json:"phone"`
	Image    string `db:"image" json:"image"`
	Status   int    `db:"status" json:"status"`
}

type Login struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}
