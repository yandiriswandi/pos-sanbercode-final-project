package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
)

func CreateUser(c *gin.Context) {
	var user models.User
	// Validasi input
	if err := c.ShouldBindJSON(&user); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if user.Name == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Name can't empty",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "failed",
			Code:    400,
			Message: validationError,
		})
		return
	}

	query := `INSERT INTO users (code, name, email, username, password, level, address, phone, image, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err := config.DB.QueryRow(query, user.Code, user.Name, user.Email, user.Username, user.Password, user.Level, user.Address, user.Phone, user.Image, user.Status).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Code:    200,
		Status:  "success",
		Message: "success to add data",
		Data:    user,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id") // ambil id dari URL parameter

	// Validasi input
	if err := c.ShouldBindJSON(&user); err != nil {
		validationError := err.Error()

		if user.Name == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Name can't empty",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "failed",
			Code:    400,
			Message: validationError,
		})
		return
	}

	query := `UPDATE users SET code = $1, name = $2, email = $3, username = $4, password = $5, level = $6, address = $7, phone = $8, image = $9, status = $10 WHERE id = $11`

	res, err := config.DB.Exec(query, user.Code, user.Name, user.Email, user.Username, user.Password, user.Level, user.Address, user.Phone, user.Image, user.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Status:  "failed",
			Message: "data not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "succes add data",
		Data:    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM user WHERE id = $1`
	res, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    404,
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Code:    200,
		Message: "success deleted data",
	})
}

func GetUsers(c *gin.Context) {
	// Ambil query parameter page dan size, default kalau kosong
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	// Hitung total data
	var totalData int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM user`).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// Ambil data user
	rows, err := config.DB.Query(`SELECT id, code, name, email, username, password, level, address, phone, image, status FROM users ORDER BY id LIMIT $1 OFFSET $2`, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Code, &user.Name, &user.Email, &user.Username, &user.Password, &user.Level, &user.Address, &user.Phone, &user.Image, &user.Status); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		users = append(users, user)
	}

	// Return the users in a successful response
	c.JSON(http.StatusOK, models.SuccessList{
		Status:    "success",
		Code:      200,
		Message:   "Success get data",
		Data:      users,
		Page:      page,
		Size:      size,
		TotalData: totalData,
	})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	query := `SELECT id, code, name, email, username, password, level, address, phone, image, status 
              FROM users WHERE id = $1`

	err := config.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Code,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Level,
		&user.Address,
		&user.Phone,
		&user.Image,
		&user.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "user not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "Success get data",
		Data:    user,
	})
}
